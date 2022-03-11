.PHONY: build
build:
	CGO_ENABLED=0 go build -v \
		-ldflags "-X nexlab.tech/core/services/auth/version.Tag=${TAG} -X nexlab.tech/core/services/auth/version.GitCommit=${GIT_COMMIT}" \
		-o auth-api ./server

.PHONY: test
test:
	go test ./...
