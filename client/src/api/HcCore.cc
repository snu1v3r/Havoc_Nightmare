#include <Common.h>
#include <Havoc.h>
#include <api/HcCore.h>

auto HcServerApiSend(
    const std::string& endpoint,
    const json&        data
) -> json {
    auto result = httplib::Result();
    auto body   = json();

    result = Havoc->ApiSend(
        endpoint,
        data
    );

    try {
        body = json::parse( result->body );
    } catch ( std::exception& e ) {
        spdlog::error( "error while parsing body response from {}:\n{}", endpoint, e.what() );
    }

    return json {
        { "status",  result->status  },
        { "reason",  result->reason  },
        { "body",    body },
    };
}

auto HcRegisterMenuAction(
    const std::string&  name,
    const std::string&  icon_path,
    const py11::object& callback
) -> void {
    auto action = new HavocClient::ActionObject();

    action->type = HavocClient::ActionObject::ActionHavoc;
    action->name = name;
    action->icon = icon_path;
    action->callback = callback;

    Havoc->AddAction( action );

    if ( icon_path.empty() ) {
        Havoc->Gui->PageAgent->AgentActionMenu->addAction( name.c_str() );
    } else {
        Havoc->Gui->PageAgent->AgentActionMenu->addAction( QIcon( icon_path.c_str() ), name.c_str() );
    }
}