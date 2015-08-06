# Tagger Package
This a golang package designed around taking in a raw text corpus, reading
in that corpus and then tagging the parts of speech for any given input
byte slice given that the tagger has already read in the corpus. The main
go file is the tagger.go and this contains the creation of the tagger
and the functions for tagging a slice of bytes. This specific tagger works
off of the verterbi algorithm and when splitting a word will split on all
symbols which could be bad for possessives, contractions, compounds, and others
but could be easily modified to split on specific symbols.

New( path to corpus for tagging (string) );

	Takes the path to the corpus to create the tagger module from. This
	must be a string and this will return an initialized tagger module
	that the following functions can be called on.

TagBytes( raw byte slice );

	Returns a slice of Tagged Word objects that have the word, part of
	speech tag, and the byte offeset in the original slice.


# Tagger Package for copyrights
This package was developed specifically for copyright notice detection;
 however, the copyright extraction and the part of speech tagging are completely
 separate from eachother meaning that any different modules/packages can
 be easily inserted, hacked, or in all possible manners merged together
 to perform other NLP functionality after the tagging is done.

# Important
The tagger.go is mostly separated from the copyright.go part of the package
except for a small optimization where the tagger module created with the
New function will save the information needed for the copyright functions
these part could be easily removed for other projects if needed/wanted.


