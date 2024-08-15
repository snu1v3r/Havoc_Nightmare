#include <Common.h>
#include <Havoc.h>
#include <api/HcCore.h>

auto HcServerApiSend(
    const std::string& endpoint,
    const json&        data
) -> json {
    auto result = httplib::Result();

    result = Havoc->ApiSend(
        endpoint,
        data
    );

    return json {
        { "version", result->version },
        { "status",  result->status  },
        { "reason",  result->reason  },
        { "body",    result->body    },
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
        spdlog::debug( "with icon: {}", icon_path );
        Havoc->Gui->PageAgent->AgentActionMenu->addAction( QIcon( icon_path.c_str() ), name.c_str() );
    }
}