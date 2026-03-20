# Mermaid リッチ表示ライブラリ

Mermaid テキストを受け取り、ズーム・パン・テーマ切替などのインタラクティブ機能付きでリッチ表示するライブラリの実装方法をまとめる。

## ファイル構成

```
tech/mermaid/
├── README.md          # このドキュメント
├── mermaid-rich.js    # ライブラリ本体
└── index.html         # デモページ
```

## 実現したいこと

| 機能 | 概要 |
|------|------|
| リッチレンダリング | mermaid.js で SVG に変換してブラウザ表示 |
| ズーム | ホイール・ピンチ・ツールバーボタンで拡大縮小 |
| パン（位置調整） | ドラッグ・スワイプで表示領域を移動 |
| フィット表示 | 図全体をビューポートに収める |
| テーマ切替 | ライト / ダーク テーマ |
| SVG ダウンロード | 現在の図を SVG ファイルとして保存 |
| 自動初期化 | `data-mermaid-rich` 属性で宣言的に使用 |

## 技術選定

### Mermaid のレンダリング

[mermaid.js](https://mermaid.js.org/) を使用する。CDN から動的ロードするため、
利用側は `<script>` タグを自前で用意しなくてよい。

```
https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.min.js
```

`mermaid.render(id, code)` は Promise を返し、`{ svg }` を resolve する。
生成された SVG 文字列をコンテナに挿入することで図を表示する。

### ズーム / パン の実現方法

**CSS Transform** を使って SVG を包む `<div>` に `transform: translate() scale()` を
適用する方式を採用する。SVG の DOM 構造を変更しないため、どんな図種 (flowchart,
sequence, class, gantt, …) でも動作する。

```
viewport (overflow: hidden)
└── diagram (transform: translate(tx,ty) scale(s))
    └── <svg> (mermaid 出力)
```

**ズーム計算（マウス位置基点）**:

```js
// マウスカーソル位置 (cx, cy) を基点にズーム
tx_new = cx - (cx - tx) * (newScale / oldScale)
ty_new = cy - (cy - ty) * (newScale / oldScale)
```

これにより、カーソル下の点が画面上で動かずスケールが変わる直感的なズームになる。

### タッチ対応

- 1 本指: パン（`touchmove` で tx/ty を更新）
- 2 本指: ピンチズーム（指間距離の比率で scale を更新）

### フィット表示

```js
const scaleX = vpWidth  / svgWidth;
const scaleY = vpHeight / svgHeight;
scale = Math.min(scaleX, scaleY) * 0.9;  // 90% に収める
tx = (vpWidth  - svgWidth  * scale) / 2;
ty = (vpHeight - svgHeight * scale) / 2;
```

## 使い方

### 1. JavaScript API

```html
<div id="my-diagram" style="height:400px;"></div>
<script src="./mermaid-rich.js"></script>
<script>
  const mr = new MermaidRich(document.getElementById("my-diagram"), {
    theme: "light",   // "light" | "dark"
  });

  mr.render(`
    flowchart LR
      A --> B --> C
  `);
</script>
```

### 2. data 属性による自動初期化

```html
<!-- height は style で指定、テーマは data-theme で指定 -->
<div data-mermaid-rich data-theme="dark" style="height:300px;">
  sequenceDiagram
    Alice->>Bob: Hello
    Bob-->>Alice: Hi!
</div>

<script src="./mermaid-rich.js"></script>
<script>
  MermaidRich.autoInit();
</script>
```

## オプション一覧

| オプション | 型 | デフォルト | 説明 |
|---|---|---|---|
| `theme` | string | `"light"` | UI テーマ (`"light"` / `"dark"`) |
| `mermaidTheme` | string | theme に連動 | mermaid 内部テーマ (`"default"`, `"dark"`, `"forest"`, `"neutral"`) |
| `minZoom` | number | `0.1` | 最小ズーム倍率 |
| `maxZoom` | number | `10` | 最大ズーム倍率 |
| `zoomStep` | number | `0.1` | ホイール 1 ステップのズーム量 |
| `toolbar` | boolean | `true` | ツールバーを表示するか |

## ツールバーボタン

| ボタン | 機能 |
|--------|------|
| ＋ | ズームイン |
| － | ズームアウト |
| ⊙ | 100% に戻す |
| ⊞ | 全体をビューポートに収める |
| 🌗 | ライト / ダーク テーマ切替 |
| ↓ | 現在の SVG をファイルとしてダウンロード |

## 対応 Mermaid 図種

mermaid.js が対応している図種はすべて動作する。

- `flowchart` / `graph`
- `sequenceDiagram`
- `classDiagram`
- `gantt`
- `erDiagram`
- `journey`
- `pie`
- `stateDiagram-v2`
- `mindmap`
- `timeline`
- その他

## デモページ

`index.html` をブラウザで開くとデモを確認できる。
GitHub Pages にホストすれば URL で共有できる。

```bash
# ローカルで簡単に確認する例
cd tech/mermaid
python3 -m http.server 8080
# → http://localhost:8080 を開く
```

## 今後の拡張アイデア

- **ノードのドラッグ移動**: `flowchart` の SVG ノード要素を個別にドラッグ可能にする。
  mermaid が出力する SVG の `.node` 要素に `mousedown` イベントを仕掛け、
  `transform` を書き換える方式が有力だが、エッジ（矢印）の再描画が課題。
- **エディタ連携**: Monaco Editor や CodeMirror と組み合わせてライブプレビューにする。
- **PNG エクスポート**: `<canvas>` に SVG を描画して `toDataURL()` で PNG 化する。
- **共有 URL**: Mermaid コードを Base64 エンコードして URL パラメータに載せる。
