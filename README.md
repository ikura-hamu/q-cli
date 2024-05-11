# q-cli

[traQ](https://github.com/traPtitech/traQ) にwebhookを使ってメッセージを投稿するCLIツールです。

## Preparation

traQのWebhookを作成します。Secure形式のものを用意してください。
traPで使用しているtraQインスタンスで使用する場合は、bot-consoleから作成できます。
Webhook IDとWebhookシークレットが必要になります。

## Installation

TODO

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
