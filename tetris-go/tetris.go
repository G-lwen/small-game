package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	GAME_PIECE      = " @ "
	GAME_BACKGROUND = "   "
)

// 清屏
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Start()
}

// 绘画图像
func paint(tg *tetrisGame) {
	face := "\n Weclcome to Tetris Game \n"
	face += "\n #  #  #  #  #  #  #  #  #  #  #  # \n"

	for y := 0; y < len(tg.panel); y++ {
		face += " # "

		for x := 0; x < len(tg.panel[y]); x++ {
			if tg.panel[y][x] == 0 {
				face += GAME_BACKGROUND
			} else {
				face += GAME_PIECE
			}
		}

		if y == 2 {
			face += " #                游戏状态："
			if tg.state == gamePrepare {
				face += " 准备中 \n"
			} else if tg.state == gameStarted {
				face += " 进行中 \n"
			} else if tg.state == gamePause {
				face += " 暂停中 \n"
			} else {
				face += " 游戏结束 \n"
			}
		} else if y == 6 {
			face += " #                游戏得分："
			face += strconv.Itoa(tg.score)
			face += "\n"
		} else if y == 10 {
			face += " #                [ 开始/暂停游戏：「s」键，重置游戏：「r」键，退出游戏：「Esc」键，操作游戏：「方向键」]\n"
		} else {
			face += " # \n"
		}
	}

	face += " #  #  #  #  #  #  #  #  #  #  #  # "
	fmt.Println(face)
}

// 运行游戏
func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)

	tg := newTetrisGame()
	paint(tg)

	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Key == termbox.KeyArrowUp:
					tg.rotatePiece()
				case ev.Key == termbox.KeyArrowDown:
					tg.downMovePiece()
				case ev.Key == termbox.KeyArrowLeft:
					tg.leftMovePiece()
				case ev.Key == termbox.KeyArrowRight:
					tg.rightMovePiece()
				case ev.Ch == 's':
					tg.startOrPause()
				case ev.Ch == 'r':
					tg.init()
				case ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD:
					return
				}
			}
		case <-tg.autoFallCheck.C:
			tg.autoFall()
		default:
			time.Sleep(10 * time.Millisecond)
			clearScreen()
			time.Sleep(10 * time.Millisecond)
			paint(tg)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
