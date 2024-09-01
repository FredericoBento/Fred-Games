BINARY_NAME=handgame.out
.PHONY: build
	
build: generate
	@go build -o build/${BINARY_NAME} cmd/main.go
	@chmod +x build/${BINARY_NAME}

generate:
	templ generate

run: build
	@build/${BINARY_NAME}

clean:
	@go clean
	@rm build/${BINARY_NAME}

test:
	@go test -v ./... -count=1
