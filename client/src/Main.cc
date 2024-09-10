#include <Havoc.h>

HavocClient* Havoc = nullptr;

auto main(
    int    argc,
    char** argv
) -> int {
    auto Application = QApplication( argc, argv );

    Havoc = new HavocClient;
    Havoc->Main( argc, argv );

    QApplication::exec();
}