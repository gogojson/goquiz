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

var (
	userInput         string
	totalProblemCount int
	rightProblemCount int
	beenPlaced        []int
	fTimeInterval     = flag.Int("time", 10, "Time interval per question, default: 10s, limit: 30s")
	fFP               = flag.String("file", "problems.csv", "CSV file path, default: problems.csv")
	shuffle           = flag.Bool("shuffle", false, "Shuffle the quizes, default is false")
)

func main() {
	flag.Parse()

	fmt.Println(*fFP)
	f, err := os.Open(*fFP)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file `%s`, please check file path", *fFP))
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	rows, err := rdr.ReadAll()
	if err != nil {
		panic("It is not a valid csv file")
	}

	// Time limit
	if *fTimeInterval > 30 {
		*fTimeInterval = 30
	}

	totalProblemCount = len(rows)

	for i := range 10000 {
		var row = []string{}
		switch {
		case *shuffle:

			for {

				randN := rand.Intn(totalProblemCount)
				fmt.Printf("Current rand:%d\n", randN)
				fmt.Printf("Rand List:%v\n", beenPlaced)
			out:
				for i := range beenPlaced {
					if i == randN {
						continue out
					}
				}
				beenPlaced = append(beenPlaced, randN)
				row = rows[randN]
				break
			}
		default:
			row = rows[i]

		}

		ch := make(chan int)
		if len(row) < 2 {
			panic("It is not a valid csv file, it has to have at lest 2 columns")
		}

		fmt.Printf("Press Enter to start quiz #%d", i+1)
		fmt.Scanln()
		fmt.Printf("what is %s: ", row[0])

		go func() {
			fmt.Scan(&userInput)
			ch <- 1
		}()

		select {
		case <-ch:
			if userInput == strings.TrimSpace(row[1]) {
				rightProblemCount++
			}
		case <-time.After(time.Second * time.Duration(*fTimeInterval)):
			fmt.Println("\nPassed time limit!")
			return
		}
	}
	fmt.Printf("Your score: %d%%\n", (rightProblemCount*100)/totalProblemCount)
	fmt.Printf("Total Problems: %d    Correct Problems: %d\n", totalProblemCount, rightProblemCount)
}
