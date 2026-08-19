package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type result struct {
	OutboundUDP bool   `json:"outbound_udp_reachable"`
	Server      string `json:"ntp_server"`
	LatencyMS   int64  `json:"latency_ms,omitempty"`
	Error       string `json:"error,omitempty"`
}

// testUDP sends a minimal NTP client request and waits for a reply.
// A successful read proves a real UDP packet left the platform AND a
// reply made it back in — the strongest signal short of testing the
// actual qURL Hub directly.
func testUDP(server string) result {
	r := result{Server: server}

	conn, err := net.DialTimeout("udp", server, 5*time.Second)
	if err != nil {
		r.Error = fmt.Sprintf("dial failed: %v", err)
		return r
	}
	defer conn.Close()

	// Minimal 48-byte NTP client request (mode 3, version 3).
	packet := make([]byte, 48)
	packet[0] = 0x1B

	start := time.Now()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	if _, err := conn.Write(packet); err != nil {
		r.Error = fmt.Sprintf("write failed: %v", err)
		return r
	}

	resp := make([]byte, 48)
	n, err := conn.Read(resp)
	if err != nil {
		r.Error = fmt.Sprintf("no reply within timeout (%v) — outbound UDP is likely blocked, or this specific server is unreachable", err)
		return r
	}

	if n > 0 {
		r.OutboundUDP = true
		r.LatencyMS = time.Since(start).Milliseconds()
	}
	return r
}

func handler(w http.ResponseWriter, req *http.Request) {
	servers := []string{
		"time.cloudflare.com:123",
		"time.google.com:123",
		"pool.ntp.org:123",
	}

	results := make([]result, 0, len(servers))
	anySuccess := false
	for _, s := range servers {
		res := testUDP(s)
		if res.OutboundUDP {
			anySuccess = true
		}
		results = append(results, res)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"any_server_reachable_via_udp": anySuccess,
		"interpretation": map[string]string{
			"if_true":  "Outbound UDP works on this Render service. Worth re-testing directly against the qURL sandbox Hub before finalizing.",
			"if_false": "All three independent NTP servers failed — strong signal outbound UDP is blocked at the platform level, not a fluke of one server.",
		},
		"results": results,
	})
}

func main() {
	http.HandleFunc("/test-udp", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "UDP egress test service. Hit /test-udp to run the check.")
	})

	port := "8080"
	fmt.Println("listening on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
