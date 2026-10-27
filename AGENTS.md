# gpui-go-base

The `base` package (module root) is the Go port of GPUI Base (gpui-component
`crates/base` at `c0820693`) over gpui-go's public API. It was split out of
gpui-go in browser-core ticket 75; the ticket history and the browser-parity
spec stay in the gpui-go repository.

## Ground rules

- The test apps (zterm, the gallery, the showcase) live only in the local
  development checkout and are never pushed. Do not add remotes, tags or
  pushes without the user's explicit go-ahead.
- The module builds against a local gpui-go checkout through the go.mod
  `replace` (`../gpui-go-wt-base`, LOCAL-ONLY; see go.mod's comment). Before
  any publish, drop the replace and require a tagged gpui-go version.
- Windows AMD64 is the only target (as gpui-go).
- Headless testing only: never set GPUI_TEST_ONSCREEN=1 or
  GPUI_REAL_INPUT_TESTS=1, and never drive real input. The pre-existing
  real-window tests (the `_real_windows_test.go` files) follow their own
  established pattern.

## Layout

- Module root: package `base` (one Go file per upstream Rust module; import
  path `github.com/hoangvu12/gpui-go-base`).
- Subpackages: `rope`, `sumtree`, `segment`, `dataurl`, `strftime`;
  `testdata` holds upstream fixtures.
- `conformance/cmd/basecoverage`: the coverage ledger generator;
  `go run ./conformance/cmd/basecoverage -check` must stay current. It scans
  this repository plus the gpui-go checkout's `gpui/` tests (resolved from
  the go.mod replace, `-gpui-root` or $GPUI_GO_ROOT override).
- `evidence/base/`: the coverage ledger, the hashed upstream source manifest
  and the per-ticket evidence.

## Working on it

- Run go commands with the same shared-PC discipline as gpui-go
  (`GOFLAGS=-p=4`, targeted tests while iterating, full suite once at the
  end).
- Upstream citations use `U:path:line` (relative to `crates/base/src`); the
  coverage ledger is regenerated, never hand-edited.
- The upstream checkout is expected at `../upstream-base/ce` (override with
  -upstream or GPUI_BASE_UPSTREAM).
