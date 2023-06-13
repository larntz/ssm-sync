# variables
BINARY_NAME=ssm-sync
docker-push docker-build: IMAGE_TAG=$(shell date +%Y%m%d)-$(shell git rev-parse --short HEAD)

all: build test
 
build:
	go build -o ./bin/${BINARY_NAME} ./cmd/ssm-sync/main.go
 
docker-build:
	docker build -t larntz/ssm-sync:$(IMAGE_TAG) --no-cache -f build/Dockerfile .

docker-push: docker-build
	docker push larntz/ssm-sync:$(IMAGE_TAG)

test:
 
clean:
	go clean
	rm ${BINARY_NAME}
