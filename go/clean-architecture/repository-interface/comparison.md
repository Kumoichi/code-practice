# 3つの実装を全部並べて比較する

このディレクトリには、**まったく同じ機能**を3通りの設計で実装したものが入っています。

| ディレクトリ | 設計 | `NewBoxUseCase`の形 |
|---|---|---|
| [1-no-di](1-no-di/) | DIをしない | `NewBoxUseCase()` |
| [2-di-concrete](2-di-concrete/) | DIするが、具体型で受け取る | `NewBoxUseCase(repository *persistence.CachedBoxRepository)` |
| [3-di-interface](3-di-interface/) | DIして、interfaceで受け取る | `NewBoxUseCase(repository domain.BoxRepository)` |

機能はどれも同じです。「箱の番号(id)を渡すと、その箱に入っている数字が4以上かどうかを返す」だけ。3つとも実行結果は次のように完全に一致します。

```
id=1 large=true
id=2 large=false
id=3 large=true
```

**外から見た挙動が同じなのに、コードの形が違うとどれだけ困るか**——それがこのディレクトリで見たいことです。

---

## 1. application層（ビジネスロジック）の比較

ここが3つの設計で最も差が出る場所です。まず全文を並べます。

### 1-no-di

```go
package application

import "code-practice/go/clean-architecture/repository-interface/1-no-di/persistence"

type BoxUseCase struct {
	repository *persistence.CachedBoxRepository
}

// 引数が0個。外から何かを渡す窓口が存在しない。
func NewBoxUseCase() *BoxUseCase {
	inner := persistence.NewDefaultBoxRepository()
	return &BoxUseCase{repository: persistence.NewCachedBoxRepository(inner)}
}

func (u *BoxUseCase) IsLarge(id int) (bool, error) {
	box, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}
	return box.Number >= 4, nil
}
```

### 2-di-concrete

```go
package application

import "code-practice/go/clean-architecture/repository-interface/2-di-concrete/persistence"

type BoxUseCase struct {
	repository *persistence.CachedBoxRepository
}

// 引数はあるが、型が具体的な構造体で固定されている。
func NewBoxUseCase(repository *persistence.CachedBoxRepository) *BoxUseCase {
	return &BoxUseCase{repository: repository}
}

func (u *BoxUseCase) IsLarge(id int) (bool, error) {
	box, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}
	return box.Number >= 4, nil
}
```

### 3-di-interface

```go
package application

import "code-practice/go/clean-architecture/repository-interface/3-di-interface/domain"

type BoxUseCase struct {
	repository domain.BoxRepository
}

// 引数の型がinterface。「Findできる何か」でありさえすればよい。
func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase {
	return &BoxUseCase{repository: repository}
}

func (u *BoxUseCase) IsLarge(id int) (bool, error) {
	box, err := u.repository.Find(id)
	if err != nil {
		return false, err
	}
	return box.Number >= 4, nil
}
```

### 違いはどこか

`IsLarge`の中身は**3つとも1文字も違いません**。違うのは次の2行だけです。

```go
// フィールドの型
repository *persistence.CachedBoxRepository   // 1-no-di, 2-di-concrete
repository domain.BoxRepository               // 3-di-interface

// コンストラクタの形
NewBoxUseCase()                                              // 1-no-di
NewBoxUseCase(repository *persistence.CachedBoxRepository)   // 2-di-concrete
NewBoxUseCase(repository domain.BoxRepository)               // 3-di-interface
```

そして注目すべきは**importしているパッケージ**です。

- `1-no-di` と `2-di-concrete` は `persistence` をimportしている → ビジネスロジックの層が、PostgreSQL実装のパッケージを直接知っている
- `3-di-interface` は `domain` をimportしている → PostgreSQLのことを一切知らない

この1行のimportの違いが、以降のすべての差を生みます。

---

## 2. 「同じ機能を後から足す」と何が起きるか

ここが一番の見どころです。3つとも**あとからキャッシュ機能を追加**しました。同じ要求に対して、それぞれ何行書き換える羽目になったかを比べます。

この比較は文章だけでなく、**gitのコミットとしても分けてあります**。いったん[3つすべてからキャッシュを取り除いた状態](https://github.com/Kumoichi/code-practice/commit/e188109)を作り、そこから同じ機能をディレクトリごとに1コミットずつ追加し直しました。各コミットの差分がそのまま答えになっているので、GitHubで開いて「変更されたファイルの一覧」を見比べるのが一番早いです。

| 設計 | コミット | 変更 |
|---|---|---|
| 3-di-interface | [b5c6223](https://github.com/Kumoichi/code-practice/commit/b5c6223) | 3ファイル（+67/-1） |
| 2-di-concrete | [2eb75ec](https://github.com/Kumoichi/code-practice/commit/2eb75ec) | 5ファイル（+58/-16） |
| 1-no-di | [ccde1f3](https://github.com/Kumoichi/code-practice/commit/ccde1f3) | 3ファイル（+38/-5） |

3-di-interfaceの`+67`はほとんどが新規ファイル2つ（キャッシュ実装とそのテスト）で、既存コードへの変更は`-1`が示す通り1行だけ、という内訳です。

### 3-di-interface: application層は無変更で済んだ

→ [コミット b5c6223](https://github.com/Kumoichi/code-practice/commit/b5c6223)

```go
// main.go — この1行だけ変更
- repository := persistence.NewBoxRepository(db)
+ repository := persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))
```

```go
// application/box_usecase.go — 変更なし（1文字も触っていない）
// domain/box_repository.go  — 変更なし
// application/box_usecase_test.go — 変更なし
```

新しく追加したのは [cached_box_repository.go](3-di-interface/persistence/cached_box_repository.go) と [そのテスト](3-di-interface/persistence/cached_box_repository_test.go) の2ファイルだけで、**既存ファイルへの変更は`main.go`の1行のみ**です。コミットをGitHubで開くと、変更ファイルの一覧に`application/box_usecase.go`が**そもそも載っていません**。

なぜ`BoxUseCase`を触らずに済んだのか。`repository`フィールドの型が`domain.BoxRepository`というinterfaceで、`CachedBoxRepository`も`Find`を持っている＝そのinterfaceを満たしているからです。`BoxUseCase`から見れば「`Find`できる何か」が来ることに変わりはなく、それがキャッシュ付きかどうかは関心の外にあります。

### 2-di-concrete: application層の型定義まで書き換えが必要だった

→ [コミット 2eb75ec](https://github.com/Kumoichi/code-practice/commit/2eb75ec)

```go
// main.go
- repository := persistence.NewBoxRepository(db)
+ repository := persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))
```

```go
// application/box_usecase.go — 型を2箇所書き換える羽目になった
  type BoxUseCase struct {
-     repository *persistence.BoxRepository
+     repository *persistence.CachedBoxRepository
  }

- func NewBoxUseCase(repository *persistence.BoxRepository) *BoxUseCase {
+ func NewBoxUseCase(repository *persistence.CachedBoxRepository) *BoxUseCase {
```

```go
// application/box_usecase_test.go — 引数の型が変わったので、呼び出し3箇所すべて書き換え
- useCase := NewBoxUseCase(persistence.NewBoxRepository(db))
+ useCase := NewBoxUseCase(persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db)))
```

DIはしているのに、**ビジネスロジックの層とそのテストにまで変更が波及**しました。`*persistence.BoxRepository`という具体的な構造体を名指ししていたせいで、「別の構造体で包む」という変更がそのまま型の不一致になってしまうからです。

### 1-no-di: 同じく application層の書き換えが必要だった

→ [コミット ccde1f3](https://github.com/Kumoichi/code-practice/commit/ccde1f3)

```go
// application/box_usecase.go
  type BoxUseCase struct {
-     repository *persistence.BoxRepository
+     repository *persistence.CachedBoxRepository
  }

  func NewBoxUseCase() *BoxUseCase {
-     return &BoxUseCase{repository: persistence.NewDefaultBoxRepository()}
+     inner := persistence.NewDefaultBoxRepository()
+     return &BoxUseCase{repository: persistence.NewCachedBoxRepository(inner)}
  }
```

型定義に加えて、**コンストラクタの中身そのもの**も書き換えています。組み立て方をUseCaseが自分で抱え込んでいるので、組み立て方が変わればUseCaseが変わる、という構造になっています。

### 変更範囲まとめ

| | main.go | application層 | テスト | 新規ファイル | 実際の差分 |
|---|---|---|---|---|---|
| 3-di-interface | 1行 | **変更なし** | **変更なし** | 2つ | [b5c6223](https://github.com/Kumoichi/code-practice/commit/b5c6223) |
| 2-di-concrete | 1行 | 型を2箇所 | 3箇所 | 1つ | [2eb75ec](https://github.com/Kumoichi/code-practice/commit/2eb75ec) |
| 1-no-di | （なし） | 型1箇所＋処理2行 | （元々テスト不能） | （同ファイル内に追加） | [ccde1f3](https://github.com/Kumoichi/code-practice/commit/ccde1f3) |

---

## 3. テストの書き方の比較

### 3-di-interface: DBなしで書ける

```go
func TestBoxUseCase_IsLarge_5は4以上なのでtrue(t *testing.T) {
	mock := &MockBoxRepository{Box: &domain.Box{ID: 1, Number: 5}}
	useCase := NewBoxUseCase(mock) // ← Mockをそのまま注入できる

	got, err := useCase.IsLarge(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Errorf("IsLarge() = false, want true")
	}
}
```

PostgreSQLを止めていても通ります。「5が返ってきたとき`IsLarge`はtrueを返すか」という**ロジックだけ**を確認しています。

### 2-di-concrete: 本物のDBが必須

```go
func TestBoxUseCase_IsLarge_5は4以上なのでtrue(t *testing.T) {
	db := openTestDB(t) // ← PostgreSQLに接続できないとここで失敗
	useCase := NewBoxUseCase(persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db)))

	got, err := useCase.IsLarge(1) // seed data: id=1, number=5
	...
}
```

Mockを渡そうとするとコンパイルエラーになります。実際に試したコードが [2-di-concrete/application/try_mock.go](2-di-concrete/application/try_mock.go) にあり、`go build`すると必ず次のエラーが出ます（このファイルは意図的にコンパイルを失敗させています）。

```
cannot use mock (variable of type *MockBoxRepository)
as *persistence.CachedBoxRepository value in argument to NewBoxUseCase
```

`MockBoxRepository`は`Find`という同じ形のメソッドを持っているのに拒否されます。Goが見ているのは「同じメソッドを持つか」ではなく「その構造体そのものか」だからです。

### 1-no-di: 接続先すら選べない

```go
func TestBoxUseCase_IsLarge_5は4以上なのでtrue(t *testing.T) {
	useCase := NewBoxUseCase() // ← 引数なし。Mockを渡す場所が存在しない

	got, err := useCase.IsLarge(1)
	if err != nil {
		t.Skipf("postgres is required for this test ...", err)
	}
	...
}
```

Mockが渡せないだけでなく、「テスト用の別DBに向ける」こともできません。接続先が`NewBoxUseCase`の**内側**に埋め込まれているからです。結果として、PostgreSQLが起動していないとテストは`t.Skipf`で黙ってスキップされます——つまり**落ちもしないが、確認もされない**状態になります。

---

## 4. 依存の向きを図で見る

```
【1-no-di】                      【2-di-concrete】               【3-di-interface】

application                     application                    application
    │                               │                              │
    │ importして                     │ importして                    │ importして
    │ 自分でnewする                   │ 引数で受け取る                  │ 引数で受け取る
    ▼                               ▼                              ▼
persistence                     persistence                     domain (interface)
    │                               │                              ▲
    ▼                               ▼                              │ 実装している
PostgreSQL                      PostgreSQL                     persistence
                                                                   │
                                                                   ▼
                                                               PostgreSQL
```

左2つは、矢印が`application`から`persistence`へ**直接**向いています。ビジネスロジックがPostgreSQL実装に依存している状態です。

右端だけ、`application`と`persistence`の間に`domain`が挟まり、しかも`persistence`からの矢印が**上向き**になっています。「persistenceがdomainのルールに合わせる」という形で、依存の向きが逆転しています。これが *Dependency Inversion*（依存性逆転）と呼ばれるものの実体です。

---

## 5. まとめ

| | Mockを注入できるか | 接続先DBを選べるか | キャッシュ追加時にapplication層を触ったか |
|---|---|---|---|
| 1-no-di | ✕ | ✕ | 触った |
| 2-di-concrete | ✕ | ○ | 触った |
| 3-di-interface | ○ | ○ | **触っていない** |

3つの差は、`repository`フィールドの型をどう書くかという、たった1行の選択から生まれています。

```go
repository *persistence.CachedBoxRepository   // 具体型: その構造体そのものしか入らない
repository domain.BoxRepository               // interface: Findを持つものなら何でも入る
```

「interfaceを使うと良い」という抽象的な話ではなく、**この1行の書き方次第で、後から機能を足すときに何ファイル書き換える羽目になるかが決まる**、というのがここで確認できることです。
