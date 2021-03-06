package testdata

// Fetcher is a main fetcher for this module
type Fetcher struct{}

// FetchOrders fetching orders for me.
// There is also second line
func (Fetcher) FetchOrders() error { return nil }

func (Fetcher) FetchNoComments() {}

type (
	// F1 comment
	F1 int
	// F2 comment
	F2 int
)

// And this comment of some variable
var someVar = ""

const (
	// Docs of constants are also available
	MyConst, My1Const = iota, iota
	// Const 2 comment
	My2Const F1 = 15
)
