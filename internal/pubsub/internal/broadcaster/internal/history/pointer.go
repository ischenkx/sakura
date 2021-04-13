package history

type SnapshotInfo struct {
	offset int
	epoch  int64
}

//todo: replace h with something like group id
type Pointer struct {
	h    *History
	info SnapshotInfo
}
