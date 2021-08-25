package main

import (
	"os"
	"testing"
)

func TestMainProgram(t *testing.T) {
	os.Args = []string{"akash"}
	main()
}
