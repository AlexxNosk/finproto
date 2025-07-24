#include "connect.h"
#include "qt_app.h"
int main(int argc, char *argv[]){
    const std::string path = "/home/alex/pro/finproto/.env";
    loadEnvFile(path);
    QApplication app(argc, argv);
    QMainWindow window;

    QWidget *mainWidget = new QWidget;
    QHBoxLayout *mainLayout = new QHBoxLayout(mainWidget);

    QCustomPlot *customPlot = new QCustomPlot;
    QLineEdit *tableNameEdit;

    QWidget *sidePanel = createSidePanel(tableNameEdit, customPlot);

    mainLayout->addWidget(sidePanel);
    mainLayout->addWidget(customPlot, 1);

    window.setCentralWidget(mainWidget);
    window.resize(1200, 700);
    window.setWindowTitle("Candlestick Viewer");

    plotData(customPlot, tableNameEdit->text());

    window.show();
    return app.exec();
    return 0;
}

