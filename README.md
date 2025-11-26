# scope
Tool for parsing and querying data from application and debug logs.

### how to use
Simply add logs to your application in the following format:

```
[event:<eventName>] {"value": 1, "test": "success", "measure": 123}
```

Scope will parse this and insert the json payload into a sql lite database (located at `/tmp/scope.db`) table matching your `eventName`.  Scope takes input via stdin, so you can pipe a flatfile into scope like this:

```
cat myText.txt | ./scope
```

or you can have logs stream indefinitely to scope with the `--watch` flag

```
adb logcat | ./scope --watch
```

**Note:** if you grep the output of adb logcat, be sure to use --line-buffered if also using --watch

Then you can query your events using sqlite's json support with `sqlite3 /tmp/scope.db`.  Feel free to point a visualization tool of your choice at this db.

```
SELECT json_extract(payload, $.test) as testVal from eventName;
```

### how to build
```
go build -o ./bin/scope ./cmd/scope/main.go
```

### ideas going forward:
- tui integration to visualize incoming data based on a provided sql statement