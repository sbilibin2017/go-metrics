.PHONY: test test-cov migrate lint mockgen

test:
	go test ./...

test-cov:
	go test -coverprofile=coverage.out -covermode=count ./...		
	go tool cover -func=coverage.out	

lint:
	staticcheck ./...

mockgen:
	@echo "Generating mock for: $(file)"
	@mockgen -source=$(file) \
		-destination=$(dir $(file))$(notdir $(basename $(file)))_mock.go \
		-package=$(shell basename $(dir $(file))) 


