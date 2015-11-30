package enum

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

var (
	label_w Enum = WHITE
	label_g Enum = GREEN
	label_y Enum = YELLOW
	label_o Enum = ORANGE
	label_r Enum = RED

	data_s Enum = STATS
	data_a Enum = ANALYTICS
	data_p Enum = PLANS

	action_no   Enum    = NOACTION
	action_st   Enum    = START
	action_sp   Enum    = STOP
	actions_all Actions = Actions{action_no.(Action), action_st.(Action), action_sp.(Action)}

	owner_l Enum = LOCAL
	owner_c Enum = CLUSTER
)

func TestValue(t *testing.T) {

	assert.Equal(t, -1.0, label_w.Value())
	assert.Equal(t, -0.5, label_g.Value())
	assert.Equal(t, 0.0, label_y.Value())
	assert.Equal(t, 0.5, label_o.Value())
	assert.Equal(t, 1.0, label_r.Value())

	assert.Equal(t, 0.0, data_s.Value())
	assert.Equal(t, 1.0, data_a.Value())
	assert.Equal(t, 2.0, data_p.Value())

	assert.Equal(t, 0.0, action_no.Value())
	assert.Equal(t, 1.0, action_st.Value())
	assert.Equal(t, 2.0, action_sp.Value())

	assert.Equal(t, 0.0, owner_l.Value())
	assert.Equal(t, 1.0, owner_c.Value())

}

func TestToString(t *testing.T) {

	assert.Equal(t, "WHITE", label_w.ToString())
	assert.Equal(t, "GREEN", label_g.ToString())
	assert.Equal(t, "YELLOW", label_y.ToString())
	assert.Equal(t, "ORANGE", label_o.ToString())
	assert.Equal(t, "RED", label_r.ToString())

	assert.Equal(t, "STATS", data_s.ToString())
	assert.Equal(t, "ANALYTICS", data_a.ToString())
	assert.Equal(t, "PLANS", data_p.ToString())

	assert.Equal(t, "NOACTION", action_no.ToString())
	assert.Equal(t, "START", action_st.ToString())
	assert.Equal(t, "STOP", action_sp.ToString())
	assert.Equal(t, []string{"NOACTION", "START", "STOP"}, actions_all.ToString())

	assert.Equal(t, "LOCAL", owner_l.ToString())
	assert.Equal(t, "CLUSTER", owner_c.ToString())
}
