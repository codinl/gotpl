package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/codinl/gotpl/template"
	"github.com/codinl/go-logger"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: gotpl [-debug] [-watch] <input dir> <output dir>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	err := initLogger()
	if err != nil {
		panic("initLogger fail")
	}

	flag.Usage = Usage
	isDebug := flag.Bool("debug", false, "use debug mode")
	isWatch := flag.Bool("watch", false, "use watch mode")

	flag.Parse()

	option := gotpl.Option{}

	if *isDebug {
		option["Debug"] = *isDebug
	}
	if *isWatch {
		option["Watch"] = *isWatch
	}

	option["Debug"] = false
//	option["Debug"] = true
	option["Watch"] = true

	if len(flag.Args()) != 2 {
		flag.Usage()
	}

	input, output := flag.Arg(0), flag.Arg(1)
//	input, output := "./tpl/", "./gen/"
	stat, err := os.Stat(input)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	if stat.IsDir() {
		err := gotpl.Generate(input, output, option)
		if err != nil {
			logger.Error(err)
		}
	}
}

func initLogger() error {
	err := logger.Init("./log", "gotpl.log", logger.DEBUG)
	if err != nil {
		fmt.Println("logger init error err=", err)
		return err
	}
	logger.SetConsole(true)
	return nil
}
