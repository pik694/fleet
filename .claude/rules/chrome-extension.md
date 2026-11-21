---
paths:
  - "internal/chrome/**"
  - "chrome-extension/**"
---

# Chrome tab control

- Chrome tab control: `p` opens PR in Chrome via extension (reuses existing tab), falls back to `open <url>` if unavailable
- Chrome extension architecture: TUI →[unix socket]→ native host (`fleet chrome-host`) →[stdio]→ Chrome service worker
- Native messaging host: `fleet chrome-host` subcommand (also auto-detected when Chrome passes `chrome-extension://...` arg)
- Unix socket: `~/.config/fleet/chrome.sock` (created by native host, mode 0600)
- NMH manifest: auto-installed to `~/Library/Application Support/Google/Chrome/NativeMessagingHosts/com.brizzai.fleet.tabcontrol.json` on TUI startup
- Chrome extension ID: `haphpcoecelhofejcklinnlbfijgdnih` (stable via `key` in manifest.json)
- Extension commands: `open_or_focus`, `close_tab`, `create_tab_group`, `ping`
- Service worker reconnects to native host on disconnect (2s delay)
