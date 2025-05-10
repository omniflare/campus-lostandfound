"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useAuth } from "@/app/contexts/AuthContext";
import { itemsApi } from "@/app/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface ItemParams {
  params: {
    id: string;
  };
}

export default function ItemDetailPage({ params }: ItemParams) {
  const { token, user } = useAuth();
  const [item, setItem] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [sending, setSending] = useState(false);
  const [success, setSuccess] = useState("");

  useEffect(() => {
    const fetchItem = async () => {
      try {
        const response = await itemsApi.getItemDetails(parseInt(params.id));
        if (response.data) {
          setItem(response.data);
        } else if (response.error) {
          setError(response.error);
        }
      } catch (err: any) {
        setError(err.message || "Failed to fetch item details");
      } finally {
        setLoading(false);
      }
    };

    fetchItem();
  }, [params.id]);

  const handleStatusUpdate = async (status: string) => {
    if (!token) return;

    try {
      const response = await itemsApi.updateItemStatus(token, parseInt(params.id), { status });
      if (response.data) {
        setItem({ ...item, status });
        setSuccess(`Item status updated to ${status}`);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err: any) {
      setError(err.message || "Failed to update status");
    }
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !message.trim() || !item?.user_id) return;

    setSending(true);
    try {
      await fetch(`http://localhost:3000/api/v1/user/messages`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          receiver_id: item.user_id,
          item_id: item.id,
          content: message,
        }),
      });
      setMessage("");
      setSuccess("Message sent successfully");
    } catch (err: any) {
      setError(err.message || "Failed to send message");
    } finally {
      setSending(false);
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <p>Loading item details...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto max-w-4xl px-4 py-8">
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <div className="mt-4">
          <Button asChild>
            <Link href="/dashboard">Back to Dashboard</Link>
          </Button>
        </div>
      </div>
    );
  }

  if (!item) {
    return (
      <div className="container mx-auto max-w-4xl px-4 py-8">
        <p>Item not found</p>
        <div className="mt-4">
          <Button asChild>
            <Link href="/dashboard">Back to Dashboard</Link>
          </Button>
        </div>
      </div>
    );
  }

  const isOwner = user?.id === item.user_id;
  const statusBadgeColor = 
    item.status === "lost" ? "bg-red-100 text-red-800" :
    item.status === "found" ? "bg-green-100 text-green-800" :
    item.status === "claimed" ? "bg-blue-100 text-blue-800" : 
    "bg-gray-100 text-gray-800";

  return (
    <div className="container mx-auto max-w-4xl px-4 py-8">
      {success && (
        <Alert className="mb-4">
          <AlertDescription>{success}</AlertDescription>
        </Alert>
      )}
      
      <div className="mb-4">
        <Link href="/dashboard" className="text-primary hover:underline">
          &larr; Back to Dashboard
        </Link>
      </div>

      <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
        <div className="md:col-span-2">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-2xl">{item.title}</CardTitle>
                <span className={`rounded-full px-4 py-1 text-sm font-medium ${statusBadgeColor}`}>
                  {item.type === "lost" ? "Lost" : "Found"}
                </span>
              </div>
              <CardDescription>
                {item.category} • {new Date(item.created_at).toLocaleDateString()}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {item.image_url && (
                <div className="overflow-hidden rounded-md">
                  <img
                    src={item.image_url}
                    alt={item.title}
                    className="h-auto w-full object-cover"
                  />
                </div>
              )}
              <div>
                <h3 className="mb-2 text-lg font-semibold">Description</h3>
                <p className="text-gray-700">{item.description}</p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h3 className="text-sm font-medium text-gray-500">Location</h3>
                  <p>{item.location}</p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500">Status</h3>
                  <p className="capitalize">{item.status}</p>
                </div>
                {item.type === "lost" && item.lost_time && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-500">Lost Time</h3>
                    <p>{new Date(item.lost_time).toLocaleString()}</p>
                  </div>
                )}
              </div>

              {isOwner && (
                <div className="border-t pt-4">
                  <h3 className="mb-2 text-lg font-semibold">Update Status</h3>
                  <div className="flex flex-wrap gap-2">
                    {item.status !== "found" && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStatusUpdate("found")}
                      >
                        Mark as Found
                      </Button>
                    )}
                    {item.status !== "claimed" && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStatusUpdate("claimed")}
                      >
                        Mark as Claimed
                      </Button>
                    )}
                    {item.status !== "resolved" && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStatusUpdate("resolved")}
                      >
                        Mark as Resolved
                      </Button>
                    )}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        <div>
          <Card>
            <CardHeader>
              <CardTitle className="text-xl">Contact</CardTitle>
              <CardDescription>
                Send a message about this item
              </CardDescription>
            </CardHeader>
            <CardContent>
              {!isOwner ? (
                <form onSubmit={handleSendMessage} className="space-y-4">
                  <Textarea
                    placeholder="Write a message to the owner..."
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    required
                    rows={4}
                  />
                  <Button type="submit" className="w-full" disabled={sending}>
                    {sending ? "Sending..." : "Send Message"}
                  </Button>
                </form>
              ) : (
                <p className="text-sm text-gray-500">
                  You are the owner of this item. Check your messages for any inquiries.
                </p>
              )}
              
              <div className="mt-4">
                <Button asChild variant="outline" className="w-full">
                  <Link href="/messages">View All Messages</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
