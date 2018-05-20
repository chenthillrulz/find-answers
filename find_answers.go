package main

import (
	"strings"
	"sort"
	"os"
	"fmt"
	"bufio"
	"io"
)

const (
	numQuestions = 5
)

/**
* The following can go into the tokenizer package
 */
// TODO This list needs to be exhaustive to ignore iirelevant words
// revising parts of speech ;)
var ignoreMap = map[string] bool {
	// Ignore words pertaining to questions
	"when" : true,
	"which" : true,
	"what" : true,
	"who" : true,

	// Ignore non key words
	"are" : true,
	"is" : true,
	"the" : true,
	"of" : true,
}

// https://github.com/c9s/inflect - converts singular <-> plural library
// nltk python project has some good utilities
func crudeStemming (token string) string {
	if (strings.HasSuffix(token, "s")) {
		token = strings.TrimSuffix(token, "s")
	} else if (strings.HasSuffix(token, "ed")) {
		token = strings.TrimSuffix(token, "ed")
	}

	return token;
}


// Remove tokens that can be ignored, duplicates. Perform crude stemming
func sanitizeTokens (tokens []string) []string {
	encountered := map[string]bool{}

	for v:= range tokens {
		tokens[v] = strings.ToLower(tokens[v])
		if (ignoreMap[tokens[v]]) {
			continue
		}

		tokens[v] = crudeStemming(tokens[v])
		encountered[tokens[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}

func GenerateTokens (sentence string) []string {
	tokens := strings.FieldsFunc(sentence, func (delim rune) bool {
		return delim == ',' || delim == ';' || delim == ' ' || delim == '?';
	})
	tokens = sanitizeTokens(tokens)

	return tokens
}

/* */


type MatchedSentenceInfo struct {
	Index      int
	NumMatches int
}

func filterSentences (sentenceTokens [][]string, numSentences int,  questionTokens []string) []MatchedSentenceInfo {
	matches := make([]MatchedSentenceInfo, numSentences)

	for i:=0; i < numSentences; i++ {
		numMatches := 0
		var match MatchedSentenceInfo

		for j:=0; j < len(questionTokens); j++ {
			for k := 0; k < len(sentenceTokens[i]); k++ {
				if (questionTokens[j] == sentenceTokens[i][k]) {
					numMatches++
				}
			}
		}

		match.Index = i
		match.NumMatches = numMatches
		matches[i] = match
	}

	// Sort results in descending order
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].NumMatches > matches[j].NumMatches
	})

	return matches
}

/**
The algorithm maps the key words from questions to the sentences and  matches the answers with the filtered sentences.
 */

func FindAnswers (inputFile *os.File) ([]string, error) {
	var matchedAnswers []string;

	reader := bufio.NewReader(inputFile)
	paragraph, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input - %v \n", err)
		return nil, err;
	}

	var questions [numQuestions]string;
	for i := 0; i < numQuestions; i++ {
		questions[i], err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading questions  - %v \n", err)
			return nil, err;
		}
	}

	var answerStr string
	answerStr, err = reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading answers - %v", err)
		return nil, err;
	}
	answers := strings.Split(answerStr, ";")
	numAnswers := len(answers)
	answers[numAnswers -1] = strings.TrimRight(answers[numAnswers-1], "\n")

	sentences := strings.Split(paragraph, ".")
	numSentences := len(sentences)

	// TODO perform stemming of tokens to filter variations in verbs eg: 'aim, aims, aimed' etc.
	// Tokenize sentences
	senTokens := make([][] string, numSentences)
	for i := 0; i < numSentences; i++ {
		senTokens[i] = GenerateTokens(sentences[i])
	}

	// Tokenize questions
	questionTokens := make([][] string, numQuestions)
	for i := 0; i < numQuestions; i++ {
		questionTokens[i] = GenerateTokens(questions[i])
	}

	// Find answers to the questions
LoopQuestions:
	for i := 0; i < numQuestions; i++ {
		// Filter sentences
		matchedSentences := filterSentences(senTokens, numSentences, questionTokens[i])

		// Find answers in top matched sentences
		for j := 0; j < len(sentences); j++ {
			if matchedSentences[j].NumMatches != 0 {

				// Iterate answers and find the top matched ones
				for k := 0; k < len(answers); k++ {

					senIndex := matchedSentences[j].Index
					if (strings.Contains(sentences[senIndex], answers[k])) {
						// TODO we might need another iteration if multiple answers match the same sentence
						// Maybe some kind of prediction algorithm ?
						matchedAnswers = append(matchedAnswers, answers[k])
						answers = append(answers[:k], answers[k+1:]...)
						continue LoopQuestions
					}
				}
			}
		}
	}

	return matchedAnswers, nil;
}


func main () {
	matchedAnswers, err := FindAnswers(os.Stdin)
	if err != nil {
		return
	}

	for i := 0; i < len(matchedAnswers); i++ {
		fmt.Println(matchedAnswers[i])
	}
}