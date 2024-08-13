package plog

import (
	"cvgo/provider/core"
	"fmt"
	"github.com/spf13/cast"
)

var plogSvc *PlogService

type PlogService struct {
	Service
	c     core.Container
	level byte
}

// 日志级别，只记录大于配置级别的日志
const (
	trace = iota // 最低级别，默认。所有日志，完整链路追踪
	debug        // 开发调试信息
	info         // 业务需要收集的有用信息，例如访客 UA、请求耗时等
	warn         // 警告
	err          // 一般运行时错误
	fatal        // 最高级别，重要性最高，记录导致应用 panic 崩溃的严重错误，
	off          // 日志开关，用于关闭日志的记录
)

type Service interface {
	Trace(output ...interface{})
	Tracef(output ...interface{})
	Debug(output ...interface{})
	Debugf(output ...interface{})
	Info(output ...interface{})
	Infof(output ...interface{})
	Warn(output ...interface{})
	Warnf(output ...interface{})
	Error(output ...interface{})
	Errorf(output ...interface{})
	Fatal(output ...interface{})
	Fatalf(output ...interface{})

	Color(color string, output interface{})
	Colorf(color string, output ...interface{})

	//P(output interface{})
}

// trace
func (self *PlogService) Trace(out ...interface{}) {
	if self.level > trace {
		return
	}
	self.output(trace, out...)
}

func (self *PlogService) Tracef(out ...interface{}) {
	if self.level > trace {
		return
	}
	self.output(trace, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// debug
func (self *PlogService) Debug(out ...interface{}) {
	if self.level > debug {
		return
	}
	self.output(debug, out...)
}

func (self *PlogService) Debugf(out ...interface{}) {
	if self.level > debug {
		return
	}
	self.output(debug, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// info
func (self *PlogService) Info(out ...interface{}) {
	if self.level > info {
		return
	}
	self.output(info, out...)
}

func (self *PlogService) Infof(out ...interface{}) {
	if self.level > info {
		return
	}
	self.output(info, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// warn
func (self *PlogService) Warn(out ...interface{}) {
	if self.level > warn {
		return
	}
	self.output(warn, out...)
}

func (self *PlogService) Warnf(out ...interface{}) {
	if self.level > warn {
		return
	}
	self.output(warn, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// err
func (self *PlogService) Error(out ...interface{}) {
	if self.level > err {
		return
	}
	self.output(err, out...)
}

func (self *PlogService) Errorf(out ...interface{}) {
	if self.level > err {
		return
	}
	self.output(err, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// fatal
func (self *PlogService) Fatal(out ...interface{}) {
	if self.level > fatal {
		return
	}
	self.output(fatal, out...)
}

func (self *PlogService) Fatalf(out ...interface{}) {
	if self.level > fatal {
		return
	}
	self.output(fatal, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// color
func (s *PlogService) Color(color string, output interface{}) {
	plogSvc.P(color, output)
}

func (s *PlogService) Colorf(color string, out ...interface{}) {
	plogSvc.P(color, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}
