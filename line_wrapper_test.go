package base

// The legacy LineWrapper's own tests, moved with the wrapper from
// gpui-go in browser-core ticket 75 (they were gpui/rich_text_test.go's
// TestIsWordChar, gpui/text_tabs_test.go's
// TestLineWrapperTabFragmentTakesItsWidth and internal/richtextspec's
// TestWrapLine / TestLineWrapperLetterSpacing).
//
// This file is the deterministic slice: test windows drawn through the
// real element pipeline with the deterministic test text system (the
// pinned TestTextSystem port: em width = font size * 0.6 per UTF-16
// unit), whose .ZedMono advance of 600/1000 em is exactly the legacy
// em width, so the boundaries are the legacy values. The real-font
// slice lives in line_wrapper_real_windows_test.go.
//
// Upstream sources (legacy = ~/.cargo/git/checkouts/
// gpui-ce-6758b1db058c6497/beb5a6d/crates/gpui/src):
//
//   - legacy text_system/line_wrapper.rs:836-982 test_wrap_line;
//   - legacy line_wrapper.rs:96-104, 135-143 (letter-spacing
//     arithmetic) and :1241-1320 test_is_word_char.

import (
	"math"
	"testing"

	"github.com/hoangvu12/gpui-go/gpui"
)

// newWrapperWindow opens a deterministic test window with an empty root
// (the richtextspec newTestWindow pattern).
func newWrapperWindow(t *testing.T) *gpui.Window {
	t.Helper()
	ta := gpui.NewTestApp()
	window := gpui.NewTestWindow(ta.App(), gpui.Size{Width: 400, Height: 200})
	t.Cleanup(func() { gpui.CloseTestWindow(window) })
	return window
}

// approxTab compares two floats within the text-tabs tolerance.
func approxTab(a, b float32) bool { return math.Abs(float64(a-b)) < 0.02 }

// TestWrapLine is the port of legacy line_wrapper.rs:836-982
// test_wrap_line.
func TestWrapLine(t *testing.T) {
	window := newWrapperWindow(t)
	wrapper := NewLineWrapper(window, gpui.FontDescriptor{Family: ".ZedMono", Weight: 400}, 16)
	text := LineFragmentText
	element := LineFragmentElement
	b := func(ix int, indent uint32) Boundary { return Boundary{Ix: ix, NextIndent: indent} }
	cases := []struct {
		name      string
		fragments []LineFragment
		want      []Boundary
	}{
		{"words", []LineFragment{text("aa bbb cccc ddddd eeee")}, []Boundary{b(7, 0), b(12, 0), b(18, 0)}},
		{"long word", []LineFragment{text("aaa aaaaaaaaaaaaaaaaaa")}, []Boundary{b(4, 0), b(11, 0), b(18, 0)}},
		{"indented", []LineFragment{text("     aaaaaaa")}, []Boundary{b(7, 5), b(9, 5), b(11, 5)}},
		{"spaces", []LineFragment{text("                            ")}, []Boundary{b(7, 0), b(14, 0), b(21, 0)}},
		{"deep indent", []LineFragment{text("          aaaaaaaaaaaaaa")}, []Boundary{b(7, 0), b(14, 3), b(18, 3), b(22, 3)}},
		{"multiple text fragments", []LineFragment{text("aa bbb "), text("cccc ddddd eeee")}, []Boundary{b(7, 0), b(12, 0), b(18, 0)}},
		{"text and elements", []LineFragment{text("aa "), element(20, 1), text(" bbb "), element(30, 1), text(" cccc")}, []Boundary{b(5, 0), b(9, 0), b(11, 0)}},
		{"leading element", []LineFragment{element(50, 1), text(" aaaa bbbb cccc dddd")}, []Boundary{b(2, 0), b(7, 0), b(12, 0), b(17, 0)}},
		{"large element", []LineFragment{text("short text "), element(100, 1), text(" more text")}, []Boundary{b(6, 0), b(11, 0), b(12, 0), b(18, 0)}},
		{"glue characters", []LineFragment{text("a b c‑d e")}, []Boundary{b(12, 0)}},
	}
	for _, c := range cases {
		got := wrapper.WrapLine(c.fragments, 72)
		if len(got) != len(c.want) {
			t.Errorf("%s: boundaries = %v, want %v", c.name, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: boundaries = %v, want %v", c.name, got, c.want)
				break
			}
		}
	}
}

// TestLineWrapperLetterSpacing pins the wrapper's tracking arithmetic
// (legacy line_wrapper.rs:96-104, 135-143): letter spacing widens every
// character after the first, so the same text wraps earlier.
func TestLineWrapperLetterSpacing(t *testing.T) {
	window := newWrapperWindow(t)
	wrapper := NewLineWrapper(window, gpui.FontDescriptor{Family: ".ZedMono", Weight: 400}, 16)
	fragments := []LineFragment{LineFragmentText("aaaaaaaaaaaa")}
	// 9.6 per char: 7 chars fit 72 (67.2), the 8th overflows.
	if got := wrapper.WrapLine(fragments, 72); len(got) == 0 || got[0].Ix != 7 {
		t.Fatalf("untracked boundaries = %v, want the first at 7", got)
	}
	// +2 per following char: 9.6 + 6*11.6 = 79.2 > 72, so 6 chars fit.
	wrapper.SetLetterSpacing(ptrOf(float32(2)))
	if got := wrapper.WrapLine(fragments, 72); len(got) == 0 || got[0].Ix != 6 {
		t.Fatalf("tracked boundaries = %v, want the first at 6", got)
	}
	wrapper.SetLetterSpacing(nil)
	if got := window.TextSystem().LayoutWidth(gpui.FontDescriptor{Family: ".ZedMono"}, 16, 'm'); !near(got, 9.6) {
		t.Fatalf("LayoutWidth('m') = %v, want the em width 9.6", got)
	}
}

// TestLineWrapperTabFragmentTakesItsWidth: a resolved tab fragment
// wraps like the '\t' text character at its given width (Base's input
// hands the wrapper its column widths). Moved from
// gpui/text_tabs_test.go with the wrapper.
func TestLineWrapperTabFragmentTakesItsWidth(t *testing.T) {
	window := newWrapperWindow(t)
	wrapper := NewLineWrapper(window, gpui.DefaultFontDescriptor(), 10)
	// The test text system's advance is 6px per character at 10px.
	text := []LineFragment{LineFragmentText("ab"), LineFragmentText("\t"), LineFragmentText("cd")}
	if got := wrapper.WrapLine(text, 30); len(got) != 0 {
		t.Errorf("text tab (6px): boundaries %v, want none at 30px", got)
	}
	tab := []LineFragment{LineFragmentText("ab"), LineFragmentTab(30), LineFragmentText("cd")}
	got := wrapper.WrapLine(tab, 30)
	// Like the text '\t', the tab is a break candidate (a non-word,
	// non-space character): the row breaks before it.
	if want := []Boundary{{Ix: 2}, {Ix: 3}}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("resolved tab (30px): boundaries %v, want %v", got, want)
	}
}

// TestIsWordChar is the port of the legacy line_wrapper.rs:1241-1320
// test_is_word_char. Moved from gpui/rich_text_test.go with the
// wrapper's word-character table.
func TestIsWordChar(t *testing.T) {
	assertWord := func(word string) {
		t.Helper()
		for _, c := range word {
			if !isWordChar(c) {
				t.Errorf("isWordChar(%q U+%04X) = false, want true (word %q)", c, c, word)
			}
		}
	}
	assertNotWord := func(word string) {
		t.Helper()
		found := false
		for _, c := range word {
			if !isWordChar(c) {
				found = true
			}
		}
		if !found {
			t.Errorf("every character of %q is a word character, want at least one that is not", word)
		}
	}

	assertWord("Hello123")
	assertWord("non-English")
	assertWord("var_name")
	assertWord("123456")
	assertWord("3.1415")
	assertWord("10^2")
	assertWord("1~2")
	assertWord("100%")
	assertWord("@mention")
	assertWord("#hashtag")
	assertWord("$variable")
	assertWord("a=1")
	assertWord("Self::is_word_char")
	assertWord("on;")
	assertWord("more⋯")
	assertWord("won’t")
	assertWord("‘twas")
	assertWord("plz!")
	assertWord("see)")
	assertWord("quoted”")
	assertWord("well…")

	// Space.
	assertNotWord("foo bar")

	// URL case.
	assertWord("github.com")
	assertNotWord("zed-industries/zed")
	assertNotWord("zed-industries\\zed")
	assertNotWord("a=1&b=2")
	assertNotWord("foo?b=2")

	// Latin-1 Supplement.
	assertWord("ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ")
	// Latin Extended-A.
	assertWord("ĀāĂăĄąĆćĈĉĊċČčĎď")
	// Latin Extended-B.
	assertWord("ƀƁƂƃƄƅƆƇƈƉƊƋƌƍƎƏ")
	// Cyrillic.
	assertWord("АБВГДЕЖЗИЙКЛМНОП")
	// Vietnamese.
	assertWord("ThậmchíđếnkhithuachạychúngcònnhẫntâmgiếtnốtsốđôngtùchínhtrịởYênBáivàCaoBằng")
	// Bengali.
	assertWord("গিয়েছিলেন")
	assertWord("ছেলে")
	assertWord("হচ্ছিল")

	// Non-word characters.
	assertNotWord("你好")
	assertNotWord("안녕하세요")
	assertNotWord("こんにちは")
	assertNotWord("😀😁😂")
	assertNotWord("()[]{}<>")

	// Non-breaking glue characters (UAX #14).
	assertWord(" ")
	assertWord(" ")
	assertWord("‑")
}
