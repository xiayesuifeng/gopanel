package caddy

import (
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/caddy/config"
	"gitlab.com/xiayesuifeng/gopanel/core"
)

const tlsApi = "/config/apps/tls"

func AddTLSPolicy(subjects []string) error {
	policy := &config.PolicyType{Subjects: subjects}

	var challenges *config.ChallengesType

	if core.Conf.Caddy.TLS.DNS != nil {
		challenges = &config.ChallengesType{
			Dns: &config.DNSProviderType{
				Provider:  core.Conf.Caddy.TLS.DNS.Provider,
				Resolvers: core.Conf.Caddy.TLS.DNS.Resolvers,
			},
		}
	}

	for _, issuerName := range core.Conf.Caddy.TLS.Issuers {
		switch issuerName {
		case "acme":
			issuer := config.NewPolicyIssuerAcmeType()

			issuer.Challenges = challenges

			policy.Issuers = append(policy.Issuers, issuer)
		case "internal":
			if challenges == nil {
				issuer := config.NewPolicyIssuerInternalType()

				policy.Issuers = append(policy.Issuers, issuer)
			}

		case "zerossl":
			issuer := config.NewPolicyIssuerZerosslType()

			issuer.Challenges = challenges

			policy.Issuers = append(policy.Issuers, issuer)
		}

	}

	tls := &config.TLSType{
		Automation: &config.AutomationType{
			Policies: []*config.PolicyType{policy},
		},
	}

	resp, err := getClient().R().SetBody(tls).Post(tlsApi)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	}

	return nil
}
