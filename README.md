# logone-go
The library is a logging library that supports structured logging by generating a single line of JSON format for your specified lifecycle.

## Caution
Log messages are temporarily stored in memory.
Be careful in below point on a huge system.
* Memory usage.
* The log is not written out until the function ends.

## Usage
### Simple Example
```
import github.com/hixi-hyi/logone-go/logone

var manager *logone.Manager
func init() {
    manager = logone.NewManagerDefault()
}
func main() {
    ctx := context.Background()
    logger, finish := manager.Recording()
    defer finish()
    logger.Info("%s", "test")
    return
}

```
If there is a possibility that `panic()` will occur, call `finish()` using the `recover ()` statement.

### Using Nested Function
```
import github.com/hixi-hyi/logone-go/logone

var manager *logone.Manager
func init() {
    manager = logone.NewManagerDefault()
}
func main() {
    ctx := context.Background()
    ctx, finish := manager.RecordingWithContext(ctx)
    // logger, finish := manager.Recording()
    // ctx = logone.NewContextWithLogger(ctx, logger)
    defer finish()
    func2(ctx)
    return
}
func func2(ctx context.Context) {
    logger, _ := lambdalog.LoggerFromContext(ctx)
    logger.Info("nested function").WithTags("call-function")
}
```


### SetLogContext
You can add infomation about your context.
If you want to add the `REQUEST_ID` from HTTP Header, you ca write the following.
```
logger.SetLogContext(&logone.LogContext{
    "REQUEST_ID": req.Header['REQUEST_ID'],
})
```

### With Config
You can change the some parameters.
```
manager = lambdalog.NewManager(&lambdalog.Config{
	Type:            "request",
	DefaultSeverity: lambdalog.SeverityDebug,
	ElapsedUnit:     time.Millisecond,
	JsonIndent:      false
})
```

### With Tags
You can add annotations to investigate the log.
If you want to investigate the tag, you can get logs using `{ $.runtime.tags.aws-sdk-error >= 1 }`
```
res, err := sns.Publish(params)
if err != nil {
    log.Error("error occured in sns.Publish: %s", err).WithTags("aws-sdk-error")
    return err
}
```
### With Attributes
You can add an annotation to the log to get more details.
The `Attributes` value must be defined as a type that can be JSON.Marshal. If you want to output value that cannot be JSON.Marshal, you use fmt.Sprintf or primitive type. (e.g. error is not struct or primitive type, you must use fmt.Sprintf("%s", err) or err.Error())
```
log.Info("publish successfully").WithArrtibutes(res)
```
### With Error
You can add an annotation to the log to get more details.
```
if err != nil {
    log.Critical("error occured").WithError(err)
}
```

## Outputs
Below is an example of output from [./example/main.go].
```
{
  "type": "request",
  "context": {
    "REQUEST_ID": "xxxxxxx"
  },
  "runtime": {
    "severity": "DEBUG",
    "startTime": "2021-06-12T04:00:07.552913+09:00",
    "endTime": "2021-06-12T04:00:07.552969+09:00",
    "elapsed": 0,
    "lines": [
      {
        "severity": "DEBUG",
        "message": "invoked",
        "time": "2021-06-12T04:00:07.552964+09:00",
        "fileline": 19,
        "funcname": "main.main",
        "tags": [
          "critical"
        ],
        "attributes": "xxxx"
      }
    ],
    "tags": {
      "critical": 1
    }
  },
  "config": {
    "elapsedUnit": "1ms"
  }
}
```

## ToDo
* godoc
* Support to OutputFunc in lambdalog.Config. It is fmt.Println() now.
* Support to OutputSeverity in lambdalog.Config. It is print all logs now.
* Support to DefaultSeverity in lambdalog.Config. It is "UNKNOWN" now, if you don't write a any log.
* Support to OutputColumns in lambdalog.Config. It is `RFILELINE | RFUNCNAME | CELAPSED_UNIT` now.
* Do not want to consider about Attributes limitation. Attributes is support to only type can be JSON.Marshal or primitive now.
