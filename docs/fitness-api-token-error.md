# Fitness API トークンエラーの対処方法

GitHub Actions で以下のようなエラーが発生した場合の原因と対処方法をまとめます。

```sh
failed to refresh access token: token endpoint returned status 400: {
  "error": "invalid_grant",
  "error_description": "Token has been expired or revoked."
}
```

## 原因

`invalid_grant` エラーは、保存している `GOOGLE_FIT_REFRESH_TOKEN` が無効になった場合に発生します。
主な原因は以下の通りです。

### 1. テストモードのアプリによるトークン有効期限切れ（最も可能性が高い）

[Google Cloud Console（OAuth consent screen）](https://console.cloud.google.com/apis/credentials/consent) で
Publishing status が **Testing** のまま運用している場合、発行されたリフレッシュトークンは **7日間で自動的に失効** します。

このリポジトリのような個人利用の自動化ツールでは、アプリを本番公開（Published）せずに使い続けるケースが多く、
定期的にトークンの再発行が必要になります。

### 2. ユーザーまたは Google によるアクセス取り消し

以下のいずれかの操作でリフレッシュトークンが失効します。

- [Google アカウントの「サードパーティアプリのアクセス」](https://myaccount.google.com/permissions) から手動でアクセスを取り消した
- Google が不審なアクティビティを検出してトークンを失効させた

### 3. Google Cloud の認証情報の変更

以下の変更を行った場合、既存のリフレッシュトークンが無効になることがあります。

- OAuth クライアント ID を削除・再作成した
- OAuth 同意画面のスコープを変更した
- Google Cloud プロジェクトに新しい API を有効化した（まれに影響する場合がある）

### 4. リフレッシュトークンの上限超過

Google は同一ユーザー・同一クライアントに対して発行できるリフレッシュトークンを **50件** に制限しています。
上限に達すると、最も古いトークンが自動的に失効します。

## 調査方法

1. [Google アカウントのセキュリティページ](https://myaccount.google.com/security) を開く
1. 「最近のセキュリティ アクティビティ」に不審なアクセス取り消しがないか確認する
1. [サードパーティアプリのアクセス](https://myaccount.google.com/permissions) を開き、
   対象の OAuth アプリ（このリポジトリで使っているアプリ）が表示されているか確認する
1. [Google Cloud Console（OAuth consent screen）](https://console.cloud.google.com/apis/credentials/consent) を開き、
   Publishing status が **Testing** になっているか確認する（Testing の場合は原因1が該当）

## 対処方法

リフレッシュトークンを再発行して GitHub Secrets を更新します。

1. `tools/prepareFit/main.py` を実行して新しいリフレッシュトークンを取得する
   （詳細手順は [tools/prepareFit/README.md](../tools/prepareFit/README.md) を参照）
1. 取得した値を GitHub Secrets の `GOOGLE_FIT_REFRESH_TOKEN` に上書き登録する
   - [Settings → Secrets and variables → Actions](https://github.com/kotaoue/kotaoue/settings/secrets/actions)
1. [GitHub Actions](https://github.com/kotaoue/kotaoue/actions/workflows/update-readme.yml) を手動実行して正常に動作するか確認する

### 再発防止策（Testing モードの場合）

テストモードによる7日失効を避けるには、以下のいずれかの対応を検討してください。

- **Publishing status を Published に変更する**: 個人利用の小規模アプリであれば、Google の審査なしに公開できる場合があります
- **定期的にトークンを再発行する**: `tools/prepareFit/main.py` を7日以内に再実行してトークンを更新する運用にする

## 参考リンク

- [Google OAuth 2.0 エラーリファレンス](https://developers.google.com/identity/protocols/oauth2/web-server#error-codes)
- [リフレッシュトークンの有効期限について](https://developers.google.com/identity/protocols/oauth2#expiration)
- [tools/prepareFit/README.md](../tools/prepareFit/README.md) - リフレッシュトークン取得手順
- [tools/fit/README.md](../tools/fit/README.md) - Fitness API ツールの概要
