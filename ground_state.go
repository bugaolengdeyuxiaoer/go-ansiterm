package ansiterm

type groundState struct {
	baseState
}

func (gs groundState) Handle(b byte) (s state, e error) {
	gs.parser.context.currentChar = b

	nextState, err := gs.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case sliceContains(printables, b):
		return gs, gs.parser.print()
	case sliceContains(executors, b):
		return gs, gs.parser.CsiXDispatcher()
	case isSearchMode(b) :
		// first search clean buffer
		return gs.parser.csiSearch, gs.parser.csiXHandler.Clean()
	case isDoubleMode(b) :
		// first search clean buffer
		return gs.parser.csiX2,nil
	case isRSearchMode(b) :
		return gs.parser.csiRSearch, gs.parser.csiXHandler.Clean()
	}

	return gs, nil
}

func isRSearchMode(b byte) bool {
	if b == 0x12 {
		return true
	}
	return false
}

func isSearchMode(b byte) bool {
	if b == 0x13 {
		return true
	}
	return false
}

func isDoubleMode(b byte) bool {
	if b == 0x18 {
		return true
	}
	return false
}

