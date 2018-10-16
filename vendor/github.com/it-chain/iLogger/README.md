# iLogger

iLooger is logger library for it-chain project.  
You don't need to set complicated logeer. Just Import this pckage and Write log.

## Getting Started with iLogger

### Installation
```
go get -u github.com/it-chain/iLogger
dep ensure
```

### Usage
#### Basic example
```
package main

import "github.com/it-chain/iLogger"

func main() {
    iLogger.Infof(nil, "This is Info log")
    iLogger.Warnf(nil, "This is Warn log")

    iLogger.SetToDebug()    // Would you use Debug logger? You should set debug state.
    iLogger.Debugf(nil, "This is Debug log")

    iLogger.Errorf(nil, "This is Error log")
    iLogger.Fatalf(nil, "This is Fatal log")
    iLogger.Panicf(nil, "This is Panic log")
}
```
![log-basic-example](./images/log_images.png)  
There are 6 Level for log
- Info
- Warn
- Debug
- Error
- Fatal
- Panic

#### Make log file
```
package main

import "github.com/it-chain/iLogger"

func main() {
    iLogger.EnableFileLogger(true, "./mylog.log")

    iLogger.Infof(nil, "This is Info log")
    iLogger.Warnf(nil, "This is Warn log")
}
```
You can write log to certain file path by calling EnableFileLogger function.

#### Print formatted string
```
package main

import "github.com/it-chain/iLogger"

func main() {
    iLogger.Info(nil, "This is Info log")
    iLogger.Infof(nil, "This is Infof log. %s", "This can use variadic arguments")
}
```
![info-infof-example](./images/info_infof.png)  
Funcfion named aftet 'f' is for variadic arguments.  
You can call the others level named after 'f'.
