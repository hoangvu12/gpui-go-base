module github.com/hoangvu12/gpui-go-base

go 1.27.0

require (
	github.com/clipperhouse/uax29/v2 v2.7.0
	github.com/hoangvu12/gpui-go v0.2.0
	github.com/yuin/goldmark v1.8.6
	golang.org/x/net v0.61.0
)

// LOCAL-ONLY replace: point gpui-go at a local checkout while the two
// repositories are developed side by side. Before any publish, remove
// this replace and require a tagged gpui-go version instead (go mod
// tidy). ../gpui-go-wt-base is the orchestrator's main worktree of
// gpui-go; after browser-core ticket 75 merged there, its main branch
// carries no Base (Base lives here).
replace github.com/hoangvu12/gpui-go => ../gpui-go-wt-base
