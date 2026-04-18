package firewall

import (
	"net"
	"os/exec"
)

// NftablesManager structures the data needed to interact with the nft cmd tool.
type NftablesManager struct {
	SetName string
	Table   string
}

func NewNftablesManager(table, setName string) *NftablesManager {
	return &NftablesManager{
		Table:   table,
		SetName: setName,
	}
}

func (m *NftablesManager) AllowIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return &FirewallError{Op: "allow", Msg: "invalid IP: " + ip}
	}
	cmd := exec.Command("nft", "add", "element", m.Table, m.SetName, "{", ip, "}")
	if err := cmd.Run(); err != nil {
		return &FirewallError{Op: "allow", Msg: err.Error()}
	}
	return nil
}

func (m *NftablesManager) RevokeIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return &FirewallError{Op: "delete", Msg: "invalid IP: " + ip}
	}
	cmd := exec.Command("nft", "delete", "element", m.Table, m.SetName, "{", ip, "}")
	if err := cmd.Run(); err != nil {
		return &FirewallError{Op: "revoke", Msg: err.Error()}
	}
	return nil
}
