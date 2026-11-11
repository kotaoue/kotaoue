---
title: "Google Fit の歩数を GitHub Profile README に毎日自動表示する"
emoji: "🐔"
type: "tech"
topics: ["githubactions", "go", "googlefit", "googleapi"]
published: false
---

## やったことの概要

GitHub の Profile README に「昨日の歩数」を毎日自動で表示できるようにしました。

Google Fit API から前日の歩数を取得し、GitHub Actions で定期実行して README のマーカー間を書き換えるという仕組みです。

完成イメージは以下のとおりです。

| 今日の東京ソング | 2月25日の歩数 |
| - | - |
| (Spotify のアルバムアート) | 1756歩 |

---

## 設定方法

### 1. Google Cloud Console で OAuth2 クライアントを作成する

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセス
2. 新規プロジェクトを作成（または既存プロジェクトを選択）
3. **APIとサービス → ライブラリ** から **Fitness API** を有効化
4. **APIとサービス → 認証情報** で「OAuth クライアント ID の作成」を選択
   - アプリケーションの種類: **デスクトップ アプリ**
5. 作成後、`client_secret.json` をダウンロード

### 2. リフレッシュトークンを取得する

`tools/prepareFit/` に `client_secret.json` を置いて以下を実行します。

```bash
cd tools/prepareFit

python3 -m venv .venv
source .venv/bin/activate
pip install google-auth-oauthlib

python3 main.py
```

ブラウザで認可フローが開くので許可すると、ターミナルに `REFRESH_TOKEN` が表示されます。

### 3. GitHub Secrets に登録する

リポジトリの **Settings → Secrets and variables → Actions** で以下の 3 つを登録します。

| Secret 名 | 値 |
| - | - |
| `GOOGLE_FIT_CLIENT_ID` | OAuth2 クライアント ID |
| `GOOGLE_FIT_CLIENT_SECRET` | OAuth2 クライアントシークレット |
| `GOOGLE_FIT_REFRESH_TOKEN` | 手順 2 で取得したリフレッシュトークン |

### 4. README にマーカーを追加する

README.md の書き換えたい箇所に HTML コメントのマーカーを埋め込みます。

```markdown
| <!-- PEDOMETER_DATE_START -->日付<!-- PEDOMETER_DATE_END --> |
| <!-- PEDOMETER_STEPS_START -->歩数<!-- PEDOMETER_STEPS_END --> |
```

### 5. GitHub Actions ワークフローで定期実行する

`.github/workflows/update-readme.yml` に以下のステップを追加します。

```yaml
- name: Update README with yesterday's step count
  working-directory: tools/fit
  env:
    GOOGLE_FIT_CLIENT_ID: ${{ secrets.GOOGLE_FIT_CLIENT_ID }}
    GOOGLE_FIT_CLIENT_SECRET: ${{ secrets.GOOGLE_FIT_CLIENT_SECRET }}
    GOOGLE_FIT_REFRESH_TOKEN: ${{ secrets.GOOGLE_FIT_REFRESH_TOKEN }}
  run: go run . -cmd update-pedometer -readme $GITHUB_WORKSPACE/README.md
```

`cron: '0 15 * * *'` (JST 0:00) で毎日深夜に実行されます。

---

## シンプルに変更したコード

### リフレッシュトークン → アクセストークン (`service/auth.go`)

```go
func refreshAccessToken(clientID, clientSecret, refreshToken string) (string, error) {
    form := url.Values{}
    form.Set("client_id", clientID)
    form.Set("client_secret", clientSecret)
    form.Set("refresh_token", refreshToken)
    form.Set("grant_type", "refresh_token")

    resp, err := http.PostForm("https://oauth2.googleapis.com/token", form)
    // ...省略
    return tr.AccessToken, nil
}
```

### Google Fit API で昨日の歩数を取得 (`service/pedometer.go`)

```go
func fetchYesterdaySteps(accessToken string) (int, error) {
    jst := time.FixedZone("JST", 9*60*60)
    now := time.Now().In(jst)
    todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
    yesterdayStart := todayStart.AddDate(0, 0, -1)

    reqBody := aggregateRequest{
        AggregateBy:  []aggregateBy{{DataTypeName: "com.google.step_count.delta"}},
        BucketByTime: bucketByTime{DurationMillis: 86400000}, // 1日分
        StartTimeMs:  yesterdayStart.UnixMilli(),
        EndTimeMs:    todayStart.UnixMilli(),
    }
    // HTTP リクエストして total を集計して返す
}
```

### README のマーカー間を書き換え

```go
func replaceBetweenMarkers(content, start, end, replacement string) (string, error) {
    startIdx := strings.Index(content, start)
    endIdx := strings.Index(content, end)
    if startIdx == -1 || endIdx == -1 {
        return "", fmt.Errorf("markers not found: %q, %q", start, end)
    }
    if startIdx >= endIdx {
        return "", fmt.Errorf("start marker must appear before end marker")
    }
    return content[:startIdx+len(start)] + replacement + content[endIdx:], nil
}
```

マーカーの開始位置から終了位置までを `replacement` で差し替えるだけのシンプルな実装です。マーカーが見つからない場合はエラーを返してREADMEの破損を防いでいます。

---

## 実際のリポジトリ

https://github.com/kotaoue/kotaoue

コードは `tools/fit/` と `tools/prepareFit/` 以下にあります。
