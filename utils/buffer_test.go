package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestPushValue(t *testing.T) {
	value := 1.0
	var b Buffer
	var res []float64

	b = BuildBuffer(10)
	res = b.PushValue(value)
	res = b.PushValue(value)
	res = b.PushValue(value)
	assert.Nil(t, res)
	assert.Len(t, b.values, 3)

	b = BuildBuffer(1)
	res = b.PushValue(value)
	assert.NotNil(t, res)
	assert.Empty(t, b.values)
}

func TestPushValues(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	var b Buffer
	var res []float64

	b = BuildBuffer(10)
	res = b.PushValues(values)
	assert.Nil(t, res)
	assert.Len(t, b.values, len(values))

	b = BuildBuffer(6)
	res = b.PushValues(values)
	assert.NotNil(t, res)
	assert.Len(t, res, len(values))
	assert.Empty(t, b.values)

	b = BuildBuffer(2)
	res = b.PushValues(values)
	assert.NotNil(t, res)
	assert.Len(t, res, len(values))
	assert.Empty(t, b.values)

	b = BuildBuffer(4)
	res = b.PushValues(values)
	assert.NotNil(t, res)
	assert.Len(t, res, b.capacity)
	assert.Len(t, b.values, len(values)-b.capacity)
}

func TestGetValues(t *testing.T) {
	b := BuildBuffer(10)
	values := []float64{1, 2, 3}
	b.values = append(b.values, values...)
	res := b.GetValues()
	assert.Len(t, res, len(values))
	assert.NotEmpty(t, b.values)
}

func TestRemoveValues(t *testing.T) {
	b := BuildBuffer(10)
	values := []float64{1, 2, 3}
	b.values = append(b.values, values...)
	res := b.RemoveValues()
	assert.Len(t, res, len(values))
	assert.Empty(t, b.values)
}

func TestGetCapacity(t *testing.T) {
	b := BuildBuffer(10)
	c := b.GetCapacity()
	assert.Equal(t, b.capacity, c)
}

func TestGetFreeSlots(t *testing.T) {
	b := BuildBuffer(10)
	values := []float64{1, 2, 3}
	var f int

	f = b.GetFreeSlots()
	assert.Equal(t, b.capacity, f)

	b.values = append(b.values, values...)
	f = b.GetFreeSlots()
	assert.Equal(t, b.capacity-len(b.values), f)
}

func TestIsFull(t *testing.T) {
	b := BuildBuffer(5)
	values := []float64{1, 2, 3, 4, 5}
	var f bool

	f = b.isFull()
	assert.False(t, f)

	b.values = append(b.values, values...)
	f = b.isFull()
	assert.True(t, f)
}

func TestClear(t *testing.T) {
	b := BuildBuffer(5)
	values := []float64{1, 2, 3, 4, 5}
	b.values = append(b.values, values...)

	b.Clear()
	assert.Empty(t, b.values)
}
