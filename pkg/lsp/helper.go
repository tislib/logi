package lsp

func pointer[T any](severityError T) *T {
	return &severityError
}
