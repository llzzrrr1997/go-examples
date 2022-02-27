package main

import (
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

var errorNotExist = errors.New("not exist")
var errorTimeOut = errors.New("time out")
var g singleflight.Group

func main() {
	//demo1()
	//demo2()
	demo3()
}

func demo1() {
	var wg sync.WaitGroup
	wg.Add(10)

	//模拟10个请求
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			//对同一个key进行查询数据
			data, _ := queryData("key")
			fmt.Printf("queryData success, ret:%s\n", data)
		}()
	}
	wg.Wait()
}

//获取数据
func queryData(key string) (string, error) {
	//从缓存中查询
	data, err := getFromCache(key)
	if err == errorNotExist { //数据不存在的err
		//从数据库中查询
		//data := getFromDB(key)
		//使用单飞模式
		v, err, _ := g.Do(key, func() (interface{}, error) {
			dbData := getFromDB(key)
			//写到缓存
			return dbData, nil
		})
		if err != nil {
			fmt.Printf("g.Do err:%s", err)
			return "", err
		}
		data = v.(string)
		return data, nil
	}
	if err != nil {
		return "", err
	}
	return data, nil
}

//从缓存中查询
func getFromCache(key string) (string, error) {
	return "", errorNotExist
}

//从数据库中查询
func getFromDB(key string) string {
	time.Sleep(1 * time.Second) //查询耗时1秒
	fmt.Printf("getFromDB key:%s\n", key)
	return "data"
}

func demo2() {
	var wg sync.WaitGroup
	wg.Add(10)

	//模拟10个请求
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			//对同一个key进行查询数据
			data, _ := forgetData("key")
			fmt.Printf("queryData success, ret:%s\n", data)
		}()
		//每100ms查询1次
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
}

//Forget方法使用
func forgetData(key string) (string, error) {
	//从缓存中查询
	data, err := getFromCache(key)
	if err == errorNotExist { //数据不存在的err
		v, err, _ := g.Do(key, func() (interface{}, error) {
			//只共享500ms内的请求结果
			go func() {
				time.Sleep(500 * time.Millisecond)
				g.Forget(key)
				fmt.Printf("forget a key\n")
			}()
			dbData := getFromDB(key)
			return dbData, nil
		})
		if err != nil {
			fmt.Printf("g.Do err:%s", err)
			return "", err
		}
		data = v.(string)
		return data, nil
	}
	if err != nil {
		return "", err
	}
	return data, nil
}

func demo3() {
	var wg sync.WaitGroup
	wg.Add(10)

	//模拟10个请求
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			//对同一个key进行查询数据
			data, err := timeOutData("key")
			if err != nil {
				fmt.Printf("queryData failed, ret:%s, err:%s\n", data, err)
				return
			}
			fmt.Printf("queryData success, ret:%s\n", data)
		}()
	}
	wg.Wait()
}

//DoChan使用
func timeOutData(key string) (string, error) {
	//从缓存中查询
	data, err := getFromCache(key)
	if err == errorNotExist { //数据不存在的err
		retCh := g.DoChan(key, func() (interface{}, error) {
			dbData := getFromDB(key)
			return dbData, nil
		})
		//500ms就超时
		timeOut := time.After(500 * time.Millisecond)
		select {
		case <-timeOut: //超时
			return "", errorTimeOut
		case ret := <-retCh: //查询到数据
			return ret.Val.(string), nil
		}
	}
	if err != nil {
		return "", err
	}
	return data, nil
}
