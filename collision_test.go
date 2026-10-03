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
