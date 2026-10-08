# Step 03: MySQL 連携

ゴール: Docker 上の MySQL に接続して CRUD の入口を作る。

## 前提

```bash
make up   # app + mysql を起動
```

MySQL は `compose.yaml` の `db` サービスです。
初期スキーマは `migrations/001_create_tasks.sql` が自動適用されます。

## やってみよう

```bash
make step3
```

別ターミナル:

```bash
curl http://localhost:8082/tasks
curl -X POST http://localhost:8082/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"MySQLに保存"}'
```

## 読むポイント

- DSN (`user:pass@tcp(host:port)/dbname`)
- `database/sql` + `go-sql-driver/mysql`
- `Query` / `Exec` / `Scan` / `LastInsertId`
- ハンドラの中に SQL が直接書かれていること（次ステップで分離する）

## 演習

1. `GET /tasks/{id}` を追加する
2. `UPDATE` で status を更新する API を追加する
3. SQL エラー時にログを出し分ける

## 次へ

[04-layered.md](./04-layered.md)
