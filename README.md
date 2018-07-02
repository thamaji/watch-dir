watch-dir
====

### Usage

```
$ watch-dir -h

Usage: watch-dir [OPTIONS] DIR [DIR...]

Watch file system events

Options:
  -e value
    	set watch event (CREATAE|WRITE|REMOVE|RENAME|CHMOD)
  -h	show help
  -v	show version

```

### Example

basic

```
$ watch-dir -e CREATE -e WRITE .
./a
./b
```

```
$ touch a
$ echo b > b
```

send desktop notification

```
$ watch-dir -e CREATE -e WRITE . | xargs -L 1 notify-send
```
