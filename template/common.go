package gotpl

var Namespace = `"github.com/codinl/gotpl/template"`

const (
	GO_EXT  = ".go"
	TPL_EXT = ".html"
	SEC_DIR = "sections/"
)

// Option have following options:
//   Debug bool
//   Watch bool
//   NameNotChange bool
type Option map[string]interface{}
