package fifty2

type CardSliceIterator interface {
	HasNext() bool
	Next() []Card
}

type CardSlice2DIterator interface {
	HasNext() bool
	Next() [][]Card
}

type comboIterator struct {
	slice  []Card
	choose int
	index  []int
	done   bool
}

func Combinations(slice []Card, choose int) CardSliceIterator {
	if choose > len(slice) {
		panic("fifty2: cannot produce combinations larger than given card slice")
	}
	iterator := &comboIterator{
		slice:  slice,
		choose: choose,
		index:  make([]int, choose),
		done:   false,
	}
	iterator.prime()
	return iterator
}

func (ci *comboIterator) prime() {
	for i := 1; i < ci.choose; i++ {
		if ci.index[i] <= ci.index[i-1] {
			ci.index[i] = ci.index[i-1] + 1
		}
	}
}

func (ci *comboIterator) moveNext() {
	inc := ci.choose - 1
	reprime := false
	for inc >= 0 {
		maxSliceIndex := len(ci.slice) - (ci.choose - inc - 1)
		ci.index[inc] = (ci.index[inc] + 1) % maxSliceIndex
		if ci.index[inc] == 0 {
			inc--
			reprime = true
		} else {
			break
		}
	}
	if inc < 0 {
		ci.done = true
	} else if reprime {
		ci.prime()
	}
}

func (ci *comboIterator) HasNext() bool {
	return !ci.done
}

func (ci *comboIterator) Next() []Card {
	combo := make([]Card, ci.choose)
	for i, sliceIndex := range ci.index {
		combo[i] = ci.slice[sliceIndex]
	}
	ci.moveNext()
	return combo
}

type comboSetIterator struct {
	slice     [][]Card
	choose    []int
	iterators []CardSliceIterator
	next      [][]Card
	done      bool
}

func MultipleCombinations(slice []Card, choose []int) CardSlice2DIterator {
	iterator := &comboSetIterator{
		slice:     make([][]Card, len(choose)),
		choose:    choose,
		iterators: make([]CardSliceIterator, len(choose)),
		next:      make([][]Card, 0, len(choose)),
		done:      false,
	}
	iterator.slice[0] = slice
	iterator.prime()
	return iterator
}

func (csi *comboSetIterator) prime() {
	for i := len(csi.next); i < len(csi.choose); i++ {
		if i > 0 {
			nextSlice := make([]Card, len(csi.slice[i-1]))
			copy(nextSlice, csi.slice[i-1])
			csi.slice[i] = Remove(nextSlice, csi.next[i-1]...)
		}
		itr := Combinations(csi.slice[i], csi.choose[i])
		csi.iterators[i] = itr
		if itr.HasNext() {
			csi.next = append(csi.next, itr.Next())
		} else {
			csi.next = append(csi.next, []Card{})
		}
	}
}

func (csi *comboSetIterator) moveNext() {
	inc := len(csi.choose) - 1
	for inc >= 0 {
		itr := csi.iterators[inc]
		if itr.HasNext() {
			csi.next[inc] = itr.Next()
			break
		} else {
			csi.next = csi.next[:inc]
			inc--
		}
	}
	if inc < 0 {
		csi.done = true
	} else {
		csi.prime()
	}
}

func (csi *comboSetIterator) HasNext() bool {
	return !csi.done
}

func (csi *comboSetIterator) Next() [][]Card {
	next := make([][]Card, len(csi.next))
	copy(next, csi.next)
	csi.moveNext()
	return next
}
