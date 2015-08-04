package gotpl

import (
	"fmt"
	"strings"
	"io/ioutil"
	"errors"
)

type Section struct {
	StartPos int
	EndPos   int
	Content string
	Name     string
}

func genSections(tplDir string) (map[string]*Section, error) {
	dir := tplDir + SEC_DIR
	if !exists(dir) {
		fmt.Println("dir", tplDir + SEC_DIR)
		return nil, errors.New("dir not exists")
	}

	files, err := getFiles(dir)
	if err != nil {
		fmt.Println("err == ", err)
		return nil, err
	}

	sections := map[string]*Section{}

	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			continue
		}

		sections, err = findSections(sections, string(content))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return sections, nil
}

func findSections(sections map[string]*Section, text string) (map[string]*Section, error) {
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '@':
			start_pos := i
			if text[i+1:i+8] == "section" { // @section
				i += 8
				i_name_start := i
				i_content_start := 0
				i_name_end := 0
				scope_1 := 0
				scope_2 := -1
				section := &Section{}
				name := ""
				for j:=i; j<len(text); j++ {
					i++
					switch text[j] {
					case '(':
						scope_1++
						if scope_1 == 1 && name == "" { // @section Name(
							i_name_end = j
							name = text[i_name_start:i_name_end]
							section.Name = strings.ToLower(strings.TrimSpace(name))
							section.StartPos = start_pos
						}
					case ')':// @section Name()
						scope_1--
					case '{':
						if scope_2 < 0 {
							scope_2 = 0
						}
						scope_2++
						if scope_2 == 1 { // @section Name() {
							i_content_start = (i+1)
						}
					case '}':  // @section Name(){ ... }
						scope_2--
					}
					if scope_2 == 0 {
						section.EndPos = j
						section.Content = text[i_content_start:j]
						sections[section.Name] = section
						break
					}
				}
			}
		}
	}

	return sections, nil
}
