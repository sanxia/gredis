package gredis

/* ================================================================================
 * redis client interface
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

type (
	IRedis interface {
		Keys(patternArgs ...string) ([]string, error)
		Exists(key string) (bool, error)
		Rename(oldKey, newKey string) error
		RenameNx(oldKey, newKey string) error
		Del(keyArgs ...string) error
		Expire(key string, seconds int) error
		Pexpire(key string, milliseconds int) error
		ExpireAt(key string, timestamp int) error
		PexpireAt(key string, timestamp int) error
		Persist(key string) (int, error)
		Ttl(key string) (int, error)
		Pttl(key string) (int, error)
		Info() (string, error)
		Type(key string) (string, error)
		Dump(key string) (string, error)

		SetData(structData interface{}, args ...interface{}) error
		GetData(structData interface{}, args ...interface{}) error
		Incr(key string, stepArgs ...int) (int, error)
		IncrByFloat(key string, value float64) (float64, error)
		Decr(key string, stepArgs ...int) (int, error)
		Set(key string, value interface{}, args ...int) error
		SetNx(key string, value interface{}) error
		SetRange(key string, index int, value interface{}) error
		Append(key string, value interface{}) error
		Get(key string) ([]byte, error)
		GetSet(key string, value interface{}) ([]byte, error)
		GetRange(key string, start, end int) ([]byte, error)
		StrLen(key string) (int, error)

		LPush(key string, value ...interface{}) error
		RPush(key string, value ...interface{}) error
		LPop(key string) (string, error)
		LRange(key string, start, end int) ([]string, error)
		LIndex(key string, index int) (string, error)
		LSet(key string, index int, value interface{}) error
		LRem(key string, value interface{}, countArgs ...int) error
		LTrim(key string, start, end int) error
		LLen(key string) (int, error)

		HSetData(structData interface{}, args ...interface{}) error
		HGetData(structData interface{}, args ...interface{}) error
		HSet(key, field string, value interface{}) error
		HSetNx(key, field string, value interface{}) error
		HMSet(key string, fields ...interface{}) error
		HGet(key string, field string) (string, error)
		HMGet(key string, fields ...interface{}) ([]string, error)
		HKeys(key string) ([]string, error)
		HVals(key string) ([]string, error)
		HIncrBy(key string, field string, value int) (int, error)
		HIncrByFloat(key string, field string, value float64) (float64, error)
		HLen(key string) (int, error)
		HStrLen(key, field string) (int, error)
		HExists(key, field string) (bool, error)
		HDel(key string, fields ...interface{}) error

		SAdd(key string, values ...interface{}) error
		SMove(srcKkey, desKey string, values ...interface{}) error
		SPop(key string, countArgs ...int) ([]string, error)
		SRem(key string, values ...interface{}) error
		SCard(key string) (int, error)
		SIsMemeber(key string, value interface{}) (bool, error)
		SMembers(key string) ([]string, error)
		SMembersInt(key string) ([]int, error)
		SMembersInt64(key string) ([]int64, error)
		SMembersFloat64(key string) ([]float64, error)
		SRandMembers(key string, countArgs ...int) ([]string, error)
		SRandMembersInt(key string, countArgs ...int) ([]int, error)
		SRandMembersInt64(key string, countArgs ...int) ([]int64, error)
		SUnion(keys ...string) ([]string, error)
		SUnionInt(keys ...string) ([]int, error)
		SInter(keys ...string) ([]string, error)
		SInterInt(keys ...string) ([]int, error)
		SDiff(keys ...string) ([]string, error)
		SDiffInt(keys ...string) ([]int, error)

		ZAdd(key string, members ...interface{}) error
		ZRange(key string, start, end int) ([]string, error)
		ZRangeWithScore(key string, start, end int) (map[string]int, error)
		ZRangeByScore(key string, min, max interface{}, limitArgs ...int) ([]string, error)
		ZRevRange(key string, start, end int) ([]string, error)
		ZRevRangeByScore(key string, min, max interface{}, limitArgs ...int) ([]string, error)
		ZRem(key string, members ...interface{}) error
		ZRemRangeByScore(key string, min, max interface{}) error
		ZRemRangeByRank(key string, start, end int) error
		ZCard(key string) (int, error)
		ZScore(key string, member interface{}) (int, error)
		ZRank(key string, member interface{}) (int, error)
		ZRevRank(key string, member interface{}) (int, error)
		ZCount(key string, min, max interface{}) (int, error)

		Pipeline(commands []map[string][]interface{}, watchKeys ...interface{}) (interface{}, error)

		SelectDb(index int) error
		BgSave() error
		FlushDb(index int) error
		FlushAll() error
	}
)
