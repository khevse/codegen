.PHONY:all
all:lint test-all build-codegen

APP_VERSION:=$(shell git tag --points-at HEAD)
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_BUILD_AT:=$(shell date -u '+%FT%T%z')

.PHONY:lint
lint:
	golangci-lint run --config ${CURDIR}/.golangci.yml

.PHONY:test-all
test-all:
	go test -count=1 ./...

.PHONY:build-codegen
build-codegen:
	CGO_ENABLED=0 go build \
	-o tmp/codegen \
	-ldflags="-X github.com/khevse/codegen/internal/pkg/application.Version=${APP_VERSION} -X github.com/khevse/codegen/internal/pkg/application.Commit=${APP_COMMIT} -X github.com/khevse/codegen/internal/pkg/application.BuildAt=${APP_BUILD_AT}" \
	./cmd/codegen/...

.PHONY:install-mockgen
install-mockgen:
	GOBIN="${CURDIR}/bin" go install github.com/gojuno/minimock/v3/cmd/minimock@latest

.PHONY:generate-mocks
generate-mocks: install-mockgen
	bin/minimock -i github.com/khevse/codegen/tests/mainpkg.IObject1 -o ./tests/mainpkg/mocks -s "_mock.go"
	bin/minimock -i github.com/khevse/codegen/tests/mainpkg.IObject2 -o ./tests/mainpkg/mocks -s "_mock.go"