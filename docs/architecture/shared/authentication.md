# 共通認証・認可アーキテクチャ

## 認証方式概要

### サポート認証方式
```
┌─── AWS Cognito (JWT) ────┐    ┌─── OAuth Providers ───┐
│  - Email/Password        │    │  - Google OAuth       │
│  - MFA対応               │    │  - LINE OAuth         │
│  - ユーザープール管理     │    │  - Facebook OAuth     │
└──────────────────────────┘    └───────────────────────┘
            │                              │
            └──────────┬───────────────────┘
                      │
            ┌─────── Gateway ──────┐
            │  - Token Validation   │
            │  - Permission Check   │
            │  - Session Mgmt      │
            └──────────────────────┘
```

## 認証フロー

### JWT Bearer認証
```typescript
interface JWTPayload {
  sub: string          // ユーザーID
  aud: string          // オーディエンス
  iss: string          // 発行者
  exp: number          // 有効期限
  iat: number          // 発行時刻
  token_use: 'access'  // トークン用途
  scope: string        // スコープ
  username: string     // ユーザー名
  
  // カスタムクレーム
  'custom:role': string
  'custom:permissions': string[]
  'custom:tenant_id'?: string
}
```

### Cookie Session認証
```typescript
interface SessionData {
  sessionId: string
  userId?: string      // 認証済みの場合
  guestId?: string     // ゲストセッション
  cartId?: string      // カートセッション
  preferences: UserPreferences
  createdAt: Date
  lastAccessedAt: Date
  expiresAt: Date
}
```

## 認可モデル

### ロールベース認可 (RBAC)
```typescript
enum UserRole {
  ADMIN = 'admin',
  COORDINATOR = 'coordinator', 
  PRODUCER = 'producer',
  MEMBER = 'member',
  GUEST = 'guest'
}

enum Permission {
  // 商品管理
  PRODUCT_READ = 'product:read',
  PRODUCT_WRITE = 'product:write',
  PRODUCT_DELETE = 'product:delete',
  
  // 注文管理
  ORDER_READ = 'order:read',
  ORDER_WRITE = 'order:write',
  ORDER_CANCEL = 'order:cancel',
  
  // ユーザー管理
  USER_READ = 'user:read',
  USER_WRITE = 'user:write',
  USER_DELETE = 'user:delete',
  
  // システム管理
  SYSTEM_CONFIG = 'system:config',
  SYSTEM_MONITOR = 'system:monitor'
}

const RolePermissions: Record<UserRole, Permission[]> = {
  [UserRole.ADMIN]: [
    Permission.PRODUCT_READ,
    Permission.PRODUCT_WRITE,
    Permission.PRODUCT_DELETE,
    Permission.ORDER_READ,
    Permission.ORDER_WRITE,
    Permission.ORDER_CANCEL,
    Permission.USER_READ,
    Permission.USER_WRITE,
    Permission.USER_DELETE,
    Permission.SYSTEM_CONFIG,
    Permission.SYSTEM_MONITOR
  ],
  [UserRole.COORDINATOR]: [
    Permission.PRODUCT_READ,
    Permission.PRODUCT_WRITE,
    Permission.ORDER_READ,
    Permission.ORDER_WRITE,
    Permission.USER_READ
  ],
  [UserRole.PRODUCER]: [
    Permission.PRODUCT_READ,
    Permission.PRODUCT_WRITE,
    Permission.ORDER_READ
  ],
  [UserRole.MEMBER]: [
    Permission.PRODUCT_READ,
    Permission.ORDER_READ,
    Permission.ORDER_WRITE
  ],
  [UserRole.GUEST]: [
    Permission.PRODUCT_READ
  ]
}
```

## ゲートウェイ認証処理

### 認証ミドルウェア
```go
// middleware/auth.go
func AuthMiddleware(schemes ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user *entity.User
        var session *entity.Session
        var err error
        
        for _, scheme := range schemes {
            switch scheme {
            case "bearerauth":
                user, err = validateJWTToken(c)
            case "cookieauth":
                session, err = validateSession(c)
            }
            
            if err == nil {
                break
            }
        }
        
        if user == nil && session == nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authentication required"
            })
            c.Abort()
            return
        }
        
        // コンテキストに設定
        if user != nil {
            c.Set("user", user)
        }
        if session != nil {
            c.Set("session", session)
        }
        
        c.Next()
    }
}
```

### 認可チェック
```go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Permission denied"
            })
            c.Abort()
            return
        }
        
        if !hasPermission(user.(*entity.User), permission) {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Insufficient permissions"
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## セッション管理

### セッション設計
```go
type SessionStore interface {
    Create(userID string, ttl time.Duration) (*Session, error)
    Get(sessionID string) (*Session, error)
    Update(sessionID string, data map[string]interface{}) error
    Delete(sessionID string) error
    Cleanup() error
}

type RedisSessionStore struct {
    client redis.Client
    prefix string
}

func (r *RedisSessionStore) Create(userID string, ttl time.Duration) (*Session, error) {
    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(ttl),
    }
    
    key := r.prefix + sessionID
    data, _ := json.Marshal(session)
    
    return session, r.client.Set(key, data, ttl).Err()
}
```

## トークン管理

### JWT検証
```go
func validateJWTToken(c *gin.Context) (*entity.User, error) {
    authHeader := c.GetHeader("Authorization")
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return nil, errors.New("invalid authorization header")
    }
    
    tokenString := authHeader[7:]
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return getJWTSecret(), nil
    })
    
    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    claims := token.Claims.(*CustomClaims)
    
    // トークンの有効性チェック
    if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
        return nil, errors.New("token expired")
    }
    
    // ユーザー情報の取得
    user, err := userService.GetByID(claims.Subject)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### リフレッシュトークン
```go
func RefreshToken(refreshToken string) (*TokenResponse, error) {
    // リフレッシュトークンの検証
    claims, err := validateRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }
    
    // 新しいアクセストークンを生成
    accessToken, err := generateAccessToken(claims.Subject)
    if err != nil {
        return nil, err
    }
    
    // 新しいリフレッシュトークンを生成
    newRefreshToken, err := generateRefreshToken(claims.Subject)
    if err != nil {
        return nil, err
    }
    
    return &TokenResponse{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
        ExpiresIn:    3600, // 1時間
    }, nil
}
```

## OAuth連携

### プロバイダー設定
```go
type OAuthConfig struct {
    Google GoogleOAuthConfig `json:"google"`
    LINE   LINEOAuthConfig   `json:"line"`
}

type GoogleOAuthConfig struct {
    ClientID     string   `json:"client_id"`
    ClientSecret string   `json:"client_secret"`
    RedirectURL  string   `json:"redirect_url"`
    Scopes       []string `json:"scopes"`
}

func (g *GoogleOAuthConfig) GetAuthURL(state string) string {
    return fmt.Sprintf(
        "https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%s",
        g.ClientID,
        strings.Join(g.Scopes, " "),
        url.QueryEscape(g.RedirectURL),
        state,
    )
}
```

### OAuth コールバック処理
```go
func HandleOAuthCallback(provider string) gin.HandlerFunc {
    return func(c *gin.Context) {
        code := c.Query("code")
        state := c.Query("state")
        
        // CSRF攻撃防止のため state を検証
        if !validateState(state) {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid state parameter"
            })
            return
        }
        
        // 認可コードをアクセストークンと交換
        oauthToken, err := exchangeCodeForToken(provider, code)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to exchange code"
            })
            return
        }
        
        // プロバイダーからユーザー情報を取得
        userInfo, err := getUserInfoFromProvider(provider, oauthToken.AccessToken)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to get user info"
            })
            return
        }
        
        // ユーザーを作成または更新
        user, err := createOrUpdateOAuthUser(provider, userInfo)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to create user"
            })
            return
        }
        
        // JWTトークンを生成
        token, err := generateJWTToken(user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to generate token"
            })
            return
        }
        
        // フロントエンドにリダイレクト
        redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", 
            config.FrontendURL, token)
        c.Redirect(http.StatusFound, redirectURL)
    }
}
```

## セキュリティ対策

### CSRF対策
```go
func CSRFProtection() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
            c.Next()
            return
        }
        
        token := c.GetHeader("X-CSRF-Token")
        if token == "" {
            token = c.PostForm("_token")
        }
        
        if !validateCSRFToken(token, c) {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "CSRF token validation failed"
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### レート制限
```go
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded"
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## ログ・監査

### 認証ログ
```go
type AuthEvent struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Result    string    `json:"result"`
    IP        string    `json:"ip"`
    UserAgent string    `json:"user_agent"`
    Timestamp time.Time `json:"timestamp"`
}

func LogAuthEvent(c *gin.Context, action, result string, userID string) {
    event := AuthEvent{
        UserID:    userID,
        Action:    action,
        Result:    result,
        IP:        c.ClientIP(),
        UserAgent: c.GetHeader("User-Agent"),
        Timestamp: time.Now(),
    }
    
    logger.Info("auth_event", zap.Any("event", event))
}
```

## エラーハンドリング

### 認証エラー
```go
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrTokenExpired      = errors.New("token expired")
    ErrInsufficientPerms = errors.New("insufficient permissions")
    ErrSessionNotFound   = errors.New("session not found")
)

func handleAuthError(c *gin.Context, err error) {
    switch err {
    case ErrInvalidCredentials:
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Code:    "INVALID_CREDENTIALS",
            Message: "The provided credentials are invalid",
        })
    case ErrTokenExpired:
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Code:    "TOKEN_EXPIRED",
            Message: "The access token has expired",
        })
    case ErrInsufficientPerms:
        c.JSON(http.StatusForbidden, ErrorResponse{
            Code:    "INSUFFICIENT_PERMISSIONS",
            Message: "You don't have permission to access this resource",
        })
    default:
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "INTERNAL_ERROR",
            Message: "An internal error occurred",
        })
    }
}
```