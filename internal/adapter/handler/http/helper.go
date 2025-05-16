package http

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any) map[string]any {
	return map[string]any{
		"meta": m,
		"data": data,
	}
}
