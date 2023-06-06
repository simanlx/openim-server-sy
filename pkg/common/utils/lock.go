package utils

import (
	"context"
	"crazy_server/pkg/common/db"
	"crazy_server/pkg/common/log"
	"errors"
	"sync"
	"time"
)

var mutex sync.Mutex

// Lock 加锁
func Lock(ctx context.Context, key string) error {
	mutex.Lock()
	defer mutex.Unlock()
	flag, err := db.DB.RDB.SetNX(ctx, key, time.Now().Unix(), 3*time.Second).Result()
	if err != nil {
		log.NewError("", "加锁失败：", key, err.Error())
		return err
	}
	if !flag {
		return errors.New("操作频繁,请稍后再试")
	}

	return nil
}

// UnLock 解锁
func UnLock(ctx context.Context, key string) int64 {
	num, err := db.DB.RDB.Del(ctx, key).Result()
	if err != nil {
		log.NewError("", "解锁失败：", key, err.Error())
		return 0
	}
	return num
}
