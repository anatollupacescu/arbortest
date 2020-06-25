package arbor

import (
	"io/ioutil"
	"log"
)

var src = `package arbor

import "testing"

func TestArbor(t *testing.T) {
	t.Fail()
}
`

func Write(fileName string) {
	err := ioutil.WriteFile(fileName, []byte(src), 0777)

	if err != nil {
		log.Fatalf("create generated test file: %v", err)
	}
}

func GenerateTest(s Suite) string {
	return src
}
