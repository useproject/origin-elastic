package psnotify

import (
	"bytes"
	"sync"
)

/* bytePool 对象池*/
var bytePool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

//获取Pool
func AcquireBuffer() *bytes.Buffer {
	t := bytePool.Get().(*bytes.Buffer)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseBuffer(b *bytes.Buffer) {
	bytePool.Put(b)
}

/* cnMsgPool 对象池*/
var cnMsgPool = sync.Pool{
	New: func() interface{} {
		return new(cnMsg)
	},
}

//获取Pool
func AcquireCnMsgPool() *cnMsg {
	t := cnMsgPool.Get().(*cnMsg)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseCnMsgPool(c *cnMsg) {
	cnMsgPool.Put(c)
}

func (c *cnMsg) Reset() {
	c.Len = 0
	c.Seq = 0
	c.Ack = 0
	c.Flags = 0
	c.Id.Idx = 0
	c.Id.Val = 0
}

/* procEventHeaderPool 对象池*/
var procEventHeaderPool = sync.Pool{
	New: func() interface{} {
		return new(procEventHeader)
	},
}

//获取Pool
func AcquireProcEvent() *procEventHeader {
	t := procEventHeaderPool.Get().(*procEventHeader)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseProcEvent(c *procEventHeader) {
	procEventHeaderPool.Put(c)
}

func (p *procEventHeader) Reset() {
	p.Cpu = 0
	p.Timestamp = 0
	p.What = 0
}

/* forkProcEvent 对象池*/
var forkProcEventPool = sync.Pool{
	New: func() interface{} {
		return new(forkProcEvent)
	},
}

//获取Pool
func AcquireForkProcEvent() *forkProcEvent {
	t := forkProcEventPool.Get().(*forkProcEvent)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseForkProcEvent(c *forkProcEvent) {
	forkProcEventPool.Put(c)
}

func (t *forkProcEvent) Reset() {
	t.ChildPid = 0
	t.ParentPid = 0
	t.ChildTgid = 0
	t.ParentTgid = 0
}

/* procEventForkPool 对象池*/
var procEventForkPool = sync.Pool{
	New: func() interface{} {
		return new(ProcEventFork)
	},
}

//获取Pool
func AcquireProcEventFork() *ProcEventFork {
	t := procEventForkPool.Get().(*ProcEventFork)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseProcEventFork(c *ProcEventFork) {
	procEventForkPool.Put(c)
}

func (t *ProcEventFork) Reset() {
	t.ChildPid = 0
	t.ParentPid = 0
}

/* execProcEventPool 对象池*/
var execProcEventPool = sync.Pool{
	New: func() interface{} {
		return new(execProcEvent)
	},
}

//获取Pool
func AcquireExecProcEvent() *execProcEvent {
	t := execProcEventPool.Get().(*execProcEvent)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseExecProcEvent(c *execProcEvent) {
	execProcEventPool.Put(c)
}

func (t *execProcEvent) Reset() {
	t.ProcessPid = 0
	t.ProcessTgid = 0
}

/* procEventExecPool 对象池*/
var procEventExecPool = sync.Pool{
	New: func() interface{} {
		return new(ProcEventExec)
	},
}

//获取Pool
func AcquireProcEventExec() *ProcEventExec {
	t := procEventExecPool.Get().(*ProcEventExec)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseProcEventExec(c *ProcEventExec) {
	procEventExecPool.Put(c)
}

func (t *ProcEventExec) Reset() {
	t.Pid = 0
}

/* exitProcEventPool 对象池*/
var exitProcEventPool = sync.Pool{
	New: func() interface{} {
		return new(exitProcEvent)
	},
}

//获取Pool
func AcquireExitProcEvent() *exitProcEvent {
	t := exitProcEventPool.Get().(*exitProcEvent)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseExitProcEvent(c *exitProcEvent) {
	exitProcEventPool.Put(c)
}

func (t *exitProcEvent) Reset() {
	t.ProcessTgid = 0
	t.ProcessPid = 0
	t.ExitCode = 0
	t.ExitSignal = 0
}

/* procEventExitPool 对象池*/
var procEventExitPool = sync.Pool{
	New: func() interface{} {
		return new(ProcEventExit)
	},
}

//获取Pool
func AcquireProcEventExit() *ProcEventExit {
	t := procEventExitPool.Get().(*ProcEventExit)
	t.Reset()
	return t
}

//  回收到Pool
func ReleaseProcEventExit(c *ProcEventExit) {
	procEventExitPool.Put(c)
}

func (t *ProcEventExit) Reset() {
	t.Pid = 0
}
