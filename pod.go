package parser

import (
	"fmt"

	"github.com/therealbill/libredis/client"
)

/*
type PodConfig struct {
	Name           string
	MasterIP       string
	MasterPort     string
	Authpass       string
	KnownSentinels []string
	KnownSlaves    []string
	Settings       map[string]string
	Quorum         string
	BadDirectives  [][]string
}
*/

// GetSlaves() returns the list of known slaves
func (p *PodConfig) GetSlaves() ([]string, error) {
	return p.KnownSlaves, nil
}

// GetSentinels() returns the list of known sentinels
func (p *PodConfig) GetSentinels() ([]string, error) {
	return p.KnownSentinels, nil
}

// ValidateMaster() connexts to the master listed in the config and verifies it
// has the role of master
func (p *PodConfig) ValidateMaster() (bool, error) {
	dc := client.DialConfig{Address: fmt.Sprintf("%s:%s", p.MasterIP, p.MasterPort), Password: p.Authpass}
	c, err := client.DialWithConfig(&dc)
	if err != nil {
		return false, err
	}
	role, err := c.RoleName()
	if err != nil {
		return false, err
	}
	if role != "master" {
		return false, nil
	}
	return true, nil
}
