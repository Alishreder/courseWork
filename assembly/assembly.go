package assembly

import (
	"courseWork/lexeme"
	"courseWork/list"
	"courseWork/token"
)

type Segment struct {
	Name      token.TypeToken
	Size      int
	LineStart int
	LineEnd   int
}

type Variable struct {
	Line    int
	Name    token.TypeToken
	Type    token.TypeToken
	Value   token.TypeToken
	VOffset int
}

type Label struct {
	Line  int
	Label token.TypeToken
}

type TypeAsm struct {
	Content          string
	List             list.TypeList
	LabelsCounter    int
	Labels           []Label
	SegmentsCounter  int
	Segments         []Segment
	VariablesCounter int
	Variables        []Variable
	LastVOffset      int
}

func CreateAssembly(data string) (assembly *TypeAsm) {
	return &TypeAsm{Content: data}
}

func (asm *TypeAsm) LexemeStructure() {

	for i, v := range asm.List.Lexeme {

		if v.ErrStr != "" {
			continue
		}

		off := 0
		if v.TokensCounter >= 2 &&
			v.Tokens[0].GetTokenType() == "TT_IDENTIFIER" &&
			v.Tokens[1].GetTokenType() == "TT_SYMBOL" &&
			v.Tokens[1].GetFirstLetterOfToken() == ":" {
			v.HasLabel = true
			off += 2 // label:
		}

		if v.TokensCounter <= off { // Only label
			v.HasInstruction = false
			asm.List.Lexeme[i] = v
			continue
		}

		if v.Tokens[off].GetTokenType() == "TT_IDENTIFIER" { // Has name
			v.HasLabel = true
			off += 1
		}

		if v.Tokens[off].GetTokenType() == "TT_INSTRUCTION" ||
			v.Tokens[off].GetTokenType() == "TT_DW_DIRECTIVE" ||
			v.Tokens[off].GetTokenType() == "TT_DB_DIRECTIVE" ||
			v.Tokens[off].GetTokenType() == "TT_DD_DIRECTIVE" ||
			v.Tokens[off].GetTokenType() == "TT_END_DIRECTIVE" ||
			v.Tokens[off].GetTokenType() == "TT_DATA_DIRECTIVE" ||
			v.Tokens[off].GetTokenType() == "TT_CODE_DIRECTIVE" ||
			v.Tokens[off].GetTokenType() == "TT_MODEL_DIRECTIVE" {
			v.HasInstruction = true
		} else {
			asm.List.Lexeme[i] = v
			continue
		}

		v.InstructionIndex = off
		off += 1

		if v.TokensCounter <= off { // Has instruction only
			asm.List.Lexeme[i] = v
			continue
		}
		v.Operands = make([]lexeme.Operand, 5)
		v.Operands[0].OperandIndex = off
		for i := off; i < v.TokensCounter; i++ {
			if v.Tokens[i].TokensType == "TT_SYMBOL" && v.Tokens[i].GetFirstLetterOfToken() == "," {
				v.Operands[v.OperandsCount+1].OperandIndex =
					v.Operands[v.OperandsCount].OperandLength + v.OperandsCount + 1 + off
				v.OperandsCount++
			} else {
				v.Operands[v.OperandsCount].OperandLength += 1
			}
		}

		v.OperandsCount++
		asm.List.Lexeme[i] = v
	}
}

func (asm *TypeAsm) GetLexemeType() {
	for i, v := range asm.List.Lexeme {

		if v.ErrStr != "" {
			continue
		}

		if v.TokensCounter == 1 &&
			v.Tokens[0].TokensType == "TT_DATA_DIRECTIVE" {
			v.Type = "LT_DATA"
		} else if v.TokensCounter == 1 &&
			v.Tokens[0].TokensType == "TT_CODE_DIRECTIVE" {
			v.Type = "LT_CODE"
		} else if v.TokensCounter == 3 &&
			v.Tokens[0].TokensType == "TT_IDENTIFIER" &&
			(v.Tokens[1].TokensType == "TT_DB_DIRECTIVE" ||
				v.Tokens[1].TokensType == "TT_DW_DIRECTIVE" ||
				v.Tokens[1].TokensType == "TT_DD_DIRECTIVE") {
			v.Type = "LT_VAR"
		} else if v.TokensCounter == 2 &&
			v.Tokens[0].TokensType == "TT_IDENTIFIER" &&
			v.Tokens[1].Str == ":" {
			v.Type = "LT_LABEL"
		} else if v.TokensCounter == 1 &&
			v.Tokens[0].TokensType == "TT_END_DIRECTIVE" {
			v.Type = "LT_END"
		} else {
			if v.HasLabel &&
				v.TokensCounter >= 3 &&
				v.Tokens[2].TokensType == "TT_INSTRUCTION" {
				v.Type = "LT_INSTRUCTION"
			} else if !v.HasLabel &&
				v.TokensCounter >= 1 &&
				v.Tokens[0].TokensType == "TT_INSTRUCTION" {
				v.Type = "LT_INSTRUCTION"
			} else {
				v.ErrStr = "Wrong token in lexeme"
			}
		}

		if v.ErrStr == "" {
			asm.List.Lexeme[i].Type = v.Type
		}
	}
}

func (asm *TypeAsm) GetSize(lex *lexeme.TypeLexeme, offset int) int {
	if lex.Type == "LT_VAR" {
		var variable Variable
		for i := 0; i < asm.VariablesCounter; i++ {
			if asm.Variables[i].Line == lex.Line {
				variable = asm.Variables[i]
			}
		}
		switch variable.Type.TokensType {
		case "TT_DB_DIRECTIVE":
			lex.Size = 1
			return 1
		case "TT_DW_DIRECTIVE":
			lex.Size = 2
			return 2
		case "TT_DD_DIRECTIVE":
			lex.Size = 4
			return 4
		}
	} else if lex.Info.Name != "" {
		size := 0
		if lex.Info.Name == "JMP" || lex.Info.Name == "JE" {
			back := true
			far := true
			jmp := true
			if lex.Info.Name == "JE" {
				jmp = false
			}
			for _, v := range asm.Labels {
				if lex.Tokens[1].Str == v.Label.Str {
					if lex.Line > v.Line {
						back = false
						if offset-asm.List.Lexeme[v.Line].Offset < 127 {
							far = false
						}
					}
				}
			}

			if back || (!back && far) {
				size = 5
				if !jmp {
					size++
				}
			} else if !back && !far {
				size = 2
			}

			return size
		}
		if lex.Info.Name == "STD" {
			return 1
		} else if lex.Info.Name == "DIV" {
			if lex.Operands[0].OperandLength == 1 {
				size = 6
				for _, v := range asm.Variables {
					if lex.Tokens[1].Str == v.Name.Str && v.Type.TokensType == "TT_DW_DIRECTIVE" {
						size += 1
						break
					}
				}
			} else if lex.Operands[0].OperandLength <= 5 {
				size = 3
				if lex.Tokens[3].TokensType == "TT_REGISTER16" {
					size += 1
				}
				if lex.Operands[1].Base.Str == "EBP" {
					size += 1
				}
			} else {
				size = 3
				if lex.Operands[0].Segment.Str != "" {
					size += 1
				}
				if lex.Operands[0].Type == "OT_MEM16" {
					size += 1
				}
				if lex.Operands[0].Base.Str == "EBP" {
					size += 1
					if lex.Operands[0].Segment.Str == "SS" {
						size -= 1
					}
				} else if lex.Operands[0].Segment.Str == "DS" {
					size -= 1
				}
			}
			return size
		} else if lex.Info.Name == "IMUL" {
			size = 2
			if lex.Tokens[1].TokensType == "TT_REGISTER16" {
				size += 1
			}
			return size
		} else if lex.Info.Name == "ADD" {
			size = 2
			if lex.Tokens[1].TokensType == "TT_REGISTER16" {
				size += 1
			}
			return size
		} else if lex.Info.Name == "AND" {
			if lex.Operands[1].OperandLength == 1 {
				size = 6
				for _, v := range asm.Variables {
					if lex.Tokens[1].Str == v.Name.Str && v.Type.TokensType == "TT_DW_DIRECTIVE" {
						size += 1
						break
					}
				}
			} else if lex.Operands[1].OperandLength <= 5 {
				size = 3
				if lex.Tokens[3].TokensType == "TT_REGISTER16" {
					size += 1
				}
				if lex.Operands[1].Base.Str == "EBP" {
					size += 1
				}
			} else {
				size = 3
				if lex.Operands[1].Segment.Str != "" {
					size += 1
				}
				if lex.Operands[1].Type == "OT_MEM16" {
					size += 1
				}
				if lex.Operands[1].Base.Str == "EBP" {
					size += 1
					if lex.Operands[1].Segment.Str == "SS" {
						size -= 1
					}
				} else if lex.Operands[1].Segment.Str == "DS" {
					size -= 1
				}
			}
			return size
		} else if lex.Info.Name == "CMP" {
			if lex.Operands[0].OperandLength == 1 {
				size = 7
				for _, v := range asm.Variables {
					if lex.Tokens[1].Str == v.Name.Str && v.Type.TokensType == "TT_DW_DIRECTIVE" {
						size += 1
						break
					}
				}
			} else {
				size = 4
				if lex.Operands[0].Segment.Str != "" {
					size += 1
				}
				if lex.Operands[0].Type == "OT_MEM16" {
					size += 1
				}
				if lex.Operands[0].Base.Str == "EBP" {
					size += 1
					if lex.Operands[0].Segment.Str == "SS" {
						size -= 1
					}
				} else if lex.Operands[0].Segment.Str == "DS" {
					size -= 1
				}
			}
			return size
		} else if lex.Info.Name == "SHL" {
			size = 3
			if lex.Tokens[1].TokensType == "TT_REGISTER16" {
				size += 1
			}
			return size
		} else if lex.Info.Name == "MOV" {
			if lex.Operands[0].OperandLength == 1 {

			} else if lex.Operands[0].OperandLength <= 5 {
				size = 3
				if lex.Tokens[3].TokensType == "TT_REGISTER16" {
					size += 1
				}
				if lex.Operands[1].Base.Str == "EBP" {
					size += 1
				}
			} else {
				size = 3
				if lex.Operands[0].Segment.Str != "" {
					size += 1
				}
				if lex.Operands[0].Type == "OT_MEM16" {
					size += 1
				}
				if lex.Operands[1].Base.Str == "EBP" {
					size += 1
				} else if lex.Operands[1].Segment.Str == "DS" {
					size -= 1
				}
				if lex.Operands[1].Segment.Str == "SS" {
					size -= 1
				}
			}
			return size
		}
	}

	return 0
}
