/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jsonpath

import "fmt"

// NodeType identifies the type of a parse tree node.
type NodeType int

// Type returns itself and provides an easy default implementation
func (t NodeType) Type() NodeType {
	return t
}

func (t NodeType) String() string {
	return NodeTypeName[t]
}

const (
	// NodeText is a text node type code
	NodeText NodeType = iota
	// NodeArray is an array node type code
	NodeArray
	// NodeList is a list node type code
	NodeList
	// NodeField is a field node type code
	NodeField
	// NodeIdentifier is an identifier node type code
	NodeIdentifier
	// NodeFilter is a filter node type code
	NodeFilter
	// NodeInt is an integer node type code
	NodeInt
	// NodeFloat is a float node type code
	NodeFloat
	// NodeWildcard is a wildcard node type code
	NodeWildcard
	// NodeRecursive is a recursive node type code
	NodeRecursive
	// NodeUnion is a union node type code
	NodeUnion
	// NodeBool is a boolean node type code
	NodeBool
)

// NodeTypeName maps node type code to node type text representation
var NodeTypeName = map[NodeType]string{
	NodeText:       "NodeText",
	NodeArray:      "NodeArray",
	NodeList:       "NodeList",
	NodeField:      "NodeField",
	NodeIdentifier: "NodeIdentifier",
	NodeFilter:     "NodeFilter",
	NodeInt:        "NodeInt",
	NodeFloat:      "NodeFloat",
	NodeWildcard:   "NodeWildcard",
	NodeRecursive:  "NodeRecursive",
	NodeUnion:      "NodeUnion",
	NodeBool:       "NodeBool",
}

// Node represents a parse tree node
type Node interface {
	Type() NodeType
	String() string
}

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Nodes []Node // The element nodes in lexical order.
}

func newList() *ListNode {
	return &ListNode{NodeType: NodeList}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) String() string {
	return l.Type().String()
}

// TextNode holds plain text.
type TextNode struct {
	NodeType
	Text string // The text; may span newlines.
}

func newText(text string) *TextNode {
	return &TextNode{NodeType: NodeText, Text: text}
}

func (t *TextNode) String() string {
	return fmt.Sprintf("%s: %s", t.Type(), t.Text)
}

// FieldNode holds field of struct
type FieldNode struct {
	NodeType
	Value string
}

func newField(value string) *FieldNode {
	return &FieldNode{NodeType: NodeField, Value: value}
}

func (f *FieldNode) String() string {
	return fmt.Sprintf("%s: %s", f.Type(), f.Value)
}

// IdentifierNode holds an identifier
type IdentifierNode struct {
	NodeType
	Name string
}

func newIdentifier(value string) *IdentifierNode {
	return &IdentifierNode{
		NodeType: NodeIdentifier,
		Name:     value,
	}
}

func (f *IdentifierNode) String() string {
	return fmt.Sprintf("%s: %s", f.Type(), f.Name)
}

// ParamsEntry holds param information for ArrayNode
type ParamsEntry struct {
	Value   int
	Known   bool // whether the value is known when parse it
	Derived bool
}

// ArrayNode holds start, end, step information for array index selection
type ArrayNode struct {
	NodeType
	Params [3]ParamsEntry // start, end, step
}

func newArray(params [3]ParamsEntry) *ArrayNode {
	return &ArrayNode{
		NodeType: NodeArray,
		Params:   params,
	}
}

func (a *ArrayNode) String() string {
	return fmt.Sprintf("%s: %v", a.Type(), a.Params)
}

// FilterNode holds operand and operator information for filter
type FilterNode struct {
	NodeType
	Left     *ListNode
	Right    *ListNode
	Operator string
}

func newFilter(left, right *ListNode, operator string) *FilterNode {
	return &FilterNode{
		NodeType: NodeFilter,
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (f *FilterNode) String() string {
	return fmt.Sprintf("%s: %s %s %s", f.Type(), f.Left, f.Operator, f.Right)
}

// IntNode holds integer value
type IntNode struct {
	NodeType
	Value int
}

func newInt(num int) *IntNode {
	return &IntNode{NodeType: NodeInt, Value: num}
}

func (i *IntNode) String() string {
	return fmt.Sprintf("%s: %d", i.Type(), i.Value)
}

// FloatNode holds float value
type FloatNode struct {
	NodeType
	Value float64
}

func newFloat(num float64) *FloatNode {
	return &FloatNode{NodeType: NodeFloat, Value: num}
}

func (i *FloatNode) String() string {
	return fmt.Sprintf("%s: %f", i.Type(), i.Value)
}

// WildcardNode means a wildcard
type WildcardNode struct {
	NodeType
}

func newWildcard() *WildcardNode {
	return &WildcardNode{NodeType: NodeWildcard}
}

func (i *WildcardNode) String() string {
	return i.Type().String()
}

// RecursiveNode means a recursive descent operator
type RecursiveNode struct {
	NodeType
}

func newRecursive() *RecursiveNode {
	return &RecursiveNode{NodeType: NodeRecursive}
}

func (r *RecursiveNode) String() string {
	return r.Type().String()
}

// UnionNode is union of ListNode
type UnionNode struct {
	NodeType
	Nodes []*ListNode
}

func newUnion(nodes []*ListNode) *UnionNode {
	return &UnionNode{NodeType: NodeUnion, Nodes: nodes}
}

func (u *UnionNode) String() string {
	return u.Type().String()
}

// BoolNode holds bool value
type BoolNode struct {
	NodeType
	Value bool
}

func newBool(value bool) *BoolNode {
	return &BoolNode{NodeType: NodeBool, Value: value}
}

func (b *BoolNode) String() string {
	return fmt.Sprintf("%s: %t", b.Type(), b.Value)
}
