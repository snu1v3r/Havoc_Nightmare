#ifndef HAVOCCLIENT_EVENTS_H
#define HAVOCCLIENT_EVENTS_H

#include <string>

namespace Event {

    namespace user {
        static std::string login   = "user::login";
        static std::string logout  = "user::logout";
        static std::string message = "user::message";
    }

    namespace agent {
        static std::string add        = "agent::register"; // this is for anonymous agents without any checkin data
        static std::string initialize = "agent::initialize";
        static std::string callback   = "agent::callback";
        static std::string console    = "agent::console";
        static std::string heartbeat  = "agent::heartbeat";
        static std::string status     = "agent::status";
        static std::string remove     = "agent::remove";
        static std::string buildlog   = "agent::build-log";
        static std::string note       = "agent::note";
    }

    namespace listener {
        static std::string add    = "listener::register";
        static std::string start  = "listener::start";
        static std::string edit   = "listener::edit";
        static std::string stop   = "listener::stop";
        static std::string status = "listener::status";
        static std::string log    = "listener::log";
        static std::string remove = "listener::remove";
    }
};


#endif //HAVOCCLIENT_EVENTS_H
