package gotpl

type Section struct {
	name    string
	text    []byte
	imports map[string]bool
	params  []string
}
