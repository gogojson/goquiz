package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	fn := flag.String("csv", "problems.csv", "CSV file path format (question, answer)")
	timer := flag.Int("timer", 30, "Set timer for each question, default value (30seconds)")
	flag.Parse()
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

	p := make([]problem, len(rec))
	if err := csvParser(rec, p); err != nil {
		fmt.Printf("Failed to parse CSV file / %e\n", err)
		os.Exit(1)
	}

	var correct int
L:
	for i, row := range p {
		ch := make(chan int)
		fmt.Print("Press 'enter' to start quiz and timer")
		fmt.Scanln()
		fmt.Printf("Q%d. %s=", i+1, row.q)

		go func() {
			var a string
			fmt.Scan(&a)
			ch <- 1
			if strings.TrimSpace(a) == row.a {
				correct++
			}
		}()
		select {
		case <-ch:
		case <-time.After(time.Second * time.Duration(*timer)):
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
