## ビルド環境
# FROM golang:1.16-alpine3.14 # alpine版だとgccがないためエラーとなる。利用しない。
FROM golang:1.16
# カレントワーキングディレクトリとして 'app' フォルダを指定する
WORKDIR /app
# カレントワーキングディレクトリ(/app)に現在のディレクトリ下のファイルをコピー
COPY . /app
# go.modを見て必要モジュールをダウンロード
RUN go mod tidy

# goを実行
CMD ["go", "run", "/app/main.go"]