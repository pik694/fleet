---
paths:
  - "internal/ui/settings.go"
  - "internal/ui/onboarding.go"
  - "internal/ui/palette.go"
  - "internal/ui/styles.go"
  - "internal/ui/sidebar_preview.go"
  - "internal/editor/**"
  - "internal/config/**"
---

# Settings, themes, editor, onboarding

- Editor: config.editor > $EDITOR > "code" (VS Code). `internal/editor` resolves the name to a command: CLI launcher on PATH if there is one, else `open -a <App>` against the installed bundle. That fallback is what makes JetBrains IDEs (GoLand, PyCharm, IntelliJ, …) work — Toolbox doesn't install `goland`/`pycharm` shims unless asked, so PATH-only lookup failed with "executable not found". Bundles are prefix-matched (`PyCharm Community Edition.app` → `pycharm`) across `/Applications`, `~/Applications`, and `~/Applications/JetBrains Toolbox`, scanned once per launch. The Settings editor cycler offers only editors this machine can actually launch (`editor.Available()`), so a preset can't be a dead option. Flags (`code -n`) are a CLI-only contract: a spec carrying them with no launcher on PATH errors rather than silently dropping them.
- Themes: fleet-pink (default, flagship brand accent `#dc88c0`), tokyo-night, catppuccin-mocha, rose-pine, nord, gruvbox — configurable via settings (S key)
- Settings dialog: S key opens a master-detail overlay — category rail (Appearance / Behavior) on the left, that category's settings on the right; `tab` switches panes, `j/k` move, `←→`/`h/l` cycle a value, auto-saves on `esc`. The **Appearance** category renders a live mock-sidebar preview (synthetic data via `RenderSidebarPreview` → the real `RenderSidebar`, so it can't drift) and holds: Theme, Status style (icon/bar), Agent icons, Slot badges, Account usage (split/grouped/off), PR badges, Dirty marker, Status pills, Header counts, Chevron style (triangle/plusminus), Density (normal/compact). Behavior holds editor, tick, auto-name, copy .claude, confirm before restart, drawer height (rows, capped at the 14-row body max), origin forget removes worktrees, enter mode, default agent, idle-session suspend (off/light/balanced/aggressive).
- Display toggles: each Appearance toggle is a `*bool`/string config field with a default-on getter (`config.go`) feeding a package global in `styles.go` (`ShowAgentGlyphs`, `ShowStatusPills`, `ShowPRBadges`, `ShowDirtyIndicator`, `ShowSlotBadges`, `ShowHeaderCounts`, `ChevronStyle`, `SidebarDensity`) — the same pattern as `StatusIndicatorMode`. Globals are synced from config via `ApplyDisplayConfig(cfg)` in `NewHome` (startup) and after each Appearance change; the sidebar `render*` funcs read them at render time.
- First-run onboarding: on first launch, a one-time theme picker (`internal/ui/onboarding.go`, gated by `display_onboarding_seen` config flag) shows the theme cycler beside an annotated sample sidebar that teaches how to read status dots, agent glyphs, branch, worktree, dirty `*`, and PR badge. `enter` keeps the previewed theme; `s`/`esc` reverts and skips. Either way it's marked seen.
