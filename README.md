# glog4logrus

Print logrus logs with glog.

## Purpose

Both `glog` and `logrus` have thousands of stars on github and large amount of users in the world.

In case writing code that intend to be imported by others, it is a question to be `glog` or to be `logrus`. 
With `glog4logrus`, you could just log with `logrus`, and set output to `glog` if work with `glog` code. 

## Usage

#### Init glog in your main:

```go
	flag.Parse()
	defer glog.Flush()
```

#### Use `glog4logrus` as Formatter/Output of `logrus`

```go
	logrus.SetLevel(GlogLevel())
	logrus.SetFormatter(&GlogFormatter{})
	logrus.SetOutput(&GlogOuptut{})

```

#### And logging with `logrus`

```go
	logrus.Debug("this is debug")
	logrus.Info("this is info")
	logrus.Printf("this is a %s", "printf")

	xlogger := logrus.WithField("logger", "xlogger")
	xlogger.Info("this is info")
	xlogger.WithField("special", "should be quote").Info("this is info again")

	xlogger.WithField("annotation", "long line").Info("here is a long line, it should be very very very very long, longer than the magic 44 characters")

	logrus.Info("")
```

#### You will got `glog` output

```
I1129 01:31:48.414746   19351 logs_test.go:19] this is from glog
I1129 01:31:48.414866   19351 logs_test.go:23] this is a printf                            
I1129 01:31:48.414874   19351 logs_test.go:26] this is info                                 logger=xlogger
I1129 01:31:48.414884   19351 logs_test.go:27] this is info again                           logger=xlogger special="should be quote"
I1129 01:31:48.414891   19351 logs_test.go:29] here is a long line, it should be very very very very long, longer than the magic 44 characters annotation="long line" logger=xlogger
I1129 01:31:48.414901   19351 logs_test.go:31]                                             
``` 

> Note: if we set `-v=1` argument, we will log the `Debug` messages
