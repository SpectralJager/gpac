package parsers

import (
	"bytes"

	"github.com/SpectralJager/gpac"
)

func Match(pattern string) gpac.ParseFunc[string] {
	parsers := []gpac.ParseFunc[byte]{}
	buff := bytes.NewBufferString(pattern)
	for {
		ch, err := buff.ReadByte()
		if err != nil {
			break
		}
		parsers = append(parsers, gpac.Char(ch))
	}
	return gpac.Map(gpac.And(parsers...), func(in []byte) string {
		acc := ""
		for _, ch := range in {
			acc += string(ch)
		}
		return acc
	})
}

func Integer() gpac.ParseFunc[string] {
	return gpac.Map(
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
		func(num []byte) string {
			return string(num)
		},
	)
}

func SignedInteger() gpac.ParseFunc[string] {
	return gpac.Map(
		gpac.And(
			gpac.Map(
				gpac.Optional(gpac.Or(
					gpac.Char('+'),
					gpac.Char('-'),
				)),
				func(in byte) string { return string(in) },
			),
			Integer(),
		),
		func(num []string) string {
			if num[0] != "-" {
				return num[1]
			}
			return num[0] + num[1]
		},
	)
}
