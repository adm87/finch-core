package events

type CommandStack struct {
	stack      []UndoRedoCommand
	stackIndex int
	stackLimit int
}

func NewCommandStack(limit uint) *CommandStack {
	if limit == 0 {
		panic("command stack limit must be greater than zero")
	}
	return &CommandStack{
		stack:      make([]UndoRedoCommand, 0),
		stackIndex: -1,
		stackLimit: int(limit),
	}
}

func (cs *CommandStack) ExecuteCommand(cmd Command) error {
	if err := cmd.Execute(); err != nil {
		return err
	}
	if undoRedoCmd, ok := cmd.(UndoRedoCommand); ok {
		cs.push_command(undoRedoCmd)
	}
	return nil
}

func (cs *CommandStack) Undo() error {
	if cs.stackIndex < 0 {
		return nil
	}
	cmd := cs.stack[cs.stackIndex]
	if err := cmd.Undo(); err != nil {
		return err
	}
	cs.stackIndex--
	return nil
}

func (cs *CommandStack) Redo() error {
	if cs.stackIndex+1 >= len(cs.stack) {
		return nil
	}
	cmd := cs.stack[cs.stackIndex+1]
	if err := cmd.Redo(); err != nil {
		return err
	}
	cs.stackIndex++
	return nil
}

func (cs *CommandStack) Clear() {
	cs.stack = make([]UndoRedoCommand, 0)
	cs.stackIndex = -1
}

func (cs *CommandStack) push_command(cmd UndoRedoCommand) {
	cs.stack = cs.stack[:cs.stackIndex+1]

	if cs.stackLimit > 0 && len(cs.stack) >= cs.stackLimit {
		cs.stack = cs.stack[1:]
		cs.stackIndex--
	}

	cs.stack = append(cs.stack, cmd)
	cs.stackIndex = len(cs.stack) - 1
}
