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

実行すると `REFRESH_TOKEN: xxxx` が表示されます。

## 取得したトークンの GitHub Secrets への登録

以下の値を GitHub Secrets に登録してください。
[Settings → Secrets and variables → Actions](https://github.com/kotaoue/kotaoue/settings/secrets/actions) から登録・更新できます。

| Secret 名 | 値の出所 |
| - | - |
| `GOOGLE_FIT_CLIENT_ID` | `client_secret.json` の `installed.client_id` |
| `GOOGLE_FIT_CLIENT_SECRET` | `client_secret.json` の `installed.client_secret` |
| `GOOGLE_FIT_REFRESH_TOKEN` | `python3 main.py` の出力 `REFRESH_TOKEN: ...` |

登録後は [GitHub Actions（update-readme）](https://github.com/kotaoue/kotaoue/actions/workflows/update-readme.yml) を手動実行して正常に動作するか確認してください。

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

### invalid_grant: Token has been expired or revoked

GitHub Actions で以下のエラーが出た場合:

```sh
failed to refresh access token: token endpoint returned status 400: {
  "error": "invalid_grant",
  "error_description": "Token has been expired or revoked."
}
```

**原因**: OAuth 同意画面の Publishing status が **Testing** の場合、リフレッシュトークンは **7日間で自動的に失効** します。

**対処方法**:

1. [Google Cloud Console（OAuth consent screen）](https://console.cloud.google.com/apis/credentials/consent) を開く
1. Publishing status を **In production**（公開済み）に変更する
   - 個人利用の小規模アプリであれば、Google の審査なしに公開できます
1. `python3 main.py` を再実行して新しいリフレッシュトークンを取得する
1. 取得した値を GitHub Secrets の `GOOGLE_FIT_REFRESH_TOKEN` に更新する

Publishing status を変更しない場合は、7日以内に `python3 main.py` を再実行して新しいリフレッシュトークンを取得し、GitHub Secrets を更新し続ける必要があります。
