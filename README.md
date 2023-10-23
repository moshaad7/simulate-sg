### Add deps
```
mkdir -p deps/src/github.com/couchbase
cd deps/src/github.com/couchbase
git clone https://github.com/moshaad7/cbgt.git
cd cbgt
git fetch origin
git checkout simulate-sg
```


### update go.mod file to use local cbgt
```
echo "replace github.com/couchbase/cbgt => ./deps/src/github.com/couchbase/cbgt" >> go.mod
go mod tidy
```

### command to run the script:
    go run main.go  >> logs.log  2>&1

### Check "pindexes to add" against "pindexes to remove"
```
./script.sh logs.log
```

In an ideal situation, on all the nodes, difference of pindexesToAdd and pindexesToRemove should be equal to 4.

