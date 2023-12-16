#include <Havoc.h>
#include <ui/HcConsole.h>

HcConsole::HcConsole(
    QWidget* parent
) : QWidget( parent ) {
    if ( objectName().isEmpty() ) {
        setObjectName( "HcConsole" );
    }

    gridLayout = new QGridLayout( this );
    gridLayout->setObjectName( "gridLayout" );
    gridLayout->setContentsMargins( 0, 0, 0, 0 );

    LabelHeader = new QLabel( this );
    LabelHeader->setObjectName( "LabelHeader" );

    Console = new QTextEdit( this );
    Console->setObjectName( "Console" );
    Console->setReadOnly( true );
    Console->setProperty( "HcConsole", "true" );

    LabelBottom = new QLabel( this );
    LabelBottom->setObjectName( "LabelBottom" );

    LabelInput = new QLabel( this );
    LabelInput->setObjectName( "LabelInput" );

    Input = new QLineEdit( this );
    Input->setObjectName( "Input" );

    QObject::connect( Input, &QLineEdit::returnPressed, this, [&] () {
        auto text = py11::str( Input->text().toStdString() );

        Input->setText( "" );

        for ( auto& callback : Havoc->MainWindows->PageScripts->PyConsole->InputCallback ) {
            callback( text );
        }
    } );

    gridLayout->addWidget( LabelHeader, 0, 0, 1, 2 );
    gridLayout->addWidget( Console,     1, 0, 1, 2 );
    gridLayout->addWidget( LabelBottom, 2, 0, 1, 2 );
    gridLayout->addWidget( LabelInput,  3, 0, 1, 1 );
    gridLayout->addWidget( Input,       3, 1, 1, 1 );

    QMetaObject::connectSlotsByName( this );
}

auto HcConsole::setHeaderLabel(
    const QString& text
) -> void {
    LabelHeader->setText( text );
}

auto HcConsole::setBottomLabel(
    const QString& text
) -> void {
    LabelBottom->setText( text );
}

auto HcConsole::setInputLabel(
    const QString& text
) -> void {
    LabelInput->setText( text );
}

auto HcConsole::appendConsole(
    const QString& text
) -> void {
    Console->append( text );
}
