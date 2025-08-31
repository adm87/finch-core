package events

type Command interface {
	Execute() error
}

type UndoRedoCommand interface {
	Command

	Undo() error
	Redo() error
}
