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

	// Getting args from cmd line
	args := flag.Args()
	if len(args) > 0 && (args[0] == ">" || args[0] == ">>"){
		operator := args[0]
		outputFileName := args[1]
		inputFiles := args[2:]
		if len(args) < 3{
			updateSingleFile(operator, outputFileName)
			return
		} else if len(args) > 2{
			mergeFiles(operator, outputFileName, inputFiles)
			return
		}
	}

	processFiles(args)
}

// Helper function to open the output file based on the operator
func openOutputFile(operator, outputFileName string) (*os.File, error) {
	switch operator {
	case ">":
		return os.Create(outputFileName) // Overwrite mode
	case ">>":
		return os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Append mode
	default:
		return nil, fmt.Errorf("invalid operator %s, use > for overwrite or >> for append", operator)
	}
}

// Update a single file
func updateSingleFile(operator, outputFileName string) {
	// Open the output file
	file, err := openOutputFile(operator, outputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file %s: %v\n", outputFileName, err)
		return
	}
	defer file.Close()

	// Read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to the file %s: %v\n", outputFileName, err)
			return
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning stdin: %v\n", err)
		return
	}

}




// Merge multiple files
func mergeFiles(operator string, outputFileName string, inputFiles []string){

}

func processFiles(files []string){
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
			case flagV && ((!unicode.IsPrint(c) && c != '\t') || c >=160 && c<=255):	// Check whether the char is unprintable or meta extended
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
func printCaretNotation(c rune) string {
	switch {
	case c < 32:
		return fmt.Sprintf("^%c", c+'@') 
	case c == 127:
		return "^?"
	case c >= 128 && c <= 159:
		return fmt.Sprintf("M-^%c", c-128+'@') 
	case c >= 160 && c <= 255:
		return fmt.Sprintf("M-%c", c-128) 
	default:
		return "O"
	}
}
