package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Generate(input string, option Option) error {
	tplMap := map[string]*Tpl{}

	paths, err := getFiles(input, TPL_EXT)
	if err != nil {
		return err
	}

	for i := 0; i < len(paths); i++ {
		path := paths[i]

		baseName := filepath.Base(path)
		name := strings.Replace(baseName, TPL_EXT, "", 1)
		fmt.Println("name=",name)
		fmt.Println("baseName=",baseName)

		tpl := &Tpl{path: path, name: name, ast: &Ast{}, tokens:[]Token{}, blocks:map[string]*Ast{}, option: option}

		err := tpl.readRaw()
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = tpl.checkExtends()
		if err != nil {
			fmt.Println(err)
			return err
		}

		tplMap[tpl.name] = tpl
	}

	for _, tpl := range tplMap {
		if tpl.parentName == "" {
			tpl.isRoot = true
			tpl.parent = nil
//			tpl.level = 0
//			Tpls[tpl.level][tpl.name] = tpl
		} else {
			tpl.parent = tplMap[tpl.parentName]
//			tpl.level = tpl.parent.level + 1
//			Tpls[tpl.level][tpl.name] = tpl
		}
	}

	for _, tpl := range tplMap {
		err = tpl.generate()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

//func Generate(input string, out string, option Option) (err error) {
//	//	if !exists(input) {
//	//		err = os.MkdirAll(input, 0755)
//	//		if err != nil {
//	//			return err
//	//		}
//	//	}
//
//	//Make it
//	if !exists(out) {
//		os.MkdirAll(out, 0775)
//	}
//
//	input_abs, _ := filepath.Abs(input)
//	out_abs, _ := filepath.Abs(out)
//
//	paths, err := getFiles(input_abs, TPL_EXT)
//	if err != nil {
//		return err
//	}
//
//	//	fun := func(path string, res chan<- string) {
//	//		//adjust with the abs path, so that we keep the same directory hierarchy
//	//		input, _ := filepath.Abs(path)
//	//		output := strings.Replace(input, input_abs, out_abs, 1)
//	//		//		output = strings.Replace(output, TMP_EXT, GO_EXT, -1)
//	//		output = strings.Replace(output, TPL_EXT, GO_EXT, -1)
//	//		err := GenFile(path, output, options)
//	//		if err != nil {
//	//			res <- fmt.Sprintf("%s -> %s", path, output)
//	//			os.Exit(2)
//	//		}
//	//		res <- fmt.Sprintf("%s -> %s", path, output)
//	//	}
//
//	//	result := make(chan string, len(paths))
//
//	//	for i := 0; i < len(paths); i++ {
//	//		<-result
//	//	}
//
//	for i := 0; i < len(paths); i++ {
//		path := paths[i]
//		input, _ := filepath.Abs(path)
//		output := strings.Replace(input, input_abs, out_abs, 1)
//		output = strings.Replace(output, TPL_EXT, GO_EXT, -1)
//
//		outdir := filepath.Dir(output)
//		if !exists(outdir) {
//			os.MkdirAll(outdir, 0775)
//		}
//
//		tpl, err := InitTpl(input, option)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		err = tpl.generate()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	//	if options["Watch"] != nil {
//	//		watchDir(incdir_abs, outdir_abs, options)
//	//	}
//
//	return
//}

type Tpl struct {
	path       string
	name       string
	isRoot	   bool
//	level      int
	parent     *Tpl
	parentName string
	raw        []byte
	result     string
	tokens     []Token
	ast        *Ast
	blocks     map[string]*Ast
	sections   map[string]*Ast
	generated  bool
	option     Option
}

//func InitTpl(name string, option Option) (*Tpl, error) {
//	tpl := &Tpl{name: name, option: option}
//
//	err := tpl.readRaw()
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//
//	return tpl, nil
//}

func (tpl *Tpl) generate() error {
	if tpl.generated {
		return nil
	}

	if tpl.isRoot {
		return tpl.gen()
	} else {
		if tpl.parent != nil {
			tpl.parent.generate()
			return tpl.gen()
		}
	}

	return nil
}

func (tpl *Tpl) gen() error {
	if tpl.generated {
		return nil
	}

	err := tpl.genToken()
	if err != nil {
		return err
	}

	err = tpl.genAst()
	if err != nil {
		return err
	}

	err = tpl.compile()
	if err != nil {
		return err
	}

//	err = tpl.fmt()
//	if err != nil {
//		return err
//	}

	err = tpl.output()
	if err != nil {
		return err
	}

	tpl.generated = true

	return nil
}

func (tpl *Tpl) genToken() error {
	lex := &Lexer{Text: string(tpl.raw), Matches: TokenMatches}

	tokens, err := lex.Scan()
	if err != nil {
		return err
	}

	tpl.tokens = tokens

	return nil
}

func (tpl *Tpl) genAst() error {
	parser := &Parser{
		ast: tpl.ast, tokens: tpl.tokens,blocks: tpl.blocks,
		preTokens: []Token{}, initMode: UNK,
	}

	// Run() -> ast
	err := parser.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (tpl *Tpl) compile() error {
	//	dir := filepath.Base(filepath.Dir(input))
	//	file := strings.Replace(filepath.Base(input), TPL_EXT, "", 1)
	//	if options["NameNotChange"] == nil {
	//		file = Capitalize(file)
	//	}

	cp := &Compiler{
		ast: tpl.ast, buf: "", firstNode: 0,
		params: []string{}, parts: []Part{},
		imports: map[string]bool{},
		//		options: options,
		//		dir:     dir,
		fileName: tpl.name,
	}

	// visit() -> cp.buf
	cp.visit()

	tpl.result = cp.buf

	return nil
}

func (tpl *Tpl) readRaw() error {
	raw, err := ioutil.ReadFile(tpl.path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tpl.raw = raw

	return nil
}

func (tpl *Tpl) checkExtends() error {
	if bytes.HasPrefix(tpl.raw, []byte("@extends")) {
		lineEnd := -1
		for i := len("@extends"); i < len(tpl.raw); i++ {
			if tpl.raw[i] == '\n' {
				lineEnd = i
				break
			}
		}
		line := tpl.raw[:lineEnd+1]
		ss := strings.Split(string(line), " ")
		if len(ss) == 2 {
			parentName := ss[1]
			fmt.Println("------parentName====", parentName)
			tpl.raw = tpl.raw[lineEnd+1:]
		}
	}
	return nil
}

func (tpl *Tpl) getOutPath() string {
	return "./gen/" + tpl.name + ".go"
}

func (tpl *Tpl) fmt() error {
	cmd := exec.Command("gofmt", "-w", tpl.getOutPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("gofmt: ", err)
		return err
	}

	return nil
}

func (tpl *Tpl) output() error {
	outDir := "./gen/"
	if !exists(outDir) {
		err := os.MkdirAll(outDir, 0755)
		if err != nil {
			return err
		}
	}

	outpath := outDir + tpl.name + ".go"
	err := ioutil.WriteFile(outpath, []byte(tpl.result), 0644)
	if err != nil {
		return err
	}
	return nil
}
