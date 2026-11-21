---
paths:
  - "internal/naming/**"
  - "internal/session/**"
---

# Session auto-naming and titles

- Auto-naming: sessions titled from the agent's own session title (`session.ReadAgentSessionName` dispatches per agent); falls back to fleet's prompt heuristic (filler stripping, word-boundary truncation) only when the agent has written no title yet
- Agent titles: Claude's come from the conversation JSONL (`ReadClaudeSessionName`). Codex's come from its state DB (`ReadCodexSessionName`) — `/rename` inside Codex runs `UPDATE threads SET title`, and the rollout JSONL carries no title at all, so SQLite is the only source. It's read **read-only** at `$CODEX_HOME/state_<N>.sqlite` (highest `N` wins; Codex bumps the suffix with its schema). Codex *seeds* `threads.title` with the **raw, untruncated first user message**, so a title still equal to `first_user_message` is not a rename and is ignored — adopting it would swap fleet's short heuristic title for a multi-thousand-char prompt. Copilot's come from `$COPILOT_HOME/session-state/<id>/workspace.yaml` (`ReadCopilotSessionName`): its `name` is model-generated after the first turn (like Claude's `ai-title`) and replaced on `/rename` (`user_named: true`), absent before that. OpenCode titles its sessions too; fleet doesn't read them yet.
- Title precedence: manual R-key rename (`ManuallyRenamed`) > the agent's own title (Claude `custom-title` from an in-session `/rename`, else Claude `ai-title` — auto, model-generated, evolves with the work; or a Codex `/rename`) > fleet prompt heuristic
- Auto-naming pipeline: worker cycle re-reads the agent's title each cycle until one appears (`agentNameFreshPollWindow`, so prompt pickup within ~2s), then ~every 30s (`agentNameRecheckInterval`) to follow `ai-title`/`custom-title`/`/rename` drift; heuristic fallback uses UserPromptSubmit hook → status file → HookWatcher → Session.FirstPrompt → naming.GenerateTitle
- Manual rename (R key) sets ManuallyRenamed flag, prevents auto-rename
