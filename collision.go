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
	bits    map[string]uint32
	bodies  map[string]body
	order   []string
	scanned int
}

// New 建矩阵。
func New() *Matrix {
	return &Matrix{
		layers: map[string]uint32{},
		bits:   map[string]uint32{},
		bodies: map[string]body{},
	}
}

// SetLayer 设置某一层允许交互的层掩码。
// 层位按登记顺序确定性分配（首层 bit0，逐层递增），同一层重复设置沿用原位。
func (m *Matrix) SetLayer(name string, mask uint32) {
	if _, ok := m.layers[name]; !ok && len(m.bits) < 32 {
		m.bits[name] = 1 << uint(len(m.bits))
	}
	m.layers[name] = mask
}

// AddBody 加一个实体。重复登记同一 id 时更新信息，不重复占位。
func (m *Matrix) AddBody(id, layer string, trigger bool) {
	if _, ok := m.bodies[id]; !ok {
		m.order = append(m.order, id)
	}
	m.bodies[id] = body{id: id, layer: layer, trigger: trigger}
}

// ShouldHit 判断两个实体是否应该检测碰撞：
// 两个实体都存在、所在层都已登记，且 a 的掩码包含 b 的层、b 的掩码也包含 a 的层。
// 触发器与否不影响判定；未知实体或未知层一律 false。
func (m *Matrix) ShouldHit(a, b string) bool {
	m.scanned++
	left, ok := m.bodies[a]
	if !ok {
		return false
	}
	right, ok := m.bodies[b]
	if !ok {
		return false
	}
	return m.layersHit(left.layer, right.layer)
}

// layersHit 做双向层掩码判定，只做常数次 map 查询，不遍历层定义。
func (m *Matrix) layersHit(leftLayer, rightLayer string) bool {
	maskLeft, ok := m.layers[leftLayer]
	if !ok {
		return false
	}
	maskRight, ok := m.layers[rightLayer]
	if !ok {
		return false
	}
	bitLeft, ok := m.bits[leftLayer]
	if !ok {
		return false
	}
	bitRight, ok := m.bits[rightLayer]
	if !ok {
		return false
	}
	return maskLeft&bitRight != 0 && maskRight&bitLeft != 0
}

// Pairs 返回所有需要检测碰撞的实体对。
// trigger 实体与未知层实体不参与；同一对只出现一次且 A 字典序小于 B，
// 结果按 A 升序、A 相同按 B 升序。
func (m *Matrix) Pairs() []Pair {
	solids := make([]body, 0, len(m.bodies))
	for _, id := range m.order {
		candidate, ok := m.bodies[id]
		if !ok || candidate.trigger {
			continue
		}
		if _, ok := m.layers[candidate.layer]; !ok {
			continue
		}
		solids = append(solids, candidate)
	}

	out := []Pair{}
	for i := 0; i < len(solids); i++ {
		for j := i + 1; j < len(solids); j++ {
			m.scanned++
			if !m.layersHit(solids[i].layer, solids[j].layer) {
				continue
			}
			a, b := solids[i].id, solids[j].id
			if b < a {
				a, b = b, a
			}
			out = append(out, Pair{A: a, B: b})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].A != out[j].A {
			return out[i].A < out[j].A
		}
		return out[i].B < out[j].B
	})
	return out
}

// Scanned 是累计扫描量（规模观测）。
func (m *Matrix) Scanned() int { return m.scanned }
