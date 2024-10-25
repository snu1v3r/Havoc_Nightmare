#include <Havoc.h>
#include <core/HcMetaWorker.h>

//
// TODO: pull agent console as well if configured
//       handle race condition when agent/listener
//       gets created inbetween the eventWorkerRun
//       and after the pulling has been done
//

HcMetaWorker::HcMetaWorker(
    bool plugin_worker
) : m_plugin( plugin_worker ) {};

HcMetaWorker::~HcMetaWorker() = default;

void HcMetaWorker::run() {
    //
    // process all the plugins, listeners,
    // agents, operator connections, etc.
    //

    if ( m_plugin ) {
        plugins();
    } else {
        listeners();
        agents();
    }

    emit Finished();

    quit();
}

auto HcMetaWorker::listeners(
    void
) -> void {
    auto result    = httplib::Result();
    auto listeners = json();

    result = Havoc->ApiSend( "/api/listener/list", {} );

    if ( result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "listener processing failure",
            "failed to pull all active listeners"
        );
        return;
    }

    try {
        if ( ( listeners = json::parse( result->body ) ).is_discarded() ) {
            spdlog::debug( "listeners processing json response has been discarded" );
            return;
        }
    } catch ( std::exception& e ) {
        spdlog::error( "failed to parse listener processing json response: \n{}", e.what() );
        return;
    }

    if ( listeners.empty() ) {
        spdlog::debug( "no listeners to process" );
        return;
    }

    if ( ! listeners.is_array() ) {
        spdlog::error( "listeners response is not an array" );
        return;
    }

    for ( auto& listener : listeners ) {
        if ( ! listener.is_object() ) {
            spdlog::debug( "warning! listener processing item is not an object" );
            continue;
        }

        spdlog::debug( "processing listener: {}", listener.dump() );

        emit AddListener( listener );
    }
}

auto HcMetaWorker::agents(
    void
) -> void {
    auto result = httplib::Result();
    auto agents = json();
    auto uuid   = std::string();

    result = Havoc->ApiSend( "/api/agent/list", {} );

    if ( result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "agent processing failure",
            "failed to pull all active agent sessions"
        );
        return;
    }

    try {
        if ( ( agents = json::parse( result->body ) ).is_discarded() ) {
            spdlog::debug( "agent processing json response has been discarded" );
            return;
        }
    } catch ( std::exception& e ) {
        spdlog::error( "failed to parse agent processing json response: \n{}", e.what() );
        return;
    }

    if ( agents.empty() ) {
        spdlog::debug( "no agents to process" );
        return;
    }

    if ( ! agents.is_array() ) {
        spdlog::error( "agents response is not an array" );
        return;
    }

    //
    // iterate over all available agents
    // and pull console logs as well
    //
    for ( auto& agent : agents ) {
        if ( ! agent.is_object() ) {
            spdlog::debug( "warning! agent processing item is not an object" );
            continue;
        }

        spdlog::debug( "processing agent: {}", agent.dump() );

        emit AddAgent( agent );

        //
        // parse the agent metadata to retrieve the uuid which will then be
        // passed to the console function to retrieve all console logs
        //
        if ( agent.contains( "uuid" ) && agent[ "uuid" ].is_string() ) {
            console(
                agent[ "uuid" ].get<std::string>()
            );
        } else {
            spdlog::error( "[HcMetaWorker::agents] agent does not contain valid uuid" );
            continue;
        }
    }
}

auto HcMetaWorker::plugins(
    void
) -> void {
    auto result  = httplib::Result();
    auto plugins = json();
    auto uuid    = std::string();
    auto name    = std::string();
    auto version = std::string();

    result = Havoc->ApiSend( "/api/plugin/list", {} );

    if ( result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "plugin processing failure",
            "failed to pull all registered plugins"
        );
        return;
    }

    try {
        if ( ( plugins = json::parse( result->body ) ).is_discarded() ) {
            spdlog::debug( "plugin processing json response has been discarded" );
            return;
        }
    } catch ( std::exception& e ) {
        spdlog::error( "failed to parse plugin processing json response: \n{}", e.what() );
        return;
    }

    if ( plugins.empty() ) {
        spdlog::debug( "no plugins to process" );
        return;
    }

    if ( ! plugins.is_array() ) {
        spdlog::error( "plugins response is not an array" );
        return;
    }

    //
    // iterate over all available agents
    // and pull console logs as well
    //
    for ( auto& plugin : plugins ) {
        if ( ! plugin.is_object() ) {
            spdlog::debug( "warning! plugin processing item is not an object" );
            continue;
        }

        spdlog::debug( "processing plugin: {}", plugin.dump() );

        if ( plugin.contains( "resources" ) && plugin[ "resources" ].is_array() ) {
            //
            // we just really assume that they
            // are contained inside the json object
            //
            name    = plugin[ "name" ].get<std::string>();
            version = plugin[ "version" ].get<std::string>();

            spdlog::debug( "{} ({}):", name, version );
            for ( const auto& res : plugin[ "resources" ].get<std::vector<std::string>>() ) {
                spdlog::debug( " - {}", res );

                if ( ! resource( name, version, res ) ) {
                    spdlog::debug( "failed to pull resources for plugin: {}", name );
                    break;
                }
            }
        }
    }
}

auto HcMetaWorker::console(
    const std::string& uuid
) -> void {
    auto result = httplib::Result();
    auto list   = json();

    result = Havoc->ApiSend( "/api/agent/console", { { "uuid", uuid } } );

    if ( result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "agent console processing failure",
            "failed to pull all agent console logs"
        );
        return;
    }

    try {
        if ( ( list = json::parse( result->body ) ).is_discarded() ) {
            spdlog::debug( "agent console processing json response has been discarded" );
            return;
        }
    } catch ( std::exception& e ) {
        spdlog::error( "failed to parse agent console processing json response: \n{}", e.what() );
        return;
    }

    if ( list.empty() ) {
        spdlog::debug( "no agent console to process" );
        return;
    }

    if ( ! list.is_array() ) {
        spdlog::error( "agents console response is not an array" );
        return;
    }

    for ( auto& console : list ) {
        if ( ! console.is_object() ) {
            spdlog::debug( "warning! agent console processing item is not an object" );
            continue;
        }

        emit AddAgentConsole( console );
    }
}

auto HcMetaWorker::resource(
    const std::string& name,
    const std::string& version,
    const std::string& resource
) -> bool {
    auto result       = httplib::Result();
    auto dir_plugin   = QDir();
    auto dir_resource = QDir();
    auto file_info    = QFileInfo();

    if ( ! ( dir_plugin = QDir( ( Havoc->directory().path().toStdString() + "/plugins/" + name + "@" + version ).c_str() ) ).exists() ) {
        if ( ! dir_plugin.mkpath( "." ) ) {
            spdlog::error( "failed to create plugin directory {}", dir_plugin.path().toStdString() );
            return false;
        }
    }

    file_info    = QFileInfo( ( dir_plugin.path().toStdString() + "/" + resource ).c_str() );
    dir_resource = file_info.absoluteDir();

    if ( ! dir_resource.exists() ) {
        if ( ! dir_resource.mkpath( dir_resource.absolutePath() ) ) {
            spdlog::error( "failed to create resource path: {}", dir_resource.absolutePath().toStdString() );
            return false;
        }
    }

    auto fil_resource = QFile( file_info.absoluteFilePath() );
    if ( ! fil_resource.exists() ) {
        result = Havoc->ApiSend( "/api/plugin/resource", {
            { "name",     name     },
            { "resource", resource },
        } );

        if ( result->status != 200 ) {
            spdlog::debug("failed to pull plugin resource {} from {}: {}", name, resource, result->body);
            return false;
        }

        spdlog::debug( "result->body {}", result->body.length() );

        if ( fil_resource.open( QIODevice::WriteOnly ) ) {
            fil_resource.write( result->body.c_str(), static_cast<qint64>( result->body.length() ) );
        } else {
            spdlog::debug( "failed to open file resource locally {}", file_info.absoluteFilePath().toStdString() );
        }
    } else {
        spdlog::debug( "file resource already exists: {}", file_info.absoluteFilePath().toStdString() );
    }

    return true;
}
