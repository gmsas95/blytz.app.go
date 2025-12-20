import { useState } from 'react'
import { Brain, Zap, Shield, Cpu, Sparkles } from 'lucide-react'

function App() {
  const [isConnected, setIsConnected] = useState(false)

  return (
    <div className="min-h-screen bg-blytz-black text-white">
      {/* Header */}
      <header className="bg-blytz-dark border-b border-blytz-neon/20">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Zap className="h-8 w-8 text-blytz-neon" />
              <h1 className="text-2xl font-bold font-display text-blytz-neon">
                Blytz.live
              </h1>
            </div>
            <nav className="hidden md:flex space-x-6">
              <a href="#" className="text-blytz-lime hover:text-blytz-neon transition-colors">
                Marketplace
              </a>
              <a href="#" className="text-blytz-lime hover:text-blytz-neon transition-colors">
                AI Assistant
              </a>
              <a href="#" className="text-blytz-lime hover:text-blytz-neon transition-colors">
                Profile
              </a>
            </nav>
            <button
              onClick={() => setIsConnected(!isConnected)}
              className={`px-4 py-2 rounded-lg font-medium transition-all ${
                isConnected
                  ? 'bg-blytz-neon text-blytz-black'
                  : 'bg-blytz-accent text-white hover:bg-purple-600'
              }`}
            >
              {isConnected ? 'Connected' : 'Connect Wallet'}
            </button>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <main className="relative overflow-hidden">
        {/* Background Effects */}
        <div className="absolute inset-0 bg-gradient-to-br from-blytz-black via-blytz-dark to-blytz-black" />
        <div className="absolute inset-0 bg-blytz-neon/5 opacity-20" />
        
        {/* Grid Pattern */}
        <div className="absolute inset-0" 
             style={{
               backgroundImage: `linear-gradient(rgba(190, 242, 100, 0.1) 1px, transparent 1px),
                              linear-gradient(90deg, rgba(190, 242, 100, 0.1) 1px, transparent 1px)`,
               backgroundSize: '50px 50px'
             }} />

        <div className="relative container mx-auto px-4 py-24">
          <div className="max-w-4xl mx-auto text-center">
            <div className="flex justify-center mb-6">
              <div className="p-4 bg-blytz-neon/10 rounded-full">
                <Cpu className="h-16 w-16 text-blytz-neon" />
              </div>
            </div>
            
            <h2 className="text-5xl md:text-7xl font-bold font-display mb-6 bg-gradient-to-r from-blytz-neon to-blytz-accent bg-clip-text text-transparent">
              Speed Commerce
            </h2>
            
            <p className="text-xl md:text-2xl text-blytz-lime mb-8 max-w-2xl mx-auto">
              Experience the future of marketplace with AI-powered trading, 
              cyberpunk aesthetics, and lightning-fast transactions.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
              <button className="px-8 py-4 bg-blytz-neon text-blytz-black rounded-lg font-bold text-lg hover:bg-blytz-lime transition-all transform hover:scale-105">
                <Sparkles className="inline h-5 w-5 mr-2" />
                Launch Marketplace
              </button>
              <button className="px-8 py-4 bg-blytz-dark text-blytz-neon rounded-lg font-bold text-lg border border-blytz-neon/30 hover:border-blytz-neon transition-all">
                <Brain className="inline h-5 w-5 mr-2" />
                Try AI Assistant
              </button>
            </div>

            {/* Feature Cards */}
            <div className="grid md:grid-cols-3 gap-6">
              <div className="bg-blytz-dark border border-blytz-neon/20 rounded-lg p-6 hover:border-blytz-neon/50 transition-all">
                <Zap className="h-12 w-12 text-blytz-neon mb-4 mx-auto" />
                <h3 className="text-xl font-bold mb-2 text-blytz-neon">
                  Lightning Fast
                </h3>
                <p className="text-blytz-lime">
                  Execute trades at the speed of light with optimized infrastructure.
                </p>
              </div>

              <div className="bg-blytz-dark border border-blytz-accent/20 rounded-lg p-6 hover:border-blytz-accent/50 transition-all">
                <Brain className="h-12 w-12 text-blytz-accent mb-4 mx-auto" />
                <h3 className="text-xl font-bold mb-2 text-blytz-accent">
                  AI-Powered
                </h3>
                <p className="text-blytz-lime">
                  Smart recommendations and automated trading with advanced AI.
                </p>
              </div>

              <div className="bg-blytz-dark border border-blytz-neon/20 rounded-lg p-6 hover:border-blytz-neon/50 transition-all">
                <Shield className="h-12 w-12 text-blytz-neon mb-4 mx-auto" />
                <h3 className="text-xl font-bold mb-2 text-blytz-neon">
                  Secure Trading
                </h3>
                <p className="text-blytz-lime">
                  Military-grade encryption protects all your transactions.
                </p>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* Status Indicator */}
      <div className="fixed bottom-4 right-4 z-50">
        <div className={`flex items-center space-x-2 px-4 py-2 rounded-full text-sm font-medium ${
          isConnected 
            ? 'bg-green-500/20 text-green-400 border border-green-500/30' 
            : 'bg-red-500/20 text-red-400 border border-red-500/30'
        }`}>
          <div className={`h-2 w-2 rounded-full ${
            isConnected ? 'bg-green-400 animate-pulse' : 'bg-red-400'
          }`} />
          {isConnected ? 'System Online' : 'System Offline'}
        </div>
      </div>
    </div>
  )
}

export default App