/*
   Author:     concurrer
   License:    GNU GENERAL PUBLIC LICENSE  Version 3, 29 June 2007

   Please refer the LICENSE file for further details
*/

package concur

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// ReadInput needs the file as argument to the main.go
func ReadInput(s string) (data [][]uint64) {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println("file read error")
	}
	data = make([][]uint64, 0)
	line := make([]uint64, 0)
	b := make([]byte, 1)
	for {
		count, err2 := file.Read(b)
		if err2 == io.EOF {
			fmt.Println("file read complete")
			break
		} else if err2 != nil {
			fmt.Println("line read error")
			break
		}
		if count > 0 {
			switch {
			case string(b[0]) == "\n":
				//append line to data
				data = append(data, line)
				//blank out line to read the next line
				line = nil
			default: // continue reading the line
				i, _ := strconv.ParseUint(string(b[0]), 10, 64)
				line = append(line, i)
			}
		} else {
			fmt.Println("count problem")
		}
	}
	return
}

// ShowBoard takes the pointer in 1st arg and prints the Board. The 2nd arg
// is a boolean that decides whether to print the alt set or not. 
func ShowBoard(board *Board, showAlt bool) {
	grid := board.grid
	rowsize := len(grid[0])
	for i := 0; i < rowsize; i++ {
		for j := 0; j < rowsize; j++ {
			fmt.Printf("%d", grid[i][j].value)
			fmt.Printf(" [")
			if showAlt && grid[i][j].value == 0 {
				for k := range grid[i][j].alt.set {
					fmt.Printf("%d", k)
				}
			}
			fmt.Printf("] ")
		}
		fmt.Println()
	}
}

// MakeBoard takes the data from the ReadInput and creates a Board of Cells
func MakeBoard(data [][]uint64) (board Board) {
	rowsize := len(data[0])
	var grid = make([][]Cell, 0)

	for i := 0; i < rowsize; i++ {
		var cellSlice []Cell // for grid
		for j := 0; j < rowsize; j++ {

			//prepare an alt for attaching to each Cell
			//make map inside struct initializer is critical otherwise the map is 'nil' so the Add() wont' work
			intSet := IntSet{set: make(map[uint64]bool)}
			for i := 1; i <= rowsize; i++ {
				intSet.Add(uint64(i)) // add 1-9 to intset.. 1-to-rowsize in this case..
			}

			cell := make([]Cell, 1)
			cell[0].value = data[i][j]
			cell[0].row = uint64(i)
			cell[0].column = uint64(j)
			cell[0].alt = intSet // each Cell gets a copy of IntSet initialized with 1-to-rowsize
			cellSlice = append(cellSlice, cell[0])
		}
		grid = append(grid, cellSlice)
	}
	board.grid = grid
	return
}
