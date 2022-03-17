package caddyManager

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"strings"
)

func newACMEIssuer(challenges *caddytls.ChallengesConfig) interface{} {
	return &struct {
		Module string `json:"module"`
		caddytls.ACMEIssuer
	}{
		Module: "acme",
		ACMEIssuer: caddytls.ACMEIssuer{
			Challenges: challenges,
		},
	}
}

func newZeroSSLIssuer(challenges *caddytls.ChallengesConfig) interface{} {
	return &struct {
		Module string `json:"module"`
		caddytls.ZeroSSLIssuer
	}{
		Module: "acme",
		ZeroSSLIssuer: caddytls.ZeroSSLIssuer{
			ACMEIssuer: &caddytls.ACMEIssuer{
				Challenges: challenges,
			},
		},
	}
}

func loadTLSConfig(domains []string, dnsChallenges map[string]caddytls.DNSChallengeConfig) *caddytls.TLS {
	var policies []*caddytls.AutomationPolicy

	for domain, challengeConfig := range dnsChallenges {
		var subjects []string

		for _, d := range domains {
			if strings.HasSuffix(d, domain) {
				subjects = append(subjects, d)
			}
		}

		challenges := &caddytls.ChallengesConfig{DNS: &challengeConfig}

		acmeIssuerRaw := caddyconfig.JSON(newACMEIssuer(challenges), nil)
		zeroSSLIssuerRaw := caddyconfig.JSON(newZeroSSLIssuer(challenges), nil)

		policies = append(policies, &caddytls.AutomationPolicy{
			Subjects:   subjects,
			IssuersRaw: []json.RawMessage{acmeIssuerRaw, zeroSSLIssuerRaw},
		})
	}

	return &caddytls.TLS{
		Automation: &caddytls.AutomationConfig{
			Policies: policies,
		},
	}
}
