package singleton

import (
	"fmt"
	"sync"
)

var(
	lock = &sync.Mutex{}
	once = &sync.Once{}
	single *Single
)
type Single struct {

}

func GetInstance() *Single  {
	if single == nil{
		lock.Lock()
		defer lock.Unlock()
		single = &Single{}
		fmt.Println("Createing single instance now")
	}else{
		fmt.Println("Single Instance already created")
	}
	return single
}

func GetInstance2() *Single {
	if single == nil{
		once.Do(func() {
			single = &Single{}
		})
		fmt.Println("Createing single instance now")
	}else{
		fmt.Println("Single Instance already created")
	}
	return single
}