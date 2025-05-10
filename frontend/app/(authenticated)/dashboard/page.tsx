"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useAuth } from "@/app/contexts/AuthContext";
import { itemsApi } from "@/app/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

interface Item {
  id: number;
  title: string;
  description: string;
  category: string;
  location: string;
  status: string;
  type: string;
  created_at: string;
  image_url?: string;
}

// ItemResponse can be either an array of items, or an object with an items property
type ItemResponse = Item[] | { items: Item[] };

export default function Dashboard() {
  const { user } = useAuth();
  const [items, setItems] = useState<Item[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const response = await itemsApi.getAllItems();
      console.log("API Response Structure:", response);
      
      if (response.data) {
        // Type assertion to handle different response formats
        const itemsData = response.data as ItemResponse;
        
        if (Array.isArray(itemsData)) {
          setItems(itemsData);
        } else if ('items' in itemsData && Array.isArray(itemsData.items)) {
          setItems(itemsData.items);
        } else {
          console.error("Unexpected response format:", itemsData);
          setItems([]);
        }
      } else {
        setItems([]);
      }
    } catch (error) {
      console.error("Failed to fetch items:", error);
      setItems([]);
    }
  };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) {
      return fetchItems();
    }

    setIsSearching(true);
    try {
      const response = await itemsApi.searchItems(searchQuery);
      if (response.data) {
        // Apply the same format handling as in fetchItems
        const searchData = response.data as ItemResponse;
        
        if (Array.isArray(searchData)) {
          setItems(searchData);
        } else if ('items' in searchData && Array.isArray(searchData.items)) {
          setItems(searchData.items);
        } else {
          console.error("Unexpected search response format:", searchData);
          setItems([]);
        }
      } else {
        setItems([]);
      }
    } catch (error) {
      console.error("Search failed:", error);
    } finally {
      setIsSearching(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "lost":
        return "bg-red-100 text-red-800";
      case "found":
        return "bg-green-100 text-green-800";
      case "claimed":
        return "bg-blue-100 text-blue-800";
      case "resolved":
        return "bg-purple-100 text-purple-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="mt-1 text-gray-600">
            Welcome, {user?.first_name} {user?.last_name}
          </p>
        </div>
        <div className="flex items-center gap-4">
          <Button asChild variant="default">
            <Link href="/items/report-lost">Report Lost Item</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/items/report-found">Report Found Item</Link>
          </Button>
        </div>
      </div>

      <div className="mb-6">
        <form onSubmit={handleSearch} className="flex max-w-lg gap-2">
          <Input
            type="text"
            placeholder="Search for items..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full"
          />
          <Button type="submit" disabled={isSearching}>
            {isSearching ? "Searching..." : "Search"}
          </Button>
        </form>
      </div>

      <div className="mb-8">
        <h2 className="mb-4 text-2xl font-semibold">Recent Items</h2>
        {items.length === 0 ? (
          <p className="text-gray-500">No items found</p>
        ) : (
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            {items.map((item) => (
              <Card key={item.id} className="overflow-hidden">
                {item.image_url && (
                  <div className="h-48 w-full overflow-hidden">
                    <img
                      src={item.image_url}
                      alt={item.title}
                      className="h-full w-full object-cover"
                    />
                  </div>
                )}
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-xl">{item.title}</CardTitle>
                    <span
                      className={`rounded-full px-3 py-1 text-xs font-medium ${getStatusColor(
                        item.status
                      )}`}
                    >
                      {item.type === "lost" ? "Lost" : "Found"}
                    </span>
                  </div>
                  <CardDescription>{item.category}</CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="mb-2 line-clamp-2 text-sm text-gray-600">
                    {item.description}
                  </p>
                  <p className="text-sm text-gray-600">
                    <strong>Location:</strong> {item.location}
                  </p>
                </CardContent>
                <CardFooter className="pt-2">
                  <Button asChild variant="outline" size="sm" className="w-full">
                    <Link href={`/items/${item.id}`}>View Details</Link>
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
