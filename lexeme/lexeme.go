package lexeme

import (
	"courseWork/token"
)

type InstructionInfoType struct {
	Name           string
	OperandCounter int
	Types          []string
}

var InstructionInfos = []InstructionInfoType{
	{Name: "STD", OperandCounter: 0},

	{Name: "DIV", OperandCounter: 1, Types: []string{"OT_MEM8"}},
	{Name: "DIV", OperandCounter: 1, Types: []string{"OT_MEM32"}},

	{Name: "IMUL", OperandCounter: 1, Types: []string{"OT_REG8"}},
	{Name: "IMUL", OperandCounter: 1, Types: []string{"OT_REG16"}},
	{Name: "IMUL", OperandCounter: 1, Types: []string{"OT_REG32"}},

	{Name: "ADD", OperandCounter: 2, Types: []string{"OT_REG8", "OT_REG8"}},
	{Name: "ADD", OperandCounter: 2, Types: []string{"OT_REG16", "OT_REG16"}},
	{Name: "ADD", OperandCounter: 2, Types: []string{"OT_REG32", "OT_REG32"}},

	{Name: "AND", OperandCounter: 2, Types: []string{"OT_REG8", "OT_MEM8"}},
	{Name: "AND", OperandCounter: 2, Types: []string{"OT_REG16", "OT_MEM16"}},
	{Name: "AND", OperandCounter: 2, Types: []string{"OT_REG32", "OT_MEM32"}},

	{Name: "CMP", OperandCounter: 2, Types: []string{"OT_MEM8", "OT_IMM8"}},
	{Name: "CMP", OperandCounter: 2, Types: []string{"OT_MEM16", "OT_IMM16"}},
	{Name: "CMP", OperandCounter: 2, Types: []string{"OT_MEM32", "OT_IMM32"}},

	{Name: "SHL", OperandCounter: 2, Types: []string{"OT_REG8", "OT_IMM8"}},
	{Name: "SHL", OperandCounter: 2, Types: []string{"OT_REG16", "OT_IMM16"}},
	{Name: "SHL", OperandCounter: 2, Types: []string{"OT_REG32", "OT_IMM32"}},

	{Name: "MOV", OperandCounter: 2, Types: []string{"OT_MEM8", "OT_REG8"}},
	{Name: "MOV", OperandCounter: 2, Types: []string{"OT_MEM16", "OT_REG16"}},
	{Name: "MOV", OperandCounter: 2, Types: []string{"OT_MEM32", "OT_REG32"}},

	{Name: "JMP", OperandCounter: 1, Types: []string{"OT_LABEL_BACKWARD"}},
	{Name: "JMP", OperandCounter: 1, Types: []string{"OT_LABEL_FORWARD"}},

	{Name: "JE", OperandCounter: 1, Types: []string{"OT_LABEL_BACKWARD"}},
	{Name: "JE", OperandCounter: 1, Types: []string{"OT_LABEL_FORWARD"}},
}

type Operand struct {
	// Index of the operand in lexeme.tokens
	OperandIndex  int
	OperandLength int
	Operand       token.TypeToken
	Type          string
	TypeKeyword   token.TypeToken
	Segment       token.TypeToken
	Base          token.TypeToken
	Index         token.TypeToken
}

type TypeLexeme struct {
	TokensCounter    int
	Tokens           []token.TypeToken
	HasLabel         bool
	HasInstruction   bool
	InstructionIndex int
	Line             int
	Size             int
	OperandsCount    int
	Source           string
	ErrStr           string
	Type             string
	// Information about operands
	Operands []Operand
	Info     InstructionInfoType
	Offset   int
}

var OTypes = []string{
	"OT_REG8",
	"OT_REG16",
	"OT_REG32",
	"OT_IMM8",
	"OT_IMM16",
	"OT_IMM32",
	"OT_LABEL_FORWARD",
	"OT_LABEL_BACKWARD",
	"OT_MEM32",
	"OT_MEM16",
	"OT_MEM8",
}

var LTypes = []string{
	"LT_DATA",
	"LT_CODE",
	"LT_MODEL",
	"LT_VAR",
	"LT_LABEL",
	"LT_INSTRUCTION",
	"LT_END",
}

func (lex *TypeLexeme) ClearLexeme() {
	lex.TokensCounter = 0
	lex.Tokens = nil
}
