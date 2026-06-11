package services

import (
	"context"

	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
)

type DashboardService struct {
	tickets  *repositories.TicketRepository
	bookings *repositories.BookingRepository
	messages *repositories.ContactMessageRepository
}

type DashboardStats struct {
	TotalBookings   int64                     `json:"total_bookings"`
	TotalRevenue    int64                     `json:"total_revenue"`
	ActiveTickets   int64                     `json:"active_tickets"`
	TotalCustomers  int64                     `json:"total_customers"`
	TotalMessages   int64                     `json:"total_messages"`
	RecentBookings  []models.Booking          `json:"recent_bookings"`
	RevenueChart    []repositories.ChartPoint `json:"revenue_chart_data"`
	BookingGrowth   []repositories.ChartPoint `json:"booking_growth_stats"`
	LowStockTickets []models.Ticket           `json:"low_stock_tickets"`
}

func NewDashboardService(tickets *repositories.TicketRepository, bookings *repositories.BookingRepository, messages *repositories.ContactMessageRepository) *DashboardService {
	return &DashboardService{tickets: tickets, bookings: bookings, messages: messages}
}

func (s *DashboardService) Stats(ctx context.Context) (*DashboardStats, error) {
	totalBookings, err := s.bookings.Count(ctx)
	if err != nil {
		return nil, err
	}
	totalRevenue, err := s.bookings.Revenue(ctx)
	if err != nil {
		return nil, err
	}
	allActiveTickets, totalActiveTickets, err := s.tickets.ListAdmin(ctx, repositories.TicketFilter{Status: "active"}, 1, 100)
	if err != nil {
		return nil, err
	}
	totalCustomers, err := s.bookings.DistinctCustomers(ctx)
	if err != nil {
		return nil, err
	}
	totalMessages, err := s.messages.Count(ctx)
	if err != nil {
		return nil, err
	}
	recent, err := s.bookings.Recent(ctx, 8)
	if err != nil {
		return nil, err
	}
	revenueChart, err := s.bookings.RevenueChart(ctx, 30)
	if err != nil {
		return nil, err
	}
	bookingGrowth, err := s.bookings.BookingGrowth(ctx, 30)
	if err != nil {
		return nil, err
	}

	lowStock := make([]models.Ticket, 0)
	for _, ticket := range allActiveTickets {
		if ticket.Stock <= 10 {
			lowStock = append(lowStock, ticket)
		}
	}

	return &DashboardStats{
		TotalBookings:   totalBookings,
		TotalRevenue:    totalRevenue,
		ActiveTickets:   totalActiveTickets,
		TotalCustomers:  totalCustomers,
		TotalMessages:   totalMessages,
		RecentBookings:  recent,
		RevenueChart:    revenueChart,
		BookingGrowth:   bookingGrowth,
		LowStockTickets: lowStock,
	}, nil
}
