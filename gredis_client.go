package gredis

import (
	"fmt"
	"time"
)

import (
	redis_go "github.com/garyburd/redigo/redis"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Redis Client impl
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
type (
	redisClient struct {
		prefixKey string
		pool      *redis_go.Pool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Redis实例
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRedis(ip string, port int, password string, db, timeout int, prefixArgs ...string) IRedis {
	client := new(redisClient)

	if len(ip) == 0 {
		ip = "127.0.0.1"
	}

	if port <= 0 {
		port = 6379
	}

	addr := fmt.Sprintf("%s:%d", ip, port)

	prefix := ""
	if len(prefixArgs) > 0 {
		prefix = prefixArgs[0]
	}

	if prefix != "" {
		client.prefixKey = prefix
	}

	client.pool = newRedisPool(addr, password, db, timeout)

	return client
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Run Command
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) command(commandName string, args ...interface{}) (interface{}, error) {
	return s.do(commandName, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Run Command
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) do(commandName string, args ...interface{}) (interface{}, error) {
	redisPool := s.pool.Get()
	defer redisPool.Close()

	return redisPool.Do(commandName, args...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String Keys
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Keys(patternArgs ...string) ([]string, error) {
	keyPattern := "*"
	if len(patternArgs) > 0 {
		keyPattern = patternArgs[0]
	}
	return redis_go.Strings(s.command(REDIS_COMMAND_KEYS, keyPattern))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Key是否存在
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Exists(key string) (bool, error) {
	return redis_go.Bool(s.command(REDIS_COMMAND_EXISTS, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 更改key名称
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Rename(oldKey, newKey string) error {
	_, err := s.command(REDIS_COMMAND_RENAME, s.GetKey(oldKey), s.GetKey(newKey))
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 若newKey不存在，则更改oldKey名称
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) RenameNx(oldKey, newKey string) error {
	_, err := s.command(REDIS_COMMAND_RENAMENX, s.GetKey(oldKey), s.GetKey(newKey))
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除指定的Key
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Del(keyArgs ...string) error {
	keys := make([]interface{}, 0)

	for _, key := range keyArgs {
		keys = append(keys, s.GetKey(key))
	}

	_, err := s.command(REDIS_COMMAND_DEL, redis_go.Args{}.Add(keys...)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置过期时间（单位秒）
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Expire(key string, seconds int) error {
	_, err := s.command(REDIS_COMMAND_EXPIRE, s.GetKey(key), seconds)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置过期时间（单位毫秒）
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Pexpire(key string, milliseconds int) error {
	_, err := s.command(REDIS_COMMAND_PEXPIRE, s.GetKey(key), milliseconds)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置在指定时间戳之后键到期/过期(Unix秒时间戳格式)
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ExpireAt(key string, timestamp int) error {
	_, err := s.command(REDIS_COMMAND_EXPIREAT, s.GetKey(key), timestamp)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置在指定时间戳之后键到期/过期(Unix毫秒时间戳格式)
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) PexpireAt(key string, timestamp int) error {
	_, err := s.command(REDIS_COMMAND_PEXPIREAT, s.GetKey(key), timestamp)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 移除key的过期时间，key将持久保持
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Persist(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_PERSIST, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Ttl
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Ttl(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_TTL, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Pttl
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Pttl(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_PTTL, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 服务器信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Info() (string, error) {
	return redis_go.String(s.command(REDIS_COMMAND_INFO))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取指定key的存储类型
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Type(key string) (string, error) {
	return redis_go.String(s.command(REDIS_COMMAND_TYPE, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Dump指定key的数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Dump(key string) (string, error) {
	return redis_go.String(s.command(REDIS_COMMAND_DUMP, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SetData(structData interface{}, args ...interface{}) error {
	key := ""
	var time int = 0

	argsCount := len(args)
	if argsCount == 0 {
		key = s.GetObjectKey(structData)
	} else if argsCount == 1 {
		switch args[0].(type) {
		case string:
			key = args[0].(string)
			break
		case int:
			key = s.GetObjectKey(structData)
			time = args[0].(int)
			break
		}
	} else if argsCount == 2 {
		key = args[0].(string)
		if timeValue, isOk := args[1].(int); isOk {
			time = timeValue
		}
	}

	if jsonString, err := glib.ToJson(structData); err != nil {
		return err
	} else if len(jsonString) > 0 {
		if time == 0 {
			if err := s.Set(key, jsonString); err != nil {
				return err
			}
		} else {
			if err := s.Set(key, jsonString, time); err != nil {
				return err
			}
		}
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) GetData(structData interface{}, args ...interface{}) error {
	key := ""
	argsCount := len(args)
	if argsCount == 0 {
		key = s.GetObjectKey(structData)
	} else if argsCount == 1 {
		key = args[0].(string)
	}

	if data, err := s.Get(key); err != nil {
		return err
	} else {
		if err := glib.FromJson(string(data), &structData); err != nil {
			return err
		}
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * INCR | INCRBY
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Incr(key string, stepArgs ...int) (int, error) {
	var step int = 1
	if len(stepArgs) > 0 {
		step = stepArgs[0]
	}

	if step == 1 {
		return redis_go.Int(s.command(REDIS_COMMAND_INCR, s.GetKey(key)))
	}

	return redis_go.Int(s.command(REDIS_COMMAND_INCRBY, s.GetKey(key), step))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * INCRBYFLOAT
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) IncrByFloat(key string, value float64) (float64, error) {
	return redis_go.Float64(s.command(REDIS_COMMAND_INCRBYFLOAT, s.GetKey(key), value))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * DECR | DECRBY
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Decr(key string, stepArgs ...int) (int, error) {
	var step int = 1
	if len(stepArgs) > 0 {
		step = stepArgs[0]
	}

	if step == 1 {
		return redis_go.Int(s.command(REDIS_COMMAND_DECR, s.GetKey(key)))
	}

	return redis_go.Int(s.command(REDIS_COMMAND_DECRBY, s.GetKey(key), step))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String SET
 * key
 * value
 * ttl value
 * 0: second | 1: millisecond
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Set(key string, value interface{}, args ...int) error {
	var err error

	argsCount := len(args)
	if argsCount == 0 {
		_, err = s.command(REDIS_COMMAND_SET, s.GetKey(key), value)
	} else {
		ttlValue := args[0]
		isMilliSecond := false

		if argsCount > 1 {
			if args[0] == 1 {
				isMilliSecond = true
			}
		}

		if isMilliSecond {
			_, err = s.command(REDIS_COMMAND_PSETEX, s.GetKey(key), ttlValue, value)
		} else {
			_, err = s.command(REDIS_COMMAND_SETEX, s.GetKey(key), ttlValue, value)
		}

	}
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String SETNX
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SetNx(key string, value interface{}) error {
	_, err := s.command(REDIS_COMMAND_SETNX, s.GetKey(key), value)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String SETRANGE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SetRange(key string, index int, value interface{}) error {
	_, err := s.command(REDIS_COMMAND_SETRANGE, s.GetKey(key), index, value)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String Append
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Append(key string, value interface{}) error {
	_, err := s.command(REDIS_COMMAND_APPEND, s.GetKey(key), value)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String GET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Get(key string) ([]byte, error) {
	value, err := redis_go.String(s.command(REDIS_COMMAND_GET, s.GetKey(key)))
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String GETSET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) GetSet(key string, value interface{}) ([]byte, error) {
	value, err := redis_go.String(s.command(REDIS_COMMAND_GETSET, s.GetKey(key), value))
	if err != nil {
		return nil, err
	}

	return []byte(value.(string)), err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String GETRANGE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) GetRange(key string, start, end int) ([]byte, error) {
	value, err := redis_go.String(s.command(REDIS_COMMAND_GETRANGE, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end)...))
	if err != nil {
		return nil, err
	}

	return []byte(value), err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * String STRLEN
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) StrLen(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_STRLEN, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LPUSH
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LPush(key string, value ...interface{}) error {
	_, err := s.command(REDIS_COMMAND_LPUSH, redis_go.Args{}.Add(s.GetKey(key)).Add(value...)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List RPUSH
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) RPush(key string, value ...interface{}) error {
	_, err := s.command(REDIS_COMMAND_RPUSH, redis_go.Args{}.Add(s.GetKey(key)).Add(value...)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LPOP
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LPop(key string) (string, error) {
	return redis_go.String(s.command(REDIS_COMMAND_LPOP, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LRANGE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LRange(key string, start, end int) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_LRANGE, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LINDEX
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LIndex(key string, index int) (string, error) {
	return redis_go.String(s.command(REDIS_COMMAND_LINDEX, redis_go.Args{}.Add(s.GetKey(key)).Add(index)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LSET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LSet(key string, index int, value interface{}) error {
	_, err := s.command(REDIS_COMMAND_LSET, redis_go.Args{}.Add(s.GetKey(key)).Add(index).Add(value)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LREM
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LRem(key string, value interface{}, countArgs ...int) error {
	var count int = 1
	if len(countArgs) > 0 {
		count = countArgs[0]
	}
	_, err := s.command(REDIS_COMMAND_LREM, redis_go.Args{}.Add(s.GetKey(key)).Add(count).Add(value)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LTRIM
 * 修建列表，使其只保存指定范围内的数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LTrim(key string, start, end int) error {
	_, err := s.command(REDIS_COMMAND_LTRIM, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * List LLEN
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) LLen(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_LLEN, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HMSET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HSetData(structData interface{}, args ...interface{}) error {
	key := ""
	var time int = 0

	argsCount := len(args)
	if argsCount == 0 {
		key = s.GetObjectKey(structData)
	} else if argsCount == 1 {
		switch args[0].(type) {
		case string:
			key = args[0].(string)
			break
		case int:
			key = s.GetObjectKey(structData)
			time = args[0].(int)
			break
		}
	} else if argsCount == 2 {
		key = args[0].(string)
		if timeValue, isOk := args[1].(int); isOk {
			time = timeValue
		}
	}

	if _, err := s.command(REDIS_COMMAND_HMSET, redis_go.Args{}.Add(s.GetKey(key)).AddFlat(structData)...); err != nil {
		return err
	}

	if time > 0 {
		if err := s.Expire(key, time); err != nil {
			return err
		}
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HGETALL
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HGetData(structData interface{}, args ...interface{}) error {
	key := ""

	argsCount := len(args)
	if argsCount == 0 {
		key = s.GetObjectKey(structData)
	} else {
		key = args[0].(string)
	}

	data, err := redis_go.Values(s.command(REDIS_COMMAND_HGETALL, s.GetKey(key)))
	if err != nil {
		return err
	}

	if err := redis_go.ScanStruct(data, structData); err != nil {
		return err
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HSET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HSet(key, field string, value interface{}) error {
	_, err := s.command(REDIS_COMMAND_HSET, redis_go.Args{}.Add(s.GetKey(key)).Add(field).Add(value)...)

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HSETNX
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HSetNx(key, field string, value interface{}) error {
	_, err := s.command(REDIS_COMMAND_HSETNX, redis_go.Args{}.Add(s.GetKey(key)).Add(field).Add(value)...)

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HMSET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HMSet(key string, fields ...interface{}) error {
	_, err := s.command(REDIS_COMMAND_HMSET, redis_go.Args{}.Add(s.GetKey(key)).Add(fields...)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HGET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HGet(key string, field string) (string, error) {
	return redis_go.String(s.command(REDIS_COMMAND_HGET, redis_go.Args{}.Add(s.GetKey(key)).Add(field)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HMGET
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HMGet(key string, fields ...interface{}) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_HMGET, redis_go.Args{}.Add(s.GetKey(key)).Add(fields...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HKEYS
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HKeys(key string) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_HKEYS, redis_go.Args{}.Add(s.GetKey(key))...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HVALS
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HVals(key string) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_HVALS, redis_go.Args{}.Add(s.GetKey(key))...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HINCRBY
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HIncrBy(key string, field string, value int) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_HINCRBY, redis_go.Args{}.Add(s.GetKey(key)).Add(field).Add(value)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HINCRBYFLOAT
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HIncrByFloat(key string, field string, value float64) (float64, error) {
	return redis_go.Float64(s.command(REDIS_COMMAND_HINCRBYFLOAT, redis_go.Args{}.Add(s.GetKey(key)).Add(field).Add(value)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HLEN
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HLen(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_HLEN, redis_go.Args{}.Add(s.GetKey(key))...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HSTRLEN
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HStrLen(key, field string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_HSTRLEN, redis_go.Args{}.Add(s.GetKey(key)).Add(field)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HEXISTS
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HExists(key, field string) (bool, error) {
	return redis_go.Bool(s.command(REDIS_COMMAND_HEXISTS, redis_go.Args{}.Add(s.GetKey(key)).Add(field)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Hash HDEL
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) HDel(key string, fields ...interface{}) error {
	if _, err := s.command(REDIS_COMMAND_HDEL, redis_go.Args{}.Add(s.GetKey(key)).Add(fields...)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Add
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SAdd(key string, values ...interface{}) error {
	if _, err := s.command(REDIS_COMMAND_SADD, redis_go.Args{}.Add(s.GetKey(key)).Add(values...)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Move
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SMove(srcKkey, desKey string, values ...interface{}) error {
	if _, err := s.command(REDIS_COMMAND_SMOVE, redis_go.Args{}.Add(s.GetKey(srcKkey)).Add(s.GetKey(desKey)).Add(values...)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Pop
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SPop(key string, countArgs ...int) ([]string, error) {
	var count int = 1
	if len(countArgs) > 0 {
		count = countArgs[0]
	}
	return redis_go.Strings(s.command(REDIS_COMMAND_SPOP, s.GetKey(key), count))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Rem
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SRem(key string, values ...interface{}) error {
	if _, err := s.command(REDIS_COMMAND_SREM, redis_go.Args{}.Add(s.GetKey(key)).Add(values...)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Card
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SCard(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_SCARD, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set IsMemeber
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SIsMemeber(key string, value interface{}) (bool, error) {
	return redis_go.Bool(s.command(REDIS_COMMAND_SISMEMBER, redis_go.Args{}.Add(s.GetKey(key)).Add(value)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Members
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SMembers(key string) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_SMEMBERS, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Members
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SMembersInt(key string) ([]int, error) {
	return redis_go.Ints(s.command(REDIS_COMMAND_SMEMBERS, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Members
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SMembersInt64(key string) ([]int64, error) {
	return redis_go.Int64s(s.command(REDIS_COMMAND_SMEMBERS, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Members
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SMembersFloat64(key string) ([]float64, error) {
	return redis_go.Float64s(s.command(REDIS_COMMAND_SMEMBERS, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set RandMember
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SRandMembers(key string, countArgs ...int) ([]string, error) {
	var count int = 1
	if len(countArgs) > 0 {
		count = countArgs[0]
	}
	return redis_go.Strings(s.command(REDIS_COMMAND_SRANDMEMBER, s.GetKey(key), count))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set RandMember
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SRandMembersInt(key string, countArgs ...int) ([]int, error) {
	var count int = 1
	if len(countArgs) > 0 {
		count = countArgs[0]
	}
	return redis_go.Ints(s.command(REDIS_COMMAND_SRANDMEMBER, s.GetKey(key), count))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set RandMember
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SRandMembersInt64(key string, countArgs ...int) ([]int64, error) {
	var count int = 1
	if len(countArgs) > 0 {
		count = countArgs[0]
	}
	return redis_go.Int64s(s.command(REDIS_COMMAND_SRANDMEMBER, s.GetKey(key), count))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Union ints
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SUnionInt(keys ...string) ([]int, error) {
	newKeys := make([]interface{}, 0)

	for _, key := range keys {
		newKeys = append(newKeys, s.GetKey(key))
	}

	return redis_go.Ints(s.command(REDIS_COMMAND_SUNION, redis_go.Args{}.Add(newKeys...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Union strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SUnion(keys ...string) ([]string, error) {
	newKeys := make([]interface{}, 0)

	for _, key := range keys {
		newKeys = append(newKeys, s.GetKey(key))
	}

	return redis_go.Strings(s.command(REDIS_COMMAND_SUNION, redis_go.Args{}.Add(newKeys...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Inter ints
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SInterInt(keys ...string) ([]int, error) {
	newKeys := make([]interface{}, 0)

	for _, key := range keys {
		newKeys = append(newKeys, s.GetKey(key))
	}

	return redis_go.Ints(s.command(REDIS_COMMAND_SINTER, redis_go.Args{}.Add(newKeys...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Inter strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SInter(keys ...string) ([]string, error) {
	newKeys := make([]interface{}, 0)

	for _, key := range keys {
		newKeys = append(newKeys, s.GetKey(key))
	}

	return redis_go.Strings(s.command(REDIS_COMMAND_SINTER, redis_go.Args{}.Add(newKeys...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Diff ints
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SDiffInt(keys ...string) ([]int, error) {
	newKeys := make([]interface{}, 0)

	for _, key := range keys {
		newKeys = append(newKeys, s.GetKey(key))
	}

	return redis_go.Ints(s.command(REDIS_COMMAND_SDIFF, redis_go.Args{}.Add(newKeys...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Set Diff strings
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SDiff(keys ...string) ([]string, error) {
	newKeys := make([]interface{}, 0)

	for _, key := range keys {
		newKeys = append(newKeys, s.GetKey(key))
	}

	return redis_go.Strings(s.command(REDIS_COMMAND_SDIFF, redis_go.Args{}.Add(newKeys...)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZADD
 * score int, value interface{}, score int, value interface{} ...
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZAdd(key string, members ...interface{}) error {
	_, err := s.command(REDIS_COMMAND_ZADD, redis_go.Args{}.Add(s.GetKey(key)).Add(members...)...)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZRANGE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRange(key string, start, end int) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_ZRANGE, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZRANGE WITHSCORES
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRangeWithScore(key string, start, end int) (map[string]int, error) {
	return redis_go.IntMap(s.command(REDIS_COMMAND_ZRANGE, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end).Add("WITHSCORES")...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZRANGEBYSCORE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRangeByScore(key string, min, max interface{}, limitArgs ...int) ([]string, error) {
	var offset int = 0
	var count int = -1
	if len(limitArgs) > 0 {
		offset = limitArgs[0]
	}
	if len(limitArgs) > 1 {
		count = limitArgs[1]
	}

	args := redis_go.Args{}.Add(s.GetKey(key)).Add(min).Add(max).Add("LIMIT").Add(offset).Add(count)

	return redis_go.Strings(s.command(REDIS_COMMAND_ZRANGEBYSCORE, args...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZREVRANGE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRevRange(key string, start, end int) ([]string, error) {
	return redis_go.Strings(s.command(REDIS_COMMAND_ZREVRANGE, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZREVRANGEBYSCORE
 * key, min, max, offset, count
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRevRangeByScore(key string, min, max interface{}, limitArgs ...int) ([]string, error) {
	var offset int = 0
	var count int = -1
	if len(limitArgs) > 0 {
		offset = limitArgs[0]
	}
	if len(limitArgs) > 1 {
		count = limitArgs[1]
	}

	args := redis_go.Args{}.Add(s.GetKey(key)).Add(max).Add(min).Add("LIMIT").Add(offset).Add(count)

	return redis_go.Strings(s.command(REDIS_COMMAND_ZREVRANGEBYSCORE, args...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZREM
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRem(key string, members ...interface{}) error {
	if _, err := s.command(REDIS_COMMAND_ZREM, redis_go.Args{}.Add(s.GetKey(key)).Add(members...)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZREMRANGEBYSCORE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRemRangeByScore(key string, min, max interface{}) error {
	if _, err := s.command(REDIS_COMMAND_ZREMRANGEBYSCORE, redis_go.Args{}.Add(s.GetKey(key)).Add(min).Add(max)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZREMRANGEBYRANK
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRemRangeByRank(key string, start, end int) error {
	if _, err := s.command(REDIS_COMMAND_ZREMRANGEBYRANK, redis_go.Args{}.Add(s.GetKey(key)).Add(start).Add(end)...); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZCARD
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZCard(key string) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_ZCARD, s.GetKey(key)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZSCORE
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZScore(key string, member interface{}) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_ZSCORE, s.GetKey(key), member))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZRank
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRank(key string, member interface{}) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_ZRANK, redis_go.Args{}.Add(s.GetKey(key)).Add(member)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZREVRANK
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZRevRank(key string, member interface{}) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_ZREVRANK, redis_go.Args{}.Add(s.GetKey(key)).Add(member)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ZSET ZCOUNT
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) ZCount(key string, min, max interface{}) (int, error) {
	return redis_go.Int(s.command(REDIS_COMMAND_ZCOUNT, redis_go.Args{}.Add(s.GetKey(key)).Add(min).Add(max)...))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Pipeline MULTI and EXEC
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) Pipeline(commands []map[string][]interface{}, watchKeys ...interface{}) (interface{}, error) {
	redisPool := s.pool.Get()
	defer redisPool.Close()

	if len(watchKeys) > 0 {
		redisPool.Send(REDIS_COMMAND_WATCH, watchKeys...)
	}

	redisPool.Send(REDIS_COMMAND_MULTI)

	for _, commandLine := range commands {
		for cmd, args := range commandLine {
			redisPool.Send(cmd, args...)
		}
	}

	return redisPool.Do(REDIS_COMMAND_EXEC)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 选择指定数据库
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) SelectDb(index int) error {
	if _, err := s.command(REDIS_COMMAND_SELECTDB, index); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 保存数据到磁盘
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) BgSave() error {
	if _, err := s.command(REDIS_COMMAND_BGSAVE); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 清理指定数据库的数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) FlushDb(index int) error {
	if _, err := s.command(REDIS_COMMAND_FLUSHDB, index); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 清理全部数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) FlushAll() error {
	if _, err := s.command(REDIS_COMMAND_FLUSHALL); err != nil {
		return err
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * redigo帮助方法包装
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
/*func (s *redisClient) Bool(reply interface{}, err error) (bool, error) {
	return redis_go.Bool(reply, err)
}

func (s *redisClient) ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redis_go.ByteSlices(reply, err)
}

func (s *redisClient) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis_go.Bytes(reply, err)
}

func (s *redisClient) Float64(reply interface{}, err error) (float64, error) {
	return redis_go.Float64(reply, err)
}

func (s *redisClient) Int(reply interface{}, err error) (int, error) {
	return redis_go.Int(reply, err)
}

func (s *redisClient) Int64(reply interface{}, err error) (int64, error) {
	return redis_go.Int64(reply, err)
}

func (s *redisClient) Int64Map(result interface{}, err error) (map[string]int64, error) {
	return redis_go.Int64Map(result, err)
}

func (s *redisClient) IntMap(result interface{}, err error) (map[string]int, error) {
	return redis_go.IntMap(result, err)
}

func (s *redisClient) Ints(reply interface{}, err error) ([]int, error) {
	return redis_go.Ints(reply, err)
}

func (s *redisClient) MultiBulk(reply interface{}, err error) ([]interface{}, error) {
	return redis_go.MultiBulk(reply, err)
}

func (s *redisClient) Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	return redis_go.Scan(src, dest...)
}

func (s *redisClient) ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error {
	return redis_go.ScanSlice(src, dest)
}

func (s *redisClient) ScanStruct(src []interface{}, dest interface{}) error {
	return redis_go.ScanStruct(src, dest)
}

func (s *redisClient) String(reply interface{}, err error) (string, error) {
	return redis_go.String(reply, err)
}

func (s *redisClient) StringMap(result interface{}, err error) (map[string]string, error) {
	return redis_go.StringMap(result, err)
}

func (s *redisClient) Strings(reply interface{}, err error) ([]string, error) {
	return redis_go.Strings(reply, err)
}

func (s *redisClient) Uint64(reply interface{}, err error) (uint64, error) {
	return redis_go.Uint64(reply, err)
}

func (s *redisClient) Values(reply interface{}, err error) ([]interface{}, error) {
	return redis_go.Values(reply, err)
}*/

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取最终的Key
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) GetKey(key string) string {
	key = fmt.Sprintf("%s%s", s.prefixKey, key)
	return key
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取对象Key
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisClient) GetObjectKey(model interface{}) string {
	key := ""

	if pgk, fieldValue, err := glib.GetStructFieldValueByName(model, "Id"); err == nil {
		if fieldValue, isOk := fieldValue.(int64); isOk {
			key = fmt.Sprintf("%s%d", pgk, fieldValue)
		}
	}

	return key
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取 RedisPool
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func newRedisPool(address, password string, db, timeout int) *redis_go.Pool {
	return &redis_go.Pool{
		MaxIdle:     8,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis_go.Conn, error) {
			return dial(address, password, db, timeout)
		},
		TestOnBorrow: func(conn redis_go.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do(REDIS_COMMAND_PING)
			return err
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 链接 Redis 服务器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func dial_old(address, password string, timeout int) (redis_go.Conn, error) {
	conn, err := redis_go.DialTimeout(
		"tcp",
		address,
		time.Duration(timeout)*time.Second,
		time.Duration(timeout)*time.Second,
		time.Duration(timeout)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	if len(password) > 0 {
		if _, err := conn.Do(REDIS_COMMAND_AUTH, password); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 链接 Redis 服务器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func dial(address, password string, db, timeout int) (redis_go.Conn, error) {
	conn, err := redis_go.Dial(
		"tcp",
		address,
		redis_go.DialConnectTimeout(time.Duration(timeout)*time.Second),
		redis_go.DialReadTimeout(time.Duration(timeout)*time.Second),
		redis_go.DialWriteTimeout(time.Duration(timeout)*time.Second),
		redis_go.DialPassword(password),
		redis_go.DialDatabase(db),
	)

	if err != nil {
		return nil, err
	}

	return conn, err
}
