package v1alpha1

import (
	"encoding/json"
)

// +k8s:deepcopy-gen=true
type EndpointEndpoint struct {
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Name string `yaml:"name" json:"name,omitempty"`
	Group string `yaml:"group,omitempty" json:"group,omitempty"`
	URL string `yaml:"url" json:"url,omitempty"`
	Method string `yaml:"method,omitempty" json:"method,omitempty"`
	Body string `yaml:"body,omitempty" json:"body,omitempty"`
	GraphQL bool `yaml:"graphql,omitempty" json:"graphql,omitempty"`
	Headers map[string]string `json:"-"`
	Interval string `yaml:"interval,omitempty" json:"interval,omitempty"`
	Conditions []string `yaml:"conditions" json:"conditions,omitempty"`
	Alerts []*AlertAlert `yaml:"alerts,omitempty" json:"alerts,omitempty"`
	DNSConfig *DnsConfig `yaml:"dns,omitempty" json:"dns,omitempty"`
	SSHConfig *SshConfig `yaml:"ssh,omitempty" json:"ssh,omitempty"`
	ClientConfig *ClientConfig `yaml:"client,omitempty" json:"client,omitempty"`
	UIConfig *UiConfig `yaml:"ui,omitempty" json:"ui,omitempty"`
	NumberOfFailuresInARow int `yaml:"-" json:"-"`
	NumberOfSuccessesInARow int `yaml:"-" json:"-"`
}

// +k8s:deepcopy-gen=true
type AlertAlert struct {
	Type string `yaml:"type" json:"type,omitempty"`
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	FailureThreshold int `yaml:"failure-threshold" json:"failure-threshold,omitempty"`
	SuccessThreshold int `yaml:"success-threshold" json:"success-threshold,omitempty"`
	Description *string `yaml:"description,omitempty" json:"description,omitempty"`
	SendOnResolved *bool `yaml:"send-on-resolved,omitempty" json:"send-on-resolved,omitempty"`
	ProviderOverride map[string]json.RawMessage `json:"-"`
	ResolveKey string `yaml:"-" json:"-"`
	Triggered bool `yaml:"-" json:"-"`
}

// +k8s:deepcopy-gen=true
type ClientOAuth2Config struct {
	TokenURL string `yaml:"token-url" json:"token-url,omitempty"`
	ClientID string `yaml:"client-id" json:"client-id,omitempty"`
	ClientSecret string `yaml:"client-secret" json:"client-secret,omitempty"`
	Scopes []string `yaml:"scopes" json:"scopes,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientTLSConfig struct {
	CertificateFile string `yaml:"certificate-file,omitempty" json:"certificate-file,omitempty"`
	PrivateKeyFile string `yaml:"private-key-file,omitempty" json:"private-key-file,omitempty"`
	RenegotiationSupport string `yaml:"renegotiation,omitempty" json:"renegotiation,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientConfig struct {
	ProxyURL string `yaml:"proxy-url,omitempty" json:"proxy-url,omitempty"`
	Insecure bool `yaml:"insecure,omitempty" json:"insecure,omitempty"`
	IgnoreRedirect bool `yaml:"ignore-redirect,omitempty" json:"ignore-redirect,omitempty"`
	Timeout string `yaml:"timeout" json:"timeout,omitempty"`
	DNSResolver string `yaml:"dns-resolver,omitempty" json:"dns-resolver,omitempty"`
	OAuth2Config *ClientOAuth2Config `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
	IAPConfig *ClientIAPConfig `yaml:"identity-aware-proxy,omitempty" json:"identity-aware-proxy,omitempty"`
	Network string `yaml:"network" json:"network,omitempty"`
	TLS *ClientTLSConfig `yaml:"tls,omitempty" json:"tls,omitempty"`
}

// +k8s:deepcopy-gen=true
type UiResponseTime struct {
	Thresholds []int `yaml:"thresholds" json:"thresholds,omitempty"`
}

// +k8s:deepcopy-gen=true
type UiConfig struct {
	HideConditions bool `yaml:"hide-conditions" json:"hide-conditions,omitempty"`
	HideHostname bool `yaml:"hide-hostname" json:"hide-hostname,omitempty"`
	HideURL bool `yaml:"hide-url" json:"hide-url,omitempty"`
	DontResolveFailedConditions bool `yaml:"dont-resolve-failed-conditions" json:"dont-resolve-failed-conditions,omitempty"`
	Badge *UiBadge `yaml:"badge" json:"badge,omitempty"`
}

// +k8s:deepcopy-gen=true
type DnsConfig struct {
	QueryType string `yaml:"query-type" json:"query-type,omitempty"`
	QueryName string `yaml:"query-name" json:"query-name,omitempty"`
}

// +k8s:deepcopy-gen=true
type SshConfig struct {
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientIAPConfig struct {
	Audience string `yaml:"audience" json:"audience,omitempty"`
}

// +k8s:deepcopy-gen=true
type UiBadge struct {
	ResponseTime *UiResponseTime `yaml:"response-time" json:"response-time,omitempty"`
}

