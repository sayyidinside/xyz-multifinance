package helpers

type contextKey string

const (
	ctxKeyIdentifier contextKey = "identifier"
	ctxKeyUsername   contextKey = "username"
	ctxKeyFunction   contextKey = "function"
)
