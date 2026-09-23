# repository-interface

テーマ：**UseCaseからInfrastructureの具象Repositoryを直接参照する場合と、Repository interfaceを挟む場合の違い**

3つのディレクトリすべてが、まったく同じ機能（箱の中の数字を取得して4以上かどうか判定する `IsLarge(id int) (bool, error)`）を実装しています。実装が違うだけで、外から見た挙動は同じです。

| ディレクトリ | 設計 | `NewBoxUseCase`の形 |
|---|---|---|
| [1-no-di](1-no-di/) | DIをしない | `NewBoxUseCase()` |
| [2-di-concrete](2-di-concrete/) | DIするが、具体型で受け取る | `NewBoxUseCase(repository *persistence.CachedBoxRepository)` |
| [3-di-interface](3-di-interface/) | DIして、interfaceで受け取る | `NewBoxUseCase(repository domain.BoxRepository)` |

**3つの全コードを並べた比較は [comparison.md](comparison.md) にあります。** まずはそちらを読むのが分かりやすいです。

## 1-no-di（DIをしない）

```text
UseCase
   ↓ (自分でDB接続を組み立てる)
persistence.BoxRepository
   ↓
PostgreSQL
```

[1-no-di/application/box_usecase.go](1-no-di/application/box_usecase.go) は、そもそもDI（Dependency Injection）自体をしていません。

```go
func NewBoxUseCase() *BoxUseCase {
    inner := persistence.NewDefaultBoxRepository()
    return &BoxUseCase{repository: persistence.NewCachedBoxRepository(inner)}
}
```

コンストラクタが引数を取らず、依存先を自分の内部で組み立てています。「外から何かを渡す窓口」自体が存在しません。詳細は [1-no-di/application/why_this_is_worse.txt](1-no-di/application/why_this_is_worse.txt) を参照してください。

## 2-di-concrete（DIするが具体型）

```text
UseCase
   ↓
persistence.CachedBoxRepository
   ↓
PostgreSQL
```

[2-di-concrete/application/box_usecase.go](2-di-concrete/application/box_usecase.go) を見ると、`BoxUseCase` が `*persistence.CachedBoxRepository` という具体的な型を直接持っています。

```go
type BoxUseCase struct {
    repository *persistence.CachedBoxRepository
}
```

UseCaseは「PostgreSQL実装そのもの」を知っている状態です。ただし依存先を外から受け取る（DI）こと自体はしているので、「どのDBに繋いだrepositoryを渡すか」は呼び出し側が選べます。

## 3-di-interface（DIしてinterfaceで受け取る）

```text
UseCase
   ↓
domain.BoxRepository (interface)
   ↑
persistence.BoxRepository / CachedBoxRepository / MockBoxRepository
```

[3-di-interface/application/box_usecase.go](3-di-interface/application/box_usecase.go) では、`BoxUseCase` は `domain.BoxRepository` というinterfaceだけを持っています。

```go
type BoxUseCase struct {
    repository domain.BoxRepository
}
```

UseCaseが知っているのは「`Find(id int) (*Box, error)` ができるRepositoryが存在する」ということだけです。実装がPostgreSQLなのかMockなのかを知りません。

## interfaceを挟む具体的なメリット

「依存したくないからinterfaceを使う」ではなく、実際に何が変わるのかをコードとテストで確認できるようにしています。

### メリット1：UseCaseをDBなしでテストできる

- 3-di-interface: [3-di-interface/application/box_usecase_test.go](3-di-interface/application/box_usecase_test.go) は `MockBoxRepository` を使っており、PostgreSQLを一切必要としません。
- 2-di-concrete: [2-di-concrete/application/box_usecase_test.go](2-di-concrete/application/box_usecase_test.go) は `*persistence.CachedBoxRepository` を直接生成するしかないため、テスト実行にPostgreSQLが必須です。
- 1-no-di: [1-no-di/application/box_usecase_test.go](1-no-di/application/box_usecase_test.go) はテスト用DBに向けることすらできません（接続先がコンストラクタ内に固定されているため）。

2-di-concreteで同じようにMockを差し込もうとすると型が合わずコンパイルできません。実際にコンパイルエラーになる例が [2-di-concrete/application/try_mock.go](2-di-concrete/application/try_mock.go) です（このファイルは意図的にコンパイルを失敗させています。詳細は [2-di-concrete/application/mock_cannot_be_injected.txt](2-di-concrete/application/mock_cannot_be_injected.txt) を参照）。

### メリット2：実装を交換できる

3-di-interfaceでは `domain.BoxRepository` を満たしてさえいれば、UseCaseを変更せずに実装を差し替えられます。

```text
BoxUseCase
      ↓
domain.BoxRepository
      ↑
 ┌────┼────────────┐
 ↓    ↓            ↓
persistence.BoxRepository (PostgreSQL)
persistence.CachedBoxRepository (キャッシュ付き)
MockBoxRepository (テスト用)
```

2-di-concreteでは `*persistence.CachedBoxRepository` という具体型に固定されているため、別の実装に差し替えるにはUseCase自体を書き換える必要があります。1-no-diではさらに、接続先DBを差し替える余地すらありません。

### メリット3：Infrastructureの変更がUseCaseに波及しにくい

これは実際にキャッシュ機能を後から追加して検証しました。3-di-interfaceでは `main.go` の1行を変えただけで、`application` パッケージには一切手を入れずに済みました。一方、2-di-concreteと1-no-diでは `application/box_usecase.go` の型定義まで書き換える必要がありました。変更範囲の比較は [comparison.md](comparison.md) を参照してください。

## 実行方法

PostgreSQL起動（リポジトリのルートで）：

```bash
docker compose up -d
```

テスト（リポジトリのルートで）：

```bash
go test ./...
```

3-di-interfaceの実行：

```bash
go run ./go/clean-architecture/repository-interface/3-di-interface
```

2-di-concreteの実行（`2-di-concrete/application` には意図的にコンパイルエラーになる `try_mock.go` が置いてあるため、パッケージ指定の `go run` はエラーになります。`main.go` 単体で試す場合は下記）：

```bash
go run ./go/clean-architecture/repository-interface/2-di-concrete/main.go
```

1-no-diの実行：

```bash
go run ./go/clean-architecture/repository-interface/1-no-di
```

## 演習

### Exercise 1: MockだけでUseCaseテストが通ることを確認する

PostgreSQLを止めた状態（`docker compose down`）で、以下を実行してください。

```bash
go test ./go/clean-architecture/repository-interface/3-di-interface/application/...
```

DBが動いていなくても成功するはずです。これが`BoxUseCase`がMockだけでテストできている証拠です。

### Exercise 2: 2-di-concreteはPostgreSQLが必須であることを確認する

同じくPostgreSQLを止めた状態で、以下を実行してください。

```bash
go test ./go/clean-architecture/repository-interface/2-di-concrete/persistence/...
```

`postgres is required for this test` のようなエラーで失敗するはずです。`docker compose up -d` してから再実行すると成功します。

### Exercise 3: MemoryBoxRepositoryを追加する

`3-di-interface/persistence/box_repository.go` とは別に、DBを使わないインメモリ実装を自分で追加してみてください。

```go
type MemoryBoxRepository struct {
    boxes map[int]*domain.Box
}

func (r *MemoryBoxRepository) Find(id int) (*domain.Box, error) {
    // ...
}
```

これが `domain.BoxRepository` を満たしていれば、`application.NewBoxUseCase` にそのまま渡せます。`3-di-interface/application/box_usecase.go` を一切変更せずに動くことを確認してください。

### Exercise 4: コードジャンプで依存関係を追う

VS Codeで `BoxUseCase` の `repository` フィールドから `Find` の呼び出しへ、`Ctrl+Click`（またはF12）でジャンプしてください。

- 3-di-interfaceでは `domain.BoxRepository` interfaceの定義にジャンプします。そこから実装（`persistence.BoxRepository` や `MockBoxRepository`）へは、さらに「実装を探す」操作（VS Codeなら `Go to Implementations`）が必要です。
- 2-di-concreteでは最初から `persistence.CachedBoxRepository` という具体型に直接ジャンプします。
- 1-no-diでは `repository` フィールドを外から渡す箇所自体が存在しません（コンストラクタの中で完結しています）。

このジャンプ先の違いが、「UseCaseが何を知っているか」の違いをそのまま表しています。
