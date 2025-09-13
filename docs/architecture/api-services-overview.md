# API Services Overview

Furumaru API の各モジュールの概要と責務を定義します。

## Service Architecture

```
Modular Monolith Architecture
├── cmd/                    # Entry points
│   ├── gateway/           # API Gateway
│   ├── media/             # Media service
│   └── messenger/         # Messaging service
├── internal/              # Business modules (Layered Architecture)
│   ├── gateway/           # API Gateway module
│   ├── user/              # User management & authentication
│   ├── store/             # E-commerce functionality
│   ├── media/             # Media management
│   └── messenger/         # Messaging & notifications
└── pkg/                   # Shared packages

Internal Communication: Direct Go function calls
```

## Module Details

### 🚪 Gateway Module
- **Responsibility**: API endpoint aggregation and routing
- **Authentication**: JWT Bearer Token / Cookie Session
- **Main Features**: 
  - Request routing to internal modules
  - Authentication & authorization
  - Response formatting
  - Cross-cutting concerns

### 🔐 User Module  
- **Responsibility**: User management and authentication
- **Database**: user_db
- **Main Entities**:
  - User: Customers (members/guests/facility users)
  - Admin: System administrators
  - Coordinator: Sales coordinators
  - Producer: Product producers
- **Main Features**:
  - OAuth authentication (Google, LINE)
  - JWT token management
  - User profile management
  - Address & notification settings

### 🛒 Store Module
- **Responsibility**: E-commerce core functionality
- **Database**: store_db
- **Main Features**:
  - Product catalog management
  - Order & payment processing
  - Shipping management
  - Inventory management
  - Cart functionality

### 📺 Media Module
- **Responsibility**: Video streaming and media management
- **Database**: media_db  
- **Main Features**:
  - Live streaming management
  - Video streaming (HLS)
  - Media file management
  - Broadcasting scheduling

### 💬 Messenger Module
- **Responsibility**: Notifications and messaging
- **Database**: messenger_db
- **Main Features**:
  - Push notifications (FCM)
  - Email notifications
  - SMS notifications  
  - Message queue processing
- **Background Workers**:
  - scheduler: Periodic tasks
  - worker: Asynchronous processing

## Technology Stack

- **Language**: Go 1.25.1
- **HTTP Framework**: Gin  
- **Internal Communication**: Direct function calls
- **ORM**: GORM
- **Database**: TiDB (MySQL compatible)
- **Authentication**: AWS Cognito + JWT
- **Message Queue**: AWS SQS
- **Notification**: Firebase Cloud Messaging

## Related Documents

- [Database Design](./database-design.md)
- [Frontend Applications](./frontend-applications.md)
- [Documentation Patterns](./documentation-patterns.md)