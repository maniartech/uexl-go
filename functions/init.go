package functions

func init() {
	var a Function = add
	Registry.Register("ADD", a)
}
