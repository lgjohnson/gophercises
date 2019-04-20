package main

import (
	"bufio"
	"os"
	"io"
	"encoding/csv"
	"fmt"
	"strings"
)

const (
	CSV_FILENAME = "test.csv"
)

type quiz struct {
	csv_filename string
	questions []string
	answers []string
	score int
}

func (q *quiz) readCsv() string {

	if q.csv_filename == "" {
		return ""
	}

	csvFile, _ := os.Open(q.csv_filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	q.questions = []string{}
	q.answers = []string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		q.questions = append(q.questions, record[0])
		q.answers = append(q.answers, record[1])
	}

	return "ok"
}

func (q *quiz) administer() string {

	q.score = 0
	stdin_reader := bufio.NewReader(os.Stdin)

	for i, question := range q.questions {
		fmt.Println(question)
		text, _ := stdin_reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == q.answers[i] {
			q.score++
		}
	}

	fmt.Printf("%d questions asked, %d questions were answered correctly.", len(q.questions), q.score)
	return "ok"
}

func main() {

	Quiz := quiz{csv_filename: CSV_FILENAME}

	err := Quiz.readCsv()
	if err != "ok" {
		fmt.Println("Error reading CSV.")
		os.Exit(1)
	}

	err = Quiz.administer()
	if err != "ok" {
		fmt.Println("Error administering quiz.")
		os.Exit(1)
	}
	os.Exit(0)
}


