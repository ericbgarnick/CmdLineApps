package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// CountUnit indicates the unit of text for which input text should be split.
type CountUnit string

// CountWords indicates that input text should be counted by whitespace-separated words.
// CountLines indicates that input text should be counted by newline-separated lines.
// CountBytes indicates that input text should be counted by bytes.
const (
	CountWords CountUnit = "w"
	CountLines CountUnit = "l"
	CountBytes CountUnit = "b"
)

func main() {
	countUnit := flag.String("unit", string(CountWords), "Count unit: 'w' for words, 'l' for lines, 'b' for bytes")
	sourceFiles := flag.String("files", "", "Source files for text to process. Filenames must be grouped within double-quotes")
	flag.Parse()

	result := 0
	var err error
	if *sourceFiles != "" {
		fileNames := strings.Split(*sourceFiles, " ")
		result, err = countFileContents(fileNames, CountUnit(*countUnit))
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		result = count(os.Stdin, CountUnit(*countUnit))
	}

	fmt.Println(result)
}

func countFileContents(fileNames []string, countUnit CountUnit) (int, error) {
	totalCount := 0
	for _, fileName := range fileNames {
		f, err := os.Open(fileName)
		if err != nil {
			return totalCount, err
		}
		totalCount += count(f, countUnit)
		f.Close()
	}
	return totalCount, nil
}

func count(r io.Reader, countUnit CountUnit) int {
	scanner := bufio.NewScanner(r)
	if countUnit == CountWords {
		scanner.Split(bufio.ScanWords)
	} else if countUnit == CountBytes {
		scanner.Split(bufio.ScanBytes)
	}
	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}
