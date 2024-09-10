#include <Common.h>
#include <Havoc.h>

HcSessionGraph::HcSessionGraph(
    QWidget* parent
) : QGraphicsView( parent ) {
    scene = new QGraphicsScene( this );
    scene->setItemIndexMethod( QGraphicsScene::BspTreeIndex );

    setScene( scene );
    setCacheMode( CacheBackground );
    setViewportUpdateMode( BoundingRectViewportUpdate );
    setTransformationAnchor( AnchorUnderMouse );
    setRenderHint( QPainter::Antialiasing );
    scaleView( 1 );

    layout = new QVBoxLayout( this );
    button = new HcSessionGraphButton( this );

    layout->setAlignment( Qt::AlignLeft | Qt::AlignTop );
    layout->addWidget( button );

    server = new HcSessionGraphItem;

    scene->addItem( server );

    server->setGraph( this );
    server->setPos( 100,50 );
}

HcSessionGraph::~HcSessionGraph() {
    delete scene;
}

auto HcSessionGraph::scaleView(
    qreal scaleFactor
) -> void {
    const auto factor = transform().scale(
        scaleFactor,
        scaleFactor
    ).mapRect( QRectF( 0, 0, 1, 1 ) ).width();

    if ( factor < 1 || factor > 50 ) {
        return;
    }

    scale( scaleFactor, scaleFactor );
}

auto HcSessionGraph::addAgent(
    HcAgent *agent
) -> HcSessionGraphItem* {
    const auto item = new HcSessionGraphItem;

    item->setAgent( agent );
    item->setGraph( this );
    item->setItemEdge( new HcSessionGraphEdge( server, item, Havoc->Theme.getGreen() ) );

    scene->addItem( item );
    scene->addItem( item->itemEdge() );
    _nodes.push_back( item );

    return item;
}

auto HcSessionGraph::itemMoved() -> void
{
    if ( ! timer_id ) {
        timer_id = startTimer( 40 );
    }
}

auto HcSessionGraph::nodes() -> std::vector<HcSessionGraphItem*>
{
    return _nodes;
}

auto HcSessionGraph::keyPressEvent(
    QKeyEvent* event
) -> void {
    switch ( event->key() ) {
        case Qt::Key_Plus: {
            zoomIn();
            break;
        }

        case Qt::Key_Minus: {
            zoomOut();
            break;
        }

        default: {
            QGraphicsView::keyPressEvent( event );
        }
    }
}

void HcSessionGraph::resizeEvent(
    QResizeEvent* event
) {
    scene->setSceneRect( 0, 0, event->size().width(), event->size().height() );

    QGraphicsView::resizeEvent( event );
}

void HcSessionGraph::wheelEvent(
    QWheelEvent* event
) {
    if ( QApplication::keyboardModifiers() & Qt::ShiftModifier ) {
        horizontalScrollBar()->event( event );
    } else if ( QApplication::keyboardModifiers() & Qt::ControlModifier ) {
        scaleView( pow( 2., event->angleDelta().y() / 500.0 ) );
    } else {
        verticalScrollBar()->event( event );
    }

    event->ignore();
}

void HcSessionGraph::timerEvent(
    QTimerEvent* event
) {
    bool itemsMoved = false;

    for ( const auto item : nodes() ) {
        item->calculateForces();

        itemsMoved = item->advancePosition();

        item->adjust();
    }

    if ( ! itemsMoved ) {
        killTimer( timer_id );
        timer_id = 0;
    }
}

void HcSessionGraph::zoomIn() {
    scaleView( qreal( 1.2 ) );
}

void HcSessionGraph::zoomOut() {
    scaleView( 1 / qreal( 1.2 ) );
}

//
// HcSessionGraphItem
//

HcSessionGraphItem::HcSessionGraphItem() {
    setFlag( ItemIsMovable );
    setFlag( ItemSendsGeometryChanges );
    setCacheMode( DeviceCoordinateCache );
    setZValue( -1 );

    rect = QRectF( -40, -50, 80, 80 );
}

HcSessionGraphItem::~HcSessionGraphItem() = default;

auto HcSessionGraphItem::setGraph(
    HcSessionGraph* graph
) -> void {
    _graph = graph;
}

auto HcSessionGraphItem::graph(
    void
) const -> HcSessionGraph* {
    return _graph;
}

auto HcSessionGraphItem::setAgent(
    HcAgent* agent
) -> void {
    _agent = agent;
}

auto HcSessionGraphItem::agent(
    void
) const -> HcAgent* {
    return _agent;
}

auto HcSessionGraphItem::addPivot(
    HcSessionGraphItem *agent
) -> void {
    children.push_back( agent );
}

auto HcSessionGraphItem::pivots(
    void
) const -> std::vector<HcSessionGraphItem*> {
    return children;
}

auto HcSessionGraphItem::addParent(
    HcSessionGraphItem* parent
) -> void {
    _parent = parent;
}

auto HcSessionGraphItem::parent(
    void
) const -> std::optional<HcSessionGraphItem*> {
    return _parent;
}

auto HcSessionGraphItem::addEdge(
    HcSessionGraphEdge* edge
) -> void {
    edges.push_back( edge );
    edge->adjust();
}

auto HcSessionGraphItem::adjust() -> void
{
    for ( const auto _edge : edges) {
        _edge->adjust();
    }
}

auto HcSessionGraphItem::setItemEdge(
    HcSessionGraphEdge* edge
) -> void {
    this->edge = edge;
}

auto HcSessionGraphItem::itemEdge(
    void
) -> HcSessionGraphEdge* {
    return edge;
}

auto HcSessionGraphItem::boundingRect(
    void
) const -> QRectF {
    return rect;
}

auto HcSessionGraphItem::paint(
    QPainter* painter,
    const QStyleOptionGraphicsItem* option,
    QWidget* widget
) -> void {
    if ( agent() == nullptr ) {
        //
        // check if the current graph item is the root item
        // in this scene (aka havoc server)
        //
        painter->drawImage( rect, QImage( ":/graph/server" ) );
    } else if ( parent().has_value() ) {
        //
        // check if a parent has been specified. if yes then
        // it means that the current item is a pivot session
        //
        spdlog::debug( "[HcSessionGraphItem::paint] pivot" );
    } else {
        //
        // it's a direct session connection
        // via no pivoting agent
        //

        // TODO: make it available over the scripting engine
        //       to specify what kind of image to render
        painter->drawImage( rect, QImage( ":/graph/win11" ) );
    }
}

auto HcSessionGraphItem::contextMenuEvent(
    QGraphicsSceneContextMenuEvent* event
) -> void {
    //
    // only process context menu if on a agent session
    //
    if ( agent() == nullptr ) {
        return QGraphicsItem::contextMenuEvent( event );
    }

    auto menu    = QMenu();
    auto actions = Havoc->Actions( HavocClient::ActionObject::ActionAgent );

    menu.setStyleSheet( HavocClient::StyleSheet() );

    //
    // if a single selected agent item then try
    // to add the registered actions as well
    //
    menu.addAction( QIcon( ":/icons/16px-agent-console" ), "Interact" );
    menu.addSeparator();

    //
    // add all agent type registered actions
    //

    for ( const auto action : actions ) {
        if ( action->agent.type == agent()->type ) {
            if ( action->icon.empty() ) {
                menu.addAction( action->name.c_str() );
            } else {
                menu.addAction( QIcon( action->icon.c_str() ), action->name.c_str() );
            }
        }
    }

    menu.addSeparator();
    menu.addAction( QIcon( ":/icons/16px-blind-white" ), "Hide" );
    menu.addAction( QIcon( ":/icons/16px-remove" ), "Remove" );

    const auto action = menu.exec( event->screenPos() );
    if ( ! action ) {
        return;
    }

    if ( action->text().compare( "Interact" ) == 0 ) {
        Havoc->Gui->PageAgent->spawnAgentConsole( agent()->uuid );
    } else if ( action->text().compare( "Remove" ) == 0 ) {
        agent()->remove();
    } else if ( action->text().compare( "Hide" ) == 0 ) {
        agent()->hide();
    } else {
        for ( const auto act : actions ) {
            if ( act->name       == action->text().toStdString() &&
                 act->agent.type == agent()->type
            ) {
                if ( agent()->interface.has_value() ) {
                    try {
                        HcPythonAcquire();

                        act->callback( agent()->interface.value() );
                    } catch ( py11::error_already_set& e ) {
                        spdlog::error( "failed to execute action callback: {}", e.what() );
                    }
                }
                return;
            }
        }

        spdlog::debug( "[ERROR] invalid action from selected agent menu" );
    }
}

auto HcSessionGraphItem::itemChange(
    GraphicsItemChange change,
    const QVariant&    value
) -> QVariant {
    switch ( change )
    {
        case ItemPositionChange: {
            adjust();
            graph()->itemMoved();
        }

        default: break;
    }

    return QGraphicsItem::itemChange( change, value );
}

auto HcSessionGraphItem::shape() -> QPainterPath
{
    auto path = QPainterPath();

    path.addEllipse( rect );

    return path;
}

auto HcSessionGraphItem::calculateForces(
    void
) -> void {
    auto xvel = qreal( 0 );
    auto yvel = qreal( 0 );

    if ( ! scene() || scene()->mouseGrabberItem() == this ) {
        _pos = pos();
        return;
    }

    if ( qAbs( xvel ) < 0.1 && qAbs( yvel ) < 0.1 ) {
        xvel = yvel = 0;
    }

    auto sceneRect = scene()->sceneRect();

    _pos = pos() + QPointF( xvel, yvel );
    _pos.setX( qMin( qMax( _pos.x(), sceneRect.left() + 10 ), sceneRect.right()  - 10 ) );
    _pos.setY( qMin( qMax( _pos.y(), sceneRect.top()  + 10 ), sceneRect.bottom() - 10 ) );
}

auto HcSessionGraphItem::advancePosition(
    void
) -> bool {
    if ( _pos == pos() ) {
        return false;
    }

    setPos( _pos );

    return true;
}

// auto HcSessionGraphItem::shuffle(
//     void
// ) const -> void {
//     for ( const auto node : graph()->nodes() ) {
//         node->item->setPos( -150 + QRandomGenerator::global()->bounded( 300 ), -150 + QRandomGenerator::global()->bounded( 300 ) );
//     }
// }

auto HcSessionGraphItem::mousePressEvent(
    QGraphicsSceneMouseEvent* event
) -> void {
    update();
    graph()->itemMoved();
    QGraphicsItem::mousePressEvent( event );
}

auto HcSessionGraphItem::mouseReleaseEvent(
    QGraphicsSceneMouseEvent* event
) -> void {
    update();
    graph()->itemMoved();
    QGraphicsItem::mouseReleaseEvent( event );
}

auto HcSessionGraphItem::mouseMoveEvent(
    QGraphicsSceneMouseEvent* event
) -> void {
    update();
    graph()->itemMoved();
    QGraphicsItem::mouseMoveEvent( event );
}

//
// HcSessionGraphEdge
//

HcSessionGraphEdge::HcSessionGraphEdge(
    HcSessionGraphItem* source,
    HcSessionGraphItem* destination,
    QColor              color
) : _source( source ), _destination( destination ), color( color ) {
    setAcceptedMouseButtons( Qt::NoButton );

    source->addEdge( this );
    destination->addEdge( this );

    adjust();
}

auto HcSessionGraphEdge::source() const -> HcSessionGraphItem *
{
    return _source;
}

auto HcSessionGraphEdge::destination() const -> HcSessionGraphItem *
{
    return _destination;
}

auto HcSessionGraphEdge::adjust() -> void
{
    if ( ! source() || ! destination() )
        return;

    auto line   = QLineF( mapFromItem( source(), 0, 0 ), mapFromItem( destination(), 0, 0 ) );
    auto length = line.length();

    prepareGeometryChange();

    if ( length > qreal( 20. ) ) {
        auto edgeSpace  = 50;
        auto edgeOffset = QPointF( ( line.dx() * edgeSpace ) / length, ( line.dy() * edgeSpace ) / length );

        sourcePoint = line.p1() + edgeOffset;
        destPoint   = line.p2() - edgeOffset;
    } else {
        sourcePoint = destPoint = line.p1();
    }
}

auto HcSessionGraphEdge::setColor(
    const QColor& color
) -> void {
    this->color = color;
}

QRectF HcSessionGraphEdge::boundingRect() const {
    const auto width = qreal( 1 );
    const auto extra = qreal( width + arrowSize ) / 2.0;

    if ( ! source() || ! destination() ) {
        return QRectF();
    }

    return QRectF( sourcePoint, QSizeF(
        destPoint.x() - sourcePoint.x(),
        destPoint.y() - sourcePoint.y())
    ).normalized().adjusted( -extra, -extra, extra, extra );
}

void HcSessionGraphEdge::paint(
    QPainter*                       painter,
    const QStyleOptionGraphicsItem* option,
    QWidget*                        widget
) {
    auto line = QLineF( sourcePoint, destPoint );

    if ( ! source() || ! destination() ) {
        return;
    }

    if ( qFuzzyCompare( line.length(), qreal( 0. ) ) ) {
        return;
    }

    auto angle         = std::atan2( -line.dy(), line.dx() );
    auto sourceArrowP1 = sourcePoint + QPointF( sin( angle + M_PI / 3 ) * arrowSize, cos( angle + M_PI / 3 ) * arrowSize );
    auto sourceArrowP2 = sourcePoint + QPointF( sin( angle + M_PI - M_PI / 3 ) * arrowSize, cos( angle + M_PI - M_PI / 3 ) * arrowSize );
    auto destArrowP1   = destPoint   + QPointF( sin( angle - M_PI / 3 ) * arrowSize, cos( angle - M_PI / 3 ) * arrowSize );
    auto destArrowP2   = destPoint   + QPointF( sin( angle - M_PI + M_PI / 3 ) * arrowSize, cos( angle - M_PI + M_PI / 3 ) * arrowSize );

    // Draw the line itself
    painter->setPen( QPen( color, 2, Qt::SolidLine, Qt::RoundCap, Qt::RoundJoin ) );
    painter->drawLine( line );
    painter->setBrush( color );

    if ( source()->agent() == nullptr ) {
        painter->drawPolygon( QPolygonF() << line.p1() << sourceArrowP1 << sourceArrowP2 );
    } else {
        painter->drawPolygon( QPolygonF() << line.p2() << destArrowP1 << destArrowP2 );
    }
}

//
// HcSessionGraphButton
//

HcSessionGraphButton::HcSessionGraphButton(
    QWidget* parent
) : QToolButton( parent ) {
    menu.addAction( "show domain" )->setCheckable( true );
    menu.addAction( "show direct only" )->setCheckable( true );
    menu.addAction( "show elevated only" )->setCheckable( true );
    menu.setStyleSheet( HavocClient::StyleSheet() );

    setAttribute( Qt::WA_TranslucentBackground );
    setStyleSheet( "background: transparent;" );

    setMenu( &menu );
    setToolButtonStyle( Qt::ToolButtonIconOnly );
    setIcon( QIcon( ":/icons/32px-align" ) );
    setProperty( "HcButton", "true" );
    setPopupMode( InstantPopup );
    setLayoutDirection( Qt::LeftToRight );
    setMaximumWidth( 90 );
}

HcSessionGraphButton::~HcSessionGraphButton() = default;
