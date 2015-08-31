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

// These are the unit tests for the testing package
// Run this while inside the go package tagger
// Ran by running: "go test" on the commandline

package tagger

import (
	"log"
	"os"
	"testing"
)

var copyrightTagger *Tagger

func dumpTransMatrix() {
	// print out the trans matrix
	for row:=0; row < numOfTags; row++ {
		for col:=0; col < numOfTags; col++ {
			log.Printf( "%.2f  ", copyrightTagger.TransMatrix[row][col] )
		}
		log.Print( "\n" )
	}
	log.Print( "\n\nTHE WORD DICTIONARY: \n" )
	// print the dictionary
	for key := range copyrightTagger.Dictionary {
		for _, tagObject := range copyrightTagger.Dictionary[key] {
			log.Printf( "%s->%s: %.2f\n", key, tagObject.tag, tagObject.freq )
		}
		log.Print( "\n" )
	}
	log.Print( "\n" )
}


func TestMain(m *testing.M) {
	copyrightTagger = New("CopyrightCorpus.in")

//	dumpTransMatrix()

	os.Exit(m.Run())
}

func TestMatch(t *testing.T) {
	type MatchTest struct {
		Expected	bool
		Text		string
	}


	tests := []MatchTest {
		{
			Expected:	true,
			Text:		"./* Decomposed printf argument list.\n"+
					" (C) 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n\n  This program is free software; you can redistribute it. and/or modify\n"+
					" it under the terms of the GNU General Public License. as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					" if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */.",
		},
		{
			Expected:	true,
			Text:		"/* Decomposed printf argument list.\n"+
					" Copyright (C) 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n"+
					"\n"+
					"  This program is free software; you can redistribute it and/or modify\n"+
					" it under the terms of the GNU General Public License as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					" if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */",
		},
		{
			Expected:	true,
			Text:		"/* Decomposed printf argument list.\n"+
					" Copyright 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n"+
					"\n"+
					"  This program is free software; you can redistribute it and/or modify\n"+
					" it under the terms of the GNU General Public License as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					" if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */",
		},
		{
			Expected:	false,
			Text:		"/* Decomposed printf argument list.\n"+
					" 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n"+
					"\n"+
					"  This program is free software; you can redistribute it and/or modify\n"+
					" it under the terms of the GNU General Public License as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					" if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */",
		},
		{
			Expected:	true,
			Text:		"/* Decomposed printf argument list.\n"+
					" (c) 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n"+
					"\n"+
					"  This program is free software; you can redistribute it and/or modify\n"+
					" it under the terms of the GNU General Public License as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					" if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */",
		},
		{
			Expected:	true,
			Text:		"/* Decomposed printf argument list.\n"+
					" laksjdf laskdj f;l © Copyright 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n"+
					"\n"+
					"  This program is free software; you can redistribute it and/or modify\n"+
					" it under the terms of the GNU General Public License as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					" if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */",
		},
		{
			Expected:	true,
			Text:		"some stuff here. \\(co   Exablox  and Pixar     2018 with the Datto corp.",
		},
		{
			// saying to save tagged array
			Expected:	true,
			Text:		"Copyright (c) IBM       Corporation, 2003,   2008.  All rights reserved.   --",
		},
		{
			// staying to save the tagged array
			Expected:	true,
			Text:		" Â© 2001-2014 Python Software Foundation</string>",
		},
		{
			Expected:	true,
			Text:		"* Copyright 2004 by Theodore Ts'o.",
		},
		{
			Expected:	false,
			Text:		" #define Copyright sign ",
		},
		{
			Expected:	false,
			Text:		"Fetched %sB in %s (%sB/s)\n",
		},
		{
			Expected:	false,
			Text:		"GNU nano version %s (compiled %s, %s)\n",
		},
		{
			Expected:	false,
			Text:		"(c)   ",
		},
		{
			Expected:	false,
			Text:		"#define c_tolower(c) \\ ",
		},
		{
			Expected:	false,
			Text:		"/* ToUnicode().  May realloc() utf8in.  Will free utf8in unconditionally. */",
		},
		{
			Expected:	false,
			Text:		"COPYRIGHT SIGN */ (1U<<_CC_GRAPH)|(1U<<_CC_PRINT)|(1U<<_CC_QUOTEMETA),",
		},
		{
			Expected:	false,
			Text:		"( C ( { ( { ( { ( ( { ( ( { ( { < < ( ( ) 0",
		},
		{
			Expected:	false,
			Text:		" # '$siteCopyrightName' on line 12, col 24",
		},
	}

	for i, test := range tests {
		r := copyrightTagger.Match([]byte(test.Text))
		if r != test.Expected {
			t.Errorf("Test %d: expected %v got %v", i, test.Expected, r)
		}
	}
}

/*
 * XXX - Tad: this should probably be more robust than just checking the length of the tagged words array that was returned
 */
func TestTagBytes(t *testing.T) {
	raw := "It's an MIT-style license.  Here goes:\n"+
		"\n"+
		"Copyright (c) 2007, 2008 Alastair Houghton\n"+
		"Permission is hereby granted, free of charge, to any person obtaining a copy\n"+
		"of this software and associated documentation files (the \"Software\"), to deal\n"+
		"in the Software without restriction, including without limitation the rights\n"+
		"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n"+
		"copies of the Software, and to permit persons to whom the Software is\n"+
		"furnished to do so, subject to the following conditions:"
	nexpected := 108

	twords := copyrightTagger.TagBytes([]byte(raw))
	if twords == nil {
		t.Fatalf("expected %d elements, got nil", nexpected)
	}

	if len(twords) != nexpected {
		t.Errorf("expected %d elements got %d", nexpected, len(twords))
	}
}

func TestFindAllIndex(t *testing.T) {
	raw := "It's an MIT-style license.  Here goes:\n"+
		"\n"+
		"Copyright (c) 2007, 2008 Alastair Houghton\n"+
		"Permission is hereby granted, free of charge, to any person obtaining a copy\n"+
		"of this software and associated documentation files (the \"Software\"), to deal\n"+
		"in the Software without restriction, including without limitation the rights\n"+
		"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n"+
		"copies of the Software, and to permit persons to whom the Software is\n"+
		"furnished to do so, subject to the following conditions:"
	expected := [][]int{{40, 94}}

	matches := copyrightTagger.FindAllIndex([]byte(raw))
	if matches == nil {
		t.Fatalf("expected array of indexes, got nil")
	}

	if len(matches) != len(expected) {
		t.Fatalf("expected %d matches got %d", len(expected), len(matches))
	}

	for i := 0; i < len(matches); i++ {
		m := matches[i]
		e := expected[i]
		if len(m) != len(e) {
			t.Errorf("offset %d: index length mismatch: expected %d got %d", i, len(e), len(m))
			continue
		}
		for j := 0; j < len(m); j++ {
			if m[j] != e[j] {
				t.Errorf("matches[%d][%d]: expected %d got %d", i, j, e[j], m[j])
			}
		}
	}
}

/*
 * XXX - Tad: Needs addition of pass/fail criteria
 */
func TestExtract(t *testing.T) {
	type ExtractTest struct {
		Expected	string
		Text		string
	}

	tests := []ExtractTest {
		{
			Expected:	"Copyright ( C ) 2007 Free Copyright ( C ) Copyright ( C )",
			Text:		"Copyright ( C ) 2007 Free Software Foundation , Inc."+
					" ( ( copyright ( ( copyright ( ( ( ( ( ( ( copyright ( ( ( ( ( ( ( ( ( ( ( ("+
					" copyright ( ( copyright ( ( ( ( ( ( ( ( Copyright ( C ) < ( < < Copyright ( C ) < (",
		},
		{
			Expected:	"Copyright 1999 , 2002 - 2003 , 2005 - 2007 , 2009 - 2011 Free",
			Text:		"/* Decomposed printf argument list.\n"+
					" laksjdf laskdj f;l © Copyright 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n"+
					"    Foundation, Inc.\n"+
					"\n"+
					"  This program is free software; you can redistribute it and/or modify\n"+
					" it under the terms of the GNU General Public License as published by\n"+
					"  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\n"+
					"any later version.\n"+
					"\n"+
					"This program is distributed in the hope that it will be useful,\n"+
					"but WITHOUT ANY WARRANTY; without even the implied warranty of\n"+
					"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n"+
					"GNU General Public License for more details.\n"+
					"  You should have received a copy of the GNU General Public License along with this program;"+
					"if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,"+
					" Boston, MA 02110-1301, USA.  */",
		},
		{
			Expected:	"Copyright ( c )",
			Text:		"Copyright (c) IBM       Corporation, 2003,   2008.  All rights reserved.   --",
		},
		{
			Expected:	"© 2001 - 2014 Python Software",
			Text:		" © 2001-2014 Python Software Foundation</string>",
		},
		{
			Expected:	"( c ) Exablox and Pixar 2018 with the Datto corp. In",
			Text:		"some stuff here. \\(co Exablox and Pixar 2018 with the Datto corp. In accordance with this laa balh",
		},
		{
			Expected:	"Â© 2001 - 2014 Python Software",
			Text:		" Â© 2001-2014 Python Software Foundation</string>",
		},
		{
			Expected:	"Copyright ( c ) 2007 , 2008 Alastair Houghton",
			Text:		"Copyright (c) 2007, 2008 Alastair Houghton",
		},
		{
			Expected:	"Copyright \\ 1989 % - 1990 \\ PKWARE \\ Inc. Self",
			Text:		"Copyright\\ 1989% -1990\\ PKWARE\\ Inc.	Self-extracting PKZIP archive",
		},
		{
			Expected:	"Copyright ( C ) 1992 - 2009 ,",
			Text:		"< < ( ( ( ( ( ( ( C ( ( C ( ( ( ( ( ( ( ( ( < < < < < < < < < < < < < < < < ( ( ( ( ( ( ("+
					" < ( Copyright ( C ) 1992 - 2009 , Free Software Foundation , Inc. - copyright"+
					" ( < ( ( ( ( ( < ( < ( ( ( ( ( < ( < ( ( ( ( ( (",
		},
		{
			Expected:	"",
			Text:		"# ifdef _SC_PAGESIZE\n"+
					"#  define getpagesize() sysconf(_SC_PAGESIZE)\n"+
					"# else /* no _SC_PAGESIZE */\n"+
					"#  ifdef HAVE_SYS_PARAM_H\n"+
					"#   include <sys/param.h>\n"+
					"#   ifdef EXEC_PAGESIZE\n"+
					"#    define getpagesize() EXEC_PAGESIZE\n"+
					"#   else /* no EXEC_PAGESIZE */\n"+
					"#    ifdef NBPG\n"+
					"#     define getpagesize() NBPG * CLSIZE\n"+
					"#     ifndef CLSIZE\n"+
					"#      define CLSIZE 1\n"+
					"#     endif /* no CLSIZE */\n"+
					"#    else /* no NBPG */\n"+
					"#     ifdef NBPC\n"+
					"#      define getpagesize() NBPC\n"+
					"#     else /* no NBPC */\n"+
					"#      ifdef PAGESIZE\n"+
					"#       define getpagesize() PAGESIZE\n"+
					"#      endif /* PAGESIZE */\n"+
					"#     endif /* no NBPC */\n"+
					"#    endif /* no NBPG */\n"+
					"#   endif /* no EXEC_PAGESIZE */\n"+
					"#  else /* no HAVE_SYS_PARAM_H */\n"+
					"#   define getpagesize() 8192	/* punt totally */\n"+
					"#  endif /* no HAVE_SYS_PARAM_H */\n"+
					"# endif /* no _SC_PAGESIZE */",
		},
		{
			Expected:	"",
			Text:		"#???????k(?P(???-[?477?????????=?|?ʡ?zRRR?p?B?|^^ޚ5k\n                                                                                                                               RQ??+W?d?͛7?\\?R?P?????}?5@4?W_ϔ?8J`?V>?h?V?\n                                 6?J??????V(x?}?l???)HC??̟?+???p?f>/?5k\n                                                                      ?8+W\"??y3+W??@?>?oaK3?𢡊?????|?ʡ?zRB?Bj???[??T`%+e?6?y%+?cCq?ĉI?&????<yR?G2???ѣ?|?j?Pv=?{?nFF?F>&&?7?0?=?ĉ>?8q?':????tݸq???a??ūV?1b???^TT4iҤ??<?\n                                                                              EEE?_?>66V;?ر#F8p???~????????W_?O7n?,?U0⑵&A??'?`=?v??S????е?Ϡ?                                                                                                           8?x1?V1b??~?\"&MB}??????v??v\n  ?MA??II!;?????op???,^Ū?P?؋(?Ĥ<?1E?z?????Ʊ?                                                                   ?!???M?W]????z\n\n                                                                                                                             ?1cƜ:u???V?>ح[?8???D???&L?????Ϛ5ˀ?????_?????????7??ꕱ?qNN???pjj?T*ճ>???>}??/???????wU??	=?!?Q??T??ݣC?H???<?B?*??~\"??????ޙ75?T??\"z?3}:_|?????tm?q9?v©?J???",
		},
	}

	for i, test := range tests {
		r := copyrightTagger.Extract([]byte(test.Text))
		if r != test.Expected {
			t.Errorf("Extract Test %d: expected result %q got %q", i, test.Expected, r)
		}
	}
}
