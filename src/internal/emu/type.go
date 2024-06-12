package emu

type StreamType uint8
type NodeType uint8
type Operation uint8
type LocationType uint8
type LocationDirection uint8

const (
	IN StreamType = iota
	OUT
)

const (
	COMPUTE NodeType = iota
	DAMAGED
)

const (
	MOV Operation = iota
	SAV
	SWP
	SUB
	ADD
	NOP
	NEG
	JEZ
	JMP
	JNZ
	JGZ
	JLZ
	JRO
	RES
)

const (
	NUMBER LocationType = iota
	ADDRESS
)

const (
	UP LocationDirection = iota
	RIGHT
	DOWN
	LEFT
	NIL
	ACC
	ANY
	LAST
)

type Location struct {
	Number    int16
	Direction LocationDirection
}

type Instruction struct {
	Operation Operation
	SrcType   LocationType
	Src       Location
	DestType  LocationType
	Dest      Location
}

type Stream struct {
	Index    uint8      `json:"index"`
	Name     string     `json:"name,omitempty"`
	Type     StreamType `json:"type,omitempty"`
	Values   []int16    `json:"values,omitempty"`
	MaxValue int16      `json:"max_value,omitempty"`
	MinValue int16      `json:"min_value,omitempty"`
}

type NodeCode struct {
	Index uint8    `json:"index"`
	Code  []string `json:"code"`
}
