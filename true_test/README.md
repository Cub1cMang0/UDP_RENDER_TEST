# Render Hub Reachability Test (Background Worker)

Reuses the exact `qurl connector run` binary already validated locally
against sandbox — this time deployed on Render itself, to get a real
answer instead of inferring one from a generic NTP probe.

## Why a Background Worker, not a Web Service
This process doesn't listen on HTTP and doesn't need to answer health
checks. A Web Service on Render expects an HTTP port bound and will be
marked unhealthy / cycled if nothing responds. A Background Worker has
no such expectation — it just runs.

## Before deploying
Confirm the Hub host/port with Ben. Two different values exist in
what you've been given:
- `hub.nhp.layerv.xyz:443` — from the internal alpha release notes,
  explicitly marked sandbox. **This is the default in render.yaml.**
- `hub.nhp.layerv.ai:62206` — from the public `qurl-go` README's
  generic quickstart example.

If they're actually the same target under different naming, fine —
but don't let this test report "UDP blocked" when the real issue is
"wrong address." Swap `QURL_CONNECTOR_HUB_HOST` / `_PORT` in the
Render dashboard env vars if Ben says otherwise.

## Deploy
1. Push this folder to a repo Render can see.
2. In Render: **New > Blueprint**, point it at the repo — `render.yaml`
   will configure the Background Worker automatically.
3. Set `QURL_API_KEY` manually in the dashboard (it's marked `sync: false`
   in the manifest on purpose, so it's never committed).
4. Deploy, then watch the **Logs** tab.

## Reading the result

**Pass** — look for this sequence, same as your local test:
```
connector: native registration succeeded   event=native_registration_succeeded
connector: NHP knock ok                    event=knock_ok
connector: tunnel login admitted           event=login_success
```
If you see `login_success`, outbound UDP to the real Hub works. That's
the conclusive answer — stronger than the NTP result, since it's the
actual protocol and actual endpoint, not a proxy for it.

**Fail — and the important part is distinguishing *why*:**
- **Hang / timeout with no `knock_ok`, no error message** → the
  strongest signal of a genuine UDP block. The knock packet went out
  and nothing came back.
- **Immediate connection error, DNS failure, or "no route to host"** →
  likely a wrong host/port, not a UDP block. Re-check the Hub values
  with Ben before concluding anything about Render's network.
- **`RESULT: INCONCLUSIVE — token mint failed`** → an auth/API-key
  problem, unrelated to UDP entirely. Check `QURL_API_KEY` and
  `QURL_ENDPOINT`.

Only the first case is real evidence against Render. The other two are
configuration issues that happen to look similar if you're only
skimming for "did it work."

## After this test
Whatever the result, this doesn't block moving on to the custom
`AgentStateStore` question — that's a separate, parallel track. A
`worker` service on Render's free tier doesn't need to stay running
indefinitely for this specific test either; it's fine to delete it
once you have your answer.
