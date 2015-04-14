package parser

// PodConfig represents te configuration of a pod based on Sentinel.
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
