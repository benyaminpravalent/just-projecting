package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	radix "github.com/mediocregopher/radix/v3"
	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/logger"
	"github.com/sirupsen/logrus"
)

type redisCache struct {
	client radix.Client
	logger *logrus.Entry
}

// NewRedis creates new redis connection and returns redis functions.
func NewRedis(ctx context.Context) (Cache, error) {
	log := logger.GetLoggerContext(ctx, "cache", "NewRedis")

	jsonByte, err := json.Marshal(config.Get("redis"))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var cfg RedisConfig
	err = json.Unmarshal(jsonByte, &cfg)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var conn radix.Client
	var opts []radix.DialOpt
	var customConnFunc func(network, addr string) (radix.Conn, error)
	//Default Pool Size
	poolSize := 5

	// This is a ConnFunc which will set up a connection which is authenticated
	// and has a timeout on all operations
	if cfg.PoolSize != 0 {
		poolSize = cfg.PoolSize
	}
	if cfg.Timeout != 0 {
		opts = append(opts, radix.DialTimeout(time.Duration(cfg.Timeout)*time.Second))
	}
	if cfg.AuthPass != "" {
		opts = append(opts, radix.DialAuthPass(cfg.AuthPass))
	}
	customConnFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr, opts...)
	}

	var popt []radix.PoolOpt
	popt = append(popt, radix.PoolConnFunc(customConnFunc))

	// this pool will use our ConnFunc for all connections it creates.
	conn, err = radix.NewPool("tcp", cfg.Server, poolSize, popt...)
	if err != nil {
		log.Fatalf("failed creating new redis pool [%v]", err)
		return nil, err
	}

	log.Info("Created new connection to redis")

	c := &redisCache{
		client: conn,
		logger: logger.GetLogger("cache", "redisFunc"),
	}
	return c, err
}

// Get the item with the provided key.
// Return nil byte if the item didn't already exist in the cache.
func (m *redisCache) Get(key string) (rcv []byte, err error) {
	err = m.client.Do(radix.Cmd(&rcv, "GET", key))
	if err != nil {
		m.logger.Error(fmt.Sprintf("%s %s %s", key, string(rcv), err.Error()))
		return
	}
	return
}

// Set writes the given item, unconditionally.
func (m *redisCache) Set(key string, val []byte, expiration time.Duration) (err error) {

	args := []string{key, string(val)}

	if expiration != 0 {
		//EX seconds -- Set the specified expire time, in seconds.
		//PX milliseconds -- Set the specified expire time, in milliseconds.
		args = append(args, "EX", fmt.Sprintf("%d", int(expiration.Seconds())))
	}

	err = m.client.Do(radix.Cmd(nil, "SET", args...))
	if err != nil {
		m.logger.Error(fmt.Sprintf("%s %s %s", key, string(val), err.Error()))
		return
	}

	return
}

// Delete deletes the item with the provided key.
// return nil error if the item didn't already exist in the cache.
func (m *redisCache) Delete(key string) (err error) {
	err = m.client.Do(radix.Cmd(nil, "DEL", key))
	if err != nil {
		m.logger.Error(fmt.Sprintf("%s %s", key, err.Error()))
		return
	}
	return
}

// Incr the item with the provided key.
// Return incremented byte if the item didn't already exist in the cache.
func (m *redisCache) Incr(key string) (rcv []byte, err error) {
	err = m.client.Do(radix.Cmd(&rcv, "INCR", key))
	if err != nil {
		m.logger.Error(fmt.Sprintf("%s %s %s", key, string(rcv), err.Error()))
		return
	}
	return
}
