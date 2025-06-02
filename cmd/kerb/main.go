package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Johniel/kerb/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  kerb is-kerb-file <file>")
		fmt.Println("  kerb sync [--remove-header] <srcDir>")
		fmt.Println("  kerb insert-header <file>")
		fmt.Println("  kerb add-header [<file>]")
		fmt.Println("  kerb replace <file> <old> <new>")
		fmt.Println("  kerb replace-all <old> <new>")
		fmt.Println("  kerb list")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "is-kerb-file":
		hasKerbHeaderCmd(os.Args[2:])
	case "sync":
		syncKerbFilesCmd(os.Args[2:])
	case "insert-header":
		insertKerbHeaderCmd(os.Args[2:])
	case "add-header":
		addKerbHeaderCmd(os.Args[2:])
	case "replace":
		replaceInKerbFileCmd(os.Args[2:])
	case "replace-all":
		replaceAllKerbFilesCmd(os.Args[2:])
	case "list":
		listKerbFilesCmd(os.Args[2:])
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func hasKerbHeaderCmd(args []string) {
	fs := flag.NewFlagSet("is-kerb-file", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb is-kerb-file <file>")
		os.Exit(1)
	}
	file := fs.Arg(0)
	has, err := internal.HasKerbHeader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	if has {
		fmt.Println("true")
		os.Exit(0)
	} else {
		fmt.Println("false")
		os.Exit(1)
	}
}

func syncKerbFilesCmd(args []string) {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	removeHeader := fs.Bool("remove-header", false, "Remove Kerb header after copying")
	fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb sync [--remove-header] <srcDir>")
		os.Exit(1)
	}
	srcDir := fs.Arg(0)
	dstDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	if err := internal.SyncKerbFiles(srcDir, dstDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	if *removeHeader {
		files, err := internal.ListKerbFiles(dstDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing kerb files: %v\n", err)
			os.Exit(2)
		}
		for _, file := range files {
			if err := internal.RemoveKerbHeader(file); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing header from %s: %v\n", file, err)
				continue
			}
			fmt.Printf("Kerb header removed from: %s\n", file)
		}
	}
}

func insertKerbHeaderCmd(args []string) {
	fs := flag.NewFlagSet("insert-header", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb insert-header <file>")
		os.Exit(1)
	}
	file := fs.Arg(0)
	has, err := internal.HasKerbHeader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	if has {
		fmt.Println("Kerb header already present.")
		return
	}
	if err := internal.InsertKerbHeader(file); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("Kerb header inserted.")
}

func addKerbHeaderCmd(args []string) {
	fs := flag.NewFlagSet("add-header", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb add-header [<file>]")
		os.Exit(1)
	}
	if fs.NArg() == 1 {
		file := fs.Arg(0)
		if err := internal.AddKerbHeader(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding header to %s: %v\n", file, err)
			os.Exit(2)
		}
		fmt.Printf("Kerb header added to: %s\n", file)
		return
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if err := internal.AddKerbHeader(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding header to %s: %v\n", path, err)
		} else {
			fmt.Printf("Kerb header added to: %s\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(2)
	}
}

func replaceInKerbFileCmd(args []string) {
	fs := flag.NewFlagSet("replace-in-kerb-file", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 3 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb replace <file> <old> <new>")
		os.Exit(1)
	}
	file := fs.Arg(0)
	old := fs.Arg(1)
	newStr := fs.Arg(2)
	has, err := internal.HasKerbHeader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	if !has {
		fmt.Println("Kerb header not present. No replacement performed.")
		return
	}
	if err := internal.ReplaceInKerbFile(file, old, newStr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("Replacement performed.")
}

func replaceAllKerbFilesCmd(args []string) {
	fs := flag.NewFlagSet("replace-all", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb replace-all <old> <new>")
		os.Exit(1)
	}
	old := fs.Arg(0)
	newStr := fs.Arg(1)
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	files, err := internal.ListKerbFiles(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	for _, file := range files {
		err := internal.ReplaceInKerbFile(file, old, newStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error replacing in %s: %v\n", file, err)
			continue
		}
		fmt.Printf("Replaced in: %s\n", file)
	}
}

func listKerbFilesCmd(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "Error: wrong number of arguments. Usage: kerb list")
		os.Exit(1)
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	files, err := internal.ListKerbFiles(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	for _, f := range files {
		fmt.Println(f)
	}
}
