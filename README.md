mockigo
=======
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/subtle-byte/mockigo?tab=overview)

`mockigo` provides the ability to easily generate type-safe mocks for golang interfaces.

> [!WARNING]
> this project still is quite experimental and the backward incompatible changes are possible.

Table of Contents
-----------------

- [Installation](#installation)
- [Usage in command line](#usage-in-command-line)
- [Generated mocks usage](#generated-mocks-usage)
- [Comparison with `mockery` and `gomock`](#comparison-with-mockery-and-gomock)


Installation
------------

You can use the `go install`:

```
go install github.com/subtle-byte/mockigo/cmd/mockigo@latest
```

Usage in command line
------------

Easiest way to create mocks: run `mockigo` without any arguments in the directory with the package you want to generate mocks for. The generated mocks will be placed in the `mocks_test.go` file in the same directory.

`mockigo` tool has following command line flags:

| flag | default value |description  |
|---|---|---|
| `--targets` | | Comma separated list of interfaces/function types to generate mocks for, if not provided then mocks are generated for all. |
| `--test-pkg` | `false` | Should the generated mocks to be placed in the package with _test suffix. |
| `--out-file` | `mocks_test.go` | Output file (with the mocks). |
| `--gogen` | `true` | Generate `//go:generate` in the output file to ease the regeneration. |

Having `//go:generate` in the generated files, it easy to regenerate all mocks using `go generate ./...` command.

Generated mocks usage
------------

For example you had the file `internal/service/service.go`:

```go
package service

type FooBar interface {
	Foo(n int) (string, error)
	Bar()
}

type MyFunc func(a, b int) int
```

Being in the `internal/service` directory you run `mockigo` and get the `mocks_test.go` file (in the same directory) with mocks.

Then you can create file `internal/service/service_test.go` and use the generated mocks:

```go
package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/subtle-byte/mockigo/match"
)

func TestMocks(t *testing.T) {
	fooBar := NewFooBarMock(t)
	fooBar.EXPECT().Foo(match.Eq(7)).Return("7", nil)
	r1, err := fooBar.Foo(7)
	require.Equal(t, "7", r1)
	require.NoError(t, err)

	myFunc := NewMyFuncMock(t)
	myFunc.EXPECT().Execute(match.Any[int](), match.Any[int]()).Return(10)
	r2 := myFunc.Execute(1, 2)
	require.Equal(t, 10, r2)
}
```

`mathc.Eq` and `match.Any` check the arguments of mock method when it is called. Everything in this testing code is typechecked in compile time using Go 1.18 generics.

More powerful usage (also everything is typechecked in compile time):

```go
package service

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/subtle-byte/mockigo/match"
	"github.com/subtle-byte/mockigo/mock"
)

func TestMocks(t *testing.T) {
	fooBar := NewFooBarMock(t)
	mock.InOrder( // Foo can be called only after at least one call of Bar
		fooBar.EXPECT().Bar(),
		fooBar.EXPECT().Foo(match.MatchedBy(func(n int) bool {
			return n > 0
		})).RunReturn(func(n int) (string, error) {
			return strconv.Itoa(n), nil
		}),
	)
	fooBar.Bar()
	r1, err := fooBar.Foo(9)
	require.Equal(t, "9", r1)
	require.NoError(t, err)
}
```

Comparison with `mockery` and `gomock`
------------

Neither `mockery` nor `gomock` generate type-safe mocks.

For example with `gomock` the following code (which generated using `FooBar` interface defined above) is valid during compilation, but it will fail in runtime:

```go
fooBar.EXPECT().Foo("not int").Return([]byte("not string"), "not error", "not allowed third return")
```

Also with `mockery` the following code is valid during compilation, but it will fail in runtime:

```go
fooBar.EXPECT().Foo("not int").Call.Return(func(notInt string) {
	// returns nothing instead of 2 expected returns
})
```
