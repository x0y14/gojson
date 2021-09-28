package gojson

import (
	"fmt"
	"strings"
)

// Json     := Object
//		     | Array

// Object   := "{" Member? "}"
// Array    := "[" Elements? "]"

// Member   := Pair
//		     | Pair "," Member
// Pair     := TString ":" Value

// Elements := Value
//           | Value "," Elements

// Value    := TString
//			 | TNumber
//			 | TBoolean
// 			 | TNull
//           | Object
//           | Array

func NewJson(node *Node) *Json {
	return &Json{
		node: node,
	}
}

type Json struct {
	node *Node
}

func (j *Json) Map() map[string]interface{} {
	return nil
}

func (j *Json) Tree() {
	rootNode := j.node
	switch rootNode.Type {
	case NDObject:
		j.ObjectTree(1, j.node)
	case NDArray:
		j.ArrayTree(1, j.node)
	}
}

func (j *Json) ObjectTree(nest int, obj *Node) {
	fmt.Printf("%v{\n", Indent(nest))
	for _, pair := range *obj.Children {
		j.PairTree(nest+1, &pair)
	}
	fmt.Printf("%v}\n", Indent(nest))
}

func (j *Json) PairTree(nest int, pair *Node) {
	fmt.Printf("%v%v :\n", Indent(nest), pair.Key)
	children := *pair.Children
	child := children[0]
	switch child.Type {
	case NDObject:
		j.ObjectTree(nest+1, &child)
	case NDArray:
		j.ArrayTree(nest+1, &child)
	default:
		j.ShowValue(nest+1, &child)
	}
}

func (j *Json) ArrayTree(nest int, arr *Node) {
	fmt.Printf("%v[\n", Indent(nest))
	for i, element := range *arr.Children {
		j.ElementTree(i, nest+1, &element)
	}
	fmt.Printf("%v]\n", Indent(nest))
}

func (j *Json) ElementTree(index int, nest int, element *Node) {
	fmt.Printf("%v%v :\n", Indent(nest), index)
	switch element.Type {
	case NDObject:
		j.ObjectTree(nest+1, element)
	case NDArray:
		j.ArrayTree(nest+1, element)
	default:
		j.ShowValue(nest+1, element)
	}
}

func (j *Json) ShowValue(nest int, val *Node) {
	switch val.Val.Type {
	case TTrue, TFalse:
		fmt.Printf("%v`%v`(%v)\n", Indent(nest), val.Val.LoadAsBoolean(), val.Val.Type.String())
	case TNull:
		fmt.Printf("%v`%v`(%v)\n", Indent(nest), val.Val.LoadAsNull(), val.Val.Type.String())
	case TNumber:
		fmt.Printf("%v`%v`(%v)\n", Indent(nest), val.Val.LoadAsFloat64(), val.Val.Type.String())
	case TString:
		fmt.Printf("%v`%v`(%v)\n", Indent(nest), val.Val.LoadAsString(), val.Val.Type.String())
	}
}

func Indent(nest int) string {
	return strings.Repeat("  ", nest)
}
