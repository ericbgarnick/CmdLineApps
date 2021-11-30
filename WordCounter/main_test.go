package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	binName   = "WordCounter"
	testFile1 = "./testdata/test1.txt"
	testFile2 = "./testdata/test2.txt"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")
	exp := 4
	res := count(b, CountWords)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("one two three\nline2\nfour five six")
	exp := 3
	res := count(b, CountLines)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.", exp, res)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("How many bytes?")
	exp := 15
	res := count(b, CountBytes)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.", exp, res)
	}
}

func TestCountFromFile(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)
	fileNames := []string{testFile1, testFile2}
	cmd := exec.Command(cmdPath, "-unit", "b", "-files", strings.Join(fileNames, " "))

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedCount := 24
	if string(bytes.TrimSpace(out)) != strconv.Itoa(expectedCount) {
		t.Errorf("Expected %d, got %s instead\n", expectedCount, out)
	}
}
