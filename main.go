package main

import (
	"bufio"
	"fmt"
	"letters/solver"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

var uniqueLines = make(map[string]struct{})

func loadWords(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Download dictionary failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Download dictionary failed. Status Code:", resp.StatusCode)
		return
	}
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			uniqueLines[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read dictionary error:", err)
		return
	}
}

func selectLinesByLength(length int) []string {
	var result []string
	for line := range uniqueLines {
		if utf8.RuneCountInString(line) == length {
			result = append(result, line)
		}
	}
	return result
}

func writeToFile(lines []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
func readFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

func main() {

	url := "https://raw.githubusercontent.com/LussRus/Rus_words/master/cp1251/txt/nouns/summary.txt"
	letters5TmpPath := "/var/tmp/letters5"
	var letters5 []string
	if _, err := os.Stat(letters5TmpPath); err == nil {
		letters5 = readFromFile(letters5TmpPath)
	} else {
		loadWords(url)
		fmt.Printf("Unique lines: %d\n", len(uniqueLines))
		letters5 = selectLinesByLength(5)
		fmt.Printf("Number of lines of 5 letters: %d\n", len(letters5))
		err = writeToFile(letters5, letters5TmpPath)
		if err != nil {
			fmt.Println("Failed to save temp file")
		}
	}

	notContains := []rune{'б', 'у', 'к', 'в', 'п', 'о', 'з', 'е', 'д', 'с'}
	contain := []rune{'р', 'и', 'а'}
	correctPositions := []solver.RunePlace{
		{Rune: 'а', Pos: 4},
	}

	s := solver.NewSolver(letters5, 5)
	s.Contains(contain)
	s.NotContain(notContains)
	s.CorrectRunePlaces(correctPositions)

	results := s.GetSuitable()
	fmt.Println("Matching lines:")
	if len(results) == 0 {
		fmt.Println("not found")
	} else {
		for _, el := range results {
			fmt.Println(el)
		}
	}

}
