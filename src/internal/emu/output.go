package emu

var output []Stream = make([]Stream, 0)

func NewOutputStream(index uint8) Stream {
	return Stream{
		Index:  index,
		Type:   OUT,
		Values: make([]int16, 0),
	}
}

func AddOutputStream(stream Stream) {
	output = append(output, stream)
}

func AddOutputValue(index uint8, value int16) bool {
	for i := range output {
		if output[i].Index == index {
			output[i].Values = append(output[i].Values, value)
			return true
		}
	}
	return false
}

func GetOutput() []Stream {
	return output
}

func ClearOutput() {
	output = make([]Stream, 0)
}
