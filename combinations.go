package fifty2

type CardSliceIterator interface {
	HasNext() bool
	Next() []Card
}

type comboIterator struct {
	slice []Card
	choose int
	index []int
	done bool
}

func Combinations(slice []Card, choose int) CardSliceIterator {
	if choose > len(slice) {
		panic("fifty2: cannot produce combinations larger than given card slice")
	}
	iterator := &comboIterator{
		slice: slice,
		choose: choose,
		index: make([]int, choose),
		done: false,
	}
	iterator.prime()
	return iterator
}

func (ci *comboIterator) prime() {
	for i := 1; i <  ci.choose; i++ {
		if ci.index[i] <= ci.index[i-1] {
			ci.index[i] = ci.index[i-1] + 1
		}
	}
}

func (ci *comboIterator) moveNext() {
	inc := ci.choose-1
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
