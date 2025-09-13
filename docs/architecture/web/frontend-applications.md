# Frontend Applications Overview

Furumaru プロジェクトのフロントエンドアプリケーション構成とアーキテクチャを定義します。

## Application Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   web/admin     │    │   web/user      │    │   web/liff      │    │  web/shared     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ Nuxt 4          │    │ Nuxt 3          │    │ Nuxt 3          │    │ Vue 3 + Vite    │
│ Vuetify 3       │    │ Tailwind CSS    │    │ Tailwind CSS    │    │ Storybook       │
│ Material Design │    │ SSR Enabled     │    │ LINE LIFF       │    │ Component Lib   │
│ SPA Mode        │    │ SEO Optimized   │    │ SPA Mode        │    │ Testing Utils   │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
         ↓                       ↓                       ↓                       ↑
   Management Tools        E-Commerce Site          LINE Mini App           Shared Components
```

## Detailed Application Specifications

### 👨‍💼 admin - Management Portal

**Target Users**: 管理者、コーディネーター、生産者

**Technology Stack**:
```json
{
  "framework": "Nuxt 4.1.1",
  "ui": "Vuetify 3.9.7",
  "state": "@pinia/nuxt 0.11.2",
  "auth": "Firebase Auth + JWT",
  "rendering": "SPA (SSR disabled)",
  "bundler": "Vite 7.1.4"
}
```

**Key Features**:
- **Rich Editor**: TipTap v2.26.1 for content creation
- **Data Visualization**: Chart.js v4.5.0 + ECharts v5.6.0
- **Real-time Updates**: Firebase integration
- **File Management**: Drag & drop uploads, sortable lists
- **Push Notifications**: Firebase Cloud Messaging
- **Error Monitoring**: Sentry integration

**Directory Structure**:
```
src/
├── components/         # Vuetify-based components
├── layouts/           # Admin layout templates
├── pages/             # Admin pages (users, products, orders, etc.)
├── store/             # Pinia stores (auth, data management)
├── plugins/           # Firebase, Sentry, API client
├── middleware/        # Authentication guards
└── types/             # TypeScript definitions
```

**Authentication Flow**:
```
Login → Firebase Auth → JWT Token → API Gateway → Admin Services
```

### 🛒 user - E-Commerce Portal

**Target Users**: 購入者（会員・ゲスト・施設利用者）

**Technology Stack**:
```json
{
  "framework": "Nuxt 3",
  "ui": "Tailwind CSS + Custom Design System",
  "state": "@pinia/nuxt 0.5.1 + Persisted State",
  "auth": "Cookie Session + JWT Bearer",
  "rendering": "SSR + Hydration",
  "seo": "Meta tags + Open Graph",
  "maps": "Google Maps Services JS v3.4.0"
}
```

**Key Features**:
- **E-Commerce**: Product catalog, cart, checkout, order management
- **Live Commerce**: HLS.js v1.4.10 for video streaming
- **Multi-language**: i18n support
- **SEO Optimization**: SSR + structured data
- **Payment Integration**: Multiple payment providers
- **Location Services**: Google Maps integration
- **Content Management**: microCMS integration

**Directory Structure**:
```
src/
├── components/         # Tailwind-based components
├── pages/             # E-commerce pages (products, cart, checkout)
├── composables/       # Business logic hooks
├── stores/            # Cart, user, product states
├── middleware/        # Auth, routing guards  
├── plugins/           # API client, i18n, analytics
└── types/             # API response types
```

**Rendering Strategy**:
- **SSR**: Product pages, category pages (SEO critical)
- **SPA**: User dashboard, cart, checkout (interactive)

### 📱 liff - LINE Mini App

**Target Users**: LINEユーザー（チャット統合購入体験）

**Technology Stack**:
```json
{
  "framework": "Nuxt 3",
  "ui": "Tailwind CSS",
  "state": "@pinia/nuxt 0.11.2",
  "liff": "@line/liff v2.27.1",
  "auth": "LINE OAuth + Session",
  "rendering": "SPA (optimized for mobile)"
}
```

**Key Features**:
- **LINE Integration**: LIFF SDK for native LINE features
- **Seamless Auth**: LINE user profile integration
- **Optimized UX**: Mobile-first, touch-optimized
- **Chat Integration**: Purchase sharing to LINE chats
- **Lightweight**: Minimal bundle size for fast loading

**Directory Structure**:
```
src/
├── components/         # Mobile-optimized components
├── pages/             # LIFF-specific user flows
├── composables/       # LINE SDK integrations
├── stores/            # User session, cart sync
└── plugins/           # LIFF initialization
```

**LIFF Integration Pattern**:
```javascript
// LIFF initialization
await liff.init({ liffId: process.env.LIFF_ID })
const profile = await liff.getProfile()
// Sync with backend user system
```

### 🧩 shared - Component Library

**Purpose**: デザインシステム・共通コンポーネント・開発効率化

**Technology Stack**:
```json
{
  "framework": "Vue 3 + TypeScript",
  "bundler": "Vite",
  "documentation": "Storybook",
  "testing": "Vitest + Vue Test Utils",
  "build": "ESM + UMD outputs"
}
```

**Component Categories**:
- **Atoms**: Button, Input, Icon, Badge
- **Molecules**: SearchBox, ProductCard, UserAvatar
- **Organisms**: Header, Footer, ProductGrid
- **Templates**: PageLayout, FormLayout
- **Utilities**: Composables, Type definitions

**Design System Features**:
- **Theme Support**: CSS Custom Properties
- **Responsive Design**: Mobile-first breakpoints
- **Accessibility**: ARIA compliance, keyboard navigation
- **Icon System**: SVG sprite optimization
- **Animation**: CSS transitions + Vue transition components

## Cross-Application Architecture

### State Management Strategy

**Pinia Store Organization**:
```typescript
// Shared store structure across apps
interface StoreStructure {
  auth: AuthStore      // User authentication state
  cart: CartStore      // Shopping cart state  
  ui: UIStore         // Global UI state
  user: UserStore     // User profile data
  [domain]: DomainStore // Domain-specific stores
}
```

**State Persistence**:
- **admin**: Session-based (no persistence)
- **user**: `@pinia-plugin-persistedstate` for cart/preferences
- **liff**: localStorage for session continuity

### API Integration Pattern

**Unified API Client**:
```typescript
// api-client plugin (shared pattern)
interface APIClient {
  auth: AuthAPI
  user: UserAPI  
  store: StoreAPI
  media: MediaAPI
  messenger: MessengerAPI
}
```

**Request Flow**:
```
Component → Composable → API Client → HTTP Client → Gateway → Service
    ↓           ↓            ↓            ↓           ↓         ↓
  UI Logic   Business    Auth/Cache   axios/fetch   Auth     gRPC
             Logic       Handling     HTTP Layer   Validation
```

### Authentication Architecture

**Multi-App Authentication**:
```
┌─── admin ────┐    ┌─── user ────┐    ┌─── liff ────┐
│ Firebase     │    │ Cookie +    │    │ LINE OAuth  │
│ JWT Bearer   │    │ JWT Bearer  │    │ + Session   │
└──────────────┘    └─────────────┘    └─────────────┘
        ↓                   ↓                  ↓
   Gateway/admin      Gateway/user      Gateway/user
        ↓                   ↓                  ↓
    Admin APIs         User APIs          User APIs
```

### Deployment Architecture

**Build Strategy**:
```yaml
# Each app builds independently
admin:
  build: "nuxt build"
  output: ".output/"
  deploy: "S3 + CloudFront"

user:
  build: "nuxt build"
  output: ".output/"
  deploy: "S3 + CloudFront"

liff:
  build: "nuxt build"  
  output: ".output/"
  deploy: "S3 + CloudFront"

shared:
  build: "vite build"
  output: "dist/"
  publish: "npm registry"
```

## Performance Considerations

### Bundle Optimization

**Code Splitting Strategy**:
- **Route-based**: Automatic page splitting (Nuxt)
- **Component-based**: Dynamic imports for heavy components
- **Vendor-based**: Third-party library separation

**Tree Shaking**:
- **Vuetify**: Component-based imports (admin)
- **Tailwind**: PurgeCSS for unused styles
- **Shared**: Selective component imports

### Caching Strategy

**Browser Caching**:
```
Assets (images, fonts): Cache-Control: max-age=31536000
JavaScript/CSS: Cache-Control: max-age=3600, stale-while-revalidate=86400
HTML: Cache-Control: no-cache
```

**CDN Strategy**:
- **Static Assets**: CloudFront + S3
- **API Responses**: API Gateway caching
- **Images**: Optimized delivery + lazy loading

## Development Workflow

### Development Commands

```bash
# Individual app development
cd web/admin && yarn dev     # http://localhost:3000
cd web/user && yarn dev      # https://localhost:3000 (SSL)
cd web/liff && yarn dev      # http://localhost:3001

# Shared library development
cd web/shared && yarn storybook  # http://localhost:6006

# Build all applications
make build-web  # Build all web apps
```

### Testing Strategy

**Unit Testing**:
- **Framework**: Vitest + Vue Test Utils
- **Coverage**: 80%+ for composables and utilities
- **Mocking**: API responses, external services

**Integration Testing**:
- **E2E**: Playwright for critical user flows
- **Component**: Storybook interaction testing
- **API**: Contract testing with backend

### Quality Assurance

**Code Quality Tools**:
```json
{
  "linting": "ESLint + TypeScript rules",
  "formatting": "Prettier",
  "css": "Stylelint (admin only)",
  "types": "TypeScript strict mode",
  "bundleAnalysis": "vite-bundle-analyzer"
}
```

**CI/CD Pipeline**:
1. **Lint & Type Check**: All apps
2. **Unit Tests**: Run test suites
3. **Build**: Generate production builds
4. **Deploy**: S3 + CloudFront invalidation

## Security Considerations

### XSS Protection
- **Vue.js**: Automatic template escaping
- **CSP**: Content Security Policy headers
- **Sanitization**: DOMPurify for rich content

### Authentication Security
- **JWT Validation**: Token expiry + refresh flow
- **CSRF Protection**: SameSite cookies
- **Session Management**: Secure cookie flags

### Data Protection
- **Environment Variables**: Sensitive config management
- **Error Handling**: No sensitive data in error messages
- **Logging**: Client-side error reporting (Sentry)

## Related Documents

- [Web Architecture README](../../architecture/web/README.md)
- [Component Design](../../architecture/web/components.md)
- [State Management](../../architecture/web/state-management.md)
- [API Integration](../../architecture/web/api-integration.md)