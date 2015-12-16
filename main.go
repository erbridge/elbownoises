package main

import (
	"encoding/json"
	"fmt"
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

	n := rand.Float32()

	for i := 0; i < count; i++ {
		letters := make([]string, 0)

		for _, r := range c.Words[rand.Intn(len(c.Words))] {
			limit := 1

			if r == 's' || r == 'k' || r == 'p' || r == 'n' || r == 'a' || r == 'i' || r == 'u' {
				n += rand.Float32() / 10

				if n > 1 {
					n = 0
				}

				if n > 0.9 {
					limit = 4
					n = rand.Float32() / 4
				} else if n > 0.75 {
					limit = 3
					n = rand.Float32() / 3
				} else if n > 0.6 {
					limit = 2
					n = rand.Float32() / 2
				}
			}

			for i := 0; i < limit; i++ {
				letters = append(letters, string(r))
			}
		}

		word := strings.Join(letters, "")

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

func postTweet(b gotwit.Bot, c corpus) {
	tweet := createTweetText(c)

	fmt.Println("Posting:", tweet)

	b.Post(tweet, false)
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

	now := time.Now()

	rand.Seed(now.UnixNano())

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour()+1,
		0,
		0,
		0,
		now.Location(),
	)

	sleep := next.Sub(now)

	fmt.Printf("%v until first tweet\n", sleep)

	time.Sleep(sleep)

	postTweet(b, c)

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		postTweet(b, c)
	}
}
