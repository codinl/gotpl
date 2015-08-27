package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type Tpl struct {
	name     string
	parent   *Tpl
	content  []byte
	tokens   []*Token
	ast      *Ast
	blocks   map[string]*Ast
	sections map[string]*Ast
}

func (tpl Tpl) Generate() error {
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
		//		file:    file,
	}

	// visit() -> cp.buf
	cp.visit()

	err := ioutil.WriteFile(tpl.outpath(), []byte(cp.buf), 0644)
	if err != nil {
		panic(err)
	}

	err = tpl.fmt()

	//	if option["Debug"] == true {
	//		content, _ := ioutil.ReadFile(output)
	//		fmt.Println(string(content))
	//	}

	return nil
}

func (tpl Tpl) genToken() error {
	lex := &Lexer{Text: tpl.content, Matches: TokenMatches}

	tokens, err := lex.Scan()
	if err != nil {
		return err
	}

	tpl.tokens = tokens

	return nil
}

func (tpl Tpl) genAst() error {
	parser := &Parser{
		ast: tpl.ast, tokens: tpl.tokens,
		preTokens: []Token{}, initMode: UNK,
		blocks: tpl.blocks,
	}

	// Run() -> ast
	err := parser.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (tpl Tpl) content() {
	content, err := ioutil.ReadFile(tpl.path())
	if err == nil {
		tpl.content = content
	}
}

func (tpl Tpl) path() string {
	return "" + tpl.name
}

func (tpl Tpl) outpath() string {
	return "" + tpl.name
}

func (tpl Tpl) name() string {
	return "" + tpl.name
}

func (tpl Tpl) fmt() error {
	cmd := exec.Command("gofmt", "-w", tpl.outpath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("gofmt: ", err)
		return err
	}

	return nil
}
