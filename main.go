package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Maze struct {
	in       []Cell
	out      []Cell
	frontier []Cell
}

func NewMaze(width, height int) *Maze {
	maze := &Maze{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			maze.out = append(maze.out, Cell{x: x, y: y})
		}
	}
	return maze
}

func (m Maze) draw() {
	for _, c := range m.in {
		c.draw()
	}
}

func (m *Maze) expand(cell *Cell) {
	for _, c := range m.in {
		if (c.x == cell.x+1 && c.y == cell.y) ||
			(c.x == cell.x-1 && c.y == cell.y) ||
			(c.x == cell.x && c.y == cell.y+1) ||
			(c.x == cell.x && c.y == cell.y-1) {
			cell.passages = append(cell.passages, c)
			break
		}
	}
	m.in = append(m.in, *cell)

	var newOut []Cell
	for _, c := range m.out {
		if (c.x == cell.x+1 && c.y == cell.y) ||
			(c.x == cell.x-1 && c.y == cell.y) ||
			(c.x == cell.x && c.y == cell.y+1) ||
			(c.x == cell.x && c.y == cell.y-1) {
			m.frontier = append(m.frontier, c)
		} else {
			newOut = append(newOut, c)
		}
	}
	m.out = newOut
}

func (m *Maze) step() {
	if len(m.out)+len(m.frontier) == 0 {
		return
	}

	if len(m.frontier) == 0 {
		n := rand.Intn(len(m.out))
		cell := m.out[n]
		m.out = append(m.out[:n], m.out[n+1:]...)
		m.expand(&cell)
		return
	}

	n := rand.Intn(len(m.frontier))
	cell := m.frontier[n]
	m.frontier = append(m.frontier[:n], m.frontier[n+1:]...)
	m.expand(&cell)
}

type Cell struct {
	x, y     int
	passages []Cell
}

func (c Cell) draw() {
	termbox.SetCell(c.x*4, c.y*2, ' ', termbox.ColorDefault, termbox.ColorYellow)
	termbox.SetCell(c.x*4+1, c.y*2, ' ', termbox.ColorDefault, termbox.ColorYellow)
	for _, pc := range c.passages {
		xoff := pc.x - c.x
		yoff := pc.y - c.y
		termbox.SetCell(c.x*4+xoff*2, c.y*2+yoff, ' ', termbox.ColorDefault, termbox.ColorYellow)
		if xoff == 0 {
			termbox.SetCell(c.x*4+1, c.y*2+yoff, ' ', termbox.ColorDefault, termbox.ColorYellow)
		} else {
			termbox.SetCell(c.x*4+xoff*2+1, c.y*2+yoff, ' ', termbox.ColorDefault, termbox.ColorYellow)
		}
	}
}

func init() {
	rand.Seed(int64(time.Now().Second()))
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	w, h := termbox.Size()
	maze := NewMaze(w/4+1, h/2+1)

	go func() {
		for {
			maze.step()
			time.Sleep(10 * time.Millisecond)
			maze.draw()
			termbox.Flush()
		}
	}()

	termbox.PollEvent()
}
