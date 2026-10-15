package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestLANURLs(t *testing.T) {
	ips := []net.IP{net.ParseIP("192.168.1.5"), net.ParseIP("10.0.0.2")}
	both := []string{"http://192.168.1.5:8080", "http://10.0.0.2:8080"}
	tests := []struct {
		addr string
		want []string
	}{
		{"127.0.0.1:8080", []string{}},
		{"[::1]:8080", []string{}},
		{"localhost:8080", []string{}}, // not an IP: no guess
		{"0.0.0.0:8080", both},
		{"[::]:8080", both},
		{":8080", both},
		{"192.168.1.5:8080", []string{"http://192.168.1.5:8080"}},
		{"bad", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			got := LANURLs(tt.addr, ips)
			if got == nil || !slices.Equal(got, tt.want) {
				t.Fatalf("LANURLs(%q) = %#v, want %#v", tt.addr, got, tt.want)
			}
		})
	}
}

func TestLANIPsPrivateOnly(t *testing.T) {
	for _, ip := range LANIPs() {
		if !ip.IsPrivate() || ip.To4() == nil {
			t.Fatalf("unexpected %v", ip)
		}
	}
}

func TestPrivateLAN(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", nil, []string{}},
		{"drops non-private", []string{"127.0.0.1", "169.254.1.2", "100.64.0.1", "100.100.1.1", "8.8.8.8", "172.32.0.1", "fd00::1"}, []string{}},
		{"order 192 then 10 then 172", []string{"172.17.0.1", "10.0.0.2", "192.168.1.5"}, []string{"192.168.1.5", "10.0.0.2", "172.17.0.1"}},
		{"stable within group", []string{"10.0.0.9", "172.31.0.1", "192.168.2.1", "10.0.0.1", "192.168.1.1", "172.16.0.1"},
			[]string{"192.168.2.1", "192.168.1.1", "10.0.0.9", "10.0.0.1", "172.31.0.1", "172.16.0.1"}},
		{"mixed with filtered", []string{"100.64.0.1", "172.20.0.1", "169.254.0.1", "192.168.0.10"}, []string{"192.168.0.10", "172.20.0.1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var in []net.IP
			for _, s := range tt.in {
				in = append(in, net.ParseIP(s))
			}
			got := []string{}
			for _, ip := range privateLAN(in) {
				got = append(got, ip.String())
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("privateLAN(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestAPILANLoopbackBind(t *testing.T) {
	s, err := Start("127.0.0.1:0", t.TempDir(), t.TempDir(), slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown(context.Background())
	url := "http://" + s.Addr() + "/api/lan"

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("status %d content-type %q", res.StatusCode, res.Header.Get("Content-Type"))
	}
	var body map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || string(body["urls"]) != "[]" {
		t.Fatalf("body %v", body)
	}

	post, err := http.Post(url, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	post.Body.Close()
	if post.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("POST status %d", post.StatusCode)
	}
}
