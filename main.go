/*
    Author:     concurrer 
    License:    GNU GENERAL PUBLIC LICENSE  Version 3, 29 June 2007

    Please refer the LICENSE file for further details
*/
package main

import (
    "fmt"
    "os"
    "sudoku/concur"
)

func main() {
    args := os.Args
    if len(args)<2 {
        fmt.Println("pass the file arg")
        return
    }
    file_input := args[1]
    fmt.Printf("reading file: %s\n",file_input)
    data := concur.ReadInput(file_input)  // [][]uint64
    board := concur.MakeBoard(data)
    // true/false in the Showboard args will decide whether to print the Cell.alt or not
    concur.ShowBoard(&board, true)

    //action 
    concur.LaunchWatchers(&board)

    // check final board
   defer concur.ShowBoard(&board, true)
}
