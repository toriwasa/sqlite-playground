# メモ
## ディレクトリ構成の考え方
- 参考: https://go.dev/doc/modules/layout

### 外部からimportするシンプルなモジュール
- ルートディレクトリに modname.go と modname_test.go を置く
- package modname を宣言する

### 単一の実行可能CLI/GUIツール
- ルートディレクトリに func main を含む main.go を置く(modname.go でもいい)
- 大きくなってきたらルートディレクトリに分割した各種モジュールを置く
- 全てのモジュールで package main を宣言する

### 外部からimportする複雑なモジュール
- ルートディレクトリに modname.go と modname_test.go を置く
- importする外部に公開したくないプライベートなAPIは internal/内部モジュール名 ディレクトリ配下に置く
- internal 配下のモジュールは package 内部モジュール名 で独立して定義する

### 単一の実行可能なCLI/GUIツール+複雑なモジュール
- ルートディレクトリに func main を含む main.go を置く(modname.go でもいい)
- importする外部に公開したくないプライベートなAPIを internal/内部モジュール名 ディレクトリ配下に置く
- internal 配下のモジュールは package 内部モジュール名 で独立して定義する

### 外部からimportする複数のモジュール
- ルートディレクトリに modname.go と modname_test.go を置く
- import させたいサブモジュールを サブモジュール名ディレクトリ 配下に置く

### 複数の実行可能なCLI/GUIツール
- 実行させたいツールごとにディレクトリを作成する
- ツール名/main.go を置く
- 複数のツールで共通のコードを使う場合は internal ディレクトリ配下に置く

## 開発の考え方
- CUIインターフェースは handler/cui/cui.go に実装する
- -> コマンドライン引数を処理してパラメータデータクラスのインスタンスを生成、コントローラー層のExecute()を呼び出す
- コントローラーは controller/controller.go に実装する
- -> Execute() は パラメータデータクラスのメソッドとして実装する
- -> usecase 配下のビジネスロジック部分の公開関数を順番に呼び出す
- usecase 配下のビジネスロジック部分は usecase/mylogic/mylogic.go に実装する
- -> 用途に応じてパッケージ名は変更する

