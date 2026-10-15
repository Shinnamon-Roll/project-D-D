package server

import (
	"encoding/json"
	"net"
	"net/http"
	"slices"
)

// LANIPs lists the private IPv4 addresses of up, non-loopback interfaces, preferred first (see privateLAN).
func LANIPs() []net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ips []net.IP
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			if n, ok := a.(*net.IPNet); ok {
				if ip4 := n.IP.To4(); ip4 != nil {
					ips = append(ips, ip4)
				}
			}
		}
	}
	return privateLAN(ips)
}

// privateLAN keeps only RFC 1918 IPv4 (drops loopback, link-local, CGNAT/Tailscale, public) and orders
// 192.168/16, then 10/8, then 172.16/12 (Docker/VM bridges usually live there); stable within a group.
func privateLAN(ips []net.IP) []net.IP {
	rank := func(ip net.IP) int {
		switch ip[0] {
		case 192:
			return 0
		case 10:
			return 1
		}
		return 2
	}
	out := []net.IP{}
	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil && ip4.IsPrivate() {
			out = append(out, ip4)
		}
	}
	slices.SortStableFunc(out, func(a, b net.IP) int { return rank(a) - rank(b) })
	return out
}

// LANURLs returns http://<ip>:<port> links other devices can use to reach a server bound to listenAddr.
// Loopback bind → none; unspecified bind (0.0.0.0, ::, "") → one per ip; specific IPv4 bind → only that one.
func LANURLs(listenAddr string, ips []net.IP) []string {
	urls := []string{}
	host, port, err := net.SplitHostPort(listenAddr)
	if err != nil {
		return urls
	}
	if host != "" {
		bound := net.ParseIP(host)
		if bound == nil || bound.IsLoopback() {
			return urls
		}
		if !bound.IsUnspecified() {
			if bound.To4() != nil {
				urls = append(urls, "http://"+net.JoinHostPort(bound.String(), port))
			}
			return urls
		}
	}
	for _, ip := range ips {
		urls = append(urls, "http://"+net.JoinHostPort(ip.String(), port))
	}
	return urls
}

func handleLAN(listenAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"urls": LANURLs(listenAddr, LANIPs())})
	}
}
