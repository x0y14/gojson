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

func NewJson(node *Node, rootNodeType NodeType) *Json {
	return &Json{
		node:         node,
		RootNodeType: rootNodeType,
	}
}

type Json struct {
	node         *Node
	RootNodeType NodeType
}

func (j *Json) Map() (map[string]interface{}, error) {
	if j.RootNodeType != NDObject {
		panic(JsonError{})
	}

	obj, err := j.ObjectMapping(j.node)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (j *Json) Array() ([]interface{}, error) {
	if j.RootNodeType != NDArray {
		panic(JsonError{})
	}

	arr, err := j.ArrayMapping(j.node)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (j *Json) ObjectMapping(obj *Node) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for _, pair := range *obj.Children {
		key, val, err := j.PairMapping(&pair)
		if err != nil {
			return nil, err
		}

		result[key] = val
	}
	return result, nil
}

func (j *Json) PairMapping(pair *Node) (string, interface{}, error) {
	key := pair.Key
	children := *pair.Children
	child := children[0]
	var value interface{}
	var err error

	switch child.Type {
	case NDObject:
		value, err = j.ObjectMapping(&child)
		if err != nil {
			return "", nil, err
		}
	case NDArray:
		value, err = j.ArrayMapping(&child)
	default:
		value = j.ValueMapping(&child)
	}
	return key, value, nil
}

func (j *Json) ArrayMapping(arr *Node) ([]interface{}, error) {
	var result []interface{}
	for _, element := range *arr.Children {
		value, err := j.ElementMapping(&element)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (j *Json) ElementMapping(element *Node) (interface{}, error) {
	switch element.Type {
	case NDObject:
		value, err := j.ObjectMapping(element)
		if err != nil {
			return nil, err
		}
		return value, nil
	case NDArray:
		value, err := j.ArrayMapping(element)
		if err != nil {
			return nil, err
		}
		return value, nil
	default:
		value := j.ValueMapping(element)
		return value, nil
	}
}

func (j *Json) ValueMapping(val *Node) interface{} {
	switch val.Val.Type {
	case TTrue, TFalse:
		return val.Val.LoadAsBoolean()
	case TNumber:
		return val.Val.LoadAsFloat64()
	case TString:
		return val.Val.LoadAsString()
	default:
		return val.Val.LoadAsNull()
	}
}

// Tree show only
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
