import Link from 'next/link';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-white to-gray-100">
      <div className="container mx-auto px-4 py-16">
        <div className="flex flex-col items-center justify-center space-y-12 text-center">
          <div className="space-y-4">
            <h1 className="text-4xl font-extrabold tracking-tight text-primary sm:text-6xl">
              Campus Lost & Found
          </h1>
          <p className="mt-6 text-lg leading-8 text-gray-600">
            The easiest way to report lost items and reconnect with your belongings on campus.
          </p>
          <div className="mt-10 flex items-center justify-center gap-x-6">
            <Link href="/auth/login">
              <Button size="lg">Get Started</Button>
            </Link>
            <Link href="/items" className="text-sm font-semibold leading-6 text-gray-900">
              Browse Items <span aria-hidden="true">→</span>
            </Link>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16">
          <Card>
            <CardHeader>
              <CardTitle>Report Lost Items</CardTitle>
              <CardDescription>Lost something on campus? Report it here.</CardDescription>
            </CardHeader>
            <CardContent>
              <p>Quickly submit details about your lost item. Our system will match it with found items.</p>
            </CardContent>
            <CardFooter>
              <Link href="/items/report-lost">
                <Button variant="outline">Report Lost Item</Button>
              </Link>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Report Found Items</CardTitle>
              <CardDescription>Found something on campus? Report it here.</CardDescription>
            </CardHeader>
            <CardContent>
              <p>Submit details and photos of items you've found. We'll help connect them with their owners.</p>
            </CardContent>
            <CardFooter>
              <Link href="/items/report-found">
                <Button variant="outline">Report Found Item</Button>
              </Link>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Browse Items</CardTitle>
              <CardDescription>Looking for your lost item? Browse through reported items.</CardDescription>
            </CardHeader>
            <CardContent>
              <p>Search through all reported lost and found items across campus.</p>
            </CardContent>
            <CardFooter>
              <Link href="/items">
                <Button variant="outline">Browse Items</Button>
              </Link>
            </CardFooter>
          </Card>
        </div>
        
        <div className="mt-16 text-center">
          <h2 className="text-2xl font-bold mb-4">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="flex flex-col items-center">
              <div className="rounded-full bg-primary/10 p-4 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </div>
              <h3 className="font-semibold">1. Report</h3>
              <p className="text-sm text-gray-600 mt-2">Report your lost item or a found item on campus</p>
            </div>
            
            <div className="flex flex-col items-center">
              <div className="rounded-full bg-primary/10 p-4 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16l2.879-2.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242zM21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="font-semibold">2. Match</h3>
              <p className="text-sm text-gray-600 mt-2">Our system matches lost items with found items</p>
            </div>
            
            <div className="flex flex-col items-center">
              <div className="rounded-full bg-primary/10 p-4 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 className="font-semibold">3. Recover</h3>
              <p className="text-sm text-gray-600 mt-2">Get notified and arrange to recover your item</p>
            </div>
          </div>
        </div>
      </div>
      </div>
    </main>
  );
}
