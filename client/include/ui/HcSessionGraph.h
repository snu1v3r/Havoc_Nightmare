#ifndef HCSESSIONGRAPH_H
#define HCSESSIONGRAPH_H

#include <QGraphicsItem>
#include <QGraphicsView>
#include <QKeyEvent>
#include <QScrollBar>
#include <QRect>

class HcSessionGraph : public QGraphicsView
{
    Q_OBJECT

    typedef struct {
        std::string name;
    } NodeEntry;

    QGraphicsScene*         scene  = {};
    NodeEntry*              server = {};
    std::vector<NodeEntry*> nodes  = {};

public:
    explicit HcSessionGraph( QWidget *parent = nullptr );
    ~HcSessionGraph() override;

    auto scaleView(
        qreal scaleFactor
    ) -> void;

protected:
    void keyPressEvent( QKeyEvent* event ) override;
    void resizeEvent( QResizeEvent* event ) override;
    void wheelEvent( QWheelEvent* event ) override;

public Q_SLOTS:
    auto zoomIn() -> void;
    auto zoomOut() -> void;
};

#endif //HCSESSIONGRAPH_H
