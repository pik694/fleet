# fleet

> This file provides context for AI-assisted development with Claude Code.

TUI tool for managing multiple Claude Code and OpenAI Codex sessions in parallel using tmux.

## Tech Stack
- Go 1.26+, Bubble Tea v2 + Lip Gloss v2 (`charm.land/.../v2`), `charmbracelet/x/vt` (drawer's live terminal), tmux, SQLite (WAL mode)

## Build
```bash
make build    # build to build/fleet
make run      # go run
make test     # go test -race
make fmt      # go fmt
make lint     # golangci-lint
make coverage # test coverage report
```

## Commit Convention
Use [conventional commits](https://www.conventionalcommits.org/). Version is auto-computed from commits:
- `fix: ...` → patch (v0.1.0 → v0.1.1)
- `feat: ...` → minor (v0.1.0 → v0.2.0)
- `feat!: ...` or `BREAKING CHANGE: ...` → major (v0.1.0 → v1.0.0)
- `chore:`, `docs:`, `refactor:`, `test:`, `style:` → patch
- Scopes are optional: `fix(hooks): ...`, `feat(ui): ...`

## Changelog
- Each PR adds a fragment file: `changelog/unreleased/<name>.md` with `type:` frontmatter (added/improved/fixed/changed/removed)
- **Keep fragments concise: 1–2 sentences, lead with the user-facing change.** No implementation detail, marketing tone, or exhaustive sub-points — a user skimming the release notes should grasp it at a glance. e.g. "Status updates now respond in ~150ms instead of up to 2s."
- **Voice — fun via honest specificity, not hype.** Fragments become the release notes (markdown, `` `code` `` + `**bold**`), so write for a reader: open each bullet with a **bold 2–4 word headline** so a release skims in one pass; write in second person ("You can now…", "Sessions now…"); name the annoyance you killed ("…instead of a cryptic 'Operation not permitted'"); numbers over adjectives ("~150ms instead of 2s"); `code`-format keys/commands/paths (`Ctrl+K`, `~/Documents`); one idea per bullet; ban hype words (seamless, powerful, delightful, blazing). Full guide + examples: `changelog/README.md`.
- Add `highlight: true` to a fragment's frontmatter to feature it. The item still renders in its normal type section; the flag *also* copies it into a `### Highlights` block that leads the release (emitted by `prepare-release.yml`). Reserve for changes a user would be glad to discover — most fragments are not highlights.
- **Never decide `highlight` on your own — always ask.** When you write a fragment, ask the user (via the `AskUserQuestion` tool) whether it should be a highlight, and give a recommendation with your reasoning. A highlight *leads the release*, so it's an editorial call about what deserves the user's attention — not a property of the diff, and not something you can infer from how big the change felt to implement. Default recommendation is **no**: most fragments are not highlights. Ask even when you're confident, and ask on the way in (when writing the fragment), not after the PR is up.
- CI checks for fragment presence; comment `/no-changelog` to skip
- At release time, fragments are merged into CHANGELOG.md and deleted

## Release
- Comment `/ship` on any issue or PR to prepare a release
- `/ship 2.0.0` to override the version
- CI opens a release PR with changelog rolled — review and merge to release
- Merging the release PR triggers GoReleaser (binaries, GitHub Release, Homebrew)
- Install: `brew install brizzai/tap/fleet` or run `bash install.sh` (requires `gh` CLI)

## Package Structure
```text
cmd/fleet/main.go      # CLI entry point
cmd/fleet/add.go       # `fleet add [path]` — session in an existing directory
cmd/fleet/session_launch.go # Shared by add + worktree: agent/account resolution, hooks, launch
cmd/fleet/worktree.go  # `fleet worktree <branch>` — worktree + session from the shell
cmd/fleet/send.go      # `fleet send <session> <message>` — message a running session from the shell (-wait: then collect the reply)
cmd/fleet/status.go    # `fleet status` — live status of every session: prompt waiting on, unsent draft, branch, PR (--json)
cmd/fleet/read.go      # `fleet read <session>` — last transcript turns, or the cleaned screen (--lines)
cmd/fleet/wait.go      # `fleet wait <session...>` — block until done/idle/waiting/running; exit 0/2/3/4
cmd/fleet/control.go   # `fleet answer|interrupt|submit|clear|restart|suspend|resume` — guarded pane actions
cmd/fleet/inspect.go   # Shared by the agent-facing commands: state loading, rendering, JSON shape
cmd/fleet/skill.go     # `fleet skill <install|uninstall|status>` — install the agent skill
internal/ticket/       # Tracker-agnostic ticket layer: Provider interface, branch naming, seed prompt, image download, keychain store, ticket.md + meta.json, Materialize
internal/linear/       # Linear provider — GraphQL client, API-key + PKCE OAuth
internal/jira/         # Jira Cloud provider — REST v3 client, API-token Basic auth, ADF→markdown
internal/ticketing/    # The router the UI calls: registry, per-repo binding, key→provider resolution, fan-out
internal/git/exclude.go      # .git/info/exclude writer (resolved via `rev-parse --git-path`)
internal/ui/ticket.go        # TUI glue: materialize on create, infer from branch, status line
internal/ui/workspace_picker_ticket.go # ticket suggestions inside the `w` dialog's New branch field
internal/ui/connect_linear.go # Connect Linear dialog (browser sign-in or pasted API key)
internal/ui/connect_jira.go   # Connect Jira dialog (site + email + API token)
internal/skill/        # Embedded SKILL.md + per-agent install targets (skill.go, SKILL.md)
internal/tmux/tmux.go        # Tmux abstraction (create, kill, capture)
internal/tmux/pty.go         # PTY-based attach with Ctrl+Q detach
internal/session/session.go  # Session model, status detection, claude --resume
internal/session/copilot.go  # Copilot status branch, conversation ownership, relaunch wait (copilot_name.go: titles)
internal/session/paneview.go # Reader's view of a pane: the prompt it waits on (Menu), the input box (draft vs suggestion), cleaned lines
internal/session/turns.go    # Last N turns (prompt / closing reply / tool count) from a Claude transcript's tail
internal/control/            # Agent-facing observer: the TUI's status pipeline run one-shot over rows + hook files + one batched capture
internal/session/storage.go  # SQLite persistence (sessions, shells, slot_bindings, pinned_repos, …)
internal/shell/shell.go      # Shell model: non-agent drawer terminals + tmux-derived status (idle/running/exited)
internal/git/git.go          # Git operations (branch, dirty, worktree)
internal/git/head.go         # ReadHead: checkout root + branch from the HEAD file, no git fork (linked worktrees followed)
internal/git/repo_info.go    # RepoInfo cache + refresh logic
internal/github/pr.go        # GitHub PR info via gh CLI
internal/hooks/              # Hook-based status detection (claude_hooks, codex/opencode/copilot hooks, hook_watcher, status_file)
internal/workspace/provider.go     # Provider interface + GitWorktreeProvider + ShellProvider
internal/workspace/repo_config.go  # Per-repo .fleet.json loading (legacy .bc.json supported) + ResolveProvider
internal/config/config.go    # JSON config (~/.config/fleet/config.json)
internal/claudeaccount/      # Multi-account Claude: per-account CLAUDE_CONFIG_DIR + login, auth-status identity, /api/oauth/usage quota, assignment strategies
internal/ui/accounts.go      # Accounts dialog (Ctrl+K → Manage Claude Accounts) + /login pane + quota poll
internal/ui/account_picker.go # Row-anchored "Move to Account…" picker (context menu → move_account)
internal/ui/accounts_header.go   # Top-right per-account weekly-quota readout
internal/naming/naming.go    # Auto-name sessions via smart heuristic (filler stripping, title-case)
internal/stats/              # Local personal analytics: stats.db recorder, transcript scan, Report/Recap (report.go is the UI contract)
internal/ui/stats_view.go    # Stats screen (`i`) — pure renderer of stats.Report
internal/ui/stats_recap.go   # Weekly "Your week" recap reel + header badge
internal/ui/stats_wiring.go  # Recorder tees, scan goroutine, Report/Recap round trips, badge state
internal/debuglog/           # slog-based debug logging to ~/.config/fleet/debug.log
internal/ui/                 # Bubble Tea TUI (app, sidebar, preview, dialogs, styles)
internal/ui/palette.go       # Theme palette definitions (5 built-in themes)
internal/ui/settings.go      # Settings dialog (S key) — master-detail: Appearance/Behavior categories + live preview
internal/ui/sidebar_preview.go   # Synthetic sidebar fixtures + RenderSidebarPreview (drives Appearance preview + onboarding)
internal/ui/onboarding.go    # First-run theme picker + annotated "how to read the sidebar" sample
internal/ui/keybindings.go   # Centralized keybinding definitions
internal/ui/workspace_picker.go  # Worktree dialog (base branch + new branch + existing worktrees)
internal/ui/workspace_create.go  # Create workspace sub-dialog + PendingWorkspace phantom entries
internal/ui/command_palette.go   # Command palette dialog (Ctrl+K) — fuzzy search over commands + repos/worktrees, renders as overlay
internal/ui/drawer.go            # Terminal drawer (` key): render, slide animation, focus/key handling, shell lifecycle + live-stream wiring (syncShellStream → OutputReader + vterm)
internal/tmux/control_output.go  # Control-mode %output reader (live shell streaming) + octal decoder
internal/vterm/vterm.go          # Insulated charmbracelet/x/vt wrapper (drawer live rendering; strips ESC k titles)
internal/frost/                  # Frost mode: freeze a rendered frame into a styled cell grid (x/vt replay) and run cell effects over it
internal/ui/frost.go             # Frost mode trigger + tick loop
internal/chrome/protocol.go      # Command/Response types, action constants, socket path
internal/chrome/native_host.go   # Native messaging host with Unix socket bridge
internal/chrome/client.go        # TUI-side client (connects to socket, sends commands)
internal/chrome/install.go       # NMH manifest auto-install to Chrome's NativeMessagingHosts dir
chrome-extension/                # Chrome MV3 extension (service worker, manifest, icons)
```

## Design System
**Read `docs/design-system.md` before any UI work.** It is prescriptive, not advisory.

The short version: **focus > selection > mode**, in that fixed order of visual weight, and a **mode never fills**. A background fill is the scarcest thing in the UI — spend it on the cursor and the caret. Roles live in `internal/ui/design.go` (`SelectionPill`/`SelectionBand`/`SelectionMarker`/`ModeOn`/`ModeOff`/`FocusCaret`/`PrimaryAction`/`NewTextInput`); build a style inline and you are inventing a fourth dialect. `Background(ColorAccent)` outside that file **fails the build** (`TestAccentFillIsConstructedInOneFile`), and so does a mode indicator that fills (`TestModeNeverFills`).

The doc also settles the three questions that kept getting answered ad hoc: pill vs band (sized by the filled region, ~40 columns), which presentation tier a new surface takes (full-screen / centered overlay with `dimBackdrop` / row-anchored dropdown with none), and whether a tab is a mode or a selection (**does switching it move the keyboard?**).

Styles in `styles.go` are declared **bare** and constructed **only** in `ApplyPalette` — an initializer alongside would be a second copy, and a style built in only one of the two silently keeps default-pink under every other theme.

## Conventions
- Tmux session prefix: `fleet_` (agent sessions); drawer shells use a distinct `fleetsh_` prefix — intentionally not a prefix of `fleet_`, so shells never leak into agent-session enumeration (`tmux.ListSessions`)
- Session ID format: `<8hex>-<unix_timestamp>`
- SQLite DB: `~/.config/fleet/state.db`
- Sessions grouped by git repo root in sidebar with tree lines (├─/└─)
- Status: Running, Waiting, Finished, Idle, Error, Starting, Suspended
- Status icons: ● (running/finished), ◐ (waiting), ○ (idle/starting), ✕ (error), · dim (suspended — same dot as idle; dim style + "suspended" label distinguish it)
- Agent glyph: each session row shows a dim, monochrome per-agent sigil between the status dot and the title — `✻` Claude, `◇` Codex, `△` OpenCode, `⚇` Copilot (`agentGlyph` + `AgentGlyphStyle` in `sidebar.go`/`styles.go`); all are width-1 glyphs from well-covered Unicode blocks (Dingbats / Geometric Shapes — same block as the status dots; `⚇` is Misc Symbols, East-Asian-Neutral and present in Menlo) so they stay aligned in base mono fonts; identity is carried by shape so the status dot keeps sole ownership of color; empty/legacy `Agent` falls back to Claude
- Repo headers show branch name (), dirty indicator (*), and PR badge (#N)
- Git info refreshes every 2s (branch/dirty), PR info every 60s via `gh` CLI
- PR badge: green ✓ (approved+CI passed), yellow ◔ (CI running) / ⌕ (review required) — stacked `◔⌕` when both, and shown only when there's nothing red to fix, red ✕ (CI fail) / ↩ (changes requested or unresolved threads), purple ⇡ (merged), gray ◌ (draft; ✕ appended on CI fail), hidden (closed)
- PR info includes unresolved review thread count via GitHub GraphQL API
- `gh` CLI optional — PR info hidden if not installed
- Preview strips OSC-8 hyperlink sequences to prevent dotted underline artifacts
- All blocking I/O (tmux, git, gh) runs in background worker goroutine, never in Bubble Tea Update()
- Debug log: `~/.config/fleet/debug.log` (slog, init in TUI and hook-handler)
- Config file: `~/.config/fleet/config.json` (tick_interval_sec, default_project_path, editor, theme, auto_name_sessions, copy_claude_settings, confirm_before_restart, origin_delete_removes_worktrees, drawer_height, session_suspend_mode, account_strategy, default_account, allowed_accounts, ticket_start_state (legacy linear_ticket_start still read), account_usage_style, stats_recap_seen_week)
- Claude Code + OpenAI Codex + OpenCode + GitHub Copilot CLI, macOS + Linux (see `docs/linux.md` for the platform split)

## Topic rules

Detailed conventions live in `.claude/rules/`, one file per topic. Each loads automatically when you read or edit a file matching its `paths:`; open one directly when planning work in that area.

- `.claude/rules/cli-launch.md` — CLI: `fleet worktree`, `fleet add`, launch flags
- `.claude/rules/cli-control-plane.md` — CLI: `fleet send` and the agent-facing control plane
- `.claude/rules/agent-skill.md` — `fleet skill` and the embedded SKILL.md
- `.claude/rules/external-sessions.md` — Adopting sessions changed by another fleet process
- `.claude/rules/ui-navigation.md` — Keybindings, help sheet, context menu, palette, tips, delete
- `.claude/rules/keyboard-layouts.md` — Non-Latin keyboard layouts
- `.claude/rules/terminal-drawer.md` — Terminal drawer
- `.claude/rules/workspace.md` — Worktrees, `.fleet.json`, base-branch fetch, worktree removal
- `.claude/rules/tmux-terminal.md` — tmux integration: attach, clipboard, extended keys, mouse, macOS TCC
- `.claude/rules/suspend-snooze.md` — Idle-session suspend and snooze
- `.claude/rules/stats.md` — Local stats screen
- `.claude/rules/tickets.md` — Tickets: Linear and Jira
- `.claude/rules/worktree-dialog.md` — The `w` worktree dialog: base-branch autocomplete and ticket suggestions
- `.claude/rules/status-detection.md` — Status detection, hooks and per-agent behaviour
- `.claude/rules/claude-accounts.md` — Multiple Claude accounts
- `.claude/rules/settings-appearance.md` — Settings, themes, editor, onboarding
- `.claude/rules/session-naming.md` — Session auto-naming and titles
- `.claude/rules/chrome-extension.md` — Chrome tab control
