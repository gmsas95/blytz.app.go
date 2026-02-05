'use client';

import { Button } from '@/components/ui/button';
import { Apple, Play } from 'lucide-react';

export function AppDownloadSection() {
  return (
    <section className="py-12">
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-secondary via-secondary/90 to-secondary/80 p-8 sm:p-12">
        <div className="relative z-10 flex flex-col lg:flex-row items-center justify-between gap-8">
          <div className="text-center lg:text-left max-w-md">
            <h2 className="text-2xl font-bold mb-3">
              Get the Blytz App
            </h2>
            <p className="text-muted-foreground mb-6">
              Never miss a live auction. Get real-time notifications and bid on the go with our mobile app.
            </p>
            <div className="flex flex-col sm:flex-row gap-3 justify-center lg:justify-start">
              <Button variant="default" className="bg-foreground text-background hover:bg-foreground/90">
                <Apple className="mr-2 h-5 w-5" />
                App Store
              </Button>
              <Button variant="default" className="bg-foreground text-background hover:bg-foreground/90">
                <Play className="mr-2 h-5 w-5" />
                Play Store
              </Button>
            </div>
          </div>
          
          <div className="relative">
            <div className="w-48 h-48 sm:w-56 sm:h-56 rounded-full bg-primary/20 flex items-center justify-center">
              <div className="w-36 h-36 sm:w-44 sm:h-44 rounded-full bg-primary/30 flex items-center justify-center">
                <div className="w-24 h-24 sm:w-32 sm:h-32 rounded-full bg-primary flex items-center justify-center text-primary-foreground text-4xl font-bold">
                  B
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
