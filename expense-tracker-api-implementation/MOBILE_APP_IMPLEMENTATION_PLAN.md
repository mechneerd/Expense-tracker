# Mobile App Implementation Plan - Expense Tracker

## Overview
Android mobile application (APK) that connects to the Go REST API backend for the family expense tracker.

## App Design Reference
Source: `C:\Users\butom\Downloads\app`
- `Entry.png` - Login/OTP verification screen
- `Home (1).png` - Main dashboard/home screen
- `Categories.png` - Category selection screen
- `Entry.png` - Transaction entry screen
- `iPhone 13 & 14 - 31.png` - iPhone mockup rendering
- `profile.png` - User profile screen

## Technology Stack

### Frontend
- **Language**: Kotlin (recommended) or Java
- **UI Framework**: Android XML layouts + Kotlin activities/fragments
- **Network**: Retrofit2 + GsonConverterFactory
- **Image Loading**: Glide or Picasso
- **Architecture**: MVVM (Model-View-ViewModel) or Clean Architecture

### Backend Connection
- **Base URL**: `http://10.0.2.2:8080` (Android Emulator)
  - `10.0.2.2` maps to `localhost` from the emulator's perspective
- **Physical Device**: Use computer's local IP address (e.g., `http://192.168.1.x:8080`)
- **HTTPS**: Configure ngrok or self-signed cert for production

### Dependencies (build.gradle)
```kotlin
dependencies {
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-gson:2.9.0")
    implementation("com.github.bumptech.glide:glide:4.15.0")
    implementation("androidx.lifecycle:lifecycle-viewmodel-ktx:2.7.0")
    implementation("androidx.lifecycle:lifecycle-livedata-ktx:2.7.0")
    implementation("com.google.android.material:material:1.12.0")
    implementation("com.auth0.android:jwtdecode:1.5.0")
}
```

## API Client Structure

### Retrofit Interface
```kotlin
interface ExpenseApi {
    @GET("/")
    fun home(): Call<HomeResponse>

    // Auth
    @POST("/api/v1/auth/google") fun googleLogin(@Body map: Map<String, Any>): Call<AuthResponse>
    @POST("/api/v1/auth/verify-otp") fun verifyOtp(@Body map: Map<String, Any>): Call<AuthResponse>
    @POST("/api/v1/auth/refresh") fun refreshToken(@Body map: Map<String, Any>): Call<AuthResponse>
    @POST("/api/v1/auth/logout") fun logout(): Call<BaseResponse>

    // Users
    @GET("/api/v1/users/me") fun getMe(): Call<UserResponse>
    @PATCH("/api/v1/users/me") fun updateMe(@Body map: Map<String, Any>): Call<UserResponse>

    // Families
    @POST("/api/v1/families") fun createFamily(@Body map: Map<String, Any>): Call<FamilyResponse>
    @GET("/api/v1/families/me") fun getMyFamily(): Call<FamilyResponse>
    @GET("/api/v1/families/invitations") fun listInvitations(): Call<InvitationResponse>
    @POST("/api/v1/families/invitations") fun sendInvite(@Body map: Map<String, Any>): Call<InvitationResponse>
    @POST("/api/v1/families/invitations/{id}/accept") fun acceptInvite(@Path("id") id: String): Call<BaseResponse>
    @POST("/api/v1/families/invitations/{id}/reject") fun rejectInvite(@Path("id") id: String): Call<BaseResponse>

    // Transactions
    @POST("/api/v1/transactions") fun createTransaction(@Body map: Map<String, Any>): Call<TransactionResponse>
    @GET("/api/v1/transactions") fun listTransactions(): Call<TransactionListResponse>
    @GET("/api/v1/transactions/me") fun listMyTransactions(): Call<TransactionListResponse>

    // Dashboard (Family Head only)
    @GET("/api/v1/dashboard") fun getDashboard(@Query("period") period: String): Call<DashboardResponse>

    // Categories
    @GET("/api/v1/categories") fun listCategories(): Call<CategoryListResponse>

    // Payment Methods
    @GET("/api/v1/payment-methods") fun listPaymentMethods(): Call<PaymentMethodListResponse>

    // UPI Apps
    @GET("/api/v1/upi-apps") fun listUPIApps(): Call<UPIAppListResponse>
}
```

## Screen Implementation Plan

### 1. Splash Screen
- Initial launch screen
- Checks if user is logged in
- Redirects to login or home

### 2. Login / OTP Verification (Entry.png)
- **Screen Components**:
  - Email input field
  - "Continue with Google" button
  - OTP verification section
  - Resend OTP button
- **Flow**:
  1. User enters email / clicks Google
  2. Backend sends OTP to user's email (via email provider)
  3. User enters OTP received
  4. Backend verifies OTP
  5. JWT access token + refresh token stored in SecurePreferences
  6. Navigate to home

### 3. Home Dashboard (Home (1).png)
- **Screen Components**:
  - Family name and unique code
  - Balance summary (Income/Expense/Net)
  - Quick stats cards
  - Navigation tabs: Transactions, Categories, Members, Dashboard
- **Data fetched**:
  - Family info from `/api/v1/families/me`
  - Recent transactions from `/api/v1/transactions/me`
  - Dashboard data from `/api/v1/dashboard` (if family head)

### 4. Transaction Entry Screen
- **Screen Components**:
  - Amount input
  - Date picker
  - Description input
  - Type selector (DEBIT/CREDIT)
  - Category selector
  - Payment method selector
  - UPI App selector (when type = DEBIT + UPI)
  - Save/Cancel buttons
- **Flow**:
  1. User fills in transaction details
  2. Validates required fields
  3. Calls `POST /api/v1/transactions`
  4. On success, refreshes transaction list
  5. On failure, shows error snackbar

### 5. Categories Screen (Categories.png)
- **Screen Components**:
  - List of expense categories
  - Add new category (Admin/Family Head only)
  - Category icons/colors
- **Data fetched**: `GET /api/v1/categories`
- **Actions**:
  - View category expenses
  - Filter transactions by category

### 6. Profile Screen (profile.png)
- **Screen Components**:
  - User name and email
  - Profile image
  - Family role (HEAD/MEMBER)
  - Logout button
- **Data fetched**: `GET /api/v1/users/me`

## Navigation Structure

```
Splash
  └─► (if logged in) LoginOTP
       └─► (if success) MainActivity
            ├─► FragmentContainer
            │   ├─► HomeFragment (dashboard overview)
            │   ├─► TransactionsFragment (list + add)
            │   ├─► CategoriesFragment (list + filter)
            │   ├─► MembersFragment (family members)
            │   └─► DashboardFragment (family head only)
            └─► (if not head) ProfileFragment
                 └─► Logout option
```

## Permission & Settings Configuration

### AndroidManifest.xml
```xml
<manifest ...>
    <uses-permission android:name="android.permission.INTERNET" />
    <uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
    <application ...>
        <activity
            android:name=".SplashActivity"
            android:exported="true">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
        <!-- Other activities -->
    </application>
</manifest>
```

### Network Configuration
- **Emulator**: `http://10.0.2.2:8080`
- **Real device**: Find local IP (`ipconfig` on Windows), use `http://192.168.1.x:8080`
- **HTTPS/Production**: Use ngrok tunnel or deploy to AWS/DigitalOcean

## API Response Handling

### Successful Response (from pkg/response)
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation successful"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message",
  "message": "Optional details"
}
```

### Kotlin Handling
```kotlin
response.execute { 
    // Parse into data class
    // Show success snackbar
    // Update UI
} catch (e: Exception) {
    // Show error snackbar
    // Handle authentication errors (refresh token, re-login)
}
```

## Data Models (Kotlin)

### Transaction Data Class
```kotlin
data class Transaction(
    val id: String,
    val familyId: String,
    val userId: String?,
    val transactionType: String, // "DEBIT" or "CREDIT"
    val categoryId: String?,
    val paymentMethodId: String?,
    val upiAppId: String?,
    val amount: Double,
    val transactionDate: String, // "YYYY-MM-DD"
    val description: String?,
    val createdAt: String,
    val updatedAt: String?,
    val deletedAt: String?
)
```

### Dashboard Data Model
```kotlin
data class DashboardData(
    val summary: DashboardSummary,
    val breakdown: ExpenseBreakdown
)

data class DashboardSummary(
    val totalIncome: Float,
    val totalExpense: Float,
    val balance: Float,
    val thisMonthIncome: Float,
    val thisMonthExpense: Float,
    val lastMonthIncome: Float,
    val lastMonthExpense: Float,
    val period: String
)

data class ExpenseBreakdown(
    val byCategory: Map<String, Float>,
    val byMember: Map<String, Float>,
    val byPaymentMethod: Map<String, Float>
)
```

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
- [ ] Set up new Android Studio project
- [ ] Add Retrofit, Gson, Glide dependencies
- [ ] Create network layer (ApiClient, ExpenseApi interface)
- [ ] Implement SplashActivity + LoginOTP screen
- [ ] Configure base URL for emulator/device
- [ ] Test API connectivity with backend

### Phase 2: Core Features (Weeks 3-5)
- [ ] Implement HomeDashboardFragment
- [ ] Implement TransactionEntry screen
- [ ] Implement Categories screen
- [ ] Implement Profile screen
- [ ] Handle JWT storage (SecurePreferences)
- [ ] Implement refresh token logic
- [ ] Test all API endpoints

### Phase 3: Family Features (Weeks 6-8)
- [ ] Implement Family creation/joining
- [ ] Implement Invitation system (enter head email, send invite)
- [ ] Implement Accept/Reject invitations
- [ ] Implement Members list screen
- [ ] Dashboard for family head only

### Phase 4: Polish & Release (Weeks 9-10)
- [ ] Material Design UI refinement
- [ ] Offline cache handling (Room database)
- [ ] Error handling & edge cases
- [ ] Localization (if needed)
- [ ] Generate release APK
- [ ] Test on both emulator and physical device
- [ ] Submit to Play Store (optional)

## Backend API Checklist for Mobile

### Must-Implement Endpoints
- [ ] `POST /api/v1/auth/google` - Google login
- [ ] `POST /api/v1/auth/verify-otp` - OTP verification
- [ ] `POST /api/v1/auth/refresh` - Token refresh
- [ ] `GET /api/v1/users/me` - User profile
- [ ] `POST /api/v1/families` - Create family
- [ ] `GET /api/v1/families/me` - My family
- [ ] `POST /api/v1/transactions` - Create transaction
- [ ] `GET /api/v1/transactions/me` - My transactions
- [ ] `GET /api/v1/categories` - All categories
- [ ] `GET /api/v1/payment-methods` - Payment methods
- [ ] `GET /api/v1/upi-apps` - UPI apps list

### Nice-to-Have
- [ ] `POST /api/v1/families/invitations` - Send invitation
- [ ] `POST /api/v1/families/invitations/{id}/accept` - Accept
- [ ] `GET /api/v1/dashboard` - Family head dashboard
- [ ] `GET /api/v1/families/invitations` - List pending invites

## Testing Strategy

### Unit Tests
- Retrofit API interface tests
- Response parsing tests
- ViewModel logic tests

### Integration Tests
- End-to-end user flows
- Network failure handling
- Offline mode behavior

### Device Testing
- Emulator: `10.0.2.2:8080`
- Physical device: Computer's local IP
- Test all user flows on real device

## Build Configuration

### gradle.properties
```properties
org.gradle.jvmargs=-Xmx512m
android.enableJetifier=true
android.androidCompileOptions.sourceCompatibility=11
android.androidCompileOptions.targetCompatibility=11
```

### build.gradle (module app)
```kotlin
android {
    compileSdk = 34
    defaultConfig {
        applicationId = "com.expensetracker.family"
        minSdk = 21
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"
    }
    buildFeatures {
        viewBinding = true
        compose = false // Using XML + Kotlin, not Jetpack Compose
    }
}
```

## Dependencies Summary

| Dependency | Purpose |
|------------|---------|
| `retrofit2:retrofit` | HTTP client |
| `retrofit2:converter-gson` | JSON -> data class |
| `com.squareup.okhttp3:logging-interceptor` | Log network requests |
| `com.github.bumptech.glide:glide` | Image loading |
| `androidx.lifecycle:lifecycle-viewmodel-ktx` | MVVM architecture |
| `androidx.lifecycle:lifecycle-livedata-ktx` | LiveData observation |
| `com.google.android.material:material` | Material Design UI |
| `com.auth0.android:jwtdecode` | JWT token decoding |

## Next Steps

1. Accept this implementation plan
2. Create the Android Studio project
3. Begin with Phase 1: Foundation
4. Connect to the Go backend API
5. Iterate through phases and deliver the APK

---
*Implementation based on the Go REST API backend at `expense-tracker-api` and UI designs at `C:\Users\butom\Downloads\app`*