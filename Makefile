VERSION := v2.0.0
APP_NAME := NoteGoat
BUILD_DIR := build

build: clean fmt tidy deps
	mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-$(VERSION).exe .

clean:
	rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...

tidy:
	go mod tidy

deps:
	go mod download