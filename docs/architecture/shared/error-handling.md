# 共通エラーハンドリング

## エラーモデル統一仕様

### 標準エラーレスポンス
```typescript
interface ErrorResponse {
  status: number           // HTTPステータスコード
  code: string            // 内部エラーコード
  message: string         // エラーメッセージ（日本語）
  details?: any          // 詳細情報（オプショナル）
  timestamp: string       // エラー発生時刻（ISO 8601）
  request_id: string      // リクエストID（トレーシング用）
}
```

### エラーコード体系
```
{SERVICE}_{CATEGORY}_{TYPE}

例：
- USER_AUTH_INVALID_CREDENTIALS
- STORE_PRODUCT_NOT_FOUND  
- MEDIA_UPLOAD_FILE_TOO_LARGE
- SYSTEM_DATABASE_CONNECTION_FAILED
```

## HTTP エラーステータス

### クライアントエラー (4xx)
```typescript
enum ClientErrorStatus {
  BAD_REQUEST = 400,           // リクエスト形式エラー
  UNAUTHORIZED = 401,          // 認証エラー
  FORBIDDEN = 403,             // 認可エラー
  NOT_FOUND = 404,             // リソース未発見
  METHOD_NOT_ALLOWED = 405,    // HTTPメソッドエラー
  CONFLICT = 409,              // リソース競合
  PRECONDITION_FAILED = 412,   // ビジネスルール違反
  UNPROCESSABLE_ENTITY = 422,  // バリデーションエラー
  TOO_MANY_REQUESTS = 429      // レート制限
}
```

### サーバーエラー (5xx)
```typescript
enum ServerErrorStatus {
  INTERNAL_SERVER_ERROR = 500, // 予期しないサーバーエラー
  NOT_IMPLEMENTED = 501,       // 未実装機能
  BAD_GATEWAY = 502,           // 外部サービスエラー
  SERVICE_UNAVAILABLE = 503,   // サービス一時停止
  GATEWAY_TIMEOUT = 504        // 外部サービスタイムアウト
}
```

## ドメイン別エラー定義

### 認証・認可エラー
```go
// codes/auth_errors.go
const (
    AuthInvalidCredentials    = "AUTH_INVALID_CREDENTIALS"
    AuthTokenExpired         = "AUTH_TOKEN_EXPIRED"
    AuthTokenInvalid         = "AUTH_TOKEN_INVALID"
    AuthInsufficientPerms    = "AUTH_INSUFFICIENT_PERMISSIONS"
    AuthAccountLocked        = "AUTH_ACCOUNT_LOCKED"
    AuthPasswordTooWeak      = "AUTH_PASSWORD_TOO_WEAK"
    AuthEmailNotVerified     = "AUTH_EMAIL_NOT_VERIFIED"
    AuthMFARequired          = "AUTH_MFA_REQUIRED"
)

var AuthErrorMessages = map[string]string{
    AuthInvalidCredentials: "認証情報が正しくありません",
    AuthTokenExpired:      "トークンの有効期限が切れています",
    AuthTokenInvalid:      "無効なトークンです",
    AuthInsufficientPerms: "この操作を実行する権限がありません",
    AuthAccountLocked:     "アカウントがロックされています",
    AuthPasswordTooWeak:   "パスワードが複雑性要件を満たしていません",
    AuthEmailNotVerified:  "メールアドレスが認証されていません",
    AuthMFARequired:       "多要素認証が必要です",
}
```

### ビジネスロジックエラー
```go
// codes/business_errors.go  
const (
    // 商品関連
    ProductNotFound          = "PRODUCT_NOT_FOUND"
    ProductOutOfStock       = "PRODUCT_OUT_OF_STOCK" 
    ProductInactive         = "PRODUCT_INACTIVE"
    ProductCategoryInvalid  = "PRODUCT_CATEGORY_INVALID"
    
    // 注文関連
    OrderNotFound           = "ORDER_NOT_FOUND"
    OrderAlreadyCancelled   = "ORDER_ALREADY_CANCELLED"
    OrderCannotCancel       = "ORDER_CANNOT_CANCEL"
    OrderPaymentFailed      = "ORDER_PAYMENT_FAILED"
    
    // カート関連
    CartItemNotFound        = "CART_ITEM_NOT_FOUND"
    CartQuantityExceeded    = "CART_QUANTITY_EXCEEDED"
    CartEmpty              = "CART_EMPTY"
    
    // 在庫関連
    StockInsufficient      = "STOCK_INSUFFICIENT"
    StockReserved          = "STOCK_RESERVED"
)
```

### システムエラー
```go
// codes/system_errors.go
const (
    SystemDatabaseError    = "SYSTEM_DATABASE_ERROR"
    SystemCacheError      = "SYSTEM_CACHE_ERROR"
    SystemExternalAPI     = "SYSTEM_EXTERNAL_API_ERROR"
    SystemFileUpload      = "SYSTEM_FILE_UPLOAD_ERROR"
    SystemQueueError      = "SYSTEM_QUEUE_ERROR"
    SystemRateLimit       = "SYSTEM_RATE_LIMIT_EXCEEDED"
)
```

## エラーハンドリング実装

### Goバックエンド実装
```go
// exception/error_handler.go
type AppError struct {
    Code       string                 `json:"code"`
    Message    string                 `json:"message"`
    Details    map[string]interface{} `json:"details,omitempty"`
    StatusCode int                    `json:"-"`
    Cause      error                  `json:"-"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAppError(code string, statusCode int, details map[string]interface{}) *AppError {
    message := getErrorMessage(code)
    return &AppError{
        Code:       code,
        Message:    message,
        StatusCode: statusCode,
        Details:    details,
    }
}

// エラーハンドリングミドルウェア
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        var err error
        
        switch e := recovered.(type) {
        case *AppError:
            handleAppError(c, e)
            return
        case error:
            err = e
        default:
            err = fmt.Errorf("%v", recovered)
        }
        
        // 予期しないエラー
        logger.Error("unexpected error", zap.Error(err))
        
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:     http.StatusInternalServerError,
            Code:       SystemInternalError,
            Message:    "内部エラーが発生しました",
            Timestamp:  time.Now().Format(time.RFC3339),
            RequestID:  getRequestID(c),
        })
    })
}

func handleAppError(c *gin.Context, err *AppError) {
    response := ErrorResponse{
        Status:     err.StatusCode,
        Code:       err.Code,
        Message:    err.Message,
        Details:    err.Details,
        Timestamp:  time.Now().Format(time.RFC3339),
        RequestID:  getRequestID(c),
    }
    
    // エラーレベルに応じてログ出力
    if err.StatusCode >= 500 {
        logger.Error("server error", 
            zap.String("code", err.Code),
            zap.Error(err.Cause))
    } else {
        logger.Warn("client error",
            zap.String("code", err.Code),
            zap.String("message", err.Message))
    }
    
    c.JSON(err.StatusCode, response)
}
```

### バリデーションエラー処理
```go
// validation/errors.go
type ValidationError struct {
    Field   string `json:"field"`
    Value   string `json:"value"`
    Message string `json:"message"`
}

func HandleValidationErrors(c *gin.Context, err error) {
    var validationErrors []ValidationError
    
    if ve, ok := err.(validator.ValidationErrors); ok {
        for _, e := range ve {
            validationErrors = append(validationErrors, ValidationError{
                Field:   getFieldName(e),
                Value:   fmt.Sprintf("%v", e.Value()),
                Message: getValidationMessage(e),
            })
        }
    }
    
    c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
        Status:    http.StatusUnprocessableEntity,
        Code:      "VALIDATION_FAILED",
        Message:   "入力データが正しくありません",
        Details:   map[string]interface{}{"errors": validationErrors},
        Timestamp: time.Now().Format(time.RFC3339),
        RequestID: getRequestID(c),
    })
}

func getValidationMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "この項目は必須です"
    case "email":
        return "有効なメールアドレスを入力してください"
    case "min":
        return fmt.Sprintf("最小値は%sです", fe.Param())
    case "max":
        return fmt.Sprintf("最大値は%sです", fe.Param())
    default:
        return "入力値が正しくありません"
    }
}
```

## フロントエンドエラーハンドリング

### Vue.js実装
```typescript
// composables/useErrorHandler.ts
interface ErrorNotification {
  type: 'error' | 'warning' | 'info'
  message: string
  duration?: number
  action?: {
    label: string
    handler: () => void
  }
}

export const useErrorHandler = () => {
  const { addNotification } = useUIStore()
  
  const handleError = (error: any) => {
    const appError = normalizeError(error)
    
    // エラーレベルに応じた通知
    const notification: ErrorNotification = {
      type: getNotificationType(appError.status),
      message: appError.message,
      duration: getNotificationDuration(appError.status)
    }
    
    // 特定エラーへの特別対応
    if (appError.code === 'AUTH_TOKEN_EXPIRED') {
      notification.action = {
        label: '再ログイン',
        handler: () => navigateTo('/login')
      }
    }
    
    addNotification(notification)
    
    // Sentry等への報告
    if (appError.status >= 500) {
      reportError(appError)
    }
  }
  
  const normalizeError = (error: any): AppError => {
    if (error.response?.data) {
      return error.response.data
    }
    
    if (error.code === 'NETWORK_ERROR') {
      return {
        status: 0,
        code: 'NETWORK_ERROR',
        message: 'ネットワークエラーが発生しました',
        timestamp: new Date().toISOString(),
        request_id: generateRequestId()
      }
    }
    
    return {
      status: 500,
      code: 'UNKNOWN_ERROR',
      message: 'エラーが発生しました',
      timestamp: new Date().toISOString(),
      request_id: generateRequestId()
    }
  }
  
  return { handleError }
}
```

### Axios インターセプター
```typescript
// lib/api/interceptors.ts
export const setupErrorInterceptor = (axios: AxiosInstance) => {
  axios.interceptors.response.use(
    (response) => response,
    (error) => {
      const { handleError } = useErrorHandler()
      
      // 自動リトライ対象のエラー
      if (shouldRetry(error)) {
        return retryRequest(error.config)
      }
      
      // 認証エラーの自動処理
      if (error.response?.status === 401) {
        handleAuthError(error)
        return Promise.reject(error)
      }
      
      // その他のエラー処理
      handleError(error)
      return Promise.reject(error)
    }
  )
}

const shouldRetry = (error: any): boolean => {
  const retryableCodes = [408, 429, 502, 503, 504]
  return retryableCodes.includes(error.response?.status)
}

const retryRequest = async (config: any): Promise<any> => {
  const maxRetries = 3
  const delay = 1000
  
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      await new Promise(resolve => setTimeout(resolve, delay * attempt))
      return await axios.request(config)
    } catch (error) {
      if (attempt === maxRetries) {
        throw error
      }
    }
  }
}
```

## ログ・監視

### 構造化ログ出力
```go
// logging/error_logger.go
func LogError(ctx context.Context, err *AppError) {
    fields := []zap.Field{
        zap.String("error_code", err.Code),
        zap.String("error_message", err.Message),
        zap.Int("status_code", err.StatusCode),
        zap.String("request_id", getRequestID(ctx)),
        zap.String("user_id", getUserID(ctx)),
        zap.String("service", "gateway"),
        zap.Time("timestamp", time.Now()),
    }
    
    if err.Details != nil {
        fields = append(fields, zap.Any("details", err.Details))
    }
    
    if err.Cause != nil {
        fields = append(fields, zap.Error(err.Cause))
    }
    
    if err.StatusCode >= 500 {
        logger.Error("server error", fields...)
    } else {
        logger.Warn("client error", fields...)
    }
}
```

### メトリクス収集
```go
// metrics/error_metrics.go
var (
    ErrorCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "furumaru_errors_total",
            Help: "Total number of errors by code and status",
        },
        []string{"error_code", "status_code", "service"},
    )
    
    ErrorDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "furumaru_error_duration_seconds",
            Help: "Duration from request start to error",
        },
        []string{"error_code", "service"},
    )
)

func RecordError(code string, statusCode int, service string, duration time.Duration) {
    ErrorCounter.WithLabelValues(
        code,
        strconv.Itoa(statusCode),
        service,
    ).Inc()
    
    ErrorDuration.WithLabelValues(code, service).Observe(duration.Seconds())
}
```

## 外部サービス連携エラー

### 外部API呼び出しエラー
```go
// integration/external_error.go
type ExternalServiceError struct {
    Service    string `json:"service"`
    Endpoint   string `json:"endpoint"`
    StatusCode int    `json:"status_code"`
    Message    string `json:"message"`
    RetryAfter *int   `json:"retry_after,omitempty"`
}

func HandleExternalError(service, endpoint string, err error) error {
    switch e := err.(type) {
    case *http.ClientError:
        if e.StatusCode == 429 {
            // レート制限エラー
            return NewAppError(
                SystemRateLimit,
                http.StatusServiceUnavailable,
                map[string]interface{}{
                    "service": service,
                    "retry_after": e.Headers.Get("Retry-After"),
                },
            )
        }
        
        return NewAppError(
            SystemExternalAPI,
            http.StatusBadGateway,
            map[string]interface{}{
                "service": service,
                "status_code": e.StatusCode,
                "message": e.Message,
            },
        )
        
    case *http.TimeoutError:
        return NewAppError(
            "EXTERNAL_TIMEOUT",
            http.StatusGatewayTimeout,
            map[string]interface{}{
                "service": service,
                "endpoint": endpoint,
            },
        )
        
    default:
        return NewAppError(
            SystemExternalAPI,
            http.StatusBadGateway,
            map[string]interface{}{
                "service": service,
                "error": err.Error(),
            },
        )
    }
}
```

## テスト戦略

### エラーハンドリングテスト
```go
func TestErrorHandler(t *testing.T) {
    tests := []struct {
        name           string
        error          *AppError
        expectedStatus int
        expectedCode   string
    }{
        {
            name: "validation error",
            error: NewAppError(
                "VALIDATION_FAILED",
                http.StatusUnprocessableEntity,
                map[string]interface{}{
                    "field": "email",
                    "message": "invalid format",
                },
            ),
            expectedStatus: http.StatusUnprocessableEntity,
            expectedCode:   "VALIDATION_FAILED",
        },
        {
            name: "authentication error",
            error: NewAppError(
                AuthInvalidCredentials,
                http.StatusUnauthorized,
                nil,
            ),
            expectedStatus: http.StatusUnauthorized,
            expectedCode:   AuthInvalidCredentials,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            
            handleAppError(c, tt.error)
            
            assert.Equal(t, tt.expectedStatus, w.Code)
            
            var response ErrorResponse
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedCode, response.Code)
        })
    }
}
```