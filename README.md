# q-cli

[traQ](https://github.com/traPtitech/traQ) にwebhookを使ってメッセージを投稿するCLIツールです。

## Preparation

traQのWebhookを作成します。Secure形式のものを用意してください。
traPで使用しているtraQインスタンスで使用する場合は、bot-consoleから作成できます。
Webhook IDとWebhookシークレットが必要になります。

## Installation

1. [ダウンロード](#ダウンロード)
2. [Webhookの設定](#webhookの設定)

### ダウンロード

#### Go

```sh
go install github.com/ikura-hamu/q-cli@{{バージョン}}
```

`$GOPATH/bin`以下に`q-cli`という名前でインストールされます。

#### GitHub Release

リリースページから該当するOS、バージョンのものを探してダウンロードしてください。
https://github.com/ikura-hamu/q-cli/releases

```sh
curl -OJL https://github.com/ikura-hamu/q-cli/releases/download/{{version}}/q-cli_{{OS}}_{{architecture}}.(tar.gz|zip)
```

ダウンロードしたら解凍します。

```sh
tar -zxvf q-cli_{{OS}}_{{architecture}}.tar.gz
```

や

```sh
unzip q-cli_{{OS}}_{{architecture}}.zip
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

以下のような`.q-cli.yml`を用意します。

```json:.q-cli.yml
webhook_host: "{{traQのドメイン}}"
webhook_id: "{{traQのWebhook ID}}"
webhook_secret: "{{traQのWebhookシークレット}}"
```

このファイルをホームディレクトリに配置します。

## Usage

```txt
Usage:
  q [message] [flags]

Flags:
  -c, --code            Send message with code block
      --config string   config file (default is $HOME/.q-cli.yaml)
  -h, --help            help for q
  -l, --lang string     Code block language
  -v, --version         Print version information and quit
```

```sh
q -h
```

で確認できます。

### Example

```sh
q Hello World
```

traQに「Hello World」と投稿されます。

---

```sh
q
```

```txt
1行目
2行目

```

メッセージが何もない場合は、標準入力からテキストを受け取り、投稿します。`EOF`(`Ctrl-D`)を受け取るまで投稿を待ちます。

以下のようにしてファイルの中身を投稿できます。

```sh
q < a.txt
```

---

```sh
q -c text
```

````md
```
text
```
````

`--code`(`-c`)オプションを使用すると、コードブロックとして投稿できます。

`--lang`(`-l`)オプションで指定した値が最初の` ``` `の後に追加され、traQ上で適切なシンタックスハイライトが付きます。

```sh
q -c -l go < main.go
```

````txt
```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, world")
}
```
````
