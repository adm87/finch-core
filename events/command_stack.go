package events

type CommandStack struct {
	stack      []UndoRedoCommand
	stackIndex int
}

func NewCommandStack() *CommandStack {
	return &CommandStack{
		stack:      make([]UndoRedoCommand, 0),
		stackIndex: -1,
	}
}

func (cs *CommandStack) ExecuteCommand(cmd Command) error {
	if err := cmd.Execute(); err != nil {
		return err
	}
	if undoRedoCmd, ok := cmd.(UndoRedoCommand); ok {
		if cs.stackIndex < len(cs.stack)-1 {
			cs.stack = cs.stack[:cs.stackIndex+1]
		}
		cs.stack = append(cs.stack, undoRedoCmd)
		cs.stackIndex++
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
