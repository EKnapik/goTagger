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

// This file is about taking a slice of TaggedWord objects
// and determining if they contain a copyright notice
// there are a few tie overs of this into the tagger.go for
// optimizations but mostly they could be removed and the tagger.go
// can be a stand alone

package tagger

import (
	"strings"
)

// tri structure defined in main tagger.go
//type Tri struct {
//	state int
//	word  string
//	pos   string
//}

const (
	START   int = iota // start state
	LPAREN             // left parenthesis
	CCHAR              // Character C
	RPAREN             // right parenthesis
	NP                 // propper noun
	COMMA              // the comma symbol
	CD                 // any number
	DASH               // dash -
	IN                 // preposition
	DT                 // determiner
	ACCEPT             // accept state
	REJECT             // reject state
	SYM                // symbol tag
	OTHER              // anything else
	CSYM               // C symbol
	LPARENC            // left parenthesis after copyright
	INTERM             // Intermidiate state
)

// Given a string this will return whethe a copyright notice
// of that string if it exists, if not the empty string is returned
// The string must be tagged and propperly delimited
// else will just tag
// USING SHIFTING WINDOW STRATEGY
func (copyrightTagger *Tagger) Match(inBytes []byte) bool {
	var curByte = 0
	var lastCheckedByte = 0
	var currWords int // the amount of words in the notice
	var taggedSent []TaggedWord

	if len(inBytes) < 15 {
		return false
	}

	for lastCheckedByte < len(inBytes) {
		// this is shifting the window based on a period followed by space
		for curByte < len(inBytes) {
			if inBytes[curByte] == byte('.') {
				if curByte+1 < len(inBytes) && inBytes[curByte+1] == byte(' ') {
					curByte++
					break
				}
			}
			curByte++
		}

		currWords = 0
		// create array of tagged words
		taggedSent = copyrightTagger.TagBytes(inBytes[lastCheckedByte:curByte])

		currentState := REJECT
		var potentialNotice []TaggedWord = make([]TaggedWord, 0)
		var extractedNotice []TaggedWord = make([]TaggedWord, 0)
		for _, taggedWord := range taggedSent {

			// Is what I have good enough to add to the extracted Notices
			if currentState == ACCEPT {
				currWords = 0
				extractedNotice = append(extractedNotice, potentialNotice...)
				potentialNotice = nil
			}
			// Transition to the next state given current 'input'
			if strings.ToLower(taggedWord.word) == "copyright" || strings.ToLower(taggedWord.word) == "c" {
				currentState = copyrightTagger.CopyrightDFA[Tri{currentState, strings.ToLower(taggedWord.word), taggedWord.tag}]
			} else if strings.Contains(taggedWord.word, "©") {
				currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "©", "sym"}]
			} else if strings.Contains(copyrightTagger.CopyrightSyms, taggedWord.tag) {
				currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "X", taggedWord.tag}]
			} else {
				currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "X", "X"}]
			}
			// Because of multiple notices right after the other here's a check...
			if currentState == START || currentState == LPAREN || currentState == CSYM {
				currWords = 1

				if len(potentialNotice) > 3 { // Does it seem like something useful has been captured
					extractedNotice = append(extractedNotice, potentialNotice...)
					potentialNotice = nil
					potentialNotice = append(potentialNotice, taggedWord)
				} else {
					potentialNotice = nil
					potentialNotice = append(potentialNotice, taggedWord)
				}
			} else if currentState != REJECT && currentState != ACCEPT {
				currWords += 1
				potentialNotice = append(potentialNotice, taggedWord)
			} else if currentState == REJECT {
				currWords = 0
				potentialNotice = nil
			}

			// This is for optimization where all I care about yes or no
			if currWords > 3 {
				return true
			}
		}
		// Do a final check to see if I might have a notice as the very last part of the string
		// Be a little more vauge here to be safe
		if currentState == ACCEPT || currentState == CD || currentState == NP {
			extractedNotice = append(extractedNotice, potentialNotice...)
		}
		// Return the extracted notice or the not found notice as a string
		if len(extractedNotice) > 3 {
			return true
		}

		lastCheckedByte = curByte
	}
	return false // no copyright notice detected
}

// creates the DFA and symbol "array" needed to test the transitions
// for when a copyright notice can be found or noticed
func mkNoticeDFA() (string, map[Tri]int) {
	symbols := "(,),cd,np,dt,in,--,.,sym,cc"
	dfa := make(map[Tri]int)

	// possible copyright symbols
	// ©, â©, Å©
	dfa[Tri{START, "copyright", "nn"}] = START
	dfa[Tri{LPAREN, "copyright", "nn"}] = START
	dfa[Tri{CCHAR, "copyright", "nn"}] = START
	dfa[Tri{RPAREN, "copyright", "nn"}] = START
	dfa[Tri{NP, "copyright", "nn"}] = START
	dfa[Tri{COMMA, "copyright", "nn"}] = START
	dfa[Tri{CD, "copyright", "nn"}] = START
	dfa[Tri{DASH, "copyright", "nn"}] = START
	dfa[Tri{IN, "copyright", "nn"}] = START
	dfa[Tri{DT, "copyright", "nn"}] = START
	dfa[Tri{ACCEPT, "copyright", "nn"}] = START
	dfa[Tri{REJECT, "copyright", "nn"}] = START
	dfa[Tri{SYM, "copyright", "nn"}] = START
	dfa[Tri{OTHER, "copyright", "nn"}] = START
	dfa[Tri{CSYM, "copyright", "nn"}] = START
	dfa[Tri{LPARENC, "copyright", "nn"}] = START

	dfa[Tri{START, "c", "nn"}] = REJECT
	dfa[Tri{LPAREN, "c", "nn"}] = CCHAR
	dfa[Tri{CCHAR, "c", "nn"}] = REJECT
	dfa[Tri{RPAREN, "c", "nn"}] = REJECT
	dfa[Tri{NP, "c", "nn"}] = REJECT
	dfa[Tri{COMMA, "c", "nn"}] = REJECT
	dfa[Tri{CD, "c", "nn"}] = REJECT
	dfa[Tri{DASH, "c", "nn"}] = REJECT
	dfa[Tri{IN, "c", "nn"}] = REJECT
	dfa[Tri{DT, "c", "nn"}] = REJECT
	dfa[Tri{ACCEPT, "c", "nn"}] = REJECT
	dfa[Tri{REJECT, "c", "nn"}] = REJECT
	dfa[Tri{SYM, "c", "nn"}] = REJECT
	dfa[Tri{OTHER, "c", "nn"}] = REJECT
	dfa[Tri{CSYM, "c", "nn"}] = REJECT
	dfa[Tri{LPARENC, "c", "nn"}] = CCHAR

	dfa[Tri{START, "©", "sym"}] = CSYM
	dfa[Tri{LPAREN, "©", "sym"}] = CSYM
	dfa[Tri{CCHAR, "©", "sym"}] = CSYM
	dfa[Tri{RPAREN, "©", "sym"}] = CSYM
	dfa[Tri{NP, "©", "sym"}] = CSYM
	dfa[Tri{COMMA, "©", "sym"}] = CSYM
	dfa[Tri{CD, "©", "sym"}] = CSYM
	dfa[Tri{DASH, "©", "sym"}] = CSYM
	dfa[Tri{IN, "©", "sym"}] = CSYM
	dfa[Tri{DT, "©", "sym"}] = CSYM
	dfa[Tri{ACCEPT, "©", "sym"}] = ACCEPT
	dfa[Tri{REJECT, "©", "sym"}] = CSYM
	dfa[Tri{SYM, "©", "sym"}] = CSYM
	dfa[Tri{OTHER, "©", "sym"}] = CSYM
	dfa[Tri{CSYM, "©", "sym"}] = CSYM
	dfa[Tri{LPARENC, "©", "sym"}] = CSYM

	dfa[Tri{START, "X", "("}] = LPARENC
	dfa[Tri{LPAREN, "X", "("}] = REJECT
	dfa[Tri{CCHAR, "X", "("}] = REJECT
	dfa[Tri{RPAREN, "X", "("}] = REJECT
	dfa[Tri{NP, "X", "("}] = REJECT
	dfa[Tri{COMMA, "X", "("}] = REJECT
	dfa[Tri{CD, "X", "("}] = REJECT
	dfa[Tri{DASH, "X", "("}] = REJECT
	dfa[Tri{IN, "X", "("}] = REJECT
	dfa[Tri{DT, "X", "("}] = REJECT
	dfa[Tri{ACCEPT, "X", "("}] = REJECT
	dfa[Tri{REJECT, "X", "("}] = LPAREN
	dfa[Tri{SYM, "X", "("}] = REJECT
	dfa[Tri{OTHER, "X", "("}] = REJECT
	dfa[Tri{CSYM, "X", "("}] = LPAREN
	dfa[Tri{LPARENC, "X", "("}] = REJECT

	dfa[Tri{START, "X", ")"}] = REJECT
	dfa[Tri{LPAREN, "X", ")"}] = REJECT
	dfa[Tri{CCHAR, "X", ")"}] = RPAREN
	dfa[Tri{RPAREN, "X", ")"}] = REJECT
	dfa[Tri{NP, "X", ")"}] = REJECT
	dfa[Tri{COMMA, "X", ")"}] = REJECT
	dfa[Tri{CD, "X", ")"}] = REJECT
	dfa[Tri{DASH, "X", ")"}] = REJECT
	dfa[Tri{IN, "X", ")"}] = REJECT
	dfa[Tri{DT, "X", ")"}] = REJECT
	dfa[Tri{ACCEPT, "X", ")"}] = REJECT
	dfa[Tri{REJECT, "X", ")"}] = REJECT
	dfa[Tri{SYM, "X", ")"}] = RPAREN
	dfa[Tri{OTHER, "X", ")"}] = REJECT
	dfa[Tri{CSYM, "X", ")"}] = RPAREN
	dfa[Tri{LPARENC, "X", ")"}] = REJECT

	dfa[Tri{START, "X", "cd"}] = CD
	dfa[Tri{LPAREN, "X", "cd"}] = REJECT
	dfa[Tri{CCHAR, "X", "cd"}] = REJECT
	dfa[Tri{RPAREN, "X", "cd"}] = CD
	dfa[Tri{NP, "X", "cd"}] = CD
	dfa[Tri{COMMA, "X", "cd"}] = CD
	dfa[Tri{CD, "X", "cd"}] = CD
	dfa[Tri{DASH, "X", "cd"}] = CD
	dfa[Tri{IN, "X", "cd"}] = REJECT
	dfa[Tri{DT, "X", "cd"}] = REJECT
	dfa[Tri{ACCEPT, "X", "cd"}] = ACCEPT
	dfa[Tri{REJECT, "X", "cd"}] = REJECT
	dfa[Tri{SYM, "X", "cd"}] = CD
	dfa[Tri{OTHER, "X", "cd"}] = CD
	dfa[Tri{CSYM, "X", "cd"}] = CD
	dfa[Tri{LPARENC, "X", "cd"}] = REJECT

	dfa[Tri{START, "X", "np"}] = NP
	dfa[Tri{LPAREN, "X", "np"}] = REJECT
	dfa[Tri{CCHAR, "X", "np"}] = REJECT
	dfa[Tri{RPAREN, "X", "np"}] = NP
	dfa[Tri{NP, "X", "np"}] = NP
	dfa[Tri{COMMA, "X", "np"}] = NP
	dfa[Tri{CD, "X", "np"}] = NP
	dfa[Tri{DASH, "X", "np"}] = NP
	dfa[Tri{IN, "X", "np"}] = NP
	dfa[Tri{DT, "X", "np"}] = NP
	dfa[Tri{ACCEPT, "X", "np"}] = ACCEPT
	dfa[Tri{REJECT, "X", "np"}] = REJECT
	dfa[Tri{SYM, "X", "np"}] = NP
	dfa[Tri{OTHER, "X", "np"}] = NP
	dfa[Tri{CSYM, "X", "np"}] = NP
	dfa[Tri{LPARENC, "X", "np"}] = REJECT

	dfa[Tri{START, "X", "dt"}] = DT
	dfa[Tri{LPAREN, "X", "dt"}] = REJECT
	dfa[Tri{CCHAR, "X", "dt"}] = REJECT
	dfa[Tri{RPAREN, "X", "dt"}] = DT
	dfa[Tri{NP, "X", "dt"}] = DT
	dfa[Tri{COMMA, "X", "dt"}] = DT
	dfa[Tri{CD, "X", "dt"}] = DT
	dfa[Tri{DASH, "X", "dt"}] = DT
	dfa[Tri{IN, "X", "dt"}] = DT
	dfa[Tri{DT, "X", "dt"}] = REJECT
	dfa[Tri{ACCEPT, "X", "dt"}] = REJECT
	dfa[Tri{REJECT, "X", "dt"}] = REJECT
	dfa[Tri{SYM, "X", "dt"}] = REJECT
	dfa[Tri{OTHER, "X", "dt"}] = DT
	dfa[Tri{CSYM, "X", "dt"}] = DT
	dfa[Tri{LPARENC, "X", "dt"}] = REJECT

	dfa[Tri{START, "X", "in"}] = REJECT
	dfa[Tri{LPAREN, "X", "in"}] = REJECT
	dfa[Tri{CCHAR, "X", "in"}] = REJECT
	dfa[Tri{RPAREN, "X", "in"}] = REJECT
	dfa[Tri{NP, "X", "in"}] = IN
	dfa[Tri{COMMA, "X", "in"}] = REJECT
	dfa[Tri{CD, "X", "in"}] = IN
	dfa[Tri{DASH, "X", "in"}] = REJECT
	dfa[Tri{IN, "X", "in"}] = REJECT
	dfa[Tri{DT, "X", "in"}] = REJECT
	dfa[Tri{ACCEPT, "X", "in"}] = REJECT
	dfa[Tri{REJECT, "X", "in"}] = REJECT
	dfa[Tri{SYM, "X", "in"}] = REJECT
	dfa[Tri{OTHER, "X", "in"}] = IN
	dfa[Tri{CSYM, "X", "in"}] = IN
	dfa[Tri{LPARENC, "X", "in"}] = REJECT

	dfa[Tri{START, "X", "--"}] = REJECT
	dfa[Tri{LPAREN, "X", "--"}] = REJECT
	dfa[Tri{CCHAR, "X", "--"}] = REJECT
	dfa[Tri{RPAREN, "X", "--"}] = REJECT
	dfa[Tri{NP, "X", "--"}] = DASH
	dfa[Tri{COMMA, "X", "--"}] = REJECT
	dfa[Tri{CD, "X", "--"}] = DASH
	dfa[Tri{DASH, "X", "--"}] = REJECT
	dfa[Tri{IN, "X", "--"}] = REJECT
	dfa[Tri{DT, "X", "--"}] = REJECT
	dfa[Tri{ACCEPT, "X", "--"}] = REJECT
	dfa[Tri{REJECT, "X", "--"}] = REJECT
	dfa[Tri{SYM, "X", "--"}] = DASH
	dfa[Tri{OTHER, "X", "--"}] = REJECT
	dfa[Tri{CSYM, "X", "--"}] = DASH
	dfa[Tri{LPARENC, "X", "--"}] = REJECT

	dfa[Tri{START, "X", ","}] = REJECT
	dfa[Tri{LPAREN, "X", ","}] = REJECT
	dfa[Tri{CCHAR, "X", ","}] = REJECT
	dfa[Tri{RPAREN, "X", ","}] = REJECT
	dfa[Tri{NP, "X", ","}] = COMMA
	dfa[Tri{COMMA, "X", ","}] = REJECT
	dfa[Tri{CD, "X", ","}] = COMMA
	dfa[Tri{DASH, "X", ","}] = REJECT
	dfa[Tri{IN, "X", ","}] = REJECT
	dfa[Tri{DT, "X", ","}] = REJECT
	dfa[Tri{ACCEPT, "X", ","}] = REJECT
	dfa[Tri{REJECT, "X", ","}] = REJECT
	dfa[Tri{SYM, "X", ","}] = REJECT
	dfa[Tri{OTHER, "X", ","}] = REJECT
	dfa[Tri{CSYM, "X", ","}] = COMMA
	dfa[Tri{LPARENC, "X", ","}] = REJECT

	dfa[Tri{START, "X", "."}] = REJECT
	dfa[Tri{LPAREN, "X", "."}] = REJECT
	dfa[Tri{CCHAR, "X", "."}] = REJECT
	dfa[Tri{RPAREN, "X", "."}] = REJECT
	dfa[Tri{NP, "X", "."}] = ACCEPT
	dfa[Tri{COMMA, "X", "."}] = REJECT
	dfa[Tri{CD, "X", "."}] = ACCEPT
	dfa[Tri{DASH, "X", "."}] = REJECT
	dfa[Tri{IN, "X", "."}] = REJECT
	dfa[Tri{DT, "X", "."}] = REJECT
	dfa[Tri{ACCEPT, "X", "."}] = REJECT
	dfa[Tri{REJECT, "X", "."}] = REJECT
	dfa[Tri{SYM, "X", "."}] = REJECT
	dfa[Tri{OTHER, "X", "."}] = REJECT
	dfa[Tri{CSYM, "X", "."}] = REJECT
	dfa[Tri{LPARENC, "X", "."}] = REJECT

	dfa[Tri{START, "X", "sym"}] = SYM
	dfa[Tri{LPAREN, "X", "sym"}] = REJECT
	dfa[Tri{CCHAR, "X", "sym"}] = REJECT
	dfa[Tri{RPAREN, "X", "sym"}] = REJECT
	dfa[Tri{NP, "X", "sym"}] = SYM
	dfa[Tri{COMMA, "X", "sym"}] = REJECT
	dfa[Tri{CD, "X", "sym"}] = SYM
	dfa[Tri{DASH, "X", "sym"}] = REJECT
	dfa[Tri{IN, "X", "sym"}] = REJECT
	dfa[Tri{DT, "X", "sym"}] = REJECT
	dfa[Tri{ACCEPT, "X", "sym"}] = ACCEPT
	dfa[Tri{REJECT, "X", "sym"}] = REJECT
	dfa[Tri{SYM, "X", "sym"}] = REJECT
	dfa[Tri{OTHER, "X", "sym"}] = REJECT
	dfa[Tri{CSYM, "X", "sym"}] = SYM
	dfa[Tri{LPARENC, "X", "sym"}] = REJECT

	dfa[Tri{START, "X", "cc"}] = REJECT
	dfa[Tri{LPAREN, "X", "cc"}] = REJECT
	dfa[Tri{CCHAR, "X", "cc"}] = REJECT
	dfa[Tri{RPAREN, "X", "cc"}] = REJECT
	dfa[Tri{NP, "X", "cc"}] = OTHER
	dfa[Tri{COMMA, "X", "cc"}] = REJECT
	dfa[Tri{CD, "X", "cc"}] = OTHER
	dfa[Tri{DASH, "X", "cc"}] = REJECT
	dfa[Tri{IN, "X", "cc"}] = OTHER
	dfa[Tri{DT, "X", "cc"}] = OTHER
	dfa[Tri{ACCEPT, "X", "cc"}] = REJECT
	dfa[Tri{REJECT, "X", "cc"}] = REJECT
	dfa[Tri{SYM, "X", "cc"}] = REJECT
	dfa[Tri{OTHER, "X", "cc"}] = REJECT
	dfa[Tri{CSYM, "X", "cc"}] = OTHER
	dfa[Tri{LPARENC, "X", "cc"}] = REJECT

	dfa[Tri{START, "X", "X"}] = REJECT
	dfa[Tri{LPAREN, "X", "X"}] = REJECT
	dfa[Tri{CCHAR, "X", "X"}] = REJECT
	dfa[Tri{RPAREN, "X", "X"}] = REJECT
	dfa[Tri{NP, "X", "X"}] = ACCEPT
	dfa[Tri{COMMA, "X", "X"}] = REJECT
	dfa[Tri{CD, "X", "X"}] = ACCEPT
	dfa[Tri{DASH, "X", "X"}] = REJECT
	dfa[Tri{IN, "X", "X"}] = REJECT
	dfa[Tri{DT, "X", "X"}] = REJECT
	dfa[Tri{ACCEPT, "X", "X"}] = REJECT
	dfa[Tri{REJECT, "X", "X"}] = REJECT
	dfa[Tri{SYM, "X", "X"}] = REJECT
	dfa[Tri{OTHER, "X", "X"}] = REJECT
	dfa[Tri{CSYM, "X", "X"}] = REJECT
	dfa[Tri{LPARENC, "X", "X"}] = REJECT

	return symbols, dfa
}

// Given a string this will return the copyright notice
// of that string if it exists, if not the empty string is returned
// The string must be tagged and propperly delimited
func (copyrightTagger *Tagger) Extract(inBytes []byte) string {
	// Before I can match for copyright notice I need the sentence tagged
	var taggedSent []TaggedWord
	taggedSent = copyrightTagger.TagBytes(inBytes)

	currentState := REJECT
	var potentialNotice []TaggedWord = make([]TaggedWord, 0)
	var extractedNotice []TaggedWord = make([]TaggedWord, 0)
	for _, taggedWord := range taggedSent {

		// Is what I have good enough to add to the extracted Notices
		if currentState == ACCEPT {
			extractedNotice = append(extractedNotice, potentialNotice...)
			potentialNotice = nil
		}
		// Transition to the next state given current 'input'
		if strings.ToLower(taggedWord.word) == "copyright" || strings.ToLower(taggedWord.word) == "c" {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, strings.ToLower(taggedWord.word), taggedWord.tag}]
		} else if strings.Contains(taggedWord.word, "©") {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "©", "sym"}]
		} else if strings.Contains(copyrightTagger.CopyrightSyms, taggedWord.tag) {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "X", taggedWord.tag}]
		} else {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "X", "X"}]
		}
		// Because of multiple notices right after the other here's a check...
		if currentState == START || currentState == LPAREN || currentState == CSYM {
			if len(potentialNotice) > 3 { // Does it seem like something useful has been captured
				extractedNotice = append(extractedNotice, potentialNotice...)
				potentialNotice = nil
				potentialNotice = append(potentialNotice, taggedWord)
			} else {
				potentialNotice = nil
				potentialNotice = append(potentialNotice, taggedWord)
			}
		} else if currentState != REJECT { // && currentState != ACCEPT
			potentialNotice = append(potentialNotice, taggedWord)
		}
	}

	// Do a final check to see if I might have a notice as the very last part of the string
	// Be a little more vauge here to be safe
	if currentState == ACCEPT || currentState == CD || currentState == NP || len(potentialNotice) > 3 {
		extractedNotice = append(extractedNotice, potentialNotice...)
	}

	if len(extractedNotice) < 1 {
		return ""
	}
	return toString(extractedNotice)
}

// similar to the regex findAllIndex, will return the byte offsets
func (copyrightTagger *Tagger) FindAllIndex(inBytes []byte) [][]int {
	// Before I can match for copyright notice I need the sentence tagged
	var taggedSent []TaggedWord
	taggedSent = copyrightTagger.TagBytes(inBytes)

	//Return array of indicies
	var indicies = make([][]int, 0)

	currentState := REJECT
	var potentialNotice []TaggedWord = make([]TaggedWord, 0)
	for _, taggedWord := range taggedSent {

		// Is what I have good enough to add to the extracted Notices
		if currentState == ACCEPT {
			indicies = append(indicies, []int{potentialNotice[0].byteStart, taggedWord.byteStart})
			potentialNotice = nil
		}
		// Transition to the next state given current 'input'
		if strings.ToLower(taggedWord.word) == "copyright" || strings.ToLower(taggedWord.word) == "c" {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, strings.ToLower(taggedWord.word), taggedWord.tag}]
		} else if strings.Contains(taggedWord.word, "©") {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "©", "sym"}]
		} else if strings.Contains(copyrightTagger.CopyrightSyms, taggedWord.tag) {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "X", taggedWord.tag}]
		} else {
			currentState = copyrightTagger.CopyrightDFA[Tri{currentState, "X", "X"}]
		}
		// Because of multiple notices right after the other here's a check...
		if currentState == START || currentState == LPAREN || currentState == CSYM {
			if len(potentialNotice) > 3 { // Does it seem like something useful has been captured
				indicies = append(indicies, []int{potentialNotice[0].byteStart, taggedWord.byteStart})
				potentialNotice = nil
				potentialNotice = append(potentialNotice, taggedWord)
			} else {
				potentialNotice = nil
				potentialNotice = append(potentialNotice, taggedWord)
			}
		} else if currentState != REJECT { // && currentState != ACCEPT
			potentialNotice = append(potentialNotice, taggedWord)
		}
	}

	// Do a final check to see if I might have a notice as the very last part of the string
	// Be a little more vauge here to be safe
	if currentState == ACCEPT || currentState == CD || currentState == NP || len(potentialNotice) > 3 {
		indicies = append(indicies, []int{potentialNotice[0].byteStart, potentialNotice[len(potentialNotice)-1].byteStart})
	}

	return indicies
}
