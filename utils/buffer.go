package utils

type Buffer struct {
	capacity int
	values   []float64
}

func BuildBuffer(size int) Buffer {
	b := Buffer{
		capacity: size,
		values:   make([]float64, 0, size),
	}

	return b
}

// return values if full
func (b *Buffer) PushValue(value float64) []float64 {
	b.values = append(b.values, value)
	if b.isFull() {
		return b.RemoveValues()
	}

	return nil
}

// return values if full
func (b *Buffer) PushValues(values []float64) []float64 {
	free := b.capacity - len(b.values)

	if len(values) < free {
		b.values = append(b.values, values...)
		return nil
	}

	if len(values) == free {
		ret := append(b.values, values...)
		b.values = b.values[:0]
		return ret
	}

	ret := append(b.values, values[:free]...)
	b.values = b.values[:0]
	rest := b.PushValues(values[free:])
	if rest != nil {
		ret = append(ret, rest...)
	}

	return ret
}

func (b *Buffer) GetValues() []float64 {
	return b.values
}

func (b *Buffer) RemoveValues() []float64 {
	ret := make([]float64, 0, b.capacity)
	ret = append(ret, b.values...)
	b.values = b.values[:0]

	return ret
}

func (b *Buffer) GetCapacity() int {
	return b.capacity
}

func (b *Buffer) GetFreeSlots() int {
	return b.capacity - len(b.values)
}

func (b *Buffer) isFull() bool {
	return b.capacity == len(b.values)
}

func (b *Buffer) Clear() {
	b.values = b.values[:0]
}
