# Render Outbound UDP Test

Confirms — rather than assumes — whether Render allows outbound UDP,
which both the Go-library and sidecar qURL paths require (per the
qurl-go README: native registration runs over "authenticated NHP UDP",
no HTTP fallback).

## What it does
Hits three independent public NTP servers (Cloudflare, Google, and the
pool.ntp.org rotation) over raw UDP and waits for a real reply. Success
on any of them proves outbound UDP genuinely leaves and returns to the
platform. Failure on all three is strong evidence of a platform-level
block, since NTP over UDP/123 is about as universally permitted as
outbound UDP traffic gets — if this fails everywhere, qURL's NHP UDP
almost certainly will too.

## Deploy to Render
1. Push this folder to a repo (or a new one) Render can see.
2. In the Render dashboard: **New > Web Service** → connect the repo.
3. Runtime: **Go** (native runtime, no Dockerfile needed).
   - Build command: `go build -o app .`
   - Start command: `./app`
4. Deploy. Once live, hit:
   ```
   curl https://<your-service>.onrender.com/test-udp
   ```

## Reading the result
```json
{
  "any_server_reachable_via_udp": true/false,
  "results": [ ... per-server detail ... ]
}
```
- **`true`** — outbound UDP works on Render. Worth one follow-up test
  directly against the qURL sandbox Hub (`hub.nhp.layerv.xyz` per the
  alpha CLI config) before fully committing, since a permissive NTP
  path doesn't guarantee every UDP port/destination is open.
- **`false` across all three** — treat as a real platform block. That
  settles the Render question independent of which qURL integration
  path (sidecar vs. embedded Go library) you pick, since both need the
  same NHP UDP transport.

## Note on timeouts
Each server gets a 5s dial+read timeout, so the whole `/test-udp` call
resolves in ~15s worst case — comfortably inside Render's request
timeout, so no need to raise any platform timeout setting for this test.
