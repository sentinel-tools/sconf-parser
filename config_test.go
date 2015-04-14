package parser

import (
	"fmt"
	"testing"
)

/*
func TestValidPodFromSentinelConfig(t *testing.T) {
	_, err := getPodInfoFromConfigFile("sentinel.conf", "pod1")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestINValidPodFromSentinelConfig(t *testing.T) {
	_, err := getPodInfoFromConfigFile("sentinel.conf", "pod2")
	if err == nil {
		t.Error(err)
		t.Fail()
	} else {
		if !strings.Contains(err.Error(), "not found in config file") {
			t.Error("Somehow found a pod which doesn't exist!")
			t.Fail()
		}
	}
}
*/

func TestValidPodFromSentinelConfig(t *testing.T) {
	_, err := ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestBadDirectives(t *testing.T) {
	sconf, err := ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(sconf.BadDirectives) == 0 {
		t.Error(fmt.Errorf("Should have had bad directives, had none."))
		t.Fail()
	}
}

func TestGetPod(t *testing.T) {
	sconf, err := ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	pod, err := sconf.GetPod("pod1")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if pod.Name != "pod1" {
		t.Error(fmt.Errorf("retreived pod is not named 'pod1'"))
		t.Fail()
	}
}
func TestPodKnownSlaves(t *testing.T) {
	sconf, err := ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	pod, err := sconf.GetPod("pod1")
	if len(pod.KnownSlaves) != 2 {
		t.Error(fmt.Errorf("Mismatched KnownSlaves. Expected 2, got %d", len(pod.KnownSlaves)))
		t.Fail()
	}
}
func TestPodKnownSentinels(t *testing.T) {
	sconf, err := ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	pod, err := sconf.GetPod("pod1")
	if len(pod.KnownSentinels) != 2 {
		t.Error(fmt.Errorf("Mismatched KnownSentinels. Expected 2, got %d", len(pod.KnownSentinels)))
		fmt.Printf("pod: %+v\n", pod)
		t.Fail()
	}
}
