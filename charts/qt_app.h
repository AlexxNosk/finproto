#include <QApplication>
#include <QMainWindow>
#include <QWidget>
#include <QLineEdit>
#include <QPushButton>
#include <QHBoxLayout>
#include <QVBoxLayout>
#include <QVector>
#include <QtSql>
#include "qcustomplot.h"

void loadDataFromPostgreSQL(const QString& tableName);
void adjustYAxisToVisibleRange(QCPFinancial *candlesticks, QCPBars *volumeBars,
                                QCPAxis *xAxis, QCPAxis *yAxis, QCPAxis *volumeYAxis);
QCPFinancial* setupCandlestickChart(QCustomPlot* plot, QCPAxisRect* axisRect);
QCPBars* setupVolumeChart(QCustomPlot* plot, QCPAxisRect* axisRect);
void plotData(QCustomPlot* customPlot, const QString& tableName);
QWidget* createSidePanel(QLineEdit*& tableNameEdit, QCustomPlot* customPlot);

