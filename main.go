package main

import (
	"flag"
	"os"
)

var (
	flagN, flagB, flagE, flagT, flagV, flagA, flagE2, flagT2, help bool
)

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
	flag.Parse()
	if help || !flag.Parsed(){
		flag.Usage()
		os.Exit(0)
	}
}