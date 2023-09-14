default:
	echo "Target not selected"

gen-fixtures:
	go generate ./internal/fixtures/...

test:
	go test -cover ./...

install:
	go install ./cmd/mockigo