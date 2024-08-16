import threading

##
## import havoc specific libs
##
from _pyhavoc import core


def HcListenerProtocolData(
    protocol: str
) -> dict:
    return core.HcListenerProtocolData( protocol )


def HcListenerAll() -> list[str]:
    return core.HcListenerAll()


def HcListenerQueryType( name: str ) -> str:
    return core.HcListenerQueryType( name )


def HcListenerRegisterMenuAction(
    type: str,
    name: str,
    callback,
    icon: str = ""
):
    def _register( function ):
        core.HcListenerRegisterMenuAction( type, name, callback, icon )

    return _register


def HcRegisterMenuAction(
    name: str,
    icon_path: str = "",
):
    def _register( function ):
        core.HcRegisterMenuAction( name, icon_path, function )

    return _register
