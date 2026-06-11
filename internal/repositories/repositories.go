package repositories

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repositories struct {
	AdminUsers *AdminUserRepository
	Tickets    *TicketRepository
	Bookings   *BookingRepository
	Messages   *ContactMessageRepository
	Settings   *SiteSettingRepository
	AuditLogs  *AuditLogRepository
	HeroSlides *HeroSlideRepository
	Media      *MediaAssetRepository
	SEO        *SEORepository
	Gallery    *GalleryRepository
	Restaurant *RestaurantRepository
	Suites     *SuiteRepository
	Halls      *HallRepository
	Offers     *OfferRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		AdminUsers: NewAdminUserRepository(db),
		Tickets:    NewTicketRepository(db),
		Bookings:   NewBookingRepository(db),
		Messages:   NewContactMessageRepository(db),
		Settings:   NewSiteSettingRepository(db),
		AuditLogs:  NewAuditLogRepository(db),
		HeroSlides: NewHeroSlideRepository(db),
		Media:      NewMediaAssetRepository(db),
		SEO:        NewSEORepository(db),
		Gallery:    NewGalleryRepository(db),
		Restaurant: NewRestaurantRepository(db),
		Suites:     NewSuiteRepository(db),
		Halls:      NewHallRepository(db),
		Offers:     NewOfferRepository(db),
	}
}

type AdminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

func (r *AdminUserRepository) FindByEmail(ctx context.Context, email string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *AdminUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *AdminUserRepository) Create(ctx context.Context, user *models.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *AdminUserRepository) Save(ctx context.Context, user *models.AdminUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *AdminUserRepository) List(ctx context.Context, search string, page int, limit int) ([]models.AdminUser, int64, error) {
	var users []models.AdminUser
	query := r.db.WithContext(ctx).Model(&models.AdminUser{})
	if search != "" {
		term := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", term, term)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&users).Error
	return users, total, err
}

func (r *AdminUserRepository) CountSuperAdmins(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.AdminUser{}).
		Where("role = ? AND is_active = ?", models.RoleSuperAdmin, true).
		Count(&total).Error
	return total, err
}

type TicketFilter struct {
	Search string
	Status string
}

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) ListPublic(ctx context.Context) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, created_at ASC").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) FindBySlug(ctx context.Context, slug string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.WithContext(ctx).
		Where("slug = ? AND is_active = ?", slug, true).
		First(&ticket).Error
	return &ticket, err
}

func (r *TicketRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error
	return &ticket, err
}

func (r *TicketRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&ticket, "id = ?", id).Error
	return &ticket, err
}

func (r *TicketRepository) ListAdmin(ctx context.Context, filter TicketFilter, page int, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	query := r.db.WithContext(ctx).Model(&models.Ticket{})
	if filter.Search != "" {
		term := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR slug ILIKE ? OR category ILIKE ?", term, term, term)
	}
	if filter.Status == "active" {
		query = query.Where("is_active = ?", true)
	}
	if filter.Status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort_order ASC, created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&tickets).Error
	return tickets, total, err
}

func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *TicketRepository) Save(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *TicketRepository) Delete(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Delete(ticket).Error
}

type BookingFilter struct {
	Search        string
	Status        string
	PaymentStatus string
	DateFrom      *time.Time
	DateTo        *time.Time
}

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *BookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).Preload("Ticket").First(&booking, "id = ?", id).Error
	return &booking, err
}

func (r *BookingRepository) FindByRazorpayOrderIDForUpdate(ctx context.Context, tx *gorm.DB, orderID string) (*models.Booking, error) {
	var booking models.Booking
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Ticket").
		Where("razorpay_order_id = ?", orderID).
		First(&booking).Error
	return &booking, err
}

func (r *BookingRepository) SaveWithTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	return tx.WithContext(ctx).Save(booking).Error
}

func (r *BookingRepository) ListAdmin(ctx context.Context, filter BookingFilter, page int, limit int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	query := r.db.WithContext(ctx).Model(&models.Booking{}).Preload("Ticket")
	if filter.Search != "" {
		term := "%" + filter.Search + "%"
		query = query.Where("booking_id ILIKE ? OR customer_name ILIKE ? OR customer_email ILIKE ? OR customer_phone ILIKE ?", term, term, term, term)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}
	if filter.DateFrom != nil {
		query = query.Where("visit_date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("visit_date <= ?", *filter.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Booking{}).Count(&total).Error
	return total, err
}

func (r *BookingRepository) Revenue(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("payment_status = ?", models.PaymentPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *BookingRepository) Recent(ctx context.Context, limit int) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.WithContext(ctx).
		Preload("Ticket").
		Order("created_at DESC").
		Limit(limit).
		Find(&bookings).Error
	return bookings, err
}

type ChartPoint struct {
	Label string `json:"label"`
	Total int64  `json:"total"`
}

func (r *BookingRepository) RevenueChart(ctx context.Context, days int) ([]ChartPoint, error) {
	var points []ChartPoint
	err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Select("TO_CHAR(DATE(created_at), 'YYYY-MM-DD') AS label, COALESCE(SUM(amount), 0) AS total").
		Where("payment_status = ? AND created_at >= ?", models.PaymentPaid, time.Now().AddDate(0, 0, -days)).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&points).Error
	return points, err
}

func (r *BookingRepository) BookingGrowth(ctx context.Context, days int) ([]ChartPoint, error) {
	var points []ChartPoint
	err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Select("TO_CHAR(DATE(created_at), 'YYYY-MM-DD') AS label, COUNT(*) AS total").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -days)).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&points).Error
	return points, err
}

func (r *BookingRepository) DistinctCustomers(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Booking{}).Distinct("customer_email").Count(&total).Error
	return total, err
}

type ContactMessageFilter struct {
	Search string
	Status string
}

type ContactMessageRepository struct {
	db *gorm.DB
}

func NewContactMessageRepository(db *gorm.DB) *ContactMessageRepository {
	return &ContactMessageRepository{db: db}
}

func (r *ContactMessageRepository) Create(ctx context.Context, message *models.ContactMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *ContactMessageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ContactMessage, error) {
	var message models.ContactMessage
	err := r.db.WithContext(ctx).First(&message, "id = ?", id).Error
	return &message, err
}

func (r *ContactMessageRepository) Save(ctx context.Context, message *models.ContactMessage) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *ContactMessageRepository) Delete(ctx context.Context, message *models.ContactMessage) error {
	return r.db.WithContext(ctx).Delete(message).Error
}

func (r *ContactMessageRepository) ListAdmin(ctx context.Context, filter ContactMessageFilter, page int, limit int) ([]models.ContactMessage, int64, error) {
	var messages []models.ContactMessage
	query := r.db.WithContext(ctx).Model(&models.ContactMessage{})
	if filter.Search != "" {
		term := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR subject ILIKE ? OR message ILIKE ?", term, term, term, term, term)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&messages).Error
	return messages, total, err
}

func (r *ContactMessageRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.ContactMessage{}).Count(&total).Error
	return total, err
}

type SiteSettingRepository struct {
	db *gorm.DB
}

func NewSiteSettingRepository(db *gorm.DB) *SiteSettingRepository {
	return &SiteSettingRepository{db: db}
}

func (r *SiteSettingRepository) First(ctx context.Context) (*models.SiteSetting, error) {
	var setting models.SiteSetting
	err := r.db.WithContext(ctx).Order("created_at ASC").First(&setting).Error
	return &setting, err
}

func (r *SiteSettingRepository) Save(ctx context.Context, setting *models.SiteSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *SiteSettingRepository) Create(ctx context.Context, setting *models.SiteSetting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}

type AuditLogRepository struct {
	db *gorm.DB
}

type AuditLogFilter struct {
	Action  string
	Module  string
	AdminID *uuid.UUID
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditLogRepository) List(ctx context.Context, filter AuditLogFilter, page int, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	query := r.db.WithContext(ctx).Model(&models.AuditLog{}).Preload("AdminUser")
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}
	if filter.AdminID != nil {
		query = query.Where("admin_user_id = ?", *filter.AdminID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&logs).Error
	return logs, total, err
}

type HeroSlideRepository struct {
	db *gorm.DB
}

func NewHeroSlideRepository(db *gorm.DB) *HeroSlideRepository {
	return &HeroSlideRepository{db: db}
}

func (r *HeroSlideRepository) ListPublic(ctx context.Context) ([]models.HeroSlide, error) {
	var slides []models.HeroSlide
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, created_at ASC").
		Find(&slides).Error
	return slides, err
}

func (r *HeroSlideRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.HeroSlide, error) {
	var slide models.HeroSlide
	err := r.db.WithContext(ctx).First(&slide, "id = ?", id).Error
	return &slide, err
}

func (r *HeroSlideRepository) Create(ctx context.Context, slide *models.HeroSlide) error {
	return r.db.WithContext(ctx).Create(slide).Error
}

func (r *HeroSlideRepository) Save(ctx context.Context, slide *models.HeroSlide) error {
	return r.db.WithContext(ctx).Save(slide).Error
}

func (r *HeroSlideRepository) Delete(ctx context.Context, slide *models.HeroSlide) error {
	return r.db.WithContext(ctx).Delete(slide).Error
}

func (r *HeroSlideRepository) ListAdmin(ctx context.Context) ([]models.HeroSlide, error) {
	var slides []models.HeroSlide
	err := r.db.WithContext(ctx).Order("sort_order ASC, created_at DESC").Find(&slides).Error
	return slides, err
}

type MediaAssetRepository struct {
	db *gorm.DB
}

func NewMediaAssetRepository(db *gorm.DB) *MediaAssetRepository {
	return &MediaAssetRepository{db: db}
}

func (r *MediaAssetRepository) Create(ctx context.Context, asset *models.MediaAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *MediaAssetRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.MediaAsset, error) {
	var asset models.MediaAsset
	err := r.db.WithContext(ctx).Preload("UploadedBy").First(&asset, "id = ?", id).Error
	return &asset, err
}

func (r *MediaAssetRepository) List(ctx context.Context, page int, limit int) ([]models.MediaAsset, int64, error) {
	var assets []models.MediaAsset
	var total int64
	query := r.db.WithContext(ctx).Model(&models.MediaAsset{}).Preload("UploadedBy")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&assets).Error
	return assets, total, err
}

func (r *MediaAssetRepository) Delete(ctx context.Context, asset *models.MediaAsset) error {
	return r.db.WithContext(ctx).Delete(asset).Error
}

// SEORepository
type SEORepository struct {
	db *gorm.DB
}

func NewSEORepository(db *gorm.DB) *SEORepository {
	return &SEORepository{db: db}
}

func (r *SEORepository) FindBySlug(ctx context.Context, slug string) (*models.SEOPage, error) {
	var page models.SEOPage
	err := r.db.WithContext(ctx).Where("page_slug = ?", slug).First(&page).Error
	return &page, err
}

func (r *SEORepository) Save(ctx context.Context, page *models.SEOPage) error {
	return r.db.WithContext(ctx).Save(page).Error
}

func (r *SEORepository) List(ctx context.Context) ([]models.SEOPage, error) {
	var pages []models.SEOPage
	err := r.db.WithContext(ctx).Order("page_slug ASC").Find(&pages).Error
	return pages, err
}

func (r *SEORepository) FindByID(ctx context.Context, id uuid.UUID) (*models.SEOPage, error) {
	var page models.SEOPage
	err := r.db.WithContext(ctx).First(&page, "id = ?", id).Error
	return &page, err
}

func (r *SEORepository) Delete(ctx context.Context, page *models.SEOPage) error {
	return r.db.WithContext(ctx).Delete(page).Error
}

// GalleryRepository
type GalleryRepository struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *GalleryRepository {
	return &GalleryRepository{db: db}
}

func (r *GalleryRepository) ListPublic(ctx context.Context) ([]models.GalleryItem, error) {
	var items []models.GalleryItem
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *GalleryRepository) ListAdmin(ctx context.Context) ([]models.GalleryItem, error) {
	var items []models.GalleryItem
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *GalleryRepository) Create(ctx context.Context, item *models.GalleryItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GalleryRepository) Save(ctx context.Context, item *models.GalleryItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GalleryRepository) Delete(ctx context.Context, item *models.GalleryItem) error {
	return r.db.WithContext(ctx).Delete(item).Error
}

func (r *GalleryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.GalleryItem, error) {
	var item models.GalleryItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	return &item, err
}

// RestaurantRepository
type RestaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) ListPublic(ctx context.Context) ([]models.RestaurantItem, error) {
	var items []models.RestaurantItem
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *RestaurantRepository) ListAdmin(ctx context.Context) ([]models.RestaurantItem, error) {
	var items []models.RestaurantItem
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *RestaurantRepository) Create(ctx context.Context, item *models.RestaurantItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RestaurantRepository) Save(ctx context.Context, item *models.RestaurantItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *RestaurantRepository) Delete(ctx context.Context, item *models.RestaurantItem) error {
	return r.db.WithContext(ctx).Delete(item).Error
}

func (r *RestaurantRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.RestaurantItem, error) {
	var item models.RestaurantItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	return &item, err
}

// SuiteRepository
type SuiteRepository struct {
	db *gorm.DB
}

func NewSuiteRepository(db *gorm.DB) *SuiteRepository {
	return &SuiteRepository{db: db}
}

func (r *SuiteRepository) ListPublic(ctx context.Context) ([]models.SuiteRoom, error) {
	var suites []models.SuiteRoom
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&suites).Error
	return suites, err
}

func (r *SuiteRepository) ListAdmin(ctx context.Context) ([]models.SuiteRoom, error) {
	var suites []models.SuiteRoom
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&suites).Error
	return suites, err
}

func (r *SuiteRepository) FindBySlug(ctx context.Context, slug string) (*models.SuiteRoom, error) {
	var suite models.SuiteRoom
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&suite).Error
	return &suite, err
}

func (r *SuiteRepository) Create(ctx context.Context, suite *models.SuiteRoom) error {
	return r.db.WithContext(ctx).Create(suite).Error
}

func (r *SuiteRepository) Save(ctx context.Context, suite *models.SuiteRoom) error {
	return r.db.WithContext(ctx).Save(suite).Error
}

func (r *SuiteRepository) Delete(ctx context.Context, suite *models.SuiteRoom) error {
	return r.db.WithContext(ctx).Delete(suite).Error
}

func (r *SuiteRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.SuiteRoom, error) {
	var suite models.SuiteRoom
	err := r.db.WithContext(ctx).First(&suite, "id = ?", id).Error
	return &suite, err
}

// HallRepository
type HallRepository struct {
	db *gorm.DB
}

func NewHallRepository(db *gorm.DB) *HallRepository {
	return &HallRepository{db: db}
}

func (r *HallRepository) ListPackagesPublic(ctx context.Context) ([]models.HallPackage, error) {
	var packages []models.HallPackage
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&packages).Error
	return packages, err
}

func (r *HallRepository) ListPackagesAdmin(ctx context.Context) ([]models.HallPackage, error) {
	var packages []models.HallPackage
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&packages).Error
	return packages, err
}

func (r *HallRepository) CreatePackage(ctx context.Context, pkg *models.HallPackage) error {
	return r.db.WithContext(ctx).Create(pkg).Error
}

func (r *HallRepository) SavePackage(ctx context.Context, pkg *models.HallPackage) error {
	return r.db.WithContext(ctx).Save(pkg).Error
}

func (r *HallRepository) DeletePackage(ctx context.Context, pkg *models.HallPackage) error {
	return r.db.WithContext(ctx).Delete(pkg).Error
}

func (r *HallRepository) FindPackageByID(ctx context.Context, id uuid.UUID) (*models.HallPackage, error) {
	var pkg models.HallPackage
	err := r.db.WithContext(ctx).First(&pkg, "id = ?", id).Error
	return &pkg, err
}

func (r *HallRepository) CreateEnquiry(ctx context.Context, enquiry *models.HallEnquiry) error {
	return r.db.WithContext(ctx).Create(enquiry).Error
}

func (r *HallRepository) ListEnquiries(ctx context.Context, page, limit int) ([]models.HallEnquiry, int64, error) {
	var enquiries []models.HallEnquiry
	var total int64
	query := r.db.WithContext(ctx).Model(&models.HallEnquiry{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&enquiries).Error
	return enquiries, total, err
}

func (r *HallRepository) FindEnquiryByID(ctx context.Context, id uuid.UUID) (*models.HallEnquiry, error) {
	var enquiry models.HallEnquiry
	err := r.db.WithContext(ctx).First(&enquiry, "id = ?", id).Error
	return &enquiry, err
}

func (r *HallRepository) SaveEnquiry(ctx context.Context, enquiry *models.HallEnquiry) error {
	return r.db.WithContext(ctx).Save(enquiry).Error
}

// OfferRepository
type OfferRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) ListActive(ctx context.Context) ([]models.Offer, error) {
	var offers []models.Offer
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND (starts_at IS NULL OR starts_at <= ?) AND (ends_at IS NULL OR ends_at >= ?)", true, now, now).
		Order("created_at DESC").
		Find(&offers).Error
	return offers, err
}

func (r *OfferRepository) ListAdmin(ctx context.Context) ([]models.Offer, error) {
	var offers []models.Offer
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&offers).Error
	return offers, err
}

func (r *OfferRepository) Create(ctx context.Context, offer *models.Offer) error {
	return r.db.WithContext(ctx).Create(offer).Error
}

func (r *OfferRepository) Save(ctx context.Context, offer *models.Offer) error {
	return r.db.WithContext(ctx).Save(offer).Error
}

func (r *OfferRepository) Delete(ctx context.Context, offer *models.Offer) error {
	return r.db.WithContext(ctx).Delete(offer).Error
}

func (r *OfferRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Offer, error) {
	var offer models.Offer
	err := r.db.WithContext(ctx).First(&offer, "id = ?", id).Error
	return &offer, err
}
