.PHONY: run build clean

# Имя нашего будущего бинарника
APP_NAME=stroydom

run:
	@go run ./cmd/app/main.go

build:
	@go build -o bin/$(APP_NAME) main.go

clean:
	@rm -rf bin/