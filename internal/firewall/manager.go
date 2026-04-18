// Package firewall defines the interface for managing firewall rules.
package firewall

// FirewallManager is the interface that firewall adapters must implement.
type FirewallManager interface {
	// AllowIP opens access for the given IP address.
	AllowIP(ip string) error
	// RevokeIP removes access for the given IP address.
	RevokeIP(ip string) error
}

// FirewallError describes a failed firewall operation.
type FirewallError struct {
	Op  string // "allow" or "revoke"
	Msg string
}

func (e *FirewallError) Error() string {
	return "firewall: " + e.Op + ": " + e.Msg
}
