# Makefile

BIN_DIR := bin
SRC_DIR := src
BINARY := dsh

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(SRC_DIR)/main.go

clean:
	rm -f $(BIN_DIR)/$(BINARY)

run:
	go run $(SRC_DIR)/main.go
