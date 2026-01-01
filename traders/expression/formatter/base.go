package formatter

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	NodeFunction NodeType = iota
	NodeValue
)

type FormatterNode struct {
	type_    NodeType
	package_ string
	value    string
	children []*FormatterNode
}

type Formatter interface {
	Format() *FormatterNode
}

func Value(package_ string, value string) *FormatterNode {
	return &FormatterNode{
		type_:    NodeValue,
		package_: package_,
		value:    value,
		children: make([]*FormatterNode, 0),
	}
}

func IntValue(value int) *FormatterNode {
	return &FormatterNode{
		type_:    NodeValue,
		package_: "",
		value:    fmt.Sprintf("%d", value),
		children: make([]*FormatterNode, 0),
	}
}

func FloatValue(value float64) *FormatterNode {
	return &FormatterNode{
		type_:    NodeValue,
		package_: "",
		value:    fmt.Sprintf("%f", value),
		children: make([]*FormatterNode, 0),
	}
}

func Function(package_ string, value string, children ...*FormatterNode) *FormatterNode {
	return &FormatterNode{
		type_:    NodeFunction,
		package_: package_,
		value:    value,
		children: children,
	}
}

func FunctionWithChildren[T Formatter](package_ string, value string, children ...T) *FormatterNode {
	node := &FormatterNode{
		type_:    NodeFunction,
		package_: package_,
		value:    value,
		children: make([]*FormatterNode, 0, len(children)),
	}

	for _, child := range children {
		node.children = append(node.children, child.Format())
	}

	return node
}

func (n *FormatterNode) fullName() string {
	if n.package_ == "" {
		return n.value
	}

	return n.package_ + "." + n.value
}

func (n *FormatterNode) Compact() string {
	switch n.type_ {
	case NodeValue:
		return n.fullName()

	case NodeFunction:
		result := n.fullName() + "("
		for i, child := range n.children {
			if i > 0 {
				result += ", "
			}
			result += child.Compact()
		}
		return result + ")"

	default:
		panic("unknown node type")
	}
}

func (n *FormatterNode) Detailed() string {
	return n.detailedWithIndent(0)
}

func (n *FormatterNode) detailedWithIndent(indent int) string {
	ind := strings.Repeat(" ", indent*2)

	switch n.type_ {
	case NodeValue:
		return ind + n.fullName()

	case NodeFunction:
		// If there is no children or only children that are values, print in one line
		allChildrenAreValues := true
		for _, child := range n.children {
			if child.type_ != NodeValue {
				allChildrenAreValues = false
				break
			}
		}
		if allChildrenAreValues {
			return ind + n.Compact()
		}

		rows := []string{
			ind + n.fullName() + "(",
		}

		for _, child := range n.children {
			rows = append(rows, child.detailedWithIndent(indent+1)+",")
		}

		rows = append(rows, ind+")")

		return strings.Join(rows, "\n")
	default:
		panic("unknown node type")
	}
}
