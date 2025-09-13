# Database Design

Furumaru プロジェクトのデータベース設計書です。モジュールごとに分離されたデータベース設計とエンティティ関係を定義します。

## Architecture

### Database Separation by Module

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    user_db      │    │   store_db      │    │   media_db      │    │ messenger_db    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ User Module     │    │ Store Module    │    │ Media Module    │    │ Messenger Mod   │
│ - Users         │    │ - Products      │    │ - Videos        │    │ - Notifications │
│ - Admins        │    │ - Orders        │    │ - Lives         │    │ - Messages      │
│ - Addresses     │    │ - Carts         │    │ - Comments      │    │ - Templates     │
│ - Auth          │    │ - Payments      │    │ - Streams       │    │ - Schedules     │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Technology Stack
- **Database**: TiDB (MySQL Compatible)
- **ORM**: GORM
- **Migration**: Custom tool (`api/hack/database-migrate-mysql`)
- **Connection Pool**: Per module configuration

## Module-wise Database Design

### 👤 user_db - User Management Database

#### Core Entities

**User (Base)**
```go
type User struct {
    ID        string    `gorm:"primaryKey"`
    Type      UserType  `gorm:"not null"`
    Status    UserStatus `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**User Types & Status**
```go
// UserType
UserTypeMember       = 1  // 会員
UserTypeGuest        = 2  // ゲスト
UserTypeFacilityUser = 3  // 施設利用者

// UserStatus  
UserStatusGuest       = 1  // 未登録
UserStatusProvisional = 2  // 仮登録
UserStatusVerified    = 3  // 認証済み
UserStatusDeactivated = 4  // 無効
```

#### Entity Relationships
```
User (1) ←→ (n) Address
User (1) ←→ (1) UserAuth  
User (1) ←→ (n) UserNotification
```

### 🛒 store_db - E-Commerce Database

#### Core Entities
- **Product**: 商品マスター
- **Order**: 注文
- **Cart**: カート
- **Payment**: 決済情報

#### Entity Relationships
```
Product (n) ←→ (n) Cart (through CartItem)
Cart (1) ←→ (1) Order
Order (1) ←→ (n) OrderItem  
Order (1) ←→ (1) OrderPayment
```

### 📺 media_db - Media Content Database

#### Core Entities
- **Video**: オンデマンド動画
- **Live**: ライブ配信
- **VideoComment**: 動画コメント
- **LiveComment**: ライブコメント

### 💬 messenger_db - Messaging & Notification Database

#### Core Entities
- **Notification**: 通知
- **Message**: メッセージ
- **NotificationTemplate**: 通知テンプレート

## Cross-Module Data Access

### Design Principles
1. **No Direct Database Access**: モジュール間の直接DB接続は禁止
2. **Function-Based Communication**: 必要な場合は関数呼び出しでデータ取得
3. **Event-Driven Updates**: 状態変更時はイベント発行で他モジュールに通知

## Related Documents

- [API Services Overview](./api-services-overview.md)
- [Frontend Applications](./frontend-applications.md)