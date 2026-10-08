# Step 05: DDD + クリーンアーキテクチャ

ゴール: このリポジトリの本番構成 (`cmd/` + `internal/`) を理解し、API を動かす。

## 全体像

```
cmd/server/main.go          # 起動・依存の組み立て
internal/
  domain/task/              # エンティティ・Repository interface（最重要）
  usecase/task/             # アプリケーションユースケース
  infrastructure/mysql/     # DB 実装
  interface/http/           # HTTP アダプタ（入出力変換）
  config/                   # 環境変数
```

依存の向きは常に **外側 → 内側** です。

```
HTTP -> UseCase -> Domain <- MySQL Repository
```

Domain は UseCase / HTTP / MySQL を import しません。

## 起動

```bash
make up
make health
make demo
```

## API

| Method | Path | 説明 |
|--------|------|------|
| GET | `/health` | ヘルスチェック |
| GET | `/tasks` | 一覧 |
| POST | `/tasks` | 作成 |
| GET | `/tasks/{id}` | 取得 |
| PUT | `/tasks/{id}` | 更新 |
| DELETE | `/tasks/{id}` | 削除 |

例:

```bash
curl -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"クリーンアーキテクチャ","description":"層の役割を説明する"}'

curl -X PUT http://localhost:8080/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"クリーンアーキテクチャ","description":"更新した","status":"doing"}'
```

## DDD っぽいところ

- `Task` エンティティが生成・更新ルールを持つ (`New`, `Update`)
- `Status` はただの文字列ではなく、許可された値だけ通す
- Repository は「永続化の契約」として domain 側に置く

## クリーンアーキテクチャっぽいところ

- UseCase が業務フローをオーケストレーション
- HTTP DTO と Domain を分離
- MySQL 実装は infrastructure に閉じ込める
- `main` だけが全体を知って配線する

## テスト

ドメインは DB なしでテストできます。

```bash
make test
```

## 演習（チーム配属前の仕上げ）

1. `Task` に `Priority`（low/medium/high）を追加する
2. `GET /tasks?status=todo` のように絞り込みを追加する
3. UseCase のモック Repository テストを書く
4. 別ドメイン（例: `User`）を同じ層構成で追加する

ここまでできれば、実務の Go バックエンドの読み書きにかなり近づいています。
