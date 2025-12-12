package compiler

import "github.com/alecthomas/participle/v2/lexer"

var LoLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Keyword", `\b(lo|for|return|struct|new)\b`},
	{"Type", `\b(int|string|bool)\b`},
	{"Ident", `[a-zA-Z_]\w*`},
	{"Number", `\d+`},
	{"String", `"(?:[^"\\]|\\.)*"`},
	{"Operator", `==|!=|<=|>=|[+\-*/<>=]`},
	{"Punct", `[(){}.,;:]`},
	{"Comment", `//[^\n]*`},
	{"Whitespace", `\s+`},
})
