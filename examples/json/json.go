package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SpectralJager/gpac"
)

const (
	ILLEGAL_AST byte = iota
	INT_AST
	STR_AST
	DIC_AST
)

type ast struct {
	kind byte
	i    int
	s    string
	d    []kv
}

type kv struct {
	k string
	v ast
}

func (node ast) String() string {
	switch node.kind {
	case INT_AST:
		return strconv.Itoa(node.i)
	case STR_AST:
		return fmt.Sprintf("\"%s\"", node.s)
	case DIC_AST:
		buf := bytes.Buffer{}
		items := []string{}
		for _, kv := range node.d {
			items = append(items, fmt.Sprintf("\"%s\":%s", kv.k, kv.v.String()))
		}
		fmt.Fprintf(&buf, "{%s}", strings.Join(items, ", "))
		return buf.String()
	default:
		return ""
	}
}

var (
	// whitespace = [\r\t\n ]+
	whitespaceParser = gpac.Map(
		gpac.ManyOrOne(
			gpac.Or(
				gpac.Char('\r'),
				gpac.Char('\n'),
				gpac.Char('\t'),
				gpac.Char(' '),
			),
		),
		func(in []byte) (byte, error) { return 0, nil },
	)
	// integer = [0-9]+
	integerParser = gpac.Map(
		gpac.ManyOrOne(
			gpac.Or(
				gpac.Char('0'),
				gpac.Char('1'),
				gpac.Char('2'),
				gpac.Char('3'),
				gpac.Char('4'),
				gpac.Char('5'),
				gpac.Char('6'),
				gpac.Char('7'),
				gpac.Char('8'),
				gpac.Char('9'),
			),
		),
		func(bytes []byte) (ast, error) {
			i, err := strconv.Atoi(string(bytes))
			if err != nil {
				return ast{}, err
			}
			node := ast{
				kind: INT_AST,
				i:    i,
			}
			return node, nil
		},
	)
	// string = '"' [a-z0-9_ ]* '"'
	stringParser = gpac.Map(
		gpac.And(
			gpac.Map(gpac.Char('"'), func(in byte) (string, error) { return string(in), nil }),
			gpac.Map(
				gpac.Many(
					gpac.Or(
						gpac.Char(' '),
						gpac.Char('_'),
						gpac.Char('0'),
						gpac.Char('1'),
						gpac.Char('2'),
						gpac.Char('3'),
						gpac.Char('4'),
						gpac.Char('5'),
						gpac.Char('6'),
						gpac.Char('7'),
						gpac.Char('8'),
						gpac.Char('9'),
						gpac.Char('q'),
						gpac.Char('w'),
						gpac.Char('e'),
						gpac.Char('r'),
						gpac.Char('t'),
						gpac.Char('y'),
						gpac.Char('i'),
						gpac.Char('o'),
						gpac.Char('p'),
						gpac.Char('a'),
						gpac.Char('s'),
						gpac.Char('d'),
						gpac.Char('f'),
						gpac.Char('g'),
						gpac.Char('h'),
						gpac.Char('k'),
						gpac.Char('l'),
						gpac.Char('z'),
						gpac.Char('x'),
						gpac.Char('c'),
						gpac.Char('v'),
						gpac.Char('b'),
						gpac.Char('n'),
						gpac.Char('m'),
					),
				),
				func(in []byte) (string, error) { return string(in), nil },
			), gpac.Map(gpac.Char('"'), func(in byte) (string, error) { return string(in), nil }),
		),
		func(in []string) (ast, error) { return ast{kind: STR_AST, s: in[1]}, nil },
	)
	// value= string | integer
	valueParser = gpac.Or(
		integerParser,
		stringParser,
	)
	// kv	= string ':' value
	kvParser = gpac.Map(
		gpac.And(
			gpac.Map(gpac.Optional(whitespaceParser), func(in byte) (ast, error) { return ast{}, nil }),
			stringParser,
			gpac.Map(gpac.Char(':'), func(in byte) (ast, error) { return ast{}, nil }),
			valueParser,
		),
		func(in []ast) (kv, error) {
			return kv{k: in[1].s, v: in[3]}, nil
		},
	)
	// mkv = (kv ',')+ kv
	mkv = gpac.Map(
		gpac.And(
			gpac.Map(
				gpac.And(
					gpac.ManyOrOne(
						gpac.Map(
							gpac.And(
								kvParser,
								gpac.Map(gpac.Char(','), func(in byte) (kv, error) { return kv{}, nil }),
							),
							func(in []kv) (kv, error) { return in[0], nil }),
					),
					gpac.Map(kvParser, func(in kv) ([]kv, error) { return []kv{in}, nil }),
				),
				func(in [][]kv) ([]kv, error) { return append(in[0], in[1]...), nil },
			),
		),
		func(in [][]kv) ([]kv, error) { return in[0], nil },
	)
	// json = '{' (mkv | kv) '}'
	jsonParser = gpac.Map(
		gpac.Map(
			gpac.And(
				gpac.Map(gpac.Char('{'), func(in byte) ([]kv, error) { return []kv{}, nil }),
				gpac.Map(whitespaceParser, func(in byte) ([]kv, error) { return []kv{}, nil }),
				gpac.Or(
					mkv,
					gpac.Map(kvParser, func(in kv) ([]kv, error) { return []kv{in}, nil }),
				),
				gpac.Map(whitespaceParser, func(in byte) ([]kv, error) { return []kv{}, nil }),
				gpac.Map(gpac.Char('}'), func(in byte) ([]kv, error) { return []kv{}, nil }),
			),
			func(in [][]kv) ([]kv, error) { return in[2], nil },
		),
		func(in []kv) (ast, error) {
			return ast{kind: DIC_AST, d: in}, nil
		},
	)
)

func main() {
	input := `{
	"int":43,
	"string":"hello world"
}`
	result := jsonParser([]byte(input))
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Println(result.Ok.String())
}
