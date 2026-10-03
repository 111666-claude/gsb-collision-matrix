// Package collision 是碰撞过滤矩阵：双向层掩码、触发器剔除与去重排序的碰撞对。
// 缺陷：只看单向掩码、触发器也算进碰撞对、同一对正反各出一条、结果不排序、
// 每次判定都要遍历全部层定义。
package collision

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
	bodies  map[string]body
	order   []string
	scanned int
}

// New 建矩阵。
func New() *Matrix {
	return &Matrix{layers: map[string]uint32{}, bodies: map[string]body{}}
}

// SetLayer 设置某一层允许交互的层掩码。
func (m *Matrix) SetLayer(name string, mask uint32) { m.layers[name] = mask }

// AddBody 加一个实体。
func (m *Matrix) AddBody(id, layer string, trigger bool) {
	m.bodies[id] = body{id: id, layer: layer, trigger: trigger}
	m.order = append(m.order, id)
}

// ShouldHit 判断两个实体是否应该检测碰撞。
// 缺陷：只查单向掩码，且只判断掩码非零。
func (m *Matrix) ShouldHit(a, b string) bool {
	m.scanned += len(m.layers)
	left := m.bodies[a]
	right := m.bodies[b]
	mask, ok := m.layers[left.layer]
	if !ok {
		return false
	}
	if _, ok := m.layers[right.layer]; !ok {
		return false
	}
	return mask != 0
}

// Pairs 返回所有需要检测碰撞的实体对。
// 缺陷：正反各出一条、包含触发器、不排序。
func (m *Matrix) Pairs() []Pair {
	m.scanned += len(m.order) * len(m.order)
	out := []Pair{}
	for _, a := range m.order {
		for _, b := range m.order {
			if a == b {
				continue
			}
			out = append(out, Pair{A: a, B: b})
		}
	}
	return out
}

// Scanned 是累计扫描量（规模观测）。
func (m *Matrix) Scanned() int { return m.scanned }
