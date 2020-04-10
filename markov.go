package markov

import (
	"fmt"
	"math/rand"
	"strings"
)

const maxPrefixLength = 3
const suffixLength = 1

// Chain is a Markov Chain Text Generator.
type Chain struct {
	chain map[string][]string

	prefixLength int
	suffixLength int

	allowedRunes map[rune]struct{}
}

// New returns a new Chain.
func New(prefixLength int) *Chain {
	if prefixLength < 1 {
		prefixLength = 1
	} else if prefixLength > maxPrefixLength {
		prefixLength = maxPrefixLength
	}

	chain := &Chain{
		chain:        make(map[string][]string),
		prefixLength: prefixLength,
		suffixLength: suffixLength,
	}
	chain.generateAllowedRunes()

	return chain
}

// Add adds the given text to the Chain.
func (c *Chain) Add(text string) {
	words := strings.Fields(text)
	for i := range words {
		// generate prefix
		prefix := []string{}
		for n := 0; n < c.prefixLength; n++ {
			delta := i + n
			if delta < len(words) {
				cleanedWord := c.cleanWord(words[delta])
				if cleanedWord == "" {
					break
				}

				prefix = append(prefix, cleanedWord)
			}
		}

		// malformed prefix due to 'unclean' word(s)
		// skip suffix to maintain chain health
		// TODO check if any non-zero length prefix is fine
		if len(prefix) != c.prefixLength {
			continue
		}

		// generate suffix
		suffix := []string{}
		for n := 0; n < c.suffixLength; n++ {
			delta := i + c.prefixLength + n
			if delta < len(words) {
				cleanedWord := c.cleanWord(words[delta])
				if cleanedWord == "" {
					break
				}

				suffix = append(suffix, cleanedWord)
			}
		}

		// add to chain
		prefixString := strings.Join(prefix, " ")
		suffixString := strings.Join(suffix, " ")

		if len(suffix) == 0 && len(suffixString) > 0 {
			fmt.Println("suffix: we have a problem here chief")
		}

		c.createPrefix(prefixString)

		if len(suffixString) > 0 {
			c.addSuffix(prefixString, suffixString)
		}
	}
}

func (c *Chain) createPrefix(prefix string) {
	if _, ok := c.chain[prefix]; !ok {
		c.chain[prefix] = []string{}
	}
}

func (c *Chain) addSuffix(prefix, suffix string) {
	c.chain[prefix] = append(c.chain[prefix], suffix)
}

// Generate generates text simulating the chain.
func (c *Chain) Generate(numWords int) string {
	for prefix, suffix := range c.chain {
		fmt.Println(prefix, suffix)
	}

	return ""
}

func (c *Chain) cleanWord(word string) string {
	return strings.ToLower(strings.Map(
		func(r rune) rune {
			if _, ok := c.allowedRunes[r]; ok {
				return r
			}

			return -1
		},
		word,
	))
}

func (c *Chain) generateAllowedRunes() {
	c.allowedRunes = map[rune]struct{}{
		'.': {}, '!': {}, '?': {},
		',': {}, '\'': {},
	}

	// lowercase letters
	for r := 'a'; r <= 'z'; r++ {
		c.allowedRunes[r] = struct{}{}
	}

	// uppercase letters
	for r := 'A'; r <= 'Z'; r++ {
		c.allowedRunes[r] = struct{}{}
	}
}

func randNum(min, max int) int {
	return rand.Intn(max-min) + min
}
