# collision-matrix

碰撞过滤矩阵：双向层掩码、触发器剔除与去重排序的碰撞对，只用 Go 标准库。

```
go run ./cmd/collision-matrix --sample=oneway
go run ./cmd/collision-matrix --sample=trigger
go run ./cmd/collision-matrix --sample=dup
go run ./cmd/collision-matrix --sample=order
go run ./cmd/collision-matrix --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **双向放行**：只有 a 的层掩码包含 b 的层、并且 b 的层掩码包含 a 的层时，
  `ShouldHit` 才返回 true；单向不算。
- **触发器**：`trigger` 为真的实体不参与 `Pairs`，但登记信息要保留。
- **去重**：同一对实体在 `Pairs` 里只出现一次，`A` 字典序小于 `B`。
- **排序**：结果按 `A` 升序、`A` 相同按 `B` 升序。
- **未知层**：没有登记过的层不参与任何碰撞。
- **不变量**：`Pairs` 里不含同一个实体和自己成对；同一串配置重复执行结果相同。
- **代价**：`ShouldHit` 不许遍历全部层定义，`Scanned` 不随层数线性增长。
- **规模**：10 万个实体、每秒 10 万次判定，单次摊还 O(1)。

## 输出契约（不改格式）

```
hit=..
pairs=..
pairs=.. first=..
scanned=..
```

## 目录

```
collision.go              掩码、触发器与碰撞对
collision_test.go         单元测试
cmd/collision-matrix/     命令行入口
```
