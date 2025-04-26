PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /bindata)

race:
	@bash -c 'for i in {1..100}; do \
		go test -race; \
	done'

cover:
	@go test -coverprofile=cover.out -coverpkg=./... ./...
	@go tool cover -html=cover.out
