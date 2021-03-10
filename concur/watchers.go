/*
   Author:     concurrer
   License:    GNU GENERAL PUBLIC LICENSE  Version 3, 29 June 2007

   Please refer the LICENSE file for further details
*/
package concur

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// LaunchWatchers does the following:
// * initiates the Board Chan (boardChan)
// * initiates the Channel Grid (chanGrid) with one chan for each Cell 
// * launches the Listeners for both boardChan and chanGrid
// * launches the CellWatchers for each Cell on the Board, 
func LaunchWatchers(board *Board) {
	rowsize := len(board.grid[0])
	var waitGroup sync.WaitGroup

	// initiate a BoardChan as a parent channel to listen for incominng events from Cell.Fill()
	// then inform the other CellChans in the ChanGrid who would be interested in that event
	boardChan := make(chan *Cell, rowsize*rowsize) // buffered channel

	// initiate channel grid with one channel for each Cell to update its alt/value based on incoming events
	// row/column/block annoucements will write to interested CellChans when a Cell is filled
	chanGrid := make([][]chan uint64, 0)
	for i := 0; i < rowsize; i++ {
		tempChanSlice := make([]chan uint64, 0)
		for j := 0; j < rowsize; j++ {
			tempChan := make(chan uint64, rowsize*rowsize)
			tempChanSlice = append(tempChanSlice, tempChan)
		}
		chanGrid = append(chanGrid, tempChanSlice)
	}

	// wait until the LaunchChannelListener is done
	waitGroup.Add(1)
	go LaunchChannelListener(boardChan, chanGrid, &waitGroup)

	// wait until the LaunchCellWatchers is done
	waitGroup.Add(1)
	go LaunchCellWatchers(boardChan, board, chanGrid, &waitGroup)

	// wait
	waitGroup.Wait()
}

// LaunchChannelListener 's boardChan waits for incoming Cells sent by CellWatcher
// whenever a Cell is Filled. The boardChan informs the CellChans on the
// ChanGrid of other "interested" Cells in each Fill event.
// i.e. for every incoming Cell on BoardChan
// * calculate the row/column/block Cells of interest
// * send the newly filled number to the respective CellChan's so as to change their altSet
func LaunchChannelListener(boardChan chan *Cell, chanGrid [][]chan uint64, waitGroup *sync.WaitGroup) {

	defer func() { waitGroup.Done() }()

	// close all channels if not solved in certain time
	boom := time.After(3 * time.Second)

	for {
		select {
		case cell := <-boardChan:
			// wait until the AnnounceToOtherCells is done
			waitGroup.Add(1)
			go AnnounceToOtherCells(cell, chanGrid, waitGroup)
		case <-boom:
			fmt.Print("\n\nHere's the final board! :\n\n")
			return
		default:
			// does nothing but continue the loop
		}
	}
}

// AnnouceToOtherCells announces a Cell.Fill event to all other Cells of interest (row/column/block)
func AnnounceToOtherCells(cell *Cell, chanGrid [][]chan uint64, waitGroup *sync.WaitGroup) {
	defer func() {
		waitGroup.Done()
	}()
	rowsize := len(chanGrid[0])
	colsize := rowsize // assuming 2 variables for convenience, given that the incoming sudoku is a square
	row := cell.row
	column := cell.column
	// inform every column in the row
	for i := 0; i < colsize; i++ {
		if uint64(i) != cell.column { // don't have to inform the same cell that originated this update
			chanGrid[row][i] <- cell.value
		}
	}
	// inform every row in the column
	for i := 0; i < rowsize; i++ {
		if uint64(i) != cell.row { // don't have to inform the same cell that originated this update
			chanGrid[i][column] <- cell.value
		}
	}
	// inform other cells in the same block, typically the size must be a perfect square,
	// hence the below  check, run this only for a perfect square
	sqrt := math.Sqrt(float64(rowsize))
	if float64(rowsize) == sqrt*sqrt {
		i, j := cell.row, cell.column
		bi, bj := i/3*3, j/3*3         //block counters
		for k := 0; k < rowsize; k++ { // loop for 'size'.. rowsize here is only a convenience variable
			kk := uint64(k)
			chanGrid[bi+kk/3][bj+kk%3] <- cell.value
		}
	}
}

// LaunchCellWatchers launches a CellWatcher AND a CellChan for each Cell
func LaunchCellWatchers(boardChan chan *Cell, board *Board, chanGrid [][]chan uint64, waitGroup *sync.WaitGroup) {
	defer func() {
		waitGroup.Done()
	}()
	rowsize := len(board.grid[0])
	for i := 0; i < rowsize; i++ {
		for j := 0; j < rowsize; j++ {
			// wait until each CellWather is done
			waitGroup.Add(1)
			go CellWatcher(boardChan, &board.grid[i][j], chanGrid[i][j], waitGroup)
		}

	}
}

// CellWatcher is launched per Cell. It loops indefinitely :
// * read from CellChan and remove that number from altSet 
// * when alt.SingleAlt returns true then Fill the Cell with the only remaining alt
// * then inform the BoardChan to annouce the Fill
func CellWatcher(boardChan chan *Cell, cell *Cell, cellChan chan uint64, waitGroup *sync.WaitGroup) {
	defer func() {
		waitGroup.Done()
	}()

	// if the Cell is already filled from input.txt, then inform the
	// boardChan to let others know... AND return
	if cell.value > uint64(0) {
		boardChan <- cell
		return
	}

	// as long as cell is not filled (==0), listen on CellChan...
	// when u get a value, it means that value is filled somewhere else
	// update alt and then check if u have SingleAlt if yes then fill, else continue loop
	boom := time.After(5 * time.Second)
	for cell.value == 0 {
		select {
		case i := <-cellChan:
			cell.alt.mutex.Lock()
			cell.Delete(uint64(i)) // delete the incoming value from altSet
			cell.alt.mutex.Unlock()

			if cell.alt.SingleAlt() {
				cell.alt.mutex.Lock()
				cell.Fill() // fill the cell with the only remaining alt
				cell.alt.mutex.Unlock()
				//also inform boardChan to let others know
				boardChan <- cell
				return
			}
		case <-boom:
			return
		default:
			// does nothing but continue the loop
		}
	}
	for cell.value > 0 { // channel is filled.. continue receiving the values on cellchan and drop them
		<-cellChan
	}
}
