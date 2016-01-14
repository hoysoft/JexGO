package JexGO
import (
	"time"
	"fmt"
)

type TaskManager struct {
	tasks map[interface{}]func() error
	breaktag bool
}

//新建task管理器，并指定最大线程数量
func NewTasks()*TaskManager{
	t:=new(TaskManager)
	t.tasks=make(map[interface{}] func() error)
	return t
}

//开始task服务，并指定最大线程数量
func (t *TaskManager)StartServe(maxGoroutineCount int){
	t.breaktag=false
	var maxGoroutine  chan int
	maxGoroutine = make(chan int, maxGoroutineCount)
	go func() {
		for {
			for k,v:=range t.tasks{
				if t.breaktag {break}

				maxGoroutine <- 1
				go func(r func() error) {
					fmt.Println("start task:",k)
					r()
					fmt.Println("finish task:",k)
					<-maxGoroutine
				}(v)
				delete(t.tasks,k)
			}
			if t.breaktag {break}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}


func (t *TaskManager)AddTask(tag interface{}, task func() error){
	t.tasks[tag]=task
}



func (t *TaskManager)Stop(){
	t.breaktag=true
}
