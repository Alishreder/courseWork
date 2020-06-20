package main

import (
	"courseWork/assembly"
	"courseWork/lexeme"
	"courseWork/list"
	"courseWork/token"
	"fmt"
	"os"
	"strings"
)

func isWhiteSpace(c string) bool {
	return strings.ContainsAny(c, "\n \t")
}

func isDelimiter(c string) bool {
	return isWhiteSpace(c) || strings.ContainsAny(c, ",*[]+-=:")
}

func Tokenize(str string) (list list.TypeList) {

	length := len(str)
	if length == 0 {
		panic("Tokenize: Passed empty file")
	}

	var buffer string
	var buffLex string
	lineCounter := 0
	var currLexeme lexeme.TypeLexeme
	currLexeme.Line = lineCounter

	for _, v := range str {

		c := string(v)
		buffLex += c

		if isDelimiter(c) {
			trimmedBuffer := strings.TrimSpace(buffer)
			if len(trimmedBuffer) != 0 {
				tk := token.CreateToken(trimmedBuffer)
				currLexeme.Tokens = append(currLexeme.Tokens, tk)
				if tk.TokensType == "TT_UNKNOWN" {
					currLexeme.ErrStr = "Unknown tokens type"
				}
				currLexeme.TokensCounter++
			}

			buffer = c
			trimmedBuffer = strings.TrimSpace(buffer)
			if len(trimmedBuffer) != 0 {
				tk := token.CreateToken(trimmedBuffer)
				currLexeme.Tokens = append(currLexeme.Tokens, tk)
				if tk.TokensType == "TT_UNKNOWN" {
					currLexeme.ErrStr = "Unknown tokens type"
				}
				currLexeme.TokensCounter++
			}
			buffer = ""
		} else if c == "." {
			buffer += strings.ToUpper(c)
		} else {
			buffer += strings.ToUpper(c)
		}

		if c == "\n" {
			lineCounter++
			if currLexeme.TokensCounter != 0 {
				currLexeme.Line = lineCounter
				currLexeme.Source = strings.TrimSpace(buffLex)
				list.Lexeme = append(list.Lexeme, currLexeme)
				currLexeme.ClearLexeme()
			}
			buffLex = ""
		}
	}

	// Last token
	trimmedBuffer := strings.TrimSpace(buffer)
	if len(trimmedBuffer) != 0 {
		tk := token.CreateToken(trimmedBuffer)
		currLexeme.Tokens = append(currLexeme.Tokens, tk)
		if tk.TokensType == "TT_UNKNOWN" {
			currLexeme.ErrStr = "Unknown tokens type"
		}
		currLexeme.TokensCounter++
	}
	if currLexeme.TokensCounter != 0 {
		currLexeme.Line++
		currLexeme.Source = strings.TrimSpace(buffLex)
		list.Lexeme = append(list.Lexeme, currLexeme)
		currLexeme.ClearLexeme()
	}

	return
}

func PrintFirstStage(asm *assembly.TypeAsm) {
	file, err := os.Create("testFirstStage.docx")
	if err != nil {
		panic("PrintFirstStage: Can't create file for write")
	}

	res := ""

	for _, v := range asm.List.Lexeme {
		res = res + "\n" + v.Source

		res += "\n"
		res += "/----------------------------------------\\\n"
		res += "| N | String        | Type               |\n"
		res += "|---+---------------+--------------------|\n"
		for i, val := range v.Tokens {
			res += fmt.Sprintf("| %-2v| %-14v| %-19v|\n", i, val.Str, val.TokensType)
		}
		res += "|----------------------------------------|\n"
		res += "| NAME | MNEM      | OPERAND1 | OPERAND2 |\n"
		res += "|------+-----------+----------+----------|\n"

		buff1 := ""
		if v.HasLabel {
			buff1 = "0"
		} else {
			buff1 = "-"
		}

		buff2 := fmt.Sprintf("%-10s", "-")
		buff3 := fmt.Sprintf("%-9s", "-")
		buff4 := fmt.Sprintf("%-9s", "-")
		if v.HasInstruction {
			buff2 = fmt.Sprintf("%-6d%-4d", v.InstructionIndex, 1)
			if v.OperandsCount >= 1 {
				buff3 = fmt.Sprintf("%-5d %-3d", v.Operands[0].OperandIndex, v.Operands[0].OperandLength)
			}
			if v.OperandsCount > 1 {
				buff4 = fmt.Sprintf("%-5d %-3d", v.Operands[1].OperandIndex, v.Operands[1].OperandLength)
			}
		}

		res += fmt.Sprintf("|%4s  | %s| %s| %s|\n", buff1, buff2, buff3, buff4)
		res += "\\----------------------------------------/\n"
	}

	_, err = file.WriteString(res)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic("PrintFirstStage: Can't close file for write")
	}
}

func FirstStage(asm *assembly.TypeAsm) error {
	errStr := ""
	if asm == nil {
		errStr += fmt.Sprintf("FirstStage: Passed nil argument")
	}

	asm.List = Tokenize(asm.Content)
	asm.LexemeStructure()

	if errStr != "" {
		return fmt.Errorf(errStr)
	}
	return nil
}
