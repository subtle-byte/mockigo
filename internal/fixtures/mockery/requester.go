package mockery

type Requester interface {
	Get(path string) (string, error)
}
