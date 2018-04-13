package logger_test

import (
	"fmt"

	"github.com/eachain/logger"
)

type simpleLogger struct{}

func (sl simpleLogger) Infof(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}

func (sl simpleLogger) Warnf(format string, a ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", a...)
}

func (sl simpleLogger) Errorf(format string, a ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", a...)
}

func ExampleLogger() {
	log := logger.WithPrefix(simpleLogger{}, "prefix: ")
	log.Infof("Hello world")

	log = logger.WithSuffix(log, ", suffix")
	log.Warnf("Hello world")

	log = logger.WithPrefix(log, "after pre: ")
	log.Errorf("Hello world")

	// Output:
	// [INFO] prefix: Hello world
	// [WARN] prefix: Hello world, suffix
	// [ERROR] prefix: after pre: Hello world, suffix
}
