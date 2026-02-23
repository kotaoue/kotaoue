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
| `GOOGLE_FIT_CLIENT_ID` | Google OAuth2 クライアント ID |
| `GOOGLE_FIT_CLIENT_SECRET` | Google OAuth2 クライアントシークレット |
| `GOOGLE_FIT_REFRESH_TOKEN` | Google OAuth2 リフレッシュトークン |

## 認証の準備

1. [Google Cloud Console](https://console.cloud.google.com/) で OAuth2 クライアント ID を作成
2. Fitness API を有効化
3. `offline` アクセスで認証し、リフレッシュトークンを取得
4. 取得した値を GitHub Secrets に登録 (`GOOGLE_FIT_CLIENT_ID`, `GOOGLE_FIT_CLIENT_SECRET`, `GOOGLE_FIT_REFRESH_TOKEN`)
