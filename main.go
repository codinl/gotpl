package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/codinl/gotpl/template"
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

	option := template.Option{}

	if *isDebug {
		option["Debug"] = *isDebug
	}
	if *isWatch {
		option["Watch"] = *isWatch
	}

	option["Debug"] = true
//	options["Watch"] = true

//	if len(flag.Args()) != 2 {
//		flag.Usage()
//	}

//	input, output := flag.Arg(0), flag.Arg(1)
//	input, output := "./tpl/", "./gennew/"
	input, output := "./tpl/", "./gen/"
	stat, err := os.Stat(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if stat.IsDir() {
		err := template.Generate(input, output, option)
		if err != nil {
			fmt.Println(err)
		}
	}
}
