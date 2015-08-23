package gotpl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Tpl struct {
	isLeaf     bool
	Parent     *Tpl
	Children   []*Tpl
	Name       string
	Content    string
	Block      *BlockNode
	Blocks     []*BlockNode
	BlockBuild bool // 是否已经生成过block
	FileBuild  bool // 是否已经生成过file
}

var tplDir string

func (tpl *Tpl) buildBlock() error {
	text := tpl.Content

	tpl.Block = &BlockNode{Name: "root", Children: map[string]*BlockNode{}}

	var inNode *BlockNode
	var parentNode *BlockNode

	tokens := map[string]*BlockToken{}

	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '@':
			start_pos := i
			if text[i+1:i+6] == "block" { // @block
				i += 6
				i_name_start := i
				i_name_end := 0
				scope := 0
				L:
				for j := i; j < len(text); j++ {
					i++
					switch text[j] {
					case '{':
						scope++
						if scope == 1 { // @block name {
							i_name_end = j
							name := text[i_name_start:i_name_end]
							if parentNode == nil {
								parentNode = tpl.Block
							}
							blockName := strings.TrimSpace(name)
							inNode = &BlockNode{StartPos: start_pos, Name: blockName, Parent: parentNode, Children: map[string]*BlockNode{}}
							parentNode.Children[blockName] = inNode
							tokens[blockName] = &BlockToken{Name: blockName, StartPos: i_name_end+1, EndPos:0}
							tpl.Blocks = append(tpl.Blocks, inNode)
						}
						tokens[inNode.Name].Scope++
					case '}':
						scope--
						if inNode != nil {
							if _, ok := tokens[inNode.Name]; ok {
								tokens[inNode.Name].Scope--
//								if tokens[inNode.Name].Scope == 0 {
//									if parentNode != nil {
//										if _, ok := tokens[parentNode.Name]; ok {
//											tokens[parentNode.Name].Scope--
//										}
//									}
//								} else {
//									tokens[inNode.Name].Scope--
//								}
							}
						}
					case '@':
						if text[j+1:j+6] == "block" {
							if _, ok := tokens[inNode.Name]; ok {
								if tokens[inNode.Name].Scope == 1 { // 上一级还没有关闭
									parentNode = inNode
								}
							}
							i = j - 1
							break L
						}
					default:

					}

					if inNode != nil {
						if _, ok := tokens[inNode.Name]; ok {
							if tokens[inNode.Name].Scope == 0 && tokens[inNode.Name].EndPos == 0 {
								tokens[inNode.Name].EndPos = j
								inNode.EndPos = j
								inNode.Content = text[tokens[inNode.Name].StartPos:j]
								inNode = inNode.Parent
							}
						}
					}
				}
			}
		}
	}

	tpl.BlockBuild = true

	return nil
}

func (tpl *Tpl) genBlock() error {
	if tpl == nil {
		return errors.New("tpl == nil")
	}

	if !tpl.isRoot() && tpl.BlockBuild == false {
		tpl.buildBlock()
	}

	if !tpl.isRoot() && !tpl.FileBuild {
		tpl.buildFile()
	}

	if tpl.isLeaf {
		if len(tpl.Parent.Children) == 1 {
			tpl.Parent.Children[0] = nil
			if !tpl.Parent.isRoot() {
				tpl.Parent.Parent.Children[0].genBlock()
			}
		} else if len(tpl.Parent.Children) > 1 {
			tpl.Parent.Children = tpl.Parent.Children[1:]
			tpl.Parent.Children[0].genBlock()
		}
	} else {
		if len(tpl.Children) > 0 { // has children
			tpl.Children[0].genBlock()
		} else { // no children
			if !tpl.isRoot() {
				if len(tpl.Parent.Children) == 1 {
					tpl.Parent.Children[0] = nil
				} else {
					tpl.Parent.Children = tpl.Parent.Children[1:]
				}
				tpl.Parent.Children[0].genBlock()
			}
		}
	}

	return nil
}

func (tpl *Tpl) buildFile() {
	if !tpl.isRoot() && !tpl.Parent.isRoot() && len(tpl.Parent.Blocks) > 0 {
		buf := ""
		replaceCount := 0
		var lastBlock *BlockNode

		for i := 0; i < len(tpl.Blocks); i++ {
			b := tpl.Blocks[i]
			targetBlock := tpl.Parent.getBlockByName(b.Name)
			if targetBlock != nil {
				if i == 0 {
					buf += tpl.Parent.Content[:targetBlock.StartPos] + b.Content
				} else {
//					fmt.Println("-------------",string(tpl.Parent.Content[lastBlock.EndPos+1:]))
					buf += tpl.Parent.Content[lastBlock.EndPos+1:targetBlock.StartPos] + b.Content
				}
				if i == len(tpl.Blocks)-1 {
					buf += tpl.Parent.Content[targetBlock.EndPos+1:]
				}
				replaceCount++
				lastBlock = targetBlock
			}
		}

		if replaceCount > 0 {
			tpl.Content = buf
		}
	}

	text := clearBlocks(tpl.Content)

	text = replaceSections(text)

	tpl.writeFile(text)
}

func (tpl *Tpl) writeFile(text string) error {
	outpath := tplDir + TMP_DIR + tpl.Name + TMP_EXT

	outdir := filepath.Dir(tplDir + TMP_DIR)
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}

	err := ioutil.WriteFile(outpath, []byte(text), 0644)
	if err != nil {
		return err
	}

	tpl.FileBuild = true

	return nil
}

func (tpl *Tpl) getBlockByName(name string) *BlockNode {
	for _, b := range tpl.Blocks {
		if b.Name == name {
			return b
		}
	}
	return nil
}

func buildTplTree(dir string) (*Tpl, error) {
	fileNames, err := getFiles(dir)
	if err != nil {
		return nil, err
	}

	tplMap := map[string]*Tpl{}

	rootTpl := &Tpl{Name: "root", isLeaf: false}
	tplMap["root"] = rootTpl

	for _, f := range fileNames {
		baseName := filepath.Base(f)
		if !strings.HasSuffix(baseName, TPL_EXT) {
			continue
		}
		filename := strings.Replace(baseName, TPL_EXT, "", -1)
		parentFileName := filename
		ss := strings.Split(filename, "_")
		if len(ss) == 3 && ss[1] == "extends" {
			filename = ss[0]
			parentFileName = strings.Replace(ss[2], TPL_EXT, "", -1)
		} else {
			parentFileName = "root"
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		tpl := &Tpl{Name: filename, isLeaf: true, Content: string(content), FileBuild: false}

		if parent, ok := tplMap[parentFileName]; ok {
			tpl.Parent = parent
			tpl.Parent.isLeaf = false
			parent.Children = append(parent.Children, tpl)
		}

		tplMap[filename] = tpl
	}

	return rootTpl, nil
}

func clearBlocks(text string) string {
	endPos := 0

	startPos := strings.Index(text, "@block")
	if startPos < 0 {
		return text
	}

	buf := text[:startPos]

	text_1 := text[startPos:]

	startPos_1 := strings.Index(text_1, "{")

	if startPos_1 < 2 {
		if startPos_1 < 0 {
			fmt.Println("no {\n")
		} else {
			fmt.Println("no block name\n")
		}
	}

	scope := 1
	for i := startPos_1 + 1; i < len(text_1); i++ {
		switch text_1[i] {
		case '{':
			scope++
		case '}':
			scope--
		}
		if scope == 0 {
			endPos = i
			break
		}
	}

	buf += text_1[startPos_1+1:endPos] + text_1[endPos+1:]

	return clearBlocks(buf)
}

func replaceSections(text string) string {
	endPos := 0

	startPos := strings.Index(text, "@section")
	if startPos < 0 {
		return text
	}

	buf := text[:startPos]

	text_1 := text[startPos:]

	startPos_1 := strings.Index(text_1, "(")

	if startPos_1 < 2 {
		if startPos_1 < 0 {
			fmt.Println("no (\n")
		} else {
			fmt.Println("no section name\n")
		}
	}

	secName := ""
	scope := -1
	for i := startPos_1; i < len(text_1); i++ {
		switch text_1[i] {
		case '(':
			if scope < 0 {
				scope = 0
			}
			scope++
			if scope == 1 { // @section Name(
				secName = strings.TrimSpace(text_1[8:startPos_1])
			}
		case ')':
			scope--
		}
		if scope == 0 {
			endPos = i
			break
		}
	}

	if secName != "" && Sections != nil {
		section, exists := Sections[strings.ToLower(secName)]
		if exists {
			buf += text_1[startPos_1+1:endPos] + section.Content + text_1[endPos+1:]
			return replaceSections(buf)
		}
	}

	buf += text_1[startPos_1+1:endPos] + text_1[endPos+1:]

	return replaceSections(buf)
}

func (tpl *Tpl) isRoot() bool {
	return tpl.Name == "root"
}
