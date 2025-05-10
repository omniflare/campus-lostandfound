"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/app/contexts/AuthContext";
import { Button } from "@/components/ui/button";

export function Navigation() {
  const pathname = usePathname();
  const { user, logout, isAdmin, isGuard } = useAuth();

  const isActive = (path: string) => {
    // Normalize path to handle route groups
    const normalizedPathname = pathname?.replace(/^\/(authenticated)/, '') || '';
    return normalizedPathname === path || normalizedPathname?.startsWith(`${path}/`);
  };

  return (
    <nav className="bg-white shadow">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <Link href="/dashboard" className="text-xl font-bold text-primary">
              Campus L&F
            </Link>
            <div className="ml-10 hidden space-x-4 md:flex">
              <Link
                href="/dashboard"
                className={`rounded-md px-3 py-2 text-sm font-medium ${
                  isActive("/dashboard")
                    ? "bg-primary text-white"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                Dashboard
              </Link>
              <Link
                href="/items/report-lost"
                className={`rounded-md px-3 py-2 text-sm font-medium ${
                  isActive("/items/report-lost")
                    ? "bg-primary text-white"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                Report Lost
              </Link>
              <Link
                href="/items/report-found"
                className={`rounded-md px-3 py-2 text-sm font-medium ${
                  isActive("/items/report-found")
                    ? "bg-primary text-white"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                Report Found
              </Link>
              <Link
                href="/messages"
                className={`rounded-md px-3 py-2 text-sm font-medium ${
                  isActive("/messages")
                    ? "bg-primary text-white"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                Messages
              </Link>
              {(isAdmin || isGuard) && (
                <Link
                  href="/admin"
                  className={`rounded-md px-3 py-2 text-sm font-medium ${
                    isActive("/admin")
                      ? "bg-primary text-white"
                      : "text-gray-700 hover:bg-gray-100"
                  }`}
                >
                  Admin Panel
                </Link>
              )}
            </div>
          </div>
          <div className="flex items-center">
            <div className="hidden md:block">
              <div className="ml-4 flex items-center md:ml-6">
                <Link
                  href="/profile"
                  className={`mx-2 rounded-md px-3 py-2 text-sm font-medium ${
                    isActive("/profile")
                      ? "bg-primary text-white"
                      : "text-gray-700 hover:bg-gray-100"
                  }`}
                >
                  {user?.username || "Profile"}
                </Link>
                <Button
                  variant="ghost"
                  onClick={logout}
                  className="mx-2 rounded-md px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100"
                >
                  Logout
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
