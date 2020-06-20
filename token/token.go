package token

import "unicode"

type TypeToken struct {
	Str string
	TokensType string
}

func (t *TypeToken) GetTokenType() string {
	return t.TokensType
}

func (t *TypeToken) GetFirstLetterOfToken() string {
	return string(t.Str[0])
}

var tokenTypes = []string {
	"TT_INSTRUCTION",

	"TT_PTR",
	"TT_BYTE",
	"TT_WORD",
	"TT_DWORD",

	"TT_DB_DIRECTIVE",
	"TT_DW_DIRECTIVE",
	"TT_DD_DIRECTIVE",

	"TT_MODEL_DIRECTIVE",
	"TT_DATA_DIRECTIVE",
	"TT_CODE_DIRECTIVE",
	"TT_END_DIRECTIVE",

	"TT_REGISTER_SEG",
	"TT_REGISTER32",
	"TT_REGISTER16",
	"TT_REGISTER8",
	"TT_SYMBOL",

	"TT_IDENTIFIER",

	"TT_NUMBER2",
	"TT_NUMBER10",
	"TT_NUMBER16",

	"TT_UNKNOWN",
}

var TTypes = map[string]string {
	"STD": "TT_INSTRUCTION",
	"DIV": "TT_INSTRUCTION",
	"IMUL": "TT_INSTRUCTION",
	"ADD": "TT_INSTRUCTION",
	"AND": "TT_INSTRUCTION",
	"CMP": "TT_INSTRUCTION",
	"SHL": "TT_INSTRUCTION",
	"MOV": "TT_INSTRUCTION",
	"JE": "TT_INSTRUCTION",
	"JMP": "TT_INSTRUCTION",

	"PTR": "TT_PTR",
	"BYTE": "TT_BYTE",
	"WORD": "TT_WORD",
	"DWORD": "TT_DWORD",

	"END": "TT_END_DIRECTIVE",
	".MODEL": "TT_MODEL_DIRECTIVE",
	".DATA": "TT_DATA_DIRECTIVE",
	".CODE": "TT_CODE_DIRECTIVE",

	"DB": "TT_DB_DIRECTIVE",
	"DW": "TT_DW_DIRECTIVE",
	"DD": "TT_DD_DIRECTIVE",

	"ES": "TT_REGISTER_SEG",
	"DS": "TT_REGISTER_SEG",
	"FS": "TT_REGISTER_SEG",
	"SS": "TT_REGISTER_SEG",
	"GS": "TT_REGISTER_SEG",
	"CS": "TT_REGISTER_SEG",

	"EAX": "TT_REGISTER32",
	"EBX": "TT_REGISTER32",
	"ECX": "TT_REGISTER32",
	"EDX": "TT_REGISTER32",
	"EDI": "TT_REGISTER32",
	"ESI": "TT_REGISTER32",
	"EBP": "TT_REGISTER32",
	"ESP": "TT_REGISTER32",

	"AX": "TT_REGISTER16",
	"BX": "TT_REGISTER16",
	"CX": "TT_REGISTER16",
	"DX": "TT_REGISTER16",
	"SI": "TT_REGISTER16",
	"DI": "TT_REGISTER16",
	"SP": "TT_REGISTER16",
	"BP": "TT_REGISTER16",

	"AL": "TT_REGISTER8",
	"AH": "TT_REGISTER8",
	"BL": "TT_REGISTER8",
	"BH": "TT_REGISTER8",
	"CL": "TT_REGISTER8",
	"CH": "TT_REGISTER8",
	"DL": "TT_REGISTER8",
	"DH": "TT_REGISTER8",

	":": "TT_SYMBOL",
	"[": "TT_SYMBOL",
	"]": "TT_SYMBOL",
	"*": "TT_SYMBOL",
	"+": "TT_SYMBOL",
	",": "TT_SYMBOL",
}

func isId(str string) bool {

	if unicode.IsDigit(rune(str[0])) {
		return false
	}

	return len(str) <= 7
}

func isBin(str string) bool {

	hasSign := string(str[0]) == "-"
	length := len(str)
	i := 0
	if hasSign {
		i = 1
	}
	c := ""
	for ; i < length - 1; i++ {
		c = string(str[i])
		if c != "0" && c != "1" {
			return false
		}
	}

	return string(str[length - 1]) == "B"
}

func isDec(str string) bool {

	hasSign := string(str[0]) == "-"
	length := len(str)
	i := 0
	if hasSign {
		i = 1
	}
	c := ""
	for ; i < length; i++ {
		c = string(str[i])
		if !(c >= "0" && c <= "9"){
			return false
		}
	}

	return true
}

func isHex(str string) bool {

	hasSign := string(str[0]) == "-"
	length := len(str)
	i := 0
	if hasSign {
		i = 1
	}
	c := ""
	for ; i < length - 1; i++ {
		c = string(str[i])
		if !((c >= "0" && c <= "9") || (c >= "A" && c<= "F")) {
			return false
		}
	}

	return string(str[length - 1]) == "H"
}

func CreateToken(buff string)  (t TypeToken) {

	t.Str = buff
	t.TokensType = "TT_UNKNOWN"

	if v, ok := TTypes[t.Str]; ok {
		t.TokensType = v
	}
	if t.TokensType == "TT_UNKNOWN" {
		if isHex(buff) {
			t.TokensType = "TT_NUMBER16"
		} else if isDec(buff) {
			t.TokensType = "TT_NUMBER10"
		} else if isBin(buff) {
			t.TokensType = "TT_NUMBER2"
		} else if isId(buff) {
			t.TokensType = "TT_IDENTIFIER"
		}
	}

	return
}
