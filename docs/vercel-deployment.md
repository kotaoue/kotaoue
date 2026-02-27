# Vercel への github-readme-stats デプロイ手順

github-readme-stats の自分専用インスタンスを Vercel にデプロイする手順です。

## 1. リポジトリを Fork する

1. https://github.com/anuraghazra/github-readme-stats を開く
2. 右上の「Fork」ボタンをクリックして自分のアカウントにフォーク

## 2. GitHub Personal Access Token を作成する

### 推奨: Fine-grained token（最小権限）

1. https://github.com/settings/tokens を開く
2. 「Generate new token」→「Generate new token (fine-grained)」をクリック
3. 有効期限を設定する
4. 「Repository access」で「All repositories」を選択
5. 「Repository permissions」で以下を `Read-only` に設定:
   - Commit statuses
   - Contents
   - Issues
   - Metadata
   - Pull requests
6. 「Generate token」をクリックしてトークンをコピーしておく

> **注意:** Fine-grained token では、「All repositories」選択時は公開リポジトリのコミットのみカウントされます。プライベートリポジトリも含めるには「Only select repositories」で対象リポジトリを選択してください。

### 代替: Classic token

`repo` は権限が広すぎるため、公開統計のみ表示する場合は **`read:user` だけ**で十分です。

1. https://github.com/settings/tokens を開く
2. 「Generate new token (classic)」をクリック
3. スコープは `read:user` のみにチェックを入れて生成
4. 生成されたトークンをコピーしておく

> **注意:** プライベートリポジトリの貢献もカウントしたい場合は `repo` も追加してください。

## 3. Vercel にデプロイする

1. https://vercel.com にサインイン（GitHub アカウントで可）
2. 「Add New Project」→「Import Git Repository」でフォークしたリポジトリを選択
3. 「Environment Variables」に以下を追加:
   - Name: `PAT_1`
   - Value: 手順2でコピーしたトークン
4. 「Deploy」をクリック

## 4. README を更新する

デプロイ完了後、Vercel が発行した URL を使って README を更新:

```html
<a href="https://github.com/YOUR-USERNAME" rel="nofollow noreferrer noopener" target="_blank"><img src="https://YOUR-VERCEL-URL.vercel.app/api?username=YOUR-USERNAME&show_icons=true&theme=apprentice"/></a>
<a href="https://github.com/YOUR-USERNAME" rel="nofollow noreferrer noopener" target="_blank"><img src="https://YOUR-VERCEL-URL.vercel.app/api/top-langs/?username=YOUR-USERNAME&layout=compact&theme=apprentice"/></a>
```

`YOUR-VERCEL-URL` を実際に Vercel が発行したサブドメインに、`YOUR-USERNAME` を GitHub ユーザー名に置き換えてください。
