package tools

type Values struct {
	data []float64
}

func NewValues(data []float64) *Values {
	return &Values{data: data}
}

func (values *Values) Current() float64 {
	return values.data[len(values.data)-1]
}

func (values *Values) Previous() float64 {
	return values.data[len(values.data)-2]
}

func (values *Values) All() []float64 {
	return values.data
}

func (values *Values) Len() int {
	return len(values.data)
}

// At returns the value at the specified index, where 0 is the current value,
// 1 is the previous value, and so on.
func (values *Values) At(index int) float64 {
	return values.data[len(values.data)-1-index]
}
