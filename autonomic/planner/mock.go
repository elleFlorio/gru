package planner

import (
	"github.com/elleFlorio/gru/autonomic/planner/policy"
	"github.com/elleFlorio/gru/enum"
	"github.com/elleFlorio/gru/storage"
)

func init() {
	storage.New("internal")
}

func StoreMockPolicy(plc policy.Policy) {
	data, _ := convertPolicyToData(&plc)
	storage.StoreLocalData(data, enum.POLICIES)
}
