# `q-cli`

`q-cli`は、[traQ](https://github.com/traQ)のWebhook機能を使ってtraQにテキストメッセージを送信するためのCLIツールです。

**このドキュメントはv0.5.0のものです。他のバージョンのドキュメントを確認したい場合は、GitHubの過去のリリースもしくは[「以前のバージョンのドキュメント」ページ](versions)を確認してください。**

## 各コマンド

- [`q`](q.md)
- [`q init`](q_init.md)
- [`q config`](q_config.md)

## 機能

- [secret detection](secret.md)

## Preparation

traQのWebhookをSecure形式で作成してください。
traPで使用しているtraQインスタンスで使用する場合は、bot-consoleから作成できます。
Webhook IDとWebhookシークレットが必要になります。

## Installation

1. [ダウンロード](#ダウンロード)
2. [Webhookの設定](#webhookの設定)

### ダウンロード

#### Go

```sh
go install github.com/ikura-hamu/q-cli@v0.5.0
```

`$GOPATH/bin`以下に`q-cli`という名前でインストールされます。

#### GitHub Release

リリースページから該当するOS、バージョンのものを探してダウンロードしてください。
https://github.com/ikura-hamu/q-cli/releases

```sh
curl -OJL https://github.com/ikura-hamu/q-cli/releases/download/v0.5.0/q-cli_{OS}_{architecture}.(tar.gz|zip)
```

ダウンロードしたら解凍します。

```sh
tar -zxvf q-cli_{OS}_{architecture}.tar.gz
```

や

```sh
unzip q-cli_{OS}_{architecture}.zip
```

など、適切な方法で解凍してください。

解凍して出てくる`q`というファイルが実行ファイルです。パスが通る場所に移動させてください。`README.md`が一緒に解凍されますが、削除して問題ありません。

### webhookの設定

環境変数を用いるか、設定ファイルを配置します。

#### 環境変数の場合

`Q_WEBHOOK_HOST`: traQのドメイン
`Q_WEBHOOK_ID`: traQのWebhook ID
`Q_WEBHOOK_SECRET`: traQのWebhook シークレット

#### 設定ファイルを使う場合

以下のような`~/.q-cli.yml`を用意します。

```json:.q-cli.yml
webhook_host: "{traQのドメイン}"
webhook_id: "{traQのWebhook ID}"
webhook_secret: "{traQのWebhookシークレット}"
channels:
  channel: "{チャンネルのUUID}"
```

または、

```sh
q init
```

を実行することで、対話形式で設定を行えます。
