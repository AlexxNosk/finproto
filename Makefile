# Compiler and tools
CXX = g++
CXXFLAGS = -Wall -std=c++17 -fPIC \
	$(shell pkg-config --cflags Qt5Widgets Qt5Sql Qt5PrintSupport) \
	-Icharts -Icharts/cgo
LDFLAGS = $(shell pkg-config --libs Qt5Widgets Qt5Sql Qt5PrintSupport) \
	-Lcharts/cgo -lcgomod

MOC = moc
BUILD_DIR = build

# Source directories
SRC_DIR = charts
CGO_CPPDIR = charts/cgo
CGO_GODIR = backend/cgo

# Files
SRC = $(SRC_DIR)/qt_app.cpp $(SRC_DIR)/main.cpp $(SRC_DIR)/qcustomplot.cpp $(SRC_DIR)/connect.cpp
MOC_HEADERS = $(SRC_DIR)/qcustomplot.h
MOC_SRCS = $(MOC_HEADERS:$(SRC_DIR)/%.h=$(BUILD_DIR)/moc_%.cpp)

OBJS = $(SRC:$(SRC_DIR)/%.cpp=$(BUILD_DIR)/%.o) \
       $(MOC_SRCS:%.cpp=%.o)

TARGET = $(BUILD_DIR)/qcustomplot_app
CGO_SO = $(CGO_CPPDIR)/libcgomod.so

# Default target
all: $(TARGET)

# Ensure build dir exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Compile MOC files
$(BUILD_DIR)/moc_%.cpp: $(SRC_DIR)/%.h | $(BUILD_DIR)
	$(MOC) $< -o $@

# Compile .cpp files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp | $(BUILD_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Compile MOC-generated .cpp files
$(BUILD_DIR)/moc_%.o: $(BUILD_DIR)/moc_%.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

# CGO shared library
$(CGO_SO): $(CGO_GODIR)/cgo_module.go
	cd $(CGO_GODIR) && go build -buildmode=c-shared -o ../../$(CGO_SO) cgo_module.go

# Link final executable
$(TARGET): $(OBJS) $(CGO_SO)
	$(CXX) $(OBJS) -o $(TARGET) $(LDFLAGS)

# Clean
clean:
	rm -rf $(BUILD_DIR) $(CGO_CPPDIR)/libcgomod.so $(CGO_CPPDIR)/libcgomod.h

.PHONY: all clean
