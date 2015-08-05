package gotpl

type BlockNode struct {
	Name     string
	StartPos int
	EndPos   int
	Content  string
	Parent   *BlockNode
	Children map[string]*BlockNode
}

type BlockToken struct {
	Name     string
	Scope    int
	StartPos int
	EndPos   int
}

const (
	NODE_TYPE_UNK = iota
	NODE_TYPE_BLOCK
)

type Node struct {
	Type     int
	StartPos int
	EndPos   int
	Content  string
	Name     string
	Children []*Node
	Parent   *Node
}
