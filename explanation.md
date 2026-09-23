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

**テストが書きづらい**

`IsLarge`が検証したいのは「`Number >= 4`かどうか」という、ただの算数です。本来なら一瞬で終わるはずのこの確認に、DIなし(`1-no-di`)だと本物のPostgreSQLへの接続が毎回必要になります。DIあり(`3-di-interface`)ならMockを注入するだけで、DBを一切使わずに検証できます。

この違いが生む具体的な困りごとは3つあります。

1. **遅くて不安定**: DB接続はメモリ上の計算より桁違いに遅く、ネットワークやDBの起動状況次第でテストが不安定になる（flakyになる）
2. **事前準備が要る**: 「id=1の箱には5が入っている」という事実を、DBに事前にINSERT(seed)しておかないとテストが意味をなさない。他のテストがそのデータを書き換えていないことまで保証する必要もある

   たとえばテストコードには`useCase.IsLarge(1)`としか書かれていなくても、これが正しく動くには裏で

   ```sql
   INSERT INTO boxes (id, number) VALUES (1, 5);
   ```

   のようなseedデータが別ファイル([01_schema.sql](docker/init/01_schema.sql))に存在していないといけません。もし誰かがこのseedデータを`(1, 3)`に変更したら、`IsLarge`自体のロジックは何も壊れていないのに、テストは黙って失敗するようになります。テストコードだけを読んでも「なぜ5という数字を期待しているのか」が分からず、常にDB側のファイルとセットで管理しないといけない、というのがこの「事前準備」の中身です。
3. **異常系が作れない**: 「DB接続が切れた場合」のようなエラーケースを試したくても、本物のDBを相手に意図的にエラーを起こすのは難しい

**Mockなら、これらが全部その場で作れます**

```go
// 「5が返ってくる」状況を、DBの中身に関係なく1行ででっち上げる
mock := &MockBoxRepository{Box: &domain.Box{ID: 1, Number: 5}}

// 「エラーが返ってくる」状況も同様に作れる
mock := &MockBoxRepository{Err: errors.New("connection refused")}
```

DBの実データがどうなっているかを気にする必要はなく、「この状況で`IsLarge`は正しく動くか」だけをピンポイントで確認できます。


## Mockがあると便利ということはわかったが、どうやったらMockを使えるようになるのか

インターフェースの実装が必須になってくる。

### なぜインターフェースが必須なのか

答えを一言でいうと: **`repository`フィールドの型を具体型で固定してしまうと、そこに入れられる実体は「その具体型そのもの」しか許されず、`Find`を呼んだ瞬間に「その具体型が持つFind」しか使いようがなくなるから**です。Mockが同じ形の`Find`を持っていても、型そのものが違う以上は差し込む余地がありません（`2-di-concrete`で実際に見たコンパイルエラーがこれでした）。

こちらのコード
```
type MockBoxRepository struct {
	Box *persistence.Box
	Err error
}

func (m *MockBoxRepository) Find(id int) (*persistence.Box, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Box, nil
}

func tryInject() {
	mock := &MockBoxRepository{Box: &persistence.Box{ID: 1, Number: 5}}
	_ = NewBoxUseCase(mock)
}
```

なので結果として`2-di-concrete/application/box_usecase_test.go`は`openTestDB(t)`で実際にPostgreSQLへ接続してからでないとテストできなくなる。「Applicationのロジック」をテストしたいだけなのに、DBが起動していないとテストができなくなる。


`repository`フィールドの型をinterfaceにすると、話が変わります。フィールドが要求するのは「`Find`という形のメソッドを持っていること」だけになるので、**そこに何の具体型を入れるかを、呼び出す側が選べるようになります**。中身が本物のDBでもMockでも、同じ`u.repository.Find(id)`という書き方で、それぞれの実体が持つ`Find`が呼ばれます。

実際にMockを注入したときの流れを追うと、こうなります。

```
① mock := &MockBoxRepository{...} を生成
   → 型は *MockBoxRepository（具体型）
   → Find(id int)(*domain.Box, error) を持つので domain.BoxRepository を実装している
↓
② NewBoxUseCase(mock) が呼ばれる
   引数の型は domain.BoxRepository（ラベルがinterfaceに切り替わる）
   実体は *MockBoxRepository のまま
↓
③ BoxUseCase{repository: mock} が生成される
   repositoryフィールド: ラベル=domain.BoxRepository、実体=*MockBoxRepository
↓
④ useCase.IsLarge(1) を呼ぶ
↓
⑤ u.repository.Find(1) が実行される
   → u.repositoryの実体は*MockBoxRepositoryなので
   → 実際に動くのは MockBoxRepository.Find
   → m.Errがnilなので m.Box（{ID:1, Number:5}）をそのまま返す
↓
⑥ IsLargeは box.Number(=5) >= 4 を判定 → true
```

⑤が核心です。`u.repository.Find(1)`という**同じ1行のコード**が、渡された実体次第で「本物のDBに繋ぐFind」にも「Mockが即座に値を返すFind」にもなります。これができるのは、`repository`フィールドがinterface型で「中身が何であってもいい」状態になっているからです。もし具体型で固定されていたら、この1行は「その具体型のFindしか呼べない1行」になってしまい、Mockを混ぜる余地はありません。


interfaceを使うと、ここまでの「Mockが使えるかどうか」以外にも2つメリットが出てきます。実際のコードで見てみます。

### メリット1: 差し替え可能性

`3-di-interface/application/box_usecase.go`は`repository`を外から受け取るだけなので、次のどれを渡してもそのまま動きます。

```go
NewBoxUseCase(persistence.NewBoxRepository(db))                                   // 本番用: 本物のDB
NewBoxUseCase(&MockBoxRepository{...})                                            // テスト用: Mock
NewBoxUseCase(persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))) // キャッシュ付き: 本物のDB + キャッシュ
```

3つ目は架空の例ではなく、実際にこのリポジトリで起きたことです。あとから「キャッシュを追加したい」となったとき、変更したのは`persistence`パッケージに`CachedBoxRepository`を1つ足したことと、`main.go`で渡す値を差し替えたことだけでした。

```go
// main.go: 変更前
repository := persistence.NewBoxRepository(db)
// main.go: 変更後
repository := persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))
```

`3-di-interface/application/box_usecase.go`と`3-di-interface/domain/box_repository.go`は、この変更で**1行も書き換えていません**。`BoxUseCase`は「`Find`できる何か」を受け取るとしか約束していないので、渡ってくる実体が「本物のDB」から「キャッシュ付きの本物のDB」に変わっても、何も気づかず動き続けます。これが「差し替え可能性」の実例です。

### メリット2: 責務の分離

`1-no-di`版のコンストラクタを見てみます。

```go
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

`IsLarge`がやりたいのは「`Number >= 4`かどうか」というビジネスロジックだけのはずです。ところが`NewBoxUseCase`の中には、`sql.Open`や接続文字列(DSN)の組み立てといった、**PostgreSQLというインフラの都合**が入り込んでしまっています。`application`パッケージを読んだ人は、本来知らなくていいはずの「どんなDBに、どんな接続文字列で繋いでいるか」まで目にすることになります。

`3-di-interface`版はこの境界が保たれています。

```go
func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase {
    return &BoxUseCase{repository: repository}
}
```

`sql.Open`もDSNも、`application`パッケージのどこにも出てきません。DBに関する知識は全て`persistence`パッケージ側（`main.go`で組み立てる側）に押し込まれていて、`application`層は「`Find`できる何か」を受け取るだけです。これが「責務の分離」で、ビジネスロジックを書く層とインフラに繋ぐ層が、コード上ではっきり切り離されている状態を指します。

---

## 3つの実装を全部並べて見る

ここまでは断片ごとに比較してきましたが、**3つの実装の全コードを層ごとに並べた比較**を [comparison.md](go/clean-architecture/repository-interface/comparison.md) に用意しました。

特に「同じ『キャッシュを追加したい』という要求に対して、3つそれぞれが何ファイル・何行書き換える羽目になったか」の比較は、interfaceの有無が実際の作業量にどう効くかが一番はっきり出る部分です。

| | main.go | application層 | テスト | 新規ファイル |
|---|---|---|---|---|
| 3-di-interface | 1行 | **変更なし** | **変更なし** | 1つ |
| 2-di-concrete | 1行 | 型を2箇所 | 3箇所 | 1つ |
| 1-no-di | （なし） | 型1箇所＋処理2行 | （元々テスト不能） | （同ファイル内に追加） |
