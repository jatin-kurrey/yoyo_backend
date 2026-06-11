# YOYO API Collection

Base URL: `http://localhost:8080/api`

Admin endpoints require:

```http
Authorization: Bearer <jwt>
```

## Health

- `GET /health`

## Admin Auth

- `POST /admin/auth/login`
  - body: `{ "email": "...", "password": "..." }`
- `GET /admin/auth/me`
- `POST /admin/auth/logout`

## Public

- `GET /tickets`
- `GET /tickets/:slug`
- `POST /contact`
  - body: `{ "name": "...", "email": "...", "phone": "...", "subject": "...", "message": "..." }`
- `GET /settings/public`

## Bookings and Payments

- `POST /bookings/create-order`
  - body: `{ "customer_name": "...", "customer_email": "...", "customer_phone": "...", "ticket_id": "...", "quantity": 2, "visit_date": "2026-05-01" }`
- `POST /bookings/verify-payment`
  - body: `{ "razorpay_order_id": "...", "razorpay_payment_id": "...", "razorpay_signature": "..." }`
- `POST /webhooks/razorpay`
  - header: `X-Razorpay-Signature`

## Admin Dashboard

- `GET /admin/dashboard/stats`

## Admin Tickets

- `GET /admin/tickets?page=1&limit=20&search=&status=active`
- `POST /admin/tickets`
- `GET /admin/tickets/:id`
- `PATCH /admin/tickets/:id`
- `DELETE /admin/tickets/:id`
- `PATCH /admin/tickets/:id/toggle-status`

Ticket body:

```json
{
  "title": "Standard Pass",
  "slug": "standard-pass",
  "description": "Single entry",
  "price": 499,
  "original_price": 699,
  "category": "general",
  "features": ["All Day Entry"],
  "validity": "Valid for selected visit date",
  "stock": 100,
  "is_active": true,
  "sort_order": 1
}
```

## Admin Bookings

- `GET /admin/bookings?page=1&limit=20&search=&status=&payment_status=&date_from=&date_to=`
- `GET /admin/bookings/:id`
- `PATCH /admin/bookings/:id/status`
  - body: `{ "status": "confirmed" }`

## Admin Messages

- `GET /admin/messages?page=1&limit=20&search=&status=`
- `GET /admin/messages/:id`
- `PATCH /admin/messages/:id/status`
  - body: `{ "status": "read" }`
- `DELETE /admin/messages/:id`

## Admin Settings

- `GET /admin/settings`
- `PATCH /admin/settings`

## Admin Users

- `GET /admin/users`
- `POST /admin/users`
- `PATCH /admin/users/:id`
- `DELETE /admin/users/:id`

## Audit Logs

- `GET /admin/audit-logs?page=1&limit=20&module=&action=&admin_id=`

## Uploads

- `POST /admin/uploads`
  - multipart field: `file`
  - local fallback returns `{ "url": "/uploads/<file>" }`
