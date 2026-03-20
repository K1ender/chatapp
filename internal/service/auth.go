package service

import (
	"chatapp/internal/model"
	"chatapp/internal/repository"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/curve25519"
)

var (
	ErrMagicLinkNotFound = fmt.Errorf("magic link not found")
	ErrMagicLinkExpired  = fmt.Errorf("magic link expired")
	ErrMagicLinkUsed     = fmt.Errorf("magic link already used")
)

type Auth interface {
	// FIXME: Perhaps this is not an AuthService task, but an EmailService task?

	// SendMagicLink sends a magic link to the user.
	SendMagicLink(ctx context.Context, email, magicLink string) error

	// VerifyMagicLink verifies a magic link.
	VerifyMagicLink(ctx context.Context, token string) (uuid.UUID, error)
}

type AuthService struct {
	EmailService  EmailService
	UserRepo      repository.User
	MagicLinkRepo repository.MagicLink
	Salt          string
}

func NewAuthService(emailService EmailService, magicLinkRepo repository.MagicLink, userRepo repository.User, salt string) Auth {
	return &AuthService{
		UserRepo:      userRepo,
		EmailService:  emailService,
		MagicLinkRepo: magicLinkRepo,
		Salt:          salt,
	}
}

// SendMagicLink implements [Auth].
func (a *AuthService) SendMagicLink(ctx context.Context, email string, magicLink string) error {
	magicLinkToken := uuid.New().String()

	_, err := a.UserRepo.FindUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("send magic link: %w", err)
	}

	if errors.Is(err, repository.ErrUserNotFound) {
		// signing key
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("send magic link: %w", err)
		}

		// encryption key
		var encPriv [32]byte
		rand.Read(encPriv[:])
		encPub, err := curve25519.X25519(encPriv[:], curve25519.Basepoint)
		if err != nil {
			return fmt.Errorf("send magic link: %w", err)
		}

		privateBundle := append(priv, encPriv[:]...)
		key := deriveKey([]byte(a.Salt), magicLinkToken)
		encryptedPrivateBundle, err := encrypt(key, privateBundle)

		a.UserRepo.CreateUser(ctx, model.User{
			Email:               email,
			SigningPublicKey:    pub,
			EncryptionPublicKey: encPub,
			EncryptedPrivateKey: encryptedPrivateBundle,
		})
	}

	err = a.EmailService.SendMagicLink(ctx, email, magicLink)
	if err != nil {
		return fmt.Errorf("send magic link: %w", err)
	}

	return nil
}

// VerifyMagicLink implements [Auth].
func (a *AuthService) VerifyMagicLink(ctx context.Context, token string) (uuid.UUID, error) {
	magicLink, err := a.MagicLinkRepo.FindMagicLinkByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrMagicLinkNotFound) {
			return uuid.Nil, ErrMagicLinkNotFound
		}

		return uuid.Nil, fmt.Errorf("verify magic link: %w", err)
	}

	if magicLink.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, ErrMagicLinkExpired
	}

	err = a.MagicLinkRepo.UseMagicLink(ctx, magicLink.ID)
	if err != nil {
		if errors.Is(err, repository.ErrMagicLinkAlreadyUsed) {
			return uuid.Nil, ErrMagicLinkUsed
		}

		return uuid.Nil, fmt.Errorf("verify magic link: %w", err)
	}

	return magicLink.UserID, nil
}

func deriveKey(salt []byte, password string) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

func encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt (aes.NewCipher): %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encrypt (cipher.NewGCM): %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("encrypt (rand.Read): %w", err)
	}

	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
