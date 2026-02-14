# aruba-go

[![GitHub Action](https://img.shields.io/github/actions/workflow/status/raviqqe/aruba-go/test.yaml?branch=main&style=flat-square)](https://github.com/raviqqe/aruba-go/actions)
[![License](https://img.shields.io/github/license/raviqqe/aruba-go.svg?style=flat-square)](https://github.com/raviqqe/aruba-go/blob/main/LICENSE)

The Go implementation of [Aruba](https://github.com/cucumber/aruba), the command-line application testing framework.

## Install

To install a standalone `agoa` command, run:

```sh
go install github.com/raviqqe/aruba-go/cmd/agoa@latest
```

To install it as a library for [`Godog`](https://github.com/cucumber/godog), run:

```sh
go get github.com/raviqqe/aruba-go@latest
```

## Usage

```sh
agoa
```

For more information, see `agoa --help`.

## License

[MIT](LICENSE)
