# gpui-go-base

gpui-go-base is the Go port of GPUI Base, the `crates/base` crate of
[gpui-component](https://github.com/gpui-ce/gpui-component) at commit
`c08206932417f863062d2ab70cc2854c6ca04238` (Apache-2.0). It is the component
library that sits on top of [gpui-go](https://github.com/hoangvu12/gpui-go),
the Go port of the GPUI-CE core: it depends only on gpui-go's public API.

Base components are unstyled: they own interaction, focus, keyboard handling
and accessibility, and you supply every visual through gpui-go's Styled API.
Call `base.Init(app)` once and make each window's root a `base.Root` around
your view, so overlays and text selection get their window layers. The
package is the module root, so the import is just:

```go
import "github.com/hoangvu12/gpui-go-base"
```

What is ported, one Go file per upstream module:

- Infrastructure: theme tokens, `Root` and its plugins, global state, history and undo, motion (easings, springs, keyframes, presence).
- Controls: button, toggle and toggle group, checkbox, radio, switch, link, tabs, toolbar, accordion, collapsible, avatar, progress, slider, pagination, table.
- Overlays: popup, popover, tooltip, hover card, dialog, alert dialog, sheet, toast — built on the core's top layer, close watchers, `ShowModal`, `Anchored` and the tooltip hook.
- Scrolling and lists: scrollbar, scroll bounce, scrollable mask, auto scroll, virtual list.
- Input and editor: single-line input, textarea, number, OTP and time fields, select, combobox, calendar, date and color pickers, and an editor mode with folding, search, diagnostics and LSP hooks. Soft wrapping runs the legacy greedy `LineWrapper` (this repository's `line_wrapper.go`).
- Composites: tree, navigation stack, resizable panels, questionnaire, and a dock with persisted layouts.
- Text: Markdown and HTML into a node tree, `TextView`, cross-element text selection, stream fade and range highlights.
- Plot: scales, bar, line, area, arc, pie, stack and sankey shapes, axes, grids and hover.
- Data structures in subpackages: `rope` (ropey subset), `sumtree` (zed-sum-tree), `segment` (word/grapheme segmentation over uax29), `dataurl`, `strftime`.

[evidence/base/COVERAGE.md](evidence/base/COVERAGE.md) accounts for every
upstream module and test. Each test is listed as ported, deviated, not
applicable or deferred, with the reason. `go run ./conformance/cmd/basecoverage`
regenerates the ledger from the hashed source manifest and the `U:path:line`
citations in the Go tests (it also follows the tests that moved into the
gpui-go core, through the local replace target). Its `-check` flag fails when
the ledger is stale or reports a gap.

## Building against the local gpui-go

Until gpui-go publishes a tagged release with the post-split core, go.mod
carries a LOCAL-ONLY replace pointing at a sibling checkout:

```
replace github.com/hoangvu12/gpui-go => ../gpui-go-wt-base
```

Clone gpui-go next to this repository (or edit the path), then:

```sh
go build ./...
go test ./...
```

Before any publish, remove the replace and require a tagged gpui-go version
instead (`go mod tidy`). Build with `CGO_ENABLED=0` like gpui-go; the native
DLL comes embedded inside the gpui-go module.

## License

gpui-go-base is licensed under the [Apache License 2.0](LICENSE), the same
license as gpui-component's crates/base. [NOTICE](NOTICE) and
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) credit the ported sources
(uax29, goldmark, golang.org/x/net and the ported crates).
