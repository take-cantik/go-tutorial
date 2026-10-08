# Go Backend Tutorial

Go バックエンド配属向けのハンズオンです。  
**タスク管理 API** を題材に、基礎 → HTTP → MySQL → レイヤー分離 → DDD/クリーンアーキテクチャ まで段階的に進めます。

## ゴール

- Docker で Go + MySQL の開発環境を扱える
- `net/http` と `database/sql` で API を書ける
- 依存性逆転を理解し、クリーンアーキテクチャ構成のコードを読める / 直せる

## 必要環境

- Docker（Compose V2: `docker compose`）
- （任意）ローカルの Go 1.22+ ※ステップ実行やテスト用

## クイックスタート（最終形）

```bash
# コンテナ起動（app + mysql）
make up

# ヘルスチェック
make health

# CRUD デモ
make demo
```

アプリ: http://localhost:8080  
MySQL: `localhost:3306` / user `tutorial` / password `tutorial` / db `go_tutorial`

停止:

```bash
make down
```

## 学習ルート

| Step | 内容 | コマンド |
|------|------|----------|
| [01](tutorials/01-basics.md) | Go の基本 | `make step1` |
| [02](tutorials/02-http.md) | HTTP JSON API | `make step2` |
| [03](tutorials/03-mysql.md) | MySQL CRUD | `make up` → `make step3` |
| [04](tutorials/04-layered.md) | Handler/Service/Repository | `make step4` |
| [05](tutorials/05-clean-architecture.md) | DDD + クリーンアーキテクチャ | `make up` / `make demo` |

おすすめの進め方:

1. `tutorials/` を上から読む
2. 各 `steps/0N-*` を動かして改造する
3. 最後に `internal/` の最終形を読む

## 最終形のディレクトリ構成

```
.
├── compose.yaml                 # Go app + MySQL
├── Dockerfile
├── cmd/server/main.go           # エントリポイント（依存の組み立て）
├── internal/
│   ├── domain/task/             # エンティティ・Repository interface
│   ├── usecase/task/            # ユースケース
│   ├── infrastructure/mysql/    # MySQL 実装
│   ├── interface/http/          # HTTP ハンドラ
│   └── config/
├── migrations/                  # DB 初期化 SQL
├── steps/                       # 学習用の小さなサンプル
└── tutorials/                   # 解説ドキュメント
```

### 依存関係（依存は内向き）

```
interface/http ──▶ usecase ──▶ domain
                      ▲
infrastructure/mysql ─┘
```

## API（最終形）

| Method | Path | 説明 |
|--------|------|------|
| GET | `/health` | 死活監視 |
| GET | `/tasks` | 一覧 |
| POST | `/tasks` | 作成 `{"title","description"}` |
| GET | `/tasks/{id}` | 取得 |
| PUT | `/tasks/{id}` | 更新 `{"title","description","status"}` |
| DELETE | `/tasks/{id}` | 削除 |

`status` は `todo` / `doing` / `done` のみ。

## よく使うコマンド

```bash
make up        # 起動
make down      # 停止
make logs      # app ログ
make test      # ユニットテスト
make tidy      # go mod tidy
```

## 補足

- Compose 定義は `compose.yaml` です（Docker Compose の正式なファイル名）
- app コンテナは [Air](https://github.com/air-verse/air) でホットリロードします
- ドメインのルールは DB なしでテストできます: `make test`
