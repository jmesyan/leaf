package network

import (
	"errors"
	"fmt"
	"github.com/jmesyan/leaf/log"
	_ "os"
	"sync"
	"time"
)

type AsyncResult struct {
	key    uint32
	block  bool
	callback func([]interface{})
	result chan []interface{}
}

func NewAsyncResult(key uint32, block bool, callback func([]interface{})) *AsyncResult {
	asyncResult := new(AsyncResult)
	asyncResult.key = key;
	asyncResult.block = block;
	if block{
		asyncResult.result = make(chan []interface{}, 1)
	} else {
		asyncResult.callback = callback
	}
	return asyncResult
}

func (this *AsyncResult) GetKey() uint32 {
	return this.key
}

func (this *AsyncResult) SetResult(data []interface{}) {
	this.result <- data
}

func (this *AsyncResult) SetCallback(data []interface{}){
	if this.callback != nil{
		this.callback(data)
	}
}

func (this *AsyncResult) GetResult(timeout time.Duration) ([]interface{}, error) {
	select {
	case <-time.After(timeout):
		log.Error(fmt.Sprintf("GetResult AsyncResult: timeout %s", this.key))
		close(this.result)
		return nil, errors.New(fmt.Sprintf("GetResult AsyncResult: timeout %s", this.key))
	case result := <-this.result:
		return result, nil
	}
	return nil, errors.New("GetResult AsyncResult error. reason: no")
}

type AsyncResultMgr struct {
	ticker   uint32
	results map[uint32]*AsyncResult
	sync.RWMutex
}

func NewAsyncResultMgr() *AsyncResultMgr {
	return &AsyncResultMgr{
		results: make(map[uint32]*AsyncResult, 0),
		ticker:   0,
	}
}

func (this *AsyncResultMgr) Add(block bool, callback func([]interface{})) (*AsyncResult, error) {
	this.Lock()
	defer this.Unlock()
	if block && callback != nil{
		return nil, errors.New("block result can't set callback function")
	}

	if callback == nil {
		return nil, errors.New("tiker need callback")
	}
	if this.ticker > 65535{
		this.ticker = 0;
	}
	this.ticker++
	var ncallback func([]interface{})
	if callback != nil{
		ncallback = func(data []interface{}){
			this.Remove(this.ticker)
			callback(data)
		}
	}
	r := NewAsyncResult(this.ticker, block, ncallback)
	this.results[r.GetKey()] = r
	return r, nil
}

func (this *AsyncResultMgr) Remove(key uint32) {
	this.Lock()
	defer this.Unlock()

	delete(this.results, key)
}

func (this *AsyncResultMgr) GetAsyncResult(key uint32) (*AsyncResult, error) {
	this.RLock()
	defer this.RUnlock()

	r, ok := this.results[key]
	if ok {
		return r, nil
	} else {
		return nil, errors.New("not found AsyncResult")
	}
}

func (this *AsyncResultMgr) FillAsyncResult(key uint32, data []interface{}) error {
	r, err := this.GetAsyncResult(key)
	if err == nil {
		this.Remove(key)
		if r.block {
			r.SetResult(data)
		} else {
			r.SetCallback(data)
		}
		return nil
	} else {
		return err
	}
}
