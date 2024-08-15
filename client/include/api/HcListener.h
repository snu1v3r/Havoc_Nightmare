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

#endif //HAVOCCLIENT_HCLISTENER_H
