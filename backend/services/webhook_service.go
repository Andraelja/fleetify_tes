package services

import (
    "backend/app/models"
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type WebhookService interface {
    SendPurchasingNotification(purchasing models.Purchasing) error
}

type webhookService struct {
    webhookURL string
}

func NewWebhookService(webhookURL string) WebhookService {
    return &webhookService{
        webhookURL: webhookURL,
    }
}

// Struktur payload webhook
type WebhookPayload struct {
    Event       string                   `json:"event"`
    Timestamp   string                   `json:"timestamp"`
    Purchasing  WebhookPurchasingData    `json:"purchasing"`
}

type WebhookPurchasingData struct {
    ID         uint                      `json:"id"`
    Date       string                    `json:"date"`
    SupplierID uint                      `json:"supplier_id"`
    Supplier   string                    `json:"supplier_name"`
    UserID     uint                      `json:"user_id"`
    Username   string                    `json:"username"`
    GrandTotal int64                     `json:"grand_total"`
    Details    []WebhookDetailData       `json:"details"`
}

type WebhookDetailData struct {
    ItemID   uint   `json:"item_id"`
    ItemName string `json:"item_name"`
    Qty      int    `json:"qty"`
    Price    int64  `json:"price"`
    SubTotal int64  `json:"sub_total"`
}

func (s *webhookService) SendPurchasingNotification(purchasing models.Purchasing) error {
    // Siapkan payload
    details := make([]WebhookDetailData, len(purchasing.Details))
    for i, detail := range purchasing.Details {
        details[i] = WebhookDetailData{
            ItemID:   detail.ItemID,
            ItemName: detail.Item.Name,
            Qty:      detail.Qty,
            Price:    detail.Price,
            SubTotal: detail.SubTotal,
        }
    }

    payload := WebhookPayload{
        Event:     "purchasing.created",
        Timestamp: time.Now().Format(time.RFC3339),
        Purchasing: WebhookPurchasingData{
            ID:         purchasing.ID,
            Date:       purchasing.Date.Format("2006-01-02"),
            SupplierID: purchasing.SupplierID,
            Supplier:   purchasing.Supplier.Name,
            UserID:     purchasing.UserID,
            Username:   purchasing.User.Username,
            GrandTotal: purchasing.GrandTotal,
            Details:    details,
        },
    }

    // Convert to JSON
    jsonData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshaling webhook payload: %v", err)
        return err
    }

    // Kirim ke webhook (asynchronous - tidak mengganggu flow utama)
    go func() {
        client := &http.Client{
            Timeout: 10 * time.Second,
        }

        req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(jsonData))
        if err != nil {
            log.Printf("Error creating webhook request: %v", err)
            return
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Webhook-Event", "purchasing.created")

        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error sending webhook: %v", err)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            log.Printf("Webhook sent successfully to %s for purchasing ID: %d", s.webhookURL, purchasing.ID)
        } else {
            log.Printf("Webhook failed with status code: %d for purchasing ID: %d", resp.StatusCode, purchasing.ID)
        }
    }()

    return nil
}