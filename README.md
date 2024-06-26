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

解凍して出てくる`q-cli`というファイルが実行ファイルです。パスが通る場所に移動させ、名前を`q`に変えてください。`README.md`が一緒に解凍されますが、削除して問題ありません。

### webhookの設定

環境変数を用いるか、設定ファイルを配置します。

#### 環境変数の場合

`Q_WEBHOOK_ID`: traQのWebhook ID
`Q_WEBHOOK_SECRET`: traQのWebhook シークレット

#### 設定ファイルを使う場合

以下のような`config.json`を用意します。

```json:config.json
{
  "webhook_id": "{{traQのWebhook ID}}",
  "webhook_secret": "{{traQのWebhookシークレット}}"
}
```

それを適切なディレクトリに配置します。以下を参考にして配置してください。(Goの`os.UserConfigDir()`を使用しています。)

> UserConfigDir returns the default root directory to use for user-specific configuration data. Users should create their own application-specific subdirectory within this one and use that.
>
> On Unix systems, it returns \$XDG_CONFIG_HOME as specified by https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if non-empty, else \$HOME/.config. On Darwin, it returns \$HOME/Library/Application Support. On Windows, it returns %AppData%. On Plan 9, it returns $home/lib.
>
> If the location cannot be determined (for example, $HOME is not defined), then it will return an error.

## Usage

```txt
Usage of q:
q [option] [message]
  -c string
        Send message as code block. Specify the language name (e.g. go, python, shell)
  -h    Print this message
  -i    Interactive mode
  -s    Accept input from stdin
  -v    Print version
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
q -i
```

```txt
q [1]: 1行目
q [1]: 2行目
q [1]: 
q [2]:
```

`-i` オプションを使用すると、ターミナルのように連続して、また、複数行のメッセージを送信できます(interactive mode)。改行では送信されず、`Ctrl-D`を入力すると送信されます。この例では、

```txt
1行目
2行目
```

のようになります。

また、何も入力されていない状態で`Ctrl-D`を入力するとinteractive modeを抜けられます。

---

```sh
echo foo | q -s
```

```txt
foo
```

`-s`オプションを使用すると、標準入力から値を受け取り、それをそのままメッセージとして投稿できます。機密情報を間違って投稿しないよう気を付けてください。

---

```sh
q -c go package main
```

````md
```go
package main
```
````

`-c`オプションを使用すると、コードブロックとして投稿できます。オプションで指定した値が最初の` ``` `の後に追加され、traQ上で適切なシンタックスハイライトが付きます。値は必ず指定しなくてはいけません。

上の`-s`オプションと組み合わせるとファイルの中身をきれいにtraQに投稿できます。

```sh
cat main.go | q -s -c go
```
