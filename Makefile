default:
	echo "Target not selected"

gen-fixtures:
	rm -rf internal/fixtures/mocks
	cd internal/fixtures && go run ../../cmd/mockigo

test:
	go test `go list ./... | grep -v internal/fixtures | grep -v cmd` github.com/subtle-byte/mockigo/internal/fixtures/tester -coverprofile cover.out
	rm cover.out
