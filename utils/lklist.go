package utils

type LkNode struct {
	Key   string
	Value interface{}
	Next  *LkNode
	Prev  *LkNode
}

type LkList struct {
	Tail  string
	Head  string
	List  map[string]*LkNode
	Limit int
}

func CreateLkList(limit int) LkList {
	l := LkList{
		Tail:  "",
		Head:  "",
		List:  make(map[string]*LkNode, limit),
		Limit: limit,
	}

	return l
}

func (lk *LkList) PushValue(key string, value interface{}) {
	if len(lk.List) == 0 {
		initializeList(key, value, lk)
	} else {
		pushValue(key, value, lk)
	}
}

func initializeList(key string, value interface{}, lk *LkList) {
	node := LkNode{
		Key:   key,
		Value: value,
		Next:  nil,
		Prev:  nil,
	}

	lk.Head = key
	lk.Tail = key
	lk.List[key] = &node
}

func pushValue(key string, value interface{}, lk *LkList) {
	if _, ok := lk.List[key]; ok {
		updateValue(key, value, lk)
	} else {
		addValue(key, value, lk)
	}
}

func updateValue(key string, value interface{}, lk *LkList) {
	node := lk.List[key]
	node.Value = value
	if lk.Tail == key {
		lk.Tail = node.Next.Key
		node.Next.Prev = nil
	} else {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}

	node.Prev = lk.List[lk.Head]
	node.Next = nil
	lk.List[lk.Head].Next = node
	lk.Head = key
}

func addValue(key string, value interface{}, lk *LkList) {
	node := LkNode{
		Key:   key,
		Value: value,
		Next:  nil,
		Prev:  nil,
	}

	node.Prev = lk.List[lk.Head]
	lk.List[lk.Head].Next = &node
	lk.Head = key
	lk.List[key] = &node

	if len(lk.List) >= lk.Limit {
		oldTail := lk.Tail
		lk.Tail = lk.List[lk.Tail].Next.Key
		lk.List[lk.Tail].Prev = nil
		delete(lk.List, oldTail)
	}
}

func (lk *LkList) GetHead() (string, interface{}) {
	if node, ok := lk.List[lk.Head]; ok {
		return node.Key, node.Value
	}

	return "", nil
}

func (lk *LkList) GetTail() (string, interface{}) {
	if node, ok := lk.List[lk.Tail]; ok {
		return node.Key, node.Value
	}

	return "", nil
}

func (lk *LkList) GetValues() map[string]interface{} {
	values := make(map[string]interface{}, lk.Limit)
	for key, node := range lk.List {
		values[key] = node.Value
	}

	return values
}

func (lk *LkList) ClearList() {
	lk.Head = ""
	lk.Tail = ""
	lk.List = make(map[string]*LkNode, lk.Limit)
}
