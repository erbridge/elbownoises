package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/erbridge/gotwit"
	"github.com/erbridge/gotwit/twitter"
)

type (
	corpus struct {
		Words    []string            `json:"words"`
		Prefixes map[string][]string `json:"prefixes"`
	}
)

func getCorpus() (c corpus, err error) {
	corpusFile, err := os.Open("data/corpus.json")

	if err != nil {
		return
	}

	parser := json.NewDecoder(corpusFile)

	if err = parser.Decode(&c); err != nil {
		return
	}

	return
}

func createTweetText(c corpus) (text string) {
	count := rand.Intn(3) + rand.Intn(3) + 1

	words := make([]string, count)

	for i := 0; i < count; i++ {
		word := c.Words[rand.Intn(len(c.Words))]

		index := rand.Intn(len(c.Prefixes))

		for k, v := range c.Prefixes {
			if index == 0 {
				if strings.HasPrefix(word, k) {
					prefix := v[rand.Intn(len(v))]

					prefixCount := rand.Intn(4) - rand.Intn(3)

					if prefixCount > 0 {
						prefixes := make([]string, prefixCount+1)

						for j := 0; j < prefixCount; j++ {
							prefixes[j] = prefix
						}

						word = strings.Join(prefixes, "-") + word
					}
				}

				break
			}

			index--
		}

		words[i] = word
	}

	text = strings.Join(words, " ")

	return
}

func main() {
	var (
		con twitter.ConsumerConfig
		acc twitter.AccessConfig
	)

	f := "secrets.json"
	if _, err := os.Stat(f); err == nil {
		con, acc, _ = twitter.LoadConfigFile(f)
	} else {
		con, acc, _ = twitter.LoadConfigEnv()
	}

	b := gotwit.NewBot("elbownoises", con, acc)

	c, err := getCorpus()

	if err != nil {
		panic(err)
	}

	go func() {
		if err = b.Start(); err != nil {
			panic(err)
		}
	}()

	rand.Seed(time.Now().UnixNano())

	tweet := createTweetText(c)

	b.Post(tweet, false)

	if err = b.Stop(); err != nil {
		panic(err)
	}
}
