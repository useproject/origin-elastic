package self

import (
	"github.com/useproject/origin-elastic/beats/libbeat/common"
	"github.com/useproject/origin-elastic/beats/libbeat/common/cfgwarn"
	"github.com/useproject/origin-elastic/beats/metricbeat/helper"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb/parse"
)

const (
	defaultScheme = "http"
	defaultPath   = "/v2/stats/self"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   defaultPath,
	}.Build()
)

func init() {
	mb.Registry.MustAddMetricSet("etcd", "self", New,
		mb.WithHostParser(hostParser),
		mb.DefaultMetricSet(),
	)
}

type MetricSet struct {
	mb.BaseMetricSet
	http *helper.HTTP
}

func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The etcd self metricset is beta")
	config := struct{}{}

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	http, err := helper.NewHTTP(base)
	if err != nil {
		return nil, err
	}
	return &MetricSet{
		base,
		http,
	}, nil
}

func (m *MetricSet) Fetch() (common.MapStr, error) {
	content, err := m.http.FetchContent()
	if err != nil {
		return nil, err
	}
	return eventMapping(content), nil
}
