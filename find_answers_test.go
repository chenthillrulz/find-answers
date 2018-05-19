package main

import (
	"testing"
	"os"
	"bufio"
	"io"
	"strings"
)

func TestFindAnswers (t *testing.T) {
	inputs := []string{"./test-data/input1.txt"}
	outputs := []string{"./test-data/output1.txt"}

	LoopFiles:
	for i := 0; i < len (inputs); i++ {
		inputFile, err := os.Open(inputs[i])
		if (err != nil ) {
			t.Errorf("TEST[%d] failed, error opening input file %s - %v ", i, inputs[i], err)
			continue
		}

		answers, err := FindAnswers(inputFile)
		if (err != nil ) {
			t.Errorf("TEST[%d] failed - %v ", i, err)
			continue
		}

		outputFile, err := os.Open(outputs[i])
		if (err != nil ) {
			t.Errorf("TEST[%d] failed, error opening ouput file %s - %v ", i, outputs[i], err)
			continue
		}

		reader := bufio.NewReader(outputFile)
		for j := 0; j < len(answers); j++ {
			var expectAnswer, err = reader.ReadString('\n')
			if (err == io.EOF) {
				err = nil
			}
			expectAnswer = strings.TrimRight(expectAnswer, "\n")
			if (err != nil || expectAnswer != answers[j]) {
				t.Errorf("TEST[%d] failed, outputs does not match %s - error - %v ", i, outputs[i], err)
				continue LoopFiles
			}
		}
	}
}