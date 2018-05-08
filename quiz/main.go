package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	src := flag.String("csv", "problems.csv", "a csv file of the quiz in the format of 'question,answer'")
	timeout := flag.Int("timeout", 30, "the time allowed to answer to the quiz in seconds")
	flag.Parse()

	quiz, err := newQuizFromCSV(*src)
	if err != nil {
		log.Fatal(err)
	}

	score := 0
	s := bufio.NewScanner(os.Stdin)
	timer := time.NewTimer(time.Duration(*timeout) * time.Second)
	for k, p := range quiz {
		fmt.Printf("Problem #%d: %v = ", k+1, p.q)
		answerCh := make(chan string)
		go func() {
			s.Scan()
			answerCh <- strings.TrimSpace(s.Text())
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTimeout! You scored %d out of %d.\n", score, len(quiz))
			return
		case a := <-answerCh:
			if a == p.a {
				score++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", score, len(quiz))
}

func newQuizFromCSV(filename string) ([]problem, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open problems file: %s", err.Error())
	}

	c := csv.NewReader(f)
	r, err := c.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse problems file: %s", err.Error())
	}

	l := make([]problem, len(r))
	for k, v := range r {
		l[k] = problem{
			v[0],
			strings.TrimSpace(v[1]),
		}
	}

	return l, nil
}
