package types

type Factory[T any] func() (T, error)
