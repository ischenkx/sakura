package protocol

type subOpts struct {
	Topics []string
	Time int64
}

type unsubOpts struct {
	Topics []string
	All bool
	Time int64
}

type pubOpts struct {
	Data []byte
	Time int64
}

