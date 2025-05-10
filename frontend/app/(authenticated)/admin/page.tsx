"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/app/contexts/AuthContext";
import { adminApi } from "@/app/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface Stats {
  total_users: number;
  total_items: number;
  lost_items: number;
  found_items: number;
  resolved_items: number;
  items_this_month: number;
}

interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  first_name: string;
  last_name: string;
  created_at: string;
}

interface Report {
  id: number;
  reporter_id: number;
  reporter_username: string;
  reported_id: number;
  reported_username: string;
  reason: string;
  status: string;
  comment: string;
  created_at: string;
}

export default function AdminDashboard() {
  const { isAdmin, isGuard, token } = useAuth();
  const router = useRouter();
  const [stats, setStats] = useState<Stats | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [reports, setReports] = useState<Report[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!isAdmin && !isGuard) {
      router.push("/dashboard");
      return;
    }
    
    fetchAdminData();
  }, [isAdmin, isGuard, router, token]);

  const fetchAdminData = async () => {
    if (!token) return;
    setLoading(true);
    
    try {
      // Fetch stats
      const statsResponse = await adminApi.getStats(token);
      if (statsResponse.data) {
        setStats(statsResponse.data);
      }
      
      // Fetch users
      const usersResponse = await adminApi.getAllUsers(token);
      if (usersResponse.data) {
        setUsers(usersResponse.data);
      }
      
      // Fetch reports
      const reportsResponse = await adminApi.getAllReports(token);
      if (reportsResponse.data) {
        setReports(reportsResponse.data);
      }
    } catch (err: any) {
      setError(err.message || "Failed to fetch admin data");
    } finally {
      setLoading(false);
    }
  };

  const handleRoleChange = async (userId: number, role: string) => {
    if (!token) return;
    
    try {
      const response = await adminApi.updateUserRole(token, userId, { role });
      if (response.data) {
        // Update the user in the local state
        setUsers(users.map(user => 
          user.id === userId ? { ...user, role } : user
        ));
      }
    } catch (error) {
      console.error("Error updating user role:", error);
    }
  };

  const handleReportUpdate = async (reportId: number, status: string, comment: string = "") => {
    if (!token) return;
    
    try {
      const response = await adminApi.updateReportStatus(token, reportId, { status, comment });
      if (response.data) {
        // Update the report in the local state
        setReports(reports.map(report => 
          report.id === reportId ? { ...report, status, comment } : report
        ));
      }
    } catch (error) {
      console.error("Error updating report status:", error);
    }
  };

  if (loading) {
    return (
      <div className="flex h-[80vh] items-center justify-center">
        <p>Loading admin dashboard...</p>
      </div>
    );
  }

  if (!isAdmin && !isGuard) {
    return null; // We're redirecting in the useEffect
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-6 text-3xl font-bold">Admin Dashboard</h1>
      
      {error && (
        <Alert variant="destructive" className="mb-6">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      <div className="mb-8 grid grid-cols-1 gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Total Users</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.total_users || 0}</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Total Items</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.total_items || 0}</p>
            <div className="mt-2 flex justify-between text-sm">
              <span>Lost: {stats?.lost_items || 0}</span>
              <span>Found: {stats?.found_items || 0}</span>
              <span>Resolved: {stats?.resolved_items || 0}</span>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Items This Month</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.items_this_month || 0}</p>
          </CardContent>
        </Card>
      </div>
      
      <Tabs defaultValue="users" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="users">Users Management</TabsTrigger>
          <TabsTrigger value="reports">Reports Management</TabsTrigger>
        </TabsList>
        
        <TabsContent value="users" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Users</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full border-collapse">
                  <thead>
                    <tr className="border-b">
                      <th className="px-4 py-2 text-left">ID</th>
                      <th className="px-4 py-2 text-left">Username</th>
                      <th className="px-4 py-2 text-left">Name</th>
                      <th className="px-4 py-2 text-left">Email</th>
                      <th className="px-4 py-2 text-left">Role</th>
                      <th className="px-4 py-2 text-left">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map((user) => (
                      <tr key={user.id} className="border-b hover:bg-gray-50">
                        <td className="px-4 py-2">{user.id}</td>
                        <td className="px-4 py-2">{user.username}</td>
                        <td className="px-4 py-2">{user.first_name} {user.last_name}</td>
                        <td className="px-4 py-2">{user.email}</td>
                        <td className="px-4 py-2 capitalize">{user.role}</td>
                        <td className="px-4 py-2">
                          <div className="flex gap-2">
                            {isAdmin && (
                              <>
                                <Button 
                                  size="sm" 
                                  variant={user.role === "student" ? "default" : "outline"} 
                                  onClick={() => handleRoleChange(user.id, "student")}
                                >
                                  Student
                                </Button>
                                <Button 
                                  size="sm" 
                                  variant={user.role === "guard" ? "default" : "outline"} 
                                  onClick={() => handleRoleChange(user.id, "guard")}
                                >
                                  Guard
                                </Button>
                                <Button 
                                  size="sm" 
                                  variant={user.role === "admin" ? "default" : "outline"} 
                                  onClick={() => handleRoleChange(user.id, "admin")}
                                >
                                  Admin
                                </Button>
                              </>
                            )}
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="reports" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Reports</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full border-collapse">
                  <thead>
                    <tr className="border-b">
                      <th className="px-4 py-2 text-left">ID</th>
                      <th className="px-4 py-2 text-left">Reporter</th>
                      <th className="px-4 py-2 text-left">Reported User</th>
                      <th className="px-4 py-2 text-left">Reason</th>
                      <th className="px-4 py-2 text-left">Status</th>
                      <th className="px-4 py-2 text-left">Created</th>
                      <th className="px-4 py-2 text-left">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {reports.map((report) => (
                      <tr key={report.id} className="border-b hover:bg-gray-50">
                        <td className="px-4 py-2">{report.id}</td>
                        <td className="px-4 py-2">{report.reporter_username}</td>
                        <td className="px-4 py-2">{report.reported_username}</td>
                        <td className="px-4 py-2">{report.reason}</td>
                        <td className="px-4 py-2 capitalize">{report.status}</td>
                        <td className="px-4 py-2">
                          {new Date(report.created_at).toLocaleDateString()}
                        </td>
                        <td className="px-4 py-2">
                          <div className="flex gap-2">
                            <Button 
                              size="sm" 
                              variant={report.status === "pending" ? "default" : "outline"} 
                              onClick={() => handleReportUpdate(report.id, "pending")}
                            >
                              Pending
                            </Button>
                            <Button 
                              size="sm" 
                              variant={report.status === "investigating" ? "default" : "outline"} 
                              onClick={() => handleReportUpdate(report.id, "investigating")}
                            >
                              Investigating
                            </Button>
                            <Button 
                              size="sm" 
                              variant={report.status === "resolved" ? "default" : "outline"} 
                              onClick={() => handleReportUpdate(report.id, "resolved", "Issue has been resolved")}
                            >
                              Resolve
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
