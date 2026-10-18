---
pagetitle: Arkadia Lang
---

::: hero
![](assets/logo.png){.hero-logo width=120 height=120}

# Arkadia Lang

[代数的エフェクトと型クラスを備えた、正格評価の静的型付き関数型プログラミング言語]{.tagline}

ネイティブバイナリにコンパイルされ、構造化された軽量タスクがすべてのプロセッサで並列に動きます。

::: actions
[はじめる](#インストール){.button .primary}
[言語ガイド](guide.html){.button}
[GitHub](https://github.com/arkadia-lang/arkadia){.button}
:::
:::

::: {.block .sample}
```haskell
import Std.State (..)
import Std.Task (parMap)

-- エフェクトは型に現れ、ハンドラ（runState）が処理する
counter : Unit -> <State Int, IO> Unit
counter () =
  let n = get ()
  println "count = ${n}"
  put (n + 1)

-- 制御は match・ガード・再帰だけ
fib : Int -> Int
fib n
  | n < 2 = n
  | otherwise = fib (n - 1) + fib (n - 2)

main () =
  let (_, final) = runState 0 (\_ -> counter (); counter ())
  println "final = ${final}"
  println "${parMap fib [25, 30, 35]}"   -- 各要素を別のタスクで並列に計算
```

::: note
Arkadia は 0.x の開発段階です。1.0 になるまでは、マイナーバージョンの更新で互換性のない変更が入ることがあります。
:::
:::

::: block
## 特徴

::: features
::: card
### 代数的エフェクト

例外・状態・ジェネレータなどをライブラリとして定義できます。エフェクトは行多相の型で追跡・推論され、処理されていないエフェクトはコンパイルエラーになります。
:::

::: card
### 構造化された並列実行

タスクは全プロセッサで並列に動き、キャンセル・タイムアウト・失敗はスコープに沿って伝播します。純粋な無限ループも外からキャンセルできます。
:::

::: card
### 表現力のある型システム

型クラス（高カインド・`deriving`）、レコード、GADT と型族によるインデックス型。網羅性はコンパイル時に検査されます。
:::

::: card
### 小さな構文

`if` もループも無く、制御は `match`・ガード・再帰だけ。インデントによるブロック、空白区切りの関数適用、全関数カリー化。
:::

::: card
### ネイティブバイナリ

Go を経由して単一の静的バイナリを生成し、Go ランタイムの GC とスケジューラの上で動きます。Go の関数も呼び出せます（FFI）。
:::

::: card
### ツール一式

ビルド・実行・テスト・フォーマッタ・ドキュメント・REPL・依存管理を `ark` コマンドひとつで。Language Server と VS Code 拡張もあります。
:::
:::
:::

::: block
## インストール

Linux と macOS（amd64・arm64）に対応しています。`ark` は Go のコードを生成して `go build` するので、**Go 1.22 以降**が必要です。

[Releases](https://github.com/arkadia-lang/arkadia/releases) から環境に合ったアーカイブを取得して展開し、`bin/` を `PATH` に加えます。

```sh
curl -LO https://github.com/arkadia-lang/arkadia/releases/download/v0.1.0/arkadia-0.1.0-linux-amd64.tar.gz
tar xzf arkadia-0.1.0-linux-amd64.tar.gz
export PATH="$PWD/arkadia-0.1.0-linux-amd64/bin:$PATH"
ark version
```

ソースからビルドするには Go と [just](https://github.com/casey/just) を使います。

```sh
git clone https://github.com/arkadia-lang/arkadia.git && cd arkadia
just build      # bin/ark0 (Go) → bin/ark1 → bin/ark
just install    # $GOBIN/ark から bin/ark へのシンボリックリンクを張る
```
:::

::: block
## はじめる

```sh
ark new hello      # プロジェクトを作成
cd hello
ark run            # ビルドして実行  → Hello, Arkadia!
ark test           # テストを実行
ark build          # bin/hello を生成
```

単一のファイルも直接実行できます（`ark run main.ark`）。

| コマンド | 説明 |
|---|---|
| `ark new <name>` | プロジェクトの雛形を作成 |
| `ark run [file] [-- args]` | ビルドして実行 |
| `ark build [file]` | ネイティブバイナリを生成 |
| `ark check [file]` | 型検査のみ |
| `ark test [file]` | テストを実行 |
| `ark fmt [path...]` | ソースを整形 |
| `ark doc <Module>` | モジュールの公開 API を表示 |
| `ark repl` | 対話環境 |
| `ark lsp` | Language Server |
:::

::: block
## さらに読む

::: features
::: card
### [言語ガイド](guide.html)

構文・型・型クラス・エフェクト・タスク・テスト・FFI のリファレンス。
:::

::: card
### [標準ライブラリ](std/index.html)

`Prelude` と `Std.*` の各モジュールの API とドキュメント。
:::

::: card
### [サンプル](https://github.com/arkadia-lang/arkadia/tree/main/examples)

並列タスク、エフェクトによる協調スレッド、JSON パーサ、Lisp インタプリタ、HTTP サーバなど。
:::
:::
:::
