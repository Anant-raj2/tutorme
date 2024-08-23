package auth

import (
	"context"
	"net/http"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/web/templa/auth"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	*db.Queries
}

func New(db *db.Queries) *Handler {
	return &Handler{
		db,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	return nil
}

func (h *Handler) RenderRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	component := auth.Register()
	component.Render(context.Background(), w)
}


// Configurable constants
const (
	TokenExpiration    = 15 * time.Minute
	RefreshExpiration  = 7 * 24 * time.Hour
	MaxLoginAttempts   = 5
	LockoutDuration    = 15 * time.Minute
	PasswordMinLength  = 12
	PasswordMaxLength  = 64
	ArgonTime          = 3
	ArgonMemory        = 64 * 1024
	ArgonThreads       = 4
	ArgonKeyLength     = 32
	BcryptCost         = 12
	MFATokenLength     = 6
	MFATokenExpiration = 5 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account is locked")
	ErrPasswordPolicy     = errors.New("password does not meet policy requirements")
	ErrMFARequired        = errors.New("multi-factor authentication required")
	ErrInvalidMFAToken    = errors.New("invalid MFA token")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
)

type User struct {
	ID             uuid.UUID
	Username       string
	Email          string
	HashedPassword string
	Salt           []byte
	Roles          []string
	MFAEnabled     bool
	MFASecret      string
	LastLogin      time.Time
	LoginAttempts  int
	LockoutUntil   time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AuthService struct {
	DB    *sql.DB
	Redis *redis.Client
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(db *sql.DB, redis *redis.Client) *AuthService {
	return &AuthService{
		DB:    db,
		Redis: redis,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	if !isValidPassword(password) {
		return ErrPasswordPolicy
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	hashedPassword, err := hashPassword(password, salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (id, username, email, hashed_password, salt, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
	`
	_, err = s.DB.ExecContext(ctx, query, uuid.New(), username, email, hashedPassword, salt, time.Now())
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errors.New("username or email already exists")
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	user, err := s.getUserByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if time.Now().Before(user.LockoutUntil) {
		return nil, ErrAccountLocked
	}

	if !verifyPassword(password, user.HashedPassword, user.Salt) {
		err = s.incrementLoginAttempts(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to increment login attempts: %w", err)
		}
		return nil, ErrInvalidCredentials
	}

	if user.MFAEnabled {
		mfaToken, err := s.generateMFAToken(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate MFA token: %w", err)
		}
		return nil, &MFARequiredError{Token: mfaToken}
	}

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) VerifyMFA(ctx context.Context, userID uuid.UUID, mfaToken string) (*TokenPair, error) {
	storedToken, err := s.Redis.Get(ctx, fmt.Sprintf("mfa:%s", userID)).Result()
	if err != nil {
		return nil, ErrInvalidMFAToken
	}

	if storedToken != mfaToken {
		return nil, ErrInvalidMFAToken
	}

	user, err := s.getUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("refresh_secret"), nil // Replace with your actual secret
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.getUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) getUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username, email, hashed_password, salt, mfa_enabled, mfa_secret, last_login, login_attempts, lockout_until, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	var user User
	err := s.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Salt, &user.MFAEnabled, &user.MFASecret,
		&user.LastLogin, &user.LoginAttempts, &user.LockoutUntil, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	rolesQuery := `SELECT role FROM user_roles WHERE user_id = $1`
	rows, err := s.DB.QueryContext(ctx, rolesQuery, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}

func (s *AuthService) getUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, username, email, hashed_password, salt, mfa_enabled, mfa_secret, last_login, login_attempts, lockout_until, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user User
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Salt, &user.MFAEnabled, &user.MFASecret,
		&user.LastLogin, &user.LoginAttempts, &user.LockoutUntil, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	rolesQuery := `SELECT role FROM user_roles WHERE user_id = $1`
	rows, err := s.DB.QueryContext(ctx, rolesQuery, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}

func (s *AuthService) incrementLoginAttempts(ctx context.Context, user *User) error {
	user.LoginAttempts++
	if user.LoginAttempts >= MaxLoginAttempts {
		user.LockoutUntil = time.Now().Add(LockoutDuration)
	}

	query := `
		UPDATE users
		SET login_attempts = $1, lockout_until = $2
		WHERE id = $3
	`
	_, err := s.DB.ExecContext(ctx, query, user.LoginAttempts, user.LockoutUntil, user.ID)
	return err
}

func (s *AuthService) generateMFAToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := make([]byte, MFATokenLength)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	mfaToken := base64.StdEncoding.EncodeToString(token)
	err := s.Redis.Set(ctx, fmt.Sprintf("mfa:%s", userID), mfaToken, MFATokenExpiration).Err()
	if err != nil {
		return "", err
	}

	return mfaToken, nil
}

func (s *AuthService) generateTokenPair(ctx context.Context, user *User) (*TokenPair, error) {
	accessToken, err := generateJWT(user, TokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := generateJWT(user, RefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateJWT(user *User, expiration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key")) // Replace with your actual secret key
}

func hashPassword(password string, salt []byte) (string, error) {
	hash := argon2.IDKey([]byte(password), salt, ArgonTime, ArgonMemory, ArgonThreads, ArgonKeyLength)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, ArgonMemory, ArgonTime, ArgonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash)), nil
}

