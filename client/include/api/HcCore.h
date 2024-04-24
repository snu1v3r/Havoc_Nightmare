#ifndef HAVOCCLIENT_API_HCCORE_H
#define HAVOCCLIENT_API_HCCORE_H

auto HcServerApiSend(
    const std::string& endpoint,
    const json&        data
) -> json;

auto HcListenerProtocolData(
    const std::string& protocol
) -> json;

auto HcListenerAll() -> std::vector<std::string>;

auto HcListenerQuery(
    const std::string& name
) -> json;

#endif //HAVOCCLIENT_API_HCCORE_H