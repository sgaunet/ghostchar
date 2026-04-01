package charset

// Detect checks if a rune matches any detectable character definition.
// Returns the matching CharDef and true if found, zero value and false otherwise.
func Detect(r rune) (CharDef, bool) {
	for _, cd := range charDefs {
		if cd.Contains(r) {
			return cd, true
		}
	}
	return CharDef{}, false
}

// DetectInCategory checks if a rune matches any character definition in the given categories.
func DetectInCategory(r rune, categories []Category) (CharDef, bool) {
	catSet := make(map[Category]bool, len(categories))
	for _, c := range categories {
		catSet[c] = true
	}
	for _, cd := range charDefs {
		if catSet[cd.Category] && cd.Contains(r) {
			return cd, true
		}
	}
	return CharDef{}, false
}
