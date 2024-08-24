package cmd

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

var (
	board      [][]bool
	shouldExit = false
	events     = make(chan tcell.Event)
	quit       = make(chan struct{})
	rootCmd    = &cobra.Command{
		Use:   "life",
		Short: "Conway's Game Of Life",
		Long:  "Conway's Game Of Life",
		Run:   liveYourLife,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "path to file containing the initial condition")
	rootCmd.Flags().IntP("sizex", "x", 150, "size of board in X direction")
	rootCmd.Flags().IntP("sizey", "y", 50, "size of board in Y direction")
	rootCmd.Flags().IntP("delay", "d", 5, "delay between each iteration in msec")
	rootCmd.Flags().IntP("steps", "s", 1000,
		"max iterations of the game. it will end early if steady state is reached")

}

func getArgs(cmd *cobra.Command) (int, int, int, int, string) {
	steps, err := cmd.Flags().GetInt("steps")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	sizeX, err := cmd.Flags().GetInt("sizex")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	sizeY, err := cmd.Flags().GetInt("sizey")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	delay, err := cmd.Flags().GetInt("delay")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	filename, err := cmd.Flags().GetString("input")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return steps, sizeX, sizeY, delay, filename
}

func listenForExit() {
	for {
		select {
		case <-quit:
			shouldExit = true
		case ev, ok := <-events:
			if ok {
				switch ev := ev.(type) {
				case *tcell.EventKey:
					switch ev.Key() {
					case tcell.KeyEscape, tcell.KeyCtrlC:
						shouldExit = true
					}
				}
			}
		}
	}
}

func liveYourLife(cmd *cobra.Command, args []string) {
	steps, sizeX, sizeY, delay, filename := getArgs(cmd)
	board = getBoard(sizeX, sizeY)
	populateBoard(board, filename)

	// tcell is sort of like curses in python or C, but for Go
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	// screen.EnableMouse()
	defer func() { screen.Fini() }()

	// in the background we listen for events that mean it's
	// time to exit. examples are user presses ESC or CTRL-C
	go screen.ChannelEvents(events, quit)
	go listenForExit()

	// here the game starts. keep going until the max or steady state is reached
	drawBorder(screen, board)
	drawBoard(screen, board)
	for i := 0; i < steps; i++ {
		// sleep a bit between generations, otherwise things end too fast!
		time.Sleep(time.Duration(delay) * time.Millisecond)
		keepGoing := iterate(board)
		drawBoard(screen, board)
		if !keepGoing {
			break
		}

		if shouldExit {
			break
		}
	}
	// all done, sleep a bit so user can look at the majesty
	time.Sleep(1000 * time.Millisecond)
}
