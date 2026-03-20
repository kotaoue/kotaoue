/**
 * mermaid-rich.js
 * Mermaid をリッチな表示にするライブラリ
 *
 * 機能:
 * - Mermaid テキストを SVG に変換してリッチ表示
 * - マウスホイール / ピンチでズーム
 * - ドラッグでパン（位置移動）
 * - ズームリセット・フィットボタン付きツールバー
 * - ライト / ダークテーマ切替
 * - SVG ダウンロード
 */

(function (global) {
  "use strict";

  // ─── スタイル ────────────────────────────────────────────────
  const STYLES = `
.mr-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 300px;
  background: var(--mr-bg, #ffffff);
  border: 1px solid var(--mr-border, #e0e0e0);
  border-radius: 8px;
  overflow: hidden;
  font-family: sans-serif;
}
.mr-wrapper.dark {
  --mr-bg: #1e1e2e;
  --mr-border: #45475a;
  --mr-toolbar-bg: #313244;
  --mr-toolbar-border: #45475a;
  --mr-btn-bg: #45475a;
  --mr-btn-hover: #585b70;
  --mr-btn-color: #cdd6f4;
  --mr-text: #cdd6f4;
}
.mr-wrapper.light {
  --mr-bg: #ffffff;
  --mr-border: #e0e0e0;
  --mr-toolbar-bg: #f5f5f5;
  --mr-toolbar-border: #e0e0e0;
  --mr-btn-bg: #e8e8e8;
  --mr-btn-hover: #d0d0d0;
  --mr-btn-color: #333333;
  --mr-text: #333333;
}
.mr-toolbar {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  display: flex;
  gap: 4px;
  background: var(--mr-toolbar-bg, #f5f5f5);
  border: 1px solid var(--mr-toolbar-border, #e0e0e0);
  border-radius: 6px;
  padding: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12);
}
.mr-btn {
  background: var(--mr-btn-bg, #e8e8e8);
  border: none;
  border-radius: 4px;
  color: var(--mr-btn-color, #333);
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
  padding: 5px 8px;
  transition: background 0.15s;
  user-select: none;
}
.mr-btn:hover {
  background: var(--mr-btn-hover, #d0d0d0);
}
.mr-btn:active {
  opacity: 0.75;
}
.mr-viewport {
  width: 100%;
  height: 100%;
  cursor: grab;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}
.mr-viewport.grabbing {
  cursor: grabbing;
}
.mr-diagram {
  display: inline-block;
  transform-origin: 0 0;
  will-change: transform;
  transition: none;
}
.mr-diagram svg {
  max-width: none !important;
  display: block;
}
.mr-zoom-label {
  position: absolute;
  bottom: 8px;
  left: 10px;
  font-size: 11px;
  color: var(--mr-text, #888);
  pointer-events: none;
  opacity: 0.7;
}
`;

  // ─── ユーティリティ ──────────────────────────────────────────

  let _styleInjected = false;
  function injectStyles() {
    if (_styleInjected) return;
    const style = document.createElement("style");
    style.textContent = STYLES;
    document.head.appendChild(style);
    _styleInjected = true;
  }

  let _mermaidReady = false;
  let _mermaidCallbacks = [];

  /**
   * Mermaid CDN を動的ロードし、初期化する
   * @param {Function} cb  ロード完了時のコールバック
   */
  function loadMermaid(cb) {
    if (_mermaidReady) { cb(); return; }
    _mermaidCallbacks.push(cb);
    if (_mermaidCallbacks.length > 1) return; // 既にロード中

    if (typeof window.mermaid !== "undefined") {
      _onMermaidLoaded();
      return;
    }

    const script = document.createElement("script");
    script.src = "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.min.js";
    script.onload = _onMermaidLoaded;
    script.onerror = function () {
      console.error("[mermaid-rich] Mermaid CDN のロードに失敗しました。");
    };
    document.head.appendChild(script);
  }

  function _onMermaidLoaded() {
    _mermaidReady = true;
    _mermaidCallbacks.forEach(function (fn) { fn(); });
    _mermaidCallbacks = [];
  }

  // ─── MermaidRich クラス ──────────────────────────────────────

  /**
   * @param {HTMLElement|string} container  描画先の要素または CSS セレクタ
   * @param {Object}             options    オプション
   * @param {string}  [options.theme="light"]   "light" | "dark" | "default" | "forest" | "neutral"
   * @param {string}  [options.mermaidTheme]    mermaid 内部テーマ（デフォルトは theme に連動）
   * @param {number}  [options.minZoom=0.1]     最小ズーム倍率
   * @param {number}  [options.maxZoom=10]      最大ズーム倍率
   * @param {number}  [options.zoomStep=0.1]    ホイール 1 ステップのズーム量
   * @param {boolean} [options.toolbar=true]    ツールバーを表示するか
   */
  function MermaidRich(container, options) {
    if (typeof container === "string") {
      container = document.querySelector(container);
    }
    if (!container) {
      throw new Error("[mermaid-rich] container が見つかりません。");
    }

    this._container = container;
    this._opts = Object.assign({
      theme: "light",
      minZoom: 0.1,
      maxZoom: 10,
      zoomStep: 0.1,
      toolbar: true,
    }, options);

    this._scale = 1;
    this._tx = 0;
    this._ty = 0;
    this._dragging = false;
    this._dragStart = null;

    injectStyles();
    this._buildDOM();
  }

  // ── DOM 構築 ────────────────────────────────────────────────

  MermaidRich.prototype._buildDOM = function () {
    const wrapper = document.createElement("div");
    wrapper.className = "mr-wrapper " + (this._opts.theme === "dark" ? "dark" : "light");
    this._wrapper = wrapper;

    const viewport = document.createElement("div");
    viewport.className = "mr-viewport";
    this._viewport = viewport;

    const diagram = document.createElement("div");
    diagram.className = "mr-diagram";
    this._diagram = diagram;

    viewport.appendChild(diagram);
    wrapper.appendChild(viewport);

    if (this._opts.toolbar) {
      wrapper.appendChild(this._buildToolbar());
    }

    const zoomLabel = document.createElement("div");
    zoomLabel.className = "mr-zoom-label";
    zoomLabel.textContent = "100%";
    this._zoomLabel = zoomLabel;
    wrapper.appendChild(zoomLabel);

    this._container.appendChild(wrapper);
    this._bindEvents();
  };

  MermaidRich.prototype._buildToolbar = function () {
    const toolbar = document.createElement("div");
    toolbar.className = "mr-toolbar";

    const buttons = [
      { label: "＋", title: "ズームイン", onClick: () => this._zoom(this._opts.zoomStep) },
      { label: "－", title: "ズームアウト", onClick: () => this._zoom(-this._opts.zoomStep) },
      { label: "⊙", title: "実際のサイズ (100%)", onClick: () => this._resetZoom() },
      { label: "⊞", title: "全体を表示", onClick: () => this._fitToView() },
      { label: "🌗", title: "テーマ切替", onClick: () => this._toggleTheme() },
      { label: "↓", title: "SVG をダウンロード", onClick: () => this._downloadSVG() },
    ];

    buttons.forEach(function (cfg) {
      const btn = document.createElement("button");
      btn.className = "mr-btn";
      btn.textContent = cfg.label;
      btn.title = cfg.title;
      btn.addEventListener("click", cfg.onClick);
      toolbar.appendChild(btn);
    });

    return toolbar;
  };

  // ── イベント ────────────────────────────────────────────────

  MermaidRich.prototype._bindEvents = function () {
    const vp = this._viewport;

    // ホイールズーム
    vp.addEventListener("wheel", (e) => {
      e.preventDefault();
      const delta = e.deltaY < 0 ? this._opts.zoomStep : -this._opts.zoomStep;
      // マウス位置基点でズーム
      const rect = vp.getBoundingClientRect();
      const mx = e.clientX - rect.left;
      const my = e.clientY - rect.top;
      this._zoomAt(delta, mx, my);
    }, { passive: false });

    // ドラッグパン
    vp.addEventListener("mousedown", (e) => {
      if (e.button !== 0) return;
      this._dragging = true;
      this._dragStart = { x: e.clientX - this._tx, y: e.clientY - this._ty };
      vp.classList.add("grabbing");
    });

    window.addEventListener("mousemove", (e) => {
      if (!this._dragging) return;
      this._tx = e.clientX - this._dragStart.x;
      this._ty = e.clientY - this._dragStart.y;
      this._applyTransform();
    });

    window.addEventListener("mouseup", () => {
      if (!this._dragging) return;
      this._dragging = false;
      this._viewport.classList.remove("grabbing");
    });

    // タッチパン
    let lastTouch = null;
    let lastDist = null;

    vp.addEventListener("touchstart", (e) => {
      if (e.touches.length === 1) {
        lastTouch = { x: e.touches[0].clientX - this._tx, y: e.touches[0].clientY - this._ty };
      } else if (e.touches.length === 2) {
        lastDist = _touchDist(e);
      }
    }, { passive: true });

    vp.addEventListener("touchmove", (e) => {
      e.preventDefault();
      if (e.touches.length === 1 && lastTouch) {
        this._tx = e.touches[0].clientX - lastTouch.x;
        this._ty = e.touches[0].clientY - lastTouch.y;
        this._applyTransform();
      } else if (e.touches.length === 2 && lastDist !== null) {
        const dist = _touchDist(e);
        const ratio = dist / lastDist;
        const newScale = Math.min(Math.max(this._scale * ratio, this._opts.minZoom), this._opts.maxZoom);
        const rect = this._viewport.getBoundingClientRect();
        const cx = (e.touches[0].clientX + e.touches[1].clientX) / 2 - rect.left;
        const cy = (e.touches[0].clientY + e.touches[1].clientY) / 2 - rect.top;
        this._tx = cx - (cx - this._tx) * (newScale / this._scale);
        this._ty = cy - (cy - this._ty) * (newScale / this._scale);
        this._scale = newScale;
        this._applyTransform();
        lastDist = dist;
      }
    }, { passive: false });

    vp.addEventListener("touchend", () => {
      lastTouch = null;
      lastDist = null;
    }, { passive: true });
  };

  function _touchDist(e) {
    const dx = e.touches[0].clientX - e.touches[1].clientX;
    const dy = e.touches[0].clientY - e.touches[1].clientY;
    return Math.sqrt(dx * dx + dy * dy);
  }

  // ── ズーム / パン ────────────────────────────────────────────

  MermaidRich.prototype._zoom = function (delta) {
    const vp = this._viewport;
    const rect = vp.getBoundingClientRect();
    this._zoomAt(delta, rect.width / 2, rect.height / 2);
  };

  MermaidRich.prototype._zoomAt = function (delta, cx, cy) {
    const newScale = Math.min(
      Math.max(this._scale + delta, this._opts.minZoom),
      this._opts.maxZoom
    );
    this._tx = cx - (cx - this._tx) * (newScale / this._scale);
    this._ty = cy - (cy - this._ty) * (newScale / this._scale);
    this._scale = newScale;
    this._applyTransform();
  };

  MermaidRich.prototype._resetZoom = function () {
    this._scale = 1;
    this._tx = 0;
    this._ty = 0;
    this._applyTransform();
    this._centerDiagram();
  };

  MermaidRich.prototype._fitToView = function () {
    const svg = this._diagram.querySelector("svg");
    if (!svg) return;

    const vpRect = this._viewport.getBoundingClientRect();
    const svgW = svg.getBoundingClientRect().width / this._scale;
    const svgH = svg.getBoundingClientRect().height / this._scale;

    const scaleX = vpRect.width / svgW;
    const scaleY = vpRect.height / svgH;
    this._scale = Math.min(scaleX, scaleY) * 0.9;

    this._tx = (vpRect.width - svgW * this._scale) / 2;
    this._ty = (vpRect.height - svgH * this._scale) / 2;
    this._applyTransform();
  };

  MermaidRich.prototype._centerDiagram = function () {
    const svg = this._diagram.querySelector("svg");
    if (!svg) return;
    const vpRect = this._viewport.getBoundingClientRect();
    const svgRect = svg.getBoundingClientRect();
    this._tx = (vpRect.width - svgRect.width) / 2;
    this._ty = (vpRect.height - svgRect.height) / 2;
    this._applyTransform();
  };

  MermaidRich.prototype._applyTransform = function () {
    this._diagram.style.transform =
      "translate(" + this._tx + "px, " + this._ty + "px) scale(" + this._scale + ")";
    if (this._zoomLabel) {
      this._zoomLabel.textContent = Math.round(this._scale * 100) + "%";
    }
  };

  // ── テーマ ───────────────────────────────────────────────────

  MermaidRich.prototype._toggleTheme = function () {
    const isDark = this._wrapper.classList.contains("dark");
    this._wrapper.classList.toggle("dark", !isDark);
    this._wrapper.classList.toggle("light", isDark);
    this._opts.theme = isDark ? "light" : "dark";
  };

  // ── SVG ダウンロード ─────────────────────────────────────────

  MermaidRich.prototype._downloadSVG = function () {
    const svg = this._diagram.querySelector("svg");
    if (!svg) return;
    const blob = new Blob([svg.outerHTML], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "diagram.svg";
    a.click();
    URL.revokeObjectURL(url);
  };

  // ── レンダリング ─────────────────────────────────────────────

  /**
   * Mermaid テキストを描画する
   * @param {string} code  Mermaid 記法のテキスト
   */
  MermaidRich.prototype.render = function (code) {
    this._code = code;

    loadMermaid(() => {
      const mermaidTheme = this._opts.mermaidTheme ||
        (this._opts.theme === "dark" ? "dark" : "default");

      window.mermaid.initialize({
        startOnLoad: false,
        theme: mermaidTheme,
        securityLevel: "loose",
      });

      const id = "mr-" + Date.now() + "-" + Math.random().toString(36).slice(2);

      window.mermaid.render(id, code).then(({ svg }) => {
        this._diagram.innerHTML = svg;
        // 初期表示: 中央揃え + フィット
        requestAnimationFrame(() => {
          this._resetZoom();
          this._fitToView();
        });
      }).catch((err) => {
        this._diagram.innerHTML =
          '<p style="color:red;padding:16px;">Mermaid パースエラー: ' +
          err.message + "</p>";
      });
    });
  };

  // ─── 公開 API ────────────────────────────────────────────────

  /**
   * ページ上の data-mermaid-rich 属性を持つ要素を自動初期化する
   *
   * 使用例:
   *   <div data-mermaid-rich data-theme="dark">
   *     graph LR
   *       A --> B
   *   </div>
   */
  MermaidRich.autoInit = function () {
    document.querySelectorAll("[data-mermaid-rich]").forEach(function (el) {
      const code = el.textContent.trim();
      el.textContent = "";
      el.style.height = el.dataset.height || "400px";

      const instance = new MermaidRich(el, {
        theme: el.dataset.theme || "light",
      });
      instance.render(code);
    });
  };

  global.MermaidRich = MermaidRich;
})(window);
