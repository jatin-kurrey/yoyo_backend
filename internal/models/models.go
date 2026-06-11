package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AdminRole string
type BookingStatus string
type MessageStatus string
type PaymentStatus string

const (
	RoleSuperAdmin AdminRole = "super_admin"
	RoleAdmin      AdminRole = "admin"
	RoleModerator  AdminRole = "moderator"
)

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentPaid     PaymentStatus = "paid"
	PaymentFailed   PaymentStatus = "failed"
	PaymentRefunded PaymentStatus = "refunded"

	BookingNew       BookingStatus = "new"
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingUsed      BookingStatus = "used"
	BookingCancelled BookingStatus = "cancelled"
	BookingRefunded  BookingStatus = "refunded"

	MessageNew      MessageStatus = "new"
	MessageRead     MessageStatus = "read"
	MessageReplied  MessageStatus = "replied"
	MessageArchived MessageStatus = "archived"
)

type AdminUser struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Email        string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         AdminRole      `gorm:"size:50;not null;default:'admin'" json:"role"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Ticket struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Slug          string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         int64          `gorm:"not null" json:"price"`          // In Paise
	OriginalPrice *int64         `json:"original_price"`                // In Paise, optional for discount badge
	Category      string         `gorm:"size:100;not null" json:"category"`
	Features      datatypes.JSON `json:"features"`                      // Array of strings
	Validity      string         `gorm:"size:100" json:"validity"`
	Stock         int            `gorm:"not null;default:0" json:"stock"`
	SoldCount     int            `gorm:"not null;default:0" json:"sold_count"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	IsBestseller  bool           `gorm:"not null;default:false" json:"is_bestseller"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Booking struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID          string         `gorm:"size:100;uniqueIndex;not null" json:"booking_id"`
	TicketID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"ticket_id"`
	Ticket             Ticket         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"ticket"`
	CustomerName       string         `gorm:"size:255;not null" json:"customer_name"`
	CustomerEmail      string         `gorm:"size:255;not null;index" json:"customer_email"`
	CustomerPhone      string         `gorm:"size:20;not null;index" json:"customer_phone"`
	VisitDate          time.Time      `gorm:"not null;index" json:"visit_date"`
	Quantity           int            `gorm:"not null" json:"quantity"`
	Amount             int64          `gorm:"not null" json:"amount"`
	Status             BookingStatus  `gorm:"size:50;not null;default:'new'" json:"status"`
	PaymentStatus      PaymentStatus  `gorm:"size:50;not null;default:'pending'" json:"payment_status"`
	RazorpayOrderID    string         `gorm:"size:255;index" json:"razorpay_order_id"`
	RazorpayPaymentID  string         `gorm:"size:255;index" json:"razorpay_payment_id"`
	RazorpaySignature  string         `gorm:"size:255" json:"-"`
	InternalNotes      string         `gorm:"type:text" json:"internal_notes"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type ContactMessage struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;not null;index" json:"email"`
	Phone     string         `gorm:"size:20;not null" json:"phone"`
	Subject   string         `gorm:"size:255;not null" json:"subject"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Status    MessageStatus  `gorm:"size:50;not null;default:'new'" json:"status"` // new, read, replied, archived
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SiteSetting struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SiteName         string         `gorm:"size:255;not null" json:"site_name"`
	LogoURL          string         `gorm:"type:text" json:"logo_url"`
	ContactEmail     string         `gorm:"size:255" json:"contact_email"`
	PhoneNumbers     datatypes.JSON `json:"phone_numbers"` // Array of strings
	WhatsAppNumber   string         `gorm:"size:20" json:"whatsapp_number"`
	Address          string         `gorm:"type:text" json:"address"`
	GoogleMapsURL    string         `gorm:"type:text" json:"google_maps_url"`
	OpeningHours     string         `gorm:"type:text" json:"opening_hours"`
	SocialLinks      datatypes.JSON `json:"social_links"`    // map[string]string
	MetaTitle        string         `gorm:"size:255" json:"meta_title"`
	MetaDescription  string         `gorm:"type:text" json:"meta_description"`
	RazorpayEnabled  bool           `json:"razorpay_enabled"`
	MaintenanceMode  bool           `json:"maintenance_mode"`
	FeatureToggles   datatypes.JSON `json:"feature_toggles"`   // map[string]bool (halls, restaurant, gallery, etc)
	HomepageSections datatypes.JSON `json:"homepage_sections"` // map[string]bool
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type AuditLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AdminUserID *uuid.UUID     `gorm:"type:uuid;index" json:"admin_user_id"`
	AdminUser   *AdminUser     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"admin_user"`
	Action      string         `gorm:"size:100;not null" json:"action"`
	Module      string         `gorm:"size:100;not null" json:"module"`
	Metadata    datatypes.JSON `json:"metadata"`
	IPAddress   string         `gorm:"size:50" json:"ip_address"`
	CreatedAt   time.Time      `json:"created_at"`
}

type HeroSlide struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title            string         `gorm:"size:255" json:"title"` // Headline
	Subtitle         string         `gorm:"size:255" json:"subtitle"`
	Description      string         `gorm:"type:text" json:"description"`
	ImageURL         string         `gorm:"type:text;not null" json:"image_url"`
	MobileImageURL   string         `gorm:"type:text" json:"mobile_image_url"`
	CTALabel         string         `gorm:"size:100" json:"cta_label"` // CTAText
	CTAURL           string         `gorm:"size:255" json:"cta_url"`
	SecondaryCTALabel string         `gorm:"size:100" json:"secondary_cta_label"`
	SecondaryCTAURL   string         `gorm:"size:255" json:"secondary_cta_url"`
	BadgeText        string         `gorm:"size:100" json:"badge_text"`
	SortOrder        int            `gorm:"not null;default:0" json:"sort_order"`
	IsActive         bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type ContentPage struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Slug            string         `gorm:"size:180;uniqueIndex;not null" json:"slug"`
	Title           string         `gorm:"size:255;not null" json:"title"`
	Content         string         `gorm:"type:text" json:"content"`
	MetaTitle       string         `gorm:"size:255" json:"meta_title"`
	MetaDescription string         `gorm:"type:text" json:"meta_description"`
	IsPublished     bool           `gorm:"not null;default:true" json:"is_published"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type MediaAsset struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	URL              string         `gorm:"type:text;not null" json:"url"`
	StorageKey       string         `gorm:"size:255;not null;index" json:"storage_key"`
	Filename         string         `gorm:"size:255;not null" json:"filename"`
	OriginalFilename string         `gorm:"size:255" json:"original_filename"`
	MimeType         string         `gorm:"size:100;not null" json:"mime_type"`
	SizeBytes        int64          `gorm:"not null" json:"size_bytes"`
	StorageProvider  string         `gorm:"size:50;not null;index" json:"storage_provider"` // local, r2
	UploadedByID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"uploaded_by_id"`
	AltText          string         `gorm:"size:255" json:"alt_text"`
	Folder           string         `gorm:"size:100;index" json:"folder"` // hero, gallery, etc
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type SEOPage struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PageSlug        string         `gorm:"size:180;uniqueIndex;not null" json:"page_slug"`
	MetaTitle       string         `gorm:"size:255" json:"meta_title"`
	MetaDescription string         `gorm:"type:text" json:"meta_description"`
	CanonicalURL    string         `gorm:"type:text" json:"canonical_url"`
	OGTitle         string         `gorm:"size:255" json:"og_title"`
	OGDescription   string         `gorm:"type:text" json:"og_description"`
	OGImage         string         `gorm:"type:text" json:"og_image"`
	RobotsIndex     bool           `gorm:"not null;default:true" json:"robots_index"`
	RobotsFollow    bool           `gorm:"not null;default:true" json:"robots_follow"`
	SchemaJSON      datatypes.JSON `json:"schema_json"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type GalleryItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:text;not null" json:"image_url"`
	Category    string         `gorm:"size:100;index" json:"category"`
	AltText     string         `gorm:"size:255" json:"alt_text"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RestaurantItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:text" json:"image_url"`
	Category    string         `gorm:"size:100;index" json:"category"`
	Price       int64          `json:"price"` // In Paise, optional
	IsFeatured  bool           `gorm:"not null;default:false" json:"is_featured"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type SuiteRoom struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Slug          string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Description   string         `gorm:"type:text" json:"description"`
	ImageURL      string         `gorm:"type:text;not null" json:"image_url"`
	Gallery       datatypes.JSON `json:"gallery"` // Array of URLs
	PricePerNight int64          `gorm:"not null" json:"price_per_night"`
	MaxGuests     int            `gorm:"not null;default:2" json:"max_guests"`
	Amenities     datatypes.JSON `json:"amenities"` // Array of strings
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type HallPackage struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	ImageURL      string         `gorm:"type:text;not null" json:"image_url"`
	Capacity      int            `json:"capacity"`
	StartingPrice int64          `json:"starting_price"`
	SuitableFor   datatypes.JSON `json:"suitable_for"` // Array of strings
	Features      datatypes.JSON `json:"features"`     // Array of strings
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type HallEnquiry struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Phone          string    `gorm:"size:20;not null;index" json:"phone"`
	Email          string    `gorm:"size:255;index" json:"email"`
	EventType      string    `gorm:"size:100" json:"event_type"`
	ExpectedGuests int       `json:"expected_guests"`
	PreferredDate  time.Time `json:"preferred_date"`
	Message        string    `gorm:"type:text" json:"message"`
	Status         string    `gorm:"size:50;not null;default:'new'" json:"status"` // new, contacted, converted, lost
	Source         string    `gorm:"size:100" json:"source"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Offer struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Code          string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	DiscountType  string         `gorm:"size:50;not null" json:"discount_type"` // percentage, fixed
	DiscountValue int64          `gorm:"not null" json:"discount_value"`
	StartsAt      *time.Time     `json:"starts_at"`
	EndsAt        *time.Time     `json:"ends_at"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// GORM Hooks

func setUUID(id *uuid.UUID) {
	if id != nil && *id == uuid.Nil {
		*id = uuid.New()
	}
}

func (m *AdminUser) BeforeCreate(tx *gorm.DB) error { setUUID(&m.ID); return nil }
func (m *Ticket) BeforeCreate(tx *gorm.DB) error    { setUUID(&m.ID); return nil }
func (m *Booking) BeforeCreate(tx *gorm.DB) error   { setUUID(&m.ID); return nil }
func (m *ContactMessage) BeforeCreate(tx *gorm.DB) error { setUUID(&m.ID); return nil }
func (m *SiteSetting) BeforeCreate(tx *gorm.DB) error    { setUUID(&m.ID); return nil }
func (m *AuditLog) BeforeCreate(tx *gorm.DB) error       { setUUID(&m.ID); return nil }
func (m *HeroSlide) BeforeCreate(tx *gorm.DB) error      { setUUID(&m.ID); return nil }
func (m *ContentPage) BeforeCreate(tx *gorm.DB) error    { setUUID(&m.ID); return nil }
func (m *MediaAsset) BeforeCreate(tx *gorm.DB) error     { setUUID(&m.ID); return nil }
func (m *SEOPage) BeforeCreate(tx *gorm.DB) error        { setUUID(&m.ID); return nil }
func (m *GalleryItem) BeforeCreate(tx *gorm.DB) error    { setUUID(&m.ID); return nil }
func (m *RestaurantItem) BeforeCreate(tx *gorm.DB) error { setUUID(&m.ID); return nil }
func (m *SuiteRoom) BeforeCreate(tx *gorm.DB) error      { setUUID(&m.ID); return nil }
func (m *HallPackage) BeforeCreate(tx *gorm.DB) error    { setUUID(&m.ID); return nil }
func (m *HallEnquiry) BeforeCreate(tx *gorm.DB) error    { setUUID(&m.ID); return nil }
func (m *Offer) BeforeCreate(tx *gorm.DB) error          { setUUID(&m.ID); return nil }
