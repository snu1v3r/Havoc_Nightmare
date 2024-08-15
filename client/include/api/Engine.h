#ifndef HAVOCCLIENT_API_ENGINE_H
#define HAVOCCLIENT_API_ENGINE_H

#include <Common.h>

#include <api/HcCore.h>
#include <api/HcScriptManager.h>
#include <api/HcAgent.h>
#include <api/HcListener.h>

inline auto HcPythonReleaseGil() -> void {
    //
    // NOTE: I know. this is insanely stupid and retarded to do. But it has gotten to
    //       a point where I stopped caring as long as shit works. And shit does work.
    //       So who am I to complain. Tho I will try to rewrite this mostly in the
    //       future to have a better implementation to safely execute python code
    //       in other threads without the gil deadlocking me. Or I am waiting til python
    //       GIL has been removed from the interpreter and the api.
    //
    if ( PyGILState_Check() ) {
        PyEval_SaveThread();
    }
}

class HcPyEngine : public QThread {
    Q_OBJECT

public:
    py11::scoped_interpreter* guard{};

    std::optional<py11::object> PyEval = {};
    std::optional<py11::object> PyLoad = {};

    explicit HcPyEngine();
    ~HcPyEngine();

    auto run() -> void;
};

#endif
