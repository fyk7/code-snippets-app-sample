# Code Snippets App

コードスニペットを管理・検索するための REST API サーバ。Go + Echo + GORM で構築し、ヘキサゴナルアーキテクチャ（クリーンアーキテクチャ）を採用しています。

---

## Architecture

### Hexagonal Architecture (Ports & Adapters)

本プロジェクトは **ヘキサゴナルアーキテクチャ** に基づき、ビジネスロジックを外部の技術的関心事から完全に分離しています。

```mermaid
graph TB
    subgraph External["External (外部)"]
        CLIENT["HTTP Client<br/>(Browser / curl)"]
        MYSQL[("MySQL<br/>Database")]
        S3["AWS S3"]
    end

    subgraph Adapters["Interface Adapter 層"]
        direction TB
        HANDLER["Handler<br/>(HTTP Adapter)"]
        MIDDLEWARE["Middleware<br/>(CORS)"]
        REPO_IMPL["Repository Impl<br/>(DB Adapter)"]
        S3_CLIENT["S3 Client<br/>(Storage Adapter)"]
    end

    subgraph Core["Application Core"]
        direction TB
        subgraph Usecase["Usecase 層"]
            SNIPPET_SVC["SnippetService"]
            TAG_SVC["TagService"]
            USER_SVC["UserService"]
        end
        subgraph Domain["Domain 層 (中心)"]
            MODEL["Model<br/>Snippet / Tag / User"]
            REPO_IF["Repository Interface<br/>(Port)"]
            ERRORS["Domain Errors"]
        end
    end

    CLIENT -->|"HTTP Request"| HANDLER
    HANDLER -->|"calls"| SNIPPET_SVC
    HANDLER -->|"calls"| TAG_SVC
    SNIPPET_SVC -->|"uses"| REPO_IF
    TAG_SVC -->|"uses"| REPO_IF
    USER_SVC -->|"uses"| REPO_IF
    REPO_IF -.->|"implemented by"| REPO_IMPL
    REPO_IMPL -->|"SQL"| MYSQL
    S3_CLIENT -->|"API"| S3

    style Domain fill:#4a90d9,stroke:#2c5f8a,color:#fff
    style Usecase fill:#7ab648,stroke:#5a8a38,color:#fff
    style Adapters fill:#f5a623,stroke:#c4841c,color:#fff
    style External fill:#9b9b9b,stroke:#6b6b6b,color:#fff
```

### Clean Architecture - レイヤー間の依存方向

依存の方向は **常に外側から内側** へ向かいます。内側の層は外側の層を一切知りません。

```mermaid
graph LR
    subgraph L4["Infrastructure"]
        DB["database.NewDB"]
        AWS["aws.S3Client"]
    end

    subgraph L3["Interface Adapter"]
        H["Handler"]
        RI["Repository Impl"]
        MW["Middleware"]
    end

    subgraph L2["Usecase"]
        SS["SnippetService"]
        TS["TagService"]
        US["UserService"]
    end

    subgraph L1["Domain (最内層)"]
        M["Model"]
        RP["Repository Interface"]
        E["Errors"]
    end

    L4 --> L3
    L3 --> L2
    L2 --> L1

    RI -.->|"implements"| RP

    style L1 fill:#4a90d9,stroke:#2c5f8a,color:#fff
    style L2 fill:#7ab648,stroke:#5a8a38,color:#fff
    style L3 fill:#f5a623,stroke:#c4841c,color:#fff
    style L4 fill:#9b9b9b,stroke:#6b6b6b,color:#fff
```

### リクエストフロー (シーケンス図)

```mermaid
sequenceDiagram
    participant C as Client
    participant MW as Middleware
    participant H as Handler
    participant V as Validator
    participant S as Service
    participant R as Repository
    participant DB as MySQL

    C->>MW: HTTP Request
    MW->>H: (CORS headers added)
    H->>H: Bind request body
    H->>V: ValidRequest(req)
    V-->>H: ok / validation error

    H->>S: Service method(ctx, ...)
    S->>S: context.WithTimeout(ctx)
    S->>R: Repository method(ctx, ...)
    R->>R: Conn.WithContext(ctx)
    R->>DB: SQL Query
    DB-->>R: Result / Error
    R->>R: toDomainError(err)
    R-->>S: domain model / domain error
    S-->>H: result / error
    H->>H: handleError(c, err)
    H-->>C: JSON Response
```

### Dependency Injection (Google Wire)

```mermaid
graph TD
    CONFIG["config.Config"] --> NEWDB["database.NewDB"]
    NEWDB --> DB["*gorm.DB"]

    DB --> NSR["NewSnippetRepository"]
    DB --> NTR["NewTagRepository"]
    DB --> NUR["NewUserRepository"]

    NSR --> SR["SnippetRepository"]
    NTR --> TR["TagRepository"]
    NUR --> UR["UserRepository"]

    SR --> NSS["NewSnippetService"]
    TR --> NTS["NewTagService"]
    UR --> NUS["NewUserService"]

    TIMEOUT["time.Duration"] --> NSS
    TIMEOUT --> NTS
    TIMEOUT --> NUS

    NSS --> SS["SnippetService"]
    NTS --> TS["TagService"]
    NUS --> US["UserService"]

    SS --> SC["ServiceContainer"]
    TS --> SC
    US --> SC
    DB --> SC

    style SC fill:#4a90d9,stroke:#2c5f8a,color:#fff
    style DB fill:#f5a623,stroke:#c4841c,color:#fff
```

### Graceful Shutdown フロー

```mermaid
graph TD
    SIG["SIGINT / SIGTERM"] --> CTX["signal.NotifyContext<br/>ctx cancelled"]
    CTX --> EG["errgroup (並列実行)"]

    EG --> SHUTDOWN["e.Shutdown(ctx)<br/>HTTP サーバ停止"]
    EG --> DBCLOSE["sqlDB.Close()<br/>DB コネクション解放"]

    SHUTDOWN --> DONE["全完了"]
    DBCLOSE --> DONE

    style EG fill:#7ab648,stroke:#5a8a38,color:#fff
    style DONE fill:#4a90d9,stroke:#2c5f8a,color:#fff
```

---

## Project Structure

```
.
├── cmd/
│   └── main.go                 # エントリポイント (graceful shutdown)
├── app/
│   ├── config/                 # 環境変数の読み込み
│   │   └── config.go
│   ├── di/                     # Dependency Injection (Google Wire)
│   │   ├── wire.go
│   │   ├── wire_gen.go
│   │   └── service_container.go
│   ├── domain/                 # Domain 層 (最内層)
│   │   ├── model/              #   エンティティ: Snippet, Tag, User
│   │   └── repository/         #   Port: Repository インターフェース
│   ├── usecase/                # Usecase 層
│   │   ├── snippet_service.go  #   ビジネスロジック + timeout
│   │   ├── tag_service.go
│   │   └── user_service.go
│   ├── interface_adapter/      # Interface Adapter 層
│   │   ├── handler/            #   HTTP Handler (Echo)
│   │   │   ├── middleware/     #     CORS middleware
│   │   │   ├── snippet_handler.go
│   │   │   ├── tag_handler.go
│   │   │   ├── handler_helper.go  # エラーハンドリング
│   │   │   └── validator.go       # リクエストバリデーション
│   │   └── repository/         #   Repository 実装 (GORM)
│   │       ├── snippet_repository.go
│   │       ├── tag_repository.go
│   │       ├── user_repository.go
│   │       └── errors.go       #   GORM → Domain エラー変換
│   └── infrastructure/         # Infrastructure 層
│       ├── database/           #   MySQL 接続 (GORM)
│       └── aws/                #   AWS S3 クライアント
├── migrations/                 # DDL マイグレーション
├── Dockerfile                  # マルチステージビルド (Go 1.23)
├── docker-compose.yml          # ローカル開発用 MySQL
├── Makefile                    # ビルド・テスト・lint コマンド
└── learnings.md                # リファクタリング学習記録
```

---

## Features

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/snippets/:snippet_id` | スニペットを ID で取得 |
| `GET` | `/snippets/search?snippet_keyword=...` | キーワードでスニペット検索 |
| `GET` | `/snippets/tags/:tag_id` | タグに紐づくスニペット一覧 |
| `POST` | `/snippets` | スニペット作成 |
| `POST` | `/snippets/associate` | スニペットとタグの紐付け |
| `GET` | `/tags/:tag_id` | タグを ID で取得 |
| `GET` | `/tags/search?tag_keyword=...` | キーワードでタグ検索 |
| `POST` | `/tags` | タグ作成 |

### Tech Stack

| Category | Technology |
|----------|-----------|
| Language | Go 1.23 |
| HTTP Framework | Echo v4 |
| ORM | GORM (MySQL) |
| DI | Google Wire |
| Logging | log/slog (JSON) |
| Validation | go-playground/validator v10 |
| Concurrency | golang.org/x/sync/errgroup |
| Container | Docker (multi-stage, scratch) |

### Key Design Decisions

- **Dependency Inversion**: Domain 層の Repository Interface を Port とし、GORM 実装を Adapter として注入
- **Error Translation**: Repository 層で GORM/DB エラーをドメインエラーに変換し、Handler 層で HTTP ステータスにマッピング
- **Context Propagation**: 全レイヤーで `context.Context` を伝搬し、タイムアウト・キャンセルを DB クエリまで到達させる
- **Graceful Shutdown**: `errgroup` で HTTP サーバ停止と DB 接続解放を並列実行
- **Structured Logging**: `log/slog` による JSON 構造化ログ

---

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- MySQL 5.7+ (or use docker-compose)

### Setup

```bash
# 1. MySQL を起動
make db-start

# 2. マイグレーション実行 (要 golang-migrate)
migrate -path migrations -database "mysql://user:password@tcp(localhost:13306)/code_snippets_db" up

# 3. .env ファイルを作成
cat <<EOF > .env
TIMEOUT_SECOND=10
DBMS=mysql
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DBHOST=localhost
MYSQL_DBPORT=13306
MYSQL_DATABASE=code_snippets_db
EOF

# 4. アプリケーション起動
make run
```

### Build & Deploy

```bash
# ローカルビルド
make local-build

# Docker イメージビルド
make docker-build

# テスト実行
make test

# Lint
make lint
```
