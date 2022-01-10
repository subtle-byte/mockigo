package fixtures

type Variadic interface {
	SomeMethod(a int, b ...string) bool
}
