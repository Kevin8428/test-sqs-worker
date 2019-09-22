package worker

type Handler interface {
	Handle(payload []byte, messageAttributes map[string]string) error
}
