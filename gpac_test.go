package gpac

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	input := `+abc`
	abc := Or(
		Char('a'),
		Char('b'),
		Char('c'),
	)
	pattern := And(
		Optional(Char('+')),
		abc,
		abc,
		abc,
	)
	result := pattern([]byte(input))
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	fmt.Println(string(result.Ok))
}
