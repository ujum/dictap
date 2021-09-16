.PHONY: build run compile clean goinstall
BUILD_DIR=./out
BINARY_NAME=${BUILD_DIR}/dictup
SOURCE_MAIN_NAME=./cmd/dictup/main.go
SWAGGER_SCAN=./internal/server/server.go

build: copy-configs swagger
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} ${SOURCE_MAIN_NAME}

copy-configs:
	mkdir -p ${BUILD_DIR} && cp -r ./configs --parents ${BUILD_DIR}

run:
	go run ${SOURCE_MAIN_NAME}

run-swag: swagger
	go run ${SOURCE_MAIN_NAME}

compile: copy-configs swagger
	# 64-Bit
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64.bin ${SOURCE_MAIN_NAME}
	# Windows
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64.exe ${SOURCE_MAIN_NAME}

test:
	go test ./... -cover

clean-build:
	#go clean
	rm -rfv ${BUILD_DIR}

goinstall:
	go install $$(go list ./...)


check-swagger:
	which swag || (go install github.com/swaggo/swag/cmd/swag)

swagger: check-swagger
	swag init -g ${SWAGGER_SCAN} -o ./api

docker-build:
	docker build -t dictup .

docker-run:
	 docker run --rm --name dictup --network host dictup

test-cover-report:
	mkdir -p ${BUILD_DIR}
	go test ./... -cover -coverprofile=${BUILD_DIR}/test-coverage.out
	go tool cover -html=${BUILD_DIR}/test-coverage.out -o ${BUILD_DIR}/test-coverage.html