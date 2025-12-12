package compiler

import (
	participle "github.com/alecthomas/participle/v2"
)

var parser = participle.MustBuild[Program](
	participle.Lexer(LoLexer),
	participle.Elide("Whitespace", "Comment"),
	participle.Unquote("String"),
	participle.UseLookahead(5),
)

func ParseString(src string) (*Program, error) {
	prog, err := parser.ParseString("", src)
	if err != nil {
		return nil, err
	}

	return prog, nil
}
