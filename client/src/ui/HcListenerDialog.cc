#include <Havoc.h>
#include <ui/HcListenerDialog.h>

#include <QtCore/QVariant>
#include <QtWidgets/QGridLayout>
#include <QtWidgets/QPushButton>
#include <QtWidgets/QWidget>

HcListenerDialog::HcListenerDialog(
    const QString& editing
) {
    if ( objectName().isEmpty() ) {
        setObjectName( QString::fromUtf8( "HcListenerDialog" ) );
    }

    gridLayout = new QGridLayout( this );
    gridLayout->setObjectName( "gridLayout" );

    InputName = new HcLineEdit( this );
    InputName->Input->setObjectName( QString::fromUtf8( "InputName" ) );

    ComboProtocol = new HcComboBox;
    ComboProtocol->Combo->setObjectName( QString::fromUtf8( "ComboProtocol" ) );

    StackedProtocols = new QStackedWidget( this );
    StackedProtocols->setObjectName( QString::fromUtf8( "HcListenerDialog.StackedProtocols" ) );

    ButtonSave = new QPushButton( this );
    ButtonSave->setObjectName( "ButtonClose" );
    ButtonSave->setProperty( "HcButton", "true" );

    ButtonClose = new QPushButton( this );
    ButtonClose->setObjectName( "ButtonClose" );
    ButtonClose->setProperty( "HcButton", "true" );

    gridLayout->addWidget( InputName,        0, 0, 1, 6 );
    gridLayout->addWidget( ComboProtocol,    1, 0, 1, 6 );
    gridLayout->addWidget( StackedProtocols, 2, 0, 1, 6 );

    horizontal[ 0 ] = new QSpacerItem( 40, 20, QSizePolicy::Expanding, QSizePolicy::Minimum );
    horizontal[ 1 ] = new QSpacerItem( 40, 20, QSizePolicy::Expanding, QSizePolicy::Minimum );
    horizontal[ 2 ] = new QSpacerItem( 40, 20, QSizePolicy::Expanding, QSizePolicy::Minimum );
    horizontal[ 3 ] = new QSpacerItem( 40, 20, QSizePolicy::Expanding, QSizePolicy::Minimum );

    gridLayout->addWidget( ButtonSave,  3, 2, 1, 1 );
    gridLayout->addWidget( ButtonClose, 3, 3, 1, 1 );

    gridLayout->addItem( horizontal[ 0 ], 3, 0, 1, 1 );
    gridLayout->addItem( horizontal[ 2 ], 3, 1, 1, 1 );
    gridLayout->addItem( horizontal[ 3 ], 3, 4, 1, 1 );
    gridLayout->addItem( horizontal[ 1 ], 3, 5, 1, 1 );

    retranslateUi();

    QMetaObject::connectSlotsByName( this );

    connect( ComboProtocol->Combo, &QComboBox::currentTextChanged, this, &HcListenerDialog::changeProtocol );
    connect( ButtonSave,  &QPushButton::clicked, this, &HcListenerDialog::save );
    connect( ButtonClose, &QPushButton::clicked, this, [&]() {
        State = Closed;
        close();
    } );

    retranslateUi();

    for ( auto& name : Havoc->Protocols() ) {
        if ( Havoc->ProtocolObject( name ).has_value() ) {
            if ( editing.isEmpty() ) {
                AddProtocol( name, Havoc->ProtocolObject( name ).value() );
            } else if ( editing.toStdString() == name ) {
                AddProtocol( name, Havoc->ProtocolObject( name ).value() );
                break;
            }
        }
    }

    QMetaObject::connectSlotsByName( this );
}

auto HcListenerDialog::retranslateUi() -> void {
    setStyleSheet( Havoc->StyleSheet() );
    setWindowTitle( "Listener" );

    ComboProtocol->Combo->clear();

    InputName->setLabelText( "Name:      " );
    ComboProtocol->setLabelText( "Protocol:  " );
    ButtonSave->setText( "Save" );
    ButtonClose->setText( "Close" );
}

auto HcListenerDialog::changeProtocol(
    const QString &text
) -> void {
    for ( int i = 0; i < Protocols.size(); i++ ) {
        if ( Protocols[ i ].name == text.toStdString() ) {
            StackedProtocols->setCurrentIndex( i );
            break;
        }
    }
}

auto HcListenerDialog::getCloseState() -> ListenerState { return State; }

auto HcListenerDialog::start() -> void {
    if ( Havoc->Protocols().empty() ) {
        Helper::MessageBox(
            QMessageBox::Icon::Warning,
            "Listener error",
            "No protocols available"
        );
        return;
    }

    ComboProtocol->Combo->setCurrentIndex( 0 );
    exec();
}

auto HcListenerDialog::save() -> void {
    auto result    = httplib::Result();
    auto data      = json();
    auto found     = false;
    auto exception = std::string();

    if ( InputName->text().isEmpty() ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "Listener failure",
            "Failed to start listener: name is empty"
        );
    }

    auto protocol = getCurrentProtocol();
    try {
        HcPythonAcquire();

        /* sanity check input */
        if ( ! protocol.instance.attr( "sanity_check" )().cast<bool>() ) {
            /* sanity check failed. exit and dont send request */
            spdlog::debug( "sanity check failed. exit and dont send request" );
            return;
        }

        data = {
            { "name",     InputName->text().toStdString() },
            { "protocol", protocol.name },
            { "config",   protocol.instance.attr( "save" )() }
        };
    } catch ( py11::error_already_set &eas ) {
        Helper::MessageBox(
            QMessageBox::Icon::Critical,
            "Listener saving error",
            std::string( eas.what() )
        );
        return;
    }

    State = Error;

    if ( edited_config.empty() ) {
        result = Havoc->ApiSend( "/api/listener/start", data );
    } else {
        result = Havoc->ApiSend( "/api/listener/edit", data );
    }

    if ( result->status != 200 ) {
        if ( result->body.empty() ) {
            goto ERROR_SERVER_RESPONSE;
        }

        try {
            if ( ( data = json::parse( result->body ) ).is_discarded() ) {
                goto ERROR_SERVER_RESPONSE;
            }
        } catch ( std::exception& e ) {
            spdlog::error( "failed to parse json response: \n{}", e.what() );
            goto ERROR_SERVER_RESPONSE;
        }

        if ( ! data.contains( "error" ) ) {
            goto ERROR_SERVER_RESPONSE;
        }

        if ( ! data[ "error" ].is_string() ) {
            goto ERROR_SERVER_RESPONSE;
        }

        Helper::MessageBox(
            QMessageBox::Critical,
            "Listener failure",
            QString( "failed to start listener: %1" ).arg( data[ "error" ].get<std::string>().c_str() ).toStdString()
        );

        State = Error;
    } else {
        State = Saved;
    }

    close();
    return;

ERROR_SERVER_RESPONSE:
    Helper::MessageBox(
        QMessageBox::Critical,
        "Listener failure",
        "listener failure: Invalid response from the server"
    );

    close();
}

auto HcListenerDialog::getCurrentProtocol() -> Protocol { return Protocols[ StackedProtocols->currentIndex() ]; }

auto HcListenerDialog::AddProtocol(
    const std::string&  name,
    const py11::object& object
) -> void {
    HcPythonAcquire();

    auto protocol = Protocol {
        .name   = name,
        .widget = new QWidget
    };

    try {
        protocol.instance = object( U_PTR( protocol.widget ), name );
        protocol.instance.attr( "_hc_main" )();
    } catch ( py11::error_already_set &eas ) {
        Helper::MessageBox(
            QMessageBox::Icon::Critical,
            "Listener loading error",
            std::string( eas.what() )
        );
        return;
    }

    ComboProtocol->Combo->addItem( protocol.name.c_str() );
    StackedProtocols->addWidget( protocol.widget );

    Protocols.push_back( protocol );
}

auto HcListenerDialog::setEditingListener(
    const QString& name,
    const json&    config
) -> void {

    InputName->setInputText( name );
    ComboProtocol->Combo->setDisabled( true );
    InputName->Input->setDisabled( true );

    edited_config = config;

    try {
        HcPythonAcquire();

        getCurrentProtocol().instance.attr( "edit" )( config );
    } catch ( py11::error_already_set& e ) {
        spdlog::debug( "failed to execute edit method of listener: \n{}", e.what() );
        Helper::MessageBox(
            QMessageBox::Icon::Critical,
            "Listener editing error",
            std::string( e.what() )
        );
    }
}
