package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RoomStatus string
type CleanStatus string
type BookingPlan string
type BookingSource string
type BookingStatusPMS string
type FolioEntryType string
type PaymentMode string
type PaymentType string
type POSTableStatus string
type POSOrderStatus string
type HKTaskType string
type HKTaskStatus string

const (
	RoomAvailable RoomStatus = "available"
	RoomOccupied  RoomStatus = "occupied"
	RoomOOO       RoomStatus = "ooo"
	RoomBlocked   RoomStatus = "blocked"

	CleanClean CleanStatus = "clean"
	CleanDirty CleanStatus = "dirty"

	PlanEP BookingPlan = "EP"
	PlanCP BookingPlan = "CP"
	PlanAP BookingPlan = "AP"

	SourceWalkIn     BookingSource = "Walk-In"
	SourceAgoda      BookingSource = "Agoda"
	SourceMMT        BookingSource = "MakeMyTrip"
	SourceBookingCom BookingSource = "Booking.com"
	SourceCorporate  BookingSource = "Corporate"

	PMSBookingHold        BookingStatusPMS = "hold"
	PMSBookingConfirmed   BookingStatusPMS = "confirmed"
	PMSBookingCheckedIn   BookingStatusPMS = "checked-in"
	PMSBookingCheckedOut  BookingStatusPMS = "checked-out"
	PMSBookingCancelled   BookingStatusPMS = "cancelled"

	FolioRoom     FolioEntryType = "room"
	FolioFood     FolioEntryType = "food"
	FolioLaundry  FolioEntryType = "laundry"
	FolioSpa      FolioEntryType = "spa"
	FolioMinibar  FolioEntryType = "minibar"
	FolioOther    FolioEntryType = "other"

	PayCash         PaymentMode = "Cash"
	PayUPI          PaymentMode = "UPI"
	PayCard         PaymentMode = "Card"
	PayBankTransfer PaymentMode = "Bank Transfer"

	PayAdvance    PaymentType = "advance"
	PaySettlement PaymentType = "settlement"
	PayRefund     PaymentType = "refund"

	POSTableVacant   POSTableStatus = "vacant"
	POSTableOccupied POSTableStatus = "occupied"
	POSTableBilled   POSTableStatus = "billed"

	POSOrderOpen     POSOrderStatus = "open"
	POSOrderPreparing POSOrderStatus = "preparing"
	POSOrderServed   POSOrderStatus = "served"
	POSOrderBilled   POSOrderStatus = "billed"

	HKClean        HKTaskType = "clean"
	HKMaintenance  HKTaskType = "maintenance"
	HKOOO          HKTaskType = "ooo"

	HKPending    HKTaskStatus = "pending"
	HKInProgress HKTaskStatus = "in-progress"
	HKCompleted  HKTaskStatus = "completed"
)

type PMSRoomCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Slug        string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	BasePrice   int64          `json:"base_price"`
	MaxGuests   int            `gorm:"not null;default:2" json:"max_guests"`
	Amenities   datatypes.JSON `json:"amenities"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type PMSRoom struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RoomNumber int            `gorm:"uniqueIndex;not null" json:"room_number"`
	Floor      int            `gorm:"not null;default:1" json:"floor"`
	CategoryID uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Category   PMSRoomCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"category"`
	Status     RoomStatus     `gorm:"size:50;not null;default:'available'" json:"status"`
	CleanStatus CleanStatus   `gorm:"size:50;not null;default:'clean'" json:"clean_status"`
	OOOReason  string         `gorm:"type:text" json:"ooo_reason"`
	Notes      string         `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type PMSBooking struct {
	ID          uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	BookingRef  string           `gorm:"size:50;uniqueIndex;not null" json:"booking_ref"`
	RoomID      uuid.UUID        `gorm:"type:uuid;not null;index" json:"room_id"`
	Room        PMSRoom          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room"`
	GuestName   string           `gorm:"size:255;not null" json:"guest_name"`
	GuestPhone  string           `gorm:"size:20;not null;index" json:"guest_phone"`
	GuestEmail  string           `gorm:"size:255" json:"guest_email"`
	Adults      int              `gorm:"not null;default:2" json:"adults"`
	Children    int              `gorm:"not null;default:0" json:"children"`
	Plan        BookingPlan      `gorm:"size:10;not null;default:'EP'" json:"plan"`
	Source      BookingSource    `gorm:"size:50;not null;default:'Walk-In'" json:"source"`
	CheckIn     time.Time        `gorm:"not null;index" json:"check_in"`
	CheckOut    time.Time        `gorm:"not null;index" json:"check_out"`
	RatePerNight int64           `json:"rate_per_night"`
	TotalAmount int64            `json:"total_amount"`
	Discount    int64            `json:"discount"`
	Tax         int64            `json:"tax"`
	PaidAmount  int64            `json:"paid_amount"`
	BalanceAmount int64          `json:"balance_amount"`
	Status      BookingStatusPMS `gorm:"size:50;not null;default:'hold'" json:"status"`
	CreatedByID *uuid.UUID       `gorm:"type:uuid;index" json:"created_by_id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

type FolioEntry struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"booking_id"`
	Booking     PMSBooking     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Type        FolioEntryType `gorm:"size:50;not null" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	Amount      int64          `gorm:"not null" json:"amount"`
	Quantity    int            `gorm:"not null;default:1" json:"quantity"`
	PostedAt    time.Time       `gorm:"not null" json:"posted_at"`
	PostedByID  *uuid.UUID     `gorm:"type:uuid;index" json:"posted_by_id"`
}

type PMSPayment struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID   uuid.UUID    `gorm:"type:uuid;not null;index" json:"booking_id"`
	Booking     PMSBooking   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Mode        PaymentMode  `gorm:"size:50;not null" json:"mode"`
	Amount      int64        `gorm:"not null" json:"amount"`
	Type        PaymentType  `gorm:"size:50;not null;default:'advance'" json:"type"`
	Reference   string       `gorm:"size:255" json:"reference"`
	ReceivedAt  time.Time    `gorm:"not null" json:"received_at"`
	ReceivedByID *uuid.UUID  `gorm:"type:uuid;index" json:"received_by_id"`
}

type POSTable struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TableNumber     int            `gorm:"uniqueIndex;not null" json:"table_number"`
	Capacity        int            `gorm:"not null;default:4" json:"capacity"`
	Area            string         `gorm:"size:100;not null" json:"area"`
	Status          POSTableStatus `gorm:"size:50;not null;default:'vacant'" json:"status"`
	GuestName       string         `gorm:"size:255" json:"guest_name"`
	CurrentOrderValue int64        `json:"current_order_value"`
	KOTCount        int            `gorm:"not null;default:0" json:"kot_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type POSOrder struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TableID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"table_id"`
	Table       POSTable       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	MenuItemID  *uuid.UUID     `gorm:"type:uuid" json:"menu_item_id"`
	ItemName    string         `gorm:"size:255;not null" json:"item_name"`
	Quantity    int            `gorm:"not null;default:1" json:"quantity"`
	Price       int64          `gorm:"not null" json:"price"`
	Notes       string         `gorm:"type:text" json:"notes"`
	KOTNumber   int            `gorm:"not null;default:0" json:"kot_number"`
	Status      POSOrderStatus `gorm:"size:50;not null;default:'open'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	ServedAt    *time.Time     `json:"served_at"`
}

type HKTask struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"room_id"`
	Room        PMSRoom        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	StaffName   string         `gorm:"size:255" json:"staff_name"`
	TaskType    HKTaskType     `gorm:"size:50;not null" json:"task_type"`
	Status      HKTaskStatus   `gorm:"size:50;not null;default:'pending'" json:"status"`
	AssignedAt  *time.Time     `json:"assigned_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *PMSRoomCategory) BeforeCreate(tx *gorm.DB) error { setUUID(&m.ID); return nil }
func (m *PMSRoom) BeforeCreate(tx *gorm.DB) error         { setUUID(&m.ID); return nil }
func (m *PMSBooking) BeforeCreate(tx *gorm.DB) error      { setUUID(&m.ID); return nil }
func (m *FolioEntry) BeforeCreate(tx *gorm.DB) error      { setUUID(&m.ID); return nil }
func (m *PMSPayment) BeforeCreate(tx *gorm.DB) error      { setUUID(&m.ID); return nil }
func (m *POSTable) BeforeCreate(tx *gorm.DB) error        { setUUID(&m.ID); return nil }
func (m *POSOrder) BeforeCreate(tx *gorm.DB) error        { setUUID(&m.ID); return nil }
func (m *HKTask) BeforeCreate(tx *gorm.DB) error          { setUUID(&m.ID); return nil }
