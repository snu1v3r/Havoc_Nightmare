#include <Havoc.h>
#include <core/HcEventWorker.h>

HcEventWorker::HcEventWorker()  = default;
HcEventWorker::~HcEventWorker() = default;

void HcEventWorker::run() {
    WebSocket    = new QWebSocket;
    auto Host    = QString( Havoc->Server().c_str() );
    auto SslConf = WebSocket->sslConfiguration();

    /* ignore annoying SSL errors */
    SslConf.setPeerVerifyMode( QSslSocket::VerifyNone );
    WebSocket->setSslConfiguration( SslConf );
    WebSocket->ignoreSslErrors();

    QObject::connect( WebSocket, &QWebSocket::connected,             this, &HcEventWorker::connected );
    QObject::connect( WebSocket, &QWebSocket::disconnected,          this, &HcEventWorker::closed );
    QObject::connect( WebSocket, &QWebSocket::binaryMessageReceived, this, &HcEventWorker::receivedEvent );

    spdlog::debug( "HcEventWorker::run" );

    WebSocket->open( "wss://" + Host + "/api/event" );
}

auto HcEventWorker::connected() -> void {
    /* show event progress dialog or something */
    auto login = json {
        { "token", Havoc->Token() }
    };

    WebSocket->sendBinaryMessage( login.dump().c_str() );
}

auto HcEventWorker::closed() -> void {
    /* this is it. clean everything up and close the Nocturn client */
    emit socketClosed( );
}

auto HcEventWorker::receivedEvent(
    const QByteArray &message
) -> void {
    /* send json event */
    emit availableEvent( message );
}