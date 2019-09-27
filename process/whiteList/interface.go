package whiteList

// WhiteList is that allow to make operation on white list data
type WhiteList interface {
	AddHash([]byte)
	RemoveHash([]byte)
	RemoveHashes([][]byte)
	IsInWhiteList([]byte) bool
}
