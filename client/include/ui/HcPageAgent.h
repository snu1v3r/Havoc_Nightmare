#ifndef HAVOCCLIENT_HCPAGEAGENT_H
#define HAVOCCLIENT_HCPAGEAGENT_H

#include <QtCore/QVariant>
#include <QtWidgets/QApplication>
#include <QtWidgets/QComboBox>
#include <QtWidgets/QGridLayout>
#include <QtWidgets/QHeaderView>
#include <QtWidgets/QLabel>
#include <QtWidgets/QSpacerItem>
#include <QtWidgets/QSplitter>
#include <QtWidgets/QTabWidget>
#include <QtWidgets/QTableWidget>
#include <QtWidgets/QWidget>
#include <QtWidgets/QPushButton>

#include <ui/HcConsole.h>
#include <core/HcAgent.h>

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
    bool SplitterMoveToggle = false;

public:
    QGridLayout*      gridLayout             = nullptr;
    QComboBox*        ComboAgentView         = nullptr;
    QSplitter*        Splitter               = nullptr;
    QTableWidget*     AgentTable             = nullptr;
    QTabWidget*       AgentTab               = nullptr;
    QLabel*           AgentDisplayerElevated = nullptr;
    QLabel*           AgentDisplayerSessions = nullptr;
    QLabel*           AgentDisplayerTargets  = nullptr;
    QLabel*           AgentDisplayerPivots   = nullptr;
    QToolButton*      AgentActionButton      = nullptr;
    QMenu*            AgentActionMenu        = nullptr;
    QAction*          ActionPayload          = nullptr;
    QAction*          ActionShowHidden       = nullptr;
    QSpacerItem*      horizontalSpacer       = nullptr;
    QTableWidgetItem* TitleAgentID           = nullptr;
    QTableWidgetItem* TitleInternal          = nullptr;
    QTableWidgetItem* TitleUsername          = nullptr;
    QTableWidgetItem* TitleHostname          = nullptr;
    QTableWidgetItem* TitleSystem            = nullptr;
    QTableWidgetItem* TitleProcessID         = nullptr;
    QTableWidgetItem* TitleProcessName       = nullptr;
    QTableWidgetItem* TitleProcessArch       = nullptr;
    QTableWidgetItem* TitleThreadID          = nullptr;
    QTableWidgetItem* TitleNote              = nullptr;
    QTableWidgetItem* TitleLastCallback      = nullptr;

    std::vector<HcAgent*> agents = {};

    explicit HcPageAgent(QWidget* parent = nullptr );
    ~HcPageAgent();

    auto retranslateUi() -> void;

    auto addTab(
        const QString& name,
        QWidget*       widget
    ) const -> void;

    auto addAgent(
        const json& metadata
    ) -> void;

    auto removeAgent(
        const std::string& uuid
    ) -> void;

    auto spawnAgentConsole(
        const std::string& uuid
    ) -> void;

    auto handleAgentMenu(
        const QPoint& pos
    ) -> void;

    auto handleAgentDoubleClick(
        const QModelIndex& index
    ) -> void;

    auto tabCloseRequested(
        int index
    ) const -> void;

    auto Agent(
        const std::string& uuid
    ) -> std::optional<HcAgent*>;

    auto AgentConsole(
        const std::string& uuid,
        const std::string& format,
        const std::string& output = ""
    ) -> void;

    auto itemChanged(
        QTableWidgetItem *item
    ) -> void;

private:
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
