package use_case_test

import (
	"encoding/json"
	"github.com/wisle25/be-template/applications/file_handling"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/entity"
	"io"
	"mime/multipart"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wisle25/be-template/commons"
)

func TestUserUseCase(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPasswordHash := new(MockPasswordHash)
	mockValidator := new(MockValidateUser)
	mockConfig := &commons.Config{
		AccessTokenExpiresIn:  time.Hour,
		RefreshTokenExpiresIn: time.Hour * 24,
		AccessTokenPrivateKey: "any",
		RefreshTokenPublicKey: "any",
	}
	mockToken := new(MockToken)
	mockCache := new(MockCache)
	mockFileUpload := new(MockFileUpload)
	mockFileProcessing := new(MockFileProcessing)

	userUseCase := use_case.NewUserUseCase(
		mockUserRepo,
		mockFileProcessing,
		mockFileUpload,
		mockPasswordHash,
		mockValidator,
		mockConfig,
		mockToken,
		mockCache,
	)

	t.Run("Execute Register", func(t *testing.T) {
		// Arrange
		payload := &entity.RegisterUserPayload{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		}

		mockValidator.On("ValidateRegisterPayload", payload).Return(nil)
		mockPasswordHash.On("Hash", payload.Password).Return("hashedpassword")
		mockUserRepo.On("RegisterUser", payload).Return()

		// Action
		userUseCase.ExecuteRegister(payload)

		// Assert
		mockValidator.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockPasswordHash.AssertExpectations(t)
	})

	t.Run("Execute Login", func(t *testing.T) {
		// Arrange
		payload := &entity.LoginUserPayload{
			Identity: "testuser",
			Password: "password123",
		}

		user := &entity.User{
			Id:       "userid123",
			Username: "testuser",
			Email:    "test@example.com",
		}

		accessTokenDetail := &entity.TokenDetail{
			TokenId:    "access_token_id",
			ExpiresIn:  time.Now().Add(time.Hour).Unix(),
			UserDetail: user,
			Token:      "access_token",
		}

		refreshTokenDetail := &entity.TokenDetail{
			TokenId:    "refresh_token_id",
			ExpiresIn:  time.Now().Add(time.Hour * 24).Unix(),
			UserDetail: user,
			Token:      "refresh_token",
		}

		mockValidator.On("ValidateLoginPayload", payload).Return(nil)
		mockUserRepo.On("GetUserForLogin", payload.Identity).Return(user, "hashedpassword")
		mockPasswordHash.On("Compare", payload.Password, "hashedpassword").Return(nil)
		mockToken.On("CreateToken", user, mockConfig.AccessTokenExpiresIn, mockConfig.AccessTokenPrivateKey).Return(accessTokenDetail)
		mockToken.On("CreateToken", user, mockConfig.RefreshTokenExpiresIn, mockConfig.RefreshTokenPrivateKey).Return(refreshTokenDetail)

		userInfoJSON, _ := json.Marshal(user)
		mockCache.On("SetCache", accessTokenDetail.TokenId, userInfoJSON, mock.Anything).Return(nil)
		mockCache.On("SetCache", refreshTokenDetail.TokenId, userInfoJSON, mock.Anything).Return(nil)

		// Action
		accessToken, refreshToken := userUseCase.ExecuteLogin(payload)

		// Assert
		assert.Equal(t, accessTokenDetail, accessToken)
		assert.Equal(t, refreshTokenDetail, refreshToken)

		mockValidator.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockPasswordHash.AssertExpectations(t)
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute Refresh Token", func(t *testing.T) {
		// Arrange
		user := entity.User{
			Id:       "userid123",
			Username: "testuser",
			Email:    "test@example.com",
		}
		refreshTokenCookie := "refresh_token123"

		accessTokenDetail := &entity.TokenDetail{
			TokenId:    "access_token_id",
			ExpiresIn:  time.Now().Add(time.Hour).Unix(),
			UserDetail: &user,
			Token:      "access_token",
		}
		refreshTokenDetail := &entity.TokenDetail{
			TokenId:    "refresh_token_id",
			ExpiresIn:  time.Now().Add(time.Hour * 24).Unix(),
			UserDetail: &user,
			Token:      "refresh_token",
		}

		userInfoJSON, _ := json.Marshal(user)

		mockToken.On("ValidateToken", refreshTokenCookie, mockConfig.RefreshTokenPublicKey).Return(refreshTokenDetail)
		mockCache.On("GetCache", refreshTokenDetail.TokenId).Return(string(userInfoJSON))
		mockToken.On("CreateToken", &user, mockConfig.AccessTokenExpiresIn, mockConfig.AccessTokenPrivateKey).Return(accessTokenDetail)
		mockCache.On("SetCache", accessTokenDetail.TokenId, string(userInfoJSON), mock.Anything).Return(nil)

		// Action
		accessTokenResponse := userUseCase.ExecuteRefreshToken(refreshTokenCookie)

		// Assert
		assert.Equal(t, accessTokenDetail, accessTokenResponse)
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute Logout", func(t *testing.T) {
		// Arrange
		refreshTokenCookie := "refresh_token123"
		accessTokenId := "access_token_id"

		refreshTokenDetail := &entity.TokenDetail{
			TokenId: "refresh_token_id",
		}

		mockToken.On("ValidateToken", refreshTokenCookie, mockConfig.RefreshTokenPublicKey).Return(refreshTokenDetail)
		mockCache.On("DeleteCache", refreshTokenDetail.TokenId).Return(nil).Once()
		mockCache.On("DeleteCache", accessTokenId).Return(nil).Once()

		// Action
		userUseCase.ExecuteLogout(refreshTokenCookie, accessTokenId)

		// Assert
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute Guard", func(t *testing.T) {
		// Arrange
		user := &entity.User{
			Id:       "userid123",
			Username: "testuser",
			Email:    "test@example.com",
		}

		accessToken := "access_token123"
		accessTokenDetail := &entity.TokenDetail{
			TokenId:    "access_token123",
			UserDetail: user,
		}

		tokenDetailJSON, _ := json.Marshal(accessTokenDetail)

		mockToken.On("ValidateToken", accessToken, mockConfig.AccessTokenPublicKey).Return(accessTokenDetail)
		mockCache.On("GetCache", accessTokenDetail.TokenId).Return(tokenDetailJSON)

		// Action
		userIdCache, tokenDetail := userUseCase.ExecuteGuard(accessToken)

		// Unmarshal the cached userId
		var cachedUserId string
		_ = json.Unmarshal(userIdCache.([]byte), &cachedUserId)

		// Assert
		assert.NotNil(t, cachedUserId)
		assert.Equal(t, accessTokenDetail, tokenDetail)
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute GetUserById", func(t *testing.T) {
		// Arrange
		userId := "userid123"
		expectedUser := &entity.User{
			Id:         userId,
			Username:   "testuser",
			Email:      "test@example.com",
			AvatarLink: "anything",
		}

		mockUserRepo.On("GetUserById", userId).Return(expectedUser)

		// Action
		user := userUseCase.ExecuteGetUserById(userId)

		// Assert
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Execute UpdateUserById", func(t *testing.T) {
		// Arrange
		userId := "userid123"
		payload := &entity.UpdateUserPayload{
			Username:        "username",
			Email:           "email",
			Password:        "any password",
			ConfirmPassword: "any password",
			Avatar: &multipart.FileHeader{
				Filename: "Any",
				Header:   nil,
				Size:     0,
			},
		}
		avatarLink := "avatar_link"
		oldAvatarLink := "old_avatar_link"
		file, _ := payload.Avatar.Open()
		avatarBuffer, _ := io.ReadAll(file)
		var compressBuffer []byte

		minioUrl := mockConfig.MinioUrl + mockConfig.MinioBucket + "/"

		mockValidator.On("ValidateUpdatePayload", payload).Return(nil)
		mockPasswordHash.On("Hash", payload.Password).Return("hashedPassword")
		mockFileProcessing.On("CompressImage", avatarBuffer, file_handling.WEBP).Return(compressBuffer, ".webp")
		mockFileUpload.On("UploadFile", compressBuffer, ".webp").Return(avatarLink)
		mockUserRepo.On("UpdateUserById", userId, payload, minioUrl+avatarLink).Return(oldAvatarLink)
		mockFileUpload.On("RemoveFile", path.Base(oldAvatarLink)).Return(nil)

		// Action
		userUseCase.ExecuteUpdateUserById(userId, payload)

		// Assert
		mockValidator.AssertExpectations(t)
		mockPasswordHash.AssertExpectations(t)
		mockFileProcessing.AssertExpectations(t)
		mockFileUpload.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
