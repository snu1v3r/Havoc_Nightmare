#include <Havoc.h>
#include <core/HcAgent.h>

auto HcAgent::remove() -> void {
    auto result = httplib::Result();

    result = Havoc->ApiSend( "/api/agent/remove", { { "uuid", uuid } } );

    if ( result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "agent removal failure",
            std::format( "failed to remove agent {}: {}", uuid, result->body )
        );

        spdlog::error( "failed to remove agent {}: {}", uuid, result->body );

        return;
    }
}

auto HcAgent::hide() -> void {
    auto result = httplib::Result();

    result = Havoc->ApiSend( "/api/agent/hide", {
        { "uuid", uuid },
        { "hide", true }
    } );

    if ( result->status != 200 ) {
        Helper::MessageBox(
            QMessageBox::Critical,
            "agent hide failure",
            std::format( "failed to hide agent {}: {}", uuid, result->body )
        );

        spdlog::error( "failed to hide agent {}: {}", uuid, result->body );

        return;
    }
}