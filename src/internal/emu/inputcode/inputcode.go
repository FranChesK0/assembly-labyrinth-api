package inputcode

type InputCode struct {
	Lines  []string
	Labels map[string]uint8
}

func NewInputCode() InputCode {
	return InputCode{
		Lines:  make([]string, 0),
		Labels: make(map[string]uint8),
	}
}

func (ic *InputCode) AddLine(line string) {
	ic.Lines = append(ic.Lines, line)
}
