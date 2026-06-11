package main

import (
	"fmt"
	"log"

	"yoyo-server/internal/config"
	"yoyo-server/internal/database"
	"yoyo-server/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	fmt.Println("Seeding Waterpark Tickets...")

	tickets := []models.Ticket{
		{
			Title:       "Adult Splash Pass",
			Slug:        "adult-splash-pass",
			Description: "Full access to all high-speed slides, wave pool, and rain dance.",
			Price:       599,
			Category:    "Individual",
			Features:    []byte(`["All Slides", "Wave Pool", "Rain Dance", "Locker Access"]`),
			SortOrder:   1,
			IsActive:    true,
		},
		{
			Title:       "Kids Adventure Pass",
			Slug:        "kids-adventure-pass",
			Description: "Special access to Kids Fantasy Zone and mini-slides. (Height < 4ft)",
			Price:       399,
			Category:    "Individual",
			Features:    []byte(`["Kids Zone", "Mini Slides", "Toddler Pool"]`),
			SortOrder:   2,
			IsActive:    true,
		},
		{
			Title:       "Family Saver Pack (4 Persons)",
			Slug:        "family-saver-pack",
			Description: "Best value for families! Includes entry for 2 adults and 2 children.",
			Price:       1799,
			Category:    "Group",
			Features:    []byte(`["4 Entry Passes", "Reserved Table", "10% Food Discount"]`),
			SortOrder:   3,
			IsActive:    true,
		},
		{
			Title:       "Student Fun Day",
			Slug:        "student-fun-day",
			Description: "Discounted entry for students with valid ID. Minimum 5 students.",
			Price:       449,
			Category:    "Group",
			Features:    []byte(`["Full Access", "Group Entry", "Student Discount"]`),
			SortOrder:   4,
			IsActive:    true,
		},
		{
			Title:       "VIP Infinity Pass",
			Slug:        "vip-infinity-pass",
			Description: "No waiting in lines! Express entry to all slides and private lounge.",
			Price:       1299,
			Category:    "Premium",
			Features:    []byte(`["Express Queue", "Private Lounge", "Premium Locker", "Free Meal"]`),
			SortOrder:   5,
			IsActive:    true,
		},
	}

	for _, t := range tickets {
		var existing models.Ticket
		if err := db.Where("slug = ?", t.Slug).First(&existing).Error; err == nil {
			t.ID = existing.ID
			db.Save(&t)
			fmt.Printf("Updated ticket: %s\n", t.Title)
		} else {
			db.Create(&t)
			fmt.Printf("Created ticket: %s\n", t.Title)
		}
	}

	fmt.Println("\nSeeding Waterpark Hero Slides...")

	slides := []models.HeroSlide{
		{
			ImageURL:    "https://images.unsplash.com/photo-1542332213-9b5a5a3fad35?q=80&w=1600&auto=format&fit=crop",
			Title:       "Dive into Central India's Most Thrilling Waterpark!",
			Subtitle:    "Over 25+ massive slides, wave pools, and kids zones for the perfect family day out.",
			CTAURL:      "/tickets",
			CTALabel:     "Book Your Splash Now",
			SortOrder:   1,
			IsActive:    true,
		},
		{
			ImageURL:    "https://images.unsplash.com/photo-1519817650390-64a93db51149?q=80&w=1600&auto=format&fit=crop",
			Title:       "Experience the Ocean in Indore - Massive Wave Pool!",
			Subtitle:    "Ride the tides with our world-class wave system. Safe, fun, and purely refreshing.",
			CTAURL:      "/gallery",
			CTALabel:     "See All Attractions",
			SortOrder:   2,
			IsActive:    true,
		},
		{
			ImageURL:    "https://images.unsplash.com/photo-1582650625119-3a31f8fa2699?q=80&w=1600&auto=format&fit=crop",
			Title:       "High-Speed Thrills: The Vertical Drop Slide",
			Subtitle:    "Are you brave enough? Test your limits on our newest vertical drop experience.",
			CTAURL:      "/tickets",
			CTALabel:     "Grab Tickets - ₹499 Onwards",
			SortOrder:   3,
			IsActive:    true,
		},
	}

	for _, s := range slides {
		db.Create(&s)
		fmt.Printf("Created Hero Slide: %s\n", s.Title)
	}

	fmt.Println("\nSeeding Content Pages...")
	pages := []models.ContentPage{
		{
			Slug:  "privacy-policy",
			Title: "Privacy Policy – YOYO Fun N Foods",
			Content: `At YOYO Fun N Foods, we value your privacy and are committed to protecting your personal information.

### 1. Information We Collect
We may collect:
- Name
- Phone number
- Email address
- Booking details
- Payment-related identifiers (via Razorpay)

We do NOT store your card or payment credentials.

### 2. How We Use Your Information
Your data is used to:
- Process bookings
- Confirm tickets
- Provide customer support
- Send important updates related to your visit

### 3. Payment Security
All payments are processed securely via Razorpay. We do not store sensitive financial data.

### 4. Data Sharing
We do NOT sell or rent your personal data.
We may share data only with:
- Payment providers
- Legal authorities (if required)

### 5. Data Storage
Your information may be stored securely for operational and legal purposes.

### 6. Cookies
We may use cookies to improve user experience and performance tracking.

### 7. Your Rights
You can request:
- Access to your data
- Correction of incorrect data
- Deletion (subject to legal requirements)

### 8. Contact Us
For privacy-related queries:
Email: business@appnity.co.in`,
			IsPublished: true,
		},
		{
			Slug:  "terms-and-conditions",
			Title: "Terms & Conditions – YOYO Fun N Foods",
			Content: `### 1. General
By accessing and using this website, you agree to comply with these terms.

### 2. Ticket Booking
- All bookings are subject to availability.
- Prices may change without prior notice.
- Entry is only valid for the selected date.

### 3. Entry Rules
- Management reserves the right to deny entry.
- Guests must follow all safety rules and instructions.
- Proper attire (as per park guidelines) is mandatory.

### 4. Safety & Liability
- Visitors must follow lifeguard and staff instructions.
- YOYO Fun N Foods is not responsible for:
  - Personal injury due to negligence
  - Loss of personal belongings

### 5. Use of Facilities
- Misuse of rides or facilities may lead to removal without refund.
- Outside food or restricted items may not be allowed.

### 6. Intellectual Property
All website content (images, branding, design) belongs to YOYO Fun N Foods.

### 7. Modifications
We reserve the right to modify these terms at any time.

### 8. Contact
For any issues:
Email: business@appnity.co.in`,
			IsPublished: true,
		},
		{
			Slug:  "refund-policy",
			Title: "Refund Policy – YOYO Fun N Foods",
			Content: `### No Refund Policy

All ticket bookings are non-refundable and non-cancellable.

Once a booking is confirmed, no refunds will be issued under any circumstances including:
- Change of plans
- Weather conditions
- Personal reasons
- Late arrival or no-show

### Exception (Mistaken Payment Only)

Refunds will ONLY be considered if:
- Payment was made mistakenly (duplicate or incorrect transaction)

Conditions:
- Must be reported within 24 hours of transaction
- Valid proof must be provided
- Refund approval is at management discretion

### Processing Time
If approved:
- Refunds may take 5–10 business days

### Contact
For refund requests:
Email: business@appnity.co.in`,
			IsPublished: true,
		},
	}

	for _, p := range pages {
		var existing models.ContentPage
		if err := db.Where("slug = ?", p.Slug).First(&existing).Error; err == nil {
			p.ID = existing.ID
			db.Save(&p)
			fmt.Printf("Updated Page: %s\n", p.Title)
		} else {
			db.Create(&p)
			fmt.Printf("Created Page: %s\n", p.Title)
		}
	}

	fmt.Println("\nSeeding Complete!")
}
