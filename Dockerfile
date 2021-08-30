## ビルド環境
# FROM golang:1.16-alpine3.14 # alpine版だとgccがないためビルドエラーとなる。利用しない。
FROM golang:1.16 as build-stage
# カレントワーキングディレクトリとして '/app' フォルダを指定する
WORKDIR /app
# カレントワーキングディレクトリ(/app)に現在のディレクトリ下のファイルをコピー
COPY . /app
# go.modを見て必要モジュールをダウンロード
RUN go mod tidy
# goをビルド
RUN go build main.go


## 本番環境
# FROM golang:1.16-alpine3.14 # alpine版ではgcoが使えないため実行エラーとなる。利用しない。
FROM golang:1.16 as production-stage
# ビルド済み資材をコピー
COPY --from=build-stage /app/main /app/main
# DBをコピー(サブディレクトリごとコピー)
COPY /db/ /app/db/
# カレントワーキングディレクトリとして 'app' フォルダを指定する
WORKDIR /app
# ビルド済みのgoを実行
# 実行時のカレントディレクトリから見た相対パスで、configで指定したパスにdbファイルが必要である。
CMD ["./main"]