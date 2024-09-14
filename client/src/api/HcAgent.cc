#include <Havoc.h>
#include <api/HcAgent.h>

/*!
 * @brief
 *  register an agent interface to the havoc client
 *
 * @param type
 *  type of the agent (name)
 *
 * @param object
 *  object to register
 */
auto HcAgentRegisterInterface(
    const std::string&  type,
    const py11::object& object
) -> void {
    Havoc->AddAgentObject( type, object );
}

/*!
 * @brief
 *  writes content to specified agent console
 *  bottom label
 *
 * @param uuid
 *  uuid of the agent
 *
 * @param content
 *  content to set as the label of the console
 */
auto HcAgentConsoleLabel(
    const std::string& uuid,
    const std::string& content
) -> void {
    auto agent = Havoc->Agent( uuid );

    if ( agent.has_value() ) {
        //
        // if we are going to append nothing then hide the label
        //
        agent.value()->console->setBottomLabel( content.c_str() );
        agent.value()->console->LabelBottom->setFixedHeight( content.empty() ? 0 : 20 );
    }
}

/*!
 * @brief
 *  writes content to specified agent console header
 *
 * @param uuid
 *  uuid of the agent
 *
 * @param content
 *  content to set as the header of the console
 */
auto HcAgentConsoleHeader(
    const std::string& uuid,
    const std::string& content
) -> void {
    auto agent = Havoc->Agent( uuid );

    if ( agent.has_value() ) {
        //
        // if we are going to append nothing then hide the label
        //
        agent.value()->console->setHeaderLabel( content.c_str() );
        agent.value()->console->LabelHeader->setFixedHeight( content.empty() ? 0 : 20 );
    }
}

/*!
 * @brief
 *  writes content to specified agent console
 *
 * @param uuid
 *  uuid of the agent
 *
 * @param content
 *  content to write to the console
 */
auto HcAgentConsoleWrite(
    const std::string& uuid,
    const std::string& content
) -> void {
    auto agent = Havoc->Agent( uuid );

    if ( agent.has_value() ) {
        emit agent.value()->ui.signal.ConsoleWrite( uuid.c_str(), content.c_str() );
    }
}

/*!
 * @brief
 *  register a callback
 *
 * @param uuid
 *  uuid of the callback to be registered
 *
 * @param callback
 *  callback function to call
 */
auto HcAgentRegisterCallback(
    const std::string&  uuid,
    const py11::object& callback
) -> void {
    Havoc->AddCallbackObject( uuid, callback );
}

/*!
 * @brief
 *  unregister a callback
 *
 * @param uuid
 *  callback uuid to remove and unregister
 */
auto HcAgentUnRegisterCallback(
    const std::string& uuid
) -> void {
    Havoc->RemoveCallbackObject( uuid );
}

/*!
 * @brief
 *  send agent command to the server implant plugin handler
 *
 * @param uuid
 *  uuid of the agent
 *
 * @param data
 *  data to send to the handler
 *
 * @param wait
 *  wait for a response
 *
 * @return
 *  response from the server implant handler
 */
auto HcAgentExecute(
    const std::string& uuid,
    const json&        data,
    const bool         wait
) -> json {
    auto future   = QFuture<json>();
    auto request  = json();
    auto result   = httplib::Result();
    auto response = json();

    //
    // build request that is going to be
    // sent to the server implant handler
    //
    request = {
        { "uuid", uuid },
        { "wait", wait },
        { "data", data }
    };

    //
    // send api request
    //
    if ( ( result = Havoc->ApiSend( "/api/agent/execute", request, true ) ) ) {
        //
        // check for valid status response
        //
        if ( result->status != 200 ) {
            spdlog::debug( "failed to send request: status code {}", result->status );
            return json {
                    { "error", "failed to send request" }
            };
        }
        //
        // check for emtpy request
        //
        if ( ! result->body.empty() ) {
            if ( ( response = json::parse( result->body ) ).is_discarded() ) {
                response[ "error" ] = "failed to parse response";
            }
        }
    }

    return response;
}

/*!
 * @brief
 *  register agent context menu action
 *
 * @param agent_type
 *  agent type to add action to
 *
 * @param name
 *  name of the action to register
 *
 * @param icon_path
 *  icon path to add to action
 *
 * @param callback
 *  python callback function
 *  to call on action trigger
 */
auto HcAgentRegisterMenuAction(
    const std::string&  agent_type,
    const std::string&  name,
    const std::string&  icon_path,
    const py11::object& callback
) -> void {
    auto action = new HavocClient::ActionObject();

    action->type       = HavocClient::ActionObject::ActionAgent;
    action->name       = name;
    action->icon       = icon_path;
    action->callback   = callback;
    action->agent.type = agent_type;

    spdlog::debug( "HcAgentRegisterMenuAction( {}, {}, {} )", agent_type, name, icon_path );

    Havoc->AddAction( action );
}