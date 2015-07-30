# Determine Copyright Within Raw Text

## Overview of Code
The overal idea is that given a piece of raw text the program will either use
its initialized tagging model or use the one it is given to:
  * Separate the sentence so it can be proccessed
  * Tag the sentence with its appropriate part of speech
  * The sentence is now in a word|~|tag pair and is delimited propperly
  * Make and run the compression of specific tag sequences on the tagged sentence
  * Make and run the Copyright notice detection DFA on the Tagged Sentence
    * This currently returns true is the sentence contains a copyright notice
    * However that is determined by seeing if the returned/extracted notice is the empty string

## How the Trainer Works
The program loads in the corpus with word|~|tag pairs then walks through each sentence. Adding
each instance of the word with its tag to the unigram word dictionary. The dictionary is a key
value pair where the key, the word that I see corresponds to a value which is slice (golang) of
TagFrequency objects (structs). What this slice is the tag and the amount of times I have seen that
tag before. The increment unigram word dictionary function in the Tagger package might help visualize
it. While that is happeing the transition matrix is being formed by remembering the previous tag to
the current tag transition. Every time I move from tag A to tag B I increment appropriatelly in the
transition matrix, which can be read as moving from the row corresponding to tag A to the column
corresponding to tag B, and accessed by transMatrix[indexOfTagA][indexOfTagB]. Once I increment
the occurance of every word and its tag and the transitions from one tag to another finishing with
the given corpus the unigram word dictionary and the transition matrix are converted into a
probabilistic format count all occurances then dividing by the total, this is were I have implemented
Laplace Smoothing into the transitionMatrix making all tag transitions possible however some are very
unlikely.

## How the Tagger Works
Given the a raw formatted sentence, the unigram word dictionary and the part of speech tag transition
matrix. The tagger creates a two dimensional array for the sentence where each column corresponds to
each word in the sentence plus one and the rows are the probability that, the specific word has that tag.
To fill in the 2d matrix the Viterbi Algorithm is used relying heavily on [Bayes' Theorem.](https://en.wikipedia.org/wiki/Bayes%27_theorem)The first column must be the start of the sentence which doesnt have the word associated with it but we
know that before every sentence what the start is. So the first tag is the start of the sentence or a
period because that would have been the previous tag. Then we look at the first word determining the
probabilty of that word having each tag by the following formula:
>_P(tag|word) = P(word|tag)*max(previousWord's tag*transition to current tag)_

Which can be read as the probability that the current word has a given tag is given by the probability
that this word has this tag (because I have see the word before) times the best transitional probability
that this tag will be next because I have seen this kind of tag transition before.
For more information look at the [forward backward algorithm](https://en.wikipedia.org/wiki/Forward–backward_algorithm)
and the [viterbi algorithm](https://en.wikipedia.org/wiki/Viterbi_algorithm).
Once every word has the probability for each tag determined walk throught the array giving the tag with
the max likelyhood to the given word.

## The Corpus Part of Speech Tags
| Tag | Description           | Examples                |
|-----|-----------------------|-------------------------|
|  $  |    the dollar sign    |            $            |
|  "  |     quotation mark    |         ", ', `         |
|  (  |    open parenthesis   |        (, <, {, [       |
|  )  |   close parenthesis   |        ), >, }, ]       |
|  ,  |         comma         |            ,            |
|  -- |         dashes        |          --, -          |
|  .  |  sentence terminator  |         ., !, ?         |
|  :  |   colon or ellipsis   |        :, ;, ...        |
|  cc |      conjunction      |       and, or, but      |
|  cd |   cardinal, numeral   |     one, 123, dozen     |
|  dt |       determiner      |        all, an, a       |
|  fw |      foreign word     |    var, bool, alsdjf    |
|  jj |       adjective       |       large, small      |
|  ls |       list item       |          A.  1.         |
|  nn |          noun         |       stove, pool       |
|  np |      propper noun     |      Eric, Exablox      |
| pos |       possessive      |            `s           |
|  pr |        pronoun        |       he, she, we       |
|  rb |         adverb        |       bigger, best      |
| sym |         symbol        |        @, ^, #, /       |
|  to |      the word to      |            to           |
|  uh |      interjection     |                         |
|  vb |          verb         |     run, drive, swim    |
|  md |    modal auxiliary    | can, should, will, must |
|  in |      preposition      |     with, as, under     |
| bos | beginning of sentence |                         |


The corpus uses the delimeter of |~| between words and tags. Between tag word
pairs there are three spaces, certain languages like python do not care how
many spaces are between words. However, golang reading in the corpus makes it
difficult if each tag word pair is not separated the same way. The reason that
the corpus is all on one line is another issue with reading in raw text from
golang that dealing would require dealing with carriage returns.
