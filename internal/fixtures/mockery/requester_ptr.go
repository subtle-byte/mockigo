package mockery

type RequesterPtr interface {
	Get(path string) (*string, error)
}
