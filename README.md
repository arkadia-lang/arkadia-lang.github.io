# arkadia-lang.github.io

The website of [Arkadia](https://github.com/arkadia-lang/arkadia), published at https://arkadia-lang.github.io/.

| ファイル | 内容 |
|---|---|
| `index.md` | トップページ |
| `template.html` | 全ページ共通の pandoc テンプレート |
| `style.css` | スタイル（ライト／ダークモード） |
| `guide.lua` | 言語ガイドを変換する pandoc の Lua フィルタ |
| `stdref/` | `ark doc` の出力から標準ライブラリのリファレンスを生成するプログラム |

言語ガイド・変更履歴・標準ライブラリのリファレンスは [arkadia-lang/arkadia](https://github.com/arkadia-lang/arkadia) の内容から生成します（`docs/guide.md`、`CHANGELOG.md`、`lib/`）。それらを直すときは arkadia 側を変更してください。

## ビルド

arkadia を隣に clone して、`nix develop` の中で `just site` を実行すると `_site/` に出力されます（Go・just・pandoc が必要）。

```sh
git clone https://github.com/arkadia-lang/arkadia.git ../arkadia
nix develop --command just site            # ../arkadia から _site/ を生成
nix develop --command just arkadia=path site   # 別の場所の arkadia を使う
```

## 公開

`.github/workflows/pages.yml` が main への push でビルドして GitHub Pages に公開します。arkadia 側の変更は 1 時間ごとの定期実行で取り込みます（arkadia の main が前回の公開から進んでいなければ何もしません）。すぐ反映したいときは Actions から手動で実行してください。
