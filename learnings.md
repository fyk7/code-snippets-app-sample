# Learnings: Code Snippets App リファクタリング記録

本ドキュメントは、ヘキサゴナルアーキテクチャ（クリーンアーキテクチャ）への改善と Go 2025 ベストプラクティス適用にあたって得られた学びをまとめたものです。

---

## 1. 致命的バグの発見と修正

### 1.1 テーブル名の誤り

**問題**: コピペに起因する、クエリ内のテーブル名間違い。

| ファイル | メソッド | 誤 | 正 |
|---------|---------|----|----|
| `tag_repository.go` | `GetByID` | `FROM snippet` | `FROM tag` |
| `tag_repository.go` | `FindByKeyWord` | `FROM snippet` | `FROM tag` |
| `user_repository.go` | `Create` | `INSERT INTO snippet` | `INSERT INTO user` |

**学び**: Repository を複製して新規作成する際は、テーブル名の置換漏れに注意。IDE の一括置換だけでなく、SQL 文のレビューを徹底する。

### 1.2 LIKE 句のパラメータバインド不具合

**問題**: `'%@keyword%'` のように記述すると、`@keyword` はリテラル文字列として扱われ、パラメータバインドが行われない。

```sql
-- NG: @keyword がリテラル扱い
WHERE title LIKE '%@keyword%'

-- OK: CONCAT で動的に結合
WHERE title LIKE CONCAT('%', @keyword, '%')
```

**学び**: GORM の Named Parameter は SQL 文の文字列リテラル内では展開されない。LIKE 句では `CONCAT` を使って外で結合する。

### 1.3 handleError の return 漏れ

**問題**: `c.Bind()` のエラー処理で `return handleError(c, err)` ではなく `handleError(c, err)` と書いていたため、エラー時も処理が続行されていた。

```go
// NG: return が無いため後続処理が実行される
if err := c.Bind(&req); err != nil {
    handleError(c, err)
}

// OK
if err := c.Bind(&req); err != nil {
    return handleError(c, err)
}
```

**学び**: Echo の handler は `error` を返す関数。エラーヘルパーの戻り値を `return` しないと、レスポンスは書き込まれるが処理は止まらない。

### 1.4 Service 層のメソッド呼び間違い

**問題**: `userService.Update()` 内部で `repo.Create()` を呼んでいた。

**学び**: Service が Repository に委譲するだけの薄いラッパーの場合でも、委譲先のメソッド名が正しいか確認する。テスト(mock)を書くことで早期発見できる。

### 1.5 タイムアウトの二重適用

**問題**: `config.LoadConf()` で `time.Duration(timeOut) * time.Second` を適用済みなのに、`main.go` でさらに `time.Duration(cfg.AppTimeOut * time.Second)` を掛けていた。

**学び**: `time.Duration` 型の値に再度 `time.Second` を掛けると、意図しない巨大な値になる。変換箇所を一箇所に集約する。

### 1.6 SQL 構文エラー（末尾カンマ）

**問題**: `user_repository.go` の INSERT 文で VALUES の末尾に不要なカンマがあった。

```sql
-- NG: 末尾カンマで SQL エラー
) VALUES (
  @userName,
  @now,
);

-- OK
) VALUES (
  @userName,
  @now
);
```

### 1.7 ルート定義とパラメータ取得の不一致

**問題**: `POST /snippets/associate` のルートにパスパラメータが無いのに、handler が `c.Param("snippet_id")` で取得しようとしていた。

**修正**: リクエストボディ（JSON）から `snippet_id`, `tag_id` を受け取るように変更。

---

## 2. クリーンアーキテクチャの改善

### 2.1 Context 伝搬の徹底

**問題**: Repository の全メソッドが `ctx context.Context` を受け取っていたが、GORM クエリに渡していなかった。

**修正**: 全クエリを `sr.Conn.WithContext(ctx).Raw(...)` / `.Exec(...)` に変更。

**効果**: タイムアウトやキャンセルが DB クエリレベルまで伝搬されるようになった。

### 2.2 不要な Handler Interface の削除

**問題**: `SnippetHandler`, `TagHandler` インターフェースが定義されていたが、どこからも参照されていなかった。

**学び**: ヘキサゴナルアーキテクチャにおいて、Interface Adapter 層（Handler）は外部からの入力を内部に変換する「アダプタ」であり、それ自体をインターフェースで抽象化する必要はない。インターフェースが必要なのは「内側の層が外側に依存しないようにする」ための Port（Repository Interface, Service Interface）。

### 2.3 ドメインエラーマッピングの導入

**問題**: Repository 層が GORM のエラーをそのまま返していたため、Handler の `errors.Is(err, model.ErrNotFound)` が永遠にマッチしなかった。

**修正**: `toDomainError()` ヘルパーを作成し、Repository 層の出口でエラー変換を行うようにした。

```
GORM Error → toDomainError() → Domain Error → Service → Handler
```

**学び**: ヘキサゴナルアーキテクチャでは、外部ライブラリのエラーをドメイン層のエラーに変換するのは Adapter 層（Repository 実装）の責務。

### 2.4 バリデーションの統合

**問題**: `ValidRequest()` 関数が定義されていたが、Handler から呼ばれていなかった。またリクエスト DTO に `validate` タグが無かった。

**修正**: リクエスト DTO に `validate:"required"` タグを追加し、全 POST/PUT Handler で `ValidRequest()` を呼び出すようにした。

---

## 3. Go 2025 ベストプラクティスの適用

### 3.1 Go バージョンの更新 (1.18 → 1.23)

- `any` 型エイリアスの活用 (`map[string]interface{}` → `map[string]any`)
- `log/slog` 標準ライブラリの利用 (Go 1.21+)
- `signal.NotifyContext` による簡潔なシグナルハンドリング

### 3.2 log/slog への移行

**Before**: `log.Fatal()` (標準) + `logrus` (外部依存、未使用)
**After**: `log/slog` (Go 1.21+ 標準) + JSON ハンドラ

```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})))
```

**メリット**:
- 外部依存の削減
- 構造化ログが標準で利用可能
- パフォーマンスが良い（遅延評価）

### 3.3 Validator v9 → v10

`gopkg.in/go-playground/validator.v9` は非推奨。`github.com/go-playground/validator/v10` に移行。

**変更点**:
- import パスの変更
- バリデータインスタンスをパッケージ変数で共有（毎回 `New()` しない）

### 3.4 Graceful Shutdown

**Before**: `log.Fatal(e.Start(":8080"))` — シグナルを受けると即座に終了。

**After**: `errgroup` + `signal.NotifyContext` による協調的シャットダウン。

```
SIGINT/SIGTERM
  ↓
signal.NotifyContext で ctx がキャンセル
  ↓
errgroup で並列に:
  - e.Shutdown() (HTTP サーバの graceful 停止)
  - sqlDB.Close() (DB コネクションプールの解放)
```

---

## 4. 並列化

### 4.1 Graceful Shutdown の並列化

HTTP サーバの停止と DB コネクションのクローズは独立した操作なので、`errgroup` で並列実行。

```go
sg, _ := errgroup.WithContext(shutdownCtx)
sg.Go(func() error { return e.Shutdown(shutdownCtx) })
sg.Go(func() error {
    sqlDB, err := serviceContainer.DB.DB()
    if err != nil { return err }
    return sqlDB.Close()
})
return sg.Wait()
```

### 4.2 サーバ起動とシャットダウン監視の並列化

メイン goroutine で HTTP サーバ起動とシャットダウンシグナル監視を `errgroup` で並列管理。どちらかがエラーを返すと全体が停止する。

---

## 5. 今後の改善候補

- [ ] テストの追加（特に Repository 層の統合テスト）
- [ ] 認証・認可の実装（現在 UserID がダミー値）
- [ ] S3Client の DI 統合（定義されているが未使用）
- [ ] DB マイグレーションの自動実行
- [ ] OpenAPI (Swagger) ドキュメントの生成
- [ ] `errgroup` を活用した並列データ取得（例: Snippet + 関連 Tag の同時取得）
