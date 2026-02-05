package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blytz/live/backend/internal/domain/auction"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	EventBidPlaced       = "bid.placed"
	EventAuctionStarted  = "auction.started"
	EventAuctionEnded    = "auction.ended"
	EventAuctionExtended = "auction.extended"
)

type EventBus struct {
	client *redis.Client
	prefix string
}

type Event struct {
	Type      string                 `json:"type"`
	AuctionID string                 `json:"auction_id"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

type Subscription struct {
	pubsub *redis.PubSub
	ch     <-chan *redis.Message
}

func (s *Subscription) Channel() <-chan *redis.Message {
	return s.ch
}

func (s *Subscription) Close() error {
	return s.pubsub.Close()
}

func NewEventBus(client *redis.Client) *EventBus {
	return &EventBus{
		client: client,
		prefix: "blytz:events:",
	}
}

func (b *EventBus) PublishBidPlaced(ctx context.Context, auctionID uuid.UUID, bid *auction.Bid) error {
	return b.publish(ctx, auctionID, EventBidPlaced, map[string]interface{}{
		"bid_id":     bid.ID.String(),
		"user_id":    bid.UserID.String(),
		"amount":     bid.Amount,
		"is_auto_bid": bid.IsAutoBid,
		"bid_time":   bid.BidTime,
	})
}

func (b *EventBus) PublishAuctionStarted(ctx context.Context, auctionID uuid.UUID) error {
	return b.publish(ctx, auctionID, EventAuctionStarted, map[string]interface{}{
		"started_at": time.Now(),
	})
}

func (b *EventBus) PublishAuctionEnded(ctx context.Context, auctionID uuid.UUID, winnerID *uuid.UUID) error {
	payload := map[string]interface{}{
		"ended_at": time.Now(),
	}
	if winnerID != nil {
		payload["winner_id"] = winnerID.String()
	}
	return b.publish(ctx, auctionID, EventAuctionEnded, payload)
}

func (b *EventBus) PublishAuctionExtended(ctx context.Context, auctionID uuid.UUID, newEndTime time.Time) error {
	return b.publish(ctx, auctionID, EventAuctionExtended, map[string]interface{}{
		"new_end_time": newEndTime,
	})
}

func (b *EventBus) publish(ctx context.Context, auctionID uuid.UUID, eventType string, payload map[string]interface{}) error {
	e := Event{
		Type:      eventType,
		AuctionID: auctionID.String(),
		Timestamp: time.Now(),
		Payload:   payload,
	}

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	channel := b.channelName(auctionID)
	if err := b.client.Publish(ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	globalChannel := b.prefix + "all"
	if err := b.client.Publish(ctx, globalChannel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish to global channel: %w", err)
	}

	return nil
}

func (b *EventBus) Subscribe(ctx context.Context, auctionID uuid.UUID) (*Subscription, error) {
	channel := b.channelName(auctionID)
	pubsub := b.client.Subscribe(ctx, channel)

	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &Subscription{
		pubsub: pubsub,
		ch:     pubsub.Channel(),
	}, nil
}

func (b *EventBus) SubscribeGlobal(ctx context.Context) (*Subscription, error) {
	channel := b.prefix + "all"
	pubsub := b.client.Subscribe(ctx, channel)

	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &Subscription{
		pubsub: pubsub,
		ch:     pubsub.Channel(),
	}, nil
}

func (b *EventBus) channelName(auctionID uuid.UUID) string {
	return fmt.Sprintf("%sauction:%s", b.prefix, auctionID.String())
}