VERSION ?= 1.0.1
REPO ?= agussyahrilmubarok

hello/run:
	go run cmd/hello/main.go

hello/build:
	go fmt ./...
	go vet ./...
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/hello-docker cmd/hello/main.go

hello/image: hello/build
	docker build -t hello-docker:$(VERSION) -t hello-docker:latest .

hello/push:
	docker tag hello-docker:$(VERSION) $(REPO)/hello-docker:$(VERSION)
	docker tag hello-docker:latest $(REPO)/hello-docker:latest
	docker push $(REPO)/hello-docker:$(VERSION)
	docker push $(REPO)/hello-docker:latest