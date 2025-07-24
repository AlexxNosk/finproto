TEMPLATE = app
TARGET = charts

QT += widgets sql printsupport

SOURCES +=  main.cpp \
            qt_app.cpp \
            qcustomplot.cpp

HEADERS += qcustomplot.h
HEADERS += connect.h
HEADERS += qt_app.h