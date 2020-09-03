package model

import "net/http"

type RedisSettings struct {
	Enable             *bool
	Address            *string
	Password           *string
	Index              *int
	PoolSize           *int
	EnableRedisCluster *bool
}

func (s *RedisSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}
	if s.Address == nil {
		s.Address = NewString("localhost:6379")
	}
	if s.PoolSize == nil {
		s.PoolSize = NewInt(100)
	}
	if s.Index == nil {
		s.Index = NewInt(0)
	}
	if s.EnableRedisCluster == nil {
		s.EnableRedisCluster = NewBool(false)
	}
}
func (rs *RedisSettings) isValid() *AppError {
	if *rs.Enable {
		if *rs.Address == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.redis_empty_address.app_error", nil, "", http.StatusBadRequest)
		}
		if *rs.Index > 10 || *rs.Index < 0 {
			return NewAppError("Config.IsValid", "model.config.is_valid.redis_index.app_error", nil, "", http.StatusBadRequest)
		}
	}
	return nil
}

type SocketExporterSettings struct {
	Enable       *bool
	Region       *string
	StreamName   *string
	PartitionKey *string
}
