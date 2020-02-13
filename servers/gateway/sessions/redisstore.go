package sessions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	store := &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
	return store
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	redisKey := sid.getRedisKey()
	state, err := json.Marshal(sessionState)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	err = rs.Client.Set(redisKey, state, cache.DefaultExpiration).Err()
	if err != nil {
		return err
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	prevSession, err := rs.Client.Get(sid.getRedisKey()).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrStateNotFound
		}
		return err
	}
	err = json.Unmarshal([]byte(prevSession), sessionState)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	rs.Client.Do("EXPIRE", sid.getRedisKey(), cache.DefaultExpiration)
	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	err := rs.Client.Del(sid.getRedisKey()).Err()
	if err != nil {
		return err
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	return "sid:" + sid.String()
}
