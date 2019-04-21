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
	"flag"
)

type quiz struct {
	csv_filename string
	questions []string
	answers []string
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

func (q *quiz) timeQuiz(sec int, c chan bool){
	timer := time.NewTimer(time.Duration(sec) * time.Second)
	<- timer.C
	c <- true
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
	
	score, asked := 0, 0
	var correct bool
	if timed {
		time_up_c := make(chan bool)
		go q.timeQuiz(30, time_up_c)
		for {
			select {
			case <- time_up_c:
				time_up_c = nil

			case correct, ok := <-c:
				if !ok {
					c = nil
				}
				asked++
				if correct {
					score++
				}
			}
			if c == nil || time_up_c == nil {
				break
			}
		}
	} else {
		for correct = range c {
			asked++
			if correct {
				score++
			}
		}
	}

	fmt.Printf("%d of %d questions asked. %d questions were answered correctly.", asked, len(q.questions), score)

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

	csvPtr := flag.String("csv", "test.csv", "Specify path to csv of questions and answers.")
	randomPtr := flag.Bool("random", false, "Randomize questions?")
	timedPtr := flag.Bool("timed", false, "Time quiz?")

	flag.Parse()

	Quiz := quiz{csv_filename: *csvPtr}

	err := Quiz.readCsv()
	if err != "ok" {
		fmt.Println("Error reading CSV.")
		os.Exit(1)
	}

	fmt.Println(*timedPtr)
	fmt.Println(*randomPtr)
	

	err = Quiz.administerQuiz(*timedPtr, *randomPtr)
	if err != "ok" {
		fmt.Println("Error administering quiz.")
		os.Exit(1)
	}
	os.Exit(0)
}
