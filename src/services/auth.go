package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"vibely-backend/src/models"
)

type AuthService interface {
	GenerateAccessToken(user models.User) (string, error)
	GenerateRefreshToken(user models.User) (string, error)
	ExtractUserIDfromAccessToken(tokenString string) (uuid.UUID, error)
	ExtractUserIDfromRefreshToken(tokenString string) (uuid.UUID, error)
	SetStudentUserRole(user *models.User)
	SetTutorUserRole(user *models.User)
}

type authService struct {
	AccessJWTSecretKey  string
	RefreshJWTSecretKey string
}

func NewAuthService(accessJWTSecretKey, refreshJWTSecretKey string) AuthService {
	return &authService{
		AccessJWTSecretKey:  accessJWTSecretKey,
		RefreshJWTSecretKey: refreshJWTSecretKey,
	}
}

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserRole string    `json:"user_role"`
	jwt.RegisteredClaims
}

func (s *authService) GenerateAccessToken(user models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		UserRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(s.AccessJWTSecretKey) // Use secret key from the struct
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed signing token")
	}
	fmt.Println(tokenString)
	return tokenString, nil
}
func (s *authService) GenerateRefreshToken(user models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		UserRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshSecretKey := []byte(s.RefreshJWTSecretKey)
	refreshTokenString, err := token.SignedString(refreshSecretKey)
	if err != nil {
		return "", errors.New("failed signing refresh token")
	}
	return refreshTokenString, nil
}

func (s *authService) ExtractUserIDfromAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.AccessJWTSecretKey), nil // Use secret key from the struct
		})
	if err != nil {
		return uuid.UUID{}, errors.New("error while parsing token")
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.UUID{}, errors.New("invalid token claims")
}
func (s *authService) ExtractUserIDfromRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.RefreshJWTSecretKey), nil // Use refresh token secret
		})
	if err != nil {
		return uuid.UUID{}, errors.New("error while parsing refresh token")
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.UUID{}, errors.New("invalid refresh token claims")
}
func (s *authService) SetStudentUserRole(user *models.User) {
	user.Role = models.UserRoleStudent
}
func (s *authService) SetTutorUserRole(user *models.User) {
	user.Role = models.UserRoleTutor
}
