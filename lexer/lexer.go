package lexer

type Lexer struct {
	input     string
	startLine int
	endLine   int
	//curLineNum  int
	//nextLineNum int
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	return l
}

func (l *Lexer) readLine() {

}
