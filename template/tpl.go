package gotpl

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"os/exec"
"github.com/codinl/go-logger"
)

type Tpl struct {
	path       string
	name       string
	isRoot     bool
	parent     *Tpl
	parentName string
	raw        []byte
	result     string
	tokens     []Token
	ast        *Ast
	blocks     map[string]*Ast
	sections   map[string]*Ast
	params     []string
	imports    map[string]bool
	generated  bool
	option     Option
	outDir    string
}

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

	//	logger.Info(tpl.name, "------------------- TOKEN START -----------------")
	//	for _, elem := range tpl.tokens {
	//		elem.debug()
	//	}
	//	logger.Info(tpl.name, "--------------------- TOKEN END -----------------\n")

	err = tpl.genAst()
	if err != nil {
		return err
	}

	//	logger.Info(tpl.name, "--------------------- AST START -----------------")
	//	tpl.ast.debug(0, 20)
	//	logger.Info(tpl.name, "--------------------- AST END -----------------\n")
	//	if tpl.ast.Mode != PRG {
	//		panic("TYPE")
	//	}

	err = tpl.genResult()
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
		ast: tpl.ast, tokens: tpl.tokens, blocks: tpl.blocks,
		preTokens: []Token{}, initMode: UNK,
	}

	// Run() -> ast
	err := parser.Run()
	if err != nil {
		logger.Info(err)
		return err
	}

	if !tpl.isRoot && tpl.parent != nil {
		firstNode := tpl.ast.Children[0]
		tpl.ast = tpl.parent.ast
		tpl.ast.Children[0] = firstNode

		if tpl.blocks != nil && len(tpl.blocks) > 0 {
			updateAst(tpl.ast, tpl.blocks)
		}
	}

	return nil
}

func (tpl *Tpl) genResult() error {
	cp := &Compiler{
		tpl: tpl,
		ast: tpl.ast, buf: "",
		params:  []string{},
		imports: map[string]bool{},
		parts:   []Part{},
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
		logger.Info(err)
		return err
	}

	tpl.raw = raw

	tpl.fmtRaw()

	return nil
}

func (tpl *Tpl) checkSection(sections map[string]*Section) error {
	idx := bytes.Index(tpl.raw, []byte("@section"))
	if idx > 0 {
		left := tpl.raw[idx+len("@section"):]
		idx_1 := bytes.Index(left, []byte("("))
		idx_2 := bytes.Index(left, []byte(")"))
		if idx_1 > 0 && idx_2 > 0 {
			name := string(bytes.TrimSpace(left[:idx_1]))
			tpl.raw = tpl.raw[:idx]
			left_1 := left[idx_2+1:]
			if sections != nil {
				if section, ok := sections[name]; ok {
					tem := append(section.text, '\n')
					left_1 = append(tem, left_1...)
				}
			}
			tpl.raw = append(tpl.raw, left_1...)
		} else {
			tpl.raw = bytes.Replace(tpl.raw, []byte("@section"), []byte(""), 1)
		}
		return tpl.checkSection(sections)
	}
	return nil
}

func (tpl *Tpl) checkExtends() error {
	tpl.parentName = ""
	tpl.isRoot = true

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
		if len(ss) == 2 && lineEnd > 0 {
			parentName := strings.TrimSpace(ss[1])
			tpl.parentName = parentName
			tpl.isRoot = false

			tpl.raw = tpl.raw[lineEnd+1:]

			tpl.fmtRaw()
		}
	}

	return nil
}

func genSection(input string) (map[string]*Section, error) {
	dir := input + SEC_DIR
	if !exists(dir) {
		return nil, nil
	}

	paths, err := getFiles(dir, TPL_EXT)
	if err != nil {
		return nil, err
	}

	sections := map[string]*Section{}

	for i := 0; i < len(paths); i++ {
		path := paths[i]

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Info(err)
			return nil, err
		}

		ss := bytes.Split(raw, []byte("@section"))
		for _, s := range ss {
			idx := bytes.Index(s, []byte("("))
			idx_2 := bytes.Index(s, []byte("{"))
			idxEnd := bytes.LastIndex(s, []byte("}"))
			if idx > 0 && idx_2 > 0 && idxEnd > 0 {
				name := string(bytes.TrimSpace(s[:idx]))
				text := bytes.TrimSpace(s[idx_2+1 : idxEnd])
				section := &Section{name: name, text: text}
				sections[name] = section
			}
		}
	}

	return sections, nil
}

func fmtCode(output string) error {
	cmd := exec.Command("gofmt", "-w", output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Info("gofmt: ", err)
		return err
	}
	return nil
}

func (tpl *Tpl) output() error {
	if !exists(tpl.outDir) {
		err := os.MkdirAll(tpl.outDir, 0755)
		if err != nil {
			return err
		}
	}

	outpath := tpl.outDir + tpl.name + ".go"
	err := ioutil.WriteFile(outpath, []byte(tpl.result), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (tpl *Tpl) fmtRaw() {
	if tpl.raw != nil {
		tpl.raw = bytes.TrimSpace(tpl.raw)
		if !bytes.HasPrefix(tpl.raw, []byte("@{")) && !bytes.HasPrefix(tpl.raw, []byte("@extends")) {
			tpl.raw = append([]byte("@{\n}\n"), tpl.raw...)
		}
	}
}

func updateAst(ast *Ast, blocks map[string]*Ast) {
	if ast.Children == nil || len(ast.Children) == 0 || blocks == nil || len(blocks) == 0 {
		return
	}
	for idx, c := range ast.Children {
		if a, ok := c.(*Ast); ok {
			if a.Mode == BLOCK {
				if b, ok := blocks[a.TagName]; ok {
					ast.Children[idx] = b
					delete(blocks, ast.TagName)
				}
			}
			updateAst(a, blocks)
		}
	}
}
