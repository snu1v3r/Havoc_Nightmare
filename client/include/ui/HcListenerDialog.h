#ifndef HAVOCCLIENT_HCLISTENERDIALOG_H
#define HAVOCCLIENT_HCLISTENERDIALOG_H

#include <Common.h>
#include <ui/HcLineEdit.h>
#include <ui/HcComboBox.h>

#include <QtCore/QVariant>
#include <QtWidgets/QApplication>
#include <QtWidgets/QComboBox>
#include <QtWidgets/QDialog>
#include <QtWidgets/QGridLayout>
#include <QtWidgets/QGroupBox>
#include <QtWidgets/QLabel>
#include <QtWidgets/QLineEdit>
#include <QtWidgets/QStackedWidget>

QT_BEGIN_NAMESPACE

struct Protocol {
    std::string  name;
    QWidget*     widget;
    py11::object instance;
};

enum ListenerState {
    None,
    Saved,
    Closed,
    Error
};

class HcListenerDialog : public QDialog
{
    ListenerState   State            = None;
    json            edited_config    = {};

public:
    QGridLayout*    gridLayout       = nullptr;
    HcLineEdit*     InputName        = nullptr;
    HcComboBox*     ComboProtocol    = nullptr;
    QStackedWidget* StackedProtocols = nullptr;
    QPushButton*    ButtonSave       = nullptr;
    QPushButton*    ButtonClose      = nullptr;
    QSpacerItem*    horizontal[ 4 ]  = {};

    std::vector<Protocol> Protocols  = {};

    explicit HcListenerDialog(
        const QString& editing = ""
    );

    ~HcListenerDialog();

    auto getCurrentProtocol() -> Protocol;
    auto getCloseState() -> ListenerState;
    auto save() -> void;

    auto setEditingListener(
        const QString& name,
        const json&    config
    ) -> void;

    auto changeProtocol(
        const QString& text
    ) -> void;

    auto AddProtocol(
        const std::string&  name,
        const py11::object& object
    ) -> void;

    auto start() -> void;
    auto retranslateUi() -> void;
};

QT_END_NAMESPACE

#endif //HAVOCCLIENT_HCLISTENERDIALOG_H
