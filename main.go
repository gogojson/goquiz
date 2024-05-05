package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	// Init flags
	fn := flag.String("csv", "problems.csv", "CSV file path format (question, answer)")
	timer := flag.Int("timer", 30, "Set timer for each question, default value (30seconds)")
	shuffle := flag.Bool("shuffle", false, "Shuffle question order, default value: false")
	flag.Parse()

	// Open file and decode csv data
	f, err := os.Open(*fn)
	if err != nil {
		fmt.Printf("Failed to open file: '%s', %v\n", *fn, err)
		os.Exit(1)
	}
	data := csv.NewReader(f)
	rec, err := data.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read csv file, %v\n", err)
		os.Exit(1)
	}

	pLen := len(rec)
	correct := 0
	ch := make(chan int)

	p := make([]problem, pLen)
	if err := csvParser(rec, p); err != nil {
		fmt.Printf("Failed to parse CSV file / %e\n", err)
		os.Exit(1)
	}

	// Start timer after confirm
	fmt.Println("after chanel insert")
	fmt.Print("Press 'enter' to start quiz and timer")
	fmt.Scanln()
	timeLimit := time.NewTimer(time.Second * time.Duration(*timer))
	defer timeLimit.Stop()

	// Make it rand order if toggled
	ps := make([]int, pLen)
	switch {
	case *shuffle:
		ps = rand.Perm(pLen)
	default:
		for i := range pLen {
			ps[i] = i
		}
	}
L:
	for i, pNum := range ps {
		fmt.Printf("Q%d. %s=", i+1, p[pNum].q)
		// time.NewTimer(time.Second * time.Duration(*timer))

		go func() {
			var a string
			fmt.Scan(&a)
			ch <- 1
			if strings.TrimSpace(a) == p[pNum].a {
				correct++
			}
		}()
		select {
		case <-ch:
			// No matter time after chanel is empty
		case <-timeLimit.C:
			fmt.Println("\nPass time limit")
			break L
		}
	}

	fmt.Printf("Quiz Finished got %d out of %d\n", correct, len(rec))
	fmt.Printf("Score: %d%%\n", (correct*100)/len(rec))
}

func csvParser(input [][]string, output []problem) error {
	for i, row := range input {
		if len(row) != 2 {
			fmt.Printf("Invalid CSV format row: %v", row)
			return fmt.Errorf(fmt.Sprintf("Invalid CSV format row: %v", row))
		}
		output[i].q = row[0]
		output[i].a = strings.TrimSpace(row[1])
	}
	return nil
}
