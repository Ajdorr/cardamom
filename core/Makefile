.PHONY: test run build
# Run tests to verify code functionality.
test:
	go test -coverprofile=test/coverage.out -coverpkg=./source/... ./source/...
	go tool cover -html=test/coverage.out
	go tool cover -func=test/coverage.out | grep "total:"

run:
	go run bin/server/server.go

build:
	go build -o cardamom-core bin/server/server.go 
	go build -o cardamom-core-migrate bin/migrate/migrate.go 
