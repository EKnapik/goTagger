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
	"fmt"
	"testing"
)

// This function is the testing function and can be the only one
// because the others need to have something passed into them
// go test does not state which functions are run when and in what
// order so having multiple functions with possible uninitialized
// variables would be an impropper use case or test case because
// I am not controlling when those said functions are given that value
func TestMain(t *testing.T) {
	// There is not a point in giving the MkTagger function something that is not
	// a pointer because the compiler will not let the code compile if a pointer
	// is not given
	copyrightTagger := New("CopyrightCorpus.in")
	/*
		// print out the trans matrix
		for row:=0; row < numOfTags; row++ {
			for col:=0; col < numOfTags; col++ {
				fmt.Printf( "%.2f  ", copyrightTagger.TransMatrix[row][col] )
			}
			fmt.Print( "\n" )
		}
		fmt.Print( "\n\nTHE WORD DICTIONARY: \n" )
		// print the dictionary
		for key := range copyrightTagger.Dictionary {
			for _, tagObject := range copyrightTagger.Dictionary[key] {
				fmt.Printf( "%s->%s: %.2f\n", key, tagObject.tag, tagObject.freq )
			}
			fmt.Print( "\n" )
		}
		fmt.Print( "\n" )
	*/

	// Given a sentence to tag lets see what it prints out:
	copyrightSymbol := copyrightTagger.Match([]byte("./* Decomposed printf argument list.\n (C) 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it. and/or modify\n it under the terms of the GNU General Public License. as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */."))
	full := copyrightTagger.Match([]byte("/* Decomposed printf argument list.\n Copyright (C) 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it and/or modify\n it under the terms of the GNU General Public License as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */"))
	half := copyrightTagger.Match([]byte("/* Decomposed printf argument list.\n Copyright 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it and/or modify\n it under the terms of the GNU General Public License as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */"))
	none := copyrightTagger.Match([]byte("/* Decomposed printf argument list.\n 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it and/or modify\n it under the terms of the GNU General Public License as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */"))
	lowerCase := copyrightTagger.Match([]byte("/* Decomposed printf argument list.\n (c) 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it and/or modify\n it under the terms of the GNU General Public License as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */"))
	flipped := copyrightTagger.Match([]byte("/* Decomposed printf argument list.\n laksjdf laskdj f;l © Copyright 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it and/or modify\n it under the terms of the GNU General Public License as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */"))
	weird := copyrightTagger.Extract([]byte("Copyright ( C ) 2007 Free Software Foundation , Inc. ( ( copyright ( ( copyright ( ( ( ( ( ( ( copyright ( ( ( ( ( ( ( ( ( ( ( ( copyright ( ( copyright ( ( ( ( ( ( ( ( Copyright ( C ) < ( < < Copyright ( C ) < ("))
	extractString := copyrightTagger.Extract([]byte("/* Decomposed printf argument list.\n laksjdf laskdj f;l © Copyright 1999, 2002-2003, 2005-2007, 2009-2011 Free Software\n    Foundation, Inc.\n\n  This program is free software; you can redistribute it and/or modify\n it under the terms of the GNU General Public License as published by\n  the Free Software Foundation; either version 3.1.2, or 9.3 (at your option)\nany later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n  You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.  */"))

	passTest1 := copyrightTagger.Match([]byte("some stuff here. \\(co   Exablox  and Pixar     2018 with the Datto corp."))
	passTest2 := copyrightTagger.Match([]byte("Copyright (c) IBM       Corporation, 2003,   2008.  All rights reserved.   --")) //saying to save tagged array
	passTest2Ext := copyrightTagger.Extract([]byte("Copyright (c) IBM       Corporation, 2003,   2008.  All rights reserved.   --"))
	passTest3 := copyrightTagger.Match([]byte(" Â© 2001-2014 Python Software Foundation</string>")) // staying to save the tagged array
	passTest4 := copyrightTagger.Match([]byte(" Â© 2001-2014 Python Software Foundation</string>"))
	passTest3Ext := copyrightTagger.Extract([]byte(" © 2001-2014 Python Software Foundation</string>"))
	complexExtract := copyrightTagger.Extract([]byte("some stuff here. \\(co Exablox and Pixar 2018 with the Datto corp. In accordance with this laa balh"))
	passTest5 := copyrightTagger.Extract([]byte(" Â© 2001-2014 Python Software Foundation</string>"))
	passTest6 := copyrightTagger.Match([]byte("* Copyright 2004 by Theodore Ts'o."))
	passTest7 := copyrightTagger.Extract([]byte("Copyright (c) 2007, 2008 Alastair Houghton"))

	failTest1 := copyrightTagger.Match([]byte(" #define Copyright sign "))
	failTest2 := copyrightTagger.Match([]byte("Fetched %sB in %s (%sB/s)\n"))
	failTest3 := copyrightTagger.Match([]byte("GNU nano version %s (compiled %s, %s)\n"))
	failTest4 := copyrightTagger.Match([]byte("(c)   "))
	failTest5 := copyrightTagger.Match([]byte("#define c_tolower(c) \\ "))
	failTest6 := copyrightTagger.Match([]byte("/* ToUnicode().  May realloc() utf8in.  Will free utf8in unconditionally. */"))
	failTest7 := copyrightTagger.Match([]byte("COPYRIGHT SIGN */ (1U<<_CC_GRAPH)|(1U<<_CC_PRINT)|(1U<<_CC_QUOTEMETA),"))
	failTest8 := copyrightTagger.Extract([]byte("Copyright\\ 1989% -1990\\ PKWARE\\ Inc.	Self-extracting PKZIP archive"))
	failTest9 := copyrightTagger.Match([]byte("( C ( { ( { ( { ( ( { ( ( { ( { < < ( ( ) 0"))
	failTest10 := copyrightTagger.Extract([]byte("< < ( ( ( ( ( ( ( C ( ( C ( ( ( ( ( ( ( ( ( < < < < < < < < < < < < < < < < ( ( ( ( ( ( ( < ( Copyright ( C ) 1992 - 2009 , Free Software Foundation , Inc. - copyright ( < ( ( ( ( ( < ( < ( ( ( ( ( < ( < ( ( ( ( ( ("))
	failTest11 := copyrightTagger.Match([]byte(" # '$siteCopyrightName' on line 12, col 24"))
	failTest13 := copyrightTagger.Extract([]byte("# ifdef _SC_PAGESIZE\n#  define getpagesize() sysconf(_SC_PAGESIZE)\n# else /* no _SC_PAGESIZE */\n#  ifdef HAVE_SYS_PARAM_H\n#   include <sys/param.h>\n#   ifdef EXEC_PAGESIZE\n#    define getpagesize() EXEC_PAGESIZE\n#   else /* no EXEC_PAGESIZE */\n#    ifdef NBPG\n#     define getpagesize() NBPG * CLSIZE\n#     ifndef CLSIZE\n#      define CLSIZE 1\n#     endif /* no CLSIZE */\n#    else /* no NBPG */\n#     ifdef NBPC\n#      define getpagesize() NBPC\n#     else /* no NBPC */\n#      ifdef PAGESIZE\n#       define getpagesize() PAGESIZE\n#      endif /* PAGESIZE */\n#     endif /* no NBPC */\n#    endif /* no NBPG */\n#   endif /* no EXEC_PAGESIZE */\n#  else /* no HAVE_SYS_PARAM_H */\n#   define getpagesize() 8192	/* punt totally */\n#  endif /* no HAVE_SYS_PARAM_H */\n# endif /* no _SC_PAGESIZE */"))
	failTest14 := copyrightTagger.Extract([]byte("#???????k(?P(???-[?477?????????=?|?ʡ?zRRR?p?B?|^^ޚ5k\n                                                                                                                               RQ??+W?d?͛7?\\?R?P?????}?5@4?W_ϔ?8J`?V>?h?V?\n                                 6?J??????V(x?}?l???)HC??̟?+???p?f>/?5k\n                                                                      ?8+W\"??y3+W??@?>?oaK3?𢡊?????|?ʡ?zRB?Bj???[??T`%+e?6?y%+?cCq?ĉI?&????<yR?G2???ѣ?|?j?Pv=?{?nFF?F>&&?7?0?=?ĉ>?8q?':????tݸq???a??ūV?1b???^TT4iҤ??<?\n                                                                              EEE?_?>66V;?ر#F8p???~????????W_?O7n?,?U0⑵&A??'?`=?v??S????е?Ϡ?                                                                                                           8?x1?V1b??~?\"&MB}??????v??v\n  ?MA??II!;?????op???,^Ū?P?؋(?Ĥ<?1E?z?????Ʊ?                                                                   ?!???M?W]????z\n\n                                                                                                                             ?1cƜ:u???V?>ح[?8???D???&L?????Ϛ5ˀ?????_?????????7??ꕱ?qNN???pjj?T*ճ>???>}??/???????wU??	=?!?Q??T??ݣC?H???<?B?*??~\"??????ޙ75?T??\"z?3}:_|?????tm?q9?v©?J???"))

	fmt.Print("\n")
	fmt.Println("True: ", copyrightSymbol)
	fmt.Println("True: ", full)
	fmt.Println("True: ", half)
	fmt.Println("True: ", lowerCase)
	fmt.Println("True: ", flipped)
	fmt.Println("True: ", passTest1)
	fmt.Println("True: ", passTest2)
	fmt.Println("True: ", passTest3)
	fmt.Println("True: ", passTest4)
	fmt.Println("True: ", passTest6)

	fmt.Println("False: ", none)
	fmt.Println("False: ", failTest1)
	fmt.Println("False: ", failTest2)
	fmt.Println("False: ", failTest3)
	fmt.Println("False: ", failTest4)
	fmt.Println("False: ", failTest5)
	fmt.Println("False: ", failTest6)
	fmt.Println("False: ", failTest7)
	fmt.Println("False: ", failTest9)
	fmt.Println("False: ", failTest11)

	fmt.Println("Extracting")
	fmt.Println(extractString)
	fmt.Println(complexExtract)
	fmt.Println(passTest2Ext)
	fmt.Println(passTest3Ext)
	fmt.Println(weird)
	fmt.Println(passTest5)
	fmt.Println(failTest10)
	fmt.Println(failTest13)
	fmt.Println("Odd symbols: ", failTest14)
	fmt.Println("Start: ", passTest7)

	fmt.Println("CURIOUS RESULTS:")
	fmt.Println("Should Fail: ", failTest7)
	fmt.Println("Should pass: ", failTest8)

	fmt.Println("FINDING INDICIES")
	raw := []byte("It's an MIT-style license.  Here goes:\n\nCopyright (c) 2007, 2008 Alastair Houghton\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:")
	taggedRaw := copyrightTagger.TagBytes(raw)
	indiciesTest1 := copyrightTagger.FindAllIndex(raw)
	fmt.Println(taggedRaw)
	fmt.Println(indiciesTest1)

}
