.PHONY: build install test test-coverage test-coverage-summary licenses

build:
	go build .

install:
	go install .

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out; \
	go tool cover -html=coverage.out

test-coverage-summary:
	go test ./... -coverprofile=coverage.out; \
	go tool cover -func=coverage.out

THIRD_PARTY_NOTICES:
	rm -rf THIRD_PARTY_NOTICES
	go run github.com/google/go-licenses@v1.0.0 \
		save ./... \
		--save_path=THIRD_PARTY_NOTICES
