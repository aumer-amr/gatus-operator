package v1alpha1

import (
	"encoding/json"
)

// +k8s:deepcopy-gen=true
type DnsConfig struct {
	QueryType string `json:"query-type,omitempty"`
	QueryName string `json:"query-name,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientConfig struct {
	ProxyURL string `json:"proxy-url,omitempty"`
	Insecure bool `json:"insecure,omitempty"`
	IgnoreRedirect bool `json:"ignore-redirect,omitempty"`
	Timeout string `json:"timeout,omitempty"`
	DNSResolver string `json:"dns-resolver,omitempty"`
	OAuth2Config *ClientOAuth2Config `json:"oauth2,omitempty"`
	IAPConfig *ClientIAPConfig `json:"identity-aware-proxy,omitempty"`
	Network string `json:"network,omitempty"`
	TLS *ClientTLSConfig `json:"tls,omitempty"`
}

// +k8s:deepcopy-gen=true
type UiConfig struct {
	HideConditions bool `json:"hide-conditions,omitempty"`
	HideHostname bool `json:"hide-hostname,omitempty"`
	HideURL bool `json:"hide-url,omitempty"`
	DontResolveFailedConditions bool `json:"dont-resolve-failed-conditions,omitempty"`
	Badge *UiBadge `json:"badge,omitempty"`
}

// +k8s:deepcopy-gen=true
type AlertAlert struct {
	Type string `json:"type,omitempty"`
	Enabled *bool `json:"enabled,omitempty"`
	FailureThreshold int `json:"failure-threshold,omitempty"`
	SuccessThreshold int `json:"success-threshold,omitempty"`
	Description *string `json:"description,omitempty"`
	SendOnResolved *bool `json:"send-on-resolved,omitempty"`
	ProviderOverride map[string]json.RawMessage `json:"-"`
	ResolveKey string `json:"-"`
	Triggered bool `json:"-"`
}

// +k8s:deepcopy-gen=true
type ClientOAuth2Config struct {
	TokenURL string `json:"token-url,omitempty"`
	ClientID string `json:"client-id,omitempty"`
	ClientSecret string `json:"client-secret,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientIAPConfig struct {
	Audience string `json:"audience,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientTLSConfig struct {
	CertificateFile string `json:"certificate-file,omitempty"`
	PrivateKeyFile string `json:"private-key-file,omitempty"`
	RenegotiationSupport string `json:"renegotiation,omitempty"`
}

// +k8s:deepcopy-gen=true
type UiResponseTime struct {
	Thresholds []int `json:"thresholds,omitempty"`
}

// +k8s:deepcopy-gen=true
type UiBadge struct {
	ResponseTime *UiResponseTime `json:"response-time,omitempty"`
}

// +k8s:deepcopy-gen=true
type EndpointEndpoint struct {
	Enabled *bool `json:"enabled,omitempty"`
	Name string `json:"name,omitempty"`
	Group string `json:"group,omitempty"`
	URL string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
	Body string `json:"body,omitempty"`
	GraphQL bool `json:"graphql,omitempty"`
	Headers map[string]string `json:"-"`
	Interval string `json:"interval,omitempty"`
	Conditions []string `json:"conditions,omitempty"`
	Alerts []*AlertAlert `json:"alerts,omitempty"`
	DNSConfig *DnsConfig `json:"dns,omitempty"`
	SSHConfig *SshConfig `json:"ssh,omitempty"`
	ClientConfig *ClientConfig `json:"client,omitempty"`
	UIConfig *UiConfig `json:"ui,omitempty"`
	NumberOfFailuresInARow int `json:"-"`
	NumberOfSuccessesInARow int `json:"-"`
}

// +k8s:deepcopy-gen=true
type SshConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

