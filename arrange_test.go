package goarrange

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func formatCode(t *testing.T, src string) string {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", []byte(src), parser.ParseComments)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	require.NoError(t, format.Node(buf, fSet, f))

	return buf.String()
}

func TestArrange(t *testing.T) {
	src := `package bita

// D comment
func D() {}

func A() {}

// B Comment 
func B() {}
`

	res, err := Arrange([]byte(src))
	require.NoError(t, err)
	assert.Equal(t, formatCode(t, `package bita

func A() {}

// B Comment 
func B() {}

// D comment
func D() {}
`), string(res))
}

func TestArrange2(t *testing.T) {
	src := `package bita

// D comment
func (*testing) D() {}

func A() {}

// B Comment 
func (a Alpha) B() {}

func (a Alpha) A() {}
`

	res, err := Arrange([]byte(src))
	require.NoError(t, err)
	assert.Equal(t, formatCode(t, `package bita

func (a Alpha) A() {}

// B Comment 
func (a Alpha) B() {}

// D comment
func (*testing) D() {}

func A() {}
`), string(res))
}
