.PHONY: build run compile clean goinstall
BUILD_DIR=./out
BINARY_NAME=${BUILD_DIR}/dictup
SOURCE_MAIN_NAME=./cmd/dictup/main.go
SWAGGER_SCAN=./internal/server/server.go

build:
	go build -o ${BINARY_NAME} ${SOURCE_MAIN_NAME}

run:
	go run ${SOURCE_MAIN_NAME}

run-swag: swagger
	go run ${SOURCE_MAIN_NAME}

compile:
	# 64-Bit
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64.bin ${SOURCE_MAIN_NAME}
	# Windows
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64.exe ${SOURCE_MAIN_NAME}

test:
	go test ./...

clean:
	#go clean
	rm -rfv ${BUILD_DIR}

goinstall:
	go install $$(go list ./...)


check-swagger:
	which swag || (go get -u github.com/swaggo/swag/cmd/swag)

swagger:
	swag init -g ${SWAGGER_SCAN} -o ./docs