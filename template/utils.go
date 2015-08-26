package gotpl

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func HTMLEscape(m interface{}) string {
	s := fmt.Sprint(m)
	return template.HTMLEscapeString(s)
}

func StrTime(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
}

func Itoa(obj int) string {
	return strconv.Itoa(obj)
}

func Capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func getFiles(dir string) ([]string, error) {
	fileNames := []string{}

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if !strings.HasSuffix(path, TPL_EXT) {
				return nil
			}
			base := filepath.Base(path)
			if strings.HasPrefix(base, "__") {
				return nil
			}
			fullpath, err := filepath.Abs(path)
			if err != nil {
				return nil
			}
			fileNames = append(fileNames, fullpath)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return nil, err
	}

	return fileNames, nil
}
