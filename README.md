# Go言語でつくるインタプリタ

[O'Reilly Japan - Go言語でつくるインタプリタ](https://www.oreilly.co.jp/books/9784873118222/)を実装するためのリポジトリです。

## 環境

```console
% brew install goenv
```

`~/.zshrc` に以下を追加し、 `source ~/.zshrc` を実行。

```shell
export GOENV_ROOT=$HOME/.goenv
export PATH=$GOENV_ROOT/bin:$PATH
eval "$(goenv init -)"
```

```console
% goenv install 1.23.1
% goenv local 1.23.1
% go version
go version go1.23.1 darwin/arm64
```

## 実行

```console
% go run ./main
```

## Monkey言語

```monkey
let one = 1;
let two = 2;

let add = fn(x, y) {
  x + y;
};

let result = add(one, two);
```
