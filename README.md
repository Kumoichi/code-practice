# code-practice

このRepositoryは、技術を「読む」のではなく「コードを書いて理解する」ためのRepositoryです。

完成されたアプリケーションを作ることは目的にしていません。代わりに、

- 実際にコードを書く
- Good / Badの実装を比較する
- テストを書く
- コードジャンプで依存関係を追う
- 実際に変更を加えてみる
- READMEに自分の理解を整理する

という方法で、一つずつ技術テーマを理解していく場所です。

## 構成

技術テーマごとにディレクトリを分けています。今回は以下のテーマのみ実装済みです。

```text
code-practice/
├── go/
│   └── clean-architecture/
│       └── repository-interface/   ← 今回のテーマ
│
└── README.md
```

今後、`go/concurrency/` や `sql/transaction/` のような形でテーマを追加していく想定です。

## 今回のテーマ

[go/clean-architecture/repository-interface/](go/clean-architecture/repository-interface/) — UseCaseからInfrastructureの具象Repositoryを直接参照する場合と、Repository interfaceを挟む場合の違い。

詳細はそのディレクトリ内のREADME.mdを参照してください。
