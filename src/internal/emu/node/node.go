package node

import (
	"errors"
	"strconv"
	"strings"

	"github.com/franchesko/assembly-labyrinth/src/internal/emu"
	"github.com/franchesko/assembly-labyrinth/src/internal/emu/inputcode"
)

type Node struct {
	Index          uint8
	Visible        bool
	Blocked        bool
	CursorPosition uint8
	Instructions   []*emu.Instruction
	ACC            int16
	BAK            int16
	OutputPort     *Node
	Last           *Node
	OutputValue    int16
	Ports          [4]*Node
}

type readResult struct {
	Blocked bool
	Value   int16
}

func NewNode() *Node {
	return &Node{
		Visible:        false,
		Blocked:        false,
		CursorPosition: 0,
		Instructions:   make([]*emu.Instruction, 0),
		ACC:            0,
		BAK:            0,
		OutputPort:     nil,
		Last:           nil,
		OutputValue:    0,
		Ports:          [4]*Node{nil, nil, nil, nil},
	}
}

func (n *Node) CreateInstruction(op emu.Operation) *emu.Instruction {
	ins := &emu.Instruction{
		Operation: op,
	}
	n.Instructions = append(n.Instructions, ins)

	return ins
}

func (n *Node) ParseCode(ic *inputcode.InputCode) error {
	for i, line := range ic.Lines {
		if ind := strings.Index(line, ":"); ind != -1 {
			label := line[:ind]
			ic.Labels[label] = uint8(i)

			rem := strings.TrimSpace(line[ind+1:])
			if len(rem) == 0 {
				rem = "NOP"
			}
			ic.Lines[i] = rem
		}
	}

	var err error
	for _, line := range ic.Lines {
		if err = n.ParseLine(ic, line); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) ParseLine(ic *inputcode.InputCode, line string) error {
	if len(line) <= 2 {
		return errors.New("invalid line length")
	}

	strIns := line[:3]
	var err error
	switch strIns {
	case "MOV":
		if err = n.parseMov(line); err != nil {
			return err
		}
	case "SUB":
		if err = n.parseOneArg(ic, line, emu.SUB); err != nil {
			return err
		}
	case "ADD":
		if err = n.parseOneArg(ic, line, emu.ADD); err != nil {
			return err
		}
	case "JEZ":
		if err = n.parseOneArg(ic, line, emu.JEZ); err != nil {
			return err
		}
	case "JMP":
		if err = n.parseOneArg(ic, line, emu.JMP); err != nil {
			return err
		}
	case "JNZ":
		if err = n.parseOneArg(ic, line, emu.JNZ); err != nil {
			return err
		}
	case "JGZ":
		if err = n.parseOneArg(ic, line, emu.JGZ); err != nil {
			return err
		}
	case "JLZ":
		if err = n.parseOneArg(ic, line, emu.JLZ); err != nil {
			return err
		}
	case "JRO":
		if err = n.parseOneArg(ic, line, emu.JRO); err != nil {
			return err
		}
	case "SAV":
		n.CreateInstruction(emu.SAV)
	case "SWP":
		n.CreateInstruction(emu.SWP)
	case "NOP":
		n.CreateInstruction(emu.NOP)
	case "NEG":
		n.CreateInstruction(emu.NEG)
	case "RES":
		n.CreateInstruction(emu.RES)
	default:
		return errors.New("invalid instruction")
	}

	return nil
}

func (n *Node) Read(locType emu.LocationType, loc emu.Location) (readResult, error) {
	res := readResult{
		Blocked: false,
	}

	if n.OutputPort != nil {
		return res, nil
	}

	if locType == emu.NUMBER {
		res.Value = loc.Number
	} else {
		switch loc.Direction {
		case emu.NIL:
			res.Value = 0
		case emu.ACC:
			res.Value = n.ACC
		case emu.UP, emu.RIGHT, emu.DOWN, emu.LEFT, emu.ANY, emu.LAST:
			readFrom := n.getInputPort(loc.Direction)
			if readFrom != nil && readFrom.OutputPort == n {
				res.Value = readFrom.OutputValue
				res.Blocked = false

				readFrom.OutputValue = 0
				readFrom.OutputPort = nil
				readFrom.MoveCursor()

				if loc.Direction == emu.ANY {
					n.Last = readFrom
				}
			} else if readFrom != nil && loc.Direction == emu.LAST {
				res.Value = 0
			} else {
				res.Blocked = true
			}
		default:
			return readResult{}, errors.New("unknown direction")
		}
	}

	return res, nil
}

func (n *Node) Write(dir emu.LocationDirection, value int16) (bool, error) {
	switch dir {
	case emu.ACC:
		n.ACC = value
	case emu.UP, emu.RIGHT, emu.DOWN, emu.LEFT, emu.ANY, emu.LAST:
		dest := n.getOutputPort(dir)
		if dest != nil && n.OutputPort == nil {
			n.OutputPort = dest
			n.OutputValue = value
			if dir == emu.ANY {
				n.Last = dest
			}
		}
		return true, nil
	case emu.NIL:
		return false, errors.New("unable to write")
	default:
		return false, errors.New("nowhere to write")
	}

	return false, nil
}

func (n *Node) MoveCursor() {
	n.CursorPosition += 1
}

func (n *Node) Tick() error {
	n.Blocked = true

	if n.CursorPosition >= uint8(len(n.Instructions)) {
		n.CursorPosition = 0
	}
	ins := n.Instructions[n.CursorPosition]
	switch ins.Operation {
	case emu.MOV:
		read, err := n.Read(ins.SrcType, ins.Src)
		if err != nil {
			return err
		}
		if read.Blocked {
			return nil
		}

		blocked, err := n.Write(ins.Dest.Direction, read.Value)
		if err != nil {
			return err
		}
		if blocked {
			return nil
		}
	case emu.ADD:
		read, err := n.Read(ins.SrcType, ins.Src)
		if err != nil {
			return err
		}
		if read.Blocked {
			return nil
		}

		n.ACC += read.Value
		n.normalizeACC()
	case emu.SUB:
		read, err := n.Read(ins.SrcType, ins.Src)
		if err != nil {
			return err
		}
		if read.Blocked {
			return nil
		}

		n.ACC -= read.Value
		n.normalizeACC()
	case emu.JMP:
		n.setCursorPosition(ins.Src.Number)
		return nil
	case emu.JRO:
		n.setCursorPosition(int16(n.CursorPosition) + ins.Src.Number)
		return nil
	case emu.JEZ:
		if n.ACC == 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case emu.JGZ:
		if n.ACC > 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case emu.JLZ:
		if n.ACC < 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case emu.JNZ:
		if n.ACC != 0 {
			n.setCursorPosition(ins.Src.Number)
		}
	case emu.SWP:
		tmp := n.BAK
		n.BAK = n.ACC
		n.ACC = tmp
	case emu.SAV:
		n.BAK = n.ACC
	case emu.NEG:
		n.ACC = -1 * n.ACC
	case emu.NOP:
	case emu.RES:
		emu.AddOutputValue(n.Index, n.ACC)
	default:
		return errors.New("unknown operation")
	}

	n.Blocked = false
	n.MoveCursor()
	return nil
}

func (n *Node) parseMov(line string) error {
	if len(line) <= 3 {
		return errors.New("wrong mov instruction format")
	}
	rem := line[4:]
	var tokens []string
	if strings.Contains(rem, ", ") {
		tokens = strings.Split(rem, ", ")
	} else if strings.Contains(rem, ",") {
		tokens = strings.Split(rem, ",")
	} else {
		tokens = strings.Split(rem, " ")
	}
	if len(tokens) != 2 {
		return errors.New("wrong mov instruction format")
	}

	ins := n.CreateInstruction(emu.MOV)
	var err error
	if err = parseLocation(tokens[0], &ins.SrcType, &ins.Src); err != nil {
		return err
	}
	if err = parseLocation(tokens[1], &ins.DestType, &ins.Dest); err != nil {
		return err
	}

	return nil
}

func (n *Node) parseOneArg(ic *inputcode.InputCode, line string, op emu.Operation) error {
	if len(line) <= 3 {
		return errors.New("wrong one arg instruction format")
	}
	rem := line[4:]
	ins := n.CreateInstruction(op)

	switch op {
	case emu.JEZ, emu.JMP, emu.JNZ, emu.JGZ, emu.JLZ:
		for label, pos := range ic.Labels {
			if rem == label {
				ins.SrcType = emu.NUMBER
				ins.Src.Number = int16(pos)
			}
		}
	default:
		if err := parseLocation(rem, &ins.SrcType, &ins.Src); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) setCursorPosition(pos int16) {
	if pos >= int16(len(n.Instructions)) || pos < 0 {
		pos = 0
	}
	n.CursorPosition = uint8(pos)
}

func (n *Node) getInputPort(dir emu.LocationDirection) *Node {
	if dir == emu.ANY {
		dirs := []emu.LocationDirection{emu.LEFT, emu.RIGHT, emu.UP, emu.DOWN}
		for _, d := range dirs {
			port := n.Ports[d]
			if port != nil && port.OutputPort == n {
				return port
			}
		}
		return nil
	} else if dir == emu.LAST {
		return n.Last
	} else {
		return n.Ports[dir]
	}
}

func (n *Node) getOutputPort(dir emu.LocationDirection) *Node {
	if dir == emu.ANY {
		dirs := []emu.LocationDirection{emu.UP, emu.LEFT, emu.RIGHT, emu.DOWN}
		for _, d := range dirs {
			port := n.Ports[d]
			if port != nil {
				ins := port.Instructions[port.CursorPosition]
				if ins.Operation == emu.MOV && ins.SrcType == emu.ADDRESS && (ins.Src.Direction == emu.ANY || port.Ports[ins.Src.Direction] == n) {
					return port
				}
			}
		}
		return nil
	} else if dir == emu.LAST {
		return n.Last
	} else {
		return n.Ports[dir]
	}
}

func (n *Node) normalizeACC() {
	if n.ACC > emu.MaxACC {
		n.ACC = emu.MaxACC
	}
	if n.ACC < emu.MinACC {
		n.ACC = emu.MinACC
	}
}

func parseLocation(strLoc string, locType *emu.LocationType, loc *emu.Location) error {
	if strLoc == "" {
		return errors.New("no source was found")
	}

	*locType = emu.ADDRESS
	switch strLoc {
	case "UP":
		loc.Direction = emu.UP
	case "DOWN":
		loc.Direction = emu.DOWN
	case "LEFT":
		loc.Direction = emu.LEFT
	case "RIGHT":
		loc.Direction = emu.RIGHT
	case "ACC":
		loc.Direction = emu.ACC
	case "NIL":
		loc.Direction = emu.NIL
	case "ANY":
		loc.Direction = emu.ANY
	case "LAST":
		loc.Direction = emu.LAST
	default:
		num, err := strconv.Atoi(strLoc)
		if err != nil {
			return err
		}

		*locType = emu.NUMBER
		loc.Number = int16(num)
	}

	return nil
}
