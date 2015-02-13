package basic

type Item map[string]interface{}

func NewItems() Item {
	it := make(map[string]interface{})
	return it
}
