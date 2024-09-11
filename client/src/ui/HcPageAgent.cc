#include <Havoc.h>
#include <ui/HcPageAgent.h>
#include <QtConcurrent/QtConcurrent>

HcPageAgent::HcPageAgent(
    QWidget* parent
) : QWidget( parent ) {

    if ( objectName().isEmpty() ) {
        setObjectName( QString::fromUtf8( "PageAgent" ) );
    }

    gridLayout = new QGridLayout( this );
    gridLayout->setObjectName( "gridLayout" );

    ComboAgentView = new QComboBox( this );
    ComboAgentView->setObjectName( "ComboAgentView" );
    ComboAgentView->view()->setProperty( "ComboBox", "true" );

    Splitter = new QSplitter( this );
    Splitter->setObjectName( QString::fromUtf8( "Splitter" ) );
    Splitter->setOrientation( Qt::Vertical );

    StackedWidget = new QStackedWidget( Splitter );
    StackedWidget->setContentsMargins( 0, 0, 0, 0 );

    AgentTable = new QTableWidget( StackedWidget );
    AgentTable->setObjectName( "AgentTable" );

    AgentGraph = new HcSessionGraph( StackedWidget );
    AgentGraph->setObjectName( "AgentGraph" );

    TitleAgentID      = new QTableWidgetItem( "UUID" );
    TitleInternal     = new QTableWidgetItem( "Internal" );
    TitleUsername     = new QTableWidgetItem( "User" );
    TitleHostname     = new QTableWidgetItem( "Host" );
    TitleSystem       = new QTableWidgetItem( "System" );
    TitleProcessID    = new QTableWidgetItem( "Pid" );
    TitleProcessName  = new QTableWidgetItem( "Process" );
    TitleProcessArch  = new QTableWidgetItem( "Arch" );
    TitleThreadID     = new QTableWidgetItem( "Tid" );
    TitleNote         = new QTableWidgetItem( "Note" );
    TitleLastCallback = new QTableWidgetItem( "Last" );

    if ( AgentTable->columnCount() < 11 ) {
        AgentTable->setColumnCount( 11 );
    }

    /* TODO: get how we should add this from the settings
     * for now we just do a default one */
    AgentTable->setHorizontalHeaderItem( 0,  TitleAgentID      );
    AgentTable->setHorizontalHeaderItem( 1,  TitleInternal     );
    AgentTable->setHorizontalHeaderItem( 2,  TitleUsername     );
    AgentTable->setHorizontalHeaderItem( 3,  TitleHostname     );
    AgentTable->setHorizontalHeaderItem( 4,  TitleProcessName  );
    AgentTable->setHorizontalHeaderItem( 5,  TitleProcessID    );
    AgentTable->setHorizontalHeaderItem( 6,  TitleThreadID     );
    AgentTable->setHorizontalHeaderItem( 7,  TitleProcessArch  );
    AgentTable->setHorizontalHeaderItem( 8,  TitleSystem       );
    AgentTable->setHorizontalHeaderItem( 9,  TitleLastCallback );
    AgentTable->setHorizontalHeaderItem( 10, TitleNote         );

    /* table settings */
    AgentTable->setEnabled( true );
    AgentTable->setShowGrid( false );
    AgentTable->setSortingEnabled( false );
    AgentTable->setWordWrap( true );
    AgentTable->setCornerButtonEnabled( true );
    AgentTable->horizontalHeader()->setVisible( true );
    AgentTable->setSelectionBehavior( QAbstractItemView::SelectRows );
    AgentTable->setContextMenuPolicy( Qt::CustomContextMenu );
    AgentTable->horizontalHeader()->setSectionResizeMode( QHeaderView::ResizeMode::Stretch );
    AgentTable->horizontalHeader()->setStretchLastSection( true );
    AgentTable->verticalHeader()->setVisible( false );
    AgentTable->setFocusPolicy( Qt::NoFocus );
    AgentTable->setAlternatingRowColors( true );

    AgentTab = new QTabWidget( Splitter );
    AgentTab->setObjectName( "AgentTab" );
    AgentTab->setTabsClosable( true );
    AgentTab->setMovable( true );
    AgentTab->tabBar()->setProperty( "HcTab", "true" );

    AgentDisplayerSessions = new QLabel( this );
    AgentDisplayerSessions->setObjectName( "LabelDisplaySessions" );
    AgentDisplayerSessions->setProperty( "HcLabelDisplay", "true" );

    AgentDisplayerTargets = new QLabel( this );
    AgentDisplayerTargets->setObjectName( "AgentDisplayTargets" );
    AgentDisplayerTargets->setProperty( "HcLabelDisplay", "true" );

    AgentDisplayerPivots = new QLabel( this );
    AgentDisplayerPivots->setObjectName( "AgentDisplayPivots" );
    AgentDisplayerPivots->setProperty( "HcLabelDisplay", "true" );

    AgentDisplayerElevated = new QLabel( this );
    AgentDisplayerElevated->setObjectName( "AgentDisplayElevated" );
    AgentDisplayerElevated->setProperty( "HcLabelDisplay", "true" );

    horizontalSpacer = new QSpacerItem( 40, 20, QSizePolicy::Expanding, QSizePolicy::Minimum );
    AgentActionMenu  = new QMenu( this );

    ActionPayload    = AgentActionMenu->addAction( QIcon( ":/icons/16px-payload-build" ), "Payload Builder" );
    ActionShowHidden = AgentActionMenu->addAction( QIcon( ":/icons/32px-eye-white" ), "Show Hidden" );
    ActionShowHidden->setCheckable( true );

    AgentActionButton = new QToolButton( this );
    AgentActionButton->setObjectName( "AgentActionButton" );
    AgentActionButton->setToolButtonStyle( Qt::ToolButtonTextBesideIcon );
    AgentActionButton->setText( "Actions" );
    AgentActionButton->setProperty( "HcButton", "true" );
    AgentActionButton->setMenu( AgentActionMenu );
    AgentActionButton->setIcon( QIcon( ":/icons/32px-flash" ) );
    AgentActionButton->setPopupMode( QToolButton::InstantPopup );
    AgentActionButton->setLayoutDirection( Qt::LeftToRight );
    AgentActionButton->setMaximumWidth( 90 );

    StackedWidget->addWidget( AgentTable );
    StackedWidget->addWidget( AgentGraph );

    Splitter->addWidget( StackedWidget );
    Splitter->addWidget( AgentTab );
    Splitter->handle( 1 )->setEnabled( SplitterMoveToggle ); /* disabled by default */

    connect( AgentTable, &QTableWidget::customContextMenuRequested, this, &HcPageAgent::handleAgentMenu );
    connect( AgentTable, &QTableWidget::doubleClicked, this, &HcPageAgent::handleAgentDoubleClick );
    connect( AgentTab, &QTabWidget::tabCloseRequested, this, &HcPageAgent::tabCloseRequested );
    connect( ActionPayload, &QAction::triggered, this, &HcPageAgent::actionPayloadBuilder );
    connect( ActionShowHidden, &QAction::triggered, this, &HcPageAgent::actionShowHidden );
    connect( AgentActionButton, &QToolButton::triggered, this, &HcPageAgent::actionTriggered );
    connect( AgentTable, &QTableWidget::itemChanged, this, &HcPageAgent::itemChanged );
    connect( ComboAgentView, &QComboBox::currentIndexChanged, this, &HcPageAgent::viewChanged );

    gridLayout->addWidget( ComboAgentView,         0, 0, 1, 1 );
    gridLayout->addWidget( Splitter,               1, 0, 1, 7 );
    gridLayout->addWidget( AgentDisplayerSessions, 0, 2, 1, 1 );
    gridLayout->addWidget( AgentDisplayerTargets,  0, 1, 1, 1 );
    gridLayout->addWidget( AgentDisplayerPivots,   0, 3, 1, 1 );
    gridLayout->addWidget( AgentDisplayerElevated, 0, 4, 1, 1 );
    gridLayout->addItem( horizontalSpacer,         0, 5, 1, 1 );
    gridLayout->addWidget( AgentActionButton,      0, 6, 1, 1 );

    retranslateUi( );

    QMetaObject::connectSlotsByName( this );
}

HcPageAgent::~HcPageAgent() = default;

auto HcPageAgent::retranslateUi() -> void {
    setStyleSheet( Havoc->StyleSheet() );

    AgentDisplayerElevated->setText( "Elevated: 0" );
    AgentDisplayerSessions->setText( "Sessions: 0" );
    AgentDisplayerTargets->setText( "Targets: 0" );
    AgentDisplayerPivots->setText( "Pivots: 0" );
    ComboAgentView->addItems( QStringList() << "Sessions" << "Sessions Graph" );
}

auto HcPageAgent::addTab(
    const QString& name,
    QWidget*       widget
) const -> void {
    if ( AgentTab->count() == 0 ) {
        Splitter->setSizes( QList<int>() << 200 << 220 );
        Splitter->handle( 1 )->setEnabled( true );
        Splitter->handle( 1 )->setCursor( Qt::SplitVCursor );
    }

    AgentTab->setCurrentIndex( AgentTab->addTab( widget, QIcon( ":/icons/32px-agent-console" ), name ) );
}

HcAgentTableItem::HcAgentTableItem(
    const QString&          value,
    const Qt::ItemFlag      flags,
    const Qt::AlignmentFlag align
)  : QTableWidgetItem() {
    setText( value );
    setTextAlignment( align );
    setFlags( this->flags() ^ flags );
}

auto HcPageAgent::addAgent(
    const json& metadata
) -> void {
    auto uuid    = QString();
    auto type    = std::string();
    auto arch    = QString();
    auto user    = QString();
    auto host    = QString();
    auto local   = QString();
    auto path    = QString();
    auto process = QString();
    auto pid     = QString();
    auto tid     = QString();
    auto system  = QString();
    auto last    = QString();
    auto note    = QString();
    auto meta    = json();
    auto row     = AgentTable->rowCount();
    auto sort    = AgentTable->isSortingEnabled();

    if ( metadata.contains( "uuid" ) && metadata[ "uuid" ].is_string() ) {
        uuid = QString( metadata[ "uuid" ].get<std::string>().c_str() );
    } else {
        spdlog::error( "[HcPageAgent::addAgent] agent does not contain valid uuid" );
        return;
    }

    if ( metadata.contains( "type" ) && metadata[ "type" ].is_string() ) {
        type = metadata[ "type" ].get<std::string>();
    } else {
        spdlog::error( "[HcPageAgent::addAgent] agent does not contain valid type" );
        return;
    }

    if ( metadata.contains( "note" ) && metadata[ "note" ].is_string() ) {
        note = QString( metadata[ "note" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain any note" );
    }

    if ( metadata.contains( "meta" ) && metadata[ "meta" ].is_object() ) {
        meta = metadata[ "meta" ].get<json>();
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta object" );
        return;
    }

    if ( meta.contains( "user" ) && meta[ "user" ].is_string() ) {
        user = QString( meta[ "user" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta user" );
    }

    if ( meta.contains( "host" ) && meta[ "host" ].is_string() ) {
        host = QString( meta[ "host" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta host" );
    }

    if ( meta.contains( "arch" ) && meta[ "arch" ].is_string() ) {
        arch = QString( meta[ "arch" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta arch" );
    }

    if ( meta.contains( "local ip" ) && meta[ "local ip" ].is_string() ) {
        local = QString( meta[ "local ip" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta local ip" );
    }

    if ( meta.contains( "process path" ) && meta[ "process path" ].is_string() ) {
        path = QString( meta[ "process path" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta process path" );
    }

    if ( meta.contains( "process name" ) && meta[ "process name" ].is_string() ) {
        process = QString( meta[ "process name" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta process name" );
    }

    if ( meta.contains( "pid" ) && meta[ "pid" ].is_number_integer() ) {
        pid = QString::number( meta[ "pid" ].get<int>() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta pid" );
    }

    if ( meta.contains( "tid" ) && meta[ "tid" ].is_number_integer() ) {
        tid = QString::number( meta[ "tid" ].get<int>() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta tid" );
    }

    if ( meta.contains( "system" ) && meta[ "system" ].is_string() ) {
        system = QString( meta[ "system" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta system" );
    }

    if ( meta.contains( "last callback" ) && meta[ "last callback" ].is_string() ) {
        last = QString( meta[ "last callback" ].get<std::string>().c_str() );
    } else {
        spdlog::debug( "[HcPageAgent::addAgent] agent does not contain valid meta last" );
    }

    auto agent = new HcAgent {
        .uuid = uuid.toStdString(),
        .type = type,
        .data = metadata,
        .last = last,
        .ui   = {
            .Uuid        = new HcAgentTableItem( uuid ),
            .Internal    = new HcAgentTableItem( local ),
            .Username    = new HcAgentTableItem( user ),
            .Hostname    = new HcAgentTableItem( host ),
            .ProcessPath = new HcAgentTableItem( path ),
            .ProcessName = new HcAgentTableItem( process ),
            .ProcessId   = new HcAgentTableItem( pid ),
            .ThreadId    = new HcAgentTableItem( tid ),
            .Arch        = new HcAgentTableItem( arch ),
            .System      = new HcAgentTableItem( system ),
            .Note        = new HcAgentTableItem( note, Qt::NoItemFlags, Qt::AlignVCenter ),
            .Last        = new HcAgentTableItem( last ),
        }
    };

    agent->console = new HcAgentConsole( agent );
    agent->console->setBottomLabel( QString( "[User: %1] [Process: %2] [Pid: %3] [Tid: %4]" ).arg( user ).arg( path ).arg( pid ).arg( tid ) );
    agent->console->setInputLabel( ">>>" );
    agent->console->LabelHeader->setFixedHeight( 0 );

    agent->ui.Uuid->agent        = agent;
    agent->ui.Internal->agent    = agent;
    agent->ui.Username->agent    = agent;
    agent->ui.Hostname->agent    = agent;
    agent->ui.ProcessPath->agent = agent;
    agent->ui.ProcessName->agent = agent;
    agent->ui.ProcessId->agent   = agent;
    agent->ui.ThreadId->agent    = agent;
    agent->ui.Arch->agent        = agent;
    agent->ui.System->agent      = agent;
    agent->ui.Note->agent        = agent;
    agent->ui.Last->agent        = agent;

    //
    // connect signals and slots
    //
    connect( & agent->emitter, & HcAgentEmit::ConsoleWrite, this, []( const QString& uuid, const QString& text ) {
        if ( const auto agent = Havoc->Agent( uuid.toStdString() ); agent.has_value() ) {
            agent.value()->console->appendConsole( HcConsole::formatString( text.toStdString() ).c_str() );
        }
    } );

    agents.push_back( agent );

    AgentTable->setRowCount( row + 1 );
    AgentTable->setSortingEnabled( false );
    AgentTable->setItem( row, 0,  agent->ui.Uuid        );
    AgentTable->setItem( row, 1,  agent->ui.Internal    );
    AgentTable->setItem( row, 2,  agent->ui.Username    );
    AgentTable->setItem( row, 3,  agent->ui.Hostname    );
    AgentTable->setItem( row, 4,  agent->ui.ProcessName );
    AgentTable->setItem( row, 5,  agent->ui.ProcessId   );
    AgentTable->setItem( row, 6,  agent->ui.ThreadId    );
    AgentTable->setItem( row, 7,  agent->ui.Arch        );
    AgentTable->setItem( row, 8,  agent->ui.System      );
    AgentTable->setItem( row, 9,  agent->ui.Last        );
    AgentTable->setItem( row, 10, agent->ui.Note        );
    AgentTable->setSortingEnabled( sort );

    //
    // if an interface has been registered then assign it to the agent
    //
    agent->interface = std::nullopt;
    if ( auto interface = Havoc->AgentObject( agent->type ) ) {
        if ( interface.has_value() ) {
            HcPythonAcquire();

            try {
                agent->interface = interface.value()( agent->uuid, agent->type, metadata[ "meta" ] );
            } catch ( py11::error_already_set &eas ) {
                spdlog::error( "failed to invoke agent interface [uuid: {}] [type: {}]: {}", agent->uuid, agent->type, eas.what() );
            }
        }
    }

    agent->node = AgentGraph->addAgent( agent );

    AgentDisplayerTargets->setText( QString( "Targets: %1" ).arg( agents.size() ) );   /* TODO: all targets (only show one host)        */
    AgentDisplayerSessions->setText( QString( "Sessions: %1" ).arg( agents.size() ) ); /* TODO: only set current alive beacons/sessions */
    AgentDisplayerPivots->setText( "Pivots: 0" );
    AgentDisplayerElevated->setText( "Elevated: 0" );
}

auto HcPageAgent::handleAgentMenu(
    const QPoint& pos
) -> void {
    auto menu          = QMenu( this );
    auto uuid          = std::string();
    auto selections    = AgentTable->selectionModel()->selectedRows();
    auto type          = std::string();
    auto agent_actions = Havoc->Actions( HavocClient::ActionObject::ActionAgent );

    /* check if we point to a session table item/agent */
    if ( ! AgentTable->itemAt( pos ) ) {
        return;
    }

    if ( selections.count() > 1 ) {
        //
        // for now we dont allow custom actions for multiple sessions
        // as we have to take the following things into consideration:
        //  - agent type
        //  - add action to the menu while other agent types are selected as well
        //  - check every selected row for agent type and add only the selected action etc.
        //
        // TODO: while this is surely possible i am going
        //       to move this into the future to implement.
        //
        menu.addAction( QIcon( ":/icons/16px-agent-console" ), "Interact" );
        menu.addSeparator();
        menu.addAction( QIcon( ":/icons/16px-blind-white" ), "Hide" );
        menu.addAction( QIcon( ":/icons/16px-remove" ), "Remove" );
    } else {
        //
        // if a single selected agent item then try
        // to add the registered actions as well
        //
        menu.addAction( QIcon( ":/icons/16px-agent-console" ), "Interact" );
        menu.addSeparator();

        //
        // add the registered agent type actions
        //
        uuid = AgentTable->item( selections.at( 0 ).row(), 0 )->text().toStdString();

        for ( auto& agent : agents ) {
            if ( agent->uuid == uuid ) {
                type = agent->type;
                break;
            }
        }

        for ( auto action : agent_actions ) {
            if ( action->agent.type == type ) {
                if ( action->icon.empty() ) {
                    menu.addAction( action->name.c_str() );
                } else {
                    menu.addAction( QIcon( action->icon.c_str() ), action->name.c_str() );
                }
            }
        }

        menu.addSeparator();
        menu.addAction( QIcon( ":/icons/16px-blind-white" ), "Hide" );
        menu.addAction( QIcon( ":/icons/16px-remove" ), "Remove" );
    }

    if ( auto action = menu.exec( AgentTable->horizontalHeader()->viewport()->mapToGlobal( pos ) ) ) {
        for ( const auto& selected : selections ) {
            uuid = AgentTable->item( selected.row(), 0 )->text().toStdString();

            if ( action->text().compare( "Interact" ) == 0 ) {
                spawnAgentConsole( uuid );
            } else if ( action->text().compare( "Remove" ) == 0 ) {
                auto agent = Agent( uuid ).value();

                agent->remove();
            } else if ( action->text().compare( "Hide" ) == 0 ) {
                auto agent = Agent( uuid ).value();

                agent->hide();
            } else {
                for ( auto agent_action : agent_actions ) {
                    if ( agent_action->name       == action->text().toStdString() &&
                         agent_action->agent.type == type
                    ) {
                        auto agent = Havoc->Agent( uuid );

                        if ( agent.has_value() && agent.value()->interface.has_value() ) {
                            try {
                                HcPythonAcquire();

                                agent_action->callback( agent.value()->interface.value() );
                            } catch ( py11::error_already_set& e ) {
                                spdlog::error( "failed to execute action callback: {}", e.what() );
                            }
                        }
                        return;
                    }
                }

                spdlog::debug( "[ERROR] invalid action from selected agent menu" );
            }
        }
    }
}

auto HcPageAgent::handleAgentDoubleClick(
    const QModelIndex& index
) -> void {
    auto uuid = AgentTable->item( index.row(), 0 )->text();

    spawnAgentConsole( uuid.toStdString() );
}

auto HcPageAgent::spawnAgentConsole(
    const std::string& uuid
) -> void {
    for ( auto& agent : agents ) {
        if ( agent->uuid == uuid ) {
            auto user  = agent->data[ "meta" ][ "user" ].get<std::string>();
            auto title = QString( "[%1] %2" ).arg( uuid.c_str() ).arg( user.c_str() );

            //
            // check if we already have the agent console open.
            // if yes then just focus on the opened tab already
            //
            for ( int i = 0; i < AgentTab->count(); i++ ) {
                if ( AgentTab->widget( i ) == agent->console ) {
                    AgentTab->setCurrentIndex( i );
                    return;
                }
            }

            //
            // no tab with the title name found so lets just add a new one
            //
            addTab( title, agent->console );

            break;
        }
    }
}

auto HcPageAgent::tabCloseRequested(
    int index
) const -> void {
    if ( index == -1 ) {
        return;
    }

    AgentTab->removeTab( index );

    if ( AgentTab->count() == 0 ) {
        Splitter->setSizes( QList<int>() << 0 );
        Splitter->handle( 1 )->setEnabled( false );
        Splitter->handle( 1 )->setCursor( Qt::ArrowCursor );
    }
}

auto HcPageAgent::AgentConsole(
    const std::string& uuid,
    const std::string& format,
    const std::string& output
) -> void {
    auto agent = Agent( uuid );

    if ( agent.has_value() ) {
        //
        // now print the content
        //
        agent.value()->console->appendConsole( HcConsole::formatString( format, output ).c_str() );
    } else {
        spdlog::debug( "[AgentConsole] agent not found: {}", uuid );
    }
}

auto HcPageAgent::Agent(
    const std::string &uuid
) -> std::optional<HcAgent*> {
    for ( auto agent : agents ) {
        if ( agent->uuid == uuid ) {
            return agent;
        }
    }

    return std::nullopt;
}

auto HcPageAgent::actionShowHidden(
    bool checked
) -> void {
    // TODO: show hidden
}

auto HcPageAgent::actionPayloadBuilder(
    bool checked
) -> void {
    auto dialog = HcDialogBuilder();

    HcPythonAcquire();

    QObject::connect( Havoc->Gui, &HcMainWindow::signalBuildLog, &dialog, &HcDialogBuilder::EventBuildLog );

    //
    // if there was an error while loading or executing
    // any scripts do not display the window
    //
    if ( ! dialog.ErrorReceived ) {
        dialog.exec();
    }
}

auto HcPageAgent::itemChanged(
    QTableWidgetItem *item
) -> void {
    auto agent_item = ( HcAgentTableItem* ) item;

    //
    // check if it is the agent item is equal to the note widget
    //
    if ( agent_item == agent_item->agent->ui.Note ) {
        //
        // check if is the note item getting created
        //
        if ( agent_item->ignore ) {
            agent_item->ignore = false;
            return;
        }

        auto note   = agent_item->text().toStdString();
        auto result = Havoc->ApiSend( "/api/agent/note", {
            { "uuid", agent_item->agent->uuid },
            { "note", note }
        } );

        spdlog::debug( "result->status: {}", result->status );
        spdlog::debug( "result->body  : {}", result->body );
    }
}

auto HcPageAgent::actionTriggered(
    QAction* triggered
) -> void {
    for ( auto action : Havoc->Actions( HavocClient::ActionObject::ActionHavoc ) ) {
        if ( action->name == triggered->text().toStdString() ) {

            try {
                HcPythonAcquire();

                action->callback();
            } catch ( py11::error_already_set& e ) {
                spdlog::error( "failed to execute action callback: {}", e.what() );
            }

            break;
        }
    }
}

auto HcPageAgent::viewChanged(
    int index
) -> void {
    StackedWidget->setCurrentIndex( index );
}

auto HcPageAgent::removeAgent(
    const std::string& uuid
) -> void {
    HcAgent* agent     = {};
    auto     item_uuid = std::string();

    //
    // remove the agent from
    // the table ui widget entry
    //
    for ( int i = 0; i < AgentTable->rowCount(); i++ ) {
        item_uuid = ( ( HcAgentTableItem* ) AgentTable->item( i, 0 ) )->agent->uuid;

        if ( item_uuid == uuid ) {
            spdlog::debug( "remove agent from table: {}", uuid );
            AgentTable->removeRow( i );
            break;
        }
    }

    //
    // remove the agent entry from
    // the table entry vector list
    //
    for ( int i = 0; i < agents.size(); i++ ) {
        if ( agents[ i ]->uuid == item_uuid ) {
            agent = agents[ i ];

            agents.erase( agents.begin() + i );
            break;
        }
    }

    if ( agent ) {
        HcPythonAcquire();

        AgentTab->removeTab( AgentTab->indexOf( agent->console ) );
        AgentGraph->removeAgent( agent );

        spdlog::debug( "agent delete objects" );


        delete agent->console;
        delete agent;
    }

    AgentDisplayerTargets->setText( QString( "Targets: %1" ).arg( agents.size() ) );   /* TODO: all targets (only show one host)        */
    AgentDisplayerSessions->setText( QString( "Sessions: %1" ).arg( agents.size() ) ); /* TODO: only set current alive beacons/sessions */
}

HcAgentConsole::HcAgentConsole(
    HcAgent* meta,
    QWidget* parent
) : HcConsole( parent ), Meta( meta ) {}

auto HcAgentConsole::inputEnter(
    void
) -> void {
    auto input = std::string();

    input = Input->text().toStdString();
    Input->clear();

    //
    // invoke the command in a separate thread
    //
    auto future = QtConcurrent::run( []( HcAgent* agent, const std::string& input ) {
        HcPythonAcquire();

        if ( agent->interface.has_value() ) {
            try {
                agent->interface.value().attr( "_input_dispatch" )( input );
            } catch ( py11::error_already_set &eas ) {
                emit agent->emitter.ConsoleWrite( agent->uuid.c_str(), eas.what() );
            }
        } else {
            emit agent->emitter.ConsoleWrite( agent->uuid.c_str(), "[!] No agent script handler registered for this type" );
        }
    }, Meta, input );
}