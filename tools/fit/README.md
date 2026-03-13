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
3. OAuth2 クライアント作成時にダウンロードした `client_secret.json` を `tools/prepareFit/client_secret.json` に配置
4. `tools/prepareFit` で `offline` アクセス認証を行い、リフレッシュトークンを取得
5. 取得した値を GitHub Secrets に登録 (`GOOGLE_FIT_CLIENT_ID`, `GOOGLE_FIT_CLIENT_SECRET`, `GOOGLE_FIT_REFRESH_TOKEN`)

### Secrets に入れる値

`GOOGLE_FIT_CLIENT_ID` と `GOOGLE_FIT_CLIENT_SECRET` は、`client_secret.json` の値をそのまま使います。

- `GOOGLE_FIT_CLIENT_ID`: `installed.client_id`
- `GOOGLE_FIT_CLIENT_SECRET`: `installed.client_secret`（JSON 全体ではなく文字列1つ）
- `GOOGLE_FIT_REFRESH_TOKEN`: `python3 main.py` の `REFRESH_TOKEN: ...` に表示された値

`GOOGLE_FIT_REFRESH_TOKEN` は、上記の `client_id` / `client_secret` と同じ OAuth クライアントで発行したものを使ってください。
クライアントを作り直した場合は、`python3 main.py` を再実行して新しい `REFRESH_TOKEN` を再発行する必要があります。
