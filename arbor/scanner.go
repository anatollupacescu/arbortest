package arbor

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func readFile(target string) string {
	var b strings.Builder

	f, err := os.Open(target)

	if err != nil {
		log.Fatalf("open file: %v\n", err)
	}

	s := bufio.NewScanner(f)

	for s.Scan() {
		b.WriteString(s.Text())
		b.WriteString("\n")
	}

	return b.String()
}
