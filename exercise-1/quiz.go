package main

import (
	"bufio"
	"os"
	"io"
	"math/rand"
	"time"
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

func (q *quiz ) administerQuestions(numbers []int, c chan bool) {

	stdin_reader := bufio.NewReader(os.Stdin)
	for _, i := range numbers {
		fmt.Println(q.questions[i])
		text, _ := stdin_reader.ReadString('\n')
		text = strings.TrimSpace(text)
		c <- text == q.answers[i]
	}
	close(c)

}

func (q *quiz) administerQuiz(timed, random bool) string {

	var numbers []int
	for i := range q.questions {
		numbers = append(numbers, i)
	}

	if random {
		Shuffle(numbers)
	} 

	c := make(chan bool)
	go q.administerQuestions(numbers, c)
	q.score = 0

	if timed {
		fmt.Println("Time not implemented yet :(")	
	} else {
		for correct := range c {
			if correct {
				q.score++
			}
		}
		fmt.Println(q.score)
	}

	fmt.Printf("%d questions asked, %d questions were answered correctly.", len(q.questions), q.score)

	return "ok"
}

func Shuffle(slice []int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}


func main() {

	Quiz := quiz{csv_filename: CSV_FILENAME}

	err := Quiz.readCsv()
	if err != "ok" {
		fmt.Println("Error reading CSV.")
		os.Exit(1)
	}

	err = Quiz.administerQuiz(false, true)
	if err != "ok" {
		fmt.Println("Error administering quiz.")
		os.Exit(1)
	}
	os.Exit(0)
}


