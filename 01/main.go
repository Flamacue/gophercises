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

// TODO(asumner): handle errors in a better manner than panic

// A config holds the configuration values for the quiz game
type config struct {
	timeLimit        int    // maximum time limit for the quiz game
	problemsFilePath string // path to the file containing the quiz questions
	random           bool
}

type problem struct {
	problem string
	answer  string
}

// Constants
const (
	timeLimitDefault        = 30                  // default time limit
	problemsFilePathDefault = "/etc/problems.csv" // default path to problems file
	randomDefault           = false
)

const welcomeText = `Welcome to the Quiz Show!

The quiz will have a time limit of %d seconds

We will be presenting problems located from %s

Do you accept this challenge? [y/N]
`

func beginQuiz(problems []*problem, timeLimit int) {
	totalProblems := len(problems)
	correct := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

Quiz:
	for i, problem := range problems {
		fmt.Printf("Problem %d: %s\n", i, problem.problem)

		answerChan := make(chan string)
		go func() {
			var answer string
			_, err := fmt.Scanf("%s", &answer)
			if err != nil {
				panic(err)
			}
			answerChan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Times up! Thanks for playing!")
			break Quiz
		case answer := <-answerChan:
			if answer == problem.answer {
				correct++
				fmt.Println("Correct!")
			} else {
				fmt.Println("Incorrect!")
			}
		}
	}

	fmt.Printf("Total Problems: %d, Number correct: %d\n", totalProblems, correct)
}

// Parses cmd line flags and returns a config
func buildConfiguration() *config {
	var timeLimit int
	var problemsFilePath string
	var random bool

	flag.IntVar(&timeLimit, "time", timeLimitDefault, "Time limit for the quiz game in seconds")
	flag.IntVar(&timeLimit, "t", timeLimitDefault, "Time limit for the quiz game in seconds")
	flag.StringVar(&problemsFilePath, "file", problemsFilePathDefault, "Time limit for the quiz game in seconds")
	flag.StringVar(&problemsFilePath, "f", problemsFilePathDefault, "Time limit for the quiz game in seconds")
	flag.BoolVar(&random, "random", randomDefault, "Randomize the question ordering")
	flag.BoolVar(&random, "r", randomDefault, "Randomize the question ordering")
	flag.Parse()

	return &config{timeLimit, problemsFilePath, random}
}

// Main entry point to our program
// Executes the main logic loop
func main() {
	conf := buildConfiguration()
	fmt.Printf(welcomeText, conf.timeLimit, conf.problemsFilePath)
	var accept string
	_, err := fmt.Scanf("%s", &accept)
	if err != nil {
		panic(err)
	}

	if strings.ToUpper(accept) != "Y" {
		fmt.Println("Not so brave now are you?")
		os.Exit(0)
	}

	fmt.Println() // Newline for formatting
	problems := parseProblems(conf.problemsFilePath)
	if conf.random {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}
	beginQuiz(problems, conf.timeLimit)
	os.Exit(0)
}

// This is a naive file read function that will load an entire csv file's contents
// into memory and return a slice of problems
func parseProblems(filePath string) []*problem {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	var problems []*problem
	for _, line := range lines {
		problems = append(problems, &problem{line[0], strings.TrimSpace(line[1])})
	}

	return problems
}
