package gotpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	GO_EXT  = ".go"
	TPL_EXT = ".tpl"
	TMP_EXT = ".t"
	TMP_DIR = "tmp/"
	GEN_DIR = "testgen/"
	SEC_DIR = "sections/"
)

var Sections map[string]*Section

func Generate(input string, output string, options Option) error {
	tplDir = input

	tpl, err := buildTplTree(input)
	if err != nil {
		return err
	}

	secs, err := genSections(tplDir)
	if err == nil {
		Sections = secs
	} else {
		Sections = nil
	}

	err = tpl.genBlock()
	if err != nil {
		return err
	}

	input = input + TMP_DIR

	err = genFolder(input, output, options)
	if err != nil {
		return err
	}

	return nil
}

// Generate from input to output file,
// gofmt will trigger an error if it fails.
func GenFile(input string, output string, options Option) error {
	outdir := filepath.Dir(output)
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}
	return generate(input, output, options)
}

// Generate from directory to directory, Find all the files with extension
// of .gohtml and generate it into target dir.
func genFolder(input string, out string, options Option) (err error) {
	if !exists(input) {
		err = os.MkdirAll(input, 0755)
		if err != nil {
			return err
		}
	}

	//Make it
	if !exists(out) {
		os.MkdirAll(out, 0775)
	}

	input_abs, _ := filepath.Abs(input)
	out_abs, _ := filepath.Abs(out)

	paths := []string{}

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//Just do file with exstension .t
			if !strings.HasSuffix(path, TMP_EXT) {
				return nil
			}
			filename := filepath.Base(path)
			if strings.HasPrefix(filename, ".#") {
				return nil
			}
			paths = append(paths, path)
		}
		return nil
	}

	fun := func(path string, res chan<- string) {
		//adjust with the abs path, so that we keep the same directory hierarchy
		input, _ := filepath.Abs(path)
		output := strings.Replace(input, input_abs, out_abs, 1)
		output = strings.Replace(output, TMP_EXT, GO_EXT, -1)
		err := GenFile(path, output, options)
		if err != nil {
			res <- fmt.Sprintf("%s -> %s", path, output)
			os.Exit(2)
		}
		res <- fmt.Sprintf("%s -> %s", path, output)
	}

	err = filepath.Walk(input, visit)
	runtime.GOMAXPROCS(runtime.NumCPU())
	result := make(chan string, len(paths))

	for w := 0; w < len(paths); w++ {
		go fun(paths[w], result)
	}

	for i := 0; i < len(paths); i++ {
		<-result
	}

	//	if options["Watch"] != nil {
	//		watchDir(incdir_abs, outdir_abs, options)
	//	}

	return
}

func generate(input string, output string, Options Option) error {
	compiler, err := run(input, Options)
	if err != nil || compiler == nil {
		panic(err)
	}

	err = ioutil.WriteFile(output, []byte(compiler.buf), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("gofmt", "-w", output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("gofmt: ", err)
		return err
	}

	if Options["Debug"] == true {
		content, _ := ioutil.ReadFile(output)
		fmt.Println(string(content))
	}

	return err
}

func run(input string, Options Option) (*Compiler, error) {
	content, err := ioutil.ReadFile(input)
	if err != nil {
		return nil, err
	}
	text := string(content)

	lex := &Lexer{Text: text, Matches: TokenMatches}

	// Scan() -> tokens
	tokens, err := lex.Scan()
	if err != nil {
		return nil, err
	}

	//DEBUG
	if Options["Debug"] == true {
		fmt.Println("------------------- TOKEN START -----------------")
		for _, elem := range tokens {
			elem.P()
		}
		fmt.Println("--------------------- TOKEN END -----------------\n")
	}

	parser := &Parser{ast: &Ast{}, rootAst: nil,
		tokens: tokens, preTokens: []Token{},
		saveTextTag: false, initMode: UNK}

	// Run() -> ast
	err = parser.Run()
	if err != nil {
		fmt.Println(input, ":", err)
		os.Exit(2)
	}

	//DEBUG
	if Options["Debug"] == true {
		fmt.Println("--------------------- AST START -----------------")
		parser.ast.debug(0, 20)
		fmt.Println("--------------------- AST END -----------------\n")
		if parser.ast.Mode != PRG {
			panic("TYPE")
		}
	}
	cp := makeCompiler(parser.ast, Options, input)

	// visit() -> cp.buf
	cp.visit()

	return cp, nil
}
