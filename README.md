# Sets Calculator

## Install

1. Install golang
```
sudo apt-get update
sudo apt-get install -y golang-go
```
2. Configure golang
```
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
```
3. Install scalc
```
go get -u github.com/mike-sul/scalc/cmd/scalc
```

## Test

```
cd $GOPATH/src/github.com/mike-sul/scalc/test
go test -v
scalc [ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]
```

## Run 

```
scalc <expression>

expression := “[“ operator set_1 set_2 set_3 … set_n “]”
operator := “SUM” | “INT” | “DIF”
set := file | expression
```

### Example

```
scalc [ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]
```



Files contains sorted integers, one integer in a line. At least one set should be given to each operator. Operator:
SUM - returns union of all sets
INT - returns intersection of all sets
DIF - returns difference of first set and the rest ones

```
$ cat a.txt
1
2
3
```
```
$ cat b.txt
2
3
4
```
```
$ cat c.txt
3
4
5
```
```
scalc [ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]
1
3
4
```