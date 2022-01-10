package foo

import "github.com/subtle-byte/mockigo/internal/fixtures/mockery/example_project/bar/foo"

type PackageNameSameAsImport interface {
	NewClient() foo.Client
}
