cmake_minimum_required( VERSION 3.20 )
project( HavocClient )

##
## cmake sets
##
set( BINARY_NAME Havoc )
set( CMAKE_CXX_STANDARD 20 )
set( CMAKE_AUTOMOC ON )
set( CMAKE_AUTORCC ON )
set( CMAKE_AUTOUIC ON )
set( QT_VERSION 6 )
set( REQUIRED_LIBS Core Gui Widgets Network WebSockets Designer Core5Compat )
set( REQUIRED_LIBS_QUALIFIED Qt6::Core Qt6::Gui Qt6::Widgets Qt6::Network Qt6::WebSockets Qt6::Designer Qt6::Core5Compat )
set( CMAKE_RUNTIME_OUTPUT_DIRECTORY .. )
set( APP_ICON_RESOURCE_WINDOWS "data/Havoc.qrc" )

##
## compiler definition
##
add_definitions( -DQT_NO_DEBUG_OUTPUT )

##
## compiler flags
##
set( CMAKE_CXX_FLAGS "-Wattributes -fvisibility-inlines-hidden -lssl -lcrypto" )

##
## include header directories
##
include_directories( include )
include_directories( extern/pybind11/include )
include_directories( extern/pybind11_json/include )
include_directories( extern/json/include )
include_directories( extern/spdlog/include )
include_directories( extern/toml )
include_directories( extern/http )

##
## Havoc include
##
set( HAVOC_INCLUDE
        include/Common.h
        include/Havoc.h
        include/Events.h

        include/core/HcHelper.h
        include/core/HcEventWorker.h
        include/core/HcAgent.h
        include/core/HcHeartbeatWorker.h
        include/core/HcStorePluginWorker.h
        include/core/HcMetaWorker.h

        include/api/Engine.h
        include/api/HcCore.h
        include/api/HcScriptManager.h
        include/api/HcAgent.h
        include/api/HcListener.h

        include/ui/HcConnectDialog.h
        include/ui/HcMainWindow.h
        include/ui/HcListenerDialog.h

        include/ui/HcPageAgent.h
        include/ui/HcDialogBuilder.h
        include/ui/HcPageListener.h
        include/ui/HcPageScript.h

        include/ui/HcWidget.h
        include/ui/HcConsole.h
        include/ui/HcLineEdit.h
        include/ui/HcComboBox.h
        include/ui/HcTheme.h
        include/ui/HcStoreWidget.h
)

##
## Havoc src
##
set( HAVOC_SRC
        src/Main.cc
        src/Havoc.cc
        src/Events.cc

        src/core/HcHelper.cc
        src/core/HcEventWorker.cc
        src/core/HcAgent.cc
        src/core/HcHeartbeatWorker.cc
        src/core/HcStorePluginWorker.cc
        src/core/HcMetaWorker.cc

        src/ui/HcConnectDialog.cc
        src/ui/HcMainWindow.cc
        src/ui/HcListenerDialog.cc

        src/ui/HcPageAgent.cc
        src/ui/HcPageListener.cc
        src/ui/HcDialogBuilder.cc
        src/ui/HcPageScript.cc

        src/ui/HcWidget.cc
        src/ui/HcConsole.cc
        src/ui/HcLineEdit.cc
        src/ui/HcComboBox.cc
        src/ui/HcTheme.cc
        src/ui/HcStoreWidget.cc

        src/api/HcEngine.cc
        src/api/HcCore.cc
        src/api/HcScriptManager.cc
        src/api/HcAgent.cc
        src/api/HcListener.cc
)

find_package( OpenSSL REQUIRED )
find_package( PythonLibs REQUIRED )

##
## Linking
##

##
## python finding and handling
##
find_package( Python 3 COMPONENTS Interpreter Development REQUIRED )
set( PYTHON_MAJOR $ENV{Python_VERSION_MAJOR} )
set( PYTHON_MINOR $ENV{Python_VERSION_MINOR} )
set( PYTHONLIBS_VERSION_STRING ${Python_VERSION} )
set( PYTHON_INCLUDE_DIR ${Python_INCLUDE_DIRS} )
set( PYTHON_LIBRARIES ${Python_LIBRARIES} )
include_directories( ${PYTHON_INCLUDE_DIRS} )

##
## qt finding and handling
##
find_package( Qt${QT_VERSION} COMPONENTS ${REQUIRED_LIBS} REQUIRED )

## 
## build the executable 
## 
add_executable( ${BINARY_NAME} ${HAVOC_INCLUDE} ${HAVOC_SRC} ${APP_ICON_RESOURCE_WINDOWS} )

##
## linking external libraries
##
target_link_libraries( ${BINARY_NAME} PRIVATE 
  ${REQUIRED_LIBS_QUALIFIED} 
  ${PYTHON_LIBRARIES} 
  OpenSSL::Crypto
)
