package gojson

type NodeType int

const (
	NDUnknown NodeType = iota
	NDJson
	NDObject
	NDPair
	NDArray
	NDValue
	NDEof
)

func (typ NodeType) String() string {
	switch typ {
	case NDJson:
		return "NDJson"
	case NDObject:
		return "NDObject"
	case NDPair:
		return "NDPair"
	case NDArray:
		return "NDArray"
	case NDValue:
		return "NDValue"
	case NDEof:
		return "NDEof"
	default:
		return "NDUnknown"
	}
}

type Node struct {
	Type NodeType
	// Children
	// on Object, Children = []Pair as Member
	// on Array, Children = []Value as Elements
	Children *[]Node
	Key      string
	Val      *Token
}

func NewNode(typ NodeType, nds *[]Node, key string, val *Token) *Node {
	return &Node{
		Type:     typ,
		Children: nds,
		Key:      key,
		Val:      val,
	}
}
