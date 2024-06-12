package nodelist

import "github.com/franchesko/assembly-labyrinth/src/internal/emu/node"

type NodeList struct {
	Node *node.Node
	Next *NodeList
}

func Append(list *NodeList, n *node.Node) *NodeList {
	tail := &NodeList{
		Node: n,
		Next: nil,
	}
	if list == nil {
		return tail
	}

	head := list
	for head.Next != nil {
		head = head.Next
	}
	head.Next = tail
	return list
}
