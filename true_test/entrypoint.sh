#!/bin/sh
set -e

echo "=================================================="
echo "HUB REACHABILITY TEST"
echo "Target Hub: $QURL_CONNECTOR_HUB_HOST:$QURL_CONNECTOR_HUB_PORT"
echo "(This is the value to confirm with Ben before trusting a failure.)"
echo "=================================================="

echo ">>> Minting enrollment token"
export QURL_CONNECTOR_TOKEN=$(curl -s -X POST "$QURL_ENDPOINT/v1/api-keys" \
  -H "Authorization: Bearer $QURL_API_KEY" -H 'content-type: application/json' \
  -d "{\"kind\":\"enrollment_token\",\"name\":\"$CONNECTOR_ID\",\"target\":\"connector\",\"claims\":[{\"type\":\"connector\",\"id\":\"$CONNECTOR_ID\"}],\"expires_in\":\"2h\"}" \
  | jq -r .data.api_key)

if [ -z "$QURL_CONNECTOR_TOKEN" ] || [ "$QURL_CONNECTOR_TOKEN" = "null" ]; then
  echo ">>> RESULT: INCONCLUSIVE — token mint failed (check QURL_API_KEY / QURL_ENDPOINT, not a UDP issue)"
  exit 1
fi

echo ">>> Token minted. Starting connector — watch for 'knock_ok' and 'login_success' below."
echo ">>> A hang or timeout here (not an auth error) is the actual UDP signal to watch for."
echo "=================================================="

# TARGET doesn't need to be a real live service for this test —
# registration, knock, and tunnel login all complete before any
# traffic needs to reach it.
exec qurl connector run --id "$CONNECTOR_ID" --target "$TARGET"
