// API client for making backend requests

const API_BASE_URL = "http://localhost:3000/api/v1";

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export async function api<T>(
  endpoint: string,
  method: string = "GET",
  body?: any,
  token?: string | null
): Promise<ApiResponse<T>> {
  try {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const options: RequestInit = {
      method,
      headers,
    };

    if (body && (method === "POST" || method === "PUT")) {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
    const data = await response.json();

    if (!response.ok) {
      return { error: data.error || "An error occurred" };
    }

    return { data };
  } catch (error: any) {
    return { error: error.message || "Network error" };
  }
}

// Authentication API
export const authApi = {
  register: (userData: any) => api("/auth/register", "POST", userData),
  login: (credentials: { username: string; password: string }) =>
    api<{ token: string }>("/auth/login", "POST", credentials),
};

// User API
export const userApi = {
  getProfile: (token: string) => api("/user/profile", "GET", undefined, token),
  updateProfile: (token: string, userData: any) =>
    api("/user/profile", "PUT", userData, token),
  changePassword: (token: string, passwordData: any) =>
    api("/user/password", "PUT", passwordData, token),
  getUserItems: (token: string) => api("/user/items", "GET", undefined, token),
  getUnreadMessages: (token: string) =>
    api("/user/messages/unread", "GET", undefined, token),
  getConversations: (token: string) =>
    api("/user/messages/conversations", "GET", undefined, token),
  getMessages: (token: string, userId: number) =>
    api(`/user/messages/${userId}`, "GET", undefined, token),
  sendMessage: (token: string, messageData: any) =>
    api("/user/messages", "POST", messageData, token),
  createReport: (token: string, reportData: any) =>
    api("/user/reports", "POST", reportData, token),
};

// Items API
export const itemsApi = {
  getAllItems: () => api("/items"),
  searchItems: (query: string) => api(`/items/search?q=${query}`),
  getItemDetails: (itemId: number) => api(`/items/${itemId}`),
  reportLostItem: (token: string, itemData: any) =>
    api("/items/lost", "POST", itemData, token),
  reportFoundItem: (token: string, itemData: any) =>
    api("/items/found", "POST", itemData, token),
  updateItemStatus: (token: string, itemId: number, statusData: any) =>
    api(`/items/${itemId}/status`, "PUT", statusData, token),
};

// Admin API
export const adminApi = {
  getAllUsers: (token: string) => api("/admin/users", "GET", undefined, token),
  updateUserRole: (token: string, userId: number, roleData: any) =>
    api(`/admin/users/${userId}/role`, "PUT", roleData, token),
  getAllReports: (token: string) =>
    api("/admin/reports", "GET", undefined, token),
  updateReportStatus: (token: string, reportId: number, statusData: any) =>
    api(`/admin/reports/${reportId}/status`, "PUT", statusData, token),
  getStats: (token: string) => api("/admin/stats", "GET", undefined, token),
};
