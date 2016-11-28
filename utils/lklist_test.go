package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestCreateLkList(t *testing.T) {
	l := CreateLkList(5)
	assert.Equal(t, 5, l.Limit)
}

func TestLkPushValue(t *testing.T) {
	keys := []string{"k1", "k2", "k3", "k4", "k5"}
	values := []float32{1, 2, 3, 4, 5}

	// Add the first value
	l := CreateLkList(5)
	assert.Equal(t, "empty", l.ToString())

	l.PushValue(keys[0], values[0])
	assert.Equal(t, l.head, keys[0])
	assert.Equal(t, l.tail, keys[0])
	assert.Equal(t, values[0], l.list[l.head].value)
	//fmt.Println(l.ToString())

	// Saturate the list
	for i := 1; i < len(keys); i++ {
		l.PushValue(keys[i], values[i])
	}
	// Test list print
	assert.Equal(t, "k5->k4->k3->k2->k1->end", l.ToString())
	assert.Equal(t, keys[4], l.head)
	assert.Equal(t, keys[0], l.tail)

	// Add one value more
	l.PushValue("k6", float32(6))
	//fmt.Println(l.ToString())
	assert.Equal(t, "k6", l.head)

	// Update value in the middle
	l.PushValue(keys[3], values[3])
	//fmt.Println(l.ToString())
	assert.Equal(t, keys[3], l.head)
	assert.Equal(t, keys[1], l.tail)

	// Update value that is tail
	l.PushValue(keys[1], values[1])
	//fmt.Println(l.ToString())
	assert.Equal(t, keys[1], l.head)

	// Update value that is head
	l.PushValue(keys[1], values[1])
	//fmt.Println(l.ToString())
	assert.Equal(t, keys[1], l.head)
}

func TestGetHead(t *testing.T) {
	l := CreateLkList(5)
	head, value := l.GetHead()
	assert.Equal(t, "", head)
	assert.Nil(t, value)

	l.PushValue("key", "value")
	head, value = l.GetHead()
	assert.Equal(t, "key", head)
	assert.NotNil(t, value)
}

func TestGetTail(t *testing.T) {
	l := CreateLkList(5)
	tail, value := l.GetTail()
	assert.Equal(t, "", tail)
	assert.Nil(t, value)

	l.PushValue("key", "value")
	tail, value = l.GetTail()
	assert.Equal(t, "key", tail)
	assert.NotNil(t, value)
}

func TestLkGetValues(t *testing.T) {
	l := CreateLkList(5)
	assert.Empty(t, l.GetValues())

	keys := []string{"k1", "k2", "k3", "k4", "k5"}
	values := []float32{1, 2, 3, 4, 5}
	for i := 0; i < len(keys); i++ {
		l.PushValue(keys[i], values[i])
	}
	assert.Equal(t, 5, len(l.list))

}

func TestClearList(t *testing.T) {
	l := CreateLkList(5)
	keys := []string{"k1", "k2", "k3", "k4", "k5"}
	values := []float32{1, 2, 3, 4, 5}
	for i := 0; i < len(keys); i++ {
		l.PushValue(keys[i], values[i])
	}

	l.ClearList()
	assert.Empty(t, l.list)

}
