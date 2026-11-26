package txt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupEmptyInput(t *testing.T) {
	for _, ls := range []string{
		``,
		"\n \n\t\t\n  ",
	} {
		block, _ := ParseBlock(ls, 0)
		assert.Nil(t, block)
	}
}

func TestParseBlock(t *testing.T) {
	for _, x := range []struct {
		text       string
		expect     []string // expected significant line contents
		expectHead int
		expectTail int
	}{
		// Single line
		{"a", []string{"a"}, 0, 0},
		{"\nfoo", []string{"foo"}, 1, 0},
		{"\n12345\n", []string{"12345"}, 1, 0},
		{"   \ntest ", []string{"test "}, 1, 0},
		{"   \na\ta\n", []string{"a\ta"}, 1, 0},
		{"\t\na1\n\t \t ", []string{"a1"}, 1, 1},
		{"\n\na1\n\n", []string{"a1"}, 2, 1},
		{"喜左衛門", []string{"喜左衛門"}, 0, 0},
		{"喜左衛門\n", []string{"喜左衛門"}, 0, 0},
		{"\n😀·½\n ", []string{"😀·½"}, 1, 1},

		// Invalid UTF-8 (e.g. a Latin-1 encoded file) at the very end
		{"caf\xe9", []string{"caf\xe9"}, 0, 0},
		{"\xe9", []string{"\xe9"}, 0, 0},
		{"a\nb\xe2\x82", []string{"a", "b\xe2\x82"}, 0, 0},

		// Multiple lines
		{"a1\na2", []string{"a1", "a2"}, 0, 0},
		{"\nasdf\nasdf", []string{"asdf", "asdf"}, 1, 0},
		{"\nHey 🥰!\n«How is it?»\n", []string{"Hey 🥰!", "«How is it?»"}, 1, 0},
		{"\n    \t\nA\nB", []string{"A", "B"}, 2, 0},
		{"\n    \t\na b c \n a b c\n  \t  \n", []string{"a b c ", " a b c"}, 2, 1},
		{"\n    \t\n       _       \n     -     \n\n", []string{"       _       ", "     -     "}, 2, 1},
		{" \t \t\nAS:FLKJH\n!(@* #&\n\t", []string{"AS:FLKJH", "!(@* #&"}, 1, 1},
		{" \n\t\n1—2\n·½⅓•ÄﬂÑ\n\n\n ", []string{"1—2", "·½⅓•ÄﬂÑ"}, 2, 3},
	} {
		b, _ := ParseBlock(x.text, 0)
		sgLines, head, tail := b.SignificantLines()

		require.NotNil(t, b)
		require.Len(t, b.Lines(), len(x.expect)+x.expectHead+x.expectTail)
		require.Len(t, sgLines, len(x.expect))
		for i, l := range x.expect {
			assert.Equal(t, l, sgLines[i].Text)
		}
		assert.Equal(t, x.expectHead, head)
		assert.Equal(t, x.expectTail, tail)
	}
}
