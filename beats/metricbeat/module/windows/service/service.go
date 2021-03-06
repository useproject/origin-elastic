// +build windows

package service

import (
	"github.com/useproject/origin-elastic/beats/libbeat/common"
	"github.com/useproject/origin-elastic/beats/libbeat/common/cfgwarn"
	"github.com/useproject/origin-elastic/beats/metricbeat/mb"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	mb.Registry.MustAddMetricSet("windows", "service", New,
		mb.DefaultMetricSet(),
	)
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	reader *ServiceReader
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The windows service metricset is beta")

	reader, err := NewServiceReader()
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		reader:        reader,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	services, err := m.reader.Read()
	if err != nil {
		return nil, err
	}

	return services, nil
}
