package policy

import (
	//"math"

	"github.com/elleFlorio/gru/autonomic/analyzer"
	//"github.com/elleFlorio/gru/enum"
	//"github.com/elleFlorio/gru/service"
)

type swapCreator struct{}

func (p *swapCreator) getPolicyName() string {
	return "swap"
}

func (p *swapCreator) createPolicy(analytics analyzer.GruAnalytics, targetList ...string) Policy {
	return Policy{}
}

func (p *swapCreator) listActions() []string {
	return []string{"stop", "start"}
}
