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

    setBackgroundBrush( QBrush( Havoc->Theme.getBackground() ) );
}

HcSessionGraph::~HcSessionGraph() {
    delete scene;
}

auto HcSessionGraph::scaleView(
    qreal scaleFactor
) -> void {
    auto factor = transform().scale(
        scaleFactor,
        scaleFactor
    ).mapRect( QRectF( 0, 0, 1, 1 ) ).width();

    if ( factor < 1 || factor > 50 ) {
        return;
    }

    scale( scaleFactor, scaleFactor );
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
        // zoom in - zoom out
        scaleView( pow( 2., event->angleDelta().y() / 500.0 ) );
    } else {
        // vertical scroll
        verticalScrollBar()->event( event );
    }

    event->ignore();
}

void HcSessionGraph::zoomIn()
{
    scaleView( qreal( 1.2 ) );
}

void HcSessionGraph::zoomOut()
{
    scaleView( 1 / qreal( 1.2 ) );
}

