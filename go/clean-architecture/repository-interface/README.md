# repository-interface

テーマ：**UseCaseからInfrastructureの具象Repositoryを直接参照する場合と、Repository interfaceを挟む場合の違い**

Good/Badどちらも、まったく同じ機能（`Score` を取得して合否判定する `CheckPass(id int) (bool, error)`）を実装しています。実装が違うだけで、外から見た挙動は同じです。

## Bad

```text
UseCase
   ↓
persistence.ScoreRepository
   ↓
PostgreSQL
```

[bad/application/score_usecase.go](bad/application/score_usecase.go) を見ると、`ScoreUseCase` が `*persistence.ScoreRepository` という具体的な型を直接持っています。

```go
type ScoreUseCase struct {
    repository *persistence.ScoreRepository
}
```

UseCaseは「PostgreSQL実装そのもの」を知っている状態です。

## Good

```text
UseCase
   ↓
domain.ScoreRepository (interface)
   ↑
persistence.ScoreRepository
```

[good/application/score_usecase.go](good/application/score_usecase.go) では、`ScoreUseCase` は `domain.ScoreRepository` というinterfaceだけを持っています。

```go
type ScoreUseCase struct {
    repository domain.ScoreRepository
}
```

UseCaseが知っているのは「`Find(id int) (*Score, error)` ができるRepositoryが存在する」ということだけです。実装がPostgreSQLなのかMockなのかを知りません。

## interfaceを挟む具体的なメリット

「依存したくないからinterfaceを使う」ではなく、実際に何が変わるのかをコードとテストで確認できるようにしています。

### メリット1：UseCaseをDBなしでテストできる

- Good: [good/application/score_usecase_test.go](good/application/score_usecase_test.go) は `MockScoreRepository` を使っており、PostgreSQLを一切必要としません。
- Bad: [bad/application/score_usecase_test.go](bad/application/score_usecase_test.go) は `*persistence.ScoreRepository` を直接生成するしかないため、テスト実行にPostgreSQLが必須です。

Badで同じようにMockを差し込もうとすると型が合わずコンパイルできません。具体的なコード例は [bad/application/mock_cannot_be_injected.txt](bad/application/mock_cannot_be_injected.txt) を参照してください。

### メリット2：実装を交換できる

Goodでは `domain.ScoreRepository` を満たしてさえいれば、UseCaseを変更せずに実装を差し替えられます。

```text
ScoreUseCase
      ↓
domain.ScoreRepository
      ↑
 ┌────┴────┐
 ↓         ↓
persistence.ScoreRepository (PostgreSQL)
MockScoreRepository (テスト用)
```

Badでは `*persistence.ScoreRepository` という具体型に固定されているため、別の実装に差し替えるにはUseCase自体を書き換える必要があります。

### メリット3：Infrastructureの変更がUseCaseに波及しにくい

Goodで `persistence.ScoreRepository`（PostgreSQL実装）を書き換えても、`domain.ScoreRepository` interfaceさえ満たしていれば `application` パッケージのコードは一切変更不要です。Badでは `persistence` の型そのものにUseCaseが依存しているため、Infrastructure側の変更がUseCaseの型定義にまで波及する可能性があります。

## 実行方法

PostgreSQL起動（リポジトリのルートで）：

```bash
docker compose up -d
```

テスト（リポジトリのルートで）：

```bash
go test ./...
```

Goodの実行：

```bash
go run ./go/clean-architecture/repository-interface/good
```

Badの実行：

```bash
go run ./go/clean-architecture/repository-interface/bad
```

## 演習

### Exercise 1: MockだけでUseCaseテストが通ることを確認する

PostgreSQLを止めた状態（`docker compose down`）で、以下を実行してください。

```bash
go test ./go/clean-architecture/repository-interface/good/application/...
```

DBが動いていなくても成功するはずです。これがGoodの`ScoreUseCase`がMockだけでテストできている証拠です。

### Exercise 2: BadはPostgreSQLが必須であることを確認する

同じくPostgreSQLを止めた状態で、以下を実行してください。

```bash
go test ./go/clean-architecture/repository-interface/bad/application/...
```

`postgres is required for this test` のようなエラーで失敗するはずです。`docker compose up -d` してから再実行すると成功します。

### Exercise 3: MemoryScoreRepositoryを追加する

`good/persistence/score_repository.go` とは別に、DBを使わないインメモリ実装を自分で追加してみてください。

```go
type MemoryScoreRepository struct {
    scores map[int]*domain.Score
}

func (r *MemoryScoreRepository) Find(id int) (*domain.Score, error) {
    // ...
}
```

これが `domain.ScoreRepository` を満たしていれば、`application.NewScoreUseCase` にそのまま渡せます。`good/application/score_usecase.go` を一切変更せずに動くことを確認してください。

### Exercise 4: コードジャンプで依存関係を追う

VS Codeで `ScoreUseCase` の `repository` フィールドから `Find` の呼び出しへ、`Ctrl+Click`（またはF12）でジャンプしてください。

- Goodでは `domain.ScoreRepository` interfaceの定義にジャンプします。そこから実装（`persistence.ScoreRepository` や `MockScoreRepository`）へは、さらに「実装を探す」操作（VS Codeなら `Go to Implementations`）が必要です。
- Badでは最初から `persistence.ScoreRepository` という具体型に直接ジャンプします。

このジャンプ先の違いが、「UseCaseが何を知っているか」の違いをそのまま表しています。
