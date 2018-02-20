licy
====

[![GoDoc Widget]][GoDoc]

licy is a tool to generate `LICENSE` files from the command line.

A list of supported licenses can be found [here](https://github.com/github/choosealicense.com/tree/gh-pages/_licenses).

Installation
------------

#### Go
```bash
$ go get github.com/while-loop/licy
```

#### cURL
```bash
$ wget -O $GOPATH/bin/licy https://github.com/while-loop/licy/releases/download/v0.0.1/licy-GOOS-GOARCH && chmod +x $GOPATH/bin/licy
```

#### Source
```bash
$ make install     # will install licy to $GOPATH/bin
```

Example
-------

Generate and populated an MIT license file (`LICENSE`) with the current year and
given organization name

```bash
$ licy gen mit "Anthony Alves"
```

Commands
--------

- `licy list`
    - view all available licenses
- `licy help`
    - view available commands
- `licy info [-v] <license>`
    - get info about a license. Get full license with `-v` option
- `licy gen [-o FILE] [-y YEAR] [-p PROJECT_NAME] <license> <org>`
    - generate a license file with populated fields

License
-------
licy is licensed under the MIT License.
See [LICENSE](LICENSE) for details.

Author
------

Anthony Alves

[GoDoc]: https://godoc.org/github.com/while-loop/licy
[GoDoc Widget]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square