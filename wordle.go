package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to go-wordle!")

	word_list := load_dictionary()

	// pick a random word from list
	rand.Seed(time.Now().UnixNano())
	var word string
	word = choose_random_word(word_list)

	fmt.Println("Random word: ", word)
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
