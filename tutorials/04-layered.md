# Step 04: レイヤー分離

ゴール: Handler / Service / Repository に分け、依存を逆転させる入口を作る。

## なぜ分ける？

Step 03 は「動く」が、HTTP・業務ルール・SQL が同じ場所にあります。
チーム開発では次の問題が出ます。

- テストしづらい（HTTP を立てないと業務ロジックを検証できない）
- DB を差し替えにくい
- 変更の影響範囲が見えない

## やってみよう

```bash
make step4
curl http://localhost:8083/tasks
```

## 構造

```
Handler  ->  Service  ->  Repository(interface)  ->  MySQL実装
```

ポイントは **Service が MySQL 構造体ではなく interface に依存している** ことです。
これがクリーンアーキテクチャの依存性逆転（DIP）の原型です。

## 読むポイント

- `TaskRepository` interface
- `TaskService` がバリデーションを持つ
- `main` で実装を組み立てる（Composition Root）

## 演習

1. Repository のモックを作って Service だけ単体テストする
2. `Update` / `Delete` を同じ形で追加する
3. Service を `usecase` という名前にリネームしてみる（最終形に近づく）

## 次へ

[05-clean-architecture.md](./05-clean-architecture.md)
