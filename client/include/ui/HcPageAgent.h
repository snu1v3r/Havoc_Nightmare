#ifndef HAVOCCLIENT_HCPAGEAGENT_H
#define HAVOCCLIENT_HCPAGEAGENT_H

#include <QStackedWidget>
#include <QtWidgets/QComboBox>
#include <QtWidgets/QSplitter>
#include <QtWidgets/QTabWidget>
#include <QtWidgets/QTableWidget>
#include <QtWidgets/QWidget>

#include <ui/HcConsole.h>
#include <ui/HcSessionGraph.h>
#include <core/HcAgent.h>

#include <DockManager.h>
#include <DockWidget.h>

QT_BEGIN_NAMESPACE

class HcAgentConsole : public HcConsole {
    HcAgent* Meta = nullptr;

public:
    explicit HcAgentConsole(
        HcAgent* meta,
        QWidget* parent = nullptr
    );

    auto inputEnter() -> void override;
};

class HcPageAgent : public QWidget
{
    bool                  SplitterMoveToggle  = false;
    ads::CDockAreaWidget* ConsoleAreaWidget   = {};

public:
    ads::CDockManager* DockManager            = {};
    QGridLayout*       gridLayout             = {};
    QSplitter*         Splitter               = {};
    QStackedWidget*    StackedWidget          = {};
    QTableWidget*      AgentTable             = {};
    HcSessionGraph*    AgentGraph             = {};
    QLabel*            AgentDisplayerElevated = {};
    QLabel*            AgentDisplayerSessions = {};
    QLabel*            AgentDisplayerTargets  = {};
    QLabel*            AgentDisplayerPivots   = {};
    QToolButton*       AgentActionButton      = {};
    QMenu*             AgentActionMenu        = {};
    QAction*           ActionPayload          = {};
    QAction*           ActionShowHidden       = {};
    QSpacerItem*       horizontalSpacer       = {};

    QTableWidgetItem*  TitleAgentID           = {};
    QTableWidgetItem*  TitleInternal          = {};
    QTableWidgetItem*  TitleUsername          = {};
    QTableWidgetItem*  TitleHostname          = {};
    QTableWidgetItem*  TitleSystem            = {};
    QTableWidgetItem*  TitleProcessID         = {};
    QTableWidgetItem*  TitleProcessName       = {};
    QTableWidgetItem*  TitleProcessArch       = {};
    QTableWidgetItem*  TitleThreadID          = {};
    QTableWidgetItem*  TitleNote              = {};
    QTableWidgetItem*  TitleLastCallback      = {};

    std::vector<HcAgent*> agents = {};

    explicit HcPageAgent(QWidget* parent = nullptr );
    ~HcPageAgent();

    auto retranslateUi() -> void;

    auto addTab(
        const QString& name,
        QWidget*       widget
    ) -> void;

    auto addAgent(
        const json& metadata
    ) -> void;

    auto removeAgent(
        const std::string& uuid
    ) -> void;

    auto spawnAgentConsole(
        const std::string& uuid
    ) -> void;

    auto Agent(
        const std::string& uuid
    ) -> std::optional<HcAgent*>;

    auto AgentConsole(
        const std::string& uuid,
        const std::string& format,
        const std::string& output = ""
    ) -> void;

private Q_SLOTS:
    auto handleAgentMenu(
        const QPoint& pos
    ) -> void;

    auto handleAgentDoubleClick(
        const QModelIndex& index
    ) -> void;

    auto itemChanged(
        QTableWidgetItem *item
    ) -> void;

    auto actionShowHidden(
        bool checked
    ) -> void;

    auto actionPayloadBuilder(
        bool checked
    ) -> void;

    auto actionTriggered(
        QAction* triggered
    ) -> void;
};

QT_END_NAMESPACE

#endif //HAVOCCLIENT_HCPAGEAGENT_H
