package core

import "fmt"

//AOI 区域管理模块

type AOIManager struct {
	//区域左边界坐标
	MinX int

	//区域右边界坐标
	MaxX int

	//X方向格子的数量
	CntsX int

	// 区域的上边界坐标
	MinY int

	// 区域的下边界坐标
	MaxY int

	//Y方向格子的数量
	CntsY int

	// 当前区域中有哪些格子map: key 格子的ID value 格子的对象
	grids map[int] *Grid

}
// 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int)*AOIManager {
	aoiMgr :=  &AOIManager{
		MinX: minX,
		MaxX: maxX,
		CntsX:cntsX,
		MinY:minY,
		MaxY:maxY,
		CntsY:cntsY,
		grids:make(map[int]*Grid),
	}

	// 给AOI初始化区域的格子所有的格子进行编号和初始化
	for y:=0; y < cntsY; y++ {
		for x:=0; x<cntsX; x++ {
			// 计算格子ID，根据x, y编号
			// 格子编号：id = idy *cntsX + idX
			gid := y * cntsX + x
			// 初始化gid格子
			aoiMgr.grids[gid] = &Grid{
				GID:gid,
				MinX:x * aoiMgr.gridWitdh(),
				MaxX: (x+1) * aoiMgr.gridWitdh(),
				MinY:y * aoiMgr.gridHeigh(),
				MaxY:(y+1) * aoiMgr.gridHeigh(),
			}
		}
	}

	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func(m* AOIManager) gridWitdh() int {
	return (m.MaxX - m.MinX)/m.CntsX
}

// 得到每个格子在Y轴方向的长度
func(m* AOIManager) gridHeigh() int {
	return (m.MaxY - m.MinY)/m.CntsY
}

// 打印格子信息
func(m *AOIManager) String()string {
	//打印AOIManager信息
	s:= fmt.Sprintf("AOIManager\n MinX: %d, MaxX: %d, cntsX:%d, MinY: %d, MaxY:%d, cntsY:%d\n Grid in AOIManager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	// 打印全部格子信息
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}