## DIとはなにか

DIとはコンストラクタみたいなもの。

振る舞い（メソッド）を持っていて、代わりに何かをやってくれるものを入れるとDIになる。

## DIの特徴

DIによる違いは、大きく分けて**「テスト」「責務と変更の波及範囲」**の2つの観点で整理できる。後者は「どの実装を使うかを決める権限がどの層にあるか（責務）」が原因で、「機能追加のたびにどこまで書き換える羽目になるか（波及範囲）」がその結果として現れる、という1つながりの話。

| 観点 | DIを使用しない | DIを使用する |
| --- | --- | --- |
| **1. テスト** | 実際のDBや外部サービスが必要になりやすい。
そのため、テストが遅い・外部環境に左右される・DBやテストデータの事前準備が必要・異常系を再現しづらい | Mockに差し替えてテストできる。
そのため、高速・安定・事前準備が少ない・エラーを自由に返して異常系を再現できる |
| **2. 責務と変更の波及範囲** | UseCaseが「どの実装を使うか」を自分で決めるため、DBやキャッシュなどインフラの都合がApplication層やテストにまで波及しやすい | UseCaseは「決める」役割を持たず、外から与えられるだけなので、実装を変更してもUseCaseやテストへの影響を抑えやすい |

## お題（箱のたとえ）

箱が6個あり、それぞれの箱には数字が入っている。

箱を押すと、その箱に入っている数字が返ってくる（`Find(id int) (*Box, error)`）。

その数字が4以上なら〇（`true`）、4未満なら×（`false`）と判定する（`IsLarge(id int) (bool, error)`）。

`id`は「箱の番号（何番の箱か）」であって、「箱の中の数字」ではない。たとえば箱1の中には5が入っている、というだけで、`1`と`5`の間に直接の関係はない。

### テストの面での比較

`IsLarge`が検証したいのは「`Number >= 4`かどうか」という、ただの算数です。本来なら一瞬で終わるはずのこの確認に、DIなし（`1-no-di`）だと本物のPostgreSQLへの接続が毎回必要になります。DIあり（`3-di-interface`）ならMockを注入するだけで、DBを一切使わずに検証できます。

この違いが生む具体的な困りごとは3つあります。

1. **遅くて不安定**: DB接続はメモリ上の計算より桁違いに遅く、ネットワークやDBの起動状況次第でテストが不安定になる（flakyになる）
2. **事前準備が要る**: 「id=1の箱には5が入っている」という事実を、DBに事前にINSERT（seed）しておかないとテストが意味をなさない。他のテストがそのデータを書き換えていないことまで保証する必要もある

たとえばテストコードには`useCase.IsLarge(1)`としか書かれていなくても、これが正しく動くには裏で、次のようなseedデータが別ファイル（01_schema.sql）に存在していないといけません。

```sql
INSERT INTO boxes (id, number) VALUES (1, 5);
```

もし誰かがこのseedデータを`(1, 3)`に変更したら、`IsLarge`自体のロジックは何も壊れていないのに、テストは黙って失敗するようになります。テストコードだけを読んでも「なぜ5という数字を期待しているのか」が分からず、常にDB側のファイルとセットで管理しないといけない、というのがこの「事前準備」の中身です。

1. **異常系が作れない**: 「DB接続が切れた場合」のようなエラーケースを試したくても、本物のDBを相手に意図的にエラーを起こすのは難しい

**DIを使用している場合、Mockを使用できます。Mockなら、これらが全部その場で作れます。**

```go
// 「5が返ってくる」状況を、DBの中身に関係なく1行ででっち上げる
mock := &MockBoxRepository{Box: &domain.Box{ID: 1, Number: 5}}

// 「エラーが返ってくる」状況も同様に作れる
mock := &MockBoxRepository{Err: errors.New("connection refused")}
```

DBの実データがどうなっているかを気にする必要はなく、「この状況で`IsLarge`は正しく動くか」だけをピンポイントで確認できます。

### インターフェースが必須となる

Mockがあると便利ということはわかったが、どうやったらMockを使えるようになるのか。

インターフェースの実装が必須になってくる。

#### なぜインターフェースが必須なのか

答えを一言でいうと、**`repository`フィールドの型を具体型で固定してしまうと、そこに入れられる実体は「その具体型そのもの」しか許されず、`Find`を呼んだ瞬間に「その具体型が持つFind」しか使いようがなくなるから**です。

Mockが同じ形の`Find`を持っていても、型そのものが違う以上は差し込む余地がありません（`2-di-concrete`で実際に見たコンパイルエラーがこれでした）。

```go
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
   → Find(id int) (*domain.Box, error) を持つので domain.BoxRepository を実装している
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
⑥ IsLargeは box.Number（=5）>= 4 を判定 → true
```

⑤が核心です。`u.repository.Find(1)`という**同じ1行のコード**が、渡された実体次第で「本物のDBに繋ぐFind」にも「Mockが即座に値を返すFind」にもなります。

これができるのは、`repository`フィールドがinterface型で「中身が何であってもいい」状態になっているからです。もし具体型で固定されていたら、この1行は「その具体型のFindしか呼べない1行」になってしまい、Mockを混ぜる余地はありません。

### 責務と変更の波及範囲の比較

先に正直なところを書いておきます。**この仕組みが実務で一番使われるのはテストのためで、本番の実装どうしを差し替える場面は稀です**。

interfaceの説明では「あとからPostgreSQLを別のDBに乗り換えられます」という売り文句をよく見かけますが、実際にそれが起きることはめったにありません。なので「将来差し替えられるから」を理由にinterfaceを切るのは、動機としては弱いです。

ここで見たいのは「差し替えられて便利」ではなく、**「どの実装を使うか」を決める権限がどの層にあるか**です。これを「責務」と呼びます。

#### 責務: 決定権がどこにあるか

`1-no-di`のコンストラクタを見てください。

```go
// application/box_usecase.go
import "..../1-no-di/persistence"   // ← application層がインフラのパッケージを知っている

func NewBoxUseCase() *BoxUseCase {
    inner := persistence.NewDefaultBoxRepository()   // ← 「どれを使うか」をapplication層が決めている
    return &BoxUseCase{repository: persistence.NewCachedBoxRepository(inner)}
}
```

`application`層が`persistence`パッケージを`import`し、`NewDefaultBoxRepository`という**具体的な関数名を名指し**しています。この関数が内部で環境変数を読もうが`sql.Open`しようが（実際、`sql.Open`自体は`persistence`層に置かれています）、「PostgreSQL版のBoxRepositoryを使う」という決定そのものは、ビジネスロジックを書く層が下しています。「DB接続のコードがどのファイルにあるか」ではなく、「どれを使うか決めているのは誰か」が責務の本体です。

`3-di-interface`はこうなっています。

```go
// application/box_usecase.go
import "..../3-di-interface/domain"   // ← persistenceを一切知らない

func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase {
    return &BoxUseCase{repository: repository}
}
```

```go
// main.go — 「どれを使うか」をここで初めて決める
repository := persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))
useCase := application.NewBoxUseCase(repository)
```

`application`層は「`Find`できる何か」としか言っておらず、`persistence`パッケージの存在すら知りません。「PostgreSQL版を使うか、キャッシュ付きにするか、Mockにするか」という決定は、すべて`main.go`（呼び出す側）に集約されています。

#### 波及範囲: 決定権の違いが、実際にどれだけの差を生むか

決定権が`application`層にあるとどうなるか。`1-no-di`で「テスト用のDBに繋ぎたい」となっても、その決定権はUseCaseの中にあるので、UseCase自身を書き換えないと接続先を変えられません。`3-di-interface`なら、決定権は呼び出す側にあるので、`main.go`（や`_test.go`）を書き換えるだけで済み、`BoxUseCase`は無傷のままです。

これを実際の機能追加で測ってみます。「キャッシュを追加したい」というお題を使います。素の`BoxRepository`をキャッシュで包む`CachedBoxRepository`を用意するところまでは3つとも共通で、違いは**それを実際に使わせるために、どこまで書き換える羽目になったか**です。

実際にgitのコミットとして分けてあります。いったん3つすべてからキャッシュを取り除いた状態を作り、そこから同じ機能をディレクトリごとに1コミットずつ追加し直しているので、**各コミットの差分がそのまま「決定権の違いが、同じ機能を足すときに何を書き換える羽目になるかの答え」**になっています。

#### 3-di-interface: application層は1文字も触っていない

**実際の差分**

```go
// main.go — この1行だけ
- repository := persistence.NewBoxRepository(db)
+ repository := persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))
```

```go
// application/box_usecase.go — 変更なし
type BoxUseCase struct {
    repository domain.BoxRepository
}

func NewBoxUseCase(repository domain.BoxRepository) *BoxUseCase {
    return &BoxUseCase{repository: repository}
}
```

テストコードも変更していません。`CachedBoxRepository`も`Find`を持っている＝`domain.BoxRepository`を満たしているので、`BoxUseCase`から見れば渡ってくるものは何も変わっていないからです。

コミットをGitHubで開くと、変更されたファイルの一覧に`application/box_usecase.go`がそもそも載っていません。新規ファイル2つ（キャッシュ実装とそのテスト）を除けば、既存コードへの変更は`main.go`の1行だけです。

#### 2-di-concrete: application層の型定義とテストを書き換えた

**実際の差分**

```go
// main.go
- repository := persistence.NewBoxRepository(db)
+ repository := persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db))
```

```go
// application/box_usecase.go — 型を2箇所書き換える必要があった
type BoxUseCase struct {
-    repository *persistence.BoxRepository
+    repository *persistence.CachedBoxRepository
}

- func NewBoxUseCase(repository *persistence.BoxRepository) *BoxUseCase {
+ func NewBoxUseCase(repository *persistence.CachedBoxRepository) *BoxUseCase {
    return &BoxUseCase{repository: repository}
}
```

```go
// application/box_usecase_test.go — 引数の型が変わったので呼び出し3箇所すべて
- useCase := NewBoxUseCase(persistence.NewBoxRepository(db))
+ useCase := NewBoxUseCase(persistence.NewCachedBoxRepository(persistence.NewBoxRepository(db)))
```

DIはしているのに、**ビジネスロジックの層とそのテストにまで変更が波及**しました。`*persistence.BoxRepository`という構造体を名指ししていたせいで、「別の構造体で包む」という変更がそのまま型の不一致になるからです。

#### 1-no-di: 型に加えてコンストラクタの中身まで書き換えた

**実際の差分**

```go
// application/box_usecase.go
type BoxUseCase struct {
-    repository *persistence.BoxRepository
+    repository *persistence.CachedBoxRepository
}

func NewBoxUseCase() *BoxUseCase {
-    return &BoxUseCase{repository: persistence.NewDefaultBoxRepository()}
+    inner := persistence.NewDefaultBoxRepository()
+    return &BoxUseCase{repository: persistence.NewCachedBoxRepository(inner)}
}
```

型定義だけでなく、**組み立ての手順そのもの**がUseCaseの中にあるので、組み立て方が変わればUseCaseが変わります。

#### 並べるとこうなる

|  | main.go | application層 | テスト | 実際の差分 |
| --- | --- | --- | --- | --- |
| 3-di-interface | 1行 | **変更なし** | **変更なし** | b5c6223 |
| 2-di-concrete | 1行 | 型を2箇所 | 3箇所 | 2eb75ec |
| 1-no-di | （そもそも無い） | 型1箇所＋処理2行 | **変更なし**（※） | ccde1f3 |

※ 1-no-diのテストが変更不要だったのは、3-di-interfaceとは理由が正反対です。3-di-interfaceは「interfaceが変更を吸収したから」ですが、1-no-diは「`NewBoxUseCase()`に引数が無く、そもそも渡すものが無いから」です。テスト自体は存在しますが、Mockを渡す余地も接続先を変える余地も無く、DBを起動していなければ`t.Skipf`で黙ってスキップされます（＝落ちもしないが、確認もされない）。

やりたかったことは3つとも「キャッシュを挟む」というまったく同じ1つの変更です。それなのに、`repository`フィールドを`domain.BoxRepository`と書いたか`*persistence.BoxRepository`と書いたかという**たった1行の違い**で、触るファイル数がここまで変わります。

#### なぜ「波及すること」が困るのか

核心は、**変わる頻度が違うものが、くっついていると困る**という点にあります。

|  | 変わる頻度 |
| --- | --- |
| 「`Number >= 4`なら〇」というビジネスルール | **めったに変わらない** |
| DB接続、キャッシュ、ライブラリ、リトライ処理 | **しょっちゅう変わる** |

インフラ側は変更の理由が次々に出てきます。「遅いからキャッシュを入れる」「たまに失敗するからリトライを足す」「ログを取りたい」「ライブラリのバージョンを上げる」。一方でビジネスルールの「4以上かどうか」は、仕様変更がない限り何年も変わりません。

問題は、**変わりやすい側の都合で、変わらないはずの側まで触らされること**です。`2-di-concrete`でキャッシュを入れたとき、ルールは何も変わっていないのに`box_usecase.go`（ビジネスルールが書いてあるファイル）を書き換えました。これが実務では次のような形で効いてきます。

1. **レビューとテストの範囲が広がる**: ビジネスロジックのファイルが変更されていると、レビューする人は「ロジックが変わったのか？」を確認しないといけません
2. **見積もりが外れる**: 「キャッシュを入れるだけなので30分」のつもりが、UseCaseとテスト3箇所を直してレビューもやり直しになる
3. **事故の確率が上がる**: 触る必要のなかったファイルを触れば、そこで壊す可能性が生まれる
4. **一番怖いのは「誰もやらなくなる」こと**: 面倒で怖い作業は後回しになり、性能問題が放置されたり、その場しのぎの雑な実装が入ったりする

なので正確には、依存先の**数**が問題なのではありません。**変わりにくいもの（ビジネスルール）が、変わりやすいもの（インフラ）に縛られている**という**向き**が問題です。

逆に「インフラ側がビジネスルールに合わせる」形になっていれば、インフラがいくら変わっても中心は無傷でいられます。

#### ただし、常に重要なわけではない

正直に言うと、これが効いてくるのは条件付きです。

- コードが長く使われる（数ヶ月で捨てるなら関係ない）
- 複数人が触る
- ビジネスロジックが複雑で、壊れると困る

自分ひとりで書き捨てるスクリプトなら、`1-no-di`で何の問題もありません。**「波及すると困る」のは、そのコードを長く触り続ける前提があるから**です。

冒頭に書いた通り、キャッシュを入れる機会は今後無いかもしれません。ですが「依存先の都合が変わる」こと自体は実務で何度も起きます。そのたびに変更が他の層へ波及するかどうかが、`repository`フィールドをどう書いたかという1行で決まっている、というのがここで確認できることです。

実際のコードは 1-no-di / 2-di-concrete / 3-di-interface にそれぞれ置いてあるので、見比べてみてください。

ここまでは「Mockが弾かれるか」「決定権がどこにあるか」「キャッシュ追加で何を触るか」のように、**論点ごとに**3つを比較してきました。

これとは別に、**3つのapplication層の全文を層ごとに並べた比較**を comparison.md に用意しています。

`IsLarge`の中身は3つとも1文字も違わず、違うのは`import`先と`repository`フィールドの型だけ、という事実から出発して、テストの書き方の差や依存の向きの図まで通しで追える構成にしてあります。