# Step 02: HTTP サーバー

ゴール: 標準ライブラリだけで JSON API を動かす。

## なぜ重要？

Go のバックエンドは、まず `net/http` を理解するとフレームワークの意図が読みやすくなります。
このリポジトリの最終形も、外部ルーターに依存せず標準ライブラリを使っています。

## やってみよう

```bash
make step2
```

別ターミナルで:

```bash
curl http://localhost:8081/tasks
curl -X POST http://localhost:8081/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"HTTPを学ぶ"}'
```

## 読むポイント

- `http.NewServeMux` と Go 1.22+ のメソッドルーティング (`GET /tasks`)
- `encoding/json` での encode / decode
- メモリストア + `sync.Mutex`（並行リクエスト対策の入口）

## 演習

1. `GET /tasks/{id}` を追加する
2. title が空のとき 400 を返す（すでに近い実装あり）
3. 削除 API `DELETE /tasks/{id}` を追加する

## 次へ

[03-mysql.md](./03-mysql.md)
