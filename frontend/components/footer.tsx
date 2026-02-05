'use client';

import Link from 'next/link';
import { Instagram, Twitter, Facebook, Youtube, Mail, MapPin, Phone } from 'lucide-react';

const footerLinks = {
  platform: [
    { name: 'Browse Auctions', href: '/auctions' },
    { name: 'Live Streams', href: '/streams' },
    { name: 'Sell Items', href: '/sell' },
    { name: 'How it Works', href: '/how-it-works' },
  ],
  support: [
    { name: 'Help Center', href: '/help' },
    { name: 'Safety Center', href: '/safety' },
    { name: 'Community', href: '/community' },
    { name: 'Guidelines', href: '/guidelines' },
  ],
  company: [
    { name: 'About Us', href: '/about' },
    { name: 'Careers', href: '/careers' },
    { name: 'Press', href: '/press' },
    { name: 'Contact', href: '/contact' },
  ],
  legal: [
    { name: 'Terms of Service', href: '/terms' },
    { name: 'Privacy Policy', href: '/privacy' },
    { name: 'Cookie Policy', href: '/cookies' },
  ],
};

const socialLinks = [
  { name: 'Instagram', icon: Instagram, href: '#' },
  { name: 'Twitter', icon: Twitter, href: '#' },
  { name: 'Facebook', icon: Facebook, href: '#' },
  { name: 'YouTube', icon: Youtube, href: '#' },
];

export function Footer() {
  return (
    <footer className="bg-black border-t border-neutral-800">
      {/* Main Footer */}
      <div className="container-modern py-16">
        <div className="grid lg:grid-cols-6 gap-12">
          {/* Brand */}
          <div className="lg:col-span-2 space-y-6">
            <Link href="/" className="flex items-center gap-3">
              <div className="w-12 h-12 bg-blytz-yellow rounded-xl flex items-center justify-center">
                <span className="text-black font-black text-2xl">B</span>
              </div>
              <span className="text-2xl font-black text-white">BLYTZ</span>
            </Link>
            <p className="text-gray-400 max-w-xs">
              Malaysia&apos;s premier live auction marketplace. Discover unique items, bid in real-time, and win amazing deals.
            </p>
            
            {/* Social Links */}
            <div className="flex gap-3">
              {socialLinks.map((social) => (
                <a
                  key={social.name}
                  href={social.href}
                  className="w-10 h-10 rounded-full bg-neutral-900 flex items-center justify-center text-gray-400 hover:bg-blytz-yellow hover:text-black transition-all"
                  aria-label={social.name}
                >
                  <social.icon className="w-5 h-5" />
                </a>
              ))}
            </div>
          </div>

          {/* Links */}
          <div>
            <h4 className="text-white font-bold mb-4">Platform</h4>
            <ul className="space-y-3">
              {footerLinks.platform.map((link) => (
                <li key={link.name}>
                  <Link href={link.href} className="text-gray-400 hover:text-blytz-yellow transition-colors text-sm">
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-white font-bold mb-4">Support</h4>
            <ul className="space-y-3">
              {footerLinks.support.map((link) => (
                <li key={link.name}>
                  <Link href={link.href} className="text-gray-400 hover:text-blytz-yellow transition-colors text-sm">
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-white font-bold mb-4">Company</h4>
            <ul className="space-y-3">
              {footerLinks.company.map((link) => (
                <li key={link.name}>
                  <Link href={link.href} className="text-gray-400 hover:text-blytz-yellow transition-colors text-sm">
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-white font-bold mb-4">Legal</h4>
            <ul className="space-y-3">
              {footerLinks.legal.map((link) => (
                <li key={link.name}>
                  <Link href={link.href} className="text-gray-400 hover:text-blytz-yellow transition-colors text-sm">
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>

      {/* Bottom Bar */}
      <div className="border-t border-neutral-900">
        <div className="container-modern py-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <p className="text-gray-500 text-sm">
              © {new Date().getFullYear()} Blytz.live. All rights reserved.
            </p>
            <div className="flex items-center gap-6">
              <div className="flex items-center gap-2 text-gray-500 text-sm">
                <MapPin className="w-4 h-4" />
                <span>Kuala Lumpur, Malaysia</span>
              </div>
              <div className="flex items-center gap-2 text-gray-500 text-sm">
                <Mail className="w-4 h-4" />
                <span>hello@blytz.live</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
