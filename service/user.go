package service

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mine/just-projecting/domain/kafka"
	"github.com/mine/just-projecting/domain/model"
	"github.com/mine/just-projecting/domain/redis"
	"github.com/mine/just-projecting/domain/repository"
)

type UserService interface {
	Find() ([]*model.User, error)
	FindOne() (*model.User, error)
	CountUsersUsecase(userID int64) (result string, err error)
	PublishMessageUser(ctx context.Context) error
}

type userService struct {
	// should be private
	UserRepository repository.UserRepository
	UserRedis      redis.UserRedisRepository
	UserKafka      kafka.UserKafkaRepository
	// cache, config, db transaction etc inject here
}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) SetUserRepository(userRepo repository.UserRepository) *userService {
	s.UserRepository = userRepo
	return s
}

func (s *userService) SetUserRedis(userRedis redis.UserRedisRepository) *userService {
	s.UserRedis = userRedis
	return s
}

func (s *userService) SetUserKafka(userKafka kafka.UserKafkaRepository) *userService {
	s.UserKafka = userKafka
	return s
}

func (s *userService) Validate() *userService {
	if s.UserRepository == nil {
		panic("handler need user repository")
	}
	return s
}

func (s *userService) Find() ([]*model.User, error) {
	/**
	 * business logic here
	 */
	return s.UserRepository.Find()
}

func (s *userService) FindOne() (*model.User, error) {
	return s.UserRepository.FindOne()
}

func (s *userService) CountUsersUsecase(userID int64) (result string, err error) {
	result, err = s.UserRepository.CountUserRepository(userID)
	if err != nil {
		panic("failed to call Repo")
	}

	return result, err
}

func (s *userService) PublishMessageUser(ctx context.Context) error {
	return s.UserKafka.Publish(ctx, model.User{
		ID:   1,
		Name: "kafka test",
		Age:  10,
	})
}
