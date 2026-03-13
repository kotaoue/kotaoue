# fit

Google Fit のデータを取得するツール

## Usage

```bash
cd tools/fit

# 昨日の歩数を取得して README.md を更新
go run . -cmd update-pedometer

go run . -cmd update-pedometer -readme ../../README.md
```

## 環境変数

| 変数名 | 説明 |
| - | - |
| `GOOGLE_CLOUD_CREDENTIALS_JSON` | Google OAuth2 認証情報 JSON (`authorized_user` 形式) |

## 認証の準備

1. [Google Cloud Console](https://console.cloud.google.com/) で **Fitness API 専用** の Google Cloud プロジェクトを作成
   - 他の API を有効化するプロジェクトとは分ける（同じプロジェクトで他の API を追加すると OAuth 同意画面の変更により既存のリフレッシュトークンが無効化されることがある）
2. そのプロジェクトで Fitness API を有効化
3. OAuth2 クライアント ID を作成（種類: **Desktop app**）
4. ダウンロードした `client_secret.json` を `tools/prepareFit/client_secret.json` に配置
5. `tools/prepareFit` で認証を行い、`GOOGLE_CLOUD_CREDENTIALS_JSON` の値を取得
6. 取得した JSON を GitHub Secrets に `GOOGLE_CLOUD_CREDENTIALS_JSON` として登録
