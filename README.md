recc
============

Record command and output.


Installation
-----------

```sh
go get github.com/pocke/recc
```

<!-- Or download a binary from [Latest release](https://github.com/pocke/recc/releases/latest). -->


Usage
-----------

```sh
# Record `ls -la`
$ recc ls -la
# Same
$ recc 'ls -la'

# Record with a glob
$ recc 'ls -la *'

# Record with StdErr (default: off)
$ recc --stderr curl https://example.com

# Specify output (default: clipboard)
$ recc --output /tmp/recc.out ls -la
```

License
-------

These codes are licensed under CC0.

[![CC0](http://i.creativecommons.org/p/zero/1.0/88x31.png "CC0")](http://creativecommons.org/publicdomain/zero/1.0/deed.en)
