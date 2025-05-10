"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/app/contexts/AuthContext";
import { userApi } from "@/app/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface Conversation {
  user_id: number;
  username: string;
  last_message: string;
  last_message_time: string;
  unread_count: number;
}

interface Message {
  id: number;
  sender_id: number;
  receiver_id: number;
  content: string;
  created_at: string;
  is_read: boolean;
  item_id?: number;
  item_title?: string;
}

export default function MessagesPage() {
  const { token, user } = useAuth();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchConversations();
  }, [token]);

  useEffect(() => {
    if (selectedUserId) {
      fetchMessages(selectedUserId);
    }
  }, [selectedUserId, token]);

  const fetchConversations = async () => {
    if (!token) return;
    setLoading(true);

    try {
      const response = await userApi.getConversations(token);
      if (response.data) {
        setConversations(response.data);
        
        // Auto-select the first conversation if any exists
        if (response.data.length > 0 && !selectedUserId) {
          setSelectedUserId(response.data[0].user_id);
        }
      }
    } catch (error) {
      console.error("Error fetching conversations:", error);
      setError("Failed to fetch conversations");
    } finally {
      setLoading(false);
    }
  };

  const fetchMessages = async (userId: number) => {
    if (!token) return;
    
    try {
      const response = await userApi.getMessages(token, userId);
      if (response.data) {
        setMessages(response.data);
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !newMessage.trim() || !selectedUserId) return;

    setSending(true);
    try {
      await fetch(`http://localhost:3000/api/v1/user/messages`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          receiver_id: selectedUserId,
          content: newMessage,
        }),
      });
      
      setNewMessage("");
      // Fetch updated messages
      fetchMessages(selectedUserId);
    } catch (error) {
      console.error("Error sending message:", error);
    } finally {
      setSending(false);
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    const today = new Date();
    
    // If same day, show just the time
    if (date.toDateString() === today.toDateString()) {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    
    // If within the last 7 days, show the day of week
    const daysAgo = Math.floor((today.getTime() - date.getTime()) / (1000 * 60 * 60 * 24));
    if (daysAgo < 7) {
      return date.toLocaleDateString([], { weekday: 'short' }) + ' ' + 
             date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    
    // Otherwise show full date
    return date.toLocaleDateString() + ' ' + 
           date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  if (loading && conversations.length === 0) {
    return (
      <div className="flex h-[80vh] items-center justify-center">
        <p>Loading conversations...</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto max-w-6xl px-4 py-8">
      <h1 className="mb-6 text-3xl font-bold">Messages</h1>
      
      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div className="md:col-span-1">
          <Card className="h-[70vh]">
            <CardHeader>
              <CardTitle>Conversations</CardTitle>
            </CardHeader>
            <CardContent className="overflow-y-auto p-0">
              {conversations.length === 0 ? (
                <p className="px-4 py-8 text-center text-sm text-gray-500">
                  You have no conversations yet
                </p>
              ) : (
                <ul className="divide-y">
                  {conversations.map((convo) => (
                    <li 
                      key={convo.user_id} 
                      className={`cursor-pointer p-4 hover:bg-gray-50 ${
                        selectedUserId === convo.user_id ? "bg-gray-100" : ""
                      }`}
                      onClick={() => setSelectedUserId(convo.user_id)}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <p className="font-medium">{convo.username}</p>
                          <p className="line-clamp-1 text-sm text-gray-600">
                            {convo.last_message}
                          </p>
                        </div>
                        <div className="ml-2 flex flex-col items-end">
                          <span className="text-xs text-gray-500">
                            {formatTime(convo.last_message_time)}
                          </span>
                          {convo.unread_count > 0 && (
                            <span className="mt-1 inline-flex items-center rounded-full bg-primary px-2 py-1 text-xs font-medium text-white">
                              {convo.unread_count}
                            </span>
                          )}
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </CardContent>
          </Card>
        </div>
        
        <div className="md:col-span-2">
          <Card className="flex h-[70vh] flex-col">
            <CardHeader className="border-b pb-3">
              <CardTitle>
                {selectedUserId
                  ? conversations.find(c => c.user_id === selectedUserId)?.username || "Chat"
                  : "Select a conversation"}
              </CardTitle>
            </CardHeader>
            
            <CardContent className="flex-1 overflow-y-auto p-4">
              {!selectedUserId ? (
                <div className="flex h-full items-center justify-center">
                  <p className="text-gray-500">Select a conversation to start messaging</p>
                </div>
              ) : messages.length === 0 ? (
                <div className="flex h-full items-center justify-center">
                  <p className="text-gray-500">No messages yet</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {messages.map((msg) => {
                    const isFromMe = msg.sender_id === user?.id;
                    
                    return (
                      <div 
                        key={msg.id} 
                        className={`flex ${isFromMe ? "justify-end" : "justify-start"}`}
                      >
                        <div 
                          className={`max-w-[75%] rounded-lg px-4 py-2 ${
                            isFromMe 
                              ? "bg-primary text-primary-foreground" 
                              : "bg-gray-100 text-gray-800"
                          }`}
                        >
                          {msg.item_id && (
                            <p className="mb-1 text-xs opacity-80">
                              Re: {msg.item_title}
                            </p>
                          )}
                          <p className="break-words">{msg.content}</p>
                          <p className="mt-1 text-right text-xs opacity-70">
                            {formatTime(msg.created_at)}
                          </p>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </CardContent>
            
            {selectedUserId && (
              <div className="border-t p-4">
                <form onSubmit={handleSendMessage} className="flex gap-2">
                  <Input
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    placeholder="Type your message..."
                    disabled={!selectedUserId}
                  />
                  <Button type="submit" disabled={sending || !newMessage.trim()}>
                    Send
                  </Button>
                </form>
              </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
}
