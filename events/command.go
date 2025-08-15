package events

type Command interface {
	Execute() error
}
