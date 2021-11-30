package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

var binName = "todo"

const (
	envVarName = "TODO_FILENAME"
	fileName   = "test-todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	if err := os.Setenv(envVarName, fileName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set env variable %s", envVarName)
		os.Exit(1)
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	if err := os.Remove(binName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove binary %s", binName)
	}
	if err := os.Remove(fileName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove temp file %s", fileName)
	}
	if err := os.Unsetenv(envVarName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unset env variable %s", envVarName)
	}
	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task1 := "Test task number 1"
	task2 := "Test task number 2"
	task3 := "Test task number 3"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)
	t.Run("AddNewTaskFromArgs", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task1)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmdAdd := exec.Command(cmdPath, "-add", task3)
		if err := cmdAdd.Run(); err != nil {
			t.Fatal(err)
		}

		cmdComplete := exec.Command(cmdPath, "-complete", "3")
		if err := cmdComplete.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListAllTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\nX 2: %s\n", task2, task3)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("ListActiveTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-active")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n", task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("ListActiveTasksVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-active", "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		timestampPattern := regexp.MustCompile(`\d+ \w{3} \d+ \d\d:\d\d`)
		matched := timestampPattern.FindAll(out, -1)
		if len(matched) != 1 {
			t.Errorf("Expected 1 datetime strings, got %d\n", len(matched))
		}
	})

	t.Run("ListAllTasksVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		timestampPattern := regexp.MustCompile(`\d+ \w{3} \d+ \d\d:\d\d`)
		matched := timestampPattern.FindAll(out, -1)
		if len(matched) != 3 {
			t.Errorf("Expected 3 datetime strings, got %d\n", len(matched))
		}
	})
}
