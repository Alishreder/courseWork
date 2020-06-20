package main

import (
	"courseWork/assembly"
	"courseWork/lexeme"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func stringToNum(str string, base int) (result int) {

	hasSign := string(str[0]) == "-"
	length := len(str)
	i := 0
	if hasSign {
		i = 1
	}
	c := ""
	for ; i < length-1; i++ {
		c = string(str[i])
		literal := -1
		if strings.ContainsAny(c, "1234567890ABCDEF") {
			literal, _ = strconv.Atoi(c)
			result += int(math.Pow(float64(base), float64(length-i-1))) * literal
		} else {
			panic("Can`t find that literal")
		}
	}
	if hasSign {
		return -result
	}

	return
}

func FetchOpInfo(assembl *assembly.TypeAsm, lex *lexeme.TypeLexeme) bool {
	asm := assembl

	if !lex.HasInstruction || lex.OperandsCount == 0 {
		return true
	}

	for i := 0; i < lex.OperandsCount; i++ {
		op := &lex.Operands[i]
		if op.OperandLength == 1 {
			op.Operand = lex.Tokens[op.OperandIndex]
			switch op.Operand.TokensType {
			case "TT_REGISTER8":
				op.Type = "OT_REG8"
			case "TT_REGISTER16":
				op.Type = "OT_REG16"
			case "TT_REGISTER32":
				op.Type = "OT_REG32"
			case "TT_NUMBER2":
				value := stringToNum(op.Operand.Str, 2)
				if math.Abs(float64(value)) < 2^8 {
					op.Type = "OT_IMM8"
				} else if math.Abs(float64(value)) < 65535 {
					op.Type = "OT_IMM16"
				} else {
					op.Type = "OT_IMM32"
				}
			case "TT_NUMBER10":
				value := stringToNum(op.Operand.Str, 10)
				if math.Abs(float64(value)) < 2^8 {
					op.Type = "OT_IMM8"
				} else if math.Abs(float64(value)) < 65535 {
					op.Type = "OT_IMM16"
				} else {
					op.Type = "OT_IMM32"
				}
			case "TT_NUMBER16":
				value := stringToNum(op.Operand.Str, 16)
				if math.Abs(float64(value)) < 2^8 {
					op.Type = "OT_IMM8"
				} else if math.Abs(float64(value)) < 65535 {
					op.Type = "OT_IMM16"
				} else {
					op.Type = "OT_IMM32"
				}
			case "TT_IDENTIFIER":
				if lex.Type != "LT_INSTRUCTION" {
					var label assembly.Label

					founded := false
					for _, v := range asm.Labels {
						if v.Label.Str == op.Operand.Str {
							founded = true
							label = v
						}
					}

					if founded {
						if label.Line > lex.Line {
							op.Type = "OT_LABEL_FORWARD"
						} else {
							op.Type = "OT_LABEL_BACKWARD"
						}
					} else {
						lex.ErrStr = "Undefined reference"
						return false
					}
				}
			default:
				lex.ErrStr = "Wrong token in operator"
			}
		} else {
			op.Type = "OT_MEM"

			offset := 0

			// Minimal length: 5 tokens
			// [ ecx + edi ]
			// Maximum length: 9 tokens
			// dword ptr ES : [ edx + esi ]
			if op.OperandLength >= 9 && op.OperandLength <= 5 {
				lex.ErrStr = "Instruction has wrong lengthhh"
				return false
			}

			opTokens := lex.Tokens
			curIndex := op.OperandIndex
			// Has type specifier
			if opTokens[curIndex].TokensType == "TT_DWORD" ||
				opTokens[curIndex].TokensType == "TT_BYTE" ||
				opTokens[curIndex].TokensType == "TT_WORD" {
				if opTokens[curIndex+1].TokensType != "TT_PTR" {
					lex.ErrStr = "Instruction has wrong format, should be PTR"
					return false
				}

				if opTokens[curIndex].TokensType == "TT_DWORD" {
					op.Type = "OT_MEM32"
				} else if opTokens[curIndex].TokensType == "TT_WORD" {
					op.Type = "OT_MEM16"
				} else {
					op.Type = "OT_MEM8"
				}

				op.TypeKeyword = opTokens[curIndex]
				offset += 2
			}
			curIndex += offset
			// Has segment prefix
			if opTokens[curIndex].TokensType == "TT_REGISTER_SEG" {
				// After SegReg always must be colon
				if !(opTokens[curIndex+1].TokensType == "TT_SYMBOL" &&
					opTokens[curIndex+1].GetFirstLetterOfToken() == ":") {
					lex.ErrStr = "Expected \":\" after register segment"
					return false
				}
				op.Segment = opTokens[curIndex]
				offset += 2
				curIndex += 2
			}

			// Check for minimal length
			if op.OperandLength-offset != 5 {
				lex.ErrStr = "Instruction has wrong lengthttt"
				fmt.Println(opTokens[curIndex].TokensType)
				return false
			}

			// Check for [REG + REG] format
			if !(opTokens[curIndex].TokensType == "TT_SYMBOL" &&
				opTokens[curIndex].GetFirstLetterOfToken() == "[") {
				lex.ErrStr = "Expected [ "
				return false
			}

			if !(opTokens[curIndex+1].TokensType == "TT_REGISTER8" ||
				opTokens[curIndex+1].TokensType == "TT_REGISTER16" ||
				opTokens[curIndex+1].TokensType == "TT_REGISTER32") {
				lex.ErrStr = "Expected register"
				return false
			}

			if !(opTokens[curIndex+2].TokensType == "TT_SYMBOL" &&
				opTokens[curIndex+2].Str == "+") {
				lex.ErrStr = "Expected register"
				return false
			}

			if !(opTokens[curIndex+3].TokensType == "TT_REGISTER8" ||
				opTokens[curIndex+3].TokensType == "TT_REGISTER16" ||
				opTokens[curIndex+3].TokensType == "TT_REGISTER32") {
				lex.ErrStr = "Expected register"
				return false
			}

			if !(opTokens[curIndex+4].TokensType == "TT_SYMBOL" &&
				opTokens[curIndex+4].GetFirstLetterOfToken() == "]") {
				lex.ErrStr = "Expected [ "
				return false
			}

			op.Base = opTokens[curIndex+1]
			op.Index = opTokens[curIndex+3]

			// Registers must have same size
			if !(op.Base.TokensType == op.Index.TokensType) {
				lex.ErrStr = "Base and index register must have same size"
				return false
			}
		}
	}

	return true
}

func AssignInstruction(lex *lexeme.TypeLexeme) bool {

	if !lex.HasInstruction {
		return true
	}

	instruction := lex.Tokens[lex.InstructionIndex]

	for i := 0; i < len(lexeme.InstructionInfos); i++ {
		info := lexeme.InstructionInfos[i]

		// Compare by name
		if instruction.Str != info.Name {
			continue
		}

		// Compare by operands count
		if lex.OperandsCount != info.OperandCounter {
			continue
		}

		// Compare by operand types. We should be careful with OT_MEM because its size undefined.
		typeMismatch := false
		for j := 0; i < info.OperandCounter; j++ {
			if lex.Operands[j].Type != info.Types[j] {
				typeMismatch = true

				if lex.Operands[j].Type == "OT_MEM" && info.Types[j] == "OT_MEM32" {
					typeMismatch = false
					continue
				}
				if lex.Operands[j].Type == "OT_MEM" && info.Types[j] == "OT_MEM16" {
					typeMismatch = false
					continue
				}
				if lex.Operands[j].Type == "OT_MEM" && info.Types[j] == "OT_MEM8" {
					typeMismatch = false
					continue
				}
				break
			}
		}

		if typeMismatch {
			continue
		}

		lex.Info = info
		return true
	}

	lex.ErrStr = "Can't find matching instruction"
	return false
}

func FirstPass(asm *assembly.TypeAsm) {
	asm.GetLexemeType()

	isSegmentCreated := false
	isSegDataExists := false
	isSegCodeExists := false
	for i := 0; i < len(asm.List.Lexeme); i++ {
		if asm.List.Lexeme[i].ErrStr != "" {
			continue
		}

		lex := &asm.List.Lexeme[i]
		if !isSegmentCreated {
			asm.Segments = make([]assembly.Segment, 2)
			isSegmentCreated = true
		}

		if lex.Type == "LT_DATA" {
			if !isSegDataExists {
				asm.Segments[0].Name = lex.Tokens[0]
				asm.Segments[0].LineStart = lex.Line

				isSegDataExists = true
			} else {
				lex.ErrStr = ".data segment already declared"
			}
		} else if lex.Type == "LT_CODE" {
			if !isSegDataExists {
				lex.ErrStr = "Firstly should be .data directive"
				continue
			}

			if !isSegCodeExists {
				asm.Segments[0].LineEnd = lex.Line - 1

				asm.Segments[1].Name = lex.Tokens[0]
				asm.Segments[1].LineStart = lex.Line

				isSegCodeExists = true
			} else {
				lex.ErrStr = ".code segment already declared"
			}
		} else if lex.Type == "LT_END" {
			if !isSegCodeExists {
				lex.ErrStr = "Missed .code segment"
				continue
			} else {
				asm.Segments[1].LineEnd = lex.Line - 1
			}
		} else if lex.Type == "LT_VAR" {
			var variable assembly.Variable
			variable.Line = lex.Line
			variable.Name = lex.Tokens[0]
			variable.Type = lex.Tokens[1]
			variable.Value = lex.Tokens[2]

			founded := false
			for _, v := range asm.Variables {
				if v.Name.Str == variable.Name.Str {
					founded = true
				}
			}

			if founded {
				lex.ErrStr = "Same variable already declared"
			} else {
				asm.Variables = append(asm.Variables, variable)
				asm.VariablesCounter++
			}
		} else if lex.Type == "LT_LABEL" || (lex.Type == "LT_INSTRUCTION" && lex.HasLabel) {
			var label assembly.Label
			label.Line = lex.Line
			label.Label = lex.Tokens[0]

			founded := false
			for _, v := range asm.Labels {
				if v.Label.Str == lex.Tokens[0].Str {
					founded = true
				}
			}

			if founded {
				lex.ErrStr = "Same label already declared"
			} else {
				asm.Labels = append(asm.Labels, label)
				asm.LabelsCounter++
			}
		} else if lex.Type == "LT_INSTRUCTION" {
			if !lex.HasInstruction {
				lex.ErrStr = "Unexpected token met"
			}
		} else if lex.Tokens[0].TokensType == "TT_MODEL_DIRECTIVE" {
			continue
		} else {
			if lex.Type != "LT_END" {

				lex.ErrStr = "Unexpected token met"
			}
		}

	}

	var activeSegment assembly.Segment

	for i := 0; i < len(asm.List.Lexeme); i++ {

		lex := &asm.List.Lexeme[i]

		if lex.ErrStr == "" {
			if lex.Type == "LT_INSTRUCTION" {
				if !FetchOpInfo(asm, lex) {
					continue
				}
				if !AssignInstruction(lex) {
					continue
				}
			}
			if lex.Type == "LT_DATA" || lex.Type == "LT_CODE" {
				for i := 0; i < asm.SegmentsCounter; i++ {
					if asm.Segments[i].LineStart <= lex.Line && lex.Line <= asm.Segments[i].LineEnd {
						activeSegment = asm.Segments[i]
					}
				}
			}
			if lex.Type == "LT_DATA" || lex.Type == "LT_CODE" {
				if lex.Type == "LT_CODE" {
					asm.Segments[0].Size = activeSegment.Size
					asm.LastVOffset = activeSegment.Size
				}
				activeSegment.Size = 0
			}
			lex.Size = asm.GetSize(lex, activeSegment.Size)
			if lex.Type == "LT_CODE" {
				asm.LastVOffset += lex.Size
			}
			lex.Offset = activeSegment.Size

			activeSegment.Size += lex.Size
			if lex.Type == "LT_DATA" || lex.Type == "LT_CODE" || lex.Type == "LT_END" {
				if lex.Type == "LT_END" {
					asm.Segments[1].Size = activeSegment.Size
				}
				activeSegment.Size = 0
			}

			if lex.ErrStr != "" {
				continue
			}
		}
	}
}

func PrintFirstPass(asm *assembly.TypeAsm) {
	file, err := os.Create("testFirstPass.docx")
	if err != nil {
		panic("PrintFirstStage: Can't create file for write")
	}

	res := ""

	res += "/------------------------------------------------------------------------\\\n"
	res += "| LN |  OFFSET  | S |                       SOURCE                       |\n"
	res += "|----+----------+---+----------------------------------------------------|\n"
	for i, v := range asm.List.Lexeme {
		errr := ""
		if v.ErrStr != "" {
			errr = "E"
			res += fmt.Sprintf(v.ErrStr)
		}
		hexOffset := fmt.Sprintf("%X", v.Offset)
		res += fmt.Sprintf("| %2d |  %-6s  | %d | %s   %-44s    |\n", i, hexOffset, v.Size, errr, v.Source)
	}
	res += "\\------------------------------------------------------------------------/\n\n"
	res += "SEGMENTS TABLE\n"
	res += "/----------------------------\\\n"
	res += "| ID |   SEGMENT  |   SIZE   |\n"
	res += "|----+------------+----------|\n"
	for i, v := range asm.Segments {
		res += fmt.Sprintf("| %2d |   %-6s   |  %-6X  |\n", i, v.Name.Str, v.Size)
	}
	res += "\\----------------------------/\n\n"

	res += "USER DEFINED SYMBOLS\n"
	res += "/----------------------------------------\\\n"
	res += "| ID |   SYMBOL   |   TYPE   |   VALUE   |\n"
	res += "|----+------------+----------+-----------|\n"
	offset := 0
	i := 0
	for _, v := range asm.Variables {
		seg := ""
		for _, val := range asm.Segments {
			if val.LineStart <= v.Line && val.LineEnd >= v.Line {
				seg = val.Name.Str
			}
		}
		offset = asm.List.Lexeme[v.Line].Offset
		if asm.List.Lexeme[v.Line].Type == "LT_CODE" {
			offset = asm.LastVOffset
		}
		typ := ""
		switch v.Type.TokensType {
		case "TT_DB_DIRECTIVE":
			typ = "BYTE"
		case "TT_DW_DIRECTIVE":
			typ = "WORD"
		case "TT_DD_DIRECTIVE":
			typ = "DWORD"
		}
		res += fmt.Sprintf("| %2d |   %-6s   |  %-6s  | %4s:%3d |\n", i, v.Name.Str, typ, seg, offset)
		i++
	}
	for _, v := range asm.Labels {
		seg := ""
		for _, val := range asm.Segments {
			if val.LineStart <= v.Line && val.LineEnd >= v.Line {
				seg = val.Name.Str
			}
		}
		res += fmt.Sprintf("| %2d |   %-6s   |  %-6s  | %4s:%3d |\n", i, v.Label.Str, "LABEL", seg, asm.List.Lexeme[v.Line].Offset)
		i++
	}
	res += "\\----------------------------------------/\n\n"

	_, err = file.WriteString(res)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic("PrintFirstStage: Can't close file for write")
	}
}
