BINARY_NAME=ssm-sync
 
all: build test
 
build:
	go build -o ${BINARY_NAME} ./cmd/ssm-sync/main.go
 
test:
 
clean:
	go clean
	rm ${BINARY_NAME}
