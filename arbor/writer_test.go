package arbor_test

import (
	"testing"

	"github.com/anatollupacescu/sandbox/arbor"
	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	s := arbor.Suite{
		"test1": []string{"provider"},
	}

	output := arbor.GenerateTest(s)

	expected := `package arbor

import "testing"

func TestArbor(t *testing.T) {
	t.Fail()
}
`
	assert.Equal(t, expected, output)
}
