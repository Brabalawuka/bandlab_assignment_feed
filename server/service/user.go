package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/service/mocks"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

// UserService defines the interface for user operations
type UserService interface {
	GetUserById(id primitive.ObjectID) (*dao.User, error)
	GetAllUsers() []*dao.User
}

var (
	userOnce sync.Once
	userSrv  UserService
)

// InitUserService initializes the user service with mock data
func InitUserService() {
	userOnce.Do(func() {
		aliceId, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
		bobId, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
		charlieId, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")

		userSrv = &UserServiceImpl{
			users: []*dao.User{
				{Id: aliceId, Name: "Alice"},
				{Id: bobId, Name: "Bob"},
				{Id: charlieId, Name: "Charlie"},
			},
		}
	})
}

// GetUserService returns the initialized user service
//go:generate mockgen -destination=./mocks/mock_user_service.go -package=mocks -source=./user.go
func GetUserService() UserService {
	return userSrv
}

// SetMockPostService For unit testing purpose only
func SetMockUserService(ctrl *gomock.Controller) *mocks.MockUserService {
	mocks := mocks.NewMockUserService(ctrl)
	userSrv = mocks
	return mocks
}

// UserServiceImpl is the implementation of UserService
type UserServiceImpl struct {
	users []*dao.User
}

// GetUserById returns a user by their Id
func (s *UserServiceImpl) GetUserById(id primitive.ObjectID) (*dao.User, error) {
	for _, user := range s.users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, errs.ErrUserNotFound
}

// GetAllUsers returns all users
func (s *UserServiceImpl) GetAllUsers() []*dao.User {
	return s.users
}
