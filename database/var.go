package database

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	MainDB      *gorm.DB
	RedisClient *redis.Client
)
