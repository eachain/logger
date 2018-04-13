# logger
用于业务逻辑的日志组件

---

## logger有什么用？

实话说，对编码没有什么实际的好处。当需要查找业务关键日志时，你会发现该工具的效果。

用该工具，可以使你以不同的维度，将关键信息打到日志文件中，常见的维度有：函数名、用户ID、Trace ID等信息。

多说无益，show you the code.

---

## 先决条件

首先假定，目标系统中，都有特定的日志组件，一般日志组件都满足`logger.Logger`接口:

```go
type Logger interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}
```

设项目日志组件为`project/common/log.Logger`，获取方式为: `project/common/log.GetLogger()`，类似于:

```go
type Logger struct {
	// some fields
}

var logger *Logger // global

func GetLogger() *Logger {
	return logger
}
```

以下代码中，将使用下面函数:

```go
func CommonLogger() logger.Logger {
	return log.GetLogger()
}
```

---

## 函数名

一般而言，对于一个系统，如果能追踪一个函数内的所有事件，那查找起问题来，将很容易看出事件的先后关系，比如以下日志：

```
[INFO] FuncABC: start do something
[INFO] FuncABC: load some base info
[WARN] FuncABC: get something unimportant failed, reason: some error msg here
[INFO] FuncABC: write the result to db: some result info
```

很明显可以看到：当查找所有`FuncABC`业务相关的问题时，查找起来只需要一个命令:

```
grep FuncABC your.log
```

该函数的所有过程，详尽地展示在你面前，问题一览无遗。

**对应本工具`logger`的用法是:**

```go
log := logger.WithPrefix(CommonLogger(), "FuncABC: ")
log.Infof("something log here")
```

---

### 用户ID

很多时候，问题并不是某个单一函数出错，而是在之前某个动作中，就已经埋下BUG了，只是刚好到这里出错了。假设某一个用户出的特定问题，就需要找到该用户的特定行为，比如以下日志:

```
[INFO] FuncXYZ: <user 123456>: change nick to 'null'
[INFO] FuncABC: <user 123456>: start do something
[INFO] FuncABC: <user 123456>: load some base info
[WARN] FuncABC: <user 123456>: get something unimportant failed, reason: some error msg here
[INFO] FuncABC: <user 123456>: write the result to db: some result info
```

通过搜索该用户的所有行为:

```
grep 'user 123456' your.log
```

可以分析出，`[WARN]`出错日志的原因是，在之前，用户将用户名改成了`'null'`，导致查数据库失败。同时，定位出问题的出处是在`FuncXYZ`处。

**对应本工具`logger`的用法是:**

```go
log := logger.WithPrefix(CommonLogger(), "FuncABC: ")
log = logger.WithPrefix(log, "<user 123456>: ")
log.Infof("something log here")
```

---

### Trace ID

当问题超出当前系统时，就需要和其他系统一起联调，以找出问题，例如用TraceID找出这条调用链，分析问题的成因:

```
A.log:[INFO] FuncXYZ: <user 123456>: change nick to 'null' [trace: 0xabcdef]
A.log:[INFO] FuncXYZ: <user 123456>: call B.FuncABC [trace: 0xabcdef]
B.log:[INFO] FuncABC: <user 123456>: start do something [trace: 0xabcdef]
B.log:[INFO] FuncABC: <user 123456>: load some base info [trace: 0xabcdef]
B.log:[WARN] FuncABC: <user 123456>: get something unimportant failed, reason: some error msg here [trace: 0xabcdef]
B.log:[INFO] FuncABC: <user 123456>: write the result to db: some result info [trace: 0xabcdef]
```

通过搜索`TraceID`:

```
grep 0xabcdef A.log B.log
```

可以很快找到，问题的成因，也许根源在于B系统的sql写得不清真，导致用户调皮地将昵称改为`'null'`后，不能找到用户的部分信息。此时，只需看对应的sql，是否需要改进。

**对应本工具`logger`的用法是:**

```go
log := logger.WithSuffix(CommonLogger(), "[trace: 0xabcdef]")
log = logger.WithPrefix(log, "FuncABC: ")
log = logger.WithPrefix(log, "<user 123456>: ")
log.Infof("something log here")
```

*注意，trace id一般写在最开始，自己系统的问题可以慢慢调，但请不要影响合作方进度。*

---

## 更多

日志的维度有很多，不同系统中也不尽相同，有些数据字段天然可以用作一个维度，如用户ID。有些则需要杜撰一些信息，比如提供长连接的服务中，可以为每个长连接生成一个随机的`ConnID`，这样就可以通过`ConnID`找出这条长连接的过程中，所有发生的事件。

日志的意义在于方便查找问题，打的太少，出了问题不方便排查，打的太多又难免废话连篇。用本工具固然可以将日志以不同维度打印出来，但并不是说，所有的维度都应该打出来。比如：有人举报，一个叫`eachain`的用户(用户ID是系统背后的数据，真实用户不一定知道自己的ID)，在中午出现违规现象。这个问题如何处理呢？

很明显，先通过已知的关键字`eachain`搜索中午的日志:

```
grep eachain your.log
```

得到:

```
12:03:56.947 [INFO] FuncXYZ: <user 123456>: 'eachain' says: 'Fu ck U' [trace: 0xabcdef]
```

接下来，你就可以通过`[trace: 0xabcdef]`找到该调用链的所有过程信息。也可以根据`<user 123456>`找到该用户今天的所有行为，判定是否要将该用户禁言。

所以该工具的正确用法是：将真正重要、方便排查问题的信息，打印出来，其他非重要信息可以通过各种组合手段得到即可，比较简单粗暴的方式如：将请求信息、响应信息全打印出来。
