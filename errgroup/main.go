package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"sync/atomic"
	"time"
)

const (
	nWorkers = 10
)

var (
	//好友数组
	users        = []string{"aa", "bb", "cc", "dd", "ee", "ff"}
	errorEnd     = errors.New("query user end")
	errorUnknown = errors.New("unknown err")
)

type User struct {
	Id   int64
	Name string
}

//朋友关系迭代器
type FriendIterator struct {
	index int64
}

//获取下一个朋友id
func (it *FriendIterator) Next(ctx context.Context) (int64, error) {
	if it.index >= int64(len(users)) {
		return 0, errorEnd
	}

	query := time.After(10 * time.Millisecond) //模拟查询耗时
	//模拟查询出错
	/*if it.index == 3 {
		return it.index, errorUnknown
	}*/

	// 10ms返回数据
	select {
	case <-ctx.Done():
		return it.index, ctx.Err()
	case <-query:
		r := it.index
		it.index++
		fmt.Printf("found friend id: %d\n", r)
		return r, nil
	}
}

func GetFriendIds(user int64) *FriendIterator {
	return &FriendIterator{}
}

//获取用户的信息
func GetUserProfile(ctx context.Context, id int64) (*User, error) {
	// 100ms返回数据
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
		if id < 0 || id >= int64(len(users)) {
			return nil, fmt.Errorf("unknown user: %d", id)
		}
		fmt.Printf("found user profile: %d\n", id)
		return &User{Id: id, Name: users[id]}, nil
	}
}

//假设现在有两个服务，一个是朋友关系服务，一个是用户信息服务。
//现在有个接口要获取某个用户的朋友信息
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	start := time.Now()
	//rsp, err := demo1(ctx, 0)
	//rsp, err := demo2(ctx,0)
	//rsp, err := demo3(ctx, 0)
	rsp, err := demo4(ctx, 0)
	if err != nil {
		fmt.Printf("error: %s", err)
	} else {
		fmt.Printf("finished in %s\n ret: %s", time.Since(start), jsonString(rsp))
	}
	fmt.Println()
	cancel()
}

func jsonString(v interface{}) string {
	if ret, err := json.Marshal(v); err != nil {
		return err.Error()
	} else {
		return string(ret)
	}
}

//串行操作  先查询朋友ID再查询朋友信息
func demo1(ctx context.Context, user int64) (map[string]*User, error) {
	// 获取朋友id
	var friendIds []int64
	for it := GetFriendIds(user); ; {
		if id, err := it.Next(ctx); err != nil {
			if err == errorEnd {
				break
			}
			return nil, fmt.Errorf("GetFriendIds %d: %s", id, err)
		} else {
			friendIds = append(friendIds, id)
		}
	}

	// 查询朋友的信息
	ret := map[string]*User{}
	for _, friendId := range friendIds {
		if friend, err := GetUserProfile(ctx, friendId); err != nil {
			return nil, fmt.Errorf("GetUserProfile %d: %s", friendId, err)
		} else {
			ret[friend.Name] = friend
		}
	}
	return ret, nil
}

//并行处理
func demo2(ctx context.Context, user int64) (map[string]*User, error) {
	friendIds := make(chan int64)

	// 有一个Goroutine去获取朋友id
	go func() {
		defer close(friendIds)
		for it := GetFriendIds(user); ; {
			if id, err := it.Next(ctx); err != nil {
				if err == errorEnd {
					break
				}
				//如果出错了，该怎么处理err呢
				log.Fatalf("GetFriendIds %d: %s", id, err)
			} else {
				friendIds <- id
			}
		}
	}()

	friends := make(chan *User)

	// 有n个Goroutine去获取朋友信息
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func() {
			defer func() {
				// 关闭channel
				if atomic.AddInt32(&workers, -1) == 0 {
					close(friends)
				}
			}()

			for id := range friendIds {
				if friend, err := GetUserProfile(ctx, id); err != nil {
					//如果出错了，该怎么处理err呢
					log.Fatalf("GetUserProfile %d: %s", id, err)
				} else {
					friends <- friend
				}
			}
		}()
	}

	//封装返回信息
	ret := map[string]*User{}
	for friend := range friends {
		ret[friend.Name] = friend
	}

	return ret, nil
}

//使用errgroup来并行处理
func demo3(ctx context.Context, user int64) (map[string]*User, error) {
	g := errgroup.Group{}
	friendIds := make(chan int64)

	// 有一个Goroutine去获取朋友id
	g.Go(func() error {
		defer close(friendIds)
		for it := GetFriendIds(user); ; {
			if id, err := it.Next(ctx); err != nil {
				if err == errorEnd {
					return nil
				}
				return fmt.Errorf("GetFriendIds %d: %s", id, err)
			} else {
				friendIds <- id
			}
		}
	})

	friends := make(chan *User)

	//有n个Goroutine去获取朋友信息
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// 关闭channel
				if atomic.AddInt32(&workers, -1) == 0 {
					close(friends)
				}
			}()

			for id := range friendIds {
				if friend, err := GetUserProfile(ctx, id); err != nil {
					return fmt.Errorf("GetUserProfile %d: %s", id, err)
				} else {
					friends <- friend
				}
			}
			return nil
		})
	}

	// 还有一个Goroutine去封装返回信息
	ret := map[string]*User{}
	g.Go(func() error {
		for friend := range friends {
			if friend.Id == 3 {
				//return errorUnknown
			}
			ret[friend.Name] = friend
		}
		return nil
	})

	// 返回数据 g.Wait()会返回Goroutine执行出错的第一个err
	return ret, g.Wait()
}

func demo4(ctx context.Context, user int64) (map[string]*User, error) {
	g, ctx := errgroup.WithContext(ctx)
	friendIds := make(chan int64)

	//有一个Goroutine去获取朋友id
	g.Go(func() error {
		defer close(friendIds)
		for it := GetFriendIds(user); ; {
			if id, err := it.Next(ctx); err != nil {
				if err == errorEnd {
					return nil
				}
				return fmt.Errorf("GetFriendIds %d: %s", id, err)
			} else {
				select { //使用select来退出Goroutine
				case <-ctx.Done():
					return ctx.Err()
				case friendIds <- id:
				}
			}
		}
	})

	friends := make(chan *User)

	//有n个Goroutine去获取朋友信息
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// 关闭channel
				if atomic.AddInt32(&workers, -1) == 0 {
					close(friends)
				}
			}()

			for id := range friendIds {
				if friend, err := GetUserProfile(ctx, id); err != nil {
					return fmt.Errorf("GetUserProfile %d: %s", id, err)
				} else {
					select { //使用select来退出Goroutine
					case <-ctx.Done():
						return ctx.Err()
					case friends <- friend:
					}
				}
			}
			return nil
		})
	}

	//还有一个Goroutine去封装返回信息
	ret := map[string]*User{}
	g.Go(func() error {
		for friend := range friends {
			if friend.Id == 3 {
				return errorUnknown
			}
			ret[friend.Name] = friend
		}
		return nil
	})

	//返回数据 g.Wait()会返回Goroutine执行出错的第一个err
	return ret, g.Wait()
}
