#ifndef HAVOCCLIENT_API_HCCORE_H
#define HAVOCCLIENT_API_HCCORE_H

auto HcServerApiSend(
    const std::string& endpoint,
    const json&        data
) -> json;

auto HcRegisterMenuAction(
    const std::string&  name,
    const std::string&  icon_path,
    const py11::object& callback
) -> void;

#endif //HAVOCCLIENT_API_HCCORE_H
