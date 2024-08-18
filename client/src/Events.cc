#include <Havoc.h>

auto HavocClient::eventDispatch(
    const json& event
) -> void {
    auto type = std::string();
    auto data = json();

    if ( ! event.contains( "type" ) ) {
        spdlog::debug( "invalid event: {}", event.dump() );
        return;
    }

    if ( ! event.contains( "data" ) ) {
        spdlog::debug( "invalid event: {}", event.dump() );
        return;
    }

    type = event[ "type" ].get<std::string>();
    data = event[ "data" ].get<json>();

    if ( type == Event::user::login )
    {
        if ( data.empty() ) {
            spdlog::error( "user::login: invalid package (data emtpy)" );
            return;
        }
    }
    else if ( type == Event::user::logout )
    {
        if ( data.empty() ) {
            spdlog::error( "user::logout: invalid package (data emtpy)" );
            return;
        }
    }
    else if ( type == Event::user::message )
    {

    }
    else if ( type == Event::listener::add )
    {
        if ( data.empty() ) {
            spdlog::error( "listener::register: invalid package (data emtpy)" );
            return;
        }

        Gui->PageListener->Protocols.push_back( data );
    }
    else if ( type == Event::listener::start )
    {
        auto name   = std::string();
        auto status = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::listener::start: invalid package (data emtpy)" );
            return;
        }

        Gui->AddListener( data );
    }
    else if ( type == Event::listener::stop )
    {
        auto name   = std::string();
        auto status = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::listener::status: invalid package (data emtpy)" );
            return;
        }

        if ( data.contains( "name" ) ) {
            if ( data[ "name" ].is_string() ) {
                name = data[ "name" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener status: \"name\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener status: \"name\" is not found" );
            return;
        }

        if ( data.contains( "status" ) ) {
            if ( data[ "status" ].is_string() ) {
                status = data[ "status" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener status: \"status\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener status: \"status\" is not found" );
            return;
        }

        Gui->PageListener->setListenerStatus( name, status );
    }
    else if ( type == Event::listener::status )
    {
        auto name = std::string();
        auto log  = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::listener::status: invalid package (data emtpy)" );
            return;
        }

        if ( data.contains( "name" ) ) {
            if ( data[ "name" ].is_string() ) {
                name = data[ "name" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener status: \"name\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener status: \"name\" is not found" );
            return;
        }

        if ( data.contains( "status" ) ) {
            if ( data[ "status" ].is_string() ) {
                log = data[ "status" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener status: \"status\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener status: \"status\" is not found" );
            return;
        }

        Gui->PageListener->setListenerStatus( name, log );
    }
    else if ( type == Event::listener::log )
    {
        auto name = std::string();
        auto log  = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::listener::log: invalid package (data emtpy)" );
            return;
        }

        if ( data.contains( "name" ) ) {
            if ( data[ "name" ].is_string() ) {
                name = data[ "name" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener log: \"name\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener log: \"name\" is not found" );
            return;
        }

        if ( data.contains( "log" ) ) {
            if ( data[ "log" ].is_string() ) {
                log = data[ "log" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener log: \"log\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener log: \"log\" is not found" );
            return;
        }

        Gui->PageListener->addListenerLog( name, log );
    }
    else if ( type == Event::listener::remove )
    {
        auto name = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::listener::remove: invalid package (data emtpy)" );
            return;
        }

        if ( data.contains( "name" ) ) {
            if ( data[ "name" ].is_string() ) {
                name = data[ "name" ].get<std::string>();
            } else {
                spdlog::error( "invalid listener remove: \"name\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid listener remove: \"name\" is not found" );
            return;
        }

        RemoveListener( name );
    }
    else if ( type == Event::agent::add )
    {
        //
        // nothing
        //
    }
    else if ( type == Event::agent::initialize )
    {
        if ( data.empty() ) {
            spdlog::error( "Event::agent::initialize: invalid package (data emtpy)" );
            return;
        }

        Gui->AddAgent( data );
    }
    else if ( type == Event::agent::callback )
    {
        auto uuid          = std::string();
        auto typ           = std::string();
        auto arg           = json();
        auto ctx           = json();
        auto pat           = std::string();
        auto out           = std::string();
        auto callback_uuid = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::agent::callback: invalid package (data emtpy)" );
            return;
        }

        if ( data.contains( "uuid" ) ) {
            if ( data[ "uuid" ].is_string() ) {
                uuid = data[ "uuid" ].get<std::string>();
            } else {
                spdlog::error( "invalid agent callback: \"uuid\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid agent callback: \"uuid\" is not found" );
            return;
        }

        if ( data.contains( "type" ) ) {
            if ( data[ "type" ].is_string() ) {
                typ = data[ "type" ].get<std::string>();
            } else {
                spdlog::error( "invalid agent callback: \"type\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid agent callback: \"type\" is not found" );
            return;
        }

        if ( typ == "console" ) {
            if ( data.contains( "data" ) ) {
                if ( data[ "data" ].is_object() ) {
                    arg = data[ "data" ].get<json>();
                } else {
                    spdlog::error( "invalid agent callback: \"data\" is not an object" );
                    return;
                }
            } else {
                spdlog::error( "invalid agent callback: \"data\" is not found" );
                return;
            }

            if ( arg.contains( "format" ) ) {
                if ( arg[ "format" ].is_string() ) {
                    pat = arg[ "format" ].get<std::string>();
                } else {
                    spdlog::error( "invalid agent callback: \"format\" is not string" );
                    return;
                }
            } else {
                spdlog::error( "invalid agent callback: \"format\" is not found" );
                return;
            }

            if ( arg.contains( "output" ) ) {
                if ( arg[ "output" ].is_string() ) {
                    out = arg[ "output" ].get<std::string>();
                } else {
                    spdlog::error( "invalid agent callback: \"output\" is not string" );
                    return;
                }
            } else {
                spdlog::error( "invalid agent callback: \"output\" is not found" );
                return;
            }

            Gui->PageAgent->AgentConsole( uuid, pat, out );
        } else if ( typ == "function" ) {
            //
            // this type indicates that a callback function should be invoked
            // with the following callback uuid (inside the uuid member)
            //
            if ( data.contains( "data" ) ) {
                if ( data[ "data" ].is_object() ) {
                    arg = data[ "data" ].get<json>();

                    //
                    // acquire python gil to interact with the
                    // python function and call the callback
                    //
                    auto gil = py11::gil_scoped_acquire();

                    //
                    // get the data to pass to the callback
                    //
                    if ( arg.contains( "data" ) ) {
                        if ( arg[ "data" ].is_string() ) {
                            out = arg[ "data" ].get<std::string>();
                        } else {
                            spdlog::error( "invalid agent callback: \"data\" is not string" );
                            return;
                        }
                    } else {
                        spdlog::error( "invalid agent callback: \"data\" is not found" );
                        return;
                    }

                    //
                    // get the task uuid to pass to the callback
                    //
                    if ( arg.contains( "callback" ) ) {
                        if ( arg[ "callback" ].is_string() ) {
                            callback_uuid = arg[ "callback" ].get<std::string>();
                        } else {
                            spdlog::error( "invalid agent callback: \"callback\" is not a string" );
                            return;
                        }
                    } else {
                        spdlog::error( "invalid agent callback: \"callback\" is not found" );
                        return;
                    }

                    //
                    // get context from the callback that needs to
                    // be passed to the callback python function
                    //
                    if ( arg.contains( "context" ) ) {
                        if ( arg[ "context" ].is_object() ) {
                            ctx = arg[ "context" ].get<json>();
                        } else {
                            spdlog::error( "invalid agent callback: \"context\" is not an object" );
                            return;
                        }
                    } else {
                        spdlog::error( "invalid agent callback: \"context\" is not found" );
                        return;
                    }

                    //
                    // get registered callback object
                    //
                    if ( auto callback = Havoc->CallbackObject( callback_uuid ) ) {
                        //
                        // did we found a callback based on the uuid ?
                        //
                        if ( callback.has_value() ) {
                            //
                            // search for the agent via the specified uuid
                            //
                            if ( auto agent = Havoc->Gui->PageAgent->Agent( uuid ) ) {
                                //
                                // check if we found the agent and if
                                // it contains a python interface
                                //
                                if ( agent.has_value() && agent.value()->interface.has_value() ) {
                                    pat = QByteArray::fromBase64( out.c_str() ).toStdString();

                                    try {
                                        //
                                        // actually invoke the callback with following arguments:
                                        //
                                        //   def callback( agent: HcAgent, data: bytes, **kwargs ):
                                        //      # agent interface can be invoked now
                                        //      return
                                        //
                                        callback.value()(
                                            agent.value()->interface.value(),
                                            py11::bytes( pat.c_str(), pat.length() ),
                                            **py11::dict( ctx )
                                        );
                                    } catch ( py11::error_already_set &eas ) {
                                        //
                                        // catch exception and print it to the agent
                                        // console as it was running under its context
                                        //
                                        emit agent.value()->emitter.ConsoleWrite( agent.value()->uuid.c_str(), eas.what() );
                                    }
                                } else {
                                    spdlog::error(
                                        "[agent.has_value(): {}] [agent.value()->interface.has_value(): {}]",
                                        agent.has_value(),
                                        agent.has_value() ? agent.value()->interface.has_value() : false
                                    );
                                }
                            }
                        } else {
                            spdlog::error( "CallbackObject has no value: {}", callback_uuid );
                        }
                    } else {
                        spdlog::error( "CallbackObject not found: {}", callback_uuid );
                    }
                } else {
                    spdlog::error( "invalid agent callback: \"data\" is not an object" );
                    return;
                }
            } else {
                spdlog::error( "invalid agent callback: \"data\" is not found" );
                return;
            }
        } else if ( typ == "heartbeat" ) {
            if ( data.contains( "data" ) ) {
                if ( data[ "data" ].is_object() ) {
                    arg = data[ "data" ].get<json>();
                } else {
                    spdlog::error( "invalid agent heartbeat: \"data\" is not an object" );
                    return;
                }
            } else {
                spdlog::error( "invalid agent heartbeat: \"data\" is not found" );
                return;
            }

            //
            // get the data to pass to the callback
            //
            if ( arg.contains( "time" ) ) {
                if ( arg[ "time" ].is_string() ) {
                    out = arg[ "time" ].get<std::string>();
                } else {
                    spdlog::error( "invalid agent heartbeat: \"time\" is not string" );
                    return;
                }
            } else {
                spdlog::error( "invalid agent heartbeat: \"time\" is not found" );
                return;
            }

            //
            // set the last callback of the agent
            //
            if ( auto agent = Agent( uuid ) ) {
                if ( agent.has_value() ) {
                    agent.value()->last = QString( out.c_str() );
                } else {
                    spdlog::error( "invalid agent heartbeat: \"uuid\" agent does not have any value" );
                }
            } else {
                spdlog::error( "invalid agent heartbeat: \"uuid\" agent not found" );
            }
        } else {
            spdlog::error( "invalid agent typ: {}", typ );
        }
    }
    else if ( type == Event::agent::console )
    {

    }
    else if ( type == Event::agent::input )
    {

    }
    else if ( type == Event::agent::task )
    {

    }
    else if ( type == Event::agent::status )
    {

    }
    else if ( type == Event::agent::remove )
    {

    }
    else if ( type == Event::agent::buildlog )
    {
        auto log = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::agent::buildlog: invalid package (data emtpy)" );
            return;
        }

        if ( data.contains( "log" ) ) {
            if ( data[ "log" ].is_string() ) {
                log = data[ "log" ].get<std::string>();

                //
                // signal opened dialog of the build message
                //
                emit Gui->signalBuildLog( log.c_str() );
            } else {
                spdlog::error( "invalid agent build log: \"log\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid agent build log: \"log\" is not found" );
            return;
        }
    } else if ( type == Event::agent::note ) {
        auto uuid = std::string();
        auto note = std::string();

        if ( data.empty() ) {
            spdlog::error( "Event::agent::note: invalid package (data emtpy)" );
            return;
        }

        spdlog::debug( "note: {}", data.dump() );

        if ( data.contains( "uuid" ) ) {
            if ( data[ "uuid" ].is_string() ) {
                uuid = data[ "uuid" ].get<std::string>();
            } else {
                spdlog::error( "invalid agent callback: \"uuid\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid agent callback: \"uuid\" is not found" );
            return;
        }

        if ( data.contains( "note" ) ) {
            if ( data[ "note" ].is_string() ) {
                note = data[ "note" ].get<std::string>();
            } else {
                spdlog::error( "invalid agent callback: \"note\" is not string" );
                return;
            }
        } else {
            spdlog::error( "invalid agent callback: \"note\" is not found" );
            return;
        }

        //
        // set the note of the agent
        //
        if ( auto agent = Agent( uuid ) ) {
            if ( agent.has_value() ) {
                agent.value()->ui.Note->ignore = true;
                agent.value()->ui.Note->setText( note.c_str() );
                agent.value()->ui.Note->ignore = false;
            } else {
                spdlog::error( "invalid agent note: \"uuid\" agent does not have any value" );
            }
        } else {
            spdlog::error( "invalid agent note: \"uuid\" agent not found" );
        }

    } else {
        spdlog::debug( "invalid event: {} not found", type );
    }
}


