
#include "qt_app.h"

// Global data vectors
QVector<double> x, open, high, low, close, volume;

// === Database and Data Loading ===
void loadDataFromPostgreSQL(const QString& tableName) {
    QString connName = "pg_conn";
    if (QSqlDatabase::contains(connName)) {
        QSqlDatabase::removeDatabase(connName);
    }

    QSqlDatabase db = QSqlDatabase::addDatabase("QPSQL", connName);
    db.setHostName("localhost");
    db.setPort(5434);
    db.setDatabaseName("finProto_db");
    db.setUserName("root");
    db.setPassword("root");

    if (!db.open()) {
        qDebug() << "Database error:" << db.lastError().text();
        return;
    }

    QSqlQuery query(db);
    query.exec(QString("SELECT EXTRACT(EPOCH FROM timestamp) AS ts, open, high, low, close, volume FROM %1 ORDER BY timestamp").arg(tableName));

    x.clear(); open.clear(); high.clear(); low.clear(); close.clear(); volume.clear();

    while (query.next()) {
        x.append(query.value("ts").toDouble());
        open.append(query.value("open").toDouble());
        high.append(query.value("high").toDouble());
        low.append(query.value("low").toDouble());
        close.append(query.value("close").toDouble());
        volume.append(query.value("volume").toDouble());
    }

    db.close();
}


// === Dynamic Y-Axis Adjustment ===
// Rescale Y axis to match visible X range
void adjustYAxisToVisibleRange(QCPFinancial *candlesticks, QCPBars *volumeBars,
                                QCPAxis *xAxis, QCPAxis *yAxis, QCPAxis *volumeYAxis)
{
    if (!candlesticks || !candlesticks->data() || candlesticks->data()->isEmpty())
        return;

    QCPRange xRange = xAxis->range();

    // Track min/max for price
    double minPrice = std::numeric_limits<double>::max();
    double maxPrice = std::numeric_limits<double>::lowest();

    // Track max for volume
    double maxVolume = 0;

    for (auto it = candlesticks->data()->constBegin(); it != candlesticks->data()->constEnd(); ++it) {
        if (it->key >= xRange.lower && it->key <= xRange.upper) {
            minPrice = std::min(minPrice, it->low);
            maxPrice = std::max(maxPrice, it->high);
        }
    }

    for (auto it = volumeBars->data()->constBegin(); it != volumeBars->data()->constEnd(); ++it) {
        if (it->key >= xRange.lower && it->key <= xRange.upper) {
            maxVolume = std::max(maxVolume, it->value);
        }
    }

    if (minPrice < maxPrice) {
        yAxis->setRange(minPrice * 0.98, maxPrice * 1.02);
    }
    if (maxVolume > 0) {
        volumeYAxis->setRange(0, maxVolume * 1.1);
    }
}


// === Plotting Functions ===
QCPFinancial* setupCandlestickChart(QCustomPlot* plot, QCPAxisRect* axisRect) {
    QCPFinancial *candlesticks = new QCPFinancial(axisRect->axis(QCPAxis::atBottom),
                                                  axisRect->axis(QCPAxis::atLeft));
    QVector<QCPFinancialData> ohlc;
    for (int i = 0; i < x.size(); ++i) {
        ohlc.append(QCPFinancialData(x[i], open[i], high[i], low[i], close[i]));
    }

    double barWidth = (x.size() >= 2) ? (x[1] - x[0]) : (60 * 60 * 24);

    candlesticks->setName("Candlestick");
    candlesticks->data()->set(ohlc);
    candlesticks->setChartStyle(QCPFinancial::csCandlestick);
    candlesticks->setWidth(barWidth);
    candlesticks->setTwoColored(true);
    candlesticks->setBrushPositive(QColor(0, 255, 0, 100));
    candlesticks->setBrushNegative(QColor(255, 0, 0, 100));
    candlesticks->rescaleAxes();

    return candlesticks;
}

QCPBars* setupVolumeChart(QCustomPlot* plot, QCPAxisRect* axisRect) {
    QCPBars *volumeBars = new QCPBars(axisRect->axis(QCPAxis::atBottom),
                                      axisRect->axis(QCPAxis::atLeft));
    double barWidth = (x.size() >= 2) ? (x[1] - x[0]) : (60 * 60 * 24);
    volumeBars->setWidth(barWidth);
    volumeBars->setData(x, volume);
    volumeBars->setBrush(QColor(100, 100, 250, 150));
    volumeBars->setPen(Qt::NoPen);
    volumeBars->rescaleAxes(true);

    return volumeBars;
}


void plotData(QCustomPlot* customPlot, const QString& tableName) {
    loadDataFromPostgreSQL(tableName);
    customPlot->clearPlottables();
    customPlot->plotLayout()->clear();
    customPlot->clearItems();

    QCPAxisRect *priceRect = new QCPAxisRect(customPlot, true);
    QCPAxisRect *volumeRect = new QCPAxisRect(customPlot, true);

    customPlot->plotLayout()->addElement(0, 0, priceRect);
    customPlot->plotLayout()->addElement(1, 0, volumeRect);

    // === Shared margin group ===
    QCPMarginGroup *marginGroup = new QCPMarginGroup(customPlot);

    priceRect->setMarginGroup(QCP::msLeft | QCP::msRight, marginGroup);
    volumeRect->setMarginGroup(QCP::msLeft | QCP::msRight, marginGroup);


    // === Date ticker ===
    QSharedPointer<QCPAxisTickerDateTime> dateTicker(new QCPAxisTickerDateTime);
    dateTicker->setDateTimeFormat("dd.MM.yyyy");
    priceRect->axis(QCPAxis::atBottom)->setTicker(dateTicker);
    volumeRect->axis(QCPAxis::atBottom)->setTicker(dateTicker);

    // === Synchronize X axes ===
    QObject::connect(priceRect->axis(QCPAxis::atBottom), SIGNAL(rangeChanged(QCPRange)),
                     volumeRect->axis(QCPAxis::atBottom), SLOT(setRange(QCPRange)));
    QObject::connect(volumeRect->axis(QCPAxis::atBottom), SIGNAL(rangeChanged(QCPRange)),
                     priceRect->axis(QCPAxis::atBottom), SLOT(setRange(QCPRange)));

    // === Interactions ===
    customPlot->setInteractions(QCP::iRangeDrag | QCP::iRangeZoom);
    priceRect->setRangeDrag(Qt::Horizontal);
    volumeRect->setRangeDrag(Qt::Horizontal);
    priceRect->setRangeZoom(Qt::Horizontal);
    volumeRect->setRangeZoom(Qt::Horizontal);

    // === Plot data ===
    QCPFinancial* candlesticks = setupCandlestickChart(customPlot, priceRect);
    QCPBars* volumeBars = setupVolumeChart(customPlot, volumeRect);

    // Connect X-axis range changes to adjust Y scaling dynamically
    QObject::connect(priceRect->axis(QCPAxis::atBottom), static_cast<void (QCPAxis::*)(const QCPRange &)>(&QCPAxis::rangeChanged),
        [=](const QCPRange&) {
            adjustYAxisToVisibleRange(candlesticks, volumeBars,
                priceRect->axis(QCPAxis::atBottom),
                priceRect->axis(QCPAxis::atLeft),
                volumeRect->axis(QCPAxis::atLeft));
        });


    // === Hide bottom axis of price chart ===
    priceRect->axis(QCPAxis::atBottom)->setVisible(true);

    // === Equal Y-axis padding ===
    priceRect->axis(QCPAxis::atLeft)->setPadding(60);
    volumeRect->axis(QCPAxis::atLeft)->setPadding(60);

    // === Enable right axes and link ranges to left axes ===
    priceRect->axis(QCPAxis::atRight)->setVisible(true);
    priceRect->axis(QCPAxis::atRight)->setTickLabels(true);
    priceRect->axis(QCPAxis::atRight)->setRange(priceRect->axis(QCPAxis::atLeft)->range());
    QObject::connect(priceRect->axis(QCPAxis::atLeft), SIGNAL(rangeChanged(QCPRange)),
                    priceRect->axis(QCPAxis::atRight), SLOT(setRange(QCPRange)));

    volumeRect->axis(QCPAxis::atRight)->setVisible(true);
    volumeRect->axis(QCPAxis::atRight)->setTickLabels(true);
    volumeRect->axis(QCPAxis::atRight)->setRange(volumeRect->axis(QCPAxis::atLeft)->range());
    QObject::connect(volumeRect->axis(QCPAxis::atLeft), SIGNAL(rangeChanged(QCPRange)),
                    volumeRect->axis(QCPAxis::atRight), SLOT(setRange(QCPRange)));

    // === Layout height ratios ===
    customPlot->plotLayout()->setRowStretchFactor(0, 3);  // price
    customPlot->plotLayout()->setRowStretchFactor(1, 1);  // volume

    customPlot->replot();
}

// === UI Side Panel ===
QWidget* createSidePanel(QLineEdit*& tableNameEdit, QCustomPlot* customPlot) {
    QWidget *panel = new QWidget;
    QVBoxLayout *layout = new QVBoxLayout(panel);

    tableNameEdit = new QLineEdit("bars_sber_misx_d");
    QPushButton *loadButton = new QPushButton("Load Table");

    layout->addWidget(tableNameEdit);
    layout->addWidget(loadButton);
    layout->addStretch();

    QObject::connect(loadButton, &QPushButton::clicked, [=]() {
        QString tableName = tableNameEdit->text().trimmed();
        if (!tableName.isEmpty()) {
            plotData(customPlot, tableName);
        }
    });

    panel->setFixedWidth(200);
    return panel;
}


