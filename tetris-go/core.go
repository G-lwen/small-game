package main

import (
	"math/rand"
	"time"
)

type gameState int

// 游戏状态
const (
	gamePrepare gameState = iota
	gamePause
	gameStarted
	gameOver
)

const (
	PANEL_HEIGHT = 20 // 游戏面板的高度
	PANEL_WEIGHT = 10 // 游戏面板的宽度

	DEFAULT_SEED = 500 * time.Millisecond // 游戏默认速度
)

// 方块：T、S、Z、J、L、I、O的x, y轴座标
var pieceX = [][]int{
	{0, 1, -1, 0},
	{0, 1, -1, -1},
	{0, 1, -1, 1},
	{0, -1, 1, 0},
	{0, 1, -1, 0},
	{0, 1, -1, -2},
	{0, 1, 1, 0},
}
var pieceY = [][]int{
	{0, 0, 0, 1},
	{0, 0, 0, 1},
	{0, 0, 0, 1},
	{0, 0, 1, 1},
	{0, 0, 1, 1},
	{0, 0, 0, 0},
	{0, 0, 1, 1},
}

// 游戏
type tetrisGame struct {
	panel [][]int // 游戏面板

	x        int   // 方块的 x 轴方位
	y        int   // 方块的 y 轴方位
	pType    int   // 方块类型
	originX  []int // 最初的 4 个格的 x 座标
	originY  []int // 最初的 4 个格的 y 座标
	prepareX []int // 备用的 4 个格的 x 座标
	prepareY []int // 备用的 4 个格的 y 座标

	state gameState // 游戏状态
	score int       // 用户得分

	autoFallCheck *time.Timer // 自动装置
}

// 初始化容器
func (tg *tetrisGame) init() {
	tg.panel = make([][]int, PANEL_HEIGHT)

	for y := 0; y < PANEL_HEIGHT; y++ {
		tg.panel[y] = make([]int, PANEL_WEIGHT)

		for x := 0; x < PANEL_WEIGHT; x++ {
			tg.panel[y][x] = 0
		}
	}

	tg.pType = 0
	tg.x = 0
	tg.y = 0
	tg.originX = []int{0, 0, 0, 0}
	tg.originY = []int{0, 0, 0, 0}
	tg.prepareX = []int{0, 0, 0, 0}
	tg.prepareY = []int{0, 0, 0, 0}

	tg.state = gamePrepare
	tg.score = 0

	tg.autoFallCheck = time.NewTimer(time.Hour)
	tg.autoFallCheck.Stop()
}

// 创建俄罗斯方块游戏
func newTetrisGame() *tetrisGame {
	tg := new(tetrisGame)
	tg.init()
	return tg
}

// 方块左移
func (tg *tetrisGame) leftMovePiece() {
	if tg.state != gameStarted {
		return
	}

	if tg.isFillPiece(tg.x-1, tg.y) {
		tg.fillPiece(0)
		tg.x--
		tg.fillPiece(tg.pType + 1)
	}
}

// 方块右移
func (tg *tetrisGame) rightMovePiece() {
	if tg.state != gameStarted {
		return
	}

	if tg.isFillPiece(tg.x+1, tg.y) {
		tg.fillPiece(0)
		tg.x++
		tg.fillPiece(tg.pType + 1)
	}
}

// 方块下落
func (tg *tetrisGame) downMovePiece() bool {
	if tg.state == gameStarted && tg.isFillPiece(tg.x, tg.y+1) {
		tg.fillPiece(0)
		tg.y++
		tg.fillPiece(tg.pType + 1)
		return true
	}
	return false
}

// 方块旋转
func (tg *tetrisGame) rotatePiece() {
	if tg.state != gameStarted {
		return
	}

	for i := 0; i < 4; i++ {
		tg.prepareX[i] = tg.originY[i]
		tg.prepareY[i] = -tg.originX[i]
	}

	if tg.isFillPiece(tg.x, tg.y) {
		tg.fillPiece(0)
		for i := 0; i < 4; i++ {
			tg.originX[i] = tg.prepareX[i]
			tg.originY[i] = tg.prepareY[i]
		}
		tg.fillPiece(tg.pType + 1)
	} else {
		for i := 0; i < 4; i++ {
			tg.prepareX[i] = tg.originX[i]
			tg.prepareY[i] = tg.originY[i]
		}
	}
}

// 获取方块
func (tg *tetrisGame) getPiece() bool {
	rand.Seed(time.Now().Unix())
	tg.pType = rand.Intn(7)

	tg.x = PANEL_WEIGHT/2 - 1
	tg.y = 0

	for i := 0; i < 4; i++ {
		tg.originX[i] = pieceX[tg.pType][i]
		tg.originY[i] = pieceY[tg.pType][i]

		tg.prepareX[i] = tg.originX[i]
		tg.prepareY[i] = tg.originY[i]
	}

	if tg.isFillPiece(tg.x, tg.y) {
		fillSign := tg.pType + 1
		tg.fillPiece(fillSign)

		return true
	}
	return false
}

// 判断方块是否可以在 x y 方向移动填充到面板
func (tg *tetrisGame) isFillPiece(x, y int) bool {
	for i := 0; i < 4; i++ {
		prepareX := x + tg.prepareX[i]
		prepareY := y + tg.prepareY[i]
		if prepareX < 0 || prepareX >= PANEL_WEIGHT || prepareY >= PANEL_HEIGHT {
			return false
		}
		if prepareY > -1 && tg.panel[prepareY][prepareX] > 0 {
			return false
		}
	}
	return true
}

// 把方格填充到面板或者擦除上一次方格移动的痕迹
func (tg *tetrisGame) fillPiece(sign int) {
	for i := 0; i < 4; i++ {
		x := tg.x + tg.originX[i]
		y := tg.y + tg.originY[i]
		if 0 <= y && y < PANEL_HEIGHT && 0 <= x && x < PANEL_WEIGHT && tg.panel[y][x] != -sign {
			tg.panel[y][x] = -sign
		}
	}
}

// 方格无法继续向下移动，被固定
func (tg *tetrisGame) fixPiece() {
	tg.signPiece()
	tg.removeLine()
	if tg.getPiece() {
		tg.autoFallCheck.Reset(tg.seed())
	} else {
		tg.state = gameOver
	}
}

// 标记被固定的方块
func (tg *tetrisGame) signPiece() {
	for i := 0; i < 4; i++ {
		x := tg.x + tg.originX[i]
		y := tg.y + tg.originY[i]
		if 0 <= y && y < PANEL_HEIGHT && 0 <= x && x < PANEL_WEIGHT {
			tg.panel[y][x] = tg.pType + 1
		}
	}
}

// 移除被填充满的行
func (tg *tetrisGame) removeLine() {
	for y := 0; y < PANEL_HEIGHT; y++ {
		isRemove := true
		for x := 0; x < PANEL_WEIGHT; x++ {
			if tg.panel[y][x] == 0 {
				isRemove = false
			}
		}

		if isRemove {
			tg.removeConfirmedLine(y)
			tg.score += 10
		}
	}
}

// 根据要行号进行移除
func (tg *tetrisGame) removeConfirmedLine(line int) {
	for y := line; y > 0; y-- {
		for x := 0; x < PANEL_WEIGHT; x++ {
			tg.panel[y][x] = tg.panel[y-1][x]
		}
	}
	for x := 0; x < PANEL_WEIGHT; x++ {
		tg.panel[0][x] = 0
	}
}

// 方块块自动落下
func (tg *tetrisGame) autoFall() {
	if tg.downMovePiece() {
		tg.autoFallCheck.Reset(tg.seed())
	} else {
		tg.fixPiece()
	}
}

// 游戏速度
func (tg *tetrisGame) seed() time.Duration {
	return DEFAULT_SEED
}

// 开始/暂停游戏
func (tg *tetrisGame) startOrPause() {
	if tg.state == gamePrepare {
		tg.state = gameStarted
		tg.getPiece()
		tg.autoFallCheck.Reset(tg.seed())
	} else if tg.state == gameStarted {
		tg.state = gamePause
		tg.autoFallCheck.Stop()
	} else if tg.state == gamePause {
		tg.state = gameStarted
		tg.autoFallCheck.Reset(tg.seed())
	}
}
