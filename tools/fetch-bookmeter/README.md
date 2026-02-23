# fetch-bookmeter

読書メーターのデータを取得するツール

## Usage

```bash
cd tools/fetch-bookmeter

# よみたい本一覧を取得
go run main.go fetch-wish

go run main.go fetch-wish -user-id 104 -output ../../wish.json

# READMEをよみたい本でアップデート
go run main.go update-readme

go run main.go update-readme -wish-file ../../wish.json -readme ../../README.md
```
