//go:build windows

package base

// Real-window slice of the legacy LineWrapper (moved from
// internal/richtextspec's TestRealRichText "LineWrapper over the real
// font" subtest in browser-core ticket 75): the greedy wrapper over the
// native shaping seam, where per-character advances come from the real
// font. The probe runs inside a prepaint pass of a real Win32 window
// (the inRealPrepaint pattern); failures fail, never skip.

import (
	"strings"
	"testing"

	"github.com/hoangvu12/gpui-go/gpui"
)

// TestRealLineWrapperWrapsWordsInsideTheWidth pins the wrapper's
// contract on a real font: every boundary starts a word after the
// previous one, and each row's content (without its trailing break
// space) fits the wrap width.
func TestRealLineWrapperWrapsWordsInsideTheWidth(t *testing.T) {
	// The probe runs on the host thread, so it reports with Errorf only.
	inRealPrepaint(t, func(w *gpui.Window, app *gpui.App) {
		const words = "alpha bravo charlie delta echo foxtrot golf"
		font := gpui.DefaultTextStyle().FontDescriptorOf()
		tsys := w.TextSystem()
		wrapper := NewLineWrapper(w, font, 16)
		const width = 90
		boundaries := wrapper.WrapLine([]LineFragment{LineFragmentText(words)}, width)
		if len(boundaries) < 2 {
			t.Errorf("boundaries = %v, want the text wrapped at %vpx", boundaries, width)
			return
		}
		start := 0
		for _, boundary := range append(boundaries, Boundary{Ix: len(words)}) {
			if boundary.Ix <= start || (boundary.Ix < len(words) && words[boundary.Ix-1] != ' ') {
				t.Errorf("boundary %d does not start a word after %d (%v)", boundary.Ix, start, boundaries)
				return
			}
			// The row's content (without its trailing break space) fits.
			row := strings.TrimRight(words[start:boundary.Ix], " ")
			var rowWidth float32
			for _, ch := range row {
				rowWidth += tsys.LayoutWidth(font, 16, ch)
			}
			if rowWidth > width+0.01 {
				t.Errorf("row %q measures %v, wider than %v", row, rowWidth, width)
			}
			start = boundary.Ix
		}
	})
}
