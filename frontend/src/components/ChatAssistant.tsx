import React, { useState, useEffect, useRef } from 'react';
import { X, Send, Bot } from 'lucide-react';
import { GoogleGenAI } from '@google/genai';
import { PRODUCTS } from '../constants';

interface ChatAssistantProps {
  onClose: () => void;
}

export const ChatAssistant: React.FC<ChatAssistantProps> = ({ onClose }) => {
  const [messages, setMessages] = useState<{role: 'user' | 'model', text: string}[]>([
    { role: 'model', text: "SYSTEM ONLINE. I am Blytz AI. searching inventory... How can I assist your acquisition today?" }
  ]);
  const [input, setInput] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
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

      <div ref={scrollRef} className="h-80 overflow-y-auto p-4 space-y-3 bg-black/80 backdrop-blur-sm">
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
