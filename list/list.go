package list

import "courseWork/lexeme"

type TypeList struct {
	Lexeme []lexeme.TypeLexeme
}

func (l *TypeList) AddLexeme(newLexeme lexeme.TypeLexeme) {
	l.Lexeme = append(l.Lexeme, newLexeme)
}
