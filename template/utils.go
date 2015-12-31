package gotpl

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func escapeHtml(m interface{}) string {
	s := fmt.Sprint(m)
	return template.HTMLEscapeString(s)
}

func timeToStr(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
}

func itoa(obj int) string {
	return strconv.Itoa(obj)
}

func capitalize(str string) string {
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

func getFiles(dir string, fileExt string) ([]string, error) {
	fileNames := []string{}
	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if !fi.IsDir() {
			if strings.HasSuffix(fi.Name(), fileExt) {
				abs, err := filepath.Abs(dir)
				//				fmt.Println(abs)
				if err == nil {
					fileNames = append(fileNames, abs+"/"+fi.Name())
				}
			}
		}
	}
	return fileNames, nil
}

func rec(reg string) *regexp.Regexp {
	return regexp.MustCompile("^" + reg)
}
