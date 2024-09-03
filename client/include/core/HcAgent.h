#ifndef HAVOCCLIENT_HCAGENT_H
#define HAVOCCLIENT_HCAGENT_H

class HcAgentConsole;
struct HcAgent;

#include <Common.h>
#include <ui/HcPageAgent.h>

class HcAgentEmit : public QWidget {
    Q_OBJECT

signals:
    auto ConsoleWrite(
        const QString& uuid,
        const QString& text
    ) -> void;

    auto HeartBeat(
        const QString& uuid,
        const QString& time
    ) -> void;
};

struct HcAgent;

class HcAgentTableItem : public QTableWidgetItem {

public:
    HcAgent* agent  = {};
    bool     ignore = true;

    explicit HcAgentTableItem(
        const QString&          value,
        const Qt::ItemFlag      flags = Qt::ItemIsEditable,
        const Qt::AlignmentFlag align = Qt::AlignCenter
    );
};

struct HcAgent {
    std::string                 uuid;
    std::string                 type;
    json                        data;
    std::optional<py11::object> interface;
    HcAgentConsole*             console;
    QString                     last;

    struct {
        HcAgentTableItem* Uuid;
        HcAgentTableItem* Internal;
        HcAgentTableItem* Username;
        HcAgentTableItem* Hostname;
        HcAgentTableItem* ProcessPath;
        HcAgentTableItem* ProcessName;
        HcAgentTableItem* ProcessId;
        HcAgentTableItem* ThreadId;
        HcAgentTableItem* Arch;
        HcAgentTableItem* System;
        HcAgentTableItem* Note;
        HcAgentTableItem* Last;
    } ui;

    //
    // signals
    //
    HcAgentEmit emitter = {};

    auto remove() -> void;
    auto hide() -> void;
};

#endif //HAVOCCLIENT_HCAGENT_H
