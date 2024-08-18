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

    QTimer* HeartbeatTimer = nullptr;

public:
    explicit HcMetaWorker();
    ~HcMetaWorker();

    /* run event thread */
    void run();

    auto listeners() -> void;
    auto agents() -> void;

Q_SIGNALS:
    auto AddListener(
        const json& listener
    )-> void;

    auto AddAgent(
        const json& agent
    )-> void;

    auto eventWorkerRun() -> void;
};

#endif //HAVOCCLIENT_HCMETAWORKER_H
