# Code Snippets API 仕様書

このディレクトリには、Code Snippets API の OpenAPI 仕様書が含まれています。

## ファイル構造

```
openapi/
├── openapi.yaml                              # メインのOpenAPI仕様書
├── README.md                                 # 本ドキュメント
├── schemas/                                  # スキーマ定義
│   ├── common/                               # 共通スキーマ
│   │   ├── Error.yaml                       # エラーレスポンス
│   │   └── OKResponse.yaml                  # 成功レスポンス
│   ├── snippets/                             # スニペット関連
│   │   ├── Snippet.yaml                     # スニペットスキーマ
│   │   ├── SnippetCreateRequest.yaml        # スニペット作成リクエスト
│   │   ├── SnippetUpdateRequest.yaml        # スニペット更新リクエスト
│   │   └── AssociateWithTagRequest.yaml     # タグ紐付けリクエスト
│   └── tags/                                 # タグ関連
│       ├── Tag.yaml                         # タグスキーマ
│       └── TagCreateRequest.yaml            # タグ作成リクエスト
└── paths/                                    # APIパス定義
    ├── snippets/                             # スニペット関連パス
    │   ├── index.yaml                       # POST /snippets, PUT /snippets
    │   ├── snippetId.yaml                   # GET /snippets/{snippet_id}
    │   ├── search.yaml                      # GET /snippets/search
    │   ├── by-tag.yaml                      # GET /snippets/tags/{tag_id}
    │   └── associate.yaml                   # POST /snippets/associate
    └── tags/                                 # タグ関連パス
        ├── index.yaml                       # POST /tags
        ├── tagId.yaml                       # GET /tags/{tag_id}
        └── search.yaml                      # GET /tags/search
```

## 設計原則

### ファイル分割の方針

1. **URL階層に基づく分割**: 同じルートパスを持つエンドポイントは同一ディレクトリにまとめる
2. **機能単位での分離**: 関連するスキーマ・パスは同じディレクトリで管理
3. **DRY原則**: 共通のレスポンス定義は `schemas/common/` で一元管理し `$ref` で参照

### 命名規約

- **ディレクトリ**: 複数形・小文字 (例: `snippets`, `tags`)
- **スキーマファイル**: PascalCase (例: `Snippet.yaml`, `TagCreateRequest.yaml`)
- **パスファイル**: kebab-case (例: `by-tag.yaml`, `associate.yaml`)
- **パスパラメータファイル**: camelCase (例: `snippetId.yaml`, `tagId.yaml`)

### $ref 参照の管理

- **メインファイルから**: `./schemas/snippets/Snippet.yaml`
- **同一ディレクトリ内**: `./FileName.yaml`
- **異なるディレクトリ**: `../../schemas/common/Error.yaml`

## 使用方法

### バリデーション

```bash
# swagger-cli でバリデーション
npx @redocly/cli lint openapi.yaml

# または
npx swagger-cli validate openapi.yaml
```

### Swagger UI でプレビュー

```bash
# Docker で Swagger UI を起動
docker run -p 8081:8080 \
  -e SWAGGER_JSON=/openapi/openapi.yaml \
  -v $(pwd)/openapi:/openapi \
  swaggerapi/swagger-ui

# http://localhost:8081 でアクセス
```

### バンドル (単一ファイル生成)

```bash
npx @redocly/cli bundle openapi.yaml -o bundled.yaml
```

## メンテナンス

### 新しいエンドポイントの追加

1. `schemas/` に必要なスキーマファイルを作成
2. `paths/` にパス定義ファイルを作成
3. `openapi.yaml` に `$ref` を追加

### スキーマの変更

1. 対象のスキーマファイルを編集
2. `npx @redocly/cli lint openapi.yaml` で検証
