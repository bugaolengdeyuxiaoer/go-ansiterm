package ansiterm

type AnsiEventHandler interface {
	// Print
	Print(b byte) error

	// Execute C0 commands
	Execute(b byte) error

	// CUrsor Up
	CUU(int) error

	// CUrsor Down
	CUD(int) error

	// CUrsor Forward
	CUF(int) error

	// CUrsor Backward
	CUB(int) error

	// Cursor to Next Line
	CNL(int) error

	// Cursor to Previous Line
	CPL(int) error

	// Cursor Horizontal position Absolute
	CHA(int) error

	// Vertical line Position Absolute
	VPA(int) error

	// CUrsor Position
	CUP(int, int) error

	// Horizontal and Vertical Position (depends on PUM)
	HVP(int, int) error

	// Text Cursor Enable Mode
	DECTCEM(bool) error

	// Origin Mode
	DECOM(bool) error

	// 132 Column Mode
	DECCOLM(bool) error

	// Erase in Display
	ED(int) error

	// Erase in Line
	EL(int) error

	// Insert Line
	IL(int) error

	// Delete Line
	DL(int) error

	// Insert Character
	ICH(int) error

	// Delete Character
	DCH(int) error

	// Set Graphics Rendition
	SGR([]int) error

	// Pan Down
	SU(int) error

	// Pan Up
	SD(int) error

	// Device Attributes
	DA([]string) error

	// Set Top and Bottom Margins
	DECSTBM(int, int) error

	// Index
	IND() error

	// Reverse Index
	RI() error

	// Flush updates from previous commands
	Flush() error
	CsiXHandler
	CsiSearchHandler
}

type direction int

const (
	forward  direction = 1
	backward direction = 2
)

type CsiXHandler interface {
	// Close quit this session
	Close() error
	// Enter is equal to isExecute
	Enter() error
	// Reset reset param
	Reset() error
	// NextCommand return next command in history
	NextCommand() error
	// PreviousCommand return previous command in history
	PreviousCommand() error
	// EnterWithRedisplay with command still display
	EnterWithRedisplay() error


	// ShowBuffer returns last item that deleted or cut
	ShowBuffer() error
	// Clean clean input buffer and search buffer,and reset status
	Clean() error

	// RemoveForwardWord remove one word forward
	RemoveForwardWord() error
	// RemoveBackwardWord remove one word backward
	RemoveBackwardWord() error

	// RemoveForwardAll remove one word forward
	RemoveForwardAll() error
	// RemoveBackwardAll remove one word backward
	RemoveBackwardAll() error

	// RemoveForwardCharacterOrClose delete any characters
	RemoveForwardCharacterOrClose() error
	// RemoveBackwardCharacter delete any characters
	RemoveBackwardCharacter() error

	// MoveForwardWord cursor move one word forward
	MoveForwardWord() error
	// MoveBackwardWord cursor move one word backward
	MoveBackwardWord() error

	// MoveLineHead cursor move to line head
	MoveLineHead() error
	// 	MoveLineEnd() error cursor move to line head
	MoveLineEnd() error

	// 	MoveForwardCharacter move cursor one character forward
	MoveForwardCharacter() error
	// 	MoveBackwardCharacter move cursor one character backward
	MoveBackwardCharacter() error

	// DoubleX cursor position switch
	DoubleX()error

	SwapLastTwoCharacter() error
}

type CsiSearchHandler interface{
	// QuitSearchMode  quit reverse search mode
	QuitSearchMode() error
	// ReverseSearch reverse search command in history
	ReverseSearch(c byte) error
	// Search search command in history
	Search(c byte) error
}
