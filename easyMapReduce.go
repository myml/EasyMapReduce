package easyMapReduce

import (
	"fmt"
	"reflect"
)

func mapreduce(source interface{}, goroutineLimit int, Map interface{}, Reduce interface{}) (err error) {
	defer func() {
		panicV := recover()
		if panicV != nil {
			err = fmt.Errorf("%v", panicV)
		}
		return
	}()
	vChan := reflect.ValueOf(source)
	if vChan.Type().Kind().String() != "chan" {
		if vChan.Len() > 0 {
			tChan := reflect.ChanOf(reflect.BothDir, vChan.Index(0).Type())
			vChan = reflect.MakeChan(tChan, 0)
			go func() {
				vSlice := reflect.ValueOf(source)
				for i, l := 0, vSlice.Len(); i < l; i++ {
					vChan.Send(vSlice.Index(i))
				}
				vChan.Close()
			}()
		}
	}
	mF := reflect.ValueOf(Map)
	rF := reflect.ValueOf(Reduce)
	exitChan := make(chan []struct{})
	retChan := make(chan []reflect.Value)
	go func() {
		for ret := range retChan {
			e := rF.Call(ret)[0]
			if !e.IsNil() {
				err = fmt.Errorf("%v", e)
				close(exitChan)
				return
			}
		}
	}()
	limitChan := make(chan struct{}, goroutineLimit)
	for {
		v, ok := vChan.Recv()
		if !ok {
			break
		}
		select {
		case <-exitChan:
			break
		case limitChan <- struct{}{}:
		}
		go func(v reflect.Value) {
			defer func() {
				<-limitChan
			}()
			ret := mF.Call([]reflect.Value{v})
			select {
			case <-exitChan:
				return
			case retChan <- ret:
			}
		}(v)
	}
	for i := 0; i < goroutineLimit; i++ {
		limitChan <- struct{}{}
	}
	close(retChan)
	return nil
}
