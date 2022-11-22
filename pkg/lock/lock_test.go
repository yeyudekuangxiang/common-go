package lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func redisLock() *RedisDistributedLock {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if rdb.Ping(context.Background()).Err() != nil {
		log.Println("redis not collection skip lock test")
		return nil
	}
	lock := RedisDistributedLock{
		Redis:  rdb,
		Prefix: "test",
	}
	return &lock
}
func TestLock(t *testing.T) {
	lock := redisLock()
	if lock == nil {
		return
	}
	lock.UnLock("unit_lock_test")
	var lockNum int64
	timer := time.NewTimer(time.Millisecond * 5000)
	for j := 0; j < 10000; j++ {
		go func() {
			if lock.Lock("unit_lock_test", time.Second*8) {
				atomic.AddInt64(&lockNum, 1)
			}
		}()
	}
	<-timer.C
	assert.Equal(t, int64(1), lockNum)
	lock.UnLock("unit_lock_test")
}
func TestDeleteLock(t *testing.T) {
	lock := redisLock()
	if lock == nil {
		return
	}
	lock.UnLock("unit_lock_test_delete_key")
	lock.Lock("unit_lock_test_delete_key", time.Second*10)
	lock.UnLock("unit_lock_test_delete_key")
	assert.Equal(t, true, lock.Lock("unit_lock_test_delete_key", time.Second*10))
	lock.UnLock("unit_lock_test_delete_key")
}
func TestLockWait(t *testing.T) {
	lock := redisLock()
	if lock == nil {
		return
	}
	lock.UnLock("unit_lock_test_lock_wait")
	timer := time.NewTimer(time.Second * 8)
	go func() {
		t2 := time.Now()
		time.Sleep(time.Millisecond * 500)
		for i := 0; i < 50; i++ {
			go func() {
				lock.LockWait("unit_lock_test_lock_wait", time.Second*5)
				assert.Equal(t, false, time.Now().Sub(t2) <= time.Second*4)
			}()
		}
	}()
	lock.LockWait("unit_lock_test_lock_wait", time.Second*5)
	<-timer.C
	assert.Equal(t, true, true)
	lock.UnLock("unit_lock_test_lock_wait")
}
func TestLockNum(t *testing.T) {
	lock := redisLock()
	if lock == nil {
		return
	}
	lock.UnLock("unit_lock_test_lock_num")
	for i := 0; i < 5; i++ {
		assert.Equal(t, true, lock.LockNum("unit_lock_test_lock_num", 5, time.Second*5))
	}
	wait := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			assert.Equal(t, false, lock.LockNum("unit_lock_test_lock_num", 5, time.Second*5))
		}()
	}
	wait.Wait()
	lock.UnLock("unit_lock_test_lock_num")
}
func TestLockNumWait(t *testing.T) {
	lock := redisLock()
	if lock == nil {
		return
	}

	lock.UnLock("unit_lock_test_lock_num_wait")
	var num int64
	for i := 0; i < 1000; i++ {
		go func() {
			lock.LockNumWait("unit_lock_test_lock_num_wait", 5, time.Second*5)
			atomic.AddInt64(&num, 1)
		}()
	}

	timer := time.NewTimer(time.Second * 9)
	<-timer.C
	assert.Equal(t, int64(10), num)
	lock.UnLock("unit_lock_test_lock_num_wait")
}
