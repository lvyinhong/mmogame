package core

import (
	"fmt"
	"sync"
)

// 一个AOI地图中的格子类型
type Grid struct {
	//格子ID
	GID int

	//格子的左边边界坐标
	MinX int
	//格子的右边边界坐标

	MaxX int
	//格子的上面边界坐标

	MinY int
	//格子的下边边界坐标

	MaxY int

	//当前格子内玩家或者物体成员的ID集合
	playerIDs map[int]bool

	//保护当前集合的锁
	pIDLock sync.RWMutex
}

func NewGid(gID, minX, maxX, minY, maxY int)*Grid {
	return &Grid{
		GID:gID,
		MinX: minX,
		MaxX:maxX,
		MinY:minY,
		MaxY:maxY,
		playerIDs:make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid)Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid)Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}

//得到当前格子中所有的玩家
func(g *Grid)GetPlayerIDs()(playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs =append(playerIDs, k)
	}
	return
}

// 调试使用-打印出格子的基本信息
func(g *Grid)String() string {
	return fmt.Sprintf("Gid:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v \n",
		g.GID,g.MinX,g.MaxX,g.MinY,g.MaxY,g.playerIDs);
}

// 根据格子GID得到周边九宫格的格子ID集合
func(m *AOIManager)GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _,ok := m.grids[gID]; !ok {
		return
	}
	// 初始化grids 返回值切片
	grids =append(grids, m.grids[gID])

	// 需要gID的左边是否有格子？右边是否有格子
	// 需要通过gID得到当前格子x轴的编号 idx = id%nx
	idx := gID % m.CntsX


	// 判断idx编号是否左边还有格子， 如果有 放在gidsX集合中
	if idx > 0 {
		grids = append(grids, m.grids[gID - 1])
	}

	// 判断idx编号是否右边还有格子， 如果有放在gidsX集合中
	if idx < m.CntsX -1 {
		grids = append(grids, m.grids[gID + 1])
	}

	// 将x 轴当前的格子都取出，进行遍历再分辨得到每个格子上下是否还有格子
	gidsX := make([]int, 0, len(grids))
	for _,v := range  grids {
		gidsX = append(gidsX, v.GID)
	}



	// 遍历gidsX集合中每一个格子的gid
	for _, v := range gidsX {
		// gid 上边是否还有格子
		idy := v / m.CntsY
		if idy > 0 {
			grids = append(grids, m.grids[v - m.CntsX])
		}
		if idy < m.CntsY - 1 {
			grids = append(grids, m.grids[v + m.CntsX])
		}
	}

	return
}

// 通过横纵坐标得到当前GIDs 的编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWitdh()
	idy := (int(y) - m.MinY) / m.gridHeigh()
	return idy*m.CntsX + idx
}

// 通过横纵坐标得到周边九宫格内全部的playerIDs
func (m *AOIManager) GetPidsByPos(x, y float32)(playerIDs []int) {
	// 得到当前玩家的GID格子id
		gID := m.GetGidByPos(x, y)
	// 通过GID得到周边九宫格信息
		grids := m.GetSurroundGridsByGid(gID)
	// 将九宫格的信息里的全部的Player的id 累加到playerIDs
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return
}

// 添加一个PlayerID 到一个格子中
func (m *AOIManager)AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的playerID
func (m *AOIManager)RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 通过GID获取全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int)(playerIDs []int)  {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 通过坐标将Player添加到一个格子中
func (m *AOIManager)AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// 通过坐标把一个Player从一个格子中删除
func (m *AOIManager) RemoveFromGridbyPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)
}