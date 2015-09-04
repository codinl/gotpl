package gotpl

import (
	"fmt"
	"regexp"
	"github.com/codinl/go-logger"
)

//------------------------------ Ast ------------------------------//
type Ast struct {
	Parent   *Ast
	Children []interface{}
	Mode     int
	TagName  string
}

func (ast *Ast) ModeStr() string {
	switch ast.Mode {
	case PRG:
		return "PROGRAM"
	case MKP:
		return "MARKUP"
	case NODE:
		return "NODE"
	case EXP:
		return "EXPRESSION"
	case BLOCK:
		return "BLOCK"
	default:
		return "UNDEF"
	}
	return "UNDEF"
}

func (ast *Ast) check() {
	if len(ast.Children) >= 100000 {
		panic("Maximum number of elements exceeded.")
	}
}

func (ast *Ast) addChild(child interface{}) {
	ast.Children = append(ast.Children, child)
	ast.check()
	if _a, ok := child.(*Ast); ok {
		_a.Parent = ast
	}
}

func (ast *Ast) addChildren(children []Token) {
	for _, c := range children {
		ast.addChild(c)
	}
}

func (ast *Ast) addAst(_ast *Ast) {
	c := _ast
	for {
		if len(c.Children) != 1 {
			break
		}
		first := c.Children[0]
		if _, ok := first.(*Ast); !ok {
			break
		}
		c = first.(*Ast)
	}
	if c.Mode != PRG {
		ast.addChild(c)
	} else {
		for _, x := range c.Children {
			ast.addChild(x)
		}
	}
}

func (ast *Ast) popChild() {
	l := len(ast.Children)
	if l == 0 {
		return
	}
	ast.Children = ast.Children[:l-1]
}

func (ast *Ast) root() *Ast {
	p := ast
	pp := ast.Parent
	for {
		if p == pp || pp == nil {
			return p
		}
		b := pp
		pp = p.Parent
		p = b
	}
	return nil
}

func (ast *Ast) beget(mode int, tag string) *Ast {
	child := &Ast{nil, []interface{}{}, mode, tag}
	ast.addChild(child)
	return child
}

func (ast *Ast) closest(mode int, tag string) *Ast {
	p := ast
	for {
		if p.TagName != tag && p.Parent != nil {
			p = p.Parent
		} else {
			break
		}
	}
	return p
}

func (ast *Ast) hasNonExp() bool {
	if ast.Mode != EXP {
		return true
	} else {
		for _, c := range ast.Children {
			if v, ok := c.(*Ast); ok {
				if v.hasNonExp() {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (ast *Ast) debug(depth int, max int) {
	if depth >= max {
		return
	}
	for i := 0; i < depth; i++ {
		fmt.Printf("%c", '-')
	}
	fmt.Printf("TagName: %s Mode: %s Children: %d [[ \n", ast.TagName, ast.ModeStr(), len(ast.Children))
	for _, a := range ast.Children {
		if _, ok := a.(*Ast); ok {
			b := (*Ast)(a.(*Ast))
			b.debug(depth+1, max)
		} else {
			if depth+1 < max {
				t := (Token)(a.(Token))
				for i := 0; i < depth+1; i++ {
					fmt.Printf("%c", '-')
				}
				t.debug()
			}
		}
	}
	for i := 0; i < depth; i++ {
		fmt.Printf("%c", '-')
	}

	logger.Info("]]")
}

func regMatch(reg string, text string) (string, error) {
	regc, err := regexp.Compile(reg)
	if err != nil {
		panic(err)
		return "", err
	}
	found := regc.FindStringIndex(text)
	if found != nil {
		return text[found[0]:found[1]], nil
	}
	return "", nil
}
