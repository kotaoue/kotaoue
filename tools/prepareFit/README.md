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

### エラー 403: access_denied

次のエラーが表示される場合:

> このアプリは現在テスト中で、デベロッパーに承認されたテスターのみがアクセスできます。

OAuth 同意画面が **Testing** 状態で、ログインした Google アカウントがテスターに登録されていない可能性が高いです。

1. [Google Cloud Console（OAuth consent screen）](https://console.cloud.google.com/apis/credentials/consent) の`対象`を開く
1. User Type が External の場合、Publishing status が Testing であることを確認
1. Test users に、`python3 main.py` 実行時に使う Google アカウント（Gmail）を追加
1. 追加後、ブラウザを開き直して `python3 main.py` を再実行

補足:

- Workspace（社内）アカウントの場合は User Type が Internal の設定になっていないかも確認してください
- 自分の Google アカウントでログインしているつもりでも、ブラウザで別アカウントが選ばれていると同じエラーになります

### redirect_uri_mismatch

ブラウザで `redirect_uri` を含むエラー (例: `redirect_uri_mismatch`) が出る場合は、
`client_secret.json` が Web アプリ用で作られている可能性が高いです。

1. [Google Cloud Console（Credentials）](https://console.cloud.google.com/apis/credentials) で OAuth クライアント ID を新規作成
1. 種類は Desktop app を選択
1. ダウンロードした JSON を `tools/prepareFit/client_secret.json` に置き換え

出力された `REFRESH_TOKEN` を `GOOGLE_FIT_REFRESH_TOKEN` として GitHub Secrets に登録してください。
