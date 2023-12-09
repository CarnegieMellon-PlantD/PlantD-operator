package tree

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type Node struct {
	Name           string
	Id             pcommon.SpanID
	StartTimestamp pcommon.Timestamp
	EndTimestamp   pcommon.Timestamp
	Children       []*Node
}

type Forest struct {
	Roots []*Node
}

func NewNode(span ptrace.Span) *Node {
	return &Node{
		Name:           span.Name(),
		Id:             span.SpanID(),
		StartTimestamp: span.StartTimestamp(),
		EndTimestamp:   span.EndTimestamp(),
	}
}

func NewNodeWithId(id pcommon.SpanID) *Node {
	return &Node{
		Id: id,
	}
}

func (n *Node) Update(span ptrace.Span) {
	n.Name = span.Name()
	n.StartTimestamp = span.StartTimestamp()
	n.EndTimestamp = span.EndTimestamp()
}

func NewForest() *Forest {
	return &Forest{
		Roots: nil,
	}
}

func (f *Forest) AddSpans(spans ptrace.SpanSlice) {
	lookup := make(map[pcommon.SpanID]*Node)

	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		spanId := span.SpanID()
		currNode, currExist := lookup[spanId]
		if currExist {
			currNode.Update(span)
		} else {
			currNode = NewNode(span)
			lookup[spanId] = currNode
		}

		parentId := span.ParentSpanID()
		if span.ParentSpanID().IsEmpty() {
			f.Roots = append(f.Roots, currNode)
		} else {
			parentNode, parentExist := lookup[parentId]
			if !parentExist {
				parentNode = NewNodeWithId(parentId)
				lookup[parentId] = parentNode
			}
			parentNode.Children = append(parentNode.Children, currNode)
		}
	}
}

func NewNodeWithName(name string) *Node {
	return &Node{
		Name: name,
	}
}

func (f *Forest) Print() {
	newLineNode := NewNodeWithName("\n")
	for _, n := range f.Roots {
		q := make([]*Node, 2)
		q[0] = n
		q[1] = newLineNode
		for len(q) > 0 {
			top := q[0]
			q = q[1:]
			for _, child := range top.Children {
				q = append(q, child)
			}
			if top.Name == "\n" && len(q) > 0 {
				q = append(q, newLineNode)
				fmt.Print("\n")
			} else {
				fmt.Printf("%s\t", top.Name)
			}
		}
		fmt.Print("\n---------------\n")
	}
}
