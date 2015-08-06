/*
Copyright (c) 2015 Eric Knapik, All Rights Reserved

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

  1. Redistributions of source code must retain the above copyright
     notice, this list of conditions and the following disclaimer.

  2. Redistributions in binary form must reproduce the above copyright
     notice, this list of conditions and the following disclaimer in the
     documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
// This file is about the creation of tables that should
// be able to become import files. And using made tables to tag the parts
// of speech of possible copyright laced text. Once the text is tagged
// it can be called on by the file copyright.go to extract a possible
// notice, using a DFA.

package tagger

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
)

// global const:
const numOfTags int = 26

// global regex
var copyright = regexp.MustCompile("(\\\\[(]co)")

// A struct/pair for the dictionary value
// The dictionary actually stores an array of these.
type TagFrequency struct {
	tag  string
	freq float32
}

// The Tagger Object
type Tagger struct {
	Dictionary  map[string][]TagFrequency
	TransMatrix [][]float32
	// for the copyright extraction
	CopyrightDFA  map[Tri]int
	CopyrightSyms string
}

type TaggedWord struct {
	word      string
	tag       string
	byteStart int
}

// three variable structure used in DFA translation
type Tri struct {
	state int
	word  string
	pos   string
}

// create maps for converting tag string to integer and vice versa
var TagStrToInt = make(map[string]int)
var TagIntToStr = make(map[int]string)

// intitalizes the map used to convert the integer and string
// representation of part of speech tags
func initTagConversionMap() {

	TagStrToInt["bos"] = 0
	TagStrToInt["$"] = 1
	TagStrToInt["\""] = 2
	TagStrToInt["("] = 3
	TagStrToInt[")"] = 4
	TagStrToInt[","] = 5
	TagStrToInt["--"] = 6
	TagStrToInt["."] = 7
	TagStrToInt[":"] = 8
	TagStrToInt["cc"] = 9
	TagStrToInt["cd"] = 10
	TagStrToInt["dt"] = 11
	TagStrToInt["fw"] = 12
	TagStrToInt["jj"] = 13
	TagStrToInt["ls"] = 14
	TagStrToInt["nn"] = 15
	TagStrToInt["np"] = 16
	TagStrToInt["pos"] = 17
	TagStrToInt["pr"] = 18
	TagStrToInt["rb"] = 19
	TagStrToInt["sym"] = 20
	TagStrToInt["to"] = 21
	TagStrToInt["uh"] = 22
	TagStrToInt["vb"] = 23
	TagStrToInt["md"] = 24
	TagStrToInt["in"] = 25

	TagIntToStr[0] = "bos"
	TagIntToStr[1] = "$"
	TagIntToStr[2] = "\""
	TagIntToStr[3] = "("
	TagIntToStr[4] = ")"
	TagIntToStr[5] = ","
	TagIntToStr[6] = "--"
	TagIntToStr[7] = "."
	TagIntToStr[8] = ":"
	TagIntToStr[9] = "cc"
	TagIntToStr[10] = "cd"
	TagIntToStr[11] = "dt"
	TagIntToStr[12] = "fw"
	TagIntToStr[13] = "jj"
	TagIntToStr[14] = "ls"
	TagIntToStr[15] = "nn"
	TagIntToStr[16] = "np"
	TagIntToStr[17] = "pos"
	TagIntToStr[18] = "pr"
	TagIntToStr[19] = "rb"
	TagIntToStr[20] = "sym"
	TagIntToStr[21] = "to"
	TagIntToStr[22] = "uh"
	TagIntToStr[23] = "vb"
	TagIntToStr[24] = "md"
	TagIntToStr[25] = "in"
}

// Initialization for the Tagger object
// Takes a file path and will create the unigram dictionary and transition
// matrix required for sentence tagging and NLP processing
func New(path string) *Tagger {

	// initialize my TagStrToInt and TagIntToStr
	initTagConversionMap()

	// initialize the dictionary
	var dictionary = make(map[string][]TagFrequency)

	// Initialize the transition Matrix,
	var transMatrix = make([][]float32, numOfTags)
	for row := range transMatrix {
		transMatrix[row] = make([]float32, numOfTags)
	}

	// read through the corpus file to populate the dictionary and transMatrix
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic("could not read the file for tagging")
	}
	rawString := string(raw[:])

	prevTag := "."
	currTag := ""
	// I need to use the split feature on the coprus. So the input Corpus must have three
	// spaces between each word|~|tag pair. Once I have each word|~|tag pair I can
	// then split on the delimeter. Assumptions are made but I am assuming a safe input
	// file which I feel is acceptable, since if you save a file you know what it will
	// look like.
	textArry := strings.Split(rawString, "   ")
	textArry = textArry[:len(textArry)-1]
	for _, word := range textArry {
		wrdArry := strings.Split(word, "|~|")
		currTag = wrdArry[1]
		incrementUnigramWrd(dictionary, wrdArry[0], currTag)
		incrementTransMatrix(&transMatrix, TagStrToInt[prevTag], TagStrToInt[currTag])
		prevTag = currTag
	}
	// everything is counted now convert the dictionary and TransMatrix to probabilistic
	convertDictToProb(dictionary)
	convertTransMatrixToProb(&transMatrix)

	// SETUP THE COPYRIGHT DFA
	symbols, dfa := mkNoticeDFA()

	return &Tagger{Dictionary: dictionary, TransMatrix: transMatrix, CopyrightDFA: dfa, CopyrightSyms: symbols}
}

// This is the counter of tag transitions. Moving from one part of speech tag
// to the other. When reading the input corpus this function is called to
// increment/make note of every part of speech tag transition.
// transitionOccurances = transMatrix[prev POS tag][current POS tag]
func incrementTransMatrix(transMatrix *[][]float32, prevTagIndex int, currTagIndex int) {
	(*transMatrix)[prevTagIndex][currTagIndex]++
}

// Given the unigram word dictionary, a word and the given part of speech
// tag for the word this will increment if the word already existed in the dictionary
// if the word did not this will create a new entry and set the times seen to 1
func incrementUnigramWrd(dictionary map[string][]TagFrequency, word string, tag string) {
	// dictionary is the map used for unigram word count/frequency
	// it is a key->slice of TagFrequency objects
	if dictionary[word] != nil {
		for i := 0; i < len(dictionary[word]); i++ {
			if tag == dictionary[word][i].tag {
				dictionary[word][i].freq++
				return
			}
		}
		dictionary[word] = append(dictionary[word], TagFrequency{tag, 1})
		return
	} else {
		dictionary[word] = append(dictionary[word], TagFrequency{tag, 1})
		return
	}
}

// This will convert the dictionary which was in the form of
// counted occurances into a dictionary of probability for each part of speech
// tag given a specific word
func convertDictToProb(dictionary map[string][]TagFrequency) {
	// dictionary is a global variable
	var total float32
	for key := range dictionary {
		total = 0
		for i := 0; i < len(dictionary[key]); i++ {
			total = total + dictionary[key][i].freq
		}
		for i := 0; i < len(dictionary[key]); i++ {
			dictionary[key][i].freq = dictionary[key][i].freq / total
		}
	}
}

// This will convert the Transition Matrix to the probability
// Transition matrix the likelyhood of a given part of speech tag transition.
// Moving from tag A to tag B will result in what probility.
// transMatrix[FromTagA][ToTagB] = Probability X
// This is where the smoothing will be implemented
// I am using Laplace Smoothing across the transitional probability
// This means that every transition has a small probability of happeing
func convertTransMatrixToProb(transMatrix *[][]float32) {
	// transMatrix is a global variable
	var total float32

	for row := 0; row < numOfTags; row++ {
		total = float32(numOfTags)
		for col := 0; col < numOfTags; col++ {
			total += (*transMatrix)[row][col]
		}

		for col := 0; col < numOfTags; col++ {
			(*transMatrix)[row][col] = ((*transMatrix)[row][col] + 1) / total
		}
	}
}

// Given a word with an unknown part of speech. Using a model based from the
// Brill tagger, Krymolowski and Roth 1998 research (http://www.aclweb.org/anthology/P98-2186)
// This returns a guessed part of speech for unknown words
func tagUnkown(word string) string {

	// perform an N for loop checking for integer ascii value
	var i int
	for i = 0; i < len(word); i++ {
		if word[i] > 47 && word[i] < 58 {
			return "cd"
		}
	}

	loWord := strings.ToLower(word)

	switch {
	case strings.HasSuffix(loWord, "able"):
		return "jj"
	case strings.HasSuffix(loWord, "ible"):
		return "jj"
	case strings.HasSuffix(loWord, "ic"):
		return "jj"
	case strings.HasSuffix(loWord, "ous"):
		return "jj"
	case strings.HasSuffix(loWord, "al"):
		return "jj"
	case strings.HasSuffix(loWord, "ful"):
		return "jj"
	case strings.HasSuffix(loWord, "less"):
		return "jj"
	case strings.HasSuffix(loWord, "ly"):
		return "rb"
	case strings.HasSuffix(loWord, "ate"):
		return "vb"
	case strings.HasSuffix(loWord, "fy"):
		return "vb"
	case strings.HasSuffix(loWord, "ize"):
		return "vb"
	}

	// perform an N for loop checking for capital letter
	for i = 0; i < len(word); i++ {
		if word[i] > 64 && word[i] < 91 {
			return "np"
		}
	}

	switch {
	case strings.HasSuffix(loWord, "ion"):
		return "nn"
	case strings.HasSuffix(loWord, "ess"):
		return "nn"
	case strings.HasSuffix(loWord, "ment"):
		return "nn"
	case strings.HasSuffix(loWord, "er"):
		return "nn"
	case strings.HasSuffix(loWord, "or"):
		return "nn"
	case strings.HasSuffix(loWord, "ist"):
		return "nn"
	case strings.HasSuffix(loWord, "ism"):
		return "nn"
	case strings.HasSuffix(loWord, "ship"):
		return "nn"
	case strings.HasSuffix(loWord, "hood"):
		return "nn"
	case strings.HasSuffix(loWord, "ology"):
		return "nn"
	case strings.HasSuffix(loWord, "ty"):
		return "nn"
	case strings.HasSuffix(loWord, "y"):
		return "nn"
	default:
		return "fw"
	}
}

// Performs several string substitutions so that the tagger has an easier job
// These calls are to substitute parts of the string for other parts
// Once the sentence is formatted correctly it returns the string
func formatSent(rawBytes []byte) []byte {
	// to ensure a propper formatting.
	// replace weird copyright symbols
	// replaces \(co with (c)
	rawBytes = copyright.ReplaceAll(rawBytes, []byte("(c) ")) // added extra space to preserve byte offset

	// replace contractions
	// for byte preservation can not do these, but for more accurate tagging
	// replacing contractions can be useful
	/*
		rawBytes = bytes.Replace(rawBytes, []byte("ain't"), []byte("are not"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("won't"), []byte("will not"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("can't"), []byte("cannot"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("n't"), []byte(" not"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'re"), []byte(" are"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'m"), []byte(" am"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'ll"), []byte(" will"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'ve"), []byte(" have"), -1)
	*/
	return rawBytes
}

// returns true if the given byte is a white space character
func isSpace(b ...byte) bool {

	return bytes.Contains([]byte(" \n\r\t"), b)
}

// returns true if the given byte is a ASCII symbolic character
func isSymbol(b ...byte) bool {
	return bytes.Contains([]byte("~!`@#$%^&*()[]_+-=|}{:;'\"/\\.?><,"), b)
}

// Given a slice of raw bytes will convert this into a slice of
// TaggedWord objects with no tag set. This slice of TaggedWord objects will
// then be given to the tagger for determining the part of speech tag
func mkWrdArray(rawBytes []byte) []TaggedWord {

	currByte := 0
	wordStart := currByte
	var taggedWords []TaggedWord = make([]TaggedWord, 0)

	for currByte < len(rawBytes) {
		if isSpace(rawBytes[currByte]) {
			if wordStart != currByte { // add the word if I can
				taggedWords = append(taggedWords, TaggedWord{word: string(rawBytes[wordStart:currByte]), tag: "", byteStart: wordStart})
			}
			currByte++
			wordStart = currByte
		} else if isSymbol(rawBytes[currByte]) {
			if wordStart != currByte { // add the word if I can
				taggedWords = append(taggedWords, TaggedWord{word: string(rawBytes[wordStart:currByte]), tag: "", byteStart: wordStart})
			}
			wordStart = currByte
			currByte++
			taggedWords = append(taggedWords, TaggedWord{word: string(rawBytes[wordStart:currByte]), tag: "", byteStart: wordStart})
			wordStart = currByte
		} else {
			currByte++
		}
	}
	taggedWords = append(taggedWords, TaggedWord{word: string(rawBytes[wordStart:currByte]), tag: "", byteStart: wordStart})
	return taggedWords
}

// Given any string this will return a slice of TaggedWord objects
// representing that word in the sentence and the part of speech for
// that word
func (copyrightTagger *Tagger) TagBytes(rawBytes []byte) []TaggedWord {
	// ERROR AND SANITIZATION CHECKS
	var wrdArry []TaggedWord = make([]TaggedWord, 0)
	if len(rawBytes) < 1 { // do I even need to do any work
		return wrdArry
	}

	// perform several regular expression subs and other so that the string is in a desired
	// form
	rawBytes = formatSent(rawBytes)
	// split the sentence propperly
	wrdArry = mkWrdArray(rawBytes)

	sentLength := len(wrdArry) + 1             // I need 1 more for the start of the sentence
	sentMatrix := make([][]float32, numOfTags) // Create the sentence Matrix
	for row := range sentMatrix {
		sentMatrix[row] = make([]float32, sentLength)
	}

	// initialize the first column
	var lastBestProb float32 = 1.0
	var lastBestTag int = TagStrToInt["."]
	var currBestProb float32 = 0.0
	var currBestTag int = 0

	sentMatrix[TagStrToInt["."]][0] = 1.0 // the max probability something can be
	for wrdIndex := 0; wrdIndex < len(wrdArry); wrdIndex++ {
		for tagIndex := 0; tagIndex < numOfTags; tagIndex++ {
			var currTrans float32 = copyrightTagger.TransMatrix[lastBestTag][tagIndex]
			var currProb float32 = lastBestProb * currTrans

			if len(copyrightTagger.Dictionary[wrdArry[wrdIndex].word]) != 0 { // has the word been seen before?
				if wrdArry[wrdIndex].word == "." || wrdArry[wrdIndex].word == "?" || wrdArry[wrdIndex].word == "!" {
					sentMatrix[TagStrToInt["."]][wrdIndex+1] = 1.0
				} else {
					for _, tagObject := range copyrightTagger.Dictionary[wrdArry[wrdIndex].word] {
						if TagIntToStr[tagIndex] == tagObject.tag {
							sentMatrix[tagIndex][wrdIndex+1] = currProb * tagObject.freq
						}
					}
				}
				// check for the word not carring about capitalization
			} else if len(copyrightTagger.Dictionary[strings.ToLower(wrdArry[wrdIndex].word)]) != 0 {
				for _, tagObject := range copyrightTagger.Dictionary[strings.ToLower(wrdArry[wrdIndex].word)] {
					if TagIntToStr[tagIndex] == tagObject.tag {
						sentMatrix[tagIndex][wrdIndex+1] = currProb * tagObject.freq
					}
				}
			} else { // Try to determine tag based on transitional probability and word itself
				if currTrans >= 0.7 {
					sentMatrix[tagIndex][wrdIndex+1] = currProb
				} else {
					likelyTag := tagUnkown(wrdArry[wrdIndex].word)
					sentMatrix[TagStrToInt[likelyTag]][wrdIndex+1] = currProb * 0.95
				}
			}
			// see if this is the best transition for next column
			if currProb > currBestProb {
				currBestProb = currProb
				currBestTag = tagIndex
			}
		}
		lastBestProb = currBestProb
		lastBestTag = currBestTag
	}
	// Sentence Matrix Created.
	// Now walk through the matrix assigning the best tag to each word
	for wrdIndex := 0; wrdIndex < len(wrdArry); wrdIndex++ {
		var tagProb float32 = 0.0
		for tagIndex := 0; tagIndex < numOfTags; tagIndex++ {
			if sentMatrix[tagIndex][wrdIndex+1] > tagProb {
				tagProb = sentMatrix[tagIndex][wrdIndex+1]
				wrdArry[wrdIndex].tag = TagIntToStr[tagIndex]
			}
		}
	}

	// compress numbers and propper nouns that might have been split
	wrdArry = compressNumInString(wrdArry)
	wrdArry = compressNP(wrdArry)

	return wrdArry
}

// The DFA required for number compression
// into one number not number period number
// START is start state
// INTERM is intermediate state
// REJECT is reject state
// ACCEPT is accept state
func mkNumCompressDFA() map[Tri]int {
	// symbols := "cd."
	dfa := make(map[Tri]int)

	dfa[Tri{state: START, word: "X", pos: "cd"}] = START
	dfa[Tri{state: INTERM, word: "X", pos: "cd"}] = ACCEPT
	dfa[Tri{state: REJECT, word: "X", pos: "cd"}] = START
	dfa[Tri{state: ACCEPT, word: "X", pos: "cd"}] = START

	dfa[Tri{state: START, word: ".", pos: "."}] = INTERM
	dfa[Tri{state: INTERM, word: ".", pos: "."}] = REJECT
	dfa[Tri{state: REJECT, word: ".", pos: "."}] = REJECT
	dfa[Tri{state: ACCEPT, word: ".", pos: "."}] = REJECT

	dfa[Tri{state: START, word: "?", pos: "."}] = REJECT
	dfa[Tri{state: INTERM, word: "?", pos: "."}] = REJECT
	dfa[Tri{state: REJECT, word: "?", pos: "."}] = REJECT
	dfa[Tri{state: ACCEPT, word: "?", pos: "."}] = REJECT

	dfa[Tri{state: START, word: "!", pos: "."}] = REJECT
	dfa[Tri{state: INTERM, word: "!", pos: "."}] = REJECT
	dfa[Tri{state: REJECT, word: "!", pos: "."}] = REJECT
	dfa[Tri{state: ACCEPT, word: "!", pos: "."}] = REJECT

	dfa[Tri{state: START, word: "X", pos: "X"}] = REJECT
	dfa[Tri{state: INTERM, word: "X", pos: "X"}] = REJECT
	dfa[Tri{state: REJECT, word: "X", pos: "X"}] = REJECT
	dfa[Tri{state: ACCEPT, word: "X", pos: "X"}] = REJECT

	return dfa
}

// Given a slice of TaggedWord objects this will
// compress floating point and numbers containing periods
// that might have been split up by the tagger's formatting
func compressNumInString(inSent []TaggedWord) []TaggedWord {
	dfa := mkNumCompressDFA()

	var finalSent []TaggedWord = make([]TaggedWord, 0)

	currentState := REJECT // the dead state
	var compNum []string = make([]string, 0)
	var saveNum []TaggedWord = make([]TaggedWord, 0)
	var saveStartByte int

	for _, taggedWord := range inSent {
		// Make the transition to the next state based on the input
		if taggedWord.tag == "." {
			currentState = dfa[Tri{state: currentState, word: taggedWord.word, pos: taggedWord.tag}]
		} else if taggedWord.tag == "cd" {
			currentState = dfa[Tri{state: currentState, word: "X", pos: taggedWord.tag}]
		} else {
			currentState = dfa[Tri{state: currentState, word: "X", pos: "X"}]
		}

		// Based on the current input decide how to save information
		if currentState == START {
			finalSent = append(finalSent, saveNum...)
			compNum = nil
			saveNum = nil
			compNum = append(compNum, taggedWord.word)
			saveStartByte = taggedWord.byteStart
			saveNum = append(saveNum, taggedWord)
		} else if currentState == INTERM {
			compNum = append(compNum, ".")
			saveNum = append(saveNum, taggedWord)
		} else if currentState == REJECT {
			finalSent = append(finalSent, saveNum...)
			saveNum = nil
			finalSent = append(finalSent, taggedWord)
		} else if currentState == ACCEPT {
			compNum = append(compNum, taggedWord.word)
			saveNum = nil
			saveNum = append(saveNum, TaggedWord{word: strings.Join(compNum, ""), tag: "cd", byteStart: saveStartByte})
			currentState = START
		}
	}
	if currentState != REJECT {
		finalSent = append(finalSent, saveNum...)
	}

	// returns the joined string with 'floating point' numbers combined
	return finalSent
}

// Similar to the compressNumInString this recompresses
// propper nouns that the tagger possibly separated to generalize
// tagging and account for words it has not seen before.
func compressNP(inSent []TaggedWord) []TaggedWord {
	var finalSent []TaggedWord = make([]TaggedWord, 0)

	prevTag := ""
	var saveWord []string = make([]string, 0)
	var saveByteStart int
	for _, taggedWord := range inSent {

		if prevTag == "np" && taggedWord.word == "." {
			saveWord = append(saveWord, ".")
			finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart})
			saveWord = nil
		} else if prevTag == "np" && taggedWord.tag == "np" {
			finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart})
			saveWord = nil
			saveWord = append(saveWord, taggedWord.word)
		} else if prevTag == "np" && taggedWord.word != "." {
			finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart}, taggedWord)
			saveWord = nil
		} else if taggedWord.tag == "np" {
			saveWord = append(saveWord, taggedWord.word)
		} else {
			finalSent = append(finalSent, taggedWord)
		}

		saveByteStart = taggedWord.byteStart
		prevTag = taggedWord.tag

	}
	if prevTag == "np" {
		finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart})
	}

	return finalSent
}

func toString(inSent []TaggedWord) string {
	var finalSent = make([]string, 0)
	for _, taggedWord := range inSent {
		finalSent = append(finalSent, taggedWord.word)
	}
	return strings.Join(finalSent, " ")
}
