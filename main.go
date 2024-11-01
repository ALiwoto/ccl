package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ALiwoto/ccl-gen/src/cclParser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided. Use 'generate' to generate code.")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		// Define flags
		fs := flag.NewFlagSet("generate", flag.ExitOnError)
		source := fs.String("source", "", "Path to the source file")
		language := fs.String("language", "", "Programming language for the generated code")
		output := fs.String("output", "", "Output path for the generated package")

		// Parse flags
		fs.Parse(os.Args[2:])

		// Validate flags
		if *source == "" || *language == "" || *output == "" {
			fmt.Println("Error: --source, --language, and --output flags are required when using 'generate'")
			fs.Usage()
			os.Exit(1)
		}

		// Handle generate command
		fmt.Printf("Generating code from %s in %s language to %s\n", *source, *language, *output)

		// check if the source file exists
		if _, err := os.Stat(*source); err != nil {
			switch {
			case os.IsNotExist(err):
				fmt.Printf("Error: source file %s does not exist\n", *source)
			case os.IsPermission(err):
				fmt.Printf("Error: permission denied for source file %s\n", *source)
			default:
				fmt.Printf("Error: failed to check source file %s: %v\n", *source, err)
			}
			os.Exit(1)
		}

		_, parseErr := cclParser.ParseCCLSourceFile(*source)
		if parseErr != nil {
			fmt.Printf("Error: failed to parse source file %s: %v\n", *source, parseErr)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("ccl version %s %s/%s\n", "1.0.0", runtime.GOOS, runtime.GOARCH)
	case "help":
		fmt.Println("ccl is a tool for managing ccl source code.\n\n" +
			"Usage:\n" +
			"\tccl <command> [arguments]\n\n" +
			"The commands are:\n\n" +
			"\tgenerate    Generate code from a ccl source file\n" +
			"\thelp        Show this help message",
		)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Use 'ccl help' for usage.")
		os.Exit(1)
	}
}