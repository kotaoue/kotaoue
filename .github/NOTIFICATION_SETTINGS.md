# GitHub通知設定：CopilotからのメールをフィルタリングするGitHub設定

## 概要

GitHub CopilotによるメンションやPR更新の通知メールを受け取らず、
チームメンバー（人間）からの通知は受け取り続けるための設定方法です。

この設定は**アカウント全体**に適用されます。

---

## 方法1: GitHubの通知設定でカスタムフィルターを作成する

GitHubの受信トレイには、複数のカスタムフィルターを作成できます。

### 手順

1. [GitHub通知設定ページ](https://github.com/settings/notifications) を開く
2. 画面左側の通知受信トレイで「**Filters**（フィルター）」を選択
3. 「**Create a filter**（フィルターを作成）」をクリック
4. 以下のフィルター条件を入力:
   - **Name（名前）**: `Copilot notifications`
   - **Filter query（クエリ）**: `author:app/copilot`
5. 保存して適用

> **注意**: このフィルターはGitHub上の受信トレイ表示の整理に使用します。メール送信を完全に止めるには、下記の方法2も組み合わせてください。

---

## 方法2: Copilotによる通知をボットとして無視する

GitHubでは特定のユーザーやボットからの通知を無視（ignore）できます。

### 手順

1. Copilotボットが行ったコメントや操作が含まれる通知を開く
2. 通知右上の「**...**」（その他のオプション）をクリック
3. 「**Ignore（無視）**」または「**Unsubscribe（購読解除）**」を選択

または、Copilotの通知が多いリポジトリ単位で設定する場合:

1. 対象リポジトリの「**Watch（ウォッチ）**」設定を変更
2. 「**Custom（カスタム）**」を選択し、受け取りたい通知タイプのみチェック

---

## 方法3: メールクライアントでフィルタリングする（最も確実）

GitHubからの通知メールには、フィルタリングに使える特定のヘッダー情報が含まれています。

### Gmailの場合

1. Gmailの検索バーに以下を入力して検索:
   ```
   from:notifications@github.com subject:copilot
   ```
2. 検索結果ページで「**このような検索条件のメールをフィルタリング**」をクリック
3. 「**ラベルを付ける**」で専用ラベル（例: `GitHub/Copilot`）を作成するか、
   「**受信トレイをスキップ（アーカイブする）**」にチェックを入れる
4. 「**フィルタを作成**」をクリック

### メールヘッダーを使った高度なフィルタリング

GitHubの通知メールには以下のカスタムヘッダーが含まれています:

| ヘッダー | 説明 |
|----------|------|
| `X-GitHub-Sender` | 通知の送信者（例: `copilot[bot]`） |
| `X-GitHub-Reason` | 通知の理由（例: `mention`, `review_requested`） |

Gmailのフィルターではヘッダーによるフィルタリングは直接できませんが、
Outlookやその他のメールクライアントでは対応しています。

---

## 方法4: GitHubの通知メール設定でボット通知を減らす

1. [GitHubの通知設定](https://github.com/settings/notifications) を開く
2. 「**Email notification preferences**」セクションを確認
3. 受け取りたい通知の種類（例: `Pull Request reviews`, `Comments`）のみにチェックを絞る

---

## まとめ

| 方法 | 効果 | 難易度 |
|------|------|--------|
| GitHub受信トレイのカスタムフィルター | 表示の整理 | 低 |
| ボットの通知を無視 | 特定通知の停止 | 低 |
| メールクライアントのフィルター | メール整理 | 中 |
| GitHub通知設定の絞り込み | 通知タイプの制限 | 低 |

最も効果的なのは **方法3（メールクライアントのフィルター）** と **方法4（GitHub設定の絞り込み）** の組み合わせです。

---

## 関連リンク

- [GitHubの通知設定](https://github.com/settings/notifications)
- [GitHubの通知管理ドキュメント（英語）](https://docs.github.com/en/subscriptions-and-notifications/get-started/configuring-notifications)
- [受信トレイの管理（英語）](https://docs.github.com/en/subscriptions-and-notifications/how-tos/viewing-and-triaging-notifications/managing-notifications-from-your-inbox)
