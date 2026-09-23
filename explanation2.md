## DIとはなにか

DIとはコンストラクタみたいなもの。

振る舞い（メソッド）を持っていて、代わりに何かをやってくれるものを、自分で作らずに外から受け取ると、それがDI(Dependency Injection)になる。

## お題（箱のたとえ）

箱が6個あり、それぞれの箱には数字が入っている。

箱を押すと、その箱に入っている数字が返ってくる（`Find(id int) (*Box, error)`）。

その数字が4以上なら〇（`true`）、4未満なら×（`false`）と判定する（`IsLarge(id int) (bool, error)`）。

`id`は「箱の番号（何番の箱か）」であって、「箱の中の数字」ではない。たとえば箱1の中には5が入っている、というだけで、`1`と`5`の間に直接の関係はない。

## DIを使わないとどうなるか

`1-no-di`版のコンストラクタは、こうなっている。

ここで実際のコードを見せながらrepositoryのコードを直接つないでいることを見せる
```go
func NewBoxUseCase() *BoxUseCase {
	inner := persistence.NewDefaultBoxRepository() // ← ここで名指しで生成
	return &BoxUseCase{repository: persistence.NewCachedBoxRepository(inner)}
}
```


`()`の中が空。呼び出す側は`NewBoxUseCase()`と書く以外の選択肢がなく、中で何を使うかは関数の中に決め打ちされている。**「何を渡すか選べる窓口」自体が存在しない**、というのがDIをしていない状態。

この状態でテストを書こうとすると、困りごとが3つ出てくる。

1. **遅くて不安定**: DB接続はメモリ上の計算より桁違いに遅く、ネットワークやDBの起動状況次第でテストが不安定になる（flakyになる）
2. **事前準備が要る**: 「id=1の箱には5が入っている」という事実を、DBに事前にINSERT(seed)しておかないとテストが意味をなさない。たとえばテストコードには`useCase.IsLarge(1)`としか書かれていなくても、これが正しく動くには裏で`INSERT INTO boxes (id, number) VALUES (1, 5);`のようなseedデータが別ファイルに存在していないといけない。もし誰かがこのseedデータを`(1, 3)`に変更したら、`IsLarge`自体のロジックは何も壊れていないのに、テストは黙って失敗する
3. **異常系が作れない**: 「DB接続が切れた場合」のようなエラーケースを試したくても、本物のDBを相手に意図的にエラーを起こすのは難しい

## DIを使うとどうなるか

`3-di-interface`版と`2-di-concrete`版は、どちらも引数を受け取る。

```go
// good
func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase

// bad
func NewBoxUseCase(repository *persistence.BoxRepository) *BoxUseCase
```

どちらも`repository`という引数があるので、呼び出す側が

```go
NewBoxUseCase(realRepo) // 本物を渡す
NewBoxUseCase(mock)     // Mockを渡す
```

のように「何を渡すか」を選べる。nodiと違って「窓口」は存在する。

ところが、実際に試してみると`3-di-interface`はMockを渡せるのに、`2-di-concrete`は同じことをするとコンパイルエラーになる。**窓口があることと、Mockが実際に通ることは別問題**、というのがここで見えてくる。

## Mockを使うためにはどんな状態である必要があるか

`2-di-concrete`側で`3-di-interface`と同じようにMockを自作して渡そうとすると、次のエラーになる。

```
cannot use mock (variable of type *MockBoxRepository)
as *persistence.BoxRepository value in argument to NewBoxUseCase
```

`MockBoxRepository`は`Find(id int) (*Box, error)`という、本物と全く同じ形のメソッドを持っているのに拒否される。つまり「Mockが使えるかどうか」を決めているのは、Mockの中身の出来ではなく、**引数の型として何を要求しているか**、ということになる。

## なぜ具体型に入れたら問題が起こるのか

`2-di-concrete`の引数の型はこうなっている。

```go
func NewBoxUseCase(repository *persistence.BoxRepository) *BoxUseCase
```

`*persistence.BoxRepository`は、`persistence`パッケージの中で定義された**具体的な1つの構造体**を指している。Goはここで「同じメソッドを持っているか」ではなく、「本当にその構造体そのものか」で判定する。`MockBoxRepository`は見た目(メソッドの形)が同じでも、Goから見れば全くの別の型なので、素通しできない。

たとえるなら、「山田さんという名前の人しか通しません」という受付があるとして、山田さんと全く同じ服装・全く同じ持ち物の人が来ても、名前が「山田」でなければ通さない、というようなもの。中身や振る舞いが似ているかどうかは一切見ていない。

## そもそもインターフェースはGoではどうやって定義されるのか（暗黙的実装の話）

`3-di-interface`版では、`repository`の型が具体的な構造体ではなく、こうなっている。

```go
type BoxRepository interface {
	Find(id int) (*Box, error)
}
```

これがinterface。「`Find(id int) (*Box, error)`という形のメソッドを持っていること」だけを条件にした、**型の集合の名前**のようなもの。

Goのinterfaceには大きな特徴がある。「このinterfaceを実装します」と宣言する必要が一切ない(暗黙的実装)。他の言語(Javaなど)でよくある`class Mock implements BoxRepository`のような宣言はGoには存在しない。

```go
type MockBoxRepository struct { ... }

func (m *MockBoxRepository) Find(id int) (*Box, error) { ... }
```

これだけ書けば、`MockBoxRepository`は「宣言なし」で自動的に`BoxRepository`interfaceを満たしたことになる。判定基準は名前ではなく、**その型が要求された形のメソッドを持っているかどうか**、それだけ。

## インターフェースがあればMockの差し込みがなぜ簡単にできるのか

さきほどの「山田さんしか通さない受付」のたとえで言うと、interfaceを使うのは受付のルールを「名前が山田であること」から「入館証を持っていること」に変えるようなもの。入館証(＝`Find(id int) (*Box, error)`というメソッド)さえ持っていれば、本物の社員(`persistence.BoxRepository`)でも、テスト用の来訪者(`MockBoxRepository`)でも、同じように通してもらえる。

```go
func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase
```

この引数が要求しているのは「`Find`を持っていること」だけなので、

- 本物: `persistence.BoxRepository`（PostgreSQLに繋ぐ）
- テスト用: `MockBoxRepository`（メモリ上で好きな値を返す）

のどちらも、**同じ条件を満たしてさえいれば**そのまま渡せる。具体型だった`2-di-concrete`では「本物の型そのものか」まで問われていたので、この余地が無かった。

これが、`3-di-interface`だけがMockを注入できて、DBなしで一瞬・安定してテストできる理由の全体像。
