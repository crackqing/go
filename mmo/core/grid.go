package core

import (
	"fmt"
	"sync"
)

//AOI 2二维数组的显示
type Grid struct{
	GID int

	MinX int //格子左边界坐标
	MaxX int //格子右边界坐标
	MinY int // 格子的上边边界
	MaxY int //格子的下边边界

	playerIDs map[int]bool

	pIDLock sync.RWMutex
}

func NewGrid(gId int ,minX int ,maxX int,minY int,maxY int) *Grid {
	return &Grid{
		GID :gId,
		MinX :minX,
		MaxX :maxX,
		MinY :minY,
		MaxY :maxY,
		playerIDs : make(map[int]bool),
	}
}

func (g *Grid) Add (playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

func (g *Grid) Remove  (playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs,playerID)
}


func (g *Grid) GetPlayerIds () (playerIDs []int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs,k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid is %d,minX %d, maxX %d,minY %d,maxY %d playerIDs %v",
			g.GID,g.MinX,g.MaxX,g.MinY,g.MaxY,g.playerIDs)
}
