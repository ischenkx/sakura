package data

import "context"

type Item interface {
	ID() string
}

type Storage[I Item, Selector any] interface {
	Insert(ctx context.Context, item I) error
	Select(ctx context.Context, selector Selector) (Set[I], error)
}

type Set[I Item] interface {
	Erase(ctx context.Context) error
	Iter(ctx context.Context, iter func(I) bool)
}
