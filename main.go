package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Initialize flag variables
var (
	flagN, flagB, flagE, flagT, flagV, flagA, flagE2, flagT2, help bool
)

// Assign flag to the corresponding variables
func init() {
	flag.BoolVar(&flagN, "n", false, "Number all lines")
	flag.BoolVar(&flagB, "b", false, "Number non empty lines")
	flag.BoolVar(&flagE, "E", false, "Display $ at the end of the line")
	flag.BoolVar(&flagT, "T", false, "Display ^I instead of tab")
	flag.BoolVar(&flagV, "v", false, "Show unprintable characters")
	flag.BoolVar(&flagA, "A", false, "Equivalent to -vET")
	flag.BoolVar(&flagE2, "e", false, "Equivalent to -vE")
	flag.BoolVar(&flagT2, "t", false, "Equivalent to -vT")
	flag.BoolVar(&help, "h", false, "show help")
}

func main() {
	// Pass command line arguments
	flag.Parse()
	if help || !flag.Parsed(){
		flag.Usage()
		os.Exit(0)
	}

	if flagA{
		flagV = true
		flagE = true
		flagT = true
	}
	if flagE2{
		flagV = true
		flagE = true
	}
	if flagT2{
		flagV = true
		flagT = true
	}

	// Getting files from cmd line
	files := flag.Args()
	if len(files) == 0{
		files = []string{"-"}
	}

	lineNum := 0
	// Go through the file names
	for _, fname := range files{
		var r io.Reader
		if fname == "-"{
			r = os.Stdin
		} else{
			file, err := os.Open(fname)
			if err != nil{
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			defer file.Close()
			r = file
		}

		// Now print the content on each file
		printFile(r, &lineNum)
	}
}

// Go through the file
func printFile(r io.Reader, lineNum *int){
	scanner := bufio.NewScanner(r)

	// Line by line
	for scanner.Scan(){
		line := scanner.Text()
		printLineNum := flagN || (flagB && strings.TrimSpace(line) != "")
		
		if printLineNum{
			*lineNum++
			fmt.Printf("%6d\t", *lineNum)
		}

		// go through the line
		for _, c := range line{
			switch{
			case c == '\t' && flagT:
				fmt.Print("^I")
			case flagV && (!unicode.IsPrint(c) && c != '\t'):
				fmt.Print(printCaretNotation(c))
			default:
				fmt.Print(string(c))	
			}
		}

		if flagE{
			fmt.Print("$")
		}
		fmt.Print("\n")
	}	
}

// To display non printable characters with corresponding ASCII
func printCaretNotation(c rune) string{
	switch {
	case c < 32:
		return fmt.Sprintf("^%c", c+'@') // ^@ to ^_
	case c == 127:
		return "^?"
	case c > 127 && c < 160:
		return fmt.Sprintf("M-^%c", c-128+'@')
	case c >= 160 && c <= 255:
		return fmt.Sprintf("M-%c", c-128)
	default:
		return string(c)
	}
}