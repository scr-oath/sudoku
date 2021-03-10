package cmd

import (
	"github.com/concurrer/sodoku/concur"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sudokuCmd represents the sudoku command
var sudokuCmd = &cobra.Command{
	Use:   "sudoku file_input",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Debugf("sudoku called with %v", args)
		fileInput := args[0]
		logrus.Debugf("reading file: %s", fileInput)
		data := concur.ReadInput(fileInput) // [][]uint64
		board := concur.MakeBoard(data)
		// true/false in the Showboard args will decide whether to print the Cell.alt or not
		concur.ShowBoard(&board, true)

		//action
		concur.LaunchWatchers(&board)

		// check final board
		concur.ShowBoard(&board, true)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(sudokuCmd)
}
