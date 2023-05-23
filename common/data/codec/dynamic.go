package codec

type dynamicConverter[F, T any] struct {
	f func(F) (T, error)
}

func (d dynamicConverter[F, T]) Convert(from F) (T, error) {
	return d.f(from)
}

func newDynamicConverter[F, T any](f func(F) (T, error)) dynamicConverter[F, T] {
	return dynamicConverter[F, T]{f: f}
}

type dynamicCodec[F, T any] struct {
	c1 Converter[F, T]
	c2 Converter[T, F]
}

func (d dynamicCodec[F, T]) Encoder() Converter[F, T] {
	return d.c1
}

func (d dynamicCodec[F, T]) Decoder() Converter[T, F] {
	return d.c2
}

func New[From, To any](forward func(From) (To, error), backward func(To) (From, error)) Codec[From, To] {
	return dynamicCodec[From, To]{
		c1: newDynamicConverter[From, To](forward),
		c2: newDynamicConverter[To, From](backward),
	}
}
