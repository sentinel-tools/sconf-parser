package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

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

// LocalSentinelConfig is a struct holding information about the sentinel RS is
// running on.
type SentinelConfig struct {
	Name              string
	Host              string
	Port              int
	ManagedPodConfigs map[string]PodConfig
	Dir               string
	KnownSentinels    []string
	BadDirectives     []string
}

func (c *SentinelConfig) GetPod(podname string) (PodConfig, error) {
	pod, exists := c.ManagedPodConfigs[podname]
	if !exists {
		return pod, fmt.Errorf("Pod '%s' not found in config", podname)
	}
	return pod, nil
}

func ParseSentinelConfig(filename string) (SentinelConfig, error) {
	var config SentinelConfig
	config.ManagedPodConfigs = make(map[string]PodConfig)
	file, err := os.Open(filename)
	if err != nil {
		log.Print(err)
		return config, err
	}
	defer file.Close()
	bf := bufio.NewReader(file)
	for {
		rawline, err := bf.ReadString('\n')
		if err == nil || err == io.EOF {
			line := strings.TrimSpace(rawline)
			// ignore comments
			if strings.Contains(line, "#") {
				continue
			}
			entries := strings.Split(line, " ")
			//Most values are key/value pairs
			switch entries[0] {
			case "sentinel": // Have a sentinel directive
				err := extractSentinelDirective(entries[1:], &config)
				if err != nil {
					config.BadDirectives = append(config.BadDirectives, line)
				}
			case "port":
				iport, _ := strconv.Atoi(entries[1])
				config.Port = iport
			case "dir":
				config.Dir = entries[1]
			case "bind":
				config.Host = entries[1]
			case "":
				if err == io.EOF {
					return config, nil
				}
			default:
				log.Printf("UNhandled Sentinel Directive: %s", line)
				//err := fmt.Errorf("Unhandled sentinel directive: %+v", entries)
				config.BadDirectives = append(config.BadDirectives, line)
			}
		} else {
			log.Print("=============== LOAD FILE ERROR ===============")
			log.Fatal(err)
		}
	}
	return config, nil
}

// extractSentinelDirective parses the sentinel directives from the
// sentinel config file
func extractSentinelDirective(entries []string, c *SentinelConfig) error {
	switch entries[0] {
	case "monitor":
		pname := entries[1]
		port := entries[3]
		quorum := entries[4]
		spc := PodConfig{Name: pname, MasterIP: entries[2], MasterPort: port, Quorum: quorum}
		//spc.KnownSentinels = make(map[string]string)
		// normally we should not see duplicate IP:PORT combos, however it
		// can happen when people do things manually and dont' clean up.
		// We need to detect them and ignore the second one if found,
		// reporting the error condition this will require tracking
		// ip:port pairs...
		addr := fmt.Sprintf("%s:%d", entries[2], port)
		_, exists := c.ManagedPodConfigs[addr]
		if !exists {
			c.ManagedPodConfigs[entries[1]] = spc
		}
		return nil

	case "auth-pass":
		pname := entries[1]
		pc := c.ManagedPodConfigs[pname]
		pc.Authpass = entries[2]
		c.ManagedPodConfigs[pname] = pc
		return nil

	case "known-sentinel":
		podname := entries[1]
		sentinel_address := entries[2] + ":" + entries[3]
		pc, exists := c.ManagedPodConfigs[podname]
		if !exists {
			err := fmt.Errorf("Had known-sentinel entries for nonexistent pod configuration")
			log.Print(err)
			return err
		}
		pc.KnownSentinels = append(pc.KnownSentinels, sentinel_address)
		c.KnownSentinels = append(pc.KnownSentinels, sentinel_address)
		c.ManagedPodConfigs[podname] = pc
		return nil

	case "known-slave":
		podname := entries[1]
		pc, exists := c.ManagedPodConfigs[podname]
		if !exists {
			err := fmt.Errorf("Had known-slave entries for nonexistent pod configuration")
			log.Print(err)
			return err
		}
		slave_address := entries[2] + ":" + entries[3]
		pc.KnownSlaves = append(pc.KnownSlaves, slave_address)
		c.ManagedPodConfigs[podname] = pc
		return nil

	case "config-epoch", "leader-epoch", "current-epoch", "down-after-milliseconds", "maxclients", "parallel-syncs", "failover-timeout":
		// We don't use these keys
		return nil

	default:
		return fmt.Errorf("Unhandled sentinel directive")
	}
}
