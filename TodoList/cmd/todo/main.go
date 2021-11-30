package main

import (
	todo "TodoList"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Default file name
var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for the Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "-add flag allows task name piped in through STDIN or included in the command after the flag")
	}
	// Parse command line flags
	task := flag.Bool("add", false, "Add task to the ToDo List")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Number of item to be marked complete")
	delItem := flag.Int("delete", 0, "Number of item to be removed from the ToDo list")
	verbose := flag.Bool("verbose", false, "Include more details when displaying item list")
	active := flag.Bool("active", false, "Omit completed items when displaying item list")

	flag.Parse()

	if filename := os.Getenv("TODO_FILENAME"); filename != "" {
		todoFileName = filename
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l.ListDisplay(*active, *verbose))
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		saveList(l)
	case *delItem > 0:
		if err := l.Delete(*delItem); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		saveList(l)
	case *task:
		tasks, err := getTasks(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, t := range tasks {
			l.Add(t)
		}
		saveList(l)
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// saveList saves the given List using its Save method and handles any resulting error.
func saveList(l *todo.List) {
	if err := l.Save(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// getTasks decides where to get the description for new tasks from: arguments or STDIN.
func getTasks(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	return getTasksFromStdin(r)
}

// getTasksFromStdin reads newline-separated tasks from STDIN and returns them in a slice of strings.
func getTasksFromStdin(r io.Reader) ([]string, error) {
	var result []string
	s := bufio.NewScanner(r)
	scanned := s.Scan()
	for scanned {
		if err := s.Err(); err != nil {
			return result, err
		}
		if len(s.Text()) == 0 {
			return result, fmt.Errorf("task cannot be blank")
		}
		result = append(result, s.Text())
		scanned = s.Scan()
	}

	return result, nil
}
