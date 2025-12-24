import React, { useState } from 'react';
import { MessageSquare } from 'lucide-react';
import { Header } from './components/Header';
import { ChatAssistant } from './components/ChatAssistant';
import { Home } from './components/Home';
import { ProductDetail } from './components/ProductDetail';
import { Checkout } from './components/Checkout';
import { Drops } from './components/Drops';
import { Account } from './components/Account';
import { CartDrawer } from './components/CartDrawer';
import { SellerDashboard } from './components/SellerDashboard';
import { Login } from './components/Login';
import { Register } from './components/Register';
import { useCartStore } from '../store/cartStore';
import { useAuthStore } from '../store/authStore';
import { ViewState } from './types';

interface AppState {
  view: ViewState;
  selectedProduct: Product | null;
  activeCategory: string;
  isChatOpen: boolean;
  isLoginOpen: boolean;
  isRegisterOpen: boolean;
}

const App: React.FC = () => {
  // Use stores
  const cart = useCartStore();
  const auth = useAuthStore();

  // App state
  const [appState, setAppState] = useState<AppState>({
    view: 'HOME',
    selectedProduct: null,
    activeCategory: 'all',
    isChatOpen: false,
    isLoginOpen: false,
    isRegisterOpen: false,
  });

  const handleNavClick = (newView: ViewState) => {
    setAppState(prev => ({ ...prev, view: newView }));
    window.scrollTo(0, 0);
  };

  const handleProductClick = (product: Product) => {
    setAppState(prev => ({ ...prev, selectedProduct: product, view: 'PRODUCT_DETAIL' }));
  };

  const handleAddToCart = (product: Product) => {
    cart.addItem(product, 1);
    cart.setIsCartOpen(true);
  };

  return (
    <div className="min-h-screen bg-blytz-black text-gray-100 font-sans selection:bg-blytz-neon selection:text-black">
      <Header
        cartCount={cart.getItemCount()}
        onCartClick={() => cart.setIsCartOpen(true)}
        onNavClick={handleNavClick}
        onLoginClick={() => setAppState(prev => ({ ...prev, isLoginOpen: true }))}
        onRegisterClick={() => setAppState(prev => ({ ...prev, isRegisterOpen: true }))}
      />

      <main>
        {appState.view === 'HOME' && (
          <Home
            onAddToCart={handleAddToCart}
            onProductClick={handleProductClick}
          />
        )}

        {appState.view === 'PRODUCT_DETAIL' && appState.selectedProduct && (
          <ProductDetail
            product={appState.selectedProduct}
            onBack={() => setAppState(prev => ({ ...prev, selectedProduct: null, view: 'HOME' }))}
            onAddToCart={handleAddToCart}
          />
        )}

        {appState.view === 'CHECKOUT' && (
          <Checkout
            onComplete={() => {
              setAppState(prev => ({ ...prev, view: 'HOME' }));
            }}
          />
        )}

        {appState.view === 'DROPS' && <Drops />}

        {appState.view === 'SELL' && (
          <SellerDashboard
            user={auth.user}
            onLoginClick={() => setAppState(prev => ({ ...prev, isLoginOpen: true }))}
          />
        )}

        {appState.view === 'ACCOUNT' && <Account onViewChange={handleNavClick} />}
      </main>

      {/* Chat FAB */}
      {!appState.isChatOpen && (
        <button
          onClick={() => setAppState(prev => ({ ...prev, isChatOpen: true }))}
          className="fixed bottom-8 right-8 w-14 h-14 bg-blytz-neon text-black rounded-full shadow-[0_0_20px_rgba(190,242,100,0.5)] flex items-center justify-center z-40 hover:scale-110 transition-transform animate-bounce"
        >
          <MessageSquare className="w-6 h-6 fill-current" />
        </button>
      )}

      {/* Chat Assistant */}
      {appState.isChatOpen && (
        <ChatAssistant onClose={() => setAppState(prev => ({ ...prev, isChatOpen: false }))} />
      )}

      {/* Cart Drawer */}
      {cart.isCartOpen && <CartDrawer view={appState.view} onViewChange={handleNavClick} />}

      {/* Login Modal */}
      {appState.isLoginOpen && (
        <Login
          onClose={() => setAppState(prev => ({ ...prev, isLoginOpen: false }))}
          onSwitchToRegister={() => setAppState(prev => ({ ...prev, isLoginOpen: false, isRegisterOpen: true }))}
        />
      )}

      {/* Register Modal */}
      {appState.isRegisterOpen && (
        <Register
          onClose={() => setAppState(prev => ({ ...prev, isRegisterOpen: false }))}
          onSwitchToLogin={() => setAppState(prev => ({ ...prev, isRegisterOpen: false, isLoginOpen: true }))}
        />
      )}

      {/* Footer */}
      <footer className="bg-black border-t border-white/10 py-12 mt-auto">
        <div className="container mx-auto px-4 flex flex-col md:flex-row justify-between items-center gap-6">
          <div className="text-2xl font-display font-bold italic text-white">
            BLYTZ<span className="text-gray-600">.LIVE</span>
          </div>
          <div className="flex gap-8 text-gray-500 text-sm">
            <button onClick={() => setAppState(prev => ({ ...prev, view: 'DROPS' }))} className="hover:text-blytz-neon">Drops</button>
            <button onClick={() => setAppState(prev => ({ ...prev, view: 'SELL' }))} className="hover:text-blytz-neon">Sell</button>
            <a href="#" className="hover:text-blytz-neon">Terms</a>
            <a href="#" className="hover:text-blytz-neon">Privacy</a>
          </div>
          <div className="text-gray-600 text-sm">
            © 2024 Blytz Commerce Protocol.
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;
