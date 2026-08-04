package service

import(
	"weblog/repo"
	"weblog/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{
	userRepo *repo.UserRepo
}

func NewAuthService(u *repo.UserRepo) *AuthService{
	return &AuthService{userRepo:u}
}

func (a *AuthService) Signup(username, password string) (*model.User, error) {
	_, exists, err := a.userRepo.FindUser(username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("this username is already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := a.userRepo.AddUser(user); err != nil {
		return nil, err
	}

	createdUser, _, err := a.userRepo.FindUser(username)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (a *AuthService) Login(username, password string) (*model.User, error){
	user, status, err := a.userRepo.FindUser(username)

	if err != nil{
		return nil, err
	}

	if !status{
		return nil, errors.New("username or password is wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("username or password is wrong")
    }

    return user, nil
}