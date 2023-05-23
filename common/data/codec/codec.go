package codec

type Converter[From, To any] interface {
	Convert(From) (To, error)
}

type Codec[T1, T2 any] interface {
	Encoder() Converter[T1, T2]
	Decoder() Converter[T2, T1]
}

type Binary[T any] Codec[T, []byte]
