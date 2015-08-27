package template

const (
	GO_EXT  = ".go"
	TPL_EXT = ".tpl"
	TMP_EXT = ".t"
	TMP_DIR = "tmp/"
	GEN_DIR = "testgen/"
	SEC_DIR = "sections/"
)

// Option have following options:
//   Debug bool
//   Watch bool
//   NameNotChange bool
type Option map[string]interface{}
