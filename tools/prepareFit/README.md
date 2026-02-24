# prepareFit

Google Fit の `GOOGLE_FIT_REFRESH_TOKEN` を取得するためのツール

`client_secret.json` は Google Cloud Project で OAuth2 クライアントを作成したときにダウンロードできる JSON ファイルです。
OAuth クライアントの種類は **Desktop app** を選んでください。
このファイルを `tools/prepareFit/client_secret.json` として保存してから実行してください。

## Usage

```bash
cd tools/prepareFit

python3 -m venv .venv
source .venv/bin/activate

python3 -m pip install google-auth-oauthlib

python3 main.py
```

## トラブルシュート

ブラウザで `redirect_uri` を含むエラー (例: `redirect_uri_mismatch`) が出る場合は、
`client_secret.json` が Web アプリ用で作られている可能性が高いです。

1. [Google Cloud Console（Credentials）](https://console.cloud.google.com/apis/credentials) で OAuth クライアント ID を新規作成
1. 種類は Desktop app を選択
1. ダウンロードした JSON を `tools/prepareFit/client_secret.json` に置き換え

出力された `REFRESH_TOKEN` を `GOOGLE_FIT_REFRESH_TOKEN` として GitHub Secrets に登録してください。
