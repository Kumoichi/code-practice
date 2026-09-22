## DIとはなにか

DIとはコンストラクタみたいなもの

振る舞い（メソッド）を持っていて、代わりに何かをやってくれるものを入れるとDIになる

## お題（箱のたとえ）

箱が6個あり、それぞれの箱には数字が入っている。

箱を押すと、その箱に入っている数字が返ってくる（`Find(id int) (*Box, error)`）。

その数字が4以上なら〇（`true`）、4未満なら×（`false`）と判定する（`IsLarge(id int) (bool, error)`）。

`id`は「箱の番号（何番の箱か）」であって、「箱の中の数字」ではない。たとえば箱1の中には5が入っている、というだけで、`1`と`5`の間に直接の関係はない。

## DIの特徴

### DIを使用しないと何が起こるか

**テストが書きずらい**

1. **テストしやすさ**: `good/application/box_usecase_test.go`はMockを注入するだけでPostgreSQLなしにロジックを検証できました。DIがない`nodi`は、テストのたびに本物のDBが必須でした。

DIありではMockを使うことができる

DIなしではユースケース層であれ直接DBにつないでテストが行われる

- 返ってくる期待結果をエラーにしたいときなどに、操作できない
- ユースケースのテストなのに、DBに接続をするコストの高い処理を行う必要がある

`IsLarge`が検証したいのは「`Number >= 4`かどうか」というただの算数です。この算数のロジックを確認するのに、毎回本物のPostgreSQLへのTCP接続が必要になる、というのが不釣り合いなコストなんです。

### DB版とMock版で何が違うのか

「箱1（例えば5が入っている）を選んだら、4以上と認識されて〇となる」というテストを書きたいとき:

- **DB版**: 「箱1の中に実際に5が入っている」という事実を、事前にDB側へ用意（seed）しておかなければならない。テストが「ロジックが正しいか」ではなく「DBの中身がテストの期待と一致しているか」まで検証してしまう。
- **Mock版**: 「箱を押したら5が返ってくる」という状況を、テストコードの中でその場で作れる。DBの実データが何であろうと関係ない。

```go
// Mock: 箱の中身がなんであれ、「5が返ってきたケース」を1行で作れる
mock := &MockBoxRepository{Box: &domain.Box{ID: 1, Number: 5}}
```

DBだと「id=1の箱には5が入っている」という事実をテストコードとDBの両方が知っていないといけない上に、他のテストがデータを書き換えていないことまで保証しないといけません。

## DIをしてMockを使うためにはインターフェースが必要になってくる

### Goのinterfaceはテストでどのように役立っているか

GoのInterfaceの特徴：暗黙的実装

Goのinterfaceは宣言不要で、必要なメソッド（ここでは`Find(id int) (*domain.Box, error)`）を持ってさえいれば、それだけで`domain.BoxRepository`を満たしたことになる。

`MockBoxRepository`も`persistence.BoxRepository`(box_repository.go)も、共に「同じ形のFindを持っている」というだけで、このinterfaceの実装者になれる。

interfaceがなかったらどうなるか:

`bad/application/box_usecase.go`ではrepositoryフィールドが`*persistence.BoxRepository`という具体型に固定されているため、Mockを差し込む余地がない。

結果として`bad/application/box_usecase_test.go`は`openTestDB(t)`で実際にPostgreSQLへ接続してからでないとテストできなくなる。「Applicationのロジック」をテストしたいだけなのに、DBが起動していないとテストができなくなる。

1. **差し替え可能性**: `good/application/box_usecase.go`は`repository`を外から受け取るだけなので、本番用(`persistence.BoxRepository`)・テスト用(`MockBoxRepository`)・将来のインメモリ実装など、コードを変更せずに中身を入れ替えられます。

1. **責務の分離**: ビジネスロジック(`IsLarge`)を書く`application`層が、`sql.Open`やDSNの組み立てといったインフラの都合を知らずに済みます。`nodi`ではこの境界が壊れ、DBの都合が`application`パッケージにまで漏れ出していました。

DIを使わない書き方

```go
type BoxUseCase struct {
    repository *persistence.BoxRepository
}

func NewBoxUseCase() *BoxUseCase {
    db, _ := sql.Open("postgres", "postgres://...")       // ← 自分でDB接続を作る
    repo := persistence.NewBoxRepository(db)               // ← 自分でrepositoryを作る
    return &BoxUseCase{repository: repo}
}

// IsLarge reports whether the number in the box identified by id is 4 or higher.
func (u *BoxUseCase) IsLarge(id int) (bool, error) {
    box, err := u.repository.Find(id)
    if err != nil {
        return false, err
    }

    return box.Number >= 4, nil
}
```

interfaceをなぜ使用する必要があるのか

暗黙的実装:

Goのinterfaceは宣言不要で、必要なメソッド（ここでは`Find(id int) (*domain.Box, error)`）を持ってさえいれば、それだけで`domain.BoxRepository`を満たしたことになる。

`MockBoxRepository`も`persistence.BoxRepository`(box_repository.go)も、共に「同じ形のFindを持っている」というだけで、このinterfaceの実装者になれる。

interfaceがなかったらどうなるか:

`bad/application/box_usecase.go`ではrepositoryフィールドが`*persistence.BoxRepository`という具体型に固定されているため、Mockを差し込む余地がない。

結果として`bad/application/box_usecase_test.go`は`openTestDB(t)`で実際にPostgreSQLへ接続してからでないとテストできなくなる。「Applicationのロジック」をテストしたいだけなのに、DBが起動していないとテストができなくなる。
