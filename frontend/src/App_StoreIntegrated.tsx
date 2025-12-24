import React, { useState, useEffect } from 'react';
import { ShoppingBag, X, Minus, Plus, Trash2, ArrowRight, Zap, ShieldCheck, Truck, RotateCcw, CreditCard, MapPin, Upload, Camera, CheckCircle, Package, User, MessageSquare, Send, Sparkles, Bot, LayoutDashboard, FileSpreadsheet, MoreVertical, Edit, Copy, BarChart3, Search, TrendingUp, Map, Users, Bell, Megaphone, Settings, Globe, DollarSign } from 'lucide-react';
import { Header } from './components/Header';
import { ProductCard } from './components/ProductCard';
import { Button, Badge, Input } from './components/UI';
import { CATEGORIES, PRODUCTS, DROPS } from './constants';
import { Product, ViewState } from './types';
import { useCartStore } from './store/cartStore';
import { useAuthStore } from './store/authStore';
import { GoogleGenAI } from '@google/genai';

// --- AI Chat Assistant Component ---
const ChatAssistant: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  const [messages, setMessages] = useState<{role: 'user' | 'model', text: string}[]>([
    { role: 'model', text: "SYSTEM ONLINE. I am Blytz AI. searching inventory... How can I assist your acquisition today?" }
  ]);
  const [input, setInput] = useState('');
  const [isTyping, setIsTyping] = useState(false);

  useEffect(() => {
    const scrollContainer = document.getElementById('chat-messages');
    if (scrollContainer) {
      scrollContainer.scrollTop = scrollContainer.scrollHeight;
    }
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim()) return;
    const userMsg = input;
    setMessages(prev => [...prev, { role: 'user', text: userMsg }]);
    setInput('');
    setIsTyping(true);

    try {
      const apiKey = import.meta.env.VITE_GEMINI_API_KEY;
      if (!apiKey) {
        setMessages(prev => [...prev, { role: 'model', text: "ERR: AI service not configured." }]);
        setIsTyping(false);
        return;
      }

      const ai = new GoogleGenAI({ apiKey });
      const inventoryContext = PRODUCTS.map(p => `${p.title} ($${p.price}, ${p.category})`).join(', ');

      const response = await ai.models.generateContent({
        model: 'gemini-2.5-flash',
        contents: `User: ${userMsg}`,
        config: {
          systemInstruction: `You are AI interface for Blytz.live, a high-speed cyberpunk marketplace.
          Your tone is concise, robotic but helpful, and energetic.
          Current Inventory: ${inventoryContext}.
          If user asks for recommendations, suggest items from inventory.
          Keep responses under 50 words. Use terminology like "Affirmative", "Scanning", "Uplink established".`,
        }
      });

      setMessages(prev => [...prev, { role: 'model', text: response.text || "Connection interrupted." }]);
    } catch (error) {
      setMessages(prev => [...prev, { role: 'model', text: "ERR: Uplink failed. Try again." }]);
    } finally {
      setIsTyping(false);
    }
  };

  return (
    <div className="fixed bottom-24 right-4 md:right-8 w-80 md:w-96 bg-blytz-dark border border-blytz-neon/30 shadow-[0_0_30px_rgba(0,0,0,0.5)] rounded-lg flex flex-col overflow-hidden animate-in slide-in-from-bottom-10 z-50">
      <div className="bg-blytz-neon/10 border-b border-blytz-neon/20 p-3 flex justify-between items-center backdrop-blur-md">
        <div className="flex items-center gap-2">
          <Bot className="w-5 h-5 text-blytz-neon" />
          <span className="font-display font-bold text-white tracking-wider">BLYTZ AI</span>
        </div>
        <button onClick={onClose} className="text-gray-400 hover:text-white"><X className="w-4 h-4" /></button>
      </div>

      <div id="chat-messages" className="h-80 overflow-y-auto p-4 space-y-3 bg-black/80 backdrop-blur-sm">
        {messages.map((m, i) => (
          <div key={i} className={`flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            <div className={`max-w-[85%] p-3 rounded-lg text-sm ${
              m.role === 'user'
                ? 'bg-white/10 text-white rounded-br-none'
                : 'bg-blytz-neon/10 text-blytz-neon border border-blytz-neon/20 rounded-bl-none font-mono'
            }`}>
              {m.text}
            </div>
          </div>
        ))}
        {isTyping && (
           <div className="flex justify-start">
            <div className="bg-blytz-neon/10 border border-blytz-neon/20 px-3 py-2 rounded-lg rounded-bl-none flex items-center gap-1">
              <span className="w-1.5 h-1.5 bg-blytz-neon rounded-full animate-bounce"></span>
              <span className="w-1.5 h-1.5 bg-blytz-neon rounded-full animate-bounce delay-75"></span>
              <span className="w-1.5 h-1.5 bg-blytz-neon rounded-full animate-bounce delay-150"></span>
            </div>
           </div>
        )}
      </div>

      <div className="p-3 bg-blytz-dark border-t border-white/10 flex gap-2">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSend()}
          placeholder="Query inventory..."
          className="flex-1 bg-black border border-white/10 rounded px-3 py-2 text-sm text-white focus:border-blytz-neon outline-none font-mono"
        />
        <button
          onClick={handleSend}
          className="bg-blytz-neon text-black p-2 rounded hover:bg-white transition-colors"
        >
          <Send className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
};

// --- Main App ---
const App: React.FC = () => {
  // Use stores instead of local state
  const cart = useCartStore();
  const auth = useAuthStore();

  const [view, setView] = useState<ViewState>('HOME');
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [activeCategory, setActiveCategory] = useState<string>('all');
  const [checkoutStep, setCheckoutStep] = useState(1);
  const [isChatOpen, setIsChatOpen] = useState(false);

  // Dashboard State
  const [dashboardTab, setDashboardTab] = useState<'OVERVIEW' | 'INVENTORY' | 'BULK' | 'ORDERS' | 'ANALYTICS' | 'MARKETING' | 'MESSAGES' | 'SETTINGS'>('OVERVIEW');

  // Sell Form State
  const [sellForm, setSellForm] = useState({
    title: '',
    category: '',
    condition: '',
    description: '',
    price: ''
  });

  // Marketing AI State
  const [adPrompt, setAdPrompt] = useState('');
  const [generatedAd, setGeneratedAd] = useState('');
  const [isGeneratingAd, setIsGeneratingAd] = useState(false);

  // --- Nav Handler ---
  const handleNavClick = (newView: ViewState) => {
    setView(newView);
    window.scrollTo(0, 0);
  };

  // --- Cart Handlers (using store) ---
  const handleAddToCart = (product: Product) => {
    cart.addItem(product, 1);
    cart.setIsCartOpen(true);
  };

  // --- Render Views ---
  const renderHome = () => (
    <>
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-blytz-black border-b border-white/10">
        <div className="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/carbon-fibre.png')] opacity-20"></div>
        <div className="absolute -right-20 -top-20 w-96 h-96 bg-blytz-neon/20 blur-[100px] rounded-full"></div>

        <div className="container mx-auto px-4 py-20 relative z-10">
          <div className="max-w-3xl">
            <Badge variant="flash">System Update: Prices Dropped</Badge>
            <h1 className="text-5xl md:text-7xl font-display font-bold text-white mt-6 mb-6 leading-[0.9] italic">
              SPEED IS THE <br/>
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-blytz-neon to-blytz-lime">
                NEW CURRENCY
              </span>
            </h1>
            <p className="text-xl text-gray-400 mb-8 max-w-xl">
              The next generation marketplace for instant gratification.
              Verified sellers. Millisecond transactions. Instant dispatch.
            </p>
            <div className="flex flex-wrap gap-4">
              <Button size="lg" onClick={() => {
                document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' });
              }}>
                Shop The Drop <ArrowRight className="ml-2 w-5 h-5" />
              </Button>
              <Button variant="secondary" size="lg">Sell Your Gear</Button>
            </div>
          </div>
        </div>
      </section>

      {/* Categories */}
      <section className="bg-blytz-black py-8 border-b border-white/10 sticky top-16 z-40 backdrop-blur-md bg-opacity-90">
        <div className="container mx-auto px-4 overflow-x-auto no-scrollbar">
          <div className="flex gap-4 min-w-max">
            <button
              onClick={() => setActiveCategory('all')}
              className={`px-6 py-2 rounded-full border transition-all ${
                activeCategory === 'all'
                  ? 'bg-white text-black border-white font-bold'
                  : 'bg-transparent text-gray-400 border-gray-800 hover:border-white hover:text-white'
              }`}
            >
              All Drops
            </button>
            {CATEGORIES.map(cat => (
              <button
                key={cat.id}
                onClick={() => setActiveCategory(cat.name)}
                className={`px-6 py-2 rounded-full border flex items-center gap-2 transition-all ${
                  activeCategory === cat.name
                    ? 'bg-blytz-neon text-black border-blytz-neon font-bold'
                    : 'bg-transparent text-gray-400 border-gray-800 hover:border-blytz-neon hover:text-blytz-neon'
                }`}
              >
                {cat.icon}
                {cat.name}
              </button>
            ))}
          </div>
        </div>
      </section>

      {/* Product Grid */}
      <section id="products" className="py-12 bg-blytz-black">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {PRODUCTS
              .filter(p => activeCategory === 'all' || p.category === activeCategory)
              .map(product => (
                <ProductCard
                  key={product.id}
                  product={product}
                  onAdd={handleAddToCart}
                  onClick={(p) => {
                    setSelectedProduct(p);
                    setView('PRODUCT_DETAIL');
                    window.scrollTo(0, 0);
                  }}
                />
            ))}
          </div>

          {PRODUCTS.filter(p => activeCategory === 'all' || p.category === activeCategory).length === 0 && (
             <div className="text-center py-20">
               <h3 className="text-2xl font-bold text-gray-600">No signals found in this sector.</h3>
               <Button variant="ghost" onClick={() => setActiveCategory('all')} className="mt-4">Reset Signal</Button>
             </div>
           )}
        </div>
      </section>

      {/* Features / Trust */}
      <section className="py-20 border-t border-white/10 bg-blytz-dark">
         <div className="container mx-auto px-4 grid grid-cols-1 md:grid-cols-3 gap-8">
           <div className="p-6 border border-white/5 bg-blytz-black hover:border-blytz-neon/30 transition-colors group">
             <Zap className="w-10 h-10 text-blytz-neon mb-4 group-hover:scale-110 transition-transform" />
             <h3 className="text-xl font-bold text-white mb-2">Instant Authentication</h3>
             <p className="text-gray-400">Every item is digitally verified in real-time before it leaves the seller.</p>
           </div>
           <div className="p-6 border border-white/5 bg-blytz-black hover:border-blytz-neon/30 transition-colors group">
             <Truck className="w-10 h-10 text-blytz-neon mb-4 group-hover:scale-110 transition-transform" />
             <h3 className="text-xl font-bold text-white mb-2">Hyper-Local Logistics</h3>
             <p className="text-gray-400">Our decentralized warehouse network ensures same-day delivery in metro areas.</p>
           </div>
           <div className="p-6 border border-white/5 bg-blytz-black hover:border-blytz-neon/30 transition-colors group">
             <RotateCcw className="w-10 h-10 text-blytz-neon mb-4 group-hover:scale-110 transition-transform" />
             <h3 className="text-xl font-bold text-white mb-2">Zero-Friction Returns</h3>
             <p className="text-gray-400">Don't like it? Scan the QR code and a drone picks it up. Instant refund.</p>
           </div>
         </div>
      </section>
    </>
  );

  return (
    <div className="min-h-screen bg-blytz-black text-gray-100 font-sans selection:bg-blytz-neon selection:text-black">
      <Header
        cartCount={cart.getItemCount()}
        onCartClick={() => cart.setIsCartOpen(true)}
        onNavClick={handleNavClick}
      />

      <main>
        {view === 'HOME' && renderHome()}
        {view === 'PRODUCT_DETAIL' && selectedProduct && (
          <div className="container mx-auto px-4 py-8 animate-in slide-in-from-right-8 duration-300">
            <button
              onClick={() => setView('HOME')}
              className="mb-6 text-gray-400 hover:text-white flex items-center gap-2"
            >
              ← Back to Marketplace
            </button>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
              <div className="space-y-4">
                <div className="aspect-square bg-blytz-dark rounded-lg overflow-hidden border border-white/10 relative group">
                  <img
                    src={selectedProduct.image}
                    alt={selectedProduct.title}
                    className="w-full h-full object-cover"
                  />
                 {selectedProduct.isFlash && (
                    <div className="absolute top-4 left-4">
                      <Badge variant="flash">Flash Deal Ends {selectedProduct.timeLeft}</Badge>
                    </div>
                 )}
                </div>
                <div className="grid grid-cols-4 gap-4">
                  {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="aspect-square bg-blytz-dark rounded border border-white/10 cursor-pointer hover:border-blytz-neon">
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex flex-col h-full">
                <div className="mb-auto">
                  <h1 className="text-4xl md:text-5xl font-display font-bold text-white mb-2 italic">
                    {selectedProduct.title}
                  </h1>
                  <div className="flex items-center gap-4 mb-6">
                    <span className="text-3xl font-bold text-blytz-neon">
                      ${selectedProduct.price.toFixed(2)}
                    </span>
                    {selectedProduct.originalPrice && (
                      <span className="text-xl text-gray-500 line-through">
                        ${selectedProduct.originalPrice.toFixed(2)}
                      </span>
                    )}
                    <div className="flex items-center gap-1 text-yellow-400 ml-4">
                      <span className="font-bold">{selectedProduct.rating}</span>
                      <div className="flex">
                        {[...Array(5)].map((_, i) => (
                          <Zap key={i} className={`w-4 h-4 ${i < Math.floor(selectedProduct.rating) ? 'fill-current' : 'text-gray-700'}`} />
                        ))}
                      </div>
                      <span className="text-gray-400 text-sm ml-2">({selectedProduct.reviews} verified)</span>
                    </div>
                  </div>

                  <p className="text-gray-300 text-lg leading-relaxed mb-8">
                    {selectedProduct.description}
                  </p>

                  <div className="space-y-6 mb-8">
                    <div className="p-4 bg-white/5 border border-white/10 rounded">
                      <h3 className="text-white font-bold mb-2 flex items-center gap-2">
                        <Truck className="w-5 h-5 text-blytz-neon" />
                        Blytz Speed Delivery
                      </h3>
                      <p className="text-sm text-gray-400">Order in the next 2 hrs for delivery by tomorrow, 10 AM.</p>
                    </div>
                  </div>
                </div>

                <div className="flex gap-4 mt-8 pt-8 border-t border-white/10">
                  <Button
                    variant="primary"
                    size="lg"
                    className="flex-1 text-lg"
                    onClick={() => handleAddToCart(selectedProduct)}
                  >
                    Add To Cart
                  </Button>
                  <Button variant="outline" size="lg" className="px-4">
                    <ShieldCheck className="w-6 h-6" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}
        {view === 'CHECKOUT' && (
          <div className="container mx-auto px-4 py-12 animate-in slide-in-from-bottom-8">
            <h1 className="text-4xl font-display font-bold text-white mb-8 italic">SECURE CHECKOUT</h1>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
              <div className="lg:col-span-2 space-y-8">
                <div className={`border ${checkoutStep >= 1 ? 'border-blytz-neon/50 bg-blytz-dark/50' : 'border-white/10 bg-transparent'} p-6 rounded transition-all`}>
                  <div className="flex items-center gap-4 mb-6">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold ${checkoutStep >= 1 ? 'bg-blytz-neon text-black' : 'bg-gray-800 text-gray-500'}`}>1</div>
                    <h2 className="text-xl font-bold text-white">SHIPPING DATA</h2>
                  </div>

                  {checkoutStep === 1 && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <Input placeholder="First Name" />
                      <Input placeholder="Last Name" />
                      <Input placeholder="Address Line 1" className="md:col-span-2" />
                      <Input placeholder="City" />
                      <Input placeholder="Zip Code" />
                      <div className="md:col-span-2 mt-4">
                         <Button onClick={() => setCheckoutStep(2)} className="w-full">Proceed to Payment</Button>
                      </div>
                    </div>
                  )}
                  {checkoutStep > 1 && <div className="text-green-400 flex items-center gap-2"><CheckCircle className="w-4 h-4"/> Data Secured</div>}
                </div>

                <div className={`border ${checkoutStep >= 2 ? 'border-blytz-neon/50 bg-blytz-dark/50' : 'border-white/10 bg-transparent'} p-6 rounded transition-all`}>
                  <div className="flex items-center gap-4 mb-6">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold ${checkoutStep >= 2 ? 'bg-blytz-neon text-black' : 'bg-gray-800 text-gray-500'}`}>2</div>
                    <h2 className="text-xl font-bold text-white">PAYMENT UPLINK</h2>
                  </div>

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
                       <Input placeholder="Card Number" />
                       <div className="grid grid-cols-2 gap-4">
                          <Input placeholder="MM/YY" />
                          <Input placeholder="CVC" />
                       </div>
                       <Button onClick={() => setCheckoutStep(3)} className="w-full mt-4">Establish Uplink (Pay ${cart.getTotal().toFixed(2)})</Button>
                    </div>
                  )}
                  {checkoutStep > 2 && <div className="text-green-400 flex items-center gap-2"><CheckCircle className="w-4 h-4"/> Transaction Verified</div>}
                </div>

                 {checkoutStep === 3 && (
                   <div className="bg-blytz-neon/10 border border-blytz-neon p-8 rounded text-center animate-pulse-fast">
                     <Zap className="w-16 h-16 text-blytz-neon mx-auto mb-4" />
                     <h2 className="text-3xl font-display font-bold text-white mb-2">ORDER CONFIRMED</h2>
                     <p className="text-gray-400 mb-6">Dispatch drones are spooling up. Estimated arrival: <span className="text-blytz-neon font-mono">T-minus 2 hours</span>.</p>
                     <Button onClick={() => { cart.clearCart(); setView('HOME'); setCheckoutStep(1); }}>Return to Base</Button>
                   </div>
                 )}

              </div>

              <div className="bg-blytz-dark border border-white/10 p-6 h-fit rounded sticky top-24">
                <h3 className="text-lg font-bold text-white mb-4 border-b border-white/10 pb-2">ORDER MANIFEST</h3>
                <div className="space-y-4 mb-6">
                  {cart.items.map(item => (
                    <div key={item.id} className="flex justify-between items-center text-sm">
                      <div className="flex items-center gap-2">
                         <div className="w-5 h-5 bg-blytz-neon text-black font-bold flex items-center justify-center rounded-sm text-xs">{item.quantity}</div>
                         <span className="text-gray-300">{item.title}</span>
                      </div>
                      <span className="text-white font-mono">${(item.price * item.quantity).toFixed(2)}</span>
                    </div>
                  ))}
                </div>
                <div className="border-t border-white/10 pt-4 space-y-2">
                  <div className="flex justify-between text-gray-400"><span>Subtotal</span><span>${cart.getTotal().toFixed(2)}</span></div>
                  <div className="flex justify-between text-blytz-neon"><span>Delivery</span><span>0.00</span></div>
                  <div className="flex justify-between text-xl font-bold text-white mt-4"><span>Total</span><span>${cart.getTotal().toFixed(2)}</span></div>
                </div>
              </div>
            </div>
          </div>
        )}
        {view === 'DROPS' && (
          <div className="container mx-auto px-4 py-12">
            <div className="text-center mb-16">
              <Badge variant="flash">INCOMING TRANSMISSIONS</Badge>
              <h1 className="text-5xl md:text-7xl font-display font-bold text-white mt-4 mb-4 italic">FUTURE DROPS</h1>
              <p className="text-gray-400 max-w-2xl mx-auto">Get notified before the masses. These limited run items will sell out in seconds.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              {DROPS.map(product => (
                <div key={product.id} className="group relative border border-white/10 bg-blytz-dark overflow-hidden hover:border-blytz-neon transition-all">
                  <div className="aspect-[4/5] relative">
                    <img src={product.image} alt={product.title} className="w-full h-full object-cover opacity-60 group-hover:opacity-80 transition-opacity grayscale group-hover:grayscale-0" />
                    <div className="absolute inset-0 flex items-center justify-center">
                       <div className="bg-black/80 backdrop-blur border border-white/20 px-6 py-3 transform rotate-[-5deg]">
                          <span className="font-mono text-2xl font-bold text-blytz-neon block text-center">{product.dropDate}</span>
                          <span className="text-xs text-gray-400 uppercase tracking-widest block text-center">Launch Date</span>
                       </div>
                    </div>
                  </div>
                  <div className="p-6">
                    <h3 className="text-2xl font-bold text-white mb-2">{product.title}</h3>
                    <p className="text-gray-400 mb-6 line-clamp-2">{product.description}</p>
                    <div className="flex items-center justify-between">
                      <span className="text-xl font-mono text-gray-300">${product.price}</span>
                      <Button variant="outline" size="sm">Notify Me</Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
        {view === 'SELL' && (
          <div className="container mx-auto px-4 py-8">
            <div className="text-center mb-12">
              <h1 className="text-4xl font-display font-bold text-white mb-4">SELLER DASHBOARD</h1>
              {auth.user ? (
                <p className="text-gray-400">Welcome back, {auth.user.first_name}</p>
              ) : (
                <Button onClick={() => alert('Login functionality coming soon!')}>Login to Start Selling</Button>
              )}
            </div>

            <div className="max-w-2xl mx-auto">
              <div className="bg-blytz-dark border border-white/10 p-8 rounded-lg">
                <h2 className="text-2xl font-bold text-white mb-6">List New Product</h2>
                <div className="space-y-4">
                  <Input
                    label="Product Title"
                    value={sellForm.title}
                    onChange={(e) => setSellForm({...sellForm, title: e.target.value})}
                    placeholder="Enter product title"
                  />
                  <select
                    className="w-full bg-black border border-white/10 rounded px-4 py-3 text-white focus:border-blytz-neon outline-none"
                    value={sellForm.category}
                    onChange={(e) => setSellForm({...sellForm, category: e.target.value})}
                  >
                    <option value="">Select Category</option>
                    {CATEGORIES.map(cat => (
                      <option key={cat.id} value={cat.name}>{cat.name}</option>
                    ))}
                  </select>
                  <select
                    className="w-full bg-black border border-white/10 rounded px-4 py-3 text-white focus:border-blytz-neon outline-none"
                    value={sellForm.condition}
                    onChange={(e) => setSellForm({...sellForm, condition: e.target.value})}
                  >
                    <option value="">Select Condition</option>
                    <option value="New">New</option>
                    <option value="Like New">Like New</option>
                    <option value="Good">Good</option>
                    <option value="Fair">Fair</option>
                  </select>
                  <textarea
                    className="w-full bg-black border border-white/10 rounded p-4 text-white focus:border-blytz-neon outline-none h-32"
                    placeholder="Product description"
                    value={sellForm.description}
                    onChange={(e) => setSellForm({...sellForm, description: e.target.value})}
                  />
                  <Input
                    label="Price (USD)"
                    type="number"
                    value={sellForm.price}
                    onChange={(e) => setSellForm({...sellForm, price: e.target.value})}
                    placeholder="0.00"
                  />
                  <div className="flex gap-4 mt-4">
                    <Button variant="outline" onClick={() => setView('HOME')}>Cancel</Button>
                    <Button onClick={() => alert('Product listing coming soon!')}>List Product</Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
        {view === 'ACCOUNT' && (
          <div className="container mx-auto px-4 py-12">
            <h1 className="text-3xl font-display font-bold text-white mb-8">MY ACCOUNT</h1>
            {auth.user ? (
              <div className="max-w-2xl">
                <div className="bg-blytz-dark border border-white/10 p-6 rounded-lg mb-6">
                  <h2 className="text-xl font-bold text-white mb-4">Profile Information</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <Input label="First Name" defaultValue={auth.user.first_name} />
                    <Input label="Last Name" defaultValue={auth.user.last_name} />
                    <div className="md:col-span-2">
                      <Input label="Email" defaultValue={auth.user.email} />
                    </div>
                    <div className="md:col-span-2">
                      <Input label="Phone" defaultValue={auth.user.phone || ''} />
                    </div>
                  </div>
                  <div className="flex gap-4 mt-6">
                    <Button>Save Changes</Button>
                    <Button variant="outline" onClick={() => auth.logout()}>Logout</Button>
                  </div>
                </div>

                <div className="bg-blytz-dark border border-white/10 p-6 rounded-lg">
                  <h2 className="text-xl font-bold text-white mb-4">Account Type</h2>
                  <div className="p-4 bg-blytz-neon/10 rounded mb-4">
                    <div className="flex items-center gap-3">
                      <ShieldCheck className="w-8 h-8 text-blytz-neon" />
                      <div>
                        <p className="font-bold text-white text-lg">{auth.user.role.toUpperCase()}</p>
                        <p className="text-gray-400 text-sm">{auth.user.email_verified ? 'Verified Account' : 'Verification Pending'}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="text-center py-20 bg-blytz-dark border border-white/10 rounded-lg">
                <User className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                <h3 className="text-xl font-bold text-white mb-2">Please login to access your account</h3>
                <Button onClick={() => setView('HOME')}>Return to Home</Button>
              </div>
            )}
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-black border-t border-white/10 py-12 mt-auto">
        <div className="container mx-auto px-4 flex flex-col md:flex-row justify-between items-center gap-6">
           <div className="text-2xl font-display font-bold italic text-white">
             BLYTZ<span className="text-gray-600">.LIVE</span>
           </div>
           <div className="flex gap-8 text-gray-500 text-sm">
             <button onClick={() => setView('DROPS')} className="hover:text-blytz-neon">Drops</button>
             <button onClick={() => setView('SELL')} className="hover:text-blytz-neon">Sell</button>
             <a href="#" className="hover:text-blytz-neon">Terms</a>
             <a href="#" className="hover:text-blytz-neon">Privacy</a>
           </div>
           <div className="text-gray-600 text-sm">
             © 2024 Blytz Commerce Protocol.
           </div>
        </div>
      </footer>

      {/* Chat Assistant */}
      {isChatOpen && <ChatAssistant onClose={() => setIsChatOpen(false)} />}

      {/* Chat FAB */}
      {!isChatOpen && (
        <button
          onClick={() => setIsChatOpen(true)}
          className="fixed bottom-8 right-8 w-14 h-14 bg-blytz-neon text-black rounded-full shadow-[0_0_20px_rgba(190,242,100,0.5)] flex items-center justify-center z-40 hover:scale-110 transition-transform animate-bounce"
        >
          <MessageSquare className="w-6 h-6 fill-current" />
        </button>
      )}

      {/* Cart Drawer */}
      {cart.isCartOpen && (
        <>
          <div
            className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 transition-opacity"
            onClick={() => cart.setIsCartOpen(false)}
          />
          <div className="fixed inset-y-0 right-0 w-full md:w-[450px] bg-blytz-dark border-l border-white/10 z-[60] shadow-2xl flex flex-col animate-in slide-in-from-right duration-300">
            <div className="p-6 border-b border-white/10 flex items-center justify-between bg-blytz-black">
              <h2 className="text-2xl font-display font-bold italic text-white">YOUR HAUL</h2>
              <button onClick={() => cart.setIsCartOpen(false)} className="text-gray-400 hover:text-white">
                <X />
              </button>
            </div>

            <div className="flex-1 overflow-y-auto p-6 space-y-6">
              {cart.isLoading ? (
                <div className="h-full flex items-center justify-center text-gray-500">
                  <div className="text-center">
                    <div className="w-8 h-8 border-2 border-blytz-neon border-t-transparent rounded-full animate-spin mb-4 mx-auto"></div>
                    <p>Loading your haul...</p>
                  </div>
                </div>
              ) : cart.items.length === 0 ? (
                <div className="h-full flex flex-col items-center justify-center text-gray-500">
                  <ShoppingBag className="w-16 h-16 mb-4 opacity-20" />
                  <p className="text-lg">Your cart is empty.</p>
                  <Button
                    variant="outline"
                    className="mt-4"
                    onClick={() => {
                      cart.setIsCartOpen(false);
                      setView('HOME');
                    }}
                  >
                    Start Shopping
                  </Button>
                </div>
              ) : (
                cart.items.map(item => (
                  <div key={item.id} className="flex gap-4 bg-white/5 p-4 rounded border border-white/5">
                    <img src={item.image} alt={item.title} className="w-20 h-20 object-cover rounded bg-black" />
                    <div className="flex-1">
                      <div className="flex justify-between items-start mb-1">
                        <h4 className="font-bold text-white line-clamp-1">{item.title}</h4>
                        <button onClick={() => cart.removeItem(item.id)} className="text-gray-500 hover:text-red-500">
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                      <p className="text-blytz-neon font-mono mb-3">${item.price.toFixed(2)}</p>

                      <div className="flex items-center gap-3">
                        <button
                          className="w-6 h-6 rounded bg-black border border-white/20 flex items-center justify-center hover:border-blytz-neon"
                          onClick={() => cart.updateQuantity(item.id, -1)}
                        >
                          <Minus className="w-3 h-3" />
                        </button>
                        <span className="text-sm font-bold w-4 text-center">{item.quantity}</span>
                        <button
                          className="w-6 h-6 rounded bg-black border border-white/20 flex items-center justify-center hover:border-blytz-neon"
                          onClick={() => cart.updateQuantity(item.id, 1)}
                        >
                          <Plus className="w-3 h-3" />
                        </button>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>

            {cart.items.length > 0 && (
              <div className="p-6 bg-blytz-black border-t border-white/10">
                <div className="flex justify-between items-center mb-2 text-gray-400">
                  <span>Subtotal</span>
                  <span>${cart.getTotal().toFixed(2)}</span>
                </div>
                <div className="flex justify-between items-center mb-6 text-gray-400">
                  <span>Shipping</span>
                  <span className="text-blytz-neon">FREE (Blytz Prime)</span>
                </div>
                <div className="flex justify-between items-center mb-6 text-xl font-bold text-white">
                  <span>Total</span>
                  <span>${cart.getTotal().toFixed(2)}</span>
                </div>
                <Button
                  className="w-full h-14 text-lg"
                  onClick={() => {
                    cart.setIsCartOpen(false);
                    setView('CHECKOUT');
                    setCheckoutStep(1);
                  }}
                >
                  SECURE CHECKOUT
                </Button>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
};

export default App;
