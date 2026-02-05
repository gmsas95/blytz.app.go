# Testing Strategy

## Overview

Comprehensive testing approach covering unit, integration, and E2E tests across frontend, backend, and mobile.

## Testing Pyramid

```
        /\
       /  \
      / E2E \          # Cypress/Detox (10%)
     /--------\
    /Integration\       # API/Component (30%)
   /--------------\
  /    Unit Tests   \    # Jest/Vitest (60%)
 /--------------------\
```

## Frontend Testing (Next.js)

### Unit Tests
**Framework:** Vitest + React Testing Library

```typescript
// components/auction/__tests__/bid-form.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { BidForm } from '../bid-form';

describe('BidForm', () => {
  const mockPlaceBid = vi.fn();
  
  it('renders bid form correctly', () => {
    render(
      <BidForm 
        currentBid={1000}
        minIncrement={50}
        onPlaceBid={mockPlaceBid}
        isActive={true}
      />
    );
    
    expect(screen.getByText('Current Bid: RM 1,000')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter bid amount')).toBeInTheDocument();
  });
  
  it('validates minimum bid amount', async () => {
    render(<BidForm currentBid={1000} minIncrement={50} onPlaceBid={mockPlaceBid} isActive={true} />);
    
    const input = screen.getByPlaceholderText('Enter bid amount');
    fireEvent.change(input, { target: { value: '1000' } });
    fireEvent.click(screen.getByText('Place Bid'));
    
    expect(await screen.findByText('Bid must be at least RM 1,050')).toBeInTheDocument();
  });
  
  it('submits bid when valid', async () => {
    render(<BidForm currentBid={1000} minIncrement={50} onPlaceBid={mockPlaceBid} isActive={true} />);
    
    fireEvent.change(screen.getByPlaceholderText('Enter bid amount'), { target: { value: '1100' } });
    fireEvent.click(screen.getByText('Place Bid'));
    
    expect(mockPlaceBid).toHaveBeenCalledWith(1100);
  });
});
```

### Component Tests
```typescript
// components/auction/__tests__/auction-card.test.tsx
import { render, screen } from '@testing-library/react';
import { AuctionCard } from '../auction-card';

describe('AuctionCard', () => {
  const mockAuction = {
    id: '1',
    product: { name: 'Test Product', imageUrl: '/test.jpg' },
    seller: { storeName: 'Test Store' },
    currentBid: 1000,
    endTime: new Date(Date.now() + 3600000),
  };
  
  it('displays auction information', () => {
    render(<AuctionCard auction={mockAuction} />);
    
    expect(screen.getByText('Test Product')).toBeInTheDocument();
    expect(screen.getByText('Test Store')).toBeInTheDocument();
    expect(screen.getByText('RM 1,000')).toBeInTheDocument();
  });
  
  it('links to auction detail page', () => {
    render(<AuctionCard auction={mockAuction} />);
    
    expect(screen.getByRole('link')).toHaveAttribute('href', '/auctions/1');
  });
});
```

### Hook Tests
```typescript
// hooks/__tests__/use-auction.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { useAuction } from '../use-auction';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const wrapper = ({ children }) => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);

describe('useAuction', () => {
  it('fetches auction data', async () => {
    const { result } = renderHook(() => useAuction('123'), { wrapper });
    
    await waitFor(() => {
      expect(result.current.auction).toBeDefined();
    });
    
    expect(result.current.auction.id).toBe('123');
  });
});
```

### E2E Tests
**Framework:** Cypress

```typescript
// cypress/e2e/auction.cy.ts
describe('Auction Flow', () => {
  beforeEach(() => {
    cy.login('test@example.com', 'password');
  });
  
  it('user can place a bid', () => {
    cy.visit('/auctions');
    cy.get('[data-testid="auction-card"]').first().click();
    cy.url().should('include', '/auctions/');
    
    cy.get('[data-testid="bid-input"]').type('1500');
    cy.get('[data-testid="place-bid-btn"]').click();
    
    cy.get('[data-testid="success-toast"]').should('contain', 'Bid placed successfully');
    cy.get('[data-testid="current-bid"]').should('contain', 'RM 1,500');
  });
  
  it('shows error for insufficient bid', () => {
    cy.visit('/auctions/123');
    
    cy.get('[data-testid="bid-input"]').type('100');
    cy.get('[data-testid="place-bid-btn"]').click();
    
    cy.get('[data-testid="error-message"]').should('contain', 'Bid too low');
  });
});
```

## Backend Testing (Go/Bun)

### Unit Tests
```go
// src/application/auction/service_test.go
package auction

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
)

type MockAuctionRepository struct {
  mock.Mock
}

func (m *MockAuctionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Auction, error) {
  args := m.Called(ctx, id)
  return args.Get(0).(*domain.Auction), args.Error(1)
}

func TestAuctionService_PlaceBid(t *testing.T) {
  repo := new(MockAuctionRepository)
  service := NewService(repo, nil, nil)
  
  auctionID := uuid.New()
  bidderID := uuid.New()
  
  auction := &domain.Auction{
    ID: auctionID,
    Status: domain.AuctionStatusActive,
    CurrentBid: 1000,
    MinIncrement: 50,
  }
  
  repo.On("GetByID", mock.Anything, auctionID).Return(auction, nil)
  repo.On("Save", mock.Anything, mock.Anything).Return(nil)
  
  err := service.PlaceBid(context.Background(), PlaceBidDTO{
    AuctionID: auctionID,
    BidderID: bidderID,
    Amount: 1100,
  })
  
  assert.NoError(t, err)
  repo.AssertExpectations(t)
}

func TestAuctionService_PlaceBid_BidTooLow(t *testing.T) {
  repo := new(MockAuctionRepository)
  service := NewService(repo, nil, nil)
  
  auction := &domain.Auction{
    CurrentBid: 1000,
    MinIncrement: 50,
  }
  
  repo.On("GetByID", mock.Anything, mock.Anything).Return(auction, nil)
  
  err := service.PlaceBid(context.Background(), PlaceBidDTO{
    Amount: 1020, // Less than minimum increment
  })
  
  assert.Error(t, err)
  assert.Equal(t, "bid too low", err.Error())
}
```

### Integration Tests
```go
// tests/integration/auction_test.go
package integration

import (
  "testing"
  "net/http"
  "net/http/httptest"
  "github.com/stretchr/testify/assert"
)

func TestAuctionAPI_CreateAuction(t *testing.T) {
  // Setup test database
  db := setupTestDB(t)
  defer teardownTestDB(t, db)
  
  // Create test user and get token
  token := createTestUserAndLogin(t, db, "seller")
  
  // Create request
  payload := map[string]interface{}{
    "product_id": uuid.New().String(),
    "start_price": 1000,
    "end_time": time.Now().Add(1 * time.Hour),
  }
  body, _ := json.Marshal(payload)
  
  req := httptest.NewRequest("POST", "/api/v1/auctions", bytes.NewReader(body))
  req.Header.Set("Authorization", "Bearer "+token)
  req.Header.Set("Content-Type", "application/json")
  
  // Execute
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)
  
  // Assert
  assert.Equal(t, http.StatusCreated, w.Code)
  
  var response map[string]interface{}
  json.Unmarshal(w.Body.Bytes(), &response)
  assert.NotNil(t, response["data"]["auction"]["id"])
}
```

### API Contract Tests
```go
// tests/contract/auction_contract_test.go
package contract

import (
  "testing"
  "github.com/pact-foundation/pact-go/dsl"
)

func TestAuctionContract(t *testing.T) {
  pact := &dsl.Pact{
    Consumer: "frontend",
    Provider: "backend",
  }
  
  pact.AddInteraction().
    Given("auction exists").
    UponReceiving("a request for auction details").
    WithRequest(dsl.Request{
      Method: "GET",
      Path: dsl.String("/api/v1/auctions/123"),
    }).
    WillRespondWith(dsl.Response{
      Status: 200,
      Body: dsl.Like(map[string]interface{}{
        "id": dsl.String("123"),
        "start_price": dsl.NumberLike(1000),
        "current_bid": dsl.NumberLike(1500),
        "status": dsl.String("active"),
      }),
    })
  
  err := pact.Verify(t)
  assert.NoError(t, err)
}
```

## Mobile Testing (React Native)

### Unit Tests
```dart
// test/domain/auction_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:blytz/domain/entities/auction.dart';

void main() {
  group('Auction', () {
    test('can place valid bid', () {
      final auction = Auction(
        id: '1',
        currentBid: 1000,
        minIncrement: 50,
        status: AuctionStatus.active,
      );
      
      final bid = Bid(amount: 1100, bidderId: 'user1');
      final result = auction.placeBid(bid);
      
      expect(result.isSuccess, true);
      expect(auction.currentBid, 1100);
    });
    
    test('rejects bid below minimum increment', () {
      final auction = Auction(
        currentBid: 1000,
        minIncrement: 50,
        status: AuctionStatus.active,
      );
      
      final bid = Bid(amount: 1020, bidderId: 'user1');
      final result = auction.placeBid(bid);
      
      expect(result.isFailure, true);
      expect(result.error, 'Bid too low');
    });
  });
}
```

### Widget Tests
```dart
// test/presentation/widgets/bid_button_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:blytz/presentation/widgets/bid_button.dart';

void main() {
  testWidgets('BidButton displays correctly', (WidgetTester tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home: BidButton(
          currentBid: 1000,
          onPressed: () {},
        ),
      ),
    );
    
    expect(find.text('Place Bid'), findsOneWidget);
    expect(find.text('RM 1,000'), findsOneWidget);
  });
  
  testWidgets('BidButton calls onPressed when tapped', (WidgetTester tester) async {
    var pressed = false;
    
    await tester.pumpWidget(
      MaterialApp(
        home: BidButton(
          currentBid: 1000,
          onPressed: () => pressed = true,
        ),
      ),
    );
    
    await tester.tap(find.text('Place Bid'));
    expect(pressed, true);
  });
}
```

### Integration Tests
```dart
// integration_test/app_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:blytz/main.dart' as app;

void main() {
  IntegrationTestWidgetsReact NativeBinding.ensureInitialized();
  
  group('End-to-End Tests', () {
    test('user can login and view auctions', () async {
      app.main();
      await tester.pumpAndSettle();
      
      // Login
      await tester.enterText(find.byKey(Key('email')), 'test@example.com');
      await tester.enterText(find.byKey(Key('password')), 'password');
      await tester.tap(find.byKey(Key('login_button')));
      await tester.pumpAndSettle();
      
      // Navigate to auctions
      await tester.tap(find.byIcon(Icons.gavel));
      await tester.pumpAndSettle();
      
      // Verify auction list displayed
      expect(find.byType(AuctionCard), findsWidgets);
    });
    
    test('user can place bid', () async {
      app.main();
      await tester.pumpAndSettle();
      
      // Navigate to auction detail
      await tester.tap(find.byType(AuctionCard).first);
      await tester.pumpAndSettle();
      
      // Enter bid amount
      await tester.enterText(find.byKey(Key('bid_input')), '1500');
      await tester.tap(find.byKey(Key('place_bid_button')));
      await tester.pumpAndSettle();
      
      // Verify success
      expect(find.text('Bid placed successfully'), findsOneWidget);
    });
  });
}
```

## Test Data

### Factories
```typescript
// tests/factories/user.ts
export const userFactory = Factory.define<User>(() => ({
  id: faker.string.uuid(),
  email: faker.internet.email(),
  firstName: faker.person.firstName(),
  lastName: faker.person.lastName(),
  role: 'buyer',
  avatarUrl: faker.image.avatar(),
  createdAt: faker.date.past(),
}));

// tests/factories/auction.ts
export const auctionFactory = Factory.define<Auction>(() => ({
  id: faker.string.uuid(),
  product: productFactory.build(),
  seller: sellerFactory.build(),
  startPrice: faker.number.int({ min: 100, max: 1000 }),
  currentBid: faker.number.int({ min: 100, max: 1000 }),
  bidCount: faker.number.int({ min: 0, max: 50 }),
  status: 'active',
  endTime: faker.date.future(),
}));
```

## Test Coverage Goals

| Layer | Target | Minimum |
|-------|--------|---------|
| Backend Unit | 80% | 70% |
| Backend Integration | 60% | 50% |
| Frontend Unit | 70% | 60% |
| Frontend E2E | Critical paths | Critical paths |
| Mobile Unit | 70% | 60% |
| Mobile E2E | Critical paths | Critical paths |

## CI/CD Integration

### GitHub Actions
```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:17
      redis:
        image: redis:8
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test ./... -coverprofile=coverage.out
      - run: go tool cover -func=coverage.out

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: npm ci
      - run: npm run test:unit -- --coverage
      - run: npm run test:e2e

  mobile-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/react-native-action@v2
      - run: react-native pub get
      - run: react-native test --coverage
```

## Testing Checklist

### Before Commit
- [ ] Unit tests pass
- [ ] No console errors
- [ ] TypeScript types correct

### Before PR
- [ ] All tests pass
- [ ] Coverage maintained
- [ ] Integration tests pass
- [ ] Manual smoke test

### Before Release
- [ ] E2E tests pass
- [ ] Performance tests pass
- [ ] Security scan clean
- [ ] Accessibility audit

---

*Last updated: 2025-02-05*
