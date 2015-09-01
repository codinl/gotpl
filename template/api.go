package gotpl

import (
	"fmt"
	"path/filepath"
	"strings"
	"os"
)

func Generate(input string, output string, option Option) error {
	sections, err := genSection(input)
	if err != nil {
		return err
	}

	tplMap := map[string]*Tpl{}

	paths, err := getFiles(input, TPL_EXT)
	if err != nil {
		return err
	}

	for i := 0; i < len(paths); i++ {
		path := paths[i]

		baseName := filepath.Base(path)
		name := strings.TrimSpace(strings.Replace(baseName, TPL_EXT, "", 1))

		tpl := &Tpl{
			path: path, name: name,
			ast: &Ast{}, tokens: []Token{},
			blocks: map[string]*Ast{}, outDir:output,
			option: option,
		}

		err := tpl.readRaw()
		if err != nil {
			fmt.Println(err)
			return err
		}

		tpl.checkSection(sections)

		err = tpl.checkExtends()
		if err != nil {
			fmt.Println(err)
			return err
		}

		tplMap[tpl.name] = tpl
	}

	for key, tpl := range tplMap {
		if !tpl.isRoot {
			if p, ok := tplMap[tpl.parentName]; ok {
				tplMap[key].parent = p
			} else {
				fmt.Println(tpl.parentName, "--parent not exists")
				delete(tplMap, key)
			}
		}
	}

	err = os.RemoveAll(output)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, tpl := range tplMap {
		err = tpl.generate()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	err = fmtCode(output)
	if err != nil {
		return err
	}

	return nil
}

