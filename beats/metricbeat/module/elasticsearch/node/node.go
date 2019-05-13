package node

import (
	"github.com/useproject/origin-elastic/beats/libbeat/common/cfgwarn"
	"github.com/useproject/origin-elastic/beats/metricbeat/helper"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb/parse"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	mb.Registry.MustAddMetricSet("elasticsearch", "node", New,
		mb.WithHostParser(hostParser),
		mb.DefaultMetricSet(),
	)
}

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: "http",
		PathConfigKey: "path",
		// This only fetches data for the local node.
		DefaultPath: "_nodes/_local",
	}.Build()
)

// MetricSet type defines all fields of the MetricSet
type MetricSet struct {
	mb.BaseMetricSet
	http *helper.HTTP
}

// New create a new instance of the MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The elasticsearch node metricset is beta")

	http, err := helper.NewHTTP(base)
	if err != nil {
		return nil, err
	}
	return &MetricSet{
		base,
		http,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch(r mb.ReporterV2) {
	content, err := m.http.FetchContent()
	if err != nil {
		r.Error(err)
		return
	}

	eventsMapping(r, content)
}
