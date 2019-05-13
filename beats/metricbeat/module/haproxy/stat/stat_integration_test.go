// +build integration

package stat

import (
	"testing"

	mbtest "github.com/useproject/origin-elastic/beats/metricbeat/mb/testing"
	"github.com/useproject/origin-elastic/beats/metricbeat/module/haproxy"
)

func TestData(t *testing.T) {
	f := mbtest.NewEventsFetcher(t, getConfig())
	err := mbtest.WriteEvents(f, t)
	if err != nil {
		t.Fatal("write", err)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "haproxy",
		"metricsets": []string{"stat"},
		"hosts":      []string{"tcp://" + haproxy.GetEnvHost() + ":" + haproxy.GetEnvPort()},
	}
}
