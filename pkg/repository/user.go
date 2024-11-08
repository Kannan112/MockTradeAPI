package repository

import (
	"context"
	"time"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/domain"
	"github.com/kannan112/mock-trading-platform-api/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB: DB}
}

func (c *userDatabase) FindUserByUserID(ctx context.Context, userID uint) (user domain.User, err error) {

	query := `SELECT * FROM users WHERE id = $1`
	err = c.DB.Raw(query, userID).Scan(&user).Error

	return user, err
}

func (c *userDatabase) FindUserByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `select exists(SELECT * FROM users WHERE email = $1)`
	err := c.DB.Raw(query, email).Scan(&exists).Error
	return exists, err
}

func (c *userDatabase) SaveUser(ctx context.Context, user request.RegisterUserRequest) (userID uint, err error) {

	//save the user details
	query := `INSERT INTO users (username,email, password, created_at) 
	VALUES ($1, $2, $3, $4 ) RETURNING id`

	createdAt := time.Now()
	err = c.DB.Raw(query, user.Username, user.Email, user.Password, createdAt).Scan(&userID).Error

	return userID, err
}

func (c *userDatabase) ExtractPassword(ctx context.Context, email string) (string, error) {
	var hash string
	query := `SELECT password from users where email=$1`
	err := c.DB.Raw(query, email).Scan(&hash).Error
	return hash, err

}

func (c *userDatabase) GetUserId(ctx context.Context, email string) (int, error) {
	var userId int
	query := `SELECT id from users where email=$1`
	err := c.DB.Raw(query, email).Scan(&userId).Error
	return userId, err
}
