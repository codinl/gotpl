package template

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var Namespace = `"github.com/codinl/gotpl/gotpl"`

const (
	PART_MKP = iota
	PART_TOKEN
	PART_STAT
)

type Part struct {
	Type  int
	value string
}

func makeCompiler(ast *Ast, options Option, input string) *Compiler {
	dir := filepath.Base(filepath.Dir(input))
	file := strings.Replace(filepath.Base(input), TPL_EXT, "", 1)
	if options["NameNotChange"] == nil {
		file = Capitalize(file)
	}
	return &Compiler{
		ast: ast, buf: "",
		params: []string{}, parts: []Part{},
		imports:  map[string]bool{},
		options:  options,
		dir:      dir,
		fileName: file,
	}
}

//------------------------------ Compiler ------------------------------ //
type Compiler struct {
	tpl       *Tpl
	ast       *Ast
	buf       string //the final result
	firstNode bool
	parts     []Part
	params    []string
	imports   map[string]bool
	options   Option
	dir       string
	fileName  string
}

func (cp *Compiler) visitAst(ast *Ast) {
	switch ast.Mode {
	case MKP:
		cp.firstNode = true
		for _, c := range ast.Children {
			if _, ok := c.(Token); ok {
				cp.visitMKP(c, ast)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	case NODE:
		if cp.firstNode == false {
			cp.firstNode = true
			cp.visitFirstNode(ast)
		} else {
			remove := false
			if len(ast.Children) >= 2 {
				first := ast.Children[0]
				last := ast.Children[len(ast.Children)-1]
				v1, ok1 := first.(Token)
				v2, ok2 := last.(Token)
				if ok1 && ok2 && v1.Text == "{" && v2.Text == "}" {
					remove = true
				}
			}
			for idx, c := range ast.Children {
				if remove && (idx == 0 || idx == len(ast.Children)-1) {
					continue
				}
				if _, ok := c.(Token); ok {
					cp.visitToken(c, ast)
				} else {
					cp.visitAst(c.(*Ast))
				}
			}
		}
	case BLOCK:
		firstBraceFind := false
		for idx, c := range ast.Children {
			if !firstBraceFind {
				if t, ok := c.(Token); ok {
					if t.Type == BRACE_OPEN {
						firstBraceFind = true
					}
				}
				continue
			}

			if idx == len(ast.Children)-1 {
				continue
			}

			if _, ok := c.(Token); ok {
				cp.visitMKP(c, ast)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	case EXP:
		cp.firstNode = true
		nonExp := ast.hasNonExp()
		for i, c := range ast.Children {
			if _, ok := c.(Token); ok {
				cp.visitExp(c, ast, i, !nonExp)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	case PRG:
		for _, c := range ast.Children {
			cp.visitAst(c.(*Ast))
		}
	}
}

func (cp *Compiler) visitMKP(child interface{}, ast *Ast) {
	cp.addPart(Part{PART_MKP, getValStr(child)})
}

func (cp *Compiler) visitToken(child interface{}, ast *Ast) {
	cp.addPart(Part{PART_TOKEN, getValStr(child)})
}

func (cp *Compiler) visitExp(child interface{}, parent *Ast, idx int, isHomo bool) {
	start := ""
	end := ""
	ppNotExp := true
	ppChildCnt := len(parent.Children)
	//	pack := cp.dir
	htmlEsc := cp.options["htmlEscape"]
	if parent.Parent != nil && parent.Parent.Mode == EXP {
		ppNotExp = false
	}
	val := getValStr(child)
	if htmlEsc == nil {
		if ppNotExp && idx == 0 && isHomo {
			start += "gotpl.HTMLEscape("
			cp.imports[Namespace] = true
		}
		if ppNotExp && idx == ppChildCnt-1 && isHomo {
			end += ")"
		}
	}

	if ppNotExp && idx == 0 {
		start = "_buffer.WriteString(" + start
	}
	if ppNotExp && idx == ppChildCnt-1 {
		end += ")\n"
	}

	v := start
	if val == "raw" {
		v += end
	} else {
		v += val + end
	}

	cp.addPart(Part{PART_STAT, v})
}

// First node contains imports and parameters, specific action for layout,
// NOTE, layout have some conventions.
func (cp *Compiler) visitFirstNode(node *Ast) {
	pre := cp.buf
	cp.buf = ""
	first := ""

	backup := cp.parts
	cp.parts = []Part{}

	cp.visitAst(node)

	cp.genPart()

	first, cp.buf = cp.buf, pre

	cp.parts = backup

	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", "package main\n"+first, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		for _, s := range f.Imports {
			v := s.Path.Value
			if s.Name != nil {
				v = s.Name.Name + " " + v
			}
			cp.imports[v] = true
		}
	}

	lines := strings.SplitN(first, "\n", -1)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "var") {
			cp.params = append(cp.params, l[4:])
		}
	}

	// extends parent
	if !cp.tpl.isRoot && cp.tpl.parent != nil {
		for k, v := range cp.tpl.parent.imports {
			if _, ok := cp.imports[k]; !ok {
				cp.imports[k] = v
			}
		}

		params := cp.tpl.parent.params
		for _, v := range cp.params {
			exists := false
			for _, vv := range cp.tpl.parent.params {
				if vv == v {
					exists = true
					break
				}
			}
			if !exists {
				params = append(params, v)
			}
		}
		cp.params = params
	}

	cp.tpl.imports = cp.imports
	cp.tpl.params = cp.params
}

func (cp *Compiler) addPart(part Part) {
	if len(cp.parts) == 0 {
		cp.parts = append(cp.parts, part)
		return
	}
	last := &cp.parts[len(cp.parts)-1]
	if last.Type == part.Type {
		last.value += part.value
	} else {
		cp.parts = append(cp.parts, part)
	}
}

func (cp *Compiler) genPart() {
	res := ""
	for _, p := range cp.parts {
		if p.Type == PART_MKP && p.value != "" {
			// do some escapings
			for strings.HasSuffix(p.value, "\n") {
				p.value = p.value[:len(p.value)-1]
			}
			if p.value != "" {
				p.value = fmt.Sprintf("%#v", p.value)
				res += "_buffer.WriteString(" + p.value + ")\n"
			}
		} else if p.Type == PART_TOKEN {
			res += p.value + "\n"
		} else {
			res += p.value
		}
	}
	cp.buf = res
}

//func (cp *Compiler) processBlock() {
//	lines := strings.SplitN(cp.buf, "\n", -1)
//	out := ""
//	scope := 0
//	for _, l := range lines {
//		l = strings.TrimSpace(l)
//		if strings.HasPrefix(l, "block") && strings.HasSuffix(l, "{") {
//			scope = 1
//			out += "\n"
//		} else if scope > 0 {
//			if strings.HasSuffix(l, "{") {
//				scope++
//			} else if strings.HasSuffix(l, "}") {
//				scope--
//			}
//			if scope == 0 {
//				scope = 0
//			} else {
//				out += l + "\n"
//			}
//		} else {
//			out += l + "\n"
//		}
//	}
//
//	cp.buf = out
//}

//func (cp *Compiler) processSection() {
//	cp.imports[Namespace] = false
//	cp.imports[`"bytes"`] = false
//	head := "package tpl\n\n import (\n"
//	for k, v := range cp.imports {
//		if v {
//			head += k + "\n"
//		}
//	}
//	head += ")\n"
//
//	lines := strings.SplitN(cp.buf, "\n", -1)
//	out := ""
//	scope := 0
//	secOut := ""
//	for _, l := range lines {
//		l = strings.TrimSpace(l)
//		if strings.HasPrefix(l, "section") && strings.HasSuffix(l, "{") {
//			name := l
//			name = strings.TrimSpace(name[7 : len(name)-1])
//			out += "\n func " + Capitalize(name) + " string {\n"
//			secOut = "return `\n"
//			scope = 1
//		} else if scope > 0 {
//			if strings.HasSuffix(l, "{") {
//				scope++
//			} else if strings.HasSuffix(l, "}") {
//				scope--
//			}
//			if scope == 0 {
//				secOut += "`\n}\n"
//				out += secOut
//				secOut = ""
//			} else {
//				secOut += l + "\n"
//			}
//		} else {
//			secOut += l + "\n"
//		}
//	}
//
//	cp.buf = head + out
//}

func (cp *Compiler) visit() {
	// visitAst(cp.ast) -> cp.parts
	cp.visitAst(cp.ast)

	// genPart() -> cp.buf
	cp.genPart()

	funcName := Capitalize(cp.fileName)

	cp.imports[`"bytes"`] = true
	head := "package tpl\n\n import (\n"
	for k, v := range cp.imports {
		if v {
			head += k + "\n"
		}
	}

	head += "\n)\n func " + funcName + "("
	for idx, p := range cp.params {
		head += p
		if idx != len(cp.params)-1 {
			head += ", "
		}
	}
	head += ") string {\n var _buffer bytes.Buffer\n"

	cp.buf = head + cp.buf

	//	cp.processBlock()

	cp.buf += "return _buffer.String()\n}"
}

func getValStr(e interface{}) string {
	switch v := e.(type) {
	case *Ast:
		return v.TagName
	case Token:
		if !(v.Type == AT || v.Type == AT_COLON) {
			return v.Text
		}
		return ""
	default:
		panic(e)
	}
}

//func watchDir(input, output string, options Option) error {
//	log.Println("Watching dir:", input, output)
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer watcher.Close()
//
//	done := make(chan bool)
//
//	output_path := func(path string) string {
//		res := strings.Replace(path, input, output, 1)
//		return res
//	}
//
//	gen := func(filename string) error {
//		outpath := output_path(filename)
//		outpath = strings.Replace(outpath, ".gohtml", ".go", 1)
//		outdir := filepath.Dir(outpath)
//		if !exists(outdir) {
//			os.MkdirAll(outdir, 0775)
//		}
//		err := GenFile(filename, outpath, options)
//		if err == nil {
//			log.Printf("%s -> %s\n", filename, outpath)
//		}
//		return err
//	}
//
//	visit_gen := func(path string, info os.FileInfo, err error) error {
//		if !info.IsDir() {
//			//Just do file with exstension .gohtml
//			if !strings.HasSuffix(path, ".gohtml") {
//				return nil
//			}
//			filename := filepath.Base(path)
//			if strings.HasPrefix(filename, ".#") {
//				return nil
//			}
//			err := gen(path)
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//
//	go func() {
//		for {
//			select {
//			case event := <-watcher.Events:
//				filename := event.Name
//				if filename == "" {
//					//should be a bug for fsnotify
//					continue
//				}
//				if event.Op&fsnotify.Remove != fsnotify.Remove &&
//					(event.Op&fsnotify.Write == fsnotify.Write ||
//						event.Op&fsnotify.Create == fsnotify.Create) {
//					stat, err := os.Stat(filename)
//					if err != nil {
//						continue
//					}
//					if stat.IsDir() {
//						log.Println("add dir:", filename)
//						watcher.Add(filename)
//						output := output_path(filename)
//						log.Println("mkdir:", output)
//						if !exists(output) {
//							os.MkdirAll(output, 0755)
//							err = filepath.Walk(filename, visit_gen)
//							if err != nil {
//								done <- true
//							}
//						}
//						continue
//					}
//					if !strings.HasPrefix(filepath.Base(filename), ".#") &&
//						strings.HasSuffix(filename, ".gohtml") {
//						gen(filename)
//					}
//				} else if event.Op&fsnotify.Remove == fsnotify.Remove ||
//					event.Op&fsnotify.Rename == fsnotify.Rename {
//					output := output_path(filename)
//					if exists(output) {
//						//shoud be dir
//						watcher.Remove(filename)
//						os.RemoveAll(output)
//						log.Println("remove dir:", output)
//					} else if strings.HasSuffix(output, ".gohtml") {
//						output = strings.Replace(output, ".gohtml", ".go", 1)
//						if exists(output) {
//							os.Remove(output)
//							log.Println("removing file:", output)
//						}
//					}
//				}
//			case err := <-watcher.Errors:
//				log.Println("error:", err)
//				continue
//			}
//		}
//	}()
//
//	visit := func(path string, info os.FileInfo, err error) error {
//		if info.IsDir() {
//			watcher.Add(path)
//		}
//		return nil
//	}
//
//	err = filepath.Walk(input, visit)
//	err = watcher.Add(input)
//	if err != nil {
//		log.Fatal(err)
//	}
//	<-done
//	return nil
//}
