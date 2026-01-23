# FoodHive Mobile Sales App - Feature Proposal

## Overview

This document outlines two implementation approaches for a mobile sales application that enables field sales representatives to create orders directly at customer locations while syncing with the FoodHive ERP system.

---

## Option A: Basic Sales App (MVP)

**Scope**: Essential features for field ordering  
**Complexity**: Low-Medium

### Core Features

#### 1. Authentication & Security
- Login with employee credentials (synced with ERP)
- Session management with auto-logout
- PIN/Biometric quick login for returning users
- Remember me functionality

#### 2. Customer Management
- View assigned customer list
- Customer search by name, code, or phone
- View customer details (address, contact, credit terms)
- Navigate to customer location (Google Maps / Apple Maps integration)
- View customer's outstanding balance (read-only)

#### 3. Product Catalog
- Browse products by category
- Product search functionality
- View product details (name, SKU, unit, description)
- Display customer-specific pricing tiers automatically
- Show product images
- Filter by category, availability

#### 4. Order Creation
- Add products to cart with quantity
- Adjust quantities (+/- buttons and manual entry)
- Apply customer's pricing tier automatically
- Order notes/special instructions field
- Delivery date selection
- Submit order to ERP system
- Order confirmation screen

#### 5. Order History
- View past orders for each customer
- Order status tracking (Pending, Confirmed, Processing, Delivered)
- Reorder from previous orders (quick order feature)
- View order details and line items

#### 6. Basic Offline Mode
- Cache customer list locally
- Cache product catalog locally
- Queue orders when offline (stored locally)
- Auto-sync when connection restored
- Visual indicator for pending sync items

#### 7. Profile & Settings
- View own profile information
- Change password
- Logout functionality
- App version info

### Technical Requirements

**Suggested Stack:**
- **Framework**: React Native or Flutter (cross-platform)
- **State Management**: Redux Toolkit / Provider
- **Local Storage**: SQLite / Realm
- **API**: REST (existing FoodHive backend endpoints)
- **Maps**: Native Maps SDK

**Required Backend Endpoints:**
```
Authentication
├── POST   /api/v1/mobile/auth/login
├── POST   /api/v1/mobile/auth/refresh
└── POST   /api/v1/mobile/auth/logout

Customers
├── GET    /api/v1/mobile/customers           (list assigned customers)
├── GET    /api/v1/mobile/customers/:id       (customer details)
└── GET    /api/v1/mobile/customers/:id/balance

Products
├── GET    /api/v1/mobile/products            (with pagination)
├── GET    /api/v1/mobile/products/:id
├── GET    /api/v1/mobile/categories
└── GET    /api/v1/mobile/pricing/:customerId (customer-specific prices)

Orders
├── POST   /api/v1/mobile/orders              (create order)
├── GET    /api/v1/mobile/orders              (order history)
└── GET    /api/v1/mobile/orders/:id          (order details)
```

### Screen List (10 Screens)
1. **Login Screen** - Email/password authentication
2. **Dashboard/Home** - Quick stats, pending orders, quick actions
3. **Customer List** - Searchable list with filters
4. **Customer Detail** - Full customer info with action buttons
5. **Product Catalog** - Grid/list view with categories
6. **Product Detail** - Full product information
7. **Shopping Cart** - Review items, adjust quantities
8. **Order Confirmation** - Summary before submission
9. **Order History** - List of past orders with status
10. **Profile/Settings** - User info and app settings

---

## Option B: Full-Featured Sales App (Enterprise)

**Scope**: Complete field sales solution  
**Complexity**: High

### Includes Everything from Basic Version PLUS:

---

### 1. Enhanced Authentication & Security
- Role-based access control (Sales Rep, Senior Rep, Supervisor, Manager)
- Multi-factor authentication (OTP via SMS/Email)
- Remote device wipe capability
- Activity logging for compliance and audit
- Device registration/verification
- Session management across multiple devices

---

### 2. Real-Time Communication Hub

#### Team Chat System
| Feature | Description |
|---------|-------------|
| Direct Messages | One-on-one chat with any team member |
| Group Chats | Team channels, regional groups |
| Read Receipts | Know when messages are read |
| File Sharing | Share images, PDFs, documents |
| Voice Messages | Record and send audio clips |
| Message Search | Search through chat history |
| Offline Queue | Messages sent when back online |

#### Supervisor-Rep Communication
- Priority messaging from supervisors
- Broadcast announcements to all reps
- Request approval for special discounts
- Ask questions and get real-time answers
- Share customer photos/situations

#### Push Notifications
- New order assignments
- Order status changes
- Chat messages (with preview)
- Daily targets and reminders
- Promotional announcements
- Approval requests/responses
- System alerts

---

### 3. Advanced Customer Management

#### Visit Management
- **Check-in**: GPS-verified arrival at customer location
- **Check-out**: Automatic visit duration calculation
- **Visit Notes**: Record observations and feedback
- **Visit Photos**: Capture storefront, shelf displays
- **Visit History**: View all past visits with details

#### Customer Intelligence
- Customer credit limit display with visual indicator
- Outstanding balance with aging breakdown
- Payment behavior history
- Last order date and amount
- Preferred products list
- Customer-specific notes
- Contact multiple people per customer

#### Digital Signatures
- Capture customer signature for orders
- Signature for delivery confirmation
- Signature for payment receipts
- Stored securely with timestamp

---

### 4. Enhanced Product Features

#### Inventory Integration
- Real-time stock levels by warehouse
- Available quantity calculation
- Low stock warnings before ordering
- Out of stock indicators
- Alternative product suggestions when out of stock

#### Product Discovery
- **Barcode Scanner**: Quick product lookup via camera
- **Voice Search**: Search products by speaking
- **Recent Products**: Quick access to frequently ordered items
- **Favorites**: Mark products as favorites per customer
- **Product Comparison**: Compare similar products

#### Promotions & Pricing
- Active promotions display
- Volume discounts auto-calculation
- Bundle deals visualization
- Limited-time offers with countdown
- Customer-specific special prices

---

### 5. Advanced Order Management

#### Enhanced Order Creation
| Feature | Description |
|---------|-------------|
| Multiple Price Tiers | Select from available pricing levels |
| Catch Weight Entry | Enter actual weight for weighted products |
| Volume Discounts | Auto-apply quantity-based discounts |
| Promotional Pricing | Apply active promotions |
| Order Templates | Save and reuse frequent orders |
| Split Delivery | Multiple delivery dates per order |
| Delivery Instructions | Specific delivery notes |
| Credit Check | Warn if order exceeds credit limit |

#### Order Workflow
- Draft orders (save for later)
- Order approval workflow (for high-value orders)
- Edit pending orders
- Cancel orders (with reason)
- Duplicate previous orders
- Combine orders

#### Returns & Issues
- Create return requests
- Report delivery issues
- Quality complaints with photos
- Track return status

---

### 6. Payment Collection

#### Payment Recording
| Payment Type | Features |
|--------------|----------|
| Cash | Amount entry, denomination breakdown |
| Cheque | Cheque number, bank, date capture |
| Bank Transfer | Reference number recording |
| Mobile Payment | QR code, reference capture |

#### Payment Features
- Link payment to specific invoices
- Partial payment support
- Generate payment receipt (PDF)
- Share receipt via WhatsApp/Email
- Payment history per customer
- Daily collection summary
- Outstanding aging report

---

### 7. GPS & Route Management

#### Location Tracking
| Feature | Description |
|---------|-------------|
| Real-time GPS | Supervisor can see rep locations on map |
| Geofence Visits | Auto-detect customer visits |
| Travel Distance | Daily/weekly distance reports |
| Time at Location | Track time spent per customer |
| Route History | View past routes traveled |

#### Route Planning & Optimization
- Daily customer visit list
- Optimized route suggestions
- Drag-and-drop route reordering
- Turn-by-turn navigation
- Traffic-aware routing
- Multi-stop optimization
- Estimated arrival times
- Route sharing with customers

#### Supervisor Map View
- See all reps on map in real-time
- Filter by team/region
- Click to see rep details and activity
- Send location-based tasks
- Monitor coverage areas

---

### 8. Reporting & Analytics

#### Sales Rep Dashboard
```
┌─────────────────────────────────────────────────┐
│  Today's Performance                            │
├─────────────────────────────────────────────────┤
│  Orders: 8/12 target    ████████░░░░  67%       │
│  Revenue: 45M/60M LAK   ███████░░░░░  75%       │
│  Visits: 10/15          ██████████░░  67%       │
│  Collections: 25M LAK   ████████████  100%      │
└─────────────────────────────────────────────────┘
```

#### Reports Available
- Daily sales summary
- Weekly performance trends
- Monthly achievement vs target
- Top selling products (personal)
- Customer order frequency
- Visit productivity metrics
- Collection efficiency
- Commission tracking (if applicable)

#### Performance Comparison
- Personal trend over time
- Rank among team members
- Regional leaderboard
- Achievement badges/recognition

---

### 9. Supervisor Features

#### Team Management Dashboard
- Real-time team overview
- Individual rep performance cards
- Attendance tracking
- Activity timeline per rep

#### Approval Workflow
| Approval Type | Threshold Example |
|---------------|-------------------|
| Discount > 15% | Requires supervisor OK |
| Order > 50M LAK | Requires supervisor OK |
| New credit terms | Requires manager OK |
| Return requests | Requires supervisor OK |

#### Broadcast & Communication
- Send announcements to all/selected reps
- Schedule messages
- Attach files/images
- Track message read status

#### Reports & Analytics
- Team performance summary
- Regional sales overview
- Route coverage analysis
- Customer visit compliance
- Export reports to Excel/PDF

---

### 10. Advanced Offline Capabilities

| Feature | Basic | Full |
|---------|-------|------|
| Customer cache | ✓ | ✓ + with more details |
| Product cache | ✓ | ✓ + with images |
| Order queue | ✓ | ✓ + with full validation |
| Payment recording | - | ✓ |
| Visit recording | - | ✓ |
| Chat queue | - | ✓ |
| Conflict resolution | Simple | Intelligent merge |
| Background sync | Basic | Priority-based |
| Sync status | Simple | Detailed progress |

---

### 11. Additional Features

#### Document Management
- View customer invoices
- Download/share statements
- Share documents via WhatsApp/Email/SMS
- View delivery notes
- Access to price lists

#### Visit Scheduling
- Calendar view of planned visits
- Create/edit visit schedules
- Recurring visit patterns
- Visit reminders (push notifications)
- Reschedule with customer notification
- Sync with device calendar

#### Survey & Market Intelligence
- Customer satisfaction surveys
- Competitor product spotting
- Market price capture
- New product opportunity notes
- Shelf share photos (merchandising)
- Customer feedback collection

#### Data Export
- Export orders to PDF
- Export customer list to CSV
- Share reports via email
- Print support (Bluetooth printers)

---

## Feature Comparison Matrix

| Category | Feature | Basic | Full |
|----------|---------|:-----:|:----:|
| **Auth** | Login/Password | ✓ | ✓ |
| | Biometric Login | ✓ | ✓ |
| | Role-based Access | - | ✓ |
| | Multi-factor Auth | - | ✓ |
| | Remote Wipe | - | ✓ |
| **Customers** | Customer List | ✓ | ✓ |
| | Customer Search | ✓ | ✓ |
| | Navigation to Customer | ✓ | ✓ |
| | Customer Balance | View Only | Full Details |
| | Visit Check-in/out | - | ✓ |
| | Customer Signatures | - | ✓ |
| | Customer Photos | - | ✓ |
| **Products** | Product Catalog | ✓ | ✓ |
| | Product Search | ✓ | ✓ |
| | Category Filter | ✓ | ✓ |
| | Barcode Scanner | - | ✓ |
| | Real-time Stock | - | ✓ |
| | Voice Search | - | ✓ |
| **Orders** | Create Order | ✓ | ✓ |
| | Customer Pricing | ✓ | ✓ |
| | Order History | ✓ | ✓ |
| | Quick Reorder | ✓ | ✓ |
| | Promotions | - | ✓ |
| | Catch Weight | - | ✓ |
| | Order Templates | - | ✓ |
| | Draft Orders | - | ✓ |
| | Approval Workflow | - | ✓ |
| **Payments** | View Balance | ✓ | ✓ |
| | Record Payment | - | ✓ |
| | Payment Receipt | - | ✓ |
| | Multiple Payment Types | - | ✓ |
| **Communication** | Push Notifications | Basic | Full |
| | Team Chat | - | ✓ |
| | Group Messaging | - | ✓ |
| | Voice Messages | - | ✓ |
| | File Sharing | - | ✓ |
| **Location** | Map Navigation | ✓ | ✓ |
| | GPS Tracking | - | ✓ |
| | Route Optimization | - | ✓ |
| | Geofence Detection | - | ✓ |
| **Reports** | Order Status | ✓ | ✓ |
| | Sales Dashboard | - | ✓ |
| | Performance Metrics | - | ✓ |
| | Team Leaderboard | - | ✓ |
| **Offline** | Basic Caching | ✓ | ✓ |
| | Order Queue | ✓ | ✓ |
| | Full Offline Mode | - | ✓ |
| | Smart Sync | - | ✓ |
| **Supervisor** | Team Map View | - | ✓ |
| | Approval Queue | - | ✓ |
| | Team Analytics | - | ✓ |
| | Broadcast Messages | - | ✓ |

---

## Integration Architecture

### Basic Version
```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Mobile     │  REST   │   FoodHive   │         │   Database   │
│    App       │◄───────►│   Backend    │◄───────►│  PostgreSQL  │
└──────────────┘  HTTPS  └──────────────┘         └──────────────┘
```

### Full Version
```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│              │  REST   │   FoodHive   │         │   Database   │
│   Mobile     │◄───────►│   Backend    │◄───────►│  PostgreSQL  │
│    App       │         └──────────────┘         └──────────────┘
│              │         ┌──────────────┐         ┌──────────────┐
│              │WebSocket│  Chat/Push   │         │    Redis     │
│              │◄───────►│   Server     │◄───────►│   (Cache)    │
└──────────────┘         └──────────────┘         └──────────────┘
                         ┌──────────────┐
                         │   Firebase   │
                         │  (FCM/Maps)  │
                         └──────────────┘
```

### Data Sync Strategy

| Data Type | Sync Direction | Frequency | Priority |
|-----------|----------------|-----------|----------|
| Customers | ERP → App | Daily + On-demand | High |
| Products | ERP → App | Daily + On-demand | High |
| Pricing | ERP → App | Real-time | Critical |
| Inventory | ERP → App | Real-time (Full) | High |
| Orders | App → ERP | Immediate | Critical |
| Payments | App → ERP | Immediate | Critical |
| Visits | App → ERP | Real-time (Full) | Medium |
| Chat | Bidirectional | Real-time | Medium |
| Location | App → Server | Every 5 min (Full) | Low |

---

## Security Considerations

### Both Versions
- HTTPS/TLS 1.3 for all API communication
- JWT token-based authentication
- Token refresh mechanism (short-lived access tokens)
- Secure local storage (Keychain/Keystore)
- No sensitive data in logs
- Input validation on all forms

### Full Version (Additional)
- Certificate pinning to prevent MITM attacks
- Device registration and verification
- Remote session termination capability
- Data encryption at rest (AES-256)
- Comprehensive audit logging
- GDPR/privacy compliance for location data
- Biometric authentication for sensitive actions
- Auto-lock after inactivity
- Jailbreak/root detection

---

## Recommended Implementation Approach

### Phase 1: Basic Version (MVP)
**Goal**: Validate core workflow and gather feedback

**Benefits of starting with Basic:**
1. Faster time to market
2. Lower initial investment
3. Early user feedback
4. Establish integration patterns
5. Train sales team gradually

**Success Criteria:**
- Sales reps can create orders in the field
- Orders sync correctly with ERP
- Customer data is accurate and up-to-date
- 80% of orders placed via mobile within 3 months

### Phase 2: Communication & Tracking
**Add after successful Phase 1 deployment:**
- Team chat functionality
- Push notifications
- GPS tracking (optional opt-in)
- Basic route planning

### Phase 3: Full Enterprise Features
**Complete the full feature set:**
- Payment collection
- Advanced reporting
- Supervisor dashboards
- Full offline capabilities
- Visit scheduling and surveys

---

## Questions to Consider Before Development

### Business Questions
1. How many sales representatives will use the app?
2. What is the average number of customers per sales rep?
3. How many orders does a rep typically create per day?
4. Is payment collection in the field required immediately?
5. What are the current pain points in the field sales process?

### Technical Questions
1. **Devices**: Company-provided or BYOD (Bring Your Own Device)?
2. **Platforms**: iOS only, Android only, or both?
3. **Connectivity**: Do sales reps often work in areas with poor internet?
4. **Existing Systems**: Any third-party systems to integrate?
5. **Data Volume**: Expected size of product catalog and customer base?

### Compliance Questions
1. Any regulatory requirements for data handling (PDPA, etc.)?
2. GPS tracking consent requirements?
3. Data retention policies?
4. Audit trail requirements?

### Operational Questions
1. Who will manage the app (IT team, vendor)?
2. How will updates be deployed?
3. What is the support model for field issues?
4. Training plan for sales team?

---

## Estimated Effort (High Level)

| Component | Basic | Full |
|-----------|-------|------|
| UI/UX Design | 2-3 weeks | 4-6 weeks |
| Mobile App Development | 6-8 weeks | 12-16 weeks |
| Backend API Development | 2-3 weeks | 6-8 weeks |
| Chat Server (Full only) | - | 3-4 weeks |
| Testing & QA | 2-3 weeks | 4-6 weeks |
| Deployment & Launch | 1 week | 2 weeks |
| **Total** | **13-18 weeks** | **31-42 weeks** |

*Note: Estimates assume dedicated development team. Actual timelines depend on team size and availability.*

---

## Next Steps

1. **Review this proposal** with stakeholders
2. **Decide on version** (Basic or Full, or phased approach)
3. **Prioritize features** if customization needed
4. **Technical feasibility review** of existing backend APIs
5. **UI/UX design phase** with mockups
6. **Development kickoff** with sprint planning

---

*Document prepared for FoodHive ERP Mobile Sales App evaluation*  
*Version: 1.0*  
*Last updated: January 2026*
