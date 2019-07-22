package main

import (
	tl "github.com/JoelOtter/termloop"
	"time"
)

type grid struct {
	*tl.Rectangle
	x, y                    int
	isWall, visited, queued bool
	game                    *game
}

func (g *grid) visit() {
	if g.visited {
		panic("repeated")

	}
	g.SetColor(tl.ColorBlue)
	g.visited = true
	// 上下左右加入队列
	ps := [][2]int{
		{g.x, g.y + 1},
		{g.x, g.y - 1},
		{g.x - 1, g.y},
		{g.x + 1, g.y},
	}
	for _, p := range ps {
		if p[0] >= 0 && p[0] < g.game.col &&
			p[1] >= 0 && p[1] < g.game.row {
			gr := g.game.gridMap[p[1]][p[0]]
			//log.Println("gr", gr)
			if !gr.visited && !gr.isWall && !gr.queued {
				time.Sleep(time.Second)
				g.game.exploreQueue.push(gr)
				gr.queued = true
				gr.SetColor(tl.ColorGreen)
			}
		} else {
		}
	}

}

type game struct {
	mazeMap [][]int
	gridMap [][]*grid
	col,
	row int
	exploreQueue queue
}

func (gm *game) start() {
	gm.exploreQueue.push(gm.gridMap[0][0])
	for !gm.exploreQueue.isEmpty() {
		gr := (*grid)(gm.exploreQueue.pop())
		if gr.x == gm.col-1 && gr.y == gm.row-1 {
			break
		}
		gr.visit()
		time.Sleep(time.Second)
	}
	//log.Println("finish")
}

func newGame(mp [][]int) *game {
	return &game{
		mazeMap: mp,
		gridMap: [][]*grid{},
		row:     len(mp),
		col:     len(mp[0]),
	}
}

func newGrid(x, y int, wall bool, game2 *game) *grid {
	g := &grid{
		x:         x,
		y:         y,
		Rectangle: tl.NewRectangle(x*2, y, 2, 1, tl.ColorDefault),
		isWall:    wall,
		visited:   false,
		game:      game2,
	}
	if wall {
		g.SetColor(tl.ColorRed)
	}
	return g
}

func main() {
	mp := [][]int{
		{0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{1, 1, 1, 0, 0},
		{0, 1, 0, 0, 1},
		{0, 1, 0, 0, 0},
	}
	gamee := newGame(mp)
	g := tl.NewGame()
	g.Screen().SetFps(60)
	l := tl.NewBaseLevel(tl.Cell{
		//Bg: tl.ColorWhite,
	})

	drawMap(gamee, g.Screen())
	go func() {
		time.Sleep(1 * time.Second)
		g.SetDebugOn(true)
		g.Log("344")
		gamee.start()

	}()
	g.Screen().SetLevel(l)
	//g.Screen().AddEntity(tl.NewFpsText(0, 0, tl.ColorRed, tl.ColorDefault, 0.5))

	g.Start()
}

func drawMap(gm *game, s *tl.Screen) {
	for r := 0; r < gm.row; r++ {
		gm.gridMap = append(gm.gridMap, []*grid{})
		for c := 0; c < gm.col; c++ {
			gr := newGrid(c, r, gm.mazeMap[r][c] == 1, gm)
			s.AddEntity(gr)
			gm.gridMap[r] = append(gm.gridMap[r], gr)
		}
	}
}

type point *grid

type queue []point

func (q *queue) isEmpty() bool {
	return len(*q) == 0
}

func (q *queue) pop() point {
	if q.isEmpty() {
		panic("queue is empty")
	}
	p := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]
	return p
}

func (q *queue) push(p point) {
	*q = append(*q, p)
	if len(*q) == 1 {
		return
	}
	_ = append((*q)[1:1], (*q)[0:len(*q)-1]...)
	(*q)[0] = p
}
