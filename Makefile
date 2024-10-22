BINARY_NAME=handgame.out
.PHONY: build
	
build: generate bundle
	@go build -o build/${BINARY_NAME} cmd/main.go
	@chmod +x build/${BINARY_NAME}

generate:
	TEMPL_EXPERIMENT=rawgo templ generate

proto:
	protoc --go_out=. --go_opt=paths=source_relative internal/models/protomodels/*.proto

bundle:
	@npx tsc

run: build
	@build/${BINARY_NAME}

clean:
	@go clean
	@rm build/${BINARY_NAME}

test:
	@go test -v ./internal/... -count=1

fuser:
	fuser -k 8080/tcp

test_coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
