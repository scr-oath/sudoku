# sudoku
Solving sudoku in a concurrent way is more human. Golang only makes it easier…  :)<br>

![chan logic](https://github.com/concurrer/sudoku/blob/master/images/sudoku-golang-channels.png)
<br>
Key points:<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1. Create a CellGrid with respective Cells having the input numbers.<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2. Create a ChannelGrid which is of same size as the CellGrid. <br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;3. For every cell on the CellGrid, create a CellChan on the ChannelGrid. <br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;4. Each Cell maintains an 'altSet' which is an unordered 'Set' of possible numbers (1-9 in this example) <br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;5. If a Cell has the initial value given, 'Announce' that value on the ChannelGrid to all the CellChans of "interested" Cells.  This is the crux of the logic.  Please check the diagram.  The "interested" cells are the ones on the same row, same column and same block. <br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;6. If a Cell doesn't have an initial value, then it keeps listening on its respective CellChan on the ChannelGrid.<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;7. If a Cell receives a number on its CellChan, it removes that number from its 'altSet'. Because that number is already filled somewhere else (either in the same row, same column or same block).<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;8. If a Cell's 'altSet' has only one value, then Fill the Cell with that value and 'Announce' it to others.<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;9. BoardChan acts as an intermediate layer between CellGrid and ChannelGrid and helps implement the "interestingness" logic.

That’s all. :)

Initial Input is a text file with numbers and zeros.  Ex:<br>
530070000<br>
600195000<br>
098000060<br>
800060003<br>
400803001<br>
700020006<br>
060000280<br>
000419005<br>
000080079<br>

Initial Board implementation with 'altSet's for each Cell. Please note, if
the Cell is having an initial value, its altSet is redundant so not shown. 

--------------------------------------------------------------------------------------------------
5 [] 3 [] 0 [613578924] 0 [124793568] 7 [] 0 [239156784] 0 [234615789] 0 [715468923] 0 [125934678]<br>
6 [] 0 [345612789] 0 [491267835] 1 [] 9 [] 5 [] 0 [124579368] 0 [613489257] 0 [892456713]<br>
0 [467912358] 9 [] 8 [] 0 [812473569] 0 [134692578] 0 [891276345] 0 [934867125] 6 [] 0 [567913824]<br>
8 [] 0 [157234689] 0 [467891352] 0 [238156794] 6 [] 0 [681259347] 0 [347891256] 0 [345916782] 3 []<br>
4 [] 0 [234791685] 0 [893456712] 8 [] 0 [356792481] 3 [] 0 [891256734] 0 [135682479] 1 []<br>
7 [] 0 [789126345] 0 [235671489] 0 [457891362] 2 [] 0 [371456892] 0 [925671348] 0 [678123459] 6 []<br>
0 [945678123] 6 [] 0 [567913428] 0 [912678345] 0 [245613789] 0 [236814579] 2 [] 8 [] 0 [734568912]<br>
0 [812567934] 0 [134568279] 0 [378124569] 4 [] 1 [] 9 [] 0 [589712346] 0 [345812679] 5 []<br>
0 [632457891] 0 [791458236] 0 [356812479] 0 [413567892] 8 [] 0 [523467891] 0 [567924813] 7 [] 9 []<br>


Here's the final board at the end:


5 [] 3 [] 4 [] 6 [] 7 [] 8 [] 9 [] 1 [] 2 []<br>
6 [] 7 [] 2 [] 1 [] 9 [] 5 [] 3 [] 4 [] 8 []<br>
1 [] 9 [] 8 [] 3 [] 4 [] 2 [] 5 [] 6 [] 7 []<br>
8 [] 5 [] 9 [] 7 [] 6 [] 1 [] 4 [] 2 [] 3 []<br>
4 [] 2 [] 6 [] 8 [] 5 [] 3 [] 7 [] 9 [] 1 []<br>
7 [] 1 [] 3 [] 9 [] 2 [] 4 [] 8 [] 5 [] 6 []<br>
9 [] 6 [] 1 [] 5 [] 3 [] 7 [] 2 [] 8 [] 4 []<br>
2 [] 8 [] 7 [] 4 [] 1 [] 9 [] 6 [] 3 [] 5 []<br>
3 [] 4 [] 5 [] 2 [] 8 [] 6 [] 1 [] 7 [] 9 []<br>

--------------------------------------------------------------------------------------------------

Advantages in this approach:<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1. Lean code with almost no resource consumption.<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2. Works with any 'size' of a sudoku (must be a perfect square though) but its easy to change the code for any different size.

Next Steps:<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1. It cannot solve 'very hard' sudoku puzzles. That’s because I don't know how to solve them personally. Once I figure that piece out, I will update the code :)<br>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2.The communication channels are currently one way. So we have to create a time tick bomb to bail the results out.  To implement a channel logic to close them when the puzzle is 'solved'.
