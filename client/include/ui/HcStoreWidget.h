#ifndef HAVOCCLIENT_HCSTOREWIDGET_H
#define HAVOCCLIENT_HCSTOREWIDGET_H

#include <QtCore/QVariant>
#include <QtWidgets/QApplication>
#include <QtWidgets/QGridLayout>
#include <QtWidgets/QLabel>
#include <QtWidgets/QLineEdit>
#include <QtWidgets/QListWidget>
#include <QtWidgets/QPushButton>
#include <QtWidgets/QSpacerItem>
#include <QtWidgets/QSplitter>
#include <QtWidgets/QTextEdit>
#include <QtWidgets/QWidget>
#include <QtWidgets/QStackedWidget>

#include <ui/HcLineEdit.h>
#include <core/HcStorePluginWorker.h>

class HcMarketPluginItem : public QWidget
{
public:
    QGridLayout* gridLayout       = { 0 };
    QLabel*      LabelDescription = { 0 };
    QLabel*      LabelName        = { 0 };
    QLabel*      LabelGroup       = { 0 };
    QSpacerItem* Spacer           = { 0 };
    QGridLayout* InstallLayout    = { 0 };
    QLabel*      LabelInstalled   = { 0 };

    explicit HcMarketPluginItem(
        const QString& name,
        const QString& description,
        const QString& group,
        bool           pro,
        QWidget*       parent
    );

    void setInstalled();

    void setNotInstalled();
};

class HcStoreWidget : public QWidget {
    Q_OBJECT

private:
    QGridLayout*             gridLayout_3;
    QSplitter*               splitter;
    QGridLayout*             gridLayout;
    QWidget*                 MarketPlaceWidget;
    QLineEdit*               MarketPlaceSearch;
    QListWidget*             MarketPlaceList;
    QStackedWidget*          PluginViewStack;
    std::vector<PluginView*> Plugins;

public:
    struct {
        QThread*             Thread;
        HcStorePluginWorker* Worker;
    } PluginWorker;

    explicit HcStoreWidget( QWidget* parent = nullptr );

    auto AddPlugin(
        const std::string& parent,
        const std::string& repo,
        PluginView*        plugin
    ) -> void;

    auto AddPluginToMarketList(
        PluginView*        plugin,
        const std::string& group = "local"
    ) -> void;

    auto QueryPluginMarket(
        const QString& query
    ) -> void;

    auto PluginQueryContainMeta(
        const PluginView* plugin,
        const QString&    query
    ) -> bool;

    static auto HttpGet(
        const std::string& url,
        const std::string& authorization = ""
    ) -> std::optional<std::string>;

signals:
    auto RegisterRepository(
        const std::string&              repository,
        const std::vector<std::string>& plugins,
        const std::string&              credentials
    ) -> void;

    auto RegisterPlugin(
        const std::string& directory,
        const json&        object
    ) -> void;

    auto PluginInstall(
        PluginView* plugin
    ) -> void;
};

#endif //HAVOCCLIENT_HCSTOREWIDGET_H
