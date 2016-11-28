package utils

import (
	"bytes"
)

type LkNode struct {
	key   string
	value interface{}
	next  *LkNode
	prev  *LkNode
}

type LkList struct {
	tail  string
	head  string
	list  map[string]*LkNode
	Limit int
}

func CreateLkList(limit int) LkList {
	l := LkList{
		tail:  "",
		head:  "",
		list:  make(map[string]*LkNode, limit+1),
		Limit: limit,
	}

	return l
}

func (lk *LkList) PushValue(key string, value interface{}) {
	if len(lk.list) == 0 {
		initializeList(key, value, lk)
	} else {
		pushValue(key, value, lk)
	}
}

func initializeList(key string, value interface{}, lk *LkList) {
	node := LkNode{
		key:   key,
		value: value,
		next:  nil,
		prev:  nil,
	}

	lk.head = key
	lk.tail = key
	lk.list[key] = &node
}

func pushValue(key string, value interface{}, lk *LkList) {
	if _, ok := lk.list[key]; ok {
		updateValue(key, value, lk)
	} else {
		addValue(key, value, lk)
	}
}

func updateValue(key string, value interface{}, lk *LkList) {
	node := lk.list[key]
	node.value = value
	if lk.head == key {
		return
	}

	if lk.tail == key {
		lk.tail = node.next.key
		node.next.prev = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}

	node.prev = lk.list[lk.head]
	node.next = nil
	lk.list[lk.head].next = node
	lk.head = key
}

func addValue(key string, value interface{}, lk *LkList) {
	node := LkNode{
		key:   key,
		value: value,
		next:  nil,
		prev:  nil,
	}

	node.prev = lk.list[lk.head]
	lk.list[lk.head].next = &node
	lk.head = key
	lk.list[key] = &node

	if len(lk.list) > lk.Limit {
		oldTail := lk.tail
		lk.tail = lk.list[lk.tail].next.key
		lk.list[lk.tail].prev = nil
		delete(lk.list, oldTail)
	}
}

func (lk *LkList) GetHead() (string, interface{}) {
	if node, ok := lk.list[lk.head]; ok {
		return node.key, node.value
	}

	return "", nil
}

func (lk *LkList) GetTail() (string, interface{}) {
	if node, ok := lk.list[lk.tail]; ok {
		return node.key, node.value
	}

	return "", nil
}

func (lk *LkList) GetValues() map[string]interface{} {
	values := make(map[string]interface{}, lk.Limit)
	for key, node := range lk.list {
		values[key] = node.value
	}

	return values
}

func (lk *LkList) ClearList() {
	lk.head = ""
	lk.tail = ""
	lk.list = make(map[string]*LkNode, lk.Limit)
}

func (lk *LkList) ToString() string {
	if len(lk.list) == 0 {
		return "empty"
	}

	var buffer bytes.Buffer
	arrow := "->"
	current := lk.list[lk.head]
	buffer.WriteString(current.key)
	for current != nil {
		buffer.WriteString(arrow)
		current = current.prev
		if current != nil {
			buffer.WriteString(current.key)
		}
	}

	buffer.WriteString("end")

	return buffer.String()
}
