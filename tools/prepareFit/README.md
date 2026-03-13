# prepareFit

Google Fit の `GOOGLE_FIT_CREDENTIALS_JSON` を取得するためのツール

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

実行すると `GOOGLE_FIT_CREDENTIALS_JSON:` に続いて JSON が出力されます。
その JSON の値を GitHub Secrets の `GOOGLE_FIT_CREDENTIALS_JSON` に登録してください。

## トラブルシュート

ブラウザで `redirect_uri` を含むエラー (例: `redirect_uri_mismatch`) が出る場合は、
`client_secret.json` が Web アプリ用で作られている可能性が高いです。

1. [Google Cloud Console（Credentials）](https://console.cloud.google.com/apis/credentials) で OAuth クライアント ID を新規作成
1. 種類は Desktop app を選択
1. ダウンロードした JSON を `tools/prepareFit/client_secret.json` に置き換え
