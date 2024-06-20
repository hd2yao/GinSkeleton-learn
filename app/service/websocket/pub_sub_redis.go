package websocket

import "goskeleton/app/utils/redis_factory"

func Pub(msg string) {
    RedisCli := redis_factory.GetOneRedisClient()
    RedisCli.Client.Do("Publish", "tv_flush", msg)
    RedisCli.ReleaseOneRedisClient()
}
