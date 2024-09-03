#include <Havoc.h>
#include <QTimer>
#include <QtCore5Compat/QTextCodec>

auto HttpErrorToString(
    const httplib::Error& error
) -> std::optional<std::string> {
    switch ( error ) {
        case httplib::Error::Unknown:
            return "Unknown";

        case httplib::Error::Connection:
            return ( "Connection" );

        case httplib::Error::BindIPAddress:
            return ( "BindIPAddress" );

        case httplib::Error::Read:
            return ( "Read" );

        case httplib::Error::Write:
            return ( "Write" );

        case httplib::Error::ExceedRedirectCount:
            return ( "ExceedRedirectCount" );

        case httplib::Error::Canceled:
            return ( "Canceled" );

        case httplib::Error::SSLConnection:
            return ( "SSLConnection" );

        case httplib::Error::SSLLoadingCerts:
            return ( "SSLLoadingCerts" );

        case httplib::Error::SSLServerVerification:
            return ( "SSLServerVerification" );

        case httplib::Error::UnsupportedMultipartBoundaryChars:
            return ( "UnsupportedMultipartBoundaryChars" );

        case httplib::Error::Compression:
            return ( "Compression" );

        case httplib::Error::ConnectionTimeout:
            return ( "ConnectionTimeout" );

        case httplib::Error::ProxyConnection:
            return ( "ProxyConnection" );

        case httplib::Error::SSLPeerCouldBeClosed_:
            return ( "SSLPeerCouldBeClosed_" );

        default: break;
    }

    return std::nullopt;
}

HavocClient::HavocClient() {
    /* initialize logger */
    spdlog::set_pattern( "[%T %^%l%$] %v" );
    spdlog::info( "Havoc Framework [{} :: {}]", HAVOC_VERSION, HAVOC_CODENAME );

    /* enabled debug messages */
    spdlog::set_level( spdlog::level::debug );

    /* TODO: read this from the config file */
    const auto family = "Monospace";
    const auto size   = 9;

#ifdef Q_OS_MAC
    //
    // orignal fix and credit: https://github.com/HavocFramework/Havoc/pull/466/commits/8b75de9b4632a266badf64e1cc22e57cc55a5b7c
    //
    QApplication::setStyle( "Fusion" );
#endif

    //
    // set font
    //
    QTextCodec::setCodecForLocale( QTextCodec::codecForName( "UTF-8" ) );
    QApplication::setFont( QFont( family, size ) );
    QTimer::singleShot( 10, [&]() { QApplication::setFont( QFont( family, size ) ); } );
}

HavocClient::~HavocClient() = default;

/*!
 * @brief
 *  this entrypoint executes the connector dialog
 *  and tries to connect and login to the teamserver
 *
 *  after connecting it is going to start an event thread
 *  and starting the Havoc MainWindow.
 */
auto HavocClient::Main(
    int    argc,
    char** argv
) -> void {
    auto Connector = new HcConnectDialog();
    auto Result    = httplib::Result();
    auto Response  = json{};
    auto Error     = std::string( "Failed to send login request: " );

    /* get provided creds */
    auto data = Connector->start();
    if ( data.empty() || ! Connector->pressedConnect() ) {
        return;
    }

    if ( ! SetupDirectory() ) {
        spdlog::error( "failed to setup configuration directory. aborting" );
        return;
    }

    /* create http client */
    auto Http = httplib::Client( "https://" + data[ "host" ].get<std::string>() + ":" + data[ "port" ].get<std::string>() );
    Http.enable_server_certificate_verification( false );

    /* send request */
    Result = Http.Post( "/api/login", data.dump(), "application/json" );

    if ( HttpErrorToString( Result.error() ).has_value() ) {
        auto error = HttpErrorToString( Result.error() ).value();
        spdlog::error( "Failed to send login request: {}", error );

        Helper::MessageBox(
            QMessageBox::Critical,
            "Login failure",
            std::format( "Failed to login: {}", error )
        );

        return;
    }

    /* 401 Unauthorized: Failed to log in */
    if ( Result->status == 401 ) {
        if ( Result->body.empty() ) {
            Helper::MessageBox(
                QMessageBox::Critical,
                "Login failure",
                "Failed to login: Unauthorized"
            );

            return;
        } else {
            if ( ( data = json::parse( Result->body ) ).is_discarded() ) {
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
                "Login failure",
                QString( "Failed to login: %1" ).arg( data[ "error" ].get<std::string>().c_str() ).toStdString()
            );

            return;
        }

    } else if ( Result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "Login failure",
            QString( "Unexpected response: Http status code %1" ).arg( Result->status ).toStdString()
        );
        return;
    }

    spdlog::debug( "Result: {}", Result->body );

    Profile.Name = data[ "name" ].get<std::string>();
    Profile.Host = data[ "host" ].get<std::string>();
    Profile.Port = data[ "port" ].get<std::string>();
    Profile.User = data[ "username" ].get<std::string>();
    Profile.Pass = data[ "password" ].get<std::string>();

    if ( Result->body.empty() ) {
        goto ERROR_SERVER_RESPONSE;
    }

    if ( ( data = json::parse( Result->body ) ).is_discarded() ) {
        goto ERROR_SERVER_RESPONSE;
    }

    if ( ! data.contains( "token" ) ) {
        goto ERROR_SERVER_RESPONSE;
    }

    if ( ! data[ "token" ].is_string() ) {
        goto ERROR_SERVER_RESPONSE;
    }

    Profile.Token = data[ "token" ].get<std::string>();

    //
    // create main window
    //
    Gui = new HcMainWindow;
    Gui->renderWindow();
    Gui->setStyleSheet( StyleSheet() );
    Theme.setStyleSheet( StyleSheet() );

    //
    // setup Python thread
    //
    PyEngine = new HcPyEngine();

    //
    // load the registered scripts
    //
    if ( Config.contains( "scripts" ) && Config.at( "scripts" ).is_table() ) {
        auto scripts_tbl = Config.at( "scripts" ).as_table();

        if ( scripts_tbl.contains( "files" ) && scripts_tbl.at( "files" ).is_array() ) {
            for ( const auto& file : scripts_tbl.at( "files" ).as_array() ) {
                if ( ! file.is_string() ) {
                    spdlog::error( "configuration scripts file value is not an string" );
                    continue;
                }

                ScriptLoad( file.as_string() );
            }
        }
    }

    //
    // set up the event thread and connect to the
    // server and dispatch all the incoming events
    //
    SetupThreads();

    QApplication::exec();
    return;

ERROR_SERVER_RESPONSE:
    return Helper::MessageBox(
        QMessageBox::Critical,
        "Login failure",
        "Failed to login: Invalid response from the server"
    );
}

auto HavocClient::Exit() -> void {
    QApplication::exit( 0 );
}

auto HavocClient::ApiSend(
    const std::string& endpoint,
    const json&        body,
    const bool         keep_alive
) const -> httplib::Result {
    auto Http   = httplib::Client( "https://" + Profile.Host + ":" + Profile.Port );
    auto Result = httplib::Result();
    auto Error  = std::string( "Failed to send api request: " );

    //
    // only way to keep the connection alive even while we have
    // "keep-alive" enabled it will shut down after 5 seconds
    //
    if ( keep_alive ) {
        Http.set_read_timeout( INT32_MAX );
        Http.set_connection_timeout( INT32_MAX );
        Http.set_write_timeout( INT32_MAX );
    }

    //
    // configure the client
    //
    Http.set_keep_alive( keep_alive );
    Http.enable_server_certificate_verification( false );
    Http.set_default_headers( {
        { "x-havoc-token", Havoc->Profile.Token }
    } );

    //
    // send the request to our endpoint
    //
    Result = Http.Post( endpoint, body.dump(), "application/json" );

    if ( HttpErrorToString( Result.error() ).has_value() ) {
        spdlog::error( "Failed to send login request: {}", HttpErrorToString( Result.error() ).value() );
    }

    return Result;
}

auto HavocClient::eventClosed() -> void {
    spdlog::error( "websocket closed" );
    Exit();
}

auto HavocClient::Server() const -> std::string { return Profile.Host + ":" + Profile.Port; }
auto HavocClient::Token()  const -> std::string { return Profile.Token; }

auto HavocClient::eventHandle(
    const QByteArray& request
) -> void {
    auto event = json::parse( request.toStdString() );

    /* check if we managed to parse the json event
     * if yes then dispatch it but if not then dismiss it */
    if ( ! event.is_discarded() ) {
        eventDispatch( event );
    } else {
        spdlog::error( "failed to parse event" );
        /* what now ?
         * I guess ignore since its not valid event
         * or debug print it I guess */
    }
}

auto HavocClient::StyleSheet(
    void
) -> QByteArray {
    if ( QFile::exists( "theme.css" ) ) {
        return Helper::FileRead( "theme.css" );
    }

    return Helper::FileRead( ":/style/default" );
}

auto HavocClient::AddProtocol(
    const std::string&  name,
    const py11::object& listener
) -> void {
    protocols.push_back( NamedObject{
        .name   = name,
        .object = listener
    } );
}

auto HavocClient::ProtocolObject(
    const std::string& name
) -> std::optional<py11::object> {
    for ( auto& listener : protocols ) {
        if ( listener.name == name ) {
            return listener.object;
        }
    }

    return std::nullopt;
}

auto HavocClient::Protocols() -> std::vector<std::string> {
    auto names = std::vector<std::string>();

    for ( auto& listener : protocols ) {
        names.push_back( listener.name );
    }

    return names;
}

auto HavocClient::SetupThreads() -> void {
    //
    // now set up the event thread and dispatcher
    //
    Events.Thread = new QThread;
    Events.Worker = new HcEventWorker;
    Events.Worker->moveToThread( Events.Thread );

    QObject::connect( Events.Thread, &QThread::started, Events.Worker, &HcEventWorker::run );
    QObject::connect( Events.Worker, &HcEventWorker::availableEvent, this, &HavocClient::eventHandle );
    QObject::connect( Events.Worker, &HcEventWorker::socketClosed, this, &HavocClient::eventClosed );

    //
    // fire up the even thread that is going to
    // process events and emit signals to the main gui thread
    //
    // NOTE: manually starting the thread is no longer needed as it is going to
    //       be started and triggered once the meta worker thread has finished running
    //
    // Events.Thread->start();

    //
    // start the heartbeat worker thread
    //
    Heartbeat.Thread = new QThread;
    Heartbeat.Worker = new HcHeartbeatWorker;
    Heartbeat.Worker->moveToThread( Heartbeat.Thread );

    QObject::connect( Heartbeat.Thread, &QThread::started, Heartbeat.Worker, &HcHeartbeatWorker::run );

    //
    // fire up the even thread that is going to
    // process heart beat events
    //
    Heartbeat.Thread->start();

    //
    // start the heartbeat worker thread
    //
    MetaWorker.Thread = new QThread;
    MetaWorker.Worker = new HcMetaWorker;
    MetaWorker.Worker->moveToThread( MetaWorker.Thread );

    QObject::connect( MetaWorker.Thread, &QThread::started, MetaWorker.Worker, &HcMetaWorker::run );

    //
    // connect methods to add listeners, agents, etc. to the user interface (ui)
    //
    QObject::connect( MetaWorker.Worker, &HcMetaWorker::AddListener,     Gui, &HcMainWindow::AddListener  );
    QObject::connect( MetaWorker.Worker, &HcMetaWorker::AddAgent,        Gui, &HcMainWindow::AddAgent     );
    QObject::connect( MetaWorker.Worker, &HcMetaWorker::AddAgentConsole, Gui, &HcMainWindow::AgentConsole );

    //
    // only start the event worker once the meta worker finished
    // pulling all the listeners, agents, etc. from the server
    //
    QObject::connect( MetaWorker.Worker, &HcMetaWorker::eventWorkerRun, this, [&](){ Events.Thread->start(); } );

    MetaWorker.Thread->start();
}

auto HavocClient::AddBuilder(
    const std::string & name,
    const py11::object& builder
) -> void {
    builders.push_back( NamedObject{
        .name   = name,
        .object = builder
    } );
}

auto HavocClient::BuilderObject(
    const std::string& name
) -> std::optional<py11::object> {

    for ( auto& builder : builders ) {
        if ( builder.name == name ) {
            return builder.object;
        }
    }

    return std::nullopt;
}

auto HavocClient::Builders() -> std::vector<std::string>
{
    auto names = std::vector<std::string>();

    for ( auto& builder : builders ) {
        names.push_back( builder.name );
    }

    return names;
}

auto HavocClient::AddListener(
    const json& listener
) -> void {
    spdlog::debug( "listener -> {}", listener.dump() );
    listeners.push_back( listener );

    Gui->PageListener->addListener( listener );
}

auto HavocClient::RemoveListener(
    const std::string& name
) -> void {
    spdlog::debug( "removing listener: {}", name );

    auto object = ListenerObject( name );

    //
    // remove listener from listener list
    //
    for ( auto iter = listeners.begin(); iter != listeners.end(); ++iter ) {
        if ( object.has_value() ) {
            auto data = object.value();
            if ( data.contains( "name" ) && data[ "name" ].is_string() ) {
                if ( data[ "name" ].get<std::string>() == name ) {
                    listeners.erase( iter );
                    break;
                }
            }
        }
    }

    Gui->PageListener->removeListener( name );
}

auto HavocClient::ListenerObject(
    const std::string &name
) -> std::optional<json> {

    for ( auto& data : listeners ) {
        if ( data.contains( "name" ) && data[ "name" ].is_string() ) {
            if ( data[ "name" ].get<std::string>() == name ) {
                return data;
            }
        }
    }

    return std::nullopt;
}

auto HavocClient::Listeners() -> std::vector<std::string>
{
    auto names = std::vector<std::string>();

    for ( auto& data : listeners ) {
        if ( data.contains( "name" ) && data[ "name" ].is_string() ) {
            names.push_back( data[ "name" ].get<std::string>() );
        }
    }

    return names;
}

auto HavocClient::Agents() -> std::vector<HcAgent *>
{
    return Gui->PageAgent->agents;
}

auto HavocClient::Agent(
    const std::string& uuid
) const -> std::optional<HcAgent*> {

    for ( auto agent : Gui->PageAgent->agents ) {
        if ( agent->uuid == uuid ) {
            return agent;
        }
    }

    return std::nullopt;
}

auto HavocClient::AddAgentObject(
    const std::string&  type,
    const py11::object& object
) -> void {
    agents.push_back( NamedObject {
        .name   = type,
        .object = object
    } );
}

auto HavocClient::AgentObject(
    const std::string& type
) -> std::optional<py11::object> {
    for ( auto object : agents ) {
        if ( object.name == type ) {
            return object.object;
        }
    }

    return std::nullopt;
}

auto HavocClient::Callbacks() -> std::vector<std::string>
{
    auto names = std::vector<std::string>();

    for ( auto& callback : callbacks ) {
        names.push_back( callback.name );
    }

    return names;
}

auto HavocClient::AddCallbackObject(
    const std::string&  uuid,
    const py11::object& callback
) -> void {
    callbacks.push_back( NamedObject{
        .name   = uuid,
        .object = callback
    } );
}

auto HavocClient::RemoveCallbackObject(
    const std::string& callback
) -> void {
    //
    // iterate over the callbacks
    //
    for ( auto iter = callbacks.begin(); iter != callbacks.end(); ++iter ) {
        if ( iter->name == callback ) {
            callbacks.erase( iter );
            break;
        }
    }
}

auto HavocClient::CallbackObject(
    const std::string& uuid
) -> std::optional<py11::object> {
    //
    // iterate over registered callbacks
    //
    for ( auto object : callbacks ) {
        if ( object.name == uuid ) {
            return object.object;
        }
    }

    return std::nullopt;
}

auto HavocClient::SetupDirectory(
    void
) -> bool {
    auto havoc_dir   = QDir( QDir::homePath() + "/.havoc" );
    auto config_path = QFile();

    if ( ! havoc_dir.exists() ) {
        if ( ! havoc_dir.mkpath( "." ) ) {
            return false;
        }
    }

    client_dir.setPath( havoc_dir.path() + "/client" );
    if ( ! client_dir.exists() ) {
        if ( ! client_dir.mkpath( "." ) ) {
            return false;
        }
    }

    config_path.setFileName( client_dir.path() + "/config.toml" );
    config_path.open( QIODevice::ReadWrite );

    try {
        Config = toml::parse( config_path.fileName().toStdString() );
    } catch ( std::exception& e ) {
        spdlog::error( "failed to parse toml configuration: {}", e.what() );
        return false;
    }

    return true;
}

auto HavocClient::AddAction(
    ActionObject* action
) -> void {
    if ( ! action ) {
        spdlog::debug( "[AddAction] action is a nullptr" );
        return;
    }

    for ( auto _action : Actions( action->type ) ) {
        if ( _action->name == action->name ) {
            spdlog::debug( "[AddAction] action with the name \"{}\" and type already exists", action->name );
            return;
        }
    }

    actions.push_back( action );
}

auto HavocClient::Actions(
    const ActionObject::ActionType& type
) -> std::vector<ActionObject*> {
    auto actions_type = std::vector<ActionObject*>();

    for ( auto action : actions ) {
        if ( action->type == type ) {
            actions_type.push_back( action );
        }
    }

    return actions_type;
}

auto HavocClient::ScriptLoad(
    const std::string& path
) -> void {
    Gui->PageScripts->LoadScript( path );
}

HavocClient::ActionObject::ActionObject() {}
