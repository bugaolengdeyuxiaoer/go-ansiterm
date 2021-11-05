package ansiterm

func (ap *AnsiParser) collectParam() error {
	currChar := ap.context.currentChar
	ap.logf("collectParam %#x", currChar)
	ap.context.paramBuffer = append(ap.context.paramBuffer, currChar)
	return nil
}

func (ap *AnsiParser) collectInter() error {
	currChar := ap.context.currentChar
	ap.logf("collectInter %#x", currChar)
	ap.context.paramBuffer = append(ap.context.interBuffer, currChar)
	return nil
}

func (ap *AnsiParser) escDispatch() error {
	cmd, _ := parseCmd(*ap.context)
	intermeds := ap.context.interBuffer
	ap.logf("escDispatch currentChar: %#x", ap.context.currentChar)
	ap.logf("escDispatch: %v(%v)", cmd, intermeds)

	switch cmd {
	case "D": // IND
		return ap.eventHandler.IND()
	case "E": // NEL, equivalent to CRLF
		err := ap.eventHandler.Execute(ANSI_CARRIAGE_RETURN)
		if err == nil {
			err = ap.eventHandler.Execute(ANSI_LINE_FEED)
		}
		return err
	case "M": // RI
		return ap.eventHandler.RI()
	case "b":
		return ap.csiXHandler.MoveBackwardWord()
	case "f":
		return ap.csiXHandler.MoveForwardWord()
	}

	return nil
}

func (ap *AnsiParser) csiDispatch() error {
	cmd, _ := parseCmd(*ap.context)
	params, _ := parseParams(ap.context.paramBuffer)
	ap.logf("Parsed params: %v with length: %d", params, len(params))

	ap.logf("csiDispatch: %v(%v)", cmd, params)

	switch cmd {
	case "@":
		return ap.eventHandler.ICH(getInt(params, 1))
	case "A":
		return ap.eventHandler.CUU(getInt(params, 1))
	case "B":
		return ap.eventHandler.CUD(getInt(params, 1))
	case "C":
		return ap.eventHandler.CUF(getInt(params, 1))
	case "D":
		return ap.eventHandler.CUB(getInt(params, 1))
	case "E":
		return ap.eventHandler.CNL(getInt(params, 1))
	case "F":
		return ap.eventHandler.CPL(getInt(params, 1))
	case "G":
		return ap.eventHandler.CHA(getInt(params, 1))
	case "H":
		ints := getInts(params, 2, 1)
		x, y := ints[0], ints[1]
		return ap.eventHandler.CUP(x, y)
	case "J":
		param := getEraseParam(params)
		return ap.eventHandler.ED(param)
	case "K":
		param := getEraseParam(params)
		return ap.eventHandler.EL(param)
	case "L":
		return ap.eventHandler.IL(getInt(params, 1))
	case "M":
		return ap.eventHandler.DL(getInt(params, 1))
	case "P":
		return ap.eventHandler.DCH(getInt(params, 1))
	case "S":
		return ap.eventHandler.SU(getInt(params, 1))
	case "T":
		return ap.eventHandler.SD(getInt(params, 1))
	case "c":
		return ap.eventHandler.DA(params)
	case "d":
		return ap.eventHandler.VPA(getInt(params, 1))
	case "f":
		ints := getInts(params, 2, 1)
		x, y := ints[0], ints[1]
		return ap.eventHandler.HVP(x, y)
	case "h":
		return ap.hDispatch(params)
	case "l":
		return ap.lDispatch(params)
	case "m":
		return ap.eventHandler.SGR(getInts(params, 1, 0))
	case "r":
		ints := getInts(params, 2, 1)
		top, bottom := ints[0], ints[1]
		return ap.eventHandler.DECSTBM(top, bottom)
	default:
		ap.logf("ERROR: Unsupported CSI command: '%s', with full context:  %v", cmd, ap.context)
		return nil
	}

}

func (ap *AnsiParser) CsiXDispatcher() error {
	switch ap.context.currentChar {
	// ctrl a, cursor move line head
	case 0x01:
		return ap.csiXHandler.MoveLineHead()
	// ctrl b, cursor move backward one character
	case 0x02:
		return ap.csiXHandler.MoveBackwardCharacter()
	// ctrl c
	case 0x03:
		return ap.csiXHandler.Clean()
	// ctrl d ,if cmd line is not empty, delete one character forward. else terminate this session
	case 0x04:
		return ap.csiXHandler.RemoveForwardCharacterOrClose()
	// ctrl e, cursor move line end
	case 0x05:
		return ap.csiXHandler.MoveLineEnd()
	// ctrl f, cursor move forward one character
	case 0x06:
		return ap.csiXHandler.MoveForwardCharacter()
	// ctrl g, quit reverse search mode
	case 0x07:
		return ap.csiXHandler.Clean()
	// ctrl h, delete character ,same as backspace
	case 0x08, 0x7f:
		return ap.csiXHandler.RemoveBackwardCharacter()
	// ctrl i,horizon table
	case 0x09:
		break
	// ctrl j, new line, same as enter
	case 0x0a:
		return ap.csiXHandler.Enter()
	// ctrl k, cut the part of the selected line after the cursor and copy it to the clipboard.
	case 0x0b:
		return ap.csiXHandler.RemoveForwardAll()
	// 	Form Feed, equal "clear" command
	case 0x0c:
		return ap.csiXHandler.Clean()
	// ctrl m, carriage return
	case 0x0d:
		return ap.csiXHandler.Enter()
	// ctrl n, next command in history
	case 0x0e:
		return ap.csiXHandler.NextCommand()
	// ctrl o ,enter with command display
	case 0x0f:
		return ap.csiXHandler.EnterWithRedisplay()
	// ctrl p , previous command in history
	case 0x10:
		return ap.csiXHandler.PreviousCommand()
	// ctrl q,
	case 0x11:
		break
	// ctrl r,reverse search history
	case 0x12:
		break
	// ctrl s, search history
	case 0x13:
		break
	// ctrl t, swap the last two characters.
	case 0x14:
		return ap.csiXHandler.SwapLastTwoCharacter()
		// ctrl u,delete all characters from cursor to beginning
	case 0x15:
		return ap.csiXHandler.RemoveBackwardAll()
	// ctrl v
	case 0x16:
		break
	// ctrl w, delete one word before cursor
	case 0x17:
		return ap.csiXHandler.RemoveBackwardWord()
		// ctrl x,double x move cursor from current to ahead of command line.
	case 0x18:
		return ap.csiXHandler.DoubleX()
		// ctrl y,retrieves last item that you deleted or cut
	case 0x19:
		return ap.csiXHandler.ShowBuffer()
	}
	return nil
}

func (ap *AnsiParser) CsiSearch() error {
	return ap.csiSearchHandler.Search(ap.context.currentChar)
}
func (ap *AnsiParser) CsiRSearch() error {
	return ap.csiSearchHandler.ReverseSearch(ap.context.currentChar)
}

func (ap *AnsiParser) enter() error {
	return ap.csiXHandler.Enter()
}

func (ap *AnsiParser) print() error {
	return ap.eventHandler.Print(ap.context.currentChar)
}

func (ap *AnsiParser) clear() error {
	ap.context = &ansiContext{}
	return nil
}

func (ap *AnsiParser) execute() error {
	return ap.eventHandler.Execute(ap.context.currentChar)
}
