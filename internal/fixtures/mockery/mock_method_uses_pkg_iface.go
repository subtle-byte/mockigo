package mockery

type Sibling interface {
	DoSomething()
}

type UsesOtherPkgIface interface {
	DoSomethingElse(obj Sibling)
}
