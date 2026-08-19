package usecase

import (
	"errors"

	"github.com/prasetyoary48/wms/internal/domain"
	"github.com/prasetyoary48/wms/internal/repository"
	"github.com/prasetyoary48/wms/pkg/hash"
	"github.com/prasetyoary48/wms/pkg/jwtutil"
)

type AuthUsecase interface {
	Login(email, password string) (token string, user *domain.User, err error)
	Register(name, email, password string, roleID uint) (*domain.User, error)
}

type authUsecase struct {
	userRepo    repository.UserRepository
	jwtSecret   string
	jwtExpireHr int
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string, jwtExpireHr int) AuthUsecase {
	return &authUsecase{userRepo: userRepo, jwtSecret: jwtSecret, jwtExpireHr: jwtExpireHr}
}

func (u *authUsecase) Login(email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("email atau password salah")
	}
	if !user.IsActive {
		return "", nil, errors.New("akun tidak aktif, hubungi admin")
	}
	if !hash.CheckPassword(password, user.PasswordHash) {
		return "", nil, errors.New("email atau password salah")
	}

	token, err := jwtutil.GenerateToken(u.jwtSecret, user.ID, user.Role.Name, u.jwtExpireHr)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Register dipakai oleh Admin untuk membuat akun user baru (bukan self-signup publik).
func (u *authUsecase) Register(name, email, password string, roleID uint) (*domain.User, error) {
	existing, _ := u.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashed,
		RoleID:       roleID,
		IsActive:     true,
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}