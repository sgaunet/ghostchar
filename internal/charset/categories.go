package charset

// Category represents a type of detectable Unicode character.
type Category string

const (
	Invisible Category = "invisible"
	PUA       Category = "pua"
	Bidi      Category = "bidi"
)

// AllCategories returns all valid category values.
func AllCategories() []Category {
	return []Category{Invisible, PUA, Bidi}
}

// CharDef defines a detectable character or range.
type CharDef struct {
	Low      rune
	High     rune
	Name     string
	Category Category
}

// IsSingle returns true if this definition covers a single codepoint.
func (c CharDef) IsSingle() bool {
	return c.Low == c.High
}

// Contains returns true if the rune falls within this definition's range.
func (c CharDef) Contains(r rune) bool {
	return r >= c.Low && r <= c.High
}

// charDefs is the master list of all detectable character definitions.
var charDefs []CharDef

func init() {
	charDefs = make([]CharDef, 0, len(invisibleChars)+len(bidiChars)+len(puaRanges))
	charDefs = append(charDefs, invisibleChars...)
	charDefs = append(charDefs, bidiChars...)
	charDefs = append(charDefs, puaRanges...)
}

// AllCharDefs returns the complete list of detectable character definitions.
func AllCharDefs() []CharDef {
	return charDefs
}

// CharDefsByCategory returns character definitions filtered by category.
func CharDefsByCategory(cat Category) []CharDef {
	var result []CharDef
	for _, cd := range charDefs {
		if cd.Category == cat {
			result = append(result, cd)
		}
	}
	return result
}

// invisibleChars defines zero-width and invisible characters.
var invisibleChars = []CharDef{
	{Low: 0x00AD, High: 0x00AD, Name: "SOFT HYPHEN", Category: Invisible},
	{Low: 0x034F, High: 0x034F, Name: "COMBINING GRAPHEME JOINER", Category: Invisible},
	{Low: 0x200B, High: 0x200B, Name: "ZERO WIDTH SPACE", Category: Invisible},
	{Low: 0x200C, High: 0x200C, Name: "ZERO WIDTH NON-JOINER", Category: Invisible},
	{Low: 0x200D, High: 0x200D, Name: "ZERO WIDTH JOINER", Category: Invisible},
	{Low: 0x2060, High: 0x2060, Name: "WORD JOINER", Category: Invisible},
	{Low: 0x2061, High: 0x2061, Name: "FUNCTION APPLICATION", Category: Invisible},
	{Low: 0x2062, High: 0x2062, Name: "INVISIBLE TIMES", Category: Invisible},
	{Low: 0x2063, High: 0x2063, Name: "INVISIBLE SEPARATOR", Category: Invisible},
	{Low: 0x2064, High: 0x2064, Name: "INVISIBLE PLUS", Category: Invisible},
	{Low: 0xFEFF, High: 0xFEFF, Name: "ZERO WIDTH NO-BREAK SPACE", Category: Invisible},
}

// bidiChars defines bidirectional control characters.
var bidiChars = []CharDef{
	{Low: 0x200E, High: 0x200E, Name: "LEFT-TO-RIGHT MARK", Category: Bidi},
	{Low: 0x200F, High: 0x200F, Name: "RIGHT-TO-LEFT MARK", Category: Bidi},
	{Low: 0x202A, High: 0x202A, Name: "LEFT-TO-RIGHT EMBEDDING", Category: Bidi},
	{Low: 0x202B, High: 0x202B, Name: "RIGHT-TO-LEFT EMBEDDING", Category: Bidi},
	{Low: 0x202C, High: 0x202C, Name: "POP DIRECTIONAL FORMATTING", Category: Bidi},
	{Low: 0x202D, High: 0x202D, Name: "LEFT-TO-RIGHT OVERRIDE", Category: Bidi},
	{Low: 0x202E, High: 0x202E, Name: "RIGHT-TO-LEFT OVERRIDE", Category: Bidi},
	{Low: 0x2066, High: 0x2066, Name: "LEFT-TO-RIGHT ISOLATE", Category: Bidi},
	{Low: 0x2067, High: 0x2067, Name: "RIGHT-TO-LEFT ISOLATE", Category: Bidi},
	{Low: 0x2068, High: 0x2068, Name: "FIRST STRONG ISOLATE", Category: Bidi},
	{Low: 0x2069, High: 0x2069, Name: "POP DIRECTIONAL ISOLATE", Category: Bidi},
	{Low: 0x061C, High: 0x061C, Name: "ARABIC LETTER MARK", Category: Bidi},
}

// puaRanges defines Private Use Area ranges.
var puaRanges = []CharDef{
	{Low: 0xE000, High: 0xF8FF, Name: "PRIVATE USE AREA", Category: PUA},
	{Low: 0xF0000, High: 0xFFFFF, Name: "SUPPLEMENTARY PRIVATE USE AREA-A", Category: PUA},
	{Low: 0x100000, High: 0x10FFFF, Name: "SUPPLEMENTARY PRIVATE USE AREA-B", Category: PUA},
}
