#include <Havoc.h>
#include <ui/HcPageAgent.h>
#include <QtConcurrent/QtConcurrent>
#include <DockManager.h>

#include "DockAreaWidget.h"

HcPageAgent::HcPageAgent(
    QWidget* parent
) : QWidget( parent ) {
    if ( objectName().isEmpty() ) {
        setObjectName( "PageAgent" );
    }

    gridLayout = new QGridLayout( this );
    gridLayout->setObjectName( "gridLayout" );

    DockManager = new ads::CDockManager( this );
    DockManager->setStyleSheet( HavocClient::StyleSheet() );

    ads::CDockManager::setConfigFlag( ads::CDockManager::AlwaysShowTabs, false );
    ads::CDockManager::setConfigFlag( ads::CDockManager::DefaultDockAreaButtons, false );
    ads::CDockManager::setConfigFlag( ads::CDockManager::ActiveTabHasCloseButton, false );
    ads::CDockManager::setConfigFlag( ads::CDockManager::AllTabsHaveCloseButton, true );
    ads::CDockManager::setConfigFlag( ads::CDockManager::DockAreaDynamicTabsMenuButtonVisibility, true );
    ads::CDockManager::setConfigFlag( ads::CDockManager::FloatingContainerHasWidgetIcon, false );

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

    AgentTable = new QTableWidget( this );
    AgentTable->setObjectName( "AgentTable" );

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

    AgentGraph = new HcSessionGraph( this );
    AgentGraph->setObjectName( "AgentGraph" );

    auto graph = new ads::CDockWidget( "Session Graph" );
    auto table = new ads::CDockWidget( "Session Table" );

    graph->setWidget( AgentGraph );
    graph->setFeatures( ads::CDockWidget::DockWidgetMovable );
    graph->setIcon( QIcon( ":/icons/32px-network" ) );

    table->setWidget( AgentTable );
    table->setFeatures( ads::CDockWidget::DockWidgetMovable );
    table->setIcon( QIcon( ":/icons/32px-table-list" ) );

    auto area = DockManager->addDockWidget( ads::TopDockWidgetArea, graph );
    DockManager->addDockWidgetTabToArea( table, area );

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

    connect( AgentTable, &QTableWidget::customContextMenuRequested, this, &HcPageAgent::handleAgentMenu );
    connect( AgentTable, &QTableWidget::doubleClicked, this, &HcPageAgent::handleAgentDoubleClick );
    connect( ActionPayload, &QAction::triggered, this, &HcPageAgent::actionPayloadBuilder );
    connect( ActionShowHidden, &QAction::triggered, this, &HcPageAgent::actionShowHidden );
    connect( AgentActionButton, &QToolButton::triggered, this, &HcPageAgent::actionTriggered );
    connect( AgentTable, &QTableWidget::itemChanged, this, &HcPageAgent::itemChanged );

    gridLayout->addWidget( AgentDisplayerSessions, 0, 1, 1, 1 );
    gridLayout->addWidget( AgentDisplayerTargets, 0, 0, 1, 1 );
    gridLayout->addWidget( AgentDisplayerPivots, 0, 2, 1, 1 );
    gridLayout->addWidget( AgentDisplayerElevated, 0, 3, 1, 1 );
    gridLayout->addWidget( AgentActionButton, 0, 5, 1, 1 );
    gridLayout->addWidget( DockManager, 1, 0, 1, 7 );
    gridLayout->addItem( horizontalSpacer, 0, 4, 1, 1 );

    retranslateUi();

    QMetaObject::connectSlotsByName( this );
}

HcPageAgent::~HcPageAgent() = default;

auto HcPageAgent::retranslateUi() -> void {
    setStyleSheet( Havoc->StyleSheet() );

    AgentDisplayerElevated->setText( "Elevated: 0" );
    AgentDisplayerSessions->setText( "Sessions: 0" );
    AgentDisplayerTargets->setText( "Targets: 0" );
    AgentDisplayerPivots->setText( "Pivots: 0" );
}

auto HcPageAgent::addTab(
    const QString& name,
    QWidget*       widget
) -> void {
    //
    // first try to find the widget and
    // if found try to make it reappear
    //
    for ( const auto dock : DockManager->dockWidgetsMap() ) {
        if ( dock->widget() == widget ) {
            dock->toggleView( true );
            return;
        }
    }

    //
    // if not found then create a new widget
    // and add it to the bottom dock area
    //

    const auto tab = new ads::CDockWidget( name );

    tab->setWidget( widget );
    tab->setIcon( QIcon( ":/icons/32px-agent-console" ) );
    tab->setFeatures(
        ads::CDockWidget::DockWidgetMovable   |
        ads::CDockWidget::DockWidgetFloatable |
        ads::CDockWidget::DockWidgetFocusable |
        ads::CDockWidget::DockWidgetClosable
    );

    if ( ! ConsoleAreaWidget ) {
        ConsoleAreaWidget = DockManager->addDockWidget( ads::BottomDockWidgetArea, tab );
    } else {
        DockManager->addDockWidgetTabToArea( tab, ConsoleAreaWidget );
    }
}

HcAgentTableItem::HcAgentTableItem(
    const QString&          value,
    HcAgent*                agent,
    const Qt::ItemFlag      flags,
    const Qt::AlignmentFlag align
)  : agent( agent ) {
    setText( value );
    setTextAlignment( align );
    setFlags( this->flags() ^ flags );
}

auto HcPageAgent::addAgent(
    HcAgent* agent
) -> void {
    const auto row  = AgentTable->rowCount();
    const auto sort = AgentTable->isSortingEnabled();

    //
    // bind the initialized agent to
    // the console write signals
    //
    connect( & agent->ui.signal, & HcAgentSignals::ConsoleWrite, this, []( const QString& uuid, const QString& text ) {
        if ( const auto _agent = Havoc->Agent( uuid.toStdString() ); _agent.has_value() ) {
            _agent.value()->console->appendConsole( HcConsole::formatString( text.toStdString() ).c_str() );
        }
    } );

    agents.push_back( agent );

    AgentTable->setRowCount( row + 1 );
    AgentTable->setSortingEnabled( false );
    AgentTable->setItem( row, 0,  agent->ui.table.Uuid        );
    AgentTable->setItem( row, 1,  agent->ui.table.Internal    );
    AgentTable->setItem( row, 2,  agent->ui.table.Username    );
    AgentTable->setItem( row, 3,  agent->ui.table.Hostname    );
    AgentTable->setItem( row, 4,  agent->ui.table.ProcessName );
    AgentTable->setItem( row, 5,  agent->ui.table.ProcessId   );
    AgentTable->setItem( row, 6,  agent->ui.table.ThreadId    );
    AgentTable->setItem( row, 7,  agent->ui.table.Arch        );
    AgentTable->setItem( row, 8,  agent->ui.table.System      );
    AgentTable->setItem( row, 9,  agent->ui.table.Last        );
    AgentTable->setItem( row, 10, agent->ui.table.Note        );
    AgentTable->setSortingEnabled( sort );

    agent->ui.node = AgentGraph->addAgent( agent );

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
            const auto user  = agent->data[ "meta" ][ "user" ].get<std::string>();
            const auto title = QString( "[%1] %2" ).arg( uuid.c_str() ).arg( user.c_str() );

            //
            // check if we already have the agent console open.
            // if yes then just focus on the opened tab already
            //
            // for ( int i = 0; i < AgentTab->count(); i++ ) {
            //     if ( AgentTab->widget( i ) == agent->console ) {
            //         AgentTab->setCurrentIndex( i );
            //         return;
            //     }
            // }

            //
            // no tab with the title name found so lets just add a new one
            //
            addTab( title, agent->console );

            break;
        }
    }
}

auto HcPageAgent::AgentConsole(
    const std::string& uuid,
    const std::string& format,
    const std::string& output
) -> void {
    if ( const auto agent = Agent( uuid ); agent.has_value() ) {
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
    HcPythonAcquire();

    auto dialog = HcDialogBuilder();

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
    if ( agent_item == agent_item->agent->ui.table.Note ) {
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

        for ( const auto dock : DockManager->dockWidgetsMap() ) {
            if ( dock->widget() == agent->console ) {
                dock->closeDockWidget();
                DockManager->removeDockWidget( dock );
                delete dock;
                return;
            }
        }

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
) : HcConsole( parent ), Meta( meta ) {
    setObjectName( "HcAgentConsole" );
}

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
                emit agent->ui.signal.ConsoleWrite( agent->uuid.c_str(), eas.what() );
            }
        } else {
            emit agent->ui.signal.ConsoleWrite( agent->uuid.c_str(), "[!] No agent script handler registered for this type" );
        }
    }, Meta, input );
}