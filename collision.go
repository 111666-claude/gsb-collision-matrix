// Package collision 是碰撞过滤矩阵：双向层掩码、触发器剔除与去重排序的碰撞对。
package collision

import "sort"

// Pair 是一对需要检测碰撞的实体。
type Pair struct {
	A string
	B string
}

type body struct {
	id      string
	layer   string
	trigger bool
}

// Matrix 是碰撞过滤矩阵。
type Matrix struct {
	layers  map[string]uint32
	bits    map[string]uint
	bodies  map[string]body
	order   []string
	scanned int
}

// New 建矩阵。
func New() *Matrix {
	return &Matrix{
		layers: map[string]uint32{},
		bits:   map[string]uint{},
		bodies: map[string]body{},
	}
}

// SetLayer 设置某一层允许交互的层掩码。首次登记时按登记顺序分配层位号。
func (m *Matrix) SetLayer(name string, mask uint32) {
	if _, ok := m.layers[name]; !ok {
		m.bits[name] = uint(len(m.bits))
	}
	m.layers[name] = mask
}

// AddBody 加一个实体。
func (m *Matrix) AddBody(id, layer string, trigger bool) {
	m.bodies[id] = body{id: id, layer: layer, trigger: trigger}
	m.order = append(m.order, id)
}

// ShouldHit 判断两个实体是否应该检测碰撞。
// 只有 a 的层掩码包含 b 的层、并且 b 的层掩码包含 a 的层时才放行。
func (m *Matrix) ShouldHit(a, b string) bool {
	m.scanned += 2
	left, okA := m.bodies[a]
	right, okB := m.bodies[b]
	if !okA || !okB {
		return false
	}
	return m.layersHit(left.layer, right.layer)
}

// layersHit 判断两个登记过的层是否双向放行，O(1) 查表不遍历层定义。
func (m *Matrix) layersHit(layerA, layerB string) bool {
	maskA, okA := m.layers[layerA]
	maskB, okB := m.layers[layerB]
	if !okA || !okB {
		return false
	}
	bitA := m.bits[layerA]
	bitB := m.bits[layerB]
	return maskA&(uint32(1)<<bitB) != 0 && maskB&(uint32(1)<<bitA) != 0
}

// Pairs 返回所有需要检测碰撞的实体对。
// 触发器不参与；同一对只出现一次且 A 字典序小于 B；结果按 A、B 升序。
func (m *Matrix) Pairs() []Pair {
	ids := make([]string, 0, len(m.bodies))
	for id := range m.bodies {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := []Pair{}
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			left := m.bodies[ids[i]]
			right := m.bodies[ids[j]]
			if left.trigger || right.trigger {
				continue
			}
			if !m.layersHit(left.layer, right.layer) {
				continue
			}
			out = append(out, Pair{A: ids[i], B: ids[j]})
		}
	}
	return out
}

// Scanned 是累计扫描量（规模观测）。
func (m *Matrix) Scanned() int { return m.scanned }
