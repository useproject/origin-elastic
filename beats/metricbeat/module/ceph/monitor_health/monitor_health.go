package monitor_health

import (
	"github.com/useproject/origin-elastic/beats/libbeat/common"
	"github.com/useproject/origin-elastic/beats/libbeat/common/cfgwarn"
	"github.com/useproject/origin-elastic/beats/metricbeat/helper"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb/parse"
)

const (
	defaultScheme = "http"
	defaultPath   = "/api/v0.1/health"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   defaultPath,
	}.Build()
)

func init() {
	mb.Registry.MustAddMetricSet("ceph", "monitor_health", New,
		mb.WithHostParser(hostParser),
		mb.DefaultMetricSet(),
	)
}

type MetricSet struct {
	mb.BaseMetricSet
	*helper.HTTP
}

func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The ceph monitor_health metricset is beta")

	http, err := helper.NewHTTP(base)
	if err != nil {
		return nil, err
	}
	http.SetHeader("Accept", "application/json")

	return &MetricSet{
		base,
		http,
	}, nil
}

func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	content, err := m.HTTP.FetchContent()
	if err != nil {
		return nil, err
	}

	return eventsMapping(content), nil
}
