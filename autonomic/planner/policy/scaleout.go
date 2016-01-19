package policy

import (
	"math"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	"github.com/elleFlorio/gru/enum"
	res "github.com/elleFlorio/gru/resources"
	"github.com/elleFlorio/gru/service"
)

const c_THRESHOLD_SCALEOUT_LOAD = 0.7
const c_THRESHOLD_SCALEOUT_CPU = 0.7

type ScaleOut struct{}

func (p *ScaleOut) Name() string {
	return "scaleout"
}

//TODO find a way to compute a label that make some sense...
func (p *ScaleOut) Weight(name string, analytics analyzer.GruAnalytics) float64 {
	srv, _ := service.GetServiceByName(name)
	inst_run := len(srv.Instances.Running)
	inst_pen := len(srv.Instances.Pending)

	if (inst_pen + inst_run) > 0 {
		return 0.0
	}

	srvAnalytics := analytics.Service[name]
	if srvAnalytics.Resources.Available < 1.0 {
		return 0.0
	}

	srvCores := srv.Docker.CpusetCpus
	if srvCores != "" {
		if !res.CheckSpecificCoresAvailable(srvCores) {
			return 0.0
		}
	}

	// LOAD
	load := srvAnalytics.Load
	value_load := math.Max(load, c_THRESHOLD_SCALEOUT_LOAD)
	weight_load := (value_load - c_THRESHOLD_SCALEOUT_LOAD) / (1 - c_THRESHOLD_SCALEOUT_LOAD)
	// CPU
	cpu := srvAnalytics.Resources.Cpu
	value_cpu := math.Max(cpu, c_THRESHOLD_SCALEOUT_CPU)
	weight_cpu := (value_cpu - c_THRESHOLD_SCALEOUT_CPU) / (1 - c_THRESHOLD_SCALEOUT_CPU)
	// MEMORY
	// TODO?

	policyValue := (weight_load + weight_cpu) / 2

	return policyValue
}

func (p *ScaleOut) Actions() enum.Actions {
	return []enum.Action{
		enum.START,
	}
}
