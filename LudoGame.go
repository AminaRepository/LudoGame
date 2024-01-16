package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

// used to check correct colour
func notInArray(value string, array []string) bool {
	for _, element := range array {
		if element == value {
			return false
		}
	}
	return true
}

type Cell struct {
	Type     int  // Default value is 0 (zero value for int)
	HasPiece bool // Default value is false (zero value for bool)
	// 0 default 1 green ♣ 2 yellow ♦ 3 red ♥ 4 blue ♠ 5 green safe 6 blue safe 7 yellow safe 8 red safe 9 path
}

type Piece struct {
	Colour      int // 0 default 1 green 2 yellow 3 red 4 blue
	InGame      bool
	InSafehouse int // checks on which safehouse cell the piece is, 0 if none
	Row         int
	Column      int
}

type Player struct {
	Colour      int // colour of the player 0 default 1 green 2 yellow 3 red 4 blue
	Index       int // index of the player (number in the order)
	PieceInGame bool
}

func makePiece(colour int, row int, column int) Piece {
	return Piece{
		Colour:      colour,
		InGame:      false,
		InSafehouse: 0,
		Row:         row,
		Column:      column,
	}
}

func makePlayer(colour int, index int, ingame bool) Player {
	return Player{
		Colour:      colour,
		Index:       index,
		PieceInGame: ingame,
	}
}

func throwDice() int {
	return rand.Intn(6) + 1
}

// we use a logger to save the progress of the game
var logger *log.Logger

var board [][]Cell          // board
var pieces []*Piece         // list of pieces (players)
var playablePieces []*Piece // keeps track of pieces the player can move
var playerSet []string      // player set up
var playerTrack []*Player   // properly ordered player list (green, yellow, red, blue)
var rollSixError = false    // error showcases piece tried to leave house while not rolling six
var rollAgain = false       // piece rolled six and player playes again
var rolledAgainOnce = false // player already rolled a six and played twice, can not play again
var stepOnSelf = false      // piece stepped on itself and needs to go back, player plays again
var piecesTried = 0         // how many pieces available the player tried to move but was unable to
var foundPiece = false      // used in loop in movePlayer, makes sure there are no errors in loop
var currentThrow = 0        // value used in main function, established due to use in other functions
var movingPiece = 0         // value used in main function, established due to use in other functions
var playablePiecesInitialized = false
var playerOut = false

func printBoard(input int) {
	if !rollSixError && !stepOnSelf {
		playablePieces = playablePieces[:0]
		for i := range playablePieces {
			playablePieces = append(playablePieces[:i], playablePieces[i+1:]...)
		}
		logger.Printf("printBoard: playablePieces = []")
		playablePiecesInitialized = false
		logger.Printf("printBoard: playablePiecesInitialized = false")
	}
	logger.Printf("Print board, case %d", input)
	num := 1
	switch input {
	case 0:
		for i := 0; i < 11; i++ {
			for j := 0; j < 11; j++ {
				if board[i][j].HasPiece {
					for _, piece := range pieces {
						if piece.Row == i && piece.Column == j {
							switch piece.Colour {
							case 1:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♣ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 2:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♦ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 3:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♥ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 4:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♠ ")
								} else {
									fmt.Printf(" ○ ")
								}
							}
						}
					}
				} else {
					switch board[i][j].Type {
					case 9, 1, 2, 3, 4:
						fmt.Printf(" ○ ")
					case 5, 6, 7, 8:
						if board[i][j].HasPiece {
							for _, piece := range pieces {
								if piece.Row == i && piece.Column == j {
									switch piece.Colour {
									case 1:
										fmt.Printf(" ♣ ")
									case 2:
										fmt.Printf(" ♦ ")
									case 3:
										fmt.Printf(" ♥ ")
									case 4:
										fmt.Printf(" ♠ ")
									}
								}
							}
						} else {
							fmt.Printf(" ● ")
						}
					default:
						fmt.Printf("   ")
					}
				}
			}
			fmt.Printf("\n")
		}
	// case for when green player chooses piece
	case 1:
		for i := 0; i < 11; i++ {
			for j := 0; j < 11; j++ {
				if board[i][j].HasPiece {
					for _, piece := range pieces {
						if piece.Row == i && piece.Column == j {
							switch piece.Colour {
							case 1:
								if piece.InGame {
									fmt.Printf(" %d ", num)
									num = num + 1
									if !playablePiecesInitialized {
										playablePieces = append(playablePieces, piece)
									}
								} else {
									fmt.Printf(" ♣ ")
								}
							case 2:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♦ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 3:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♥ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 4:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♠ ")
								} else {
									fmt.Printf(" ○ ")
								}
							}
						}
					}
				} else {
					switch board[i][j].Type {
					case 9, 1, 2, 3, 4:
						fmt.Printf(" ○ ")
					case 5, 6, 7, 8:
						fmt.Printf(" ● ")
					default:
						fmt.Printf("   ")
					}
				}
			}
			fmt.Printf("\n")
		}
	// case for when yellow player chooses piece
	case 2:
		for i := 0; i < 11; i++ {
			for j := 0; j < 11; j++ {
				if board[i][j].HasPiece {
					for _, piece := range pieces {
						if piece.Row == i && piece.Column == j {
							switch piece.Colour {
							case 1:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♣ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 2:
								if piece.InGame {
									fmt.Printf(" %d ", num)
									num = num + 1
									if !playablePiecesInitialized {
										playablePieces = append(playablePieces, piece)
									}
								} else {
									fmt.Printf(" ♦ ")
								}
							case 3:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♥ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 4:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♠ ")
								} else {
									fmt.Printf(" ○ ")
								}
							}
						}
					}
				} else {
					switch board[i][j].Type {
					case 9, 1, 2, 3, 4:
						fmt.Printf(" ○ ")
					case 5, 6, 7, 8:
						fmt.Printf(" ● ")
					default:
						fmt.Printf("   ")
					}
				}
			}
			fmt.Printf("\n")
		}
	// case for when red player chooses piece
	case 3:
		for i := 0; i < 11; i++ {
			for j := 0; j < 11; j++ {
				if board[i][j].HasPiece {
					for _, piece := range pieces {
						if piece.Row == i && piece.Column == j {
							switch piece.Colour {
							case 1:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♣ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 2:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♦ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 3:
								if piece.InGame {
									fmt.Printf(" %d ", num)
									num = num + 1
									if !playablePiecesInitialized {
										playablePieces = append(playablePieces, piece)
									}
								} else {
									fmt.Printf(" ♥ ")
								}
							case 4:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♠ ")
								} else {
									fmt.Printf(" ○ ")
								}
							}
						}
					}
				} else {
					switch board[i][j].Type {
					case 9, 1, 2, 3, 4:
						fmt.Printf(" ○ ")
					case 5, 6, 7, 8:
						fmt.Printf(" ● ")
					default:
						fmt.Printf("   ")
					}
				}
			}
			fmt.Printf("\n")
		}
	// case for when blue player chooses piece
	case 4:
		for i := 0; i < 11; i++ {
			for j := 0; j < 11; j++ {
				if board[i][j].HasPiece {
					for _, piece := range pieces {
						if piece.Row == i && piece.Column == j {
							switch piece.Colour {
							case 1:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♣ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 2:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♦ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 3:
								if piece.InGame || piece.InSafehouse != 0 {
									fmt.Printf(" ♥ ")
								} else {
									fmt.Printf(" ○ ")
								}
							case 4:
								if piece.InGame {
									fmt.Printf(" %d ", num)
									num = num + 1
									if !playablePiecesInitialized {
										playablePieces = append(playablePieces, piece)
									}
								} else {
									fmt.Printf(" ♠ ")
								}
							}
						}
					}
				} else {
					switch board[i][j].Type {
					case 9, 1, 2, 3, 4:
						fmt.Printf(" ○ ")
					case 5, 6, 7, 8:
						fmt.Printf(" ● ")
					default:
						fmt.Printf("   ")
					}
				}
			}
			fmt.Printf("\n")
		}
	}
	logger.Printf("printBoard: playablePieces %v", playablePieces)
}

func movePiece(piece *Piece, dice int, pieces []*Piece, board [][]Cell) {
	logger.Printf("movePiece: roll = %d", dice)
	if dice == 6 {
		logger.Printf("movePiece: rolledAgain = true")
		rollAgain = true // value to cause player to roll again if they rolled a six
	}

	board[piece.Row][piece.Column].HasPiece = false
	startRow := piece.Row
	startColumn := piece.Column // save start row & column in case player is set back

	if board[piece.Row][piece.Column].Type == 1 || board[piece.Row][piece.Column].Type == 2 || board[piece.Row][piece.Column].Type == 3 || board[piece.Row][piece.Column].Type == 4 {
		if dice == 6 {
			switch piece.Colour {
			case 1:
				piece.Row = 4
				piece.Column = 0
			case 2:
				piece.Row = 0
				piece.Column = 6
			case 3:
				piece.Row = 6
				piece.Column = 10
			case 4:
				piece.Row = 10
				piece.Column = 4
			}
		} else {
			fmt.Printf("You can only exit the house if you roll a six!\n")
			rollSixError = true
			piecesTried = piecesTried + 1
			logger.Printf("movePiece: rollSixError = true")
			logger.Printf("movePiece: piecesTried += 1")
		}
	} else {
		for i := 0; i < dice; i++ {
			switch piece.Row {
			// we check rows to find out the players next move
			case 4:
				switch piece.Column {
				// on certain columns player moves "right" along the columns
				// otherwise they move "up" or "down" along the rows
				case 0, 1, 2, 3, 6, 7, 8, 9:
					piece.Column = piece.Column + 1 // move right
				case 4:
					piece.Row = piece.Row - 1 // move up
				case 10:
					piece.Row = piece.Row + 1 // move down
				}
			case 5:
				switch piece.Column {
				// in row five player moves "up" along the rows on column 0 or "down" on column 10
				case 0:
					piece.Row = piece.Row - 1
				case 10:
					piece.Row = piece.Row + 1
				}
			case 6:
				switch piece.Column {
				// row six works like row four, but opposite
				case 0:
					piece.Row = piece.Row - 1 // player moves "up" along the rows
				case 1, 2, 3, 4, 7, 8, 9, 10:
					piece.Column = piece.Column - 1 // player moves "left" along the columns
				case 6:
					piece.Row = piece.Row + 1 // player moves "down" along the rows
				}
			default:
				// if player is not on row 4-6 we check columns
				switch piece.Column {
				case 4:
					switch piece.Row {
					case 0:
						piece.Column = piece.Column + 1 // player moves "right" along the columns
					case 1, 2, 3, 4, 7, 8, 9, 10:
						piece.Row = piece.Row - 1 // player moves "up" along the rows
					case 6:
						piece.Column = piece.Column - 1 // player moves "left" along the columns
					}
				case 5:
					switch piece.Row {
					case 0:
						piece.Column = piece.Column + 1 // player moves "right" along the columns
					case 10:
						piece.Column = piece.Column - 1 // player moves "left" along the columns
					}
				case 6:
					switch piece.Row {
					case 0, 1, 2, 3, 6, 7, 8, 9:
						piece.Row = piece.Row + 1 // player moves "down" along the rows
					case 4:
						piece.Column = piece.Column + 1 // player moves "right" along the columns
					case 10:
						piece.Column = piece.Column - 1 // player moves "left" along the column
					}
				}
			}
			// we check if player is in front of the safehouse
			// if so we move it in as far as possible with the remaining steps (dice roll - amount of steps taken already)
			switch piece.Colour {
			case 1:
				if piece.Row == 5 && piece.Column == 0 {
					piece.Column = dice - i
					if piece.Column > 4 {
						piece.Column = 4
					}
					for board[piece.Row][piece.Column].HasPiece == true {
						piece.Column = piece.Column - 1
					}
					piece.InGame = false
				}
			case 2:
				if piece.Row == 0 && piece.Column == 5 {
					piece.Row = dice - i
					if piece.Row > 4 {
						piece.Row = 4
					}
					for board[piece.Row][piece.Column].HasPiece == true {
						piece.Row = piece.Row - 1
					}
					piece.InGame = false
				}
			case 3:
				if piece.Row == 5 && piece.Column == 10 {
					piece.Column = 10 - (dice - i)
					if piece.Column < 6 {
						piece.Column = 6
					}
					for board[piece.Row][piece.Column].HasPiece == true {
						piece.Column = piece.Column + 1
					}
					piece.InGame = false
				}
			case 4:
				if piece.Row == 10 && piece.Column == 5 {
					piece.Row = 10 - (dice - i)
					if piece.Row < 6 {
						piece.Row = 6
					}
					for board[piece.Row][piece.Column].HasPiece == true {
						piece.Row = piece.Row + 1
					}
					piece.InGame = false
				}
			}
		}
	}

	if board[piece.Row][piece.Column].HasPiece == true {
		// if we start searching for player, we set foundPlayer as false
		// prevents error of repeating function twice with the same or different players
		logger.Printf("movePiece: foundPIece = false")
		foundPiece = false
	outerLoop:
		for _, pl := range pieces {
			if pl.Row == piece.Row && pl.Column == piece.Column && &pl != &piece {
				foundPiece = true
				logger.Printf("movePiece: foundPiece = true")
				if pl.Colour == piece.Colour {
					// if the player steps on its own colour we move the piece back to its original cell
					// we set stepOnSelf as true to signal that the turn must be repeated with a different player
					piece.Row = startRow
					piece.Column = startColumn
					fmt.Printf("You can not step on your own pieces!\n")
					stepOnSelf = true
					piecesTried = piecesTried + 1
					logger.Printf("movePiece: stepOnSelf = true")
					logger.Printf("movePiece: piecesTried += 1")
				} else {
					switch pl.Colour {
					case 1:
						if board[0][0].HasPiece == false {
							pl.Row = 0
							pl.Column = 0
							board[0][0].HasPiece = true
						} else if board[0][1].HasPiece == false {
							pl.Row = 0
							pl.Column = 1
							board[0][1].HasPiece = true
						} else if board[1][0].HasPiece == false {
							pl.Row = 1
							pl.Column = 0
							board[1][0].HasPiece = true
						} else if board[1][1].HasPiece == false {
							pl.Row = 1
							pl.Column = 1
							board[1][1].HasPiece = true
						}
					case 2:
						if board[0][9].HasPiece == false {
							pl.Row = 0
							pl.Column = 9
							board[0][9].HasPiece = true
						} else if board[0][10].HasPiece == false {
							pl.Row = 0
							pl.Column = 10
							board[0][10].HasPiece = true
						} else if board[1][9].HasPiece == false {
							pl.Row = 1
							pl.Column = 9
							board[1][9].HasPiece = true
						} else if board[1][10].HasPiece == false {
							pl.Row = 1
							pl.Column = 10
							board[1][10].HasPiece = true
						}
					case 3:
						if board[9][9].HasPiece == false {
							pl.Row = 9
							pl.Column = 9
							board[9][9].HasPiece = true
						} else if board[9][10].HasPiece == false {
							pl.Row = 9
							pl.Column = 10
							board[9][10].HasPiece = true
						} else if board[10][9].HasPiece == false {
							pl.Row = 10
							pl.Column = 9
							board[10][9].HasPiece = true
						} else if board[10][10].HasPiece == false {
							pl.Row = 10
							pl.Column = 10
							board[10][10].HasPiece = true
						}
					case 4:
						if board[9][0].HasPiece == false {
							pl.Row = 9
							pl.Column = 0
							board[9][0].HasPiece = true
						} else if board[9][1].HasPiece == false {
							pl.Row = 9
							pl.Column = 1
							board[9][1].HasPiece = true
						} else if board[10][0].HasPiece == false {
							pl.Row = 10
							pl.Column = 0
							board[10][0].HasPiece = true
						} else if board[10][1].HasPiece == false {
							pl.Row = 10
							pl.Column = 1
							board[10][1].HasPiece = true
						}
					}
				}
			}
			if foundPiece == true {
				break outerLoop
			}
		}
	}
	board[piece.Row][piece.Column].HasPiece = true
}

func playRound(colour int) {
	logger.Printf("playRound for %d", colour)
	switch colour {
	case 1:
		printBoard(1)
	case 2:
		printBoard(2)
	case 3:
		printBoard(3)
	case 4:
		printBoard(4)
	}
	movingPiece = 0
	fmt.Scanln(&movingPiece) // player enters which piece to move
	// we look at the list of playable pieces and enter the appropriate index
	// pieces are numbered 1-4, playablePiece index is 0-3, the index is movingPiece - 1
	movePiece(playablePieces[movingPiece-1], currentThrow, pieces, board)
	logger.Printf("playRound: movePiece")
	if !rollSixError && !stepOnSelf {
		playablePieces = playablePieces[:0]
		for i := range playablePieces {
			playablePieces = append(playablePieces[:i], playablePieces[i+1:]...)
		}
		logger.Printf("playRound: playablePieces = []")
	}
}

func main() {
	f, err := os.OpenFile(".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	logger = log.New(f, "[Game] ", log.Ldate|log.Ltime)
	logger.Printf("\n\n\n\n- - - - - - - - - - - - - -LOG START")

	// initialize needed values

	rollSixError = false
	rollAgain = false
	rolledAgainOnce = false
	stepOnSelf = false
	logger.Printf("rollSixError = false")
	logger.Printf("rollAgain = false")
	logger.Printf("rolledAgainOnce = false")
	logger.Printf("stepOnSelf = false")

	board = make([][]Cell, 11)
	for i := range board {
		board[i] = make([]Cell, 11)
	}

	board[0][0].Type = 1
	board[0][1].Type = 1
	board[0][4].Type = 9
	board[0][5].Type = 9
	board[0][6].Type = 9
	board[0][9].Type = 3
	board[0][10].Type = 3

	board[1][0].Type = 1
	board[1][1].Type = 1
	board[1][4].Type = 9
	board[1][5].Type = 7
	board[1][6].Type = 9
	board[1][9].Type = 3
	board[1][10].Type = 3

	board[2][4].Type = 9
	board[2][5].Type = 7
	board[2][6].Type = 9

	board[3][4].Type = 9
	board[3][5].Type = 7
	board[3][6].Type = 9

	board[4][0].Type = 9
	board[4][1].Type = 9
	board[4][2].Type = 9
	board[4][3].Type = 9
	board[4][4].Type = 9
	board[4][5].Type = 7
	board[4][6].Type = 9
	board[4][7].Type = 9
	board[4][8].Type = 9
	board[4][9].Type = 9
	board[4][10].Type = 9

	board[5][0].Type = 9
	board[5][1].Type = 5
	board[5][2].Type = 5
	board[5][3].Type = 5
	board[5][4].Type = 5
	board[5][6].Type = 8
	board[5][7].Type = 8
	board[5][8].Type = 8
	board[5][9].Type = 8
	board[5][10].Type = 9

	board[6][0].Type = 9
	board[6][1].Type = 9
	board[6][2].Type = 9
	board[6][3].Type = 9
	board[6][4].Type = 9
	board[6][5].Type = 6
	board[6][6].Type = 9
	board[6][7].Type = 9
	board[6][8].Type = 9
	board[6][9].Type = 9
	board[6][10].Type = 9

	board[7][4].Type = 9
	board[7][5].Type = 6
	board[7][6].Type = 9

	board[8][4].Type = 9
	board[8][5].Type = 6
	board[8][6].Type = 9

	board[9][0].Type = 2
	board[9][1].Type = 2
	board[9][4].Type = 9
	board[9][5].Type = 6
	board[9][6].Type = 9
	board[9][9].Type = 4
	board[9][10].Type = 4

	board[10][0].Type = 2
	board[10][1].Type = 2
	board[10][4].Type = 9
	board[10][5].Type = 9
	board[10][6].Type = 9
	board[10][9].Type = 4
	board[10][10].Type = 4

	pl1g := makePiece(1, 0, 0)
	pieces = append(pieces, &pl1g)
	pl2g := makePiece(1, 0, 1)
	pieces = append(pieces, &pl2g)
	pl3g := makePiece(1, 1, 0)
	pieces = append(pieces, &pl3g)
	pl4g := makePiece(1, 1, 1)
	pieces = append(pieces, &pl4g)

	pl1y := makePiece(2, 0, 9)
	pieces = append(pieces, &pl1y)
	pl2y := makePiece(2, 0, 10)
	pieces = append(pieces, &pl2y)
	pl3y := makePiece(2, 1, 9)
	pieces = append(pieces, &pl3y)
	pl4y := makePiece(2, 1, 10)
	pieces = append(pieces, &pl4y)

	pl1r := makePiece(3, 9, 9)
	pieces = append(pieces, &pl1r)
	pl2r := makePiece(3, 9, 10)
	pieces = append(pieces, &pl2r)
	pl3r := makePiece(3, 10, 9)
	pieces = append(pieces, &pl3r)
	pl4r := makePiece(3, 10, 10)
	pieces = append(pieces, &pl4r)

	pl1b := makePiece(4, 9, 0)
	pieces = append(pieces, &pl1b)
	pl2b := makePiece(4, 10, 0)
	pieces = append(pieces, &pl2b)
	pl3b := makePiece(4, 9, 1)
	pieces = append(pieces, &pl3b)
	pl4b := makePiece(4, 10, 1)
	pieces = append(pieces, &pl4b)

	for _, player := range pieces {
		board[player.Row][player.Column].HasPiece = true
	}

	fmt.Printf("Legend: \n")
	fmt.Printf("Green: ♣ \n")
	fmt.Printf("Yellow: ♦ \n")
	fmt.Printf("Red: ♥ \n")
	fmt.Printf("Blue: ♠ \n")
	fmt.Printf("How many players are playing?\n")

	playerNum := 0
	currentTurn := 0

	fmt.Scanln(&playerNum)

	gameRun := true

	fmt.Printf("Choose the colour:\n")
	for i := 1; i <= playerNum; i++ {
		fmt.Printf("Player %d: ", i)
		playerColour := ""
		fmt.Scanln(&playerColour)
		for notInArray(playerColour, []string{"green", "blue", "yellow", "red"}) {
			fmt.Printf("The only acceptable colours are green, blue, red and yellow, try again:\n")
			fmt.Scanln(&playerColour)
		}
		for _, colour := range playerSet {
			if colour == playerColour {
				fmt.Printf("You can not pick the same colour, choose another: \n")
				fmt.Scanln(&playerColour)
			}
		}
		playerSet = append(playerSet, playerColour)
	}

	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			for x, plColour := range playerSet {
				if plColour == "green" {
					Player1 := makePlayer(1, x, false)
					playerTrack = append(playerTrack, &Player1)
				}
			}
		case 1:
			for x, plColour := range playerSet {
				if plColour == "yellow" {
					Player2 := makePlayer(2, x, false)
					playerTrack = append(playerTrack, &Player2)
				}
			}
		case 2:
			for x, plColour := range playerSet {
				if plColour == "red" {
					Player3 := makePlayer(3, x, false)
					playerTrack = append(playerTrack, &Player3)
				}
			}
		case 3:
			for x, plColour := range playerSet {
				if plColour == "blue" {
					Player4 := makePlayer(4, x, false)
					playerTrack = append(playerTrack, &Player4)
				}
			}
		}
	}

	for _, player := range playerTrack {
		logger.Printf("InGame = true for colour %d", player.Colour)
		switch player.Colour {
		case 1:
			pl1g.InGame = true
			pl2g.InGame = true
			pl3g.InGame = true
			pl4g.InGame = true
		case 2:
			pl1y.InGame = true
			pl2y.InGame = true
			pl3y.InGame = true
			pl4y.InGame = true
		case 3:
			pl1r.InGame = true
			pl2r.InGame = true
			pl3r.InGame = true
			pl4r.InGame = true
		case 4:
			pl1b.InGame = true
			pl2b.InGame = true
			pl3b.InGame = true
			pl4b.InGame = true
		}
	}

	fmt.Printf("\nThe order goes green, yellow, red, blue. Start!\n")
	printBoard(0)

	// first round each player throws dice three times
	for _, player := range playerTrack {
		currentTurn = currentTurn + 1
		logger.Printf("Turn: %d", currentTurn)
		fmt.Printf("\nTurn %d\n", currentTurn)
		fmt.Printf("Player %d turn:\n", player.Index+1)
		for j := 1; j <= 3; j++ {
			currentThrow = throwDice()
			playablePiecesInitialized = false
			logger.Printf("playablePiecesInitialized = false")
			if currentThrow == 6 && !player.PieceInGame {
				fmt.Printf("You rolled a %d! Pick a piece to move:\n", currentThrow)
				playRound(player.Colour)
				playablePiecesInitialized = true
				logger.Printf("ln 897: playablePiecesInitialized = true")
				playerOut = true
				player.PieceInGame = true
				logger.Printf("ln 899: playerOut = true")
				if stepOnSelf && (piecesTried != len(playablePieces)) {
					logger.Printf("ln 901: stepOnSelf && (piecesTried != len(playablePieces)) == %t", stepOnSelf && (piecesTried != len(playablePieces)))
					fmt.Printf("Pick a different piece to move:\n")
					playRound(player.Colour)
				}
			} else if player.PieceInGame {
				fmt.Printf("You rolled a %d! Pick a piece to move:\n", currentThrow)
				playRound(player.Colour)

				for rollSixError {
					rollSixError = false
					logger.Printf("ln 911: rollSixError = false")
					fmt.Printf("Pick a different piece:\n")
					playRound(player.Colour)
					logger.Printf("ln 914: player tried to leave house without six, replay")
				}
				rollSixError = false
				logger.Printf("ln 917: rollSixError = false")
			} else {
				fmt.Printf("You rolled a %d!\n", currentThrow)
			}
		}
	}
	rolledAgainOnce = false
	logger.Printf("rolledAgainOnce = false")

	for gameRun {
		for _, player := range playerTrack {
			currentTurn = currentTurn + 1
			logger.Printf("gameRun: Turn: %d", currentTurn)
			fmt.Printf("\nTurn %d\n", currentTurn)
			fmt.Printf("Player %d turn:\n", player.Index+1)
			currentThrow = throwDice()
			piecesTried = 0
			logger.Printf("gameRun: piecesTried = 0")

			if !player.PieceInGame && currentThrow != 6 {
				fmt.Printf("You can only exit house with a 6!\n")
			} else {
				fmt.Printf("You rolled a %d! Pick a piece to move:\n", currentThrow)
				playablePiecesInitialized = false
				logger.Printf("gameRun: playablePiecesInitialized = false")
				playRound(player.Colour)
				playablePiecesInitialized = true
				logger.Printf("gameRun: playablePiecesInitialized = true")
				player.PieceInGame = true

				for rollSixError && (piecesTried != len(playablePieces)) {
					rollSixError = false
					logger.Printf("gameRun: rollSixError = false")
					fmt.Printf("Pick a different piece:\n")
					playRound(player.Colour)
				}
				rollSixError = false
				logger.Printf("gameRun: rollSixError = false")

				if rollAgain && !rolledAgainOnce {
					rollAgain = false
					logger.Printf("gameRun: rollAgain = false")
					rolledAgainOnce = true
					logger.Printf("gameRun: rolledAgainOnce = true")
					currentThrow = throwDice()
					fmt.Printf("You rolled a %d! Pick a piece to move:\n", currentThrow)
					playRound(player.Colour)
				}
				rollAgain = false
				rolledAgainOnce = false
				logger.Printf("gameRun: rolledAgainOnce = false")
				logger.Printf("gameRun: rolledAgainOnce = false")

				for stepOnSelf && (piecesTried != len(playablePieces)) {
					stepOnSelf = false
					logger.Printf("gameRun: stepOnSelf = false")
					fmt.Printf("Pick a different piece to move:\n")
					playRound(player.Colour)
				}

				for _, currentPlayer := range playerTrack {
					logger.Printf("gameRun: foundPiece = false")
					foundPiece = false
					switch currentPlayer.Colour {
					case 1:
						for _, player := range pieces {
							if player.Colour == 1 && player.InGame {
								foundPiece = true
								logger.Printf("gameRun: Green piece found, foundPiece = true")
								break
							}
						}
						if !foundPiece {
							gameRun = false
							logger.Printf("gameRun: No green piece found, gameRun = stop")
							fmt.Printf("Congratulations, green won! ♣")
						}
					case 2:
						for _, player := range pieces {
							if player.Colour == 2 && player.InGame {
								foundPiece = true
								logger.Printf("gameRun: Yellow piece found, foundPiece = true")
								break
							}
						}
						if !foundPiece {
							gameRun = false
							logger.Printf("gameRun: No yellow piece found, gameRun = stop")
							fmt.Printf("Congratulations, yellow won! ♦")
						}
					case 3:
						for _, player := range pieces {
							if player.Colour == 3 && player.InGame {
								foundPiece = true
								logger.Printf("gameRun: Red piece found, foundPiece = true")
								break
							}
						}
						if !foundPiece {
							gameRun = false
							logger.Printf("gameRun: No red piece found, gameRun = stop")
							fmt.Printf("Congratulations, red won! ♥")
						}
					case 4:
						for _, player := range pieces {
							if player.Colour == 4 && player.InGame {
								foundPiece = true
								logger.Printf("gameRun: Blue piece found, foundPiece = true")
								break
							}
						}
						if !foundPiece {
							gameRun = false
							logger.Printf("gameRun: No blue piece found, gameRun = stop")
							fmt.Printf("Congratulations, blue won! ♠")
						}
					}
				}
				playablePieces = playablePieces[:0]
				for i := range playablePieces {
					playablePieces = append(playablePieces[:i], playablePieces[i+1:]...)
				}
				logger.Printf("gameRun: playablePieces = []")
				playablePiecesInitialized = false
				logger.Printf("gameRun: playablePiecesInitialized = false")

			}
		}
	}
	logger.Printf("\n- - - - - - - - - - - - - -LOG END")
}
