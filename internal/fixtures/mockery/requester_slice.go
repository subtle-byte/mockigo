package mockery

type RequesterSlice interface {
	Get(path string) ([]string, error)
}
