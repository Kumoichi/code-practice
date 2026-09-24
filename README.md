# code-practice

このRepositoryは、技術を「読む」のではなく「コードを書いて理解する」ためのRepositoryです。

完成されたアプリケーションを作ることは目的にしていません。代わりに、

- 実際にコードを書く
- 良い実装と悪い実装を並べて比較する
- テストを書く
- コードジャンプで依存関係を追う
- 実際に変更を加えてみる
- 自分の理解を文章に整理する

という方法で、一つずつ技術テーマを理解していく場所です。

## 構成

技術テーマごとにディレクトリを分けています。今回は以下のテーマのみ実装済みです。

```text
code-practice/
├── explanation.md                   ← 今回のテーマの解説（まずここ）
├── docker-compose.yml               ← 動作確認用のPostgreSQL
├── docker/init/01_schema.sql        ← テスト用のテーブルとseedデータ
└── go/
    └── clean-architecture/
        └── repository-interface/    ← 今回のテーマのコード
            ├── README.md            ← コードの歩き方と演習
            ├── comparison.md        ← 3実装の全コード比較
            ├── 1-no-di/             ← DIをしない
            ├── 2-di-concrete/       ← DIするが具体型で受け取る
            ├── 3-di-interface/      ← DIしてinterfaceで受け取る
            └── 4-di-wire/           ← 3-di-interfaceの組み立てをWireで生成する
```

今後、`go/concurrency/` や `sql/transaction/` のような形でテーマを追加していく想定です。

## 今回のテーマ

**DIとRepository interface** — UseCaseからInfrastructureの具象Repositoryを直接参照する場合と、Repository interfaceを挟む場合の違い。

まったく同じ機能を3通りの設計で実装し、「テストの書きやすさ」と「機能追加時に書き換える範囲」がどう変わるかを比較しています。

読む順番：

1. [explanation.md](explanation.md) — なぜDIとinterfaceが要るのかの解説。まずここから
2. [go/clean-architecture/repository-interface/README.md](go/clean-architecture/repository-interface/README.md) — 実行方法と、手を動かす演習
3. [go/clean-architecture/repository-interface/comparison.md](go/clean-architecture/repository-interface/comparison.md) — 3つの実装の全コードを並べた比較
