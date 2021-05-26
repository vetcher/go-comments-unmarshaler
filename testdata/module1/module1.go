package module1

// Module1Func is a module one function
func Module1Func() {}

var _ = Module1Func
var _ = module1PrivateFunc

// module1PrivateFunc is module private func.
func module1PrivateFunc() {}
