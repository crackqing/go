package core

import "fmt"

type AOIManager struct {
	MinX int
	MaxX int
	CntsX int

	MinY int
	MaxY int
	CntsY int

	grids map[int] * Grid
}

func NewAOIManager(MinX int,MaxX int,CntsX int,MinY int,MaxY int,CntsY int) *AOIManager{
	aoiMgr := &AOIManager{
		MinX:  MinX,
		MaxX:  MaxX,
		CntsX: CntsX,
		MinY:  MinY,
		MaxY:  MaxY,
		CntsY: CntsY,
		grids: make(map[int] *Grid),
	}
	//AOI格子编号处理 9宫格区域的所有信息
	for y:=0; y < CntsY; y++ {
		for x :=0;  x < CntsX; x++ {
			//id = idy * cntX + idx
			gid := y * CntsX + x

			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX + x * aoiMgr.gridWidth(),
				aoiMgr.MinX + (x + 1) * aoiMgr.gridWidth(),
				aoiMgr.MinY + y + aoiMgr.gridLength(),
				aoiMgr.MinY + (y + 1)+ aoiMgr.gridLength())
		}
	}
	return aoiMgr
}


//格子X轴方向的 宽度           + 偏移 基本等于X轴上坐标的所有位置
func  (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}


//格子Y轴方向的 宽度      + 偏移 基本等于X轴上坐标的所有位置
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

func (m * AOIManager) String() string {
	s := fmt.Sprintf("AOIManager: \n Minx: %d MaxX: %d conts %d minY %d maxY %d cntsY %d \n Grids in AOIManager",
				m.MinX,m.MaxX,m.CntsX,m.MinY,m.MaxY,m.CntsY)

	for _,grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}






