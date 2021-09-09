.PHONY: build run compile clean goinstall
BUILD_DIR=./out
BINARY_NAME=${BUILD_DIR}/dictup
SOURCE_MAIN_NAME=./cmd/dictup/main.go
SWAGGER_SCAN=./internal/server/server.go
PORT=8080

build: copy-configs
	go build -o ${BINARY_NAME} ${SOURCE_MAIN_NAME}


copy-configs:
	cp -r ./configs ${BUILD_DIR}

run:
	go run ${SOURCE_MAIN_NAME}

run-swag: swagger
	go run ${SOURCE_MAIN_NAME}

compile: copy-configs
	# 64-Bit
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64.bin ${SOURCE_MAIN_NAME}
	# Windows
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64.exe ${SOURCE_MAIN_NAME}

test:
	go test ./...

clean-build:
	#go clean
	rm -rfv ${BUILD_DIR}

goinstall:
	go install $$(go list ./...)


check-swagger:
	which swag || (go get -u github.com/swaggo/swag/cmd/swag)

swagger:
	swag init -g ${SWAGGER_SCAN} -o ./api

docker-build:
	docker build -t dictup .

docker-run:
	 docker run --rm --name dictup --network host dictup
