package enum

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

var (
	data_s Enum = STATS
	data_a Enum = ANALYTICS
	data_p Enum = POLICIES
	data_i Enum = INFO

	action_no   Enum    = NOACTION
	action_st   Enum    = START
	action_sp   Enum    = STOP
	action_rm   Enum    = REMOVE
	actions_all Actions = Actions{action_no.(Action), action_st.(Action), action_sp.(Action), action_rm.(Action)}

	owner_l Enum = LOCAL
	owner_c Enum = CLUSTER
)

func TestValue(t *testing.T) {

	assert.Equal(t, 0.0, data_s.Value())
	assert.Equal(t, 1.0, data_a.Value())
	assert.Equal(t, 2.0, data_p.Value())
	assert.Equal(t, 3.0, data_i.Value())

	assert.Equal(t, 0.0, action_no.Value())
	assert.Equal(t, 1.0, action_st.Value())
	assert.Equal(t, 2.0, action_sp.Value())
	assert.Equal(t, 3.0, action_rm.Value())

	assert.Equal(t, 0.0, owner_l.Value())
	assert.Equal(t, 1.0, owner_c.Value())

}

func TestToString(t *testing.T) {

	assert.Equal(t, "STATS", data_s.ToString())
	assert.Equal(t, "ANALYTICS", data_a.ToString())
	assert.Equal(t, "POLICIES", data_p.ToString())
	assert.Equal(t, "INFO", data_i.ToString())

	assert.Equal(t, "NOACTION", action_no.ToString())
	assert.Equal(t, "START", action_st.ToString())
	assert.Equal(t, "STOP", action_sp.ToString())
	assert.Equal(t, "REMOVE", action_rm.ToString())
	assert.Equal(t, []string{"NOACTION", "START", "STOP", "REMOVE"}, actions_all.ToString())

	assert.Equal(t, "LOCAL", owner_l.ToString())
	assert.Equal(t, "CLUSTER", owner_c.ToString())
}
