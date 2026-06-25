package seeds

import (
	"context"
	"encoding/json"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Run(ctx context.Context, cfg *config.Config, db *gorm.DB, svc *services.Services) error {
	// Clean up old location references from the database
	db.Exec("UPDATE hero_slides SET title = REPLACE(title, 'Indore''s', 'Bhilai Durg''s'), subtitle = REPLACE(subtitle, 'Indore''s', 'Bhilai Durg''s'), description = REPLACE(description, 'Indore''s', 'Bhilai Durg''s')")
	db.Exec("UPDATE hero_slides SET title = REPLACE(title, 'Indore', 'Bhilai Durg'), subtitle = REPLACE(subtitle, 'Indore', 'Bhilai Durg'), description = REPLACE(description, 'Indore', 'Bhilai Durg')")
	db.Exec("UPDATE attractions SET title = REPLACE(title, 'Indore''s', 'Bhilai Durg''s'), description = REPLACE(description, 'Indore''s', 'Bhilai Durg''s')")
	db.Exec("UPDATE attractions SET title = REPLACE(title, 'Indore', 'Bhilai Durg'), description = REPLACE(description, 'Indore', 'Bhilai Durg')")
	db.Exec("UPDATE gallery_items SET title = REPLACE(title, 'Indore''s', 'Bhilai Durg''s'), description = REPLACE(description, 'Indore''s', 'Bhilai Durg''s')")
	db.Exec("UPDATE gallery_items SET title = REPLACE(title, 'Indore', 'Bhilai Durg'), description = REPLACE(description, 'Indore', 'Bhilai Durg')")
	db.Exec("UPDATE tickets SET title = REPLACE(title, 'Indore''s', 'Bhilai Durg''s'), description = REPLACE(description, 'Indore''s', 'Bhilai Durg''s')")
	db.Exec("UPDATE tickets SET title = REPLACE(title, 'Indore', 'Bhilai Durg'), description = REPLACE(description, 'Indore', 'Bhilai Durg')")

	if err := svc.Auth.EnsureSuperAdmin(ctx); err != nil {
		return err
	}
	if _, err := svc.Settings.Get(ctx); err != nil {
		return err
	}
	if err := seedHeroSlides(ctx, cfg, db); err != nil {
		return err
	}
	if err := seedTickets(ctx, cfg, db); err != nil {
		return err
	}
	if err := seedAttractions(ctx, cfg, db); err != nil {
		return err
	}
	if err := seedRestaurantItems(ctx, cfg, db); err != nil {
		return err
	}
	if err := seedSuiteRooms(ctx, cfg, db); err != nil {
		return err
	}
	if err := seedGallery(ctx, cfg, db); err != nil {
		return err
	}
	if err := seedHallPackages(ctx, cfg, db); err != nil {
		return err
	}
	return seedPMSTables(ctx, db)
}

func seedTickets(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title         string
		Description   string
		Price         int64
		OriginalPrice *int64
		Category      string
		Stock         int
		SortOrder     int
		Features      []string
	}{
		{
			Title:       "Standard Pass",
			Description: "Single entry to the park with all-day access to core attractions.",
			Price:       499,
			Category:    "general",
			Stock:       100,
			SortOrder:   1,
			Features:    []string{"All Day Entry", "Locker Access", "Safety Gear Included"},
		},
		{
			Title:       "VIP Pass",
			Description: "Skip the line experience with a complimentary drink for a smoother day out.",
			Price:       999,
			Category:    "vip",
			Stock:       50,
			SortOrder:   2,
			Features:    []string{"Priority Entry", "Free Drink", "Locker Access", "Safety Gear Included"},
		},
		{
			Title:       "Family Bundle",
			Description: "Entry for 4 people, designed for families and small groups.",
			Price:       2999,
			Category:    "family",
			Stock:       30,
			SortOrder:   3,
			Features:    []string{"Entry for 4", "Family Check-in", "Locker Access", "Safety Gear Included"},
		},
	}

	for _, item := range defaults {
		slug := utils.Slugify(item.Title)
		var existing models.Ticket
		err := db.WithContext(ctx).Where("slug = ?", slug).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		features, err := json.Marshal(item.Features)
		if err != nil {
			return err
		}
		ticket := models.Ticket{
			Title:         item.Title,
			Slug:          slug,
			Description:   item.Description,
			Price:         item.Price,
			OriginalPrice: item.OriginalPrice,
			Category:      item.Category,
			Features:      datatypes.JSON(features),
			Validity:      "Valid for selected visit date",
			Stock:         item.Stock,
			IsActive:      true,
			SortOrder:     item.SortOrder,
		}
		if err := db.WithContext(ctx).Create(&ticket).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAttractions(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title       string
		Description string
		ImageURL    string
		IconName    string
		Tag         string
		SortOrder   int
	}{
		{
			Title:       "Giant Slides",
			Description: "Experience 50+ feet of pure adrenaline with our high-speed vertical drops.",
			ImageURL:    "https://images.unsplash.com/photo-1542332213-9b5a5a3fad35?q=80&w=800&auto=format&fit=crop",
			IconName:    "Zap",
			Tag:         "Thrills",
			SortOrder:   1,
		},
		{
			Title:       "Massive Wave Pool",
			Description: "Bhilai Durg's largest wave pool with state-of-the-art ocean tide simulation.",
			ImageURL:    "https://images.unsplash.com/photo-1519817650390-64a93db51149?q=80&w=800&auto=format&fit=crop",
			IconName:    "Waves",
			Tag:         "Family",
			SortOrder:   2,
		},
		{
			Title:       "Kids Fantasy Zone",
			Description: "A safe, magical water playground designed exclusively for our little guests.",
			ImageURL:    "https://images.unsplash.com/photo-1582650625119-3a31f8fa2699?q=80&w=800&auto=format&fit=crop",
			IconName:    "Heart",
			Tag:         "Kids",
			SortOrder:   3,
		},
		{
			Title:       "Cyclone Funnel",
			Description: "Spiral through the massive funnel at high speeds for a dizzying splashdown.",
			ImageURL:    "https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?q=80&w=800&auto=format&fit=crop",
			IconName:    "Star",
			Tag:         "Thrills",
			SortOrder:   4,
		},
		{
			Title:       "Lazy River",
			Description: "Relax and float along our 400ft tropical river with gentle currents.",
			ImageURL:    "https://images.unsplash.com/photo-1576013551627-0cc20b96c2a7?q=80&w=800&auto=format&fit=crop",
			IconName:    "Anchor",
			Tag:         "Relax",
			SortOrder:   5,
		},
		{
			Title:       "Rain Dance Arena",
			Description: "Dance to the latest hits under high-tech water sprinklers and disco lights.",
			ImageURL:    "https://images.unsplash.com/photo-1533174072545-7a4b6ad7a6c3?q=80&w=800&auto=format&fit=crop",
			IconName:    "Sun",
			Tag:         "Fun",
			SortOrder:   6,
		},
	}

	for _, item := range defaults {
		var existing models.Attraction
		err := db.WithContext(ctx).Where("title = ?", item.Title).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		attraction := models.Attraction{
			Title:       item.Title,
			Description: item.Description,
			ImageURL:    item.ImageURL,
			IconName:    item.IconName,
			Tag:         item.Tag,
			IsActive:    true,
			SortOrder:   item.SortOrder,
		}
		if err := db.WithContext(ctx).Create(&attraction).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedRestaurantItems(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title       string
		Description string
		ImageURL    string
		IconName    string
		Category    string
		SortOrder   int
	}{
		{
			Title:       "Family Restaurant",
			Description: "Spacious dining area with multi-cuisine options perfect for family lunches.",
			ImageURL:    "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?q=80&w=800&auto=format&fit=crop",
			IconName:    "Users",
			Category:    "Restaurant",
			SortOrder:   1,
		},
		{
			Title:       "Snacks & Beverages",
			Description: "Quick bites, refreshing mocktails, and seasonal drinks to keep you energized.",
			ImageURL:    "https://images.unsplash.com/photo-1559339352-11d035aa65de?q=80&w=800&auto=format&fit=crop",
			IconName:    "Coffee",
			Category:    "Snacks",
			SortOrder:   2,
		},
		{
			Title:       "Group Meal Packages",
			Description: "Cost-effective and delicious buffet options for school trips and corporate groups.",
			ImageURL:    "https://images.unsplash.com/photo-1547573854-74d2a71d0826?q=80&w=800&auto=format&fit=crop",
			IconName:    "Utensils",
			Category:    "Packages",
			SortOrder:   3,
		},
		{
			Title:       "Hygienic Kitchen",
			Description: "Strict quality controls and fresh ingredients for a safe dining experience.",
			ImageURL:    "https://images.unsplash.com/photo-1556910103-1c02745aae4d?q=80&w=800&auto=format&fit=crop",
			IconName:    "ShieldCheck",
			Category:    "Kitchen",
			SortOrder:   4,
		},
	}

	for _, item := range defaults {
		var existing models.RestaurantItem
		err := db.WithContext(ctx).Where("title = ?", item.Title).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		restaurantItem := models.RestaurantItem{
			Title:       item.Title,
			Description: item.Description,
			ImageURL:    item.ImageURL,
			IconName:    item.IconName,
			Category:    item.Category,
			IsActive:    true,
			SortOrder:   item.SortOrder,
		}
		if err := db.WithContext(ctx).Create(&restaurantItem).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedSuiteRooms(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title       string
		Slug        string
		IconName    string
		Description string
		ImageURL    string
		SortOrder   int
	}{
		{
			Title:       "AC Premium Rooms",
			Slug:        "ac-premium-rooms",
			IconName:    "Zap",
			Description: "Luxury air-conditioned rooms for a refreshing rest after your waterpark adventure.",
			ImageURL:    "https://images.unsplash.com/photo-1618773928121-c32242e63f39?q=80&w=800&auto=format&fit=crop",
			SortOrder:   1,
		},
		{
			Title:       "Family Suites",
			Slug:        "family-suites",
			IconName:    "Home",
			Description: "Spacious suites designed for larger families with all modern amenities.",
			ImageURL:    "https://images.unsplash.com/photo-1566665797739-1674de7a421a?q=80&w=800&auto=format&fit=crop",
			SortOrder:   2,
		},
		{
			Title:       "Clean Washrooms",
			Slug:        "clean-washrooms",
			IconName:    "ShieldCheck",
			Description: "Maintaining the highest standards of sanitation and hygiene for your comfort.",
			ImageURL:    "https://images.unsplash.com/photo-1584622650111-993a426fbf0a?q=80&w=800&auto=format&fit=crop",
			SortOrder:   3,
		},
		{
			Title:       "Easy Booking",
			Slug:        "easy-booking",
			IconName:    "Hotel",
			Description: "Hassle-free reservation process with instant confirmation and support.",
			ImageURL:    "https://images.unsplash.com/photo-1563911302283-d2bc129e7570?q=80&w=800&auto=format&fit=crop",
			SortOrder:   4,
		},
	}

	for _, item := range defaults {
		var existing models.SuiteRoom
		err := db.WithContext(ctx).Where("slug = ?", item.Slug).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		suiteRoom := models.SuiteRoom{
			Title:         item.Title,
			Slug:          item.Slug,
			IconName:      item.IconName,
			Description:   item.Description,
			ImageURL:      item.ImageURL,
			Gallery:       []byte(`[]`),
			Amenities:     []byte(`[]`),
			PricePerNight: 0,
			MaxGuests:     2,
			IsActive:      true,
			SortOrder:     item.SortOrder,
		}
		if err := db.WithContext(ctx).Create(&suiteRoom).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedGallery(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title       string
		Category    string
		ImageURL    string
		Description string
		SortOrder   int
	}{
		{
			Title:       "Main Wave Pool",
			Category:    "Water Park",
			ImageURL:    "https://images.unsplash.com/photo-1582650625119-3a31f8fa2699?q=80&w=1200&auto=format&fit=crop",
			Description: "Bhilai Durg's largest ocean tide simulation pool.",
			SortOrder:   1,
		},
		{
			Title:       "Family Dining Area",
			Category:    "Food",
			ImageURL:    "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?q=80&w=600&auto=format&fit=crop",
			Description: "Spacious multi-cuisine family dining experience.",
			SortOrder:   2,
		},
		{
			Title:       "Giant Slide Tower",
			Category:    "Water Park",
			ImageURL:    "https://images.unsplash.com/photo-1542332213-9b5a5a3fad35?q=80&w=800&auto=format&fit=crop",
			Description: "Pure speed and high-adrenaline drops.",
			SortOrder:   3,
		},
		{
			Title:       "Luxury Resort Suite",
			Category:    "Stay",
			ImageURL:    "https://images.unsplash.com/photo-1618773928121-c32242e63f39?q=80&w=1200&auto=format&fit=crop",
			Description: "Elegant suites for relaxing rest.",
			SortOrder:   4,
		},
		{
			Title:       "Night View of Park",
			Category:    "Events",
			ImageURL:    "https://images.unsplash.com/photo-1533174072545-7a4b6ad7a6c3?q=80&w=600&auto=format&fit=crop",
			Description: "Vibrant evening celebration vibes under disco lights.",
			SortOrder:   5,
		},
		{
			Title:       "Kids Activity Zone",
			Category:    "Play",
			ImageURL:    "https://images.unsplash.com/photo-1596464716127-f2a82984de30?q=80&w=600&auto=format&fit=crop",
			Description: "Safe and magical pool playground for our little guests.",
			SortOrder:   6,
		},
		{
			Title:       "Poolside Cafe",
			Category:    "Food",
			ImageURL:    "https://images.unsplash.com/photo-1559339352-11d035aa65de?q=80&w=600&auto=format&fit=crop",
			Description: "Delicious snacks and mocktails.",
			SortOrder:   7,
		},
		{
			Title:       "Lazy River Walk",
			Category:    "Water Park",
			ImageURL:    "https://images.unsplash.com/photo-1629113645366-3d71206637ba?q=80&w=1200&auto=format&fit=crop",
			Description: "Gentle currents along tropical river path.",
			SortOrder:   8,
		},
		{
			Title:       "Birthday Celebration",
			Category:    "Events",
			ImageURL:    "https://images.unsplash.com/photo-1464366400600-7168b8af9bc3?q=80&w=600&auto=format&fit=crop",
			Description: "Special events and celebrations designed for you.",
			SortOrder:   9,
		},
	}

	for _, item := range defaults {
		var existing models.GalleryItem
		err := db.WithContext(ctx).Where("title = ?", item.Title).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		galleryItem := models.GalleryItem{
			Title:       item.Title,
			Description: item.Description,
			ImageURL:    item.ImageURL,
			Category:    item.Category,
			SortOrder:   item.SortOrder,
			IsActive:    true,
		}
		if err := db.WithContext(ctx).Create(&galleryItem).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedHallPackages(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title         string
		Description   string
		ImageURL      string
		Capacity      int
		StartingPrice int64
		SuitableFor   []string
		Features      []string
		SortOrder     int
	}{
		{
			Title:         "Imperial Grand Ballroom",
			Description:   "Our signature luxury hall with high-ceilings, premium crystal chandeliers, and state-of-the-art acoustics. Ideal for royal weddings and mega conferences.",
			ImageURL:      "https://images.unsplash.com/photo-1519167758481-83f550bb49b3?q=80&w=800&auto=format&fit=crop",
			Capacity:      500,
			StartingPrice: 7500000,
			SuitableFor:   []string{"Weddings", "Receptions", "Conferences"},
			Features:      []string{"Central AC", "LED Screens", "Stage Decor Included", "Premium Sound System", "VIP Valet Parking"},
			SortOrder:     1,
		},
		{
			Title:         "Executive Boardroom",
			Description:   "An elegant, soundproof space equipped with high-speed internet, smart projectors, and ergonomic seating. Best suited for business meetings and corporate seminars.",
			ImageURL:      "https://images.unsplash.com/photo-1431540015161-0bf868a2d407?q=80&w=800&auto=format&fit=crop",
			Capacity:      50,
			StartingPrice: 2000000,
			SuitableFor:   []string{"Corporate Meetings", "Seminars", "Press Releases"},
			Features:      []string{"Ultra HD Projector", "Smart Board", "Soundproof Walls", "Video Conferencing Support", "Hi-Tea Catering Option"},
			SortOrder:     2,
		},
		{
			Title:         "Palms Open Air Lawn",
			Description:   "A gorgeous green lawn decorated with tropical lighting and private gazebos. Perfect for cocktail parties, social gatherings, and lively birthday events under the stars.",
			ImageURL:      "https://images.unsplash.com/photo-1533105079780-92b9be482077?q=80&w=800&auto=format&fit=crop",
			Capacity:      350,
			StartingPrice: 4500000,
			SuitableFor:   []string{"Birthdays", "Theme Parties", "Social Gatherings"},
			Features:      []string{"Outdoor Ambient Lighting", "Live BBQ Counter Setup", "Dance Floor Area", "Power Backup", "Valet Parking"},
			SortOrder:     3,
		},
	}

	for _, item := range defaults {
		var existing models.HallPackage
		err := db.WithContext(ctx).Where("title = ?", item.Title).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		suitableJson, err := json.Marshal(item.SuitableFor)
		if err != nil {
			return err
		}
		featuresJson, err := json.Marshal(item.Features)
		if err != nil {
			return err
		}

		pkg := models.HallPackage{
			Title:         item.Title,
			Description:   item.Description,
			ImageURL:      item.ImageURL,
			Capacity:      item.Capacity,
			StartingPrice: item.StartingPrice,
			SuitableFor:   datatypes.JSON(suitableJson),
			Features:      datatypes.JSON(featuresJson),
			IsActive:      true,
			SortOrder:     item.SortOrder,
		}
		if err := db.WithContext(ctx).Create(&pkg).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedHeroSlides(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title             string
		Subtitle          string
		Description       string
		ImageURL          string
		CTALabel          string
		CTAURL            string
		SecondaryCTALabel string
		SecondaryCTAURL   string
		BadgeText         string
		SortOrder         int
	}{
		{
			Title:             "Best Water Park in Bhilai Durg for Family Fun & Kids Safety",
			Subtitle:          "Slides, wave pools, kids zone & full-day fun — starting at ₹499",
			Description:       "Experience the ultimate adrenaline rush with 15+ world-class attractions.",
			ImageURL:          "https://images.unsplash.com/photo-1708157730402-67cc5b19e335?q=80&w=1600&auto=format&fit=crop",
			CTALabel:          "Book Tickets Now",
			CTAURL:            "/tickets",
			SecondaryCTALabel: "Explore Attractions",
			SecondaryCTAURL:   "/gallery",
			BadgeText:         "Open Today: 10 AM - 6 PM",
			SortOrder:         1,
		},
		{
			Title:             "Thrilling Slides & Massive Wave Pools Await You",
			Subtitle:          "Bhilai Durg's largest wave pool with state-of-the-art ocean tide simulation.",
			Description:       "Experience the ultimate adrenaline rush with 15+ world-class attractions.",
			ImageURL:          "https://images.unsplash.com/photo-1739295194212-0602c4d1e797?q=80&w=1600&auto=format&fit=crop",
			CTALabel:          "View Attractions",
			CTAURL:            "/gallery",
			SecondaryCTALabel: "Explore Pricing",
			SecondaryCTAURL:   "/tickets",
			BadgeText:         "Top Rated Family Destination",
			SortOrder:         2,
		},
	}

	for _, item := range defaults {
		var existing models.HeroSlide
		err := db.WithContext(ctx).Where("title = ?", item.Title).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		slide := models.HeroSlide{
			Title:             item.Title,
			Subtitle:          item.Subtitle,
			Description:       item.Description,
			ImageURL:          item.ImageURL,
			CTALabel:          item.CTALabel,
			CTAURL:            item.CTAURL,
			SecondaryCTALabel: item.SecondaryCTALabel,
			SecondaryCTAURL:   item.SecondaryCTAURL,
			BadgeText:         item.BadgeText,
			SortOrder:         item.SortOrder,
			IsActive:          true,
		}
		if err := db.WithContext(ctx).Create(&slide).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedPMSTables(ctx context.Context, db *gorm.DB) error {
	var count int64
	if err := db.WithContext(ctx).Model(&models.PMSRoomCategory{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		categories := []models.PMSRoomCategory{
			{
				ID:          uuid.New(),
				Name:        "Super Deluxe",
				Slug:        "super-deluxe",
				Description: "Super Deluxe Rooms",
				BasePrice:   4000,
				MaxGuests:   2,
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Name:        "Family Suite",
				Slug:        "family-suite",
				Description: "Family Suite Rooms",
				BasePrice:   6000,
				MaxGuests:   4,
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Name:        "Executive Pack",
				Slug:        "executive-pack",
				Description: "Executive Pack Rooms",
				BasePrice:   8000,
				MaxGuests:   5,
				IsActive:    true,
			},
		}
		for i := range categories {
			if err := db.WithContext(ctx).Create(&categories[i]).Error; err != nil {
				return err
			}
		}
	}

	if err := db.WithContext(ctx).Model(&models.PMSRoom{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		var superDeluxeCat, familySuiteCat, execPackCat models.PMSRoomCategory
		if err := db.WithContext(ctx).First(&superDeluxeCat, "slug = ?", "super-deluxe").Error; err != nil {
			return err
		}
		if err := db.WithContext(ctx).First(&familySuiteCat, "slug = ?", "family-suite").Error; err != nil {
			return err
		}
		if err := db.WithContext(ctx).First(&execPackCat, "slug = ?", "executive-pack").Error; err != nil {
			return err
		}

		rooms := []models.PMSRoom{
			{ID: uuid.New(), RoomNumber: 101, Floor: 1, CategoryID: superDeluxeCat.ID, Status: "available", CleanStatus: "clean"},
			{ID: uuid.New(), RoomNumber: 102, Floor: 1, CategoryID: superDeluxeCat.ID, Status: "available", CleanStatus: "dirty"},
			{ID: uuid.New(), RoomNumber: 103, Floor: 1, CategoryID: superDeluxeCat.ID, Status: "available", CleanStatus: "clean"},
			{ID: uuid.New(), RoomNumber: 104, Floor: 1, CategoryID: superDeluxeCat.ID, Status: "available", CleanStatus: "clean"},
			{ID: uuid.New(), RoomNumber: 105, Floor: 1, CategoryID: superDeluxeCat.ID, Status: "available", CleanStatus: "clean"},

			{ID: uuid.New(), RoomNumber: 201, Floor: 2, CategoryID: familySuiteCat.ID, Status: "available", CleanStatus: "clean"},
			{ID: uuid.New(), RoomNumber: 202, Floor: 2, CategoryID: familySuiteCat.ID, Status: "available", CleanStatus: "clean"},
			{ID: uuid.New(), RoomNumber: 203, Floor: 2, CategoryID: familySuiteCat.ID, Status: "available", CleanStatus: "dirty"},

			{ID: uuid.New(), RoomNumber: 301, Floor: 3, CategoryID: execPackCat.ID, Status: "available", CleanStatus: "clean"},
			{ID: uuid.New(), RoomNumber: 302, Floor: 3, CategoryID: execPackCat.ID, Status: "available", CleanStatus: "clean"},
		}
		for i := range rooms {
			if err := db.WithContext(ctx).Create(&rooms[i]).Error; err != nil {
				return err
			}
		}
	}

	if err := db.WithContext(ctx).Model(&models.POSTable{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		for i := 1; i <= 10; i++ {
			table := models.POSTable{
				ID:          uuid.New(),
				TableNumber: i,
				Capacity:    4,
				Area:        "Main Dining",
				Status:      "vacant",
			}
			if err := db.WithContext(ctx).Create(&table).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

