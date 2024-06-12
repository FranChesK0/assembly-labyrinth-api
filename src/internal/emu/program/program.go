package program

import (
	"errors"
	"strings"

	"github.com/franchesko/assembly-labyrinth/src/internal/emu"
	"github.com/franchesko/assembly-labyrinth/src/internal/emu/inputcode"
	"github.com/franchesko/assembly-labyrinth/src/internal/emu/node"
	"github.com/franchesko/assembly-labyrinth/src/internal/emu/nodelist"
)

type Program struct {
	Nodes       []*node.Node
	NodeList    *nodelist.NodeList
	ActiveNodes *nodelist.NodeList
}

func Run(streams []emu.Stream, nodesCode []emu.NodeCode) ([]emu.Stream, error) {
	prog := NewProgram()
	if err := prog.LoadStreams(streams); err != nil {
		return nil, err
	}
	if err := prog.LoadCode(nodesCode); err != nil {
		return nil, err
	}

	dtmCount := 0
	for dtmCount < 5 {
		allBlocked, err := prog.Tick()
		if err != nil {
			emu.ClearOutput()
			return nil, err
		}

		if allBlocked {
			dtmCount++
		} else {
			dtmCount = 0
		}
	}

	output := emu.GetOutput()
	emu.ClearOutput()
	return output, nil
}

func NewProgram() *Program {
	nodes := make([]*node.Node, 0, emu.NodesNumber)
	var n *node.Node
	for i := range emu.NodesNumber {
		n = node.NewNode()
		n.Visible = true
		n.Index = uint8(i)
		nodes = append(nodes, n)
	}
	p := &Program{
		Nodes:       nodes,
		NodeList:    nil,
		ActiveNodes: nil,
	}

	for i := range p.Nodes {
		if i != 8 && i != 9 && i != 10 && i != 11 {
			p.Nodes[i].Ports[emu.DOWN] = p.Nodes[i+4]
		}
		if i != 0 && i != 1 && i != 2 && i != 3 {
			p.Nodes[i].Ports[emu.UP] = p.Nodes[i-4]
		}
		if i != 3 && i != 7 && i != 11 {
			p.Nodes[i].Ports[emu.RIGHT] = p.Nodes[i+1]
		}
		if i != 0 && i != 4 && i != 8 {
			p.Nodes[i].Ports[emu.LEFT] = p.Nodes[i-1]
		}
	}

	return p
}

func (p *Program) Tick() (bool, error) {
	allBlocked := true
	var err error
	for list := p.ActiveNodes; list != nil; list = list.Next {
		if err = list.Node.Tick(); err != nil {
			return false, err
		}
		allBlocked = allBlocked && list.Node.Blocked
	}
	return allBlocked, nil
}

func (p *Program) LoadStreams(streams []emu.Stream) error {
	for _, stream := range streams {
		if stream.Type == emu.IN {
			n := p.createInputNode(stream)
			p.ActiveNodes = nodelist.Append(p.ActiveNodes, n)
		} else if stream.Type == emu.OUT {
			n := p.createOutputNode(stream)
			p.ActiveNodes = nodelist.Append(p.ActiveNodes, n)
		} else {
			return errors.New("unkown stream type")
		}
	}
	return nil
}

func (p *Program) LoadCode(nodesCode []emu.NodeCode) error {
	if len(nodesCode) != emu.NodesNumber {
		return errors.New("wrong nodes number")
	}

	allInput := make([]inputcode.InputCode, 0)
	for range emu.NodesNumber {
		allInput = append(allInput, inputcode.NewInputCode())
	}

	for i, nodeCode := range nodesCode {
		for _, line := range nodeCode.Code {
			formatted := strings.ToUpper(strings.TrimSpace(line))
			allInput[i].AddLine(formatted)
		}
	}

	for _, n := range p.Nodes {
		if err := n.ParseCode(&allInput[n.Index]); err != nil {
			return err
		}
		if len(n.Instructions) > 0 {
			p.ActiveNodes = nodelist.Append(p.ActiveNodes, n)
		}
	}

	return nil
}

func (p *Program) createNode() *node.Node {
	n := node.NewNode()
	p.NodeList = nodelist.Append(p.NodeList, n)
	return n
}

func (p *Program) createInputNode(stream emu.Stream) *node.Node {
	inputNode := p.createNode()
	belowNode := p.Nodes[stream.Index]

	inputNode.Ports[emu.DOWN] = belowNode
	belowNode.Ports[emu.UP] = inputNode

	for _, value := range stream.Values {
		ins := inputNode.CreateInstruction(emu.MOV)
		ins.SrcType = emu.NUMBER
		ins.Src.Number = value
		ins.DestType = emu.ADDRESS
		ins.Dest.Direction = emu.DOWN
	}

	ins := inputNode.CreateInstruction(emu.JRO)
	ins.SrcType = emu.NUMBER
	ins.Src.Number = 0

	return inputNode
}

func (p *Program) createOutputNode(stream emu.Stream) *node.Node {
	outputNode := p.createNode()
	outputNode.Index = stream.Index
	aboveNode := p.Nodes[stream.Index]

	outputNode.Ports[emu.UP] = aboveNode
	aboveNode.Ports[emu.DOWN] = outputNode

	ins := outputNode.CreateInstruction(emu.MOV)
	ins.SrcType = emu.ADDRESS
	ins.Src.Direction = emu.UP
	ins.DestType = emu.ADDRESS
	ins.Dest.Direction = emu.ACC
	outputNode.CreateInstruction(emu.RES)

	emu.AddOutputStream(emu.NewOutputStream(stream.Index))

	return outputNode
}
