package lexer

type Lexer struct {
	strs      []string
	startLine int
	endLine   int
	//curLineNum  int
	//nextLineNum int
}

func New(strs []string) *Lexer {
	l := &Lexer{
		strs: strs,
	}
	return l
}

func (l *Lexer) readLine() {

}
