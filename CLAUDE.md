# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go Todo API — DDD (Domain-Driven Design) アーキテクチャによるTodo REST APIとバッチCLIツール。

- モジュール: `github.com/knjname/go-todo-api`
- Go: 1.25.5
- APIフレームワーク: Huma v2 (OpenAPI自動生成)
- DB: PostgreSQL 16 (pgx/v5ドライバ)
- マイグレーション: Goose v3
- DI: kessoku (コード生成ベースのDIコンテナ)

## Common Commands

```bash
# ローカルPostgreSQL起動/停止
make docker-up
make docker-down

# マイグレーション
make migrate-up
make migrate-down

# ビルド (bin/api, bin/batch を生成)
make build

# APIサーバー起動
make run

# テスト
make test              # 全テスト (race detector有効、integration含む)
make test-short        # ユニットテストのみ (integrationスキップ)
make test-integration  # repository層のintegrationテストのみ

# 単一テスト実行例
go test -race -count=1 -run TestCreateTodo ./internal/usecase/...

# リント
make lint              # golangci-lint実行

# コード生成 (モック + DI)
make generate          # go generate (mockery + kessoku)
```

## Architecture

4層のDDDレイヤー構成。依存方向は handler → usecase → domain ← repository。

- **`internal/domain/`** — エンティティとビジネスルール。`Todo`構造体、バリデーション(`MaxTitleLength=200`)、ドメインエラー(`ErrNotFound`, `ErrValidation`)
- **`internal/usecase/`** — ビジネスロジックのオーケストレーション。`TodoRepository`インターフェース定義もここ(`interfaces.go`)。mockeryで自動生成されるモックは`mocks/`配下
- **`internal/repository/postgres/`** — PostgreSQL実装。SQLクエリは`queries.go`に定数定義。DBエラーをドメインエラーへマッピング
- **`internal/handler/`** — HTTPハンドラ。Huma v2でOpenAPIスキーマ付きエンドポイント登録。ドメインエラーからHTTPステータスへの変換は`errors.go`

支援層:
- **`internal/di/`** — DI構成。kessokuによるコード生成で依存解決。`api.go`/`batch.go`が定義、`*_band.go`が生成コード。プロバイダ関数は`providers.go`
- **`internal/middleware/`** — ロギング、パニックリカバリ、リクエストID
- **`internal/config/`** — 環境変数読み込み (`PORT`, `DATABASE_URL`, `LOG_LEVEL`)
- **`internal/server/`** — HTTPサーバー初期化とgraceful shutdown

エントリポイント:
- **`cmd/api/`** — REST APIサーバー
- **`cmd/batch/`** — CLIツール (migrate up/down, list, complete-all)

## Testing Patterns

- **Domain/Usecase/Handler**: ユニットテスト。`testify/assert`+`testify/require`使用
- **Usecase**: `mockery`生成モックで`TodoRepository`をモック化
- **Handler**: `humatest`でHTTPリクエスト/レスポンスをテスト
- **Repository**: `testcontainers-go`で実PostgreSQLコンテナを起動するintegrationテスト。`-short`フラグでスキップされる
- テスト内ロガーは`slog.DiscardHandler`で出力抑制

## API Endpoints

| Method | Path | 概要 |
|--------|------|------|
| POST | `/todos` | Todo作成 |
| GET | `/todos` | Todo一覧 |
| GET | `/todos/{id}` | Todo取得 |
| PUT | `/todos/{id}` | Todo更新 |
| DELETE | `/todos/{id}` | Todo削除 |
| POST | `/todos/{id}/complete` | 完了マーク |
| POST | `/todos/complete-all` | 全件完了 |

## Environment Variables

| 変数 | デフォルト値 |
|------|------------|
| `PORT` | `8080` |
| `DATABASE_URL` | `postgres://todo:todo@localhost:5432/todo?sslmode=disable` |
| `LOG_LEVEL` | `info` |

## CI

GitHub Actions (`.github/workflows/ci.yml`) が `main`へのpush/PRで実行:
- lint (`golangci-lint`)
- test (PostgreSQLサービスコンテナ付き)
- build
- generate check (`go generate`の結果が最新であることを検証)
