.PHONY: clean
clean:
	go mod tidy

# run all tests
.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...