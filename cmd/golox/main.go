package main

import (
	"fmt"
	"os"

	"github.com/LightningDev1/golox"
)

const (
	EX_USAGE    = 64
	EX_DATAERR  = 65
	EX_SOFTWARE = 70
	EX_IOERR    = 74
)

func main() {
	runtime := golox.NewRuntime()

	switch len(os.Args) {
	case 1:
		if err := runtime.RunPrompt(); err != nil {
			fmt.Fprintf(os.Stderr, "REPL error: %v\n", err)
			os.Exit(1)
		}
	case 2:
		if err := runtime.RunFile(os.Args[1]); err != nil {
			if runtime.HadError {
				os.Exit(EX_DATAERR)
			}
			if runtime.HadRuntimeError {
				os.Exit(EX_SOFTWARE)
			}

			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(EX_IOERR)
		}
	default:
		fmt.Println("Usage: golox [script]")
		os.Exit(EX_USAGE)
	}
}
