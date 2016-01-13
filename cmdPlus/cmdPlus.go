package spider
import (
	"bufio"
	"io"
	"regexp"
	"fmt"
	"os/exec"
	"time"
)

type CmdPlus struct {
	Cmd *exec.Cmd
	ReturnCode int
	TimedOut   bool
	Elapsed    time.Duration
	regexpKeys  []string
	TriggerKeyCallback TriggerKeyCallbackFunc ////正则表达式匹配时触发回调
	FinishCallback FinishCallbackFunc
	OutPutCallback OutPutCallbackFunc
}

type FinishCallbackFunc func (error);
type TriggerKeyCallbackFunc func (map[string]string);
type OutPutCallbackFunc func (string);

func NewCmdPlus(cmd *exec.Cmd, workDir string)*CmdPlus{
	c:=&CmdPlus{Cmd:cmd}
	c.Cmd.Dir=workDir
	return c
}

func (this *CmdPlus)Exec(){
   var err error

	ch := make(chan string)
	go func() {
		err = this.runCommandCh(ch)
	}()

	go func() {
		for v := range ch {
			//正则表达式提取关键字触发回调
			this.regexpTriggerKeys(v)
			//输出信息
			if this.OutPutCallback!=nil {
				this.OutPutCallback(v)
			}
		}

		if this.FinishCallback != nil {
			this.FinishCallback(err)
		}
	}()
}

func (this *CmdPlus)runCommandCh(stdoutCh chan <- string) error {
	output, err := this.Cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("RunCommand: cmd.StdoutPipe(): %v", err)
	}
	outstderr, err := this.Cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("RunCommand: cmd.StderrPipe(): %v", err)
	}
	if err := this.Cmd.Start(); err != nil {
		return fmt.Errorf("RunCommand: cmd.Start(): %v", err)
	}

	this.parsLineData(stdoutCh,output)
	this.parsLineData(stdoutCh,outstderr)
//	defer close(stdoutCh)

	if err := this.Cmd.Wait(); err != nil {
		//		if exiterr, ok := err.(*exec.ExitError); ok {
		//			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
		//				this.ReturnCode = status.ExitStatus()
		//				//fmt.Println("ccc:",status.ExitStatus())
		//			}
		//		} else {
		close(stdoutCh)
		return fmt.Errorf("RunCommand: cmd.Wait(): %v", err)
		//log.Fatalf("cmd.Wait return invalid result: %v\n%s\n", err, debug.Stack())
		//		}
	}
	close(stdoutCh)
	return nil

}

//正则表达式提取关键字触发回调
func (this *CmdPlus)regexpTriggerKeys(line string){
	if this.TriggerKeyCallback==nil {return }
	for _,v:=range this.regexpKeys{
		var digitsRegexp = regexp.MustCompile(v)
		m:=FindStringSubmatchMap(digitsRegexp,line)
		if m!=nil && len(m)>0  {
			this.TriggerKeyCallback(m)
		}
	}

}

//解析行数据
func (this *CmdPlus)parsLineData(stdoutCh chan <- string,output io.ReadCloser ) {
	go func() {
		for {
			r := bufio.NewReader(output)
			line, isPrefix, err := r.ReadLine()
			if err == nil  && !isPrefix {
				stdoutCh <- string(line)
			}
			if err == io.EOF {break}
		}
	}()
}

//设置正则表达式触发回调
func (this *CmdPlus)SetTriggerRegexpKeys(m ...string)*CmdPlus {
	this.regexpKeys=m
	return this
}

// 设置全部结束回调
func (self *CmdPlus) SetFinishCallback(callback  FinishCallbackFunc) {
	self.FinishCallback = callback
}