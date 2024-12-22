package compiler

type TokenType int

const (
	Boolean TokenType = iota
	Char
	Class
	Constructor
	Do
	Else
	False
	Field
	Function
	If
	Int
	Let
	Method
	Null
	Return
	Static
	This
	True
	Var
	void
	while
	Amp
	ParenthesisLeft
	ParenthesisRight
	Asterisk
	Plus
	Comma
	Minus
	Dot
	Slash
	Semicolon
	LessThan
	Equals
	GreaterThan
	BracketLeft
	BracketRight
	BraceLeft
	Pipe
	BraceRight
	Tidle
	Identifier
	IntConst
	StringConst
	Eof
	StrayChar
	UnterminatedComment
	UnterminatedStringConst
)

var (
	keywords = map[string]TokenType{
		"boolean":     Boolean,
		"char":        Char,
		"class":       Class,
		"constructor": Constructor,
	}
)

type Token struct {
	tokenType TokenType
	strValue  string
	intValue  int
}
