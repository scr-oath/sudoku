/*
   Author:     concurrer
   License:    GNU GENERAL PUBLIC LICENSE  Version 3, 29 June 2007

   Please refer the LICENSE file for further details
*/
package concur

import (
	"fmt"
	"sync"
)

// Board
type Board struct {
	grid [][]Cell
}

// Cell
type Cell struct {
	row, column uint64
	value       uint64
	alt         IntSet
}

// Fill is called when a Cell has just one alt.. i.e the alt.SingleAlt is true
func (cell *Cell) Fill() {
	for k := range cell.alt.set { // it has only one value at this point
		cell.value = k          // fill the Cell with that value
		delete(cell.alt.set, k) // and empty the alt
	}
}

func (cell *Cell) Delete(k uint64) {
	delete(cell.alt.set, k)
}

// Stringer for Cell
func (cell *Cell) String() (s string) {
	for k := range cell.alt.set {
		s += fmt.Sprint(k)
	}
	return fmt.Sprintf(s)
}

// IntSet emulates a Set of Integers 
type IntSet struct {
	set   map[uint64]bool
	mutex sync.Mutex
}

func (intSet *IntSet) Add(u uint64) {
	intSet.set[u] = true
}

func (intSet *IntSet) SingleAlt() bool {
	if size := len(intSet.set); size == 1 {
		return true
	}
	return false
}
