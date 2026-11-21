---
type: changed
---
**State files are private.** `debug.log`, hook status files (which hold your prompts), snapshots and the databases under `~/.config/fleet` are now created `0600` in `0700` directories instead of world-readable.
