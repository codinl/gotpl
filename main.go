package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codinl/gotpl/gotpl"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: gotpl [-debug] [-watch] <input dir> <output dir>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = Usage
	isDebug := flag.Bool("debug", false, "use debug mode")
	isWatch := flag.Bool("watch", false, "use watch mode")

	flag.Parse()

	options := gotpl.Option{}

	if *isDebug {
		options["Debug"] = *isDebug
	}
	if *isWatch {
		options["Watch"] = *isWatch
	}

	options["Debug"] = false
	options["Watch"] = false

	if len(flag.Args()) != 2 {
		flag.Usage()
	}

	input, output := flag.Arg(0), flag.Arg(1)
	stat, err := os.Stat(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if stat.IsDir() {
		err := gotpl.Generate(input, output, options)
		if err != nil {
			fmt.Println(err)
		}
	}
}
