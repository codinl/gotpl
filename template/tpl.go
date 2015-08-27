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
	raw 	 []byte
	result 	 string
	tokens   []*Token
	ast      *Ast
	blocks   map[string]*Ast
	sections map[string]*Ast
	option Option
}

func InitTpl(name string, option Option) (*Tpl, error) {
	tpl := &Tpl{name:name,option:option}

	err := tpl.content()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tpl, nil
}

func (tpl Tpl) Generate(option Option) error {
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

	err = tpl.fmt()
	if err != nil {
		return err
	}

	err = tpl.output()
	if err != nil {
		return err
	}

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

func (tpl Tpl) compile() error {
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

	tpl.result = cp.buf

	return nil
}

func (tpl Tpl) content() error {
	content, err := ioutil.ReadFile(tpl.path())
	if err != nil {
		return err
	}
	tpl.content = content
	return nil
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

func (tpl Tpl) output() error {
	err := ioutil.WriteFile(tpl.outpath(), []byte(tpl.result), 0644)
	if err != nil {
		return  err
	}
	return nil
}
