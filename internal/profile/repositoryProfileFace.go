package profile

type ProfileRepository interface {
}

//Функционал:
//✅ Полный профиль пользователя с настройками
//✅ Смена пароля с валидацией
//✅ Верификация email
//✅ Запрос на сброс пароля
//✅ Деактивация аккаунта
//✅ История активности пользователя
//✅ Социальные ссылки в профиле type UserRepository interface {
//    // Основные CRUD операции
//    Create(ctx context.Context, user *User) error
//    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
//    FindByEmail(ctx context.Context, email string) (*User, error)
//    Update(ctx context.Context, user *User) error
//    Delete(ctx context.Context, id uuid.UUID) error
//    SoftDelete(ctx context.Context, id uuid.UUID) error
//
//    // Поиск и фильтрация
//    FindAll(ctx context.Context, filter UserFilter, pagination Pagination) ([]*User, int64, error)
//    FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*User, error)
//
//    // Профиль
//    GetProfile(ctx context.Context, userID uuid.UUID) (*Profile, error)
//    UpdateProfile(ctx context.Context, profile *Profile) error
//
//    // Верификация
//    SetEmailVerified(ctx context.Context, userID uuid.UUID) error
//    SetPhoneVerified(ctx context.Context, userID uuid.UUID) error
//    SetEmailVerifyToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
//    SetResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
//
//    // Активность
//    RecordActivity(ctx context.Context, activity *Activity) error
//    GetActivities(ctx context.Context, userID uuid.UUID, limit int) ([]*Activity, error)
//    UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
//    UpdateLastActivity(ctx context.Context, userID uuid.UUID) error
// GET /api/v1/users/profile - получение профиля
// PUT /api/v1/users/profile - обновление профиля
// POST /api/v1/users/avatar - загрузка аватара
// POST /api/v1/users/verify-email - подтверждение email
// POST /api/v1/users/change-password - смена пароля
//POST /api/v1/users/enable-2fa        # Включение 2FA
// Настроим отправку email через SMTP или внешний сервис (например, SendGrid).
