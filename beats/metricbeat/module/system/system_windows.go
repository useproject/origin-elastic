package system

import (
	"github.com/useproject/origin-elastic/beats/libbeat/logp"
	"github.com/useproject/origin-elastic/beats/metricbeat/helper"
)

func initModule() {
	if err := helper.CheckAndEnableSeDebugPrivilege(); err != nil {
		logp.Warn("%v", err)
	}
}
