package probes

import (
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var dnsInitFn = sync.OnceValues(func() (*dns.ClientConfig, error) {
	return dns.ClientConfigFromFile("/etc/resolv.conf")
})

func DNS(typ, val string) (string, error) {
	config, err := dnsInitFn()
	if err != nil {
		return "", err
	}
	c := new(dns.Client)
	m := new(dns.Msg)
	switch strings.ToUpper(typ) {
	case "A":
		m.SetQuestion(val, dns.TypeA)
	case "AAAA":
		m.SetQuestion(val, dns.TypeAAAA)
	case "CNAME":
		m.SetQuestion(val, dns.TypeCNAME)
	case "MX":
		m.SetQuestion(val, dns.TypeMX)
	case "NS":
		m.SetQuestion(val, dns.TypeNS)
	case "PTR":
		m.SetQuestion(val, dns.TypePTR)
	case "TXT":
		m.SetQuestion(val, dns.TypeTXT)
	case "SOA":
		m.SetQuestion(val, dns.TypeSOA)
	case "SRV":
		m.SetQuestion(val, dns.TypeSRV)
	default:
		return "", fmt.Errorf("unsupported DNS resource record type")
	}
	m.RecursionDesired = true
	slog.Debug("querying DNS", "rr_type", typ, "rr_value", val)
	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		slog.Error("querying DNS", "error", err)
		return "", err
	}
	if r.Rcode != dns.RcodeSuccess {
		return "", fmt.Errorf("DNS record not found")
	}
	var sb strings.Builder
	for _, a := range r.Answer {
		sb.WriteString(a.String())
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
