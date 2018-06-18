package sabakan

import "regexp"

const (
	// BmcIdrac9 is BMC type for iDRAC-9
	BmcIdrac9 = "iDRAC-9"
	// BmcIpmi2 is BMC type for IPMI-2.0
	BmcIpmi2 = "IPMI-2.0"
)

// Machine represents a server hardware.
type Machine struct {
	Serial      string     `json:"serial"`
	Product     string     `json:"product"`
	Datacenter  string     `json:"datacenter"`
	Rack        uint       `json:"rack"`
	IndexInRack uint       `json:"index-in-rack"`
	Role        string     `json:"role"`
	IPv4        []string   `json:"ipv4"`
	IPv6        []string   `json:"ipv6"`
	BMC         MachineBMC `json:"bmc"`
}

var (
	reValidRole = regexp.MustCompile(`^[0-9a-zA-Z._-]+$`)
)

// IsValidRole returns true if role is valid as machine role
func IsValidRole(role string) bool {
	return reValidRole.MatchString(role)
}

// MachineBMC is a bmc interface struct for Machine
type MachineBMC struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
	Type string `json:"type"`
}
