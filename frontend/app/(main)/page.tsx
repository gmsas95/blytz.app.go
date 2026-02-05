import { Metadata } from 'next';
import { HeroSection } from '@/components/sections/hero-section';
import { FeaturesSection } from '@/components/sections/features-section';
import { CategoriesSection } from '@/components/sections/categories-section';
import { FeaturedAuctionsSection } from '@/components/sections/featured-auctions-section';
import { LiveStreamsSection } from '@/components/sections/live-streams-section';
import { AppDownloadSection } from '@/components/sections/app-download-section';

export const metadata: Metadata = {
  title: 'Blytz - Live Auction Marketplace',
  description: 'Discover unique products through live auctions and buy-now deals',
};

export default function HomePage() {
  return (
    <div className="bg-black">
      <HeroSection />
      <FeaturesSection />
      <LiveStreamsSection />
      <FeaturedAuctionsSection />
      <CategoriesSection />
      <AppDownloadSection />
    </div>
  );
}
