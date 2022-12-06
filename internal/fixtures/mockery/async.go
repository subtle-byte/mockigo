package mockery

type AsyncProducer interface {
	Input() chan<- bool
	Output() <-chan bool
	Whatever() chan bool
}
