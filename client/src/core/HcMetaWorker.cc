#include <Havoc.h>
#include <core/HcMetaWorker.h>

HcMetaWorker::HcMetaWorker()  = default;
HcMetaWorker::~HcMetaWorker() = default;

//
// NOTE: start the HcMetaWorker before the HcEventWorker or apply
//       a lock inside the event worker to wait til this thread has finished
//       or add a connect signal to signal or emit an call such as
//       emit HcMetaWorker::eventWorkerStart -> HcEventWorker::run
//

//
// TODO: pull agent console as well if configured
//       handle race condition when agent/listener
//       gets created inbetween the eventWorkerRun
//       and after the pulling has been done
//

void HcMetaWorker::run() {
    //
    // process all the listeners, agents,
    // operator connections, etc.
    //
    listeners();
    agents();

    //
    // we finished running the meta worker thread now we are going
    // to fire up the event worker to listen for incoming events
    //
    emit eventWorkerRun();

    quit();
}

auto HcMetaWorker::listeners() -> void
{
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

auto HcMetaWorker::agents() -> void
{
    auto result = httplib::Result();
    auto agents = json();

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

    for ( auto& agent : agents ) {
        if ( ! agent.is_object() ) {
            spdlog::debug( "warning! agent processing item is not an object" );
            continue;
        }

        spdlog::debug( "processing agent: {}", agent.dump() );

        emit AddAgent( agent );
    }
}
