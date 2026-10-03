package collision

import "testing"

func TestEmptyMatrixHasNoPairs(t *testing.T) {
	if len(New().Pairs()) != 0 {
		t.Fatal("空矩阵不应该有碰撞对")
	}
}

func TestUnknownBodyDoesNotHit(t *testing.T) {
	if New().ShouldHit("ghost", "ghost") {
		t.Fatal("不存在的实体不应该产生碰撞")
	}
}

func TestLayerMaskIsStored(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 2)
	if matrix.layers["player"] != 2 {
		t.Fatal("层掩码应该能读回来")
	}
}

func TestOrderRecordsBodies(t *testing.T) {
	matrix := New()
	matrix.AddBody("a", "player", false)
	if len(matrix.order) != 1 {
		t.Fatal("AddBody 应该记录顺序")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestShouldHitRequiresBothDirections(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 2) // player 的掩码包含 mob
	matrix.SetLayer("mob", 0)    // mob 的掩码不含 player
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m", "mob", false)
	if matrix.ShouldHit("p", "m") || matrix.ShouldHit("m", "p") {
		t.Fatal("单向掩码不应该放行")
	}
	matrix.SetLayer("mob", 1) // 双向都包含对方
	if !matrix.ShouldHit("p", "m") || !matrix.ShouldHit("m", "p") {
		t.Fatal("双向掩码都包含对方时应该放行")
	}
}

func TestTriggerExcludedFromPairsButRegistered(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 2)
	matrix.SetLayer("mob", 1)
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m", "mob", true)
	if got := len(matrix.Pairs()); got != 0 {
		t.Fatalf("触发器不应参与碰撞对，得到 %d 条", got)
	}
	registered, ok := matrix.bodies["m"]
	if !ok || !registered.trigger || registered.layer != "mob" {
		t.Fatal("触发器的登记信息应该保留")
	}
}

func TestPairsAreDedupedWithCanonicalOrder(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 2)
	matrix.SetLayer("mob", 1)
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m", "mob", false)
	pairs := matrix.Pairs()
	if len(pairs) != 1 {
		t.Fatalf("同一对只应出现一次，得到 %d 条", len(pairs))
	}
	if pairs[0].A != "m" || pairs[0].B != "p" {
		t.Fatalf("A 应字典序小于 B，得到 %+v", pairs[0])
	}
}

func TestPairsAreSorted(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 3)
	matrix.SetLayer("mob", 3)
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m2", "mob", false)
	matrix.AddBody("m1", "mob", false)
	pairs := matrix.Pairs()
	want := []Pair{{A: "m1", B: "m2"}, {A: "m1", B: "p"}, {A: "m2", B: "p"}}
	if len(pairs) != len(want) {
		t.Fatalf("应有 %d 条碰撞对，得到 %d 条", len(want), len(pairs))
	}
	for i := range want {
		if pairs[i] != want[i] {
			t.Fatalf("第 %d 条应为 %+v，得到 %+v", i, want[i], pairs[i])
		}
	}
}

func TestUnknownLayerNeverCollides(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 1)
	matrix.AddBody("p", "player", false)
	matrix.AddBody("g", "ghost", false)
	if matrix.ShouldHit("p", "g") || matrix.ShouldHit("g", "p") {
		t.Fatal("未知层不应参与任何碰撞")
	}
	if got := len(matrix.Pairs()); got != 0 {
		t.Fatalf("未知层实体不应出现在碰撞对里，得到 %d 条", got)
	}
}

func TestSameConfigGivesSamePairs(t *testing.T) {
	build := func() []Pair {
		matrix := New()
		matrix.SetLayer("player", 3)
		matrix.SetLayer("mob", 3)
		matrix.AddBody("p", "player", false)
		matrix.AddBody("m2", "mob", false)
		matrix.AddBody("m1", "mob", false)
		return matrix.Pairs()
	}
	first, second := build(), build()
	if len(first) != len(second) {
		t.Fatal("同一串配置重复执行结果应相同")
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatal("同一串配置重复执行结果应相同")
		}
	}
}

func TestScannedDoesNotGrowWithLayerCount(t *testing.T) {
	matrix := New()
	for i := 0; i < 3000; i++ {
		matrix.SetLayer(string(rune('a'+i%26))+string(rune('0'+i/26%10))+string(rune('0'+i/260)), 1)
	}
	matrix.AddBody("a", "a00", false)
	matrix.AddBody("b", "b00", false)
	for i := 0; i < 3000; i++ {
		matrix.ShouldHit("a", "b")
	}
	if matrix.Scanned() > 6000 {
		t.Fatalf("Scanned 不应随层数线性增长，得到 %d", matrix.Scanned())
	}
}
