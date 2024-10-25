#ifndef HAVOCCLIENT_HCMETAWORKER_H
#define HAVOCCLIENT_HCMETAWORKER_H

#include <Common.h>

#include <QThread>
#include <QString>

//
// meta worker is a thread that handles the pulling
// of metadata, listeners, agent sessions, etc.
//

class HcMetaWorker : public QThread {
    Q_OBJECT

    bool m_plugin = false;

public:
    explicit HcMetaWorker(
        bool plugin_worker = false
    );

    ~HcMetaWorker();

    /* run event thread */
    void run();

    auto listeners(
        void
    ) -> void;

    auto agents(
        void
    ) -> void;

    auto plugins(
        void
    ) -> void;

    auto console(
        const std::string& uuid
    ) -> void;

    auto resource(
        const std::string& name,
        const std::string& version,
        const std::string& resource
    ) -> bool;


Q_SIGNALS:
    auto Finished() -> bool;

    auto AddListener(
        const json& listener
    )-> void;

    auto AddAgent(
        const json& agent
    )-> void;

    auto AddAgentConsole(
        const json& data
    )-> void;
};

#endif //HAVOCCLIENT_HCMETAWORKER_H
