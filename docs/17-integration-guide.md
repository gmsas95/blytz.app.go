# Integration Guide

## Overview

Third-party service integrations for the Blytz platform.

## Stripe (Payments)

### Setup
1. Create Stripe account: https://dashboard.stripe.com
2. Enable payment methods:
   - Cards (default)
   - FPX (Malaysia bank transfer)
   - GrabPay
   - Touch 'n Go eWallet

### Configuration
```bash
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

### Backend Implementation

#### Payment Intent Creation
```go
// integrations/stripe/payment.go
package stripe

import (
  "github.com/stripe/stripe-go/v76"
  "github.com/stripe/stripe-go/v76/paymentintent"
)

type PaymentService struct {
  secretKey string
}

func NewPaymentService(secretKey string) *PaymentService {
  stripe.Key = secretKey
  return &PaymentService{secretKey: secretKey}
}

func (s *PaymentService) CreatePaymentIntent(order *domain.Order) (*stripe.PaymentIntent, error) {
  params := &stripe.PaymentIntentParams{
    Amount:   stripe.Int64(int64(order.TotalAmount * 100)), // Convert to cents
    Currency: stripe.String(string(stripe.CurrencyMYR)),
    Metadata: map[string]string{
      "order_id": order.ID.String(),
      "customer_id": order.BuyerID.String(),
    },
    PaymentMethodTypes: stripe.StringSlice([]string{
      "card",
      "fpx",
      "grabpay",
    }),
  }
  
  return paymentintent.New(params)
}

func (s *PaymentService) ConfirmPayment(paymentIntentID string) (*stripe.PaymentIntent, error) {
  params := &stripe.PaymentIntentConfirmParams{}
  return paymentintent.Confirm(paymentIntentID, params)
}

func (s *PaymentService) Refund(paymentIntentID string, amount float64) (*stripe.Refund, error) {
  params := &stripe.RefundParams{
    PaymentIntent: stripe.String(paymentIntentID),
    Amount:        stripe.Int64(int64(amount * 100)),
  }
  return refund.New(params)
}
```

#### Webhook Handler
```go
// interfaces/http/handlers/stripe_webhook.go
func (h *StripeWebhookHandler) HandleWebhook(c *gin.Context) {
  payload, err := io.ReadAll(c.Request.Body)
  if err != nil {
    c.AbortWithStatus(400)
    return
  }
  
  sigHeader := c.GetHeader("Stripe-Signature")
  event, err := webhook.ConstructEvent(payload, sigHeader, h.webhookSecret)
  if err != nil {
    c.AbortWithStatus(400)
    return
  }
  
  switch event.Type {
  case "payment_intent.succeeded":
    var paymentIntent stripe.PaymentIntent
    err := json.Unmarshal(event.Data.Raw, &paymentIntent)
    if err == nil {
      h.handlePaymentSuccess(&paymentIntent)
    }
    
  case "payment_intent.payment_failed":
    var paymentIntent stripe.PaymentIntent
    err := json.Unmarshal(event.Data.Raw, &paymentIntent)
    if err == nil {
      h.handlePaymentFailure(&paymentIntent)
    }
  }
  
  c.Status(200)
}

func (h *StripeWebhookHandler) handlePaymentSuccess(pi *stripe.PaymentIntent) {
  orderID := pi.Metadata["order_id"]
  
  // Update order status
  order, _ := h.orderService.GetByID(orderID)
  order.Status = domain.OrderStatusPaid
  h.orderService.Update(order)
  
  // Create shipment
  h.shippingService.CreateShipment(order)
  
  // Notify buyer
  h.notificationService.Send(order.BuyerID, "payment_success", map[string]interface{}{
    "order_id": order.ID,
    "amount": float64(pi.Amount) / 100,
  })
}
```

### Frontend Implementation
```typescript
// components/payment/stripe-form.tsx
import { loadStripe } from '@stripe/stripe-js';
import { Elements, PaymentElement, useStripe, useElements } from '@stripe/react-stripe-js';

const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!);

function PaymentForm({ clientSecret }: { clientSecret: string }) {
  const stripe = useStripe();
  const elements = useElements();
  const [isProcessing, setIsProcessing] = useState(false);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!stripe || !elements) return;
    
    setIsProcessing(true);
    
    const { error, paymentIntent } = await stripe.confirmPayment({
      elements,
      confirmParams: {
        return_url: `${window.location.origin}/checkout/success`,
      },
    });
    
    if (error) {
      toast.error(error.message);
    }
    
    setIsProcessing(false);
  };
  
  return (
    <form onSubmit={handleSubmit}>
      <PaymentElement />
      <Button type="submit" disabled={!stripe || isProcessing}>
        {isProcessing ? 'Processing...' : 'Pay Now'}
      </Button>
    </form>
  );
}

// Usage
export function CheckoutPayment({ orderId }: { orderId: string }) {
  const [clientSecret, setClientSecret] = useState('');
  
  useEffect(() => {
    api.post('/payments/intent', { order_id: orderId })
      .then(res => setClientSecret(res.data.client_secret));
  }, [orderId]);
  
  if (!clientSecret) return <Loading />;
  
  return (
    <Elements stripe={stripePromise} options={{ clientSecret }}>
      <PaymentForm clientSecret={clientSecret} />
    </Elements>
  );
}
```

---

## LiveKit (Streaming)

### Setup
1. Sign up: https://livekit.io
2. Create project
3. Get API credentials

### Configuration
```bash
LIVEKIT_API_KEY=...
LIVEKIT_API_SECRET=...
LIVEKIT_WS_URL=wss://your-project.livekit.cloud
```

### Backend Implementation

#### Token Generation
```go
// integrations/livekit/client.go
package livekit

import (
  "github.com/livekit/server-sdk-go/v2/auth"
  lksdk "github.com/livekit/server-sdk-go/v2"
)

type Client struct {
  apiKey    string
  apiSecret string
  wsURL     string
}

func (c *Client) GenerateToken(room string, identity string, isPublisher bool) (string, error) {
  at := auth.NewAccessToken(c.apiKey, c.apiSecret)
  
  grant := &auth.VideoGrant{
    RoomJoin: true,
    Room:     room,
  }
  
  if isPublisher {
    grant.RoomCreate = true
    grant.RoomAdmin = true
    grant.CanPublish = true
    grant.CanSubscribe = true
  } else {
    grant.CanPublish = false
    grant.CanSubscribe = true
  }
  
  at.AddGrant(grant).
    SetIdentity(identity).
    SetValidFor(24 * time.Hour)
  
  return at.ToJWT()
}

func (c *Client) CreateRoom(ctx context.Context, name string) (*livekit.Room, error) {
  roomClient := lksdk.NewRoomServiceClient(c.wsURL, c.apiKey, c.apiSecret)
  
  return roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
    Name:            name,
    EmptyTimeout:    600,  // 10 minutes
    MaxParticipants: 1000,
  })
}
```

### Frontend Implementation
```typescript
// hooks/use-livekit.ts
import { Room, RoomEvent } from 'livekit-client';

export function useLiveKit(roomName: string, token: string) {
  const [room, setRoom] = useState<Room | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [participants, setParticipants] = useState<Participant[]>([]);
  
  useEffect(() => {
    const room = new Room({
      adaptiveStream: true,
      dynacast: true,
    });
    
    room.on(RoomEvent.Connected, () => setIsConnected(true));
    room.on(RoomEvent.Disconnected, () => setIsConnected(false));
    room.on(RoomEvent.ParticipantConnected, (p) => {
      setParticipants(prev => [...prev, p]);
    });
    room.on(RoomEvent.ParticipantDisconnected, (p) => {
      setParticipants(prev => prev.filter(x => x.sid !== p.sid));
    });
    
    room.connect(process.env.NEXT_PUBLIC_LIVEKIT_URL!, token);
    setRoom(room);
    
    return () => {
      room.disconnect();
    };
  }, [roomName, token]);
  
  const toggleCamera = async () => {
    if (!room) return;
    await room.localParticipant.setCameraEnabled(!room.localParticipant.isCameraEnabled);
  };
  
  const toggleMic = async () => {
    if (!room) return;
    await room.localParticipant.setMicrophoneEnabled(!room.localParticipant.isMicrophoneEnabled);
  };
  
  return {
    room,
    isConnected,
    participants,
    toggleCamera,
    toggleMic,
  };
}
```

---

## Socket.io (Real-time Chat)

### Backend Implementation
```go
// infrastructure/websocket/hub.go
package websocket

import (
  socketio "github.com/googollee/go-socket.io"
)

type Hub struct {
  server *socketio.Server
  redis  *redis.Client
}

func NewHub(redisClient *redis.Client) *Hub {
  server := socketio.NewServer(nil)
  
  hub := &Hub{
    server: server,
    redis:  redisClient,
  }
  
  server.OnConnect("/", func(s socketio.Conn) error {
    log.Printf("Client connected: %s", s.ID())
    return nil
  })
  
  server.OnEvent("/", "stream:join", func(s socketio.Conn, data map[string]string) {
    streamID := data["stream_id"]
    s.Join(streamID)
    
    // Update viewer count
    hub.incrementViewerCount(streamID)
    
    // Broadcast to room
    server.BroadcastToRoom("/", streamID, "user:joined", map[string]string{
      "user_id": s.ID(),
    })
  })
  
  server.OnEvent("/", "chat:message", func(s socketio.Conn, data map[string]string) {
    streamID := data["stream_id"]
    message := data["content"]
    
    // Save to database
    msg := hub.saveMessage(s.Context(), streamID, s.ID(), message)
    
    // Broadcast to room
    server.BroadcastToRoom("/", streamID, "chat:message", map[string]interface{}{
      "id": msg.ID,
      "user_id": s.ID(),
      "content": message,
      "created_at": msg.CreatedAt,
    })
  })
  
  server.OnEvent("/", "auction:bid", func(s socketio.Conn, data map[string]interface{}) {
    // Handle bid via auction service
    bid, err := hub.auctionService.PlaceBid(s.Context(), data)
    if err != nil {
      s.Emit("auction:error", err.Error())
      return
    }
    
    // Broadcast to auction room
    auctionID := data["auction_id"].(string)
    server.BroadcastToRoom("/", "auction:"+auctionID, "auction:bid", bid)
  })
  
  server.OnDisconnect("/", func(s socketio.Conn, reason string) {
    log.Printf("Client disconnected: %s", s.ID())
  })
  
  return hub
}

func (h *Hub) Start() {
  go h.server.Serve()
  defer h.server.Close()
}
```

### Frontend Implementation
```typescript
// lib/socket.ts
import { io, Socket } from 'socket.io-client';

let socket: Socket | null = null;

export function getSocket(): Socket {
  if (!socket) {
    socket = io(process.env.NEXT_PUBLIC_SOCKET_URL!, {
      auth: {
        token: getAuthToken(),
      },
    });
  }
  return socket;
}

// hooks/use-chat.ts
export function useChat(streamId: string) {
  const [messages, setMessages] = useState<Message[]>([]);
  const socket = getSocket();
  
  useEffect(() => {
    socket.emit('stream:join', { stream_id: streamId });
    
    socket.on('chat:message', (msg: Message) => {
      setMessages(prev => [...prev, msg]);
    });
    
    socket.on('user:joined', (data) => {
      // Show notification
    });
    
    return () => {
      socket.off('chat:message');
      socket.off('user:joined');
    };
  }, [streamId]);
  
  const sendMessage = (content: string) => {
    socket.emit('chat:message', {
      stream_id: streamId,
      content,
    });
  };
  
  return { messages, sendMessage };
}
```

---

## NinjaVan (Shipping)

### Setup
1. Register: https://www.ninjavan.co
2. Get API credentials
3. Configure webhooks

### Configuration
```bash
NINJAVAN_BASE_URL=https://api.ninjavan.co
NINJAVAN_API_KEY=...
NINJAVAN_API_SECRET=...
NINJAVAN_WEBHOOK_SECRET=...
```

### Backend Implementation
```go
// integrations/ninjavan/client.go
package ninjavan

import (
  "net/http"
  "encoding/json"
)

type Client struct {
  baseURL    string
  apiKey     string
  apiSecret  string
  httpClient *http.Client
}

func (c *Client) CreateOrder(ctx context.Context, order *domain.Order) (*OrderResponse, error) {
  payload := map[string]interface{}{
    "service_type": "Parcel",
    "service_level": "Standard",
    "from": map[string]string{
      "name": order.Seller.StoreName,
      "phone": order.Seller.Phone,
      "address": order.Seller.Address,
      "city": order.Seller.City,
      "state": order.Seller.State,
      "postcode": order.Seller.PostalCode,
      "country": "MY",
    },
    "to": map[string]string{
      "name": order.ShippingAddress.RecipientName,
      "phone": order.ShippingAddress.Phone,
      "address": order.ShippingAddress.Line1,
      "city": order.ShippingAddress.City,
      "state": order.ShippingAddress.State,
      "postcode": order.ShippingAddress.PostalCode,
      "country": "MY",
    },
    "parcel_job": map[string]interface{}{
      "pickup_service_type": "Scheduled",
      "pickup_service_level": "Standard",
      "pickup_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
      "delivery_start_date": time.Now().Add(48 * time.Hour).Format("2006-01-02"),
      "delivery_timeslot": map[string]string{
        "start_time": "09:00",
        "end_time": "22:00",
        "timezone": "Asia/Kuala_Lumpur",
      },
      "dimensions": map[string]float64{
        "weight": order.WeightKg,
      },
      "cash_on_delivery": order.TotalAmount,
      "insured_value": order.TotalAmount,
    },
  }
  
  reqBody, _ := json.Marshal(payload)
  req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/4.2/orders", bytes.NewReader(reqBody))
  req.Header.Set("Authorization", "Bearer "+c.getAccessToken())
  req.Header.Set("Content-Type", "application/json")
  
  resp, err := c.httpClient.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  
  var result OrderResponse
  if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
    return nil, err
  }
  
  return &result, nil
}

func (c *Client) GetTracking(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
  req, _ := http.NewRequestWithContext(ctx, "GET", 
    c.baseURL+"/4.2/tracking?tracking_number="+trackingNumber, nil)
  req.Header.Set("Authorization", "Bearer "+c.getAccessToken())
  
  resp, err := c.httpClient.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  
  var result TrackingInfo
  json.NewDecoder(resp.Body).Decode(&result)
  return &result, nil
}
```

### Webhook Handler
```go
func (h *NinjaVanWebhookHandler) HandleWebhook(c *gin.Context) {
  // Verify webhook signature
  signature := c.GetHeader("X-NinjaVan-Signature")
  payload, _ := io.ReadAll(c.Request.Body)
  
  if !verifySignature(payload, signature, h.webhookSecret) {
    c.AbortWithStatus(401)
    return
  }
  
  var event NinjaVanEvent
  json.Unmarshal(payload, &event)
  
  switch event.Status {
  case "Picked_Up":
    h.handlePickup(event)
  case "In_Transit":
    h.handleTransit(event)
  case "Out_For_Delivery":
    h.handleOutForDelivery(event)
  case "Delivered":
    h.handleDelivery(event)
  }
  
  c.Status(200)
}
```

---

## Cloudflare R2 (Storage)

### Configuration
```bash
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET_NAME=blytz-storage
R2_PUBLIC_URL=https://cdn.blytz.app
```

### Backend Implementation
```go
// infrastructure/storage/r2.go
package storage

import (
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/credentials"
  "github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Storage struct {
  client *s3.Client
  bucket string
  publicURL string
}

func NewR2Storage(accountID, accessKey, secretKey, bucket string) *R2Storage {
  r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
    return aws.Endpoint{
      URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
    }, nil
  })
  
  cfg := aws.Config{
    EndpointResolverWithOptions: r2Resolver,
    Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
    Region: "auto",
  }
  
  return &R2Storage{
    client: s3.NewFromConfig(cfg),
    bucket: bucket,
    publicURL: fmt.Sprintf("https://cdn.blytz.app"),
  }
}

func (s *R2Storage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
  _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
    Bucket:      aws.String(s.bucket),
    Key:         aws.String(key),
    Body:        bytes.NewReader(data),
    ContentType: aws.String(contentType),
  })
  
  if err != nil {
    return "", err
  }
  
  return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}

func (s *R2Storage) GeneratePresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
  presignClient := s3.NewPresignClient(s.client)
  
  req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
    Bucket: aws.String(s.bucket),
    Key:    aws.String(key),
  }, s3.WithPresignExpires(expiry))
  
  if err != nil {
    return "", err
  }
  
  return req.URL, nil
}
```

---

*Last updated: 2025-02-05*


---

## WebSocket (Real-time Bidding)

### Overview
WebSocket provides real-time bid updates and auction events to connected clients.

### Connection
```javascript
const ws = new WebSocket('wss://api.blytz.app/ws/auctions/{auction_id}');

ws.onopen = () => {
  console.log('Connected to auction');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleMessage(message);
};

ws.onclose = () => {
  console.log('Disconnected');
};
```

### Message Types

#### Bid Update
```json
{
  "type": "bid",
  "auction_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "bid_id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440002",
    "amount": 1250.00,
    "bid_count": 15,
    "is_auto_bid": false
  },
  "timestamp": "2025-02-05T10:30:00Z"
}
```

#### Auction Started
```json
{
  "type": "auction_started",
  "auction_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "start_time": "2025-02-05T10:00:00Z",
    "end_time": "2025-02-05T11:00:00Z"
  },
  "timestamp": "2025-02-05T10:00:00Z"
}
```

#### Auction Ended
```json
{
  "type": "auction_ended",
  "auction_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "winner_id": "550e8400-e29b-41d4-a716-446655440002",
    "winning_amount": 1500.00,
    "total_bids": 25
  },
  "timestamp": "2025-02-05T11:00:00Z"
}
```

#### Viewer Count
```json
{
  "type": "viewer_count",
  "auction_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "count": 45
  },
  "timestamp": "2025-02-05T10:30:00Z"
}
```

#### Auction Extended
```json
{
  "type": "auction_extended",
  "auction_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "new_end_time": "2025-02-05T11:05:00Z",
    "extension_seconds": 300
  },
  "timestamp": "2025-02-05T11:00:00Z"
}
```

### Backend Architecture

The WebSocket system uses Redis Pub/Sub for cross-instance communication:

```
Client A (Server 1) -----> Redis Pub/Sub <----- Client B (Server 2)
      |                                               |
      |---> Bid placed ---> Broadcast to room -------->|
```

### Implementation

```go
// infrastructure/websocket/hub.go

type Hub struct {
  rooms map[string]*Room
  redisClient *redis.Client
  eventBus *redisMessaging.EventBus
}

func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request, auctionID, userID string) {
  // Upgrade to WebSocket
  conn, err := h.upgrader.Upgrade(w, r, nil)
  if err != nil {
    return
  }
  
  // Get or create room
  room := h.getOrCreateRoom(auctionID)
  
  // Create client
  client := &Client{
    hub:       h,
    conn:      conn,
    room:      room,
    userID:    userID,
    auctionID: auctionID,
  }
  
  // Register and broadcast viewer count
  room.addClient(client)
  h.broadcastViewerCount(auctionID)
  
  // Start pumps
  go client.writePump()
  go client.readPump()
}

func (h *Hub) handleEvent(event redisMessaging.Event) {
  switch event.Type {
  case redisMessaging.EventBidPlaced:
    h.broadcastToRoom(event.AuctionID, Message{
      Type: "bid",
      Data: event.Payload,
    })
  // ... other events
  }
}
```

---

## Cloudflare R2 (Image Storage)

### Setup
1. Create Cloudflare account: https://dash.cloudflare.com
2. Enable R2: https://dash.cloudflare.com/?to=/:account/r2
3. Create bucket
4. Create API token with R2 permissions

### Configuration
```bash
R2_ACCOUNT_ID=your-account-id
R2_ACCESS_KEY_ID=your-access-key
R2_SECRET_ACCESS_KEY=your-secret-key
R2_BUCKET_NAME=blytz-storage
R2_PUBLIC_URL=https://pub-xxx.r2.dev
R2_CDN_URL=https://cdn.blytz.app  # Optional: custom domain
```

### Backend Implementation

```go
// infrastructure/storage/r2/client.go

type Client struct {
  s3Client  *s3.Client
  bucket    string
  publicURL string
  cdnURL    string
}

func NewClient(cfg Config) (*Client, error) {
  endpointResolver := aws.EndpointResolverWithOptionsFunc(
    func(service, region string, options ...interface{}) (aws.Endpoint, error) {
      return aws.Endpoint{
        URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
      }, nil
    })

  awsCfg := aws.Config{
    EndpointResolverWithOptions: endpointResolver,
    Credentials: credentials.NewStaticCredentialsProvider(
      cfg.AccessKeyID, cfg.SecretAccessKey, ""),
    Region: "auto",
  }

  return &Client{
    s3Client: s3.NewFromConfig(awsCfg),
    bucket:   cfg.BucketName,
  }, nil
}

func (c *Client) Upload(ctx context.Context, reader io.Reader, 
  filename string, opts UploadOptions) (*UploadResult, error) {
  
  // Generate unique key
  key := fmt.Sprintf("%s/%s%s", opts.Folder, uuid.New().String(), 
    filepath.Ext(filename))
  
  // Upload
  _, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
    Bucket:      aws.String(c.bucket),
    Key:         aws.String(key),
    Body:        reader,
    ContentType: aws.String(opts.ContentType),
  })
  
  return &UploadResult{
    URL: c.getPublicURL(key),
    Key: key,
  }, nil
}
```

### Upload Service

```go
// application/upload/service.go

type Service struct {
  r2Client *r2.Client
}

func (s *Service) UploadProductImage(ctx context.Context, file io.Reader, 
  filename string, size int64) (*UploadResult, error) {
  
  // Validate
  if size > 10*1024*1024 {
    return nil, fmt.Errorf("file too large: max 10MB")
  }
  
  result, err := s.r2Client.Upload(ctx, file, filename, r2.UploadOptions{
    Folder:            "products",
    GenerateThumbnail: true,
    AllowedTypes:      []string{"image/"},
  })
  
  return &UploadResult{
    URL:          result.URL,
    ThumbnailURL: result.ThumbnailURL,
  }, nil
}
```

### API Endpoints

#### Upload Product Image
```bash
POST /api/v1/uploads/product-image
Content-Type: multipart/form-data
Authorization: Bearer {token}

file: [binary image data]
```

**Response:**
```json
{
  "data": {
    "url": "https://cdn.blytz.app/products/xxx.jpg",
    "thumbnail_url": "https://cdn.blytz.app/products/xxx-thumb.jpg",
    "key": "products/xxx.jpg",
    "content_type": "image/jpeg",
    "size": 2048000
  }
}
```

#### Upload Avatar
```bash
POST /api/v1/uploads/avatar
Content-Type: multipart/form-data
Authorization: Bearer {token}

file: [binary image data]
```

#### Upload Stream Thumbnail
```bash
POST /api/v1/uploads/stream-thumbnail
Content-Type: multipart/form-data
Authorization: Bearer {token}

file: [binary image data]
```

#### Delete File
```bash
DELETE /api/v1/uploads
Authorization: Bearer {token}
Content-Type: application/json

{
  "key": "products/xxx.jpg"
}
```

### File Structure in R2
```
blytz-storage/
├── products/
│   ├── xxx.jpg
│   ├── yyy.png
│   └── zzz.webp
├── avatars/
│   ├── user-1.jpg
│   └── user-2.png
├── streams/
│   ├── stream-1-thumb.jpg
│   └── stream-2-thumb.jpg
└── temp/
    └── [temporary uploads]
```

### Frontend Usage

```typescript
// Upload product image
async function uploadProductImage(file: File): Promise<UploadResult> {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('/api/v1/uploads/product-image', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: formData,
  });
  
  return response.json();
}

// Use in product creation
async function createProductWithImage(productData, imageFile) {
  // 1. Upload image first
  const uploadResult = await uploadProductImage(imageFile);
  
  // 2. Create product with image URL
  const product = await createProduct({
    ...productData,
    images: [{ url: uploadResult.data.url, is_primary: true }],
  });
  
  return product;
}
```

### Image Optimization (Optional)

For production, use Cloudflare Images or Image Resizing:

```html
<!-- Original image -->
<img src="https://cdn.blytz.app/products/xxx.jpg" />

<!-- Resized via Cloudflare -->
<img src="https://cdn.blytz.app/cdn-cgi/image/width=500/products/xxx.jpg" />

<!-- With WebP conversion -->
<img src="https://cdn.blytz.app/cdn-cgi/image/width=500,format=webp/products/xxx.jpg" />
```

---

*Last updated: 2025-02-05*
