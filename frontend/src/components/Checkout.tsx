import React, { useState } from 'react';
import { CheckCircle, CreditCard, Zap, ShoppingBag } from 'lucide-react';
import { Button, Badge, Input } from './UI';
import { useCartStore } from '../../store/cartStore';
import { useAuthStore } from '../../store/authStore';

interface CheckoutProps {
  onComplete: () => void;
  onBack: () => void;
}

export const Checkout: React.FC<CheckoutProps> = ({ onComplete, onBack }) => {
  const cart = useCartStore();
  const auth = useAuthStore();

  const [checkoutStep, setCheckoutStep] = useState(1);
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    address: '',
    city: '',
    zipCode: '',
    cardNumber: '',
    cardExpiry: '',
    cvv: ''
  });

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleStepChange = (step: number) => {
    setCheckoutStep(step);
  };

  const handleSubmit = async () => {
    if (!auth.user) {
      alert('Please login to complete your purchase');
      return;
    }

    if (checkoutStep === 3) {
      try {
        await cart.clearCart();
        onComplete();
      } catch (error) {
        console.error('Checkout failed:', error);
      }
    } else {
      handleStepChange(checkoutStep + 1);
    }
  };

  return (
    <div className="container mx-auto px-4 py-12 animate-in slide-in-from-bottom-8">
      <h1 className="text-4xl font-display font-bold text-white mb-8 italic">SECURE CHECKOUT</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
        <div className="lg:col-span-2 space-y-8">
          {checkoutStep === 1 && (
            <div className="border border-blytz-neon/50 bg-blytz-dark/50 p-6 rounded transition-all">
              <div className="flex items-center gap-4 mb-6">
                <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold ${checkoutStep >= 1 ? 'bg-blytz-neon text-black' : 'bg-gray-800 text-gray-500'}`}>1</div>
                <h2 className="text-xl font-bold text-white">SHIPPING DATA</h2>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Input
                  label="First Name"
                  value={formData.firstName}
                  onChange={(e) => handleInputChange('firstName', e.target.value)}
                  placeholder="John"
                />
                <Input
                  label="Last Name"
                  value={formData.lastName}
                  onChange={(e) => handleInputChange('lastName', e.target.value)}
                  placeholder="Doe"
                />
                <Input
                  label="Address Line 1"
                  className="md:col-span-2"
                  value={formData.address}
                  onChange={(e) => handleInputChange('address', e.target.value)}
                  placeholder="123 Market Street"
                />
                <Input
                  label="City"
                  value={formData.city}
                  onChange={(e) => handleInputChange('city', e.target.value)}
                  placeholder="San Francisco"
                />
                <Input
                  label="Zip Code"
                  value={formData.zipCode}
                  onChange={(e) => handleInputChange('zipCode', e.target.value)}
                  placeholder="94102"
                />
                <div className="md:col-span-2 mt-4">
                  <Button onClick={() => handleStepChange(2)} className="w-full">Proceed to Payment</Button>
                </div>
              </div>
            </div>
          )}

          {checkoutStep === 2 && (
            <div className={`border ${checkoutStep >= 2 ? 'border-blytz-neon/50 bg-blytz-dark/50' : 'border-white/10 bg-transparent'} p-6 rounded transition-all`}>
              <div className="flex items-center gap-4 mb-6">
                <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold ${checkoutStep >= 2 ? 'bg-blytz-neon text-black' : 'bg-gray-800 text-gray-500'}`}>2</div>
                <h2 className="text-xl font-bold text-white">PAYMENT UPLINK</h2>
              </div>

              {checkoutStep > 2 && <div className="text-green-400 flex items-center gap-2"><CheckCircle className="w-4 h-4" /> Data Secured</div>}

              {checkoutStep === 2 && (
                <div className="space-y-4">
                  <div className="grid grid-cols-3 gap-4 mb-4">
                    <button className="border border-blytz-neon bg-blytz-neon/10 text-blytz-neon p-4 rounded flex flex-col items-center gap-2 hover:bg-blytz-neon/20">
                      <CreditCard />
                      <span className="text-xs font-bold">Credit Card</span>
                    </button>
                    <button className="border border-white/10 text-gray-400 p-4 rounded flex flex-col items-center gap-2 hover:border-white hover:text-white">
                      <Zap />
                      <span className="text-xs font-bold">Crypto</span>
                    </button>
                    <button className="border border-white/10 text-gray-400 p-4 rounded flex flex-col items-center gap-2 hover:border-white hover:text-white">
                      <ShoppingBag />
                      <span className="text-xs font-bold">Pay Later</span>
                    </button>
                  </div>

                  <Input
                    label="Card Number"
                    value={formData.cardNumber}
                    onChange={(e) => handleInputChange('cardNumber', e.target.value)}
                    placeholder="4242 4242 4242 4242"
                  />
                  <div className="grid grid-cols-2 gap-4">
                    <Input
                      label="MM/YY"
                      value={formData.cardExpiry}
                      onChange={(e) => handleInputChange('cardExpiry', e.target.value)}
                      placeholder="12/28"
                    />
                    <Input
                      label="CVC"
                      value={formData.cvv}
                      onChange={(e) => handleInputChange('cvv', e.target.value)}
                      placeholder="123"
                    />
                  </div>

                  <Button onClick={() => handleStepChange(3)} className="w-full mt-4">
                    Establish Uplink (Pay ${cart.getTotal().toFixed(2)})
                  </Button>
                </div>
              )}
            </div>
          )}

          {checkoutStep === 3 && (
            <div className="bg-blytz-neon/10 border border-blytz-neon p-8 rounded text-center animate-pulse-fast">
              <Zap className="w-16 h-16 text-blytz-neon mx-auto mb-4" />
              <h2 className="text-3xl font-display font-bold text-white mb-2">ORDER CONFIRMED</h2>
              <p className="text-gray-400 mb-6">
                {auth.user ? `${auth.user.first_name}, your ` : 'Your'} haul has been secured!
                <br />
                Dispatch drones are spooling up. Estimated arrival:{' '}
                <span className="text-blytz-neon font-mono">T-minus 2 hours</span>.
              </p>
              <Button onClick={handleSubmit} className="mt-4">
                {auth.user ? 'Return to Home' : 'Login to Continue'}
              </Button>
            </div>
          )}
        </div>

        <div className="bg-blytz-dark border border-white/10 p-6 h-fit rounded sticky top-24">
          <h3 className="text-lg font-bold text-white mb-4 border-b border-white/10 pb-2">ORDER MANIFEST</h3>
          <div className="space-y-4 mb-6">
            {cart.isLoading ? (
              <div className="text-center py-20">
                <div className="inline-block w-8 h-8 border-2 border-blytz-neon border-t-transparent rounded-full animate-spin mb-4 mx-auto"></div>
                <p>Loading your haul...</p>
              </div>
            ) : cart.items.length === 0 ? (
              <div className="h-full flex flex-col items-center justify-center text-gray-500 py-20">
                <ShoppingBag className="w-16 h-16 mb-4 opacity-20" />
                <p className="text-lg">Your cart is empty.</p>
                <Button
                  variant="outline"
                  onClick={onBack}
                  className="mt-4"
                >
                  Start Shopping
                </Button>
              </div>
            ) : (
              cart.items.map(item => (
                <div key={item.id} className="flex justify-between items-center text-sm">
                  <div className="flex items-center gap-2">
                    <div className="w-5 h-5 bg-blytz-neon text-black font-bold flex items-center justify-center rounded-sm text-xs">
                      {item.quantity}
                    </div>
                    <div>
                      <span className="text-gray-300">{item.title}</span>
                    </div>
                  </div>
                  <span className="text-white font-mono">${(item.price * item.quantity).toFixed(2)}</span>
                </div>
              ))
            )}
          </div>

          {cart.items.length > 0 && (
            <div className="border-t border-white/10 pt-4 space-y-2">
              <div className="flex justify-between items-center mb-2 text-gray-400">
                <span>Subtotal</span>
                <span>${cart.getTotal().toFixed(2)}</span>
              </div>
              <div className="flex justify-between items-center mb-2 text-gray-400">
                <span>Shipping</span>
                <span className="text-blytz-neon">FREE (Blytz Prime)</span>
              </div>
              <div className="flex justify-between items-center mb-6 text-xl font-bold text-white mt-4">
                <span>Total</span>
                <span>${cart.getTotal().toFixed(2)}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
