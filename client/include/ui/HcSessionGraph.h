#ifndef HCSESSIONGRAPH_H
#define HCSESSIONGRAPH_H

#include <core/HcAgent.h>

#include <QGraphicsItem>
#include <QGraphicsView>
#include <QGraphicsProxyWidget>
#include <QGraphicsSceneContextMenuEvent>
#include <QPainterPath>
#include <QKeyEvent>
#include <QScrollBar>
#include <QRect>

class HcSessionGraphButton;
class HcSessionGraphItem;
class HcSessionGraphEdge;
class HcSessionGraph;

//
// HcSessionGraphItem
//

class HcSessionGraphItem final : public QGraphicsItem
{
    HcAgent*                           _agent   = {};
    std::optional<HcSessionGraphItem*> _parent  = {};
    HcSessionGraph*                    _graph   = {};
    QPointF                            _pos     = {};
    HcSessionGraphEdge*                edge     = {};
    std::vector<HcSessionGraphItem*>   children = {};
    QRectF                             rect     = {};

    std::vector<HcSessionGraphEdge*> edges;
public:
    explicit HcSessionGraphItem();
    ~HcSessionGraphItem() override;

    auto setGraph(
        HcSessionGraph* graph
    ) -> void;

    auto graph(
        void
    ) const -> HcSessionGraph*;

    auto setAgent(
        HcAgent* agent
    ) -> void;

    auto agent(
        void
    ) const -> HcAgent*;

    auto addPivot(
        HcSessionGraphItem* agent
    ) -> void;

    auto pivots(
        void
    ) const -> std::vector<HcSessionGraphItem*>;

    auto addParent(
        HcSessionGraphItem* parent
    ) -> void;

    auto parent(
        void
    ) const -> std::optional<HcSessionGraphItem*>;

    auto addEdge(
        HcSessionGraphEdge* edge
    ) -> void;

    auto adjust() -> void;

    auto setItemEdge(
        HcSessionGraphEdge* edge
    ) -> void;

    auto itemEdge() -> HcSessionGraphEdge*;

    auto calculateForces() -> void;
    auto advancePosition() -> bool;

// public Q_SLOTS:
//     auto shuffle() const -> void;

protected:
    auto boundingRect(
        void
    ) const -> QRectF override;

    auto paint(
        QPainter* painter,
        const QStyleOptionGraphicsItem* option,
        QWidget* widget
    ) -> void override;

    auto contextMenuEvent(
        QGraphicsSceneContextMenuEvent* event
    ) -> void override;

    auto itemChange(
        GraphicsItemChange change,
        const QVariant&    value
    ) -> QVariant override;

    auto shape() -> QPainterPath;
    auto mousePressEvent( QGraphicsSceneMouseEvent* event ) -> void override;
    auto mouseReleaseEvent( QGraphicsSceneMouseEvent* event ) -> void override;
    auto mouseMoveEvent( QGraphicsSceneMouseEvent* event ) -> void override;
};

//
// HcSessionGraphEdge
//

class HcSessionGraphEdge final : public QGraphicsItem
{
    HcSessionGraphItem* _source      = {};
    HcSessionGraphItem* _destination = {};
    QColor              color        = {};
    QPointF             sourcePoint  = {};
    QPointF             destPoint    = {};
    qreal               arrowSize    = 10;

public:
    explicit HcSessionGraphEdge(
        HcSessionGraphItem* source,
        HcSessionGraphItem* destination,
        QColor              color
    );

    auto source() const -> HcSessionGraphItem*;
    auto destination() const -> HcSessionGraphItem*;

    auto adjust() -> void;
    auto setColor(
        const QColor& color
    ) -> void;

    enum { Type = UserType + 2 };
    int type() const override { return Type; }

protected:
    QRectF boundingRect() const override;
    void paint( QPainter* painter, const QStyleOptionGraphicsItem* option, QWidget* widget ) override;
};

class HcSessionGraphButton final : public QToolButton
{
    Q_OBJECT

    QMenu menu;

public:
    explicit HcSessionGraphButton( QWidget* parent = nullptr );
    ~HcSessionGraphButton() override;
};

//
// HcSessionGraph
//

class HcSessionGraph final : public QGraphicsView
{
    Q_OBJECT

    std::vector<HcSessionGraphItem*> _nodes   = {};
    QGraphicsScene*                  scene    = {};
    HcSessionGraphItem*              server   = {};
    HcSessionGraphButton*            button   = {};
    QVBoxLayout*                     layout   = {};
    int                              timer_id = 0;

public:
    explicit HcSessionGraph( QWidget *parent = nullptr );
    ~HcSessionGraph() override;

    auto scaleView(
        qreal scaleFactor
    ) -> void;

    auto addAgent(
        HcAgent* agent
    ) -> HcSessionGraphItem*;

    auto itemMoved() -> void;

    auto nodes() -> std::vector<HcSessionGraphItem*>;

protected:
    void keyPressEvent( QKeyEvent* event ) override;
    void resizeEvent( QResizeEvent* event ) override;
    void wheelEvent( QWheelEvent* event ) override;
    void timerEvent( QTimerEvent* event ) override;

public Q_SLOTS:
    auto zoomIn() -> void;
    auto zoomOut() -> void;
};

#endif //HCSESSIONGRAPH_H
