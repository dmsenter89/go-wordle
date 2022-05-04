// Package main(wordle) implements a simple wordle game.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to go-wordle!")

	word_list := load_dictionary()

	// pick a random word from list
	rand.Seed(time.Now().UnixNano())
	var word, guess, colored_comparison, new_round string
	var comparison []int

	// name the game loop so that we can break out of it
gameLoop:
	for {
		word = choose_random_word(word_list)
		fmt.Println("I picked a random 5-letter word. Try to guess it.")
		var correct = false

		for i := 1; i <= 6; i++ {
			fmt.Printf("Guess %d/6. ", i)
			guess = user_input()

			comparison = compare_answer(guess, word)
			colored_comparison = color_comparison(guess, comparison)
			fmt.Println(colored_comparison, " - your guess compared:", comparison)
			if guess == word {
				correct = true
				break
			}
		}

		if correct == true {
			fmt.Println("Congrats! You guessed the correct word.")
		} else {
			fmt.Println("You didn't guess the right word. It was", word)
		}

		fmt.Print("Press 'y(es)' for another round, else exit. ")
		fmt.Scanln(&new_round)
		switch strings.ToLower(new_round) {
		case "y", "yes":
			fmt.Println("---- new game ----")
			continue
		default:
			break gameLoop
		}
	}
}

// load_dictionary loads a wordle dictionary.
// If a dictionary is not found, one will be downloaded. The dictionary
// will be parsed to return a string slice consisting of 5-letter words only.
func load_dictionary() []string {
	// read wordle dict if present,
	body, err := ioutil.ReadFile("wordle.dict") // just pass the file name
	if err != nil {
		fmt.Printf("Dictionary not found. Error:\n\t%s.\n", err)
		fmt.Println("Downloading dictionary.")
		body, err = download_dictionary()
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body))) // f is the *os.File
	word_list := make([]string, 0, 20000)                        // pre-allocate room for 15k words
	var curr_word string

	for scanner.Scan() {
		curr_word = scanner.Text()
		if len(curr_word) == 5 {
			word_list = append(word_list, strings.ToUpper(curr_word))
		}
	}

	return word_list
}

// download_dictionary downloads a complete English language dictionary.
// Note that in the current state, the full-dictionary is saved and so
// load_dictionary must parse out the 5-letter words during each program start.
func download_dictionary() ([]byte, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// create tools to read and iterate over results
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("wordle.dict", body, 0644)
	return body, err
}

// choose_random_word uses a random int to find a
// random word in our word_list of 5-letter words
func choose_random_word(word_list []string) string {
	randidx := rand.Intn(len(word_list))
	return word_list[randidx]
}

// user_input gets a word guess from user, requires a length of 5.
// Does not validate if the user input is in the dictionary. Assumes ASCII answer,
// This does not work properly with unicode answers,
// but those wouldn't be in the dictionary we use in the first place.
func user_input() string {
	var guess string

	for {
		fmt.Print("Please enter a 5-letter word: ")
		fmt.Scan(&guess)

		// validate the length of the guess
		if len(guess) == 5 {
			break
		} else {
			fmt.Print("Invalid input. You must enter a 5-letter word. ")
		}
	}

	return strings.ToUpper(guess)
}

// compare_answer compares the user's guess to the answer and returns
// an integer slice where a 0 indicates the letter in
// that position does not occur in the answer,
// a 1 indicates a match between the letter and position,
// and a 2 indicates that letter exists in the answer, but in
// a different position.
func compare_answer(guess_str string, answer string) []int {
	comparison := make([]int, 5, 5)
	indices := make([]int, 0, 5)
	guess := []rune(guess_str)

	// one-by-one comparison to get the "green" letters
	for i, letter := range answer {
		if letter == guess[i] {
			comparison[i] = 1
		} else {
			// while iterating, create index of non-matches
			indices = append(indices, i)
		}
	}

	// iterate over non-exact matches, highlight those that
	// occur in string at another position
	var letter rune
	for _, j := range indices {
		letter = guess[j]
		if strings.ContainsRune(answer, letter) {
			comparison[j] = 2
		}
	}

	return comparison
}

// color_comparison creates a string using ANSI escape codes that
// simulates the graphical output from a regular game of Wordle.
// The ANSI codes were taking from a this [gist](https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797)
//
func color_comparison(guess string, comparison []int) string {
	// declare the ANSI escape codes
	var green, yellow, reset string = "\x1b[30;42;1m", "\x1b[30;43;1m", "\x1b[0m"

	// if running windows, don't use ANSI codes
	if strings.Contains(runtime.GOOS, "windows") {
		green, yellow, reset = "", "", ""
	}

	var letter, answer string // the letter as it will be displayed

	// iterate over the guess, building the colored output
	for i, j := range comparison {
		letter = string(guess[i])
		switch j {
		case 0:
			answer += letter
		case 1:
			answer += green + letter + reset
		case 2:
			answer += yellow + letter + reset
		}
	}

	return string(answer)
}
