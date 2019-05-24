# 跳出Go module的泥潭

Go 1.11 前天已经正式发布了，这个版本包含了两个最重要的feature就是 `module`和`web assembly`。虽然也有一些简单的教程介绍了`go module`的特性，但是基本上都是`hello world`的例子，在实践的过程中， 很多人都在“拼命的挣扎”，包括我自己， 从一些qq群、github的issue, twitter上都可以看到大家茫然或者抱怨的语句。

虽然有三个帮助文件`go help mod`、`go help modules`、`go help module-get`可以了解一些go module的用法，但是感觉Go开发组对`module`这一特性还是没有很好的做一个全面的介绍，很多情况还得靠大家看源代码或者去猜，比如module下载的文件夹、版本格式的完整声明，`module`的最佳实践等，并且当前Go 1.11的实现中还有一些bug,给大家在使用的过程中带来了很大的困难。

我也在摸索中前行， 记录了摸索过程中的一些总结，希望能给还在挣扎中的Gopher一些帮助。

[Introduction to Go Modules](https://roberto.selbach.ca/intro-to-go-modules/) 是一篇很好的go module 入门介绍， 如果你仔细阅读了它，应该就不需要看本文了。



### GO111MODULE

要使用`go module`,首先要设置`GO111MODULE=on`,这没什么可说的，如果没设置，执行命令的时候会有提示，这个大家应该都了解了。

### 既有项目

假设你已经有了一个go 项目， 比如在`$GOPATH/github.com/smallnest/rpcx`下， 你可以使用`go mod init github.com/smallnest/rpcx`在这个文件夹下创建一个空的`go.mod` (只有第一行 `module github.com/smallnest/rpcx`)。

然后你可以通过 `go get ./...`让它查找依赖，并记录在`go.mod`文件中(你还可以指定 `-tags`,这样可以把tags的依赖都查找到)。

通过`go mod tidy`也可以用来为`go.mod`增加丢失的依赖，删除不需要的依赖，但是我不确定它怎么处理`tags`。

执行上面的命令会把`go.mod`的`latest`版本换成实际的最新的版本，并且会生成一个`go.sum`记录每个依赖库的版本和哈希值。

### 新的项目

你可以在`GOPATH`之外创建新的项目。

`go mod init packagename`可以创建一个空的`go.mod`,然后你可以在其中增加`require github.com/smallnest/rpcx latest`依赖，或者像上面一样让go自动发现和维护。

`go mod download`可以下载所需要的依赖，但是依赖并不是下载到`$GOPATH`中，而是`$GOPATH/pkg/mod`中，多个项目可以共享缓存的module。

### go mod命令

```
download    download modules to local cache (下载依赖的module到本地cache))
edit        edit go.mod from tools or scripts (编辑go.mod文件)
graph       print module requirement graph (打印模块依赖图))
init        initialize new module in current directory (再当前文件夹下初始化一个新的module, 创建go.mod文件))
tidy        add missing and remove unused modules (增加丢失的module，去掉未用的module)
vendor      make vendored copy of dependencies (将依赖复制到vendor下)
verify      verify dependencies have expected content (校验依赖)
why         explain why packages or modules are needed (解释为什么需要依赖)
```

有些命令还有bug, 比如`go mod download -dir`:

```
go mod download -dir /tmp
flag provided but not defined: -dir
usage: go mod download [-dir] [-json] [modules]
Run 'go help mod download' for details.
```

帮助里明明说可以设置`dir`,但是实际却不支持`dir`参数。

看这些命令的帮助已经比较容易了解命令的功能。

### 翻墙

在国内访问`golang.org/x`的各个包都需要翻墙，你可以在`go.mod`中使用`replace`替换成github上对应的库。

```
replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
```

依赖库中的`replace`对你的主`go.mod`不起作用，比如`github.com/smallnest/rpcx`的`go.mod`已经增加了`replace`,但是你的`go.mod`虽然`require`了`rpcx`的库，但是没有设置`replace`的话， `go get`还是会访问`golang.org/x`。

所以如果想编译那个项目，就在哪个项目中增加`replace`。

### 版本格式

下面的版本都是合法的：

```
gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
gopkg.in/vmihailenco/msgpack.v2 v2.9.1
gopkg.in/yaml.v2 <=v2.2.1
github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
latest
```

### go get 升级

- 运行 `go get -u` 将会升级到最新的次要版本或者修订版本(x.y.z, z是修订版本号， y是次要版本号)
- 运行 `go get -u=patch` 将会升级到最新的修订版本
- 运行 `go get package@version` 将会升级到指定的版本号`version`

### go mod vendor

`go mod vendor` 会复制modules下载到vendor中, 貌似只会下载你代码中引用的库，而不是go.mod中定义全部的module。

### go module, vendor 和 Travis CI

https://arslan.io/2018/08/26/using-go-modules-with-vendor-support-on-travis-ci/