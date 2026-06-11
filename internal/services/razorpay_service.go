package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"yoyo-server/internal/config"
)

type RazorpayService struct {
	cfg    *config.Config
	client *http.Client
}

type RazorpayOrder struct {
	ID       string `json:"id"`
	Entity   string `json:"entity,omitempty"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
	Status   string `json:"status"`
}

func NewRazorpayService(cfg *config.Config) *RazorpayService {
	return &RazorpayService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *RazorpayService) CreateOrder(ctx context.Context, amount int64, receipt string) (*RazorpayOrder, error) {
	if !s.cfg.RazorpayEnabled || s.cfg.RazorpayKeyID == "" || s.cfg.RazorpayKeySecret == "" {
		return nil, ErrPaymentGatewayDisabled
	}

	payload := map[string]interface{}{
		"amount":          amount,
		"currency":        "INR",
		"receipt":         receipt,
		"payment_capture": true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.razorpay.com/v1/orders", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.cfg.RazorpayKeyID, s.cfg.RazorpayKeySecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("razorpay order failed with status %d", resp.StatusCode)
	}

	var order RazorpayOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *RazorpayService) VerifyPaymentSignature(orderID string, paymentID string, signature string) bool {
	payload := orderID + "|" + paymentID
	return verifyHMACSHA256(payload, signature, s.cfg.RazorpayKeySecret)
}

func (s *RazorpayService) VerifyWebhookSignature(body []byte, signature string) bool {
	return verifyHMACSHA256(string(body), signature, s.cfg.RazorpayWebhookSecret)
}

func verifyHMACSHA256(payload string, signature string, secret string) bool {
	if secret == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
