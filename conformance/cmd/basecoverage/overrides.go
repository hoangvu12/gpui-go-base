package main

// override resolves an upstream test that no Go file cites by line
// or name (or, with Force, records a deviation the ticket evidence
// names but the citing Go comment does not spell out). Every entry
// names its ticket and paraphrases that ticket's evidence JSON or the
// porting Go file; the tool lists entries that match nothing as unused.
type override struct {
	Status string
	Ticket string
	Reason string
	Force  bool
}

const (
	openedURL = "not applicable: asserts cx.opened_url(), which gpui-go's TestApp does not record (ticket28 evidence text_view.rs; text_view_upstream_test.go header)"
)

// overrides is keyed "module::fn" (module relative to crates/base/src).
var overrides = map[string]override{
	// Deleted with the browser-core migrations onto core primitives:
	// the positioner tests' module (base/positioner.go) went with the
	// core's element-anchored gpui.Anchored, the tooltip provider with
	// the core tooltip hook, the reduce-motion reader with the core's
	// App.MotionReduced gate. The tickets' Answers (gpui-go
	// .scratch/browser-core/issues/) record each deletion.
	"positioner.rs::prefers_the_requested_side_when_it_fits":                          {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — placement resolves through the core's element-anchored gpui.Anchored (side/alignment, gap, window margin, AnchorShift fallback; the corner rule clamps without flipping); the placed results the popup still owns are pinned by popup_test.go (U:positioner.rs:95-109, 303-325) and the core anchor fixtures", false},
	"positioner.rs::flips_to_the_opposite_side_when_the_preferred_side_does_not_fit":    {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — placement resolves through the core's element-anchored gpui.Anchored; the corner rule clamps without flipping, and the core anchor fixtures cover the placement space", false},
	"positioner.rs::clamps_into_the_viewport_while_keeping_the_flipped_side":            {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — the window margin clamp lives in the core's gpui.Anchored placement", false},
	"positioner.rs::alignment_selects_the_leading_center_or_trailing_edge":             {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — the alignment mapping lives in the popup's Anchored placement", false},
	"positioner.rs::side_offset_adds_a_gap_between_trigger_and_popup":                  {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — the offset maps onto the anchor gap in the core's gpui.Anchored", false},
	"positioner.rs::corner_positioning_places_the_named_corner_and_never_reports_a_side": {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — the corner anchors map onto the core's Anchored corner rule (clamps, never flips); popup_test.go pins the placed results (U:positioner.rs:303-325)", false},
	"positioner.rs::corner_positioning_clamps_but_does_not_flip":                       {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — the corner clamp is the core Anchored rule; popup_test.go pins it (U:positioner.rs:95-109)", false},
	"positioner.rs::corner_position_can_be_updated_without_rebuilding_the_positioner": {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — placement is a render-time core Anchored call, so there is no positioner object to rebuild", false},
	"positioner.rs::frame_insets_apply_the_client_inset_only_on_untiled_edges":         {statusNA, "72", "not ported: windowDecorations/frameInsets/tiling were deleted with base/positioner.go in browser-core ticket 72 — the core placement takes the window margin only", false},
	"positioner.rs::clamping_keeps_only_the_margin_on_a_tiled_edge":                    {statusNA, "72", "not ported: windowDecorations/frameInsets/tiling were deleted with base/positioner.go in browser-core ticket 72 — the core placement takes the window margin only", false},
	"positioner.rs::prefers_above_when_space_allows":                                   {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — the side preference lives in the popup's Anchored placement", false},
	"positioner.rs::flips_and_clamps_on_each_axis":                                     {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — placement resolves through the core's gpui.Anchored; the core anchor fixtures cover the placement space", false},
	"positioner.rs::places_tooltip_to_the_right":                                       {statusNA, "72", "not ported: the tooltip placement this test drove belongs to the deleted positioner (browser-core ticket 72); the core tooltip hook places the hint layer (ticket 24)", false},
	"positioner.rs::right_placement_clamps_vertical_edges":                             {statusNA, "72", "not ported: the tooltip placement this test drove belongs to the deleted positioner (browser-core ticket 72); the core tooltip hook places the hint layer (ticket 24)", false},
	"positioner.rs::uses_larger_side_when_neither_vertical_side_fits":                   {statusNA, "72", "not ported: base/positioner.go was deleted in browser-core ticket 72 — placement resolves through the core's gpui.Anchored", false},
	"tooltip.rs::provider_owns_grace_switch_and_dismiss":                               {statusNA, "72", "not ported: the per-window TooltipOverlay provider was deleted in browser-core ticket 72 — the lifecycle is the core tooltip hook (DivElement.TooltipView as a Tooltip surface, ticket 24) with Chrome's interest delays and grace", false},
	"tooltip.rs::tooltip_priority_exceeds_popup_layer":                                 {statusNA, "72", "not ported: the per-window TooltipOverlay provider was deleted in browser-core ticket 72 — the hint layer's paint order is the core TopLayer's", false},
	"tooltip.rs::disabled_provider_ignores_delayed_and_immediate_requests":             {statusNA, "72", "not ported: the per-window TooltipOverlay provider was deleted in browser-core ticket 72 — there is no provider to disable; the core tooltip hook owns the delays", false},
	"reduce_motion.rs::a_system_preference_for_reduced_motion_sets_the_flag":           {statusNA, "74", "not ported: the one-shot system-preference reader (base/reduce_motion*.go) was deleted in browser-core ticket 74 — the single gate is the core's App.MotionReduced (App.ReduceMotion || the live SystemPreferences record, WM_SETTINGCHANGE); motion_reduced_test.go pins the gate", false},
	"reduce_motion.rs::an_unknown_system_preference_leaves_the_flag_alone":             {statusNA, "74", "not ported: the one-shot reader went with base/reduce_motion*.go in browser-core ticket 74; the live preference record replaces the parsed string", false},
	"reduce_motion.rs::the_system_drives_the_flag_until_the_application_sets_it":       {statusNA, "74", "not ported: the one-shot reader went with base/reduce_motion*.go in browser-core ticket 74 — the OS preference and the App flag now OR (the core's policy since ticket 44), so the application no longer un-reduces a system-reduced machine", false},
	"reduce_motion.rs::the_test_scheduler_is_never_asked_for_the_platform_preference":  {statusNA, "74", "not ported: the reader went with base/reduce_motion*.go in browser-core ticket 74 — no code asks the scheduler for the preference; the host caches per preference version (tickets 27/29)", false},
	"reduce_motion.rs::a_flag_the_application_set_before_the_first_reading_is_kept":    {statusNA, "74", "not ported: the reader went with base/reduce_motion*.go in browser-core ticket 74 — SetReduceMotion writes the App flag the core's MotionReduced ORs in", false},

	// Ported, but the Go tests cite stale line numbers
	// (U:text/selection.rs:198-222 / 224-253 in a 172-line file).
	"text/selection.rs::test_char_type_from_char": {statusPorted, "29", "ported as TestTextCharTypeFromChar (text_selection_word_test.go; its U: line citation is stale); ticket29 evidence text/selection.rs 2/2", false},
	"text/selection.rs::test_word_range_at":       {statusPorted, "29", "ported as TestTextSelectionWordRangeAt (text_selection_word_test.go; its U: line citation is stale); ticket29 evidence text/selection.rs 2/2", false},

	// crates/base/tests/element_ext.rs: the property is checked by a
	// Go test marked Go-own that cites the src module instead.
	"../tests/element_ext.rs::prepaint_callback_observes_the_parent_bounds": {statusDeviated, "13", "covered by TestElementExtOnPrepaintReportsBounds (window_test.go, cites U:element_ext.rs:15-28): the prepaint callback reports the element's bounds (a 40x30 child at y=10 instead of an 80x80 root)", false},

	// Not applicable, per the ticket evidence.
	"input/base/state.rs::test_input_does_not_invalidate_cached_parent_during_paint": {statusNA, "20", "not applicable: the cached-parent test needs AnyView::cached, absent in the Go core (ticket20 evidence not_applicable_reasons)", false},
	"text/inline.rs::shortening_an_inline_flow_releases_obsolete_fragment_layouts":   {statusNA, "28", "not applicable: the Go retained-layout table is weak-keyed, so a dropped fragment's entry lives until the next sweep; the live-entry count is ported in inline_flow_fragments_keep_their_layouts_across_frames (ticket28 evidence inline.rs records it deferred; ticket 30 did not take it up; text_inline_test.go header)", false},
	"text/text_view.rs::list_item_renders_fenced_code_block_at_document_width":       {statusNA, "28", "not applicable: needs h_resizable/resizable_panel inside a TextView list item, not ported (ticket28 evidence text_view.rs '1 h_resizable not ported')", false},
	"text/text_view.rs::clipped_markdown_link_does_not_open":                         {statusNA, "28", openedURL, false},
	"text/text_view.rs::right_click_does_not_open_url_without_handler":               {statusNA, "28", openedURL, false},

	// Deferred to the open ticket 37 (request_autoscroll).

	// Deviations the evidence names that the citing Go comment does not.
	"global_state.rs::initialization_is_idempotent_and_suppression_can_be_reset":  {statusDeviated, "13", "adapted: drop(token) becomes Release() (ticket13 evidence upstream_ported_detail)", true},
	"global_state.rs::a_dropped_registration_closes_the_deferred_context":         {statusDeviated, "13", "adapted: drop(token) becomes Release() (ticket13 evidence upstream_ported_detail)", true},
	"observe.rs::observation_preserves_identity_and_native_type_in_normal_builds": {statusDeviated, "13", "adapted: the normal-build branch only (no test-support feature) (ticket13 evidence upstream_ported_detail)", true},
	"motion.rs::status_transition_reports_delay_running_and_finished":             {statusDeviated, "14", "passes only with ticket14 deviation 1 (f64 duration scaling) (ticket14 evidence deviated_list)", true},
	"motion.rs::a_completed_transition_stops_requesting_frames":                   {statusDeviated, "14", "passes only with ticket14 deviation 1 (f64 duration scaling) (ticket14 evidence deviated_list)", true},
	"motion.rs::presence_enters_exits_and_only_unmounts_after_exit":               {statusDeviated, "14", "passes only with ticket14 deviation 1 (f64 duration scaling) (ticket14 evidence deviated_list)", true},
	"avatar.rs::image_slot_takes_precedence_over_fallback":                        {statusDeviated, "17", "the core image element has no debug selector; asserts root size and fallback absence instead of reading avatar-image (ticket17 evidence deviated_list)", true},
	"avatar.rs::fallback_renders_without_an_image":                                {statusDeviated, "17", "the core image element has no debug selector; asserts root size and fallback absence instead of reading avatar-image (ticket17 evidence deviated_list)", true},
	"input/base/state.rs::only_a_multi_line_input_paints_scrollbars":              {statusDeviated, "20", "the single-line half only (ticket20 evidence deviations; input_state_ops_test.go)", true},
	"text/text_view.rs::inline_hover_card_is_lazy_and_retains_atomic_copy":        {statusDeviated, "29", "the hover-card content is a plain div (ticket29 evidence, 28's deferred tests)", true},
}
