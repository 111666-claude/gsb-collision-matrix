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

func TestShouldHitRequiresBothWays(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 2) // player 的掩码只含 mob
	matrix.SetLayer("mob", 0)    // mob 的掩码不含 player
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m", "mob", false)
	if matrix.ShouldHit("p", "m") || matrix.ShouldHit("m", "p") {
		t.Fatal("单向掩码不应该放行")
	}
	matrix.SetLayer("mob", 1) // mob 也放行 player 后才算双向
	if !matrix.ShouldHit("p", "m") || !matrix.ShouldHit("m", "p") {
		t.Fatal("双向掩码应该放行")
	}
}

func TestShouldHitRejectsUnknownLayer(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", ^uint32(0))
	matrix.AddBody("p", "player", false)
	matrix.AddBody("g", "ghost", false)
	if matrix.ShouldHit("p", "g") {
		t.Fatal("未知层不应该参与碰撞")
	}
}

func TestTriggerExcludedFromPairsButRegistered(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 2)
	matrix.SetLayer("mob", 1)
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m", "mob", true)
	if len(matrix.Pairs()) != 0 {
		t.Fatal("触发器不应该出现在 Pairs 里")
	}
	if _, ok := matrix.bodies["m"]; !ok {
		t.Fatal("触发器的登记信息应该保留")
	}
}

func TestPairsAreDedupedAndOrdered(t *testing.T) {
	matrix := New()
	matrix.SetLayer("player", 3)
	matrix.SetLayer("mob", 3)
	matrix.AddBody("p", "player", false)
	matrix.AddBody("m2", "mob", false)
	matrix.AddBody("m1", "mob", false)
	pairs := matrix.Pairs()
	want := []Pair{
		{A: "m1", B: "m2"},
		{A: "m1", B: "p"},
		{A: "m2", B: "p"},
	}
	if len(pairs) != len(want) {
		t.Fatalf("期望 %d 对，实际 %d 对", len(want), len(pairs))
	}
	for i, pair := range pairs {
		if pair != want[i] {
			t.Fatalf("第 %d 对期望 %v，实际 %v", i, want[i], pair)
		}
		if pair.A >= pair.B {
			t.Fatalf("A 应该字典序小于 B：%v", pair)
		}
	}
}

func TestPairsAreDeterministic(t *testing.T) {
	build := func() []Pair {
		matrix := New()
		matrix.SetLayer("player", 3)
		matrix.SetLayer("mob", 3)
		matrix.AddBody("p", "player", false)
		matrix.AddBody("m2", "mob", false)
		matrix.AddBody("m1", "mob", false)
		return matrix.Pairs()
	}
	first := build()
	for i := 0; i < 10; i++ {
		got := build()
		if len(got) != len(first) {
			t.Fatal("重复执行结果应该一致")
		}
		for j := range got {
			if got[j] != first[j] {
				t.Fatal("重复执行结果应该一致")
			}
		}
	}
}

func TestScannedDoesNotGrowWithLayers(t *testing.T) {
	matrix := New()
	for i := 0; i < 3000; i++ {
		matrix.SetLayer(string(rune('a'+i%26))+string(rune('a'+i/26)), 1)
	}
	matrix.AddBody("a", "aa", false)
	matrix.AddBody("b", "ab", false)
	for i := 0; i < 3000; i++ {
		matrix.ShouldHit("a", "b")
	}
	if matrix.Scanned() > 6000 {
		t.Fatalf("Scanned 不应该随层数线性增长，实际 %d", matrix.Scanned())
	}
}
