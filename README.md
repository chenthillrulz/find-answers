# find-answers
Find the answers using a paragraph content and questions.

The algorithm uses key words from the questions and picks the sentences that could be the probable answers.
Then it matches the answers with the filtered sentences to find the best matched answer.
More in the lines of how we approach comprehension questions from GRE or CAT.

# Usage
$ go build find_answers.go
 
$ ./find_answers test-data/input1.txt

To run units,
$ go test


More files could be added to test-data directory. To include them in tests, modify the inputs and outputs array.