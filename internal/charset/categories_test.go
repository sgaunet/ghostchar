package charset

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		wantOK   bool
		wantCat  Category
		wantName string
	}{
		{"ZWSP", 0x200B, true, Invisible, "ZERO WIDTH SPACE"},
		{"ZWNJ", 0x200C, true, Invisible, "ZERO WIDTH NON-JOINER"},
		{"ZWJ", 0x200D, true, Invisible, "ZERO WIDTH JOINER"},
		{"soft hyphen", 0x00AD, true, Invisible, "SOFT HYPHEN"},
		{"BOM/ZWNBSP", 0xFEFF, true, Invisible, "ZERO WIDTH NO-BREAK SPACE"},
		{"word joiner", 0x2060, true, Invisible, "WORD JOINER"},
		{"combining grapheme joiner", 0x034F, true, Invisible, "COMBINING GRAPHEME JOINER"},
		{"RLO", 0x202E, true, Bidi, "RIGHT-TO-LEFT OVERRIDE"},
		{"PDF", 0x202C, true, Bidi, "POP DIRECTIONAL FORMATTING"},
		{"LRM", 0x200E, true, Bidi, "LEFT-TO-RIGHT MARK"},
		{"RLM", 0x200F, true, Bidi, "RIGHT-TO-LEFT MARK"},
		{"FSI", 0x2068, true, Bidi, "FIRST STRONG ISOLATE"},
		{"PUA start", 0xE000, true, PUA, "PRIVATE USE AREA"},
		{"PUA mid", 0xF000, true, PUA, "PRIVATE USE AREA"},
		{"PUA end", 0xF8FF, true, PUA, "PRIVATE USE AREA"},
		{"SPUA-A", 0xF0000, true, PUA, "SUPPLEMENTARY PRIVATE USE AREA-A"},
		{"SPUA-B", 0x100000, true, PUA, "SUPPLEMENTARY PRIVATE USE AREA-B"},
		{"normal letter", 'A', false, "", ""},
		{"normal space", ' ', false, "", ""},
		{"emoji", 0x1F600, false, "", ""},
		{"newline", '\n', false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, ok := Detect(tt.r)
			if ok != tt.wantOK {
				t.Fatalf("Detect(%#U): got ok=%v, want %v", tt.r, ok, tt.wantOK)
			}
			if ok {
				if cd.Category != tt.wantCat {
					t.Errorf("Detect(%#U): got category=%q, want %q", tt.r, cd.Category, tt.wantCat)
				}
				if cd.Name != tt.wantName {
					t.Errorf("Detect(%#U): got name=%q, want %q", tt.r, cd.Name, tt.wantName)
				}
			}
		})
	}
}

func TestDetectInCategory(t *testing.T) {
	// ZWSP should be found when looking for invisible
	cd, ok := DetectInCategory(0x200B, []Category{Invisible})
	if !ok {
		t.Fatal("expected ZWSP to match invisible category")
	}
	if cd.Category != Invisible {
		t.Errorf("expected invisible, got %s", cd.Category)
	}

	// ZWSP should NOT be found when looking for bidi only
	_, ok = DetectInCategory(0x200B, []Category{Bidi})
	if ok {
		t.Error("ZWSP should not match bidi category")
	}

	// RLO should match bidi
	cd, ok = DetectInCategory(0x202E, []Category{Bidi})
	if !ok {
		t.Fatal("expected RLO to match bidi category")
	}
	if cd.Category != Bidi {
		t.Errorf("expected bidi, got %s", cd.Category)
	}

	// PUA should match pua
	cd, ok = DetectInCategory(0xE000, []Category{PUA})
	if !ok {
		t.Fatal("expected PUA start to match pua category")
	}
	if cd.Category != PUA {
		t.Errorf("expected pua, got %s", cd.Category)
	}

	// Multiple categories
	cd, ok = DetectInCategory(0x200B, []Category{Invisible, Bidi})
	if !ok {
		t.Fatal("expected ZWSP to match when invisible is in category list")
	}
	if cd.Category != Invisible {
		t.Errorf("expected invisible, got %s", cd.Category)
	}
}

func TestAllCharDefs(t *testing.T) {
	defs := AllCharDefs()
	if len(defs) == 0 {
		t.Fatal("AllCharDefs returned empty list")
	}

	// Verify we have all three categories
	cats := make(map[Category]int)
	for _, d := range defs {
		cats[d.Category]++
	}
	if cats[Invisible] != 11 {
		t.Errorf("expected 11 invisible chars, got %d", cats[Invisible])
	}
	if cats[Bidi] != 12 {
		t.Errorf("expected 12 bidi chars, got %d", cats[Bidi])
	}
	if cats[PUA] != 3 {
		t.Errorf("expected 3 PUA ranges, got %d", cats[PUA])
	}
}
