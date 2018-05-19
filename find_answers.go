package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"sort"
	"io"
)

const (
	numQuestions = 5
)

type MatchedSentenceInfo struct {
	Index int
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

// Remove tokens that can be ignored, duplicates
func cleanUpTokens (elements []string) []string {
	encountered := map[string]bool{}

	for v:= range elements {
		if (ignoreMap[elements[v]]) {
			continue
		}

		// TODO Perform Stemming
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}

func main () {
	if (len (os.Args) != 2) {
		fmt.Println("Usage: ./find_answers <full path for filename>")
		return
	}

	inputFile, err := os.Open(os.Args[1]);
	if err != nil {
		fmt.Printf("Unable to open the file, error - %v", err)
		return;
	}

	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	paragraph, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input - %v \n", err)
		return;
	}

	var questions [numQuestions]string;
	for i := 0; i < numQuestions; i++ {
		questions[i], err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading questions  - %v \n", err)
			return
		}
	}

	var answerStr string
	answerStr, err = reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading answers - %v", err)
		return
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
		//fmt.Printf("SENTENCE[%d] - %s\n", i, sentences[i])
		senTokens[i] = strings.FieldsFunc(sentences[i], func (delim rune) bool {
			return delim == ',' || delim == ';' || delim == ' ';
		})

		senTokens[i] = cleanUpTokens(senTokens[i])
	}

	// Tokenize questions
	questionTokens := make([][] string, numQuestions)
	for i := 0; i < numQuestions; i++ {
		questionTokens[i] = strings.FieldsFunc(questions[i], func (delim rune) bool {
			return delim == ',' || delim == ';' || delim == ' ' || delim == '?';
		})
		questionTokens[i] = cleanUpTokens(questionTokens[i])
	}

	// Find answers to the questions
LoopQuestions:
	for i := 0; i < numQuestions; i++ {
		// Filter sentences
		matchedSentences := filterSentences(senTokens, numSentences, questionTokens[i])

		// Find answers in top matched sentences
		for j := 0; j < len(matchedSentences); j++ {
			if matchedSentences[j].NumMatches != 0 {
				for k := 0; k < len(answers); k++ {
					senIndex := matchedSentences[j].Index
					//fmt.Printf("SENTENCE[%d] - %s\n", senIndex, sentences[senIndex])
					//fmt.Printf("ANSWER - %s\n", answers[k])
					if (strings.Contains(sentences[senIndex], answers[k])) {
						//fmt.Println(questions[i])
						//fmt.Printf("SENTENCE[%d] - %s\n", senIndex, sentences[senIndex])
						//fmt.Printf("ANSWER - %s\n", answers[k])

						// TODO we might need another iteration if multiple answers match the same sentence
						// Maybe some kind of prediction algorithm ?
						fmt.Println(answers[k])
						answers = append(answers[:k], answers[k+1:]...)
						continue LoopQuestions
					}
				}
			}
		}
	}
}