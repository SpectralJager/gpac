package gpac

import (
	"bytes"
	"fmt"
)

type Result[a any] struct {
	Remaining []byte
	Ok        a
	Error     error
}

type ParseFunc[a any] func(input []byte) Result[a]

func Any[a any]() ParseFunc[a] {
	return func(input []byte) Result[a] {
		return Result[a]{Remaining: input}
	}
}

func None[a any]() ParseFunc[a] {
	return func(input []byte) Result[a] {
		return Result[a]{Remaining: input, Error: fmt.Errorf("none")}
	}
}

func Char(match byte) ParseFunc[byte] {
	return func(input []byte) Result[byte] {
		switch {
		case len(input) == 0:
			return Result[byte]{Remaining: input, Error: fmt.Errorf("empty input")}
		case input[0] == match:
			return Result[byte]{Remaining: input[1:], Ok: match}
		default:
			return Result[byte]{Remaining: input, Error: fmt.Errorf("char mismatched, expect %s, got %s", string(match), string(input[0]))}
		}
	}
}

func And[a any](patterns ...ParseFunc[a]) ParseFunc[[]a] {
	return func(input []byte) Result[[]a] {
		acc := []a{}
		for _, pattern := range patterns {
			result := pattern(input)
			if result.Error != nil {
				return Result[[]a]{Remaining: input, Error: result.Error}
			}
			input = result.Remaining
			acc = append(acc, result.Ok)
		}
		return Result[[]a]{Remaining: input, Ok: acc}
	}
}

func Or[a any](patterns ...ParseFunc[a]) ParseFunc[a] {
	return func(input []byte) Result[a] {
		for _, pattern := range patterns {
			result := pattern(input)
			switch {
			case result.Error != nil:
				continue
			default:
				return result
			}
		}
		return Result[a]{Remaining: input, Error: fmt.Errorf("can't match any of patterns")}
	}
}

func Optional[a any](pattern ParseFunc[a]) ParseFunc[a] {
	return Or(
		pattern,
		Any[a](),
	)
}

func Many[a any](pattern ParseFunc[a]) ParseFunc[[]a] {
	return func(input []byte) Result[[]a] {
		acc := []a{}
		for {
			result := pattern(input)
			if result.Error != nil {
				return Result[[]a]{Remaining: input, Ok: acc}
			}
			input = result.Remaining
			acc = append(acc, result.Ok)
		}
	}
}

func ManyOrOne[a any](pattern ParseFunc[a]) ParseFunc[[]a] {
	return func(input []byte) Result[[]a] {
		result := Many(pattern)(input)
		if len(result.Ok) != 0 {
			return result
		}
		return None[[]a]()(input)
	}
}

func Map[a, b any](pattern ParseFunc[a], mapper func(in a) (b, error)) ParseFunc[b] {
	return func(input []byte) Result[b] {
		result := pattern(input)
		if result.Error != nil {
			return Result[b]{Remaining: input, Error: result.Error}
		}
		val, err := mapper(result.Ok)
		if err != nil {
			return Result[b]{Remaining: input, Error: err}
		}
		return Result[b]{Remaining: result.Remaining, Ok: val}
	}
}

func Error[a any](pattern ParseFunc[a], callback func(Result[a]) error) ParseFunc[a] {
	return func(input []byte) Result[a] {
		result := pattern(input)
		if result.Error == nil {
			return result
		}
		err := callback(result)
		return Result[a]{Remaining: input, Error: err}
	}
}

func Match(pattern string) ParseFunc[bool] {
	parsers := []ParseFunc[byte]{}
	buff := bytes.NewBufferString(pattern)
	for {
		ch, err := buff.ReadByte()
		if err != nil {
			break
		}
		parsers = append(parsers, Char(ch))
	}
	return Map(And(parsers...), func(in []byte) (bool, error) {
		return true, nil
	})
}
