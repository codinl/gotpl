package gotpl

import (
	"github.com/codinl/go-logger"
	"os"
	"path/filepath"
	"strings"
	"errors"
	"strconv"
	"time"
	"html/template"
	"fmt"
)

func Generate(input string, output string, option Option) error {
	err := initLogger()
	if err != nil {
		panic("initLogger fail")
	}

	sections, err := genSection(input)
	if err != nil {
		logger.Error(err)
		return err
	}

	tplMap := map[string]*Tpl{}

	paths, err := getFiles(input, TPL_EXT)
	if err != nil {
		logger.Error(err)
		return err
	}

	for i := 0; i < len(paths); i++ {
		path := paths[i]

		baseName := filepath.Base(path)
		name := strings.TrimSpace(strings.Replace(baseName, TPL_EXT, "", 1))

		tpl := &Tpl{
			path: path, name: name,
			ast: &Ast{}, tokens: []Token{},
			blocks: map[string]*Ast{}, outDir: output,
			option: option,
		}

		err := tpl.readRaw()
		if err != nil {
			logger.Error(err)
			return err
		}

		tpl.checkSection(sections)

		err = tpl.checkExtends()
		if err != nil {
			logger.Error(err)
			return err
		}

		tplMap[tpl.name] = tpl
	}

	if len(tplMap) == 0 {
		return errors.New("tpl is empty")
	}

	for key, tpl := range tplMap {
		if !tpl.isRoot {
			if p, ok := tplMap[tpl.parentName]; ok {
				tplMap[key].parent = p
			} else {
				logger.Info(tpl.parentName, "--parent not exists")
				delete(tplMap, key)
			}
		}
	}

	// clean output direct
	err = os.RemoveAll(output)
	if err != nil {
		logger.Error(err)
	}

	for _, tpl := range tplMap {
		err = tpl.generate()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func HTMLEscape(m interface{}) string {
	s := fmt.Sprint(m)
	return template.HTMLEscapeString(s)
}

func TimeToStr(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
}

func Itoa(obj int) string {
	return strconv.Itoa(obj)
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
