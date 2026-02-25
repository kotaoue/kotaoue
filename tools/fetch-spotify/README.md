# fetch-spotify

SpotifyのプレイリストからトラックデータをJSON形式で取得するツール

## Usage

```bash
cd tools/fetch-spotify

# プレイリストのトラック一覧を取得
SPOTIFY_CLIENT_ID=<your_client_id> SPOTIFY_CLIENT_SECRET=<your_client_secret> \
  go run . fetch-playlist

# オプション指定
SPOTIFY_CLIENT_ID=<your_client_id> SPOTIFY_CLIENT_SECRET=<your_client_secret> \
  go run . fetch-playlist -playlist-id 3aARAs2A4PgdkgYzcyYPgI -output ../../playlist.json
```

## 環境変数

| 変数名 | 説明 |
|---|---|
| `SPOTIFY_CLIENT_ID` | Spotify Developer Dashboard で取得したクライアントID |
| `SPOTIFY_CLIENT_SECRET` | Spotify Developer Dashboard で取得したクライアントシークレット |

## 出力形式

```json
[
  {
    "no": 1,
    "title": "トラック名",
    "url": "https://open.spotify.com/track/...",
    "artist": "アーティスト名",
    "thumb": "https://i.scdn.co/image/...",
    "date": "2024-01-01T00:00:00Z"
  }
]
```
