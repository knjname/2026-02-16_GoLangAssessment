# Go Todo API

Go 言語のスタック検証用プロジェクト。DDD レイヤードアーキテクチャで構成された Todo REST API とバッチ CLI ツール。

## 技術スタック

| カテゴリ | 技術 |
|----------|------|
| 言語 | Go 1.25.5 |
| API フレームワーク | [Huma v2](https://github.com/danielgtaylor/huma) (OpenAPI 自動生成) |
| DB | PostgreSQL 16 ([pgx/v5](https://github.com/jackc/pgx) ドライバ) |
| マイグレーション | [Goose v3](https://github.com/pressly/goose) |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| DI | [kessoku](https://github.com/mazrean/kessoku) (コンパイル時コード生成) |
| テスト | [testify](https://github.com/stretchr/testify), [humatest](https://github.com/danielgtaylor/huma), [testcontainers-go](https://github.com/testcontainers/testcontainers-go) |
| モック | [mockery](https://github.com/vektra/mockery) |
| リント | [golangci-lint](https://github.com/golangci/golangci-lint) |

## アーキテクチャ

4 層の DDD レイヤー構成。依存方向は `handler → usecase → domain ← repository`。

```
internal/
├── domain/          エンティティ・ビジネスルール (Todo, バリデーション, ドメインエラー)
├── usecase/         ビジネスロジックのオーケストレーション + TodoRepository インターフェース
├── repository/
│   └── postgres/    PostgreSQL 実装 (SQL クエリ定数, エラーマッピング)
├── handler/         HTTP ハンドラ (Huma v2 エンドポイント登録, ドメインエラー→HTTPステータス変換)
├── di/              kessoku による DI 定義 + 自動生成コード
├── server/          HTTP サーバー初期化・graceful shutdown
├── middleware/      ロギング, パニックリカバリ, リクエストID
└── config/          環境変数読み込み
```

## API エンドポイント

| Method | Path | 概要 |
|--------|------|------|
| `POST` | `/todos` | Todo 作成 |
| `GET` | `/todos` | Todo 一覧 |
| `GET` | `/todos/{id}` | Todo 取得 |
| `PUT` | `/todos/{id}` | Todo 更新 |
| `DELETE` | `/todos/{id}` | Todo 削除 |
| `POST` | `/todos/{id}/complete` | 完了マーク |
| `POST` | `/todos/complete-all` | 全件完了 |

## セットアップ

```bash
# ローカル PostgreSQL 起動
make docker-up

# マイグレーション実行
make migrate-up

# API サーバー起動
make run
```

## コマンド一覧

```bash
# ビルド (bin/api, bin/batch を生成)
make build

# テスト
make test              # 全テスト (race detector 有効, integration 含む)
make test-short        # ユニットテストのみ
make test-integration  # repository 層の integration テストのみ

# リント
make lint

# コード生成 (モック + DI)
make generate

# Docker
make docker-up         # PostgreSQL 起動
make docker-down       # PostgreSQL 停止
```

### バッチ CLI

```bash
go run ./cmd/batch migrate up       # マイグレーション適用
go run ./cmd/batch migrate down     # ロールバック
go run ./cmd/batch list             # Todo 一覧表示
go run ./cmd/batch complete-all     # 全件完了
```

## 環境変数

| 変数 | デフォルト値 | 説明 |
|------|------------|------|
| `PORT` | `8080` | API サーバーポート |
| `DATABASE_URL` | `postgres://todo:todo@localhost:5432/todo?sslmode=disable` | PostgreSQL 接続文字列 |
| `LOG_LEVEL` | `info` | ログレベル |

## テスト戦略

| レイヤー | 手法 |
|----------|------|
| Domain | ユニットテスト (`testify/assert`) |
| Usecase | mockery 生成モックで `TodoRepository` をモック化 |
| Handler | `humatest` で HTTP リクエスト/レスポンスをテスト |
| Repository | `testcontainers-go` で実 PostgreSQL コンテナを起動する integration テスト |

## DI (依存性注入)

[kessoku](https://github.com/mazrean/kessoku) によるコンパイル時 DI を採用。`internal/di/` にプロバイダ定義を集約し、`go generate` で injector 関数を自動生成する。

- **API 用**: `InitializeAPI(ctx)` — Config, Logger, Pool, Repository, UseCase を解決
- **Batch 用**: `InitializeBatch(ctx)` — 上記に加え `*sql.DB` (マイグレーション用) を解決。Pool と StdDB は `Async` で並列初期化

## CI

GitHub Actions (`.github/workflows/ci.yml`) が `main` への push / PR で実行:

1. **lint** — `golangci-lint`
2. **test** — PostgreSQL サービスコンテナ付きで全テスト実行
3. **build** — `bin/api`, `bin/batch` のビルド確認
4. **generate-check** — `go generate` の結果が最新であることを検証

## Docker

```bash
# イメージビルド
docker build -t go-todo-api .

# 実行
docker run -p 8080:8080 -e DATABASE_URL=postgres://... go-todo-api
```

マルチステージビルドで distroless イメージを使用。`/bin/api`, `/bin/batch`, `/migrations` を含む。
