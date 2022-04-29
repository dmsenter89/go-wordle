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
)

func main() {
	fmt.Println("hello, world!")

	word_list := load_dictionary()

	// pick a random word from list
	randidx := rand.Intn(len(word_list))
	fmt.Println("Random word: ", word_list[randidx])
}

// load_dictionary loads a world dictionary.
// If the wordle dictionary is not found, one will be created.
// Returns a string slice consisting of 5-letter words.
func load_dictionary() []string {
	// read wordle dict if present,
	body, err := ioutil.ReadFile("wordle.dict") // just pass the file name
	if err != nil {
		fmt.Println("Dictionary not found (", err, ").\nDownloading dictionary.")
		body, err = download_dictionary()
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body))) // f is the *os.File
	word_list := make([]string, 0, 15000)                        // pre-allocate room for 15k words
	var curr_word string

	for scanner.Scan() {
		curr_word = scanner.Text()
		if len(curr_word) == 5 {
			word_list = append(word_list, strings.ToUpper(curr_word))
		}
	}

	return word_list
}

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
