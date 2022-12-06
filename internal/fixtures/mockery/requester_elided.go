package mockery

type RequesterElided interface {
	Get(path, url string) error
}
