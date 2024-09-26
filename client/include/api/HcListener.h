#ifndef HAVOCCLIENT_HCLISTENER_H
#define HAVOCCLIENT_HCLISTENER_H

#include <Common.h>
#include <api/Engine.h>

auto HcListenerProtocolData(
    const std::string& protocol
) -> json;

auto HcListenerAll(
    void
) -> std::vector<std::string>;

auto HcListenerQueryType(
    const std::string& name
) -> std::string;

auto HcListenerRegisterMenuAction(
    const std::string&  type,
    const std::string&  name,
    const std::string&  icon,
    const py11::object& callback
) -> void;

auto HcListenerPopupSelect(
    const std::string& protocol
) -> json;

#endif //HAVOCCLIENT_HCLISTENER_H
