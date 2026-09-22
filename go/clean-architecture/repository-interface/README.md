# repository-interface

テーマ：**UseCaseからInfrastructureの具象Repositoryを直接参照する場合と、Repository interfaceを挟む場合の違い**

Good/Bad/nodiのすべてが、まったく同じ機能（箱の中の数字を取得して4以上かどうか判定する `IsLarge(id int) (bool, error)`）を実装しています。実装が違うだけで、外から見た挙動は同じです。

## Bad

```text
UseCase
   ↓
persistence.BoxRepository
   ↓
PostgreSQL
```

[bad/application/box_usecase.go](bad/application/box_usecase.go) を見ると、`BoxUseCase` が `*persistence.BoxRepository` という具体的な型を直接持っています。

```go
type BoxUseCase struct {
    repository *persistence.BoxRepository
}
```

UseCaseは「PostgreSQL実装そのもの」を知っている状態です。ただし依存先を外から受け取る（DI）こと自体はしています。

## Good

```text
UseCase
   ↓
domain.BoxRepository (interface)
   ↑
persistence.BoxRepository
```

[good/application/box_usecase.go](good/application/box_usecase.go) では、`BoxUseCase` は `domain.BoxRepository` というinterfaceだけを持っています。

```go
type BoxUseCase struct {
    repository domain.BoxRepository
}
```

UseCaseが知っているのは「`Find(id int) (*Box, error)` ができるRepositoryが存在する」ということだけです。実装がPostgreSQLなのかMockなのかを知りません。

## nodi（アンチパターン）

```text
UseCase
   ↓ (自分でDB接続を組み立てる)
persistence.BoxRepository
   ↓
PostgreSQL
```

[nodi/application/box_usecase.go](nodi/application/box_usecase.go) は、そもそもDI（Dependency Injection）自体をしていません。

```go
func NewBoxUseCase() *BoxUseCase {
    dsn := os.Getenv("DATABASE_URL")
    db, _ := sql.Open("postgres", dsn)
    repository := persistence.NewBoxRepository(db)
    return &BoxUseCase{repository: repository}
}
```

コンストラクタが引数を取らず、DB接続を自分の内部で組み立てています。Badは「具体型を外から注入する」ことはできていましたが、nodiは「注入する」という発想自体がありません。詳細は [nodi/application/why_this_is_worse.txt](nodi/application/why_this_is_worse.txt) を参照してください。

## interfaceを挟む具体的なメリット

「依存したくないからinterfaceを使う」ではなく、実際に何が変わるのかをコードとテストで確認できるようにしています。

### メリット1：UseCaseをDBなしでテストできる

- Good: [good/application/box_usecase_test.go](good/application/box_usecase_test.go) は `MockBoxRepository` を使っており、PostgreSQLを一切必要としません。
- Bad: [bad/application/box_usecase_test.go](bad/application/box_usecase_test.go) は `*persistence.BoxRepository` を直接生成するしかないため、テスト実行にPostgreSQLが必須です。
- nodi: [nodi/application/box_usecase_test.go](nodi/application/box_usecase_test.go) はテスト用DBに向けることすらできません（接続先がコンストラクタ内に固定されているため）。

Badで同じようにMockを差し込もうとすると型が合わずコンパイルできません。実際にコンパイルエラーになる例が [bad/application/try_mock.go](bad/application/try_mock.go) です（このファイルは意図的にコンパイルを失敗させています。詳細は [bad/application/mock_cannot_be_injected.txt](bad/application/mock_cannot_be_injected.txt) を参照）。

### メリット2：実装を交換できる

Goodでは `domain.BoxRepository` を満たしてさえいれば、UseCaseを変更せずに実装を差し替えられます。

```text
BoxUseCase
      ↓
domain.BoxRepository
      ↑
 ┌────┴────┐
 ↓         ↓
persistence.BoxRepository (PostgreSQL)
MockBoxRepository (テスト用)
```

Badでは `*persistence.BoxRepository` という具体型に固定されているため、別の実装に差し替えるにはUseCase自体を書き換える必要があります。nodiではさらに、接続先DBを差し替える余地すらありません。

### メリット3：Infrastructureの変更がUseCaseに波及しにくい

Goodで `persistence.BoxRepository`（PostgreSQL実装）を書き換えても、`domain.BoxRepository` interfaceさえ満たしていれば `application` パッケージのコードは一切変更不要です。Badでは `persistence` の型そのものにUseCaseが依存しているため、Infrastructure側の変更がUseCaseの型定義にまで波及する可能性があります。

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

Badの実行（`bad/application` には意図的にコンパイルエラーになる `try_mock.go` が置いてあるため、通常の `go run ./bad` はエラーになります。`main.go` 単体で試す場合は下記）：

```bash
go run ./go/clean-architecture/repository-interface/bad/main.go
```

nodiの実行：

```bash
go run ./go/clean-architecture/repository-interface/nodi
```

## 演習

### Exercise 1: MockだけでUseCaseテストが通ることを確認する

PostgreSQLを止めた状態（`docker compose down`）で、以下を実行してください。

```bash
go test ./go/clean-architecture/repository-interface/good/application/...
```

DBが動いていなくても成功するはずです。これがGoodの`BoxUseCase`がMockだけでテストできている証拠です。

### Exercise 2: BadはPostgreSQLが必須であることを確認する

同じくPostgreSQLを止めた状態で、以下を実行してください。

```bash
go test ./go/clean-architecture/repository-interface/bad/persistence/...
```

`postgres is required for this test` のようなエラーで失敗するはずです。`docker compose up -d` してから再実行すると成功します。

### Exercise 3: MemoryBoxRepositoryを追加する

`good/persistence/box_repository.go` とは別に、DBを使わないインメモリ実装を自分で追加してみてください。

```go
type MemoryBoxRepository struct {
    boxes map[int]*domain.Box
}

func (r *MemoryBoxRepository) Find(id int) (*domain.Box, error) {
    // ...
}
```

これが `domain.BoxRepository` を満たしていれば、`application.NewBoxUseCase` にそのまま渡せます。`good/application/box_usecase.go` を一切変更せずに動くことを確認してください。

### Exercise 4: コードジャンプで依存関係を追う

VS Codeで `BoxUseCase` の `repository` フィールドから `Find` の呼び出しへ、`Ctrl+Click`（またはF12）でジャンプしてください。

- Goodでは `domain.BoxRepository` interfaceの定義にジャンプします。そこから実装（`persistence.BoxRepository` や `MockBoxRepository`）へは、さらに「実装を探す」操作（VS Codeなら `Go to Implementations`）が必要です。
- Badでは最初から `persistence.BoxRepository` という具体型に直接ジャンプします。
- nodiでは `repository` フィールドを外から渡す箇所自体が存在しません（コンストラクタの中で完結しています）。

このジャンプ先の違いが、「UseCaseが何を知っているか」の違いをそのまま表しています。
