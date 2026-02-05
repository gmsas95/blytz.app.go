import { BottomNav } from '@/components/navigation/bottom-nav';
import { Navbar } from '@/components/navigation/navbar';

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-background pb-safe">
      <Navbar />
      <main className="container mx-auto max-w-7xl px-4 pt-16 pb-20 sm:px-6 lg:px-8">
        {children}
      </main>
      <BottomNav />
    </div>
  );
}
