package obj

// StringValue string value
type StringValue string

// NewString create string object
func NewString(str string) *Obj {
	return &Obj{
		T:     ObjString,
		Value: StringValue(str),
	}
}
