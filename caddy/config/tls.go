package config

import "encoding/json"

type TLSType struct {
	Automation *AutomationType `json:"automation,omitempty"`
}

type AutomationType struct {
	Policies []*PolicyType `json:"policies,omitempty"`
}

type PolicyType struct {
	Subjects []string      `json:"subjects,omitempty"`
	Issuers  []interface{} `json:"issuers,omitempty"`
}

type PolicyIssuerType struct {
	Module      string          `json:"module"`
	CA          string          `json:"ca,omitempty"`
	TestCA      string          `json:"test_ca,omitempty"`
	Email       string          `json:"email,omitempty"`
	AccountKey  string          `json:"account_key,omitempty"`
	AcmeTimeout int             `json:"acme_timeout,omitempty"`
	Challenges  *ChallengesType `json:"challenges,omitempty"`
}

type ChallengesType struct {
	Dns *DNSProviderType `json:"dns,omitempty"`
}

type DNSProviderType struct {
	Provider  json.RawMessage `json:"provider"`
	Resolvers []string        `json:"resolvers"`
}

type PolicyIssuerInternalType struct {
	Module       string `json:"module,omitempty"`
	CA           string `json:"ca,omitempty"`
	Lifetime     int    `json:"lifetime,omitempty"`
	SignWithRoot bool   `json:"sign_with_root,omitempty"`
}

func NewPolicyIssuerAcmeType() *PolicyIssuerType {
	return &PolicyIssuerType{
		Module: "acme",
	}
}

func NewPolicyIssuerZerosslType() *PolicyIssuerType {
	return &PolicyIssuerType{
		Module: "zerossl",
	}
}

func NewPolicyIssuerInternalType() *PolicyIssuerInternalType {
	return &PolicyIssuerInternalType{
		Module: "internal",
	}
}
