import { useState } from "react";
import { Link, useLocation } from "wouter";
import { 
  LayoutDashboard, 
  Users, 
  Building2, 
  ShieldCheck, 
  ShoppingBag, 
  Truck, 
  Package, 
  FileText, 
  ShoppingCart, 
  Menu,
  Search,
  LogOut,
  Settings,
  User,
  Warehouse,
  ChevronDown,
  DollarSign,
  Scale,
  Tag,
  BookOpen,
  Receipt,
  CreditCard,
  ClipboardList,
  Boxes
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Sheet, SheetContent, SheetTrigger, SheetTitle } from "@/components/ui/sheet";
import { VisuallyHidden } from "@radix-ui/react-visually-hidden";
import { ScrollArea } from "@/components/ui/scroll-area";
import { NotificationCenter } from "@/components/NotificationCenter";
import { useAuth } from "@/contexts/AuthContext";

interface LayoutProps {
  children: React.ReactNode;
}

const navItems = [
  // Overview
  { name: "Dashboard", href: "/", icon: LayoutDashboard, section: "Overview" },
  
  // Master Data
  { name: "Customers", href: "/customers", icon: ShoppingBag, section: "Master Data" },
  { name: "Vendors", href: "/vendors", icon: Truck, section: "Master Data" },
  { name: "Products", href: "/products", icon: Package, section: "Master Data" },
  { name: "Employees", href: "/employees", icon: Users, section: "Master Data" },
  
  // Transactions
  { name: "Sales Orders", href: "/sales-orders", icon: ClipboardList, section: "Transactions" },
  { name: "Purchase Orders", href: "/purchase-orders", icon: ShoppingCart, section: "Transactions" },
  { name: "Inventory", href: "/inventory", icon: Boxes, section: "Transactions" },
  
  // Financials
  { name: "Accounts Receivable", href: "/financials/ar", icon: Receipt, section: "Financials" },
  { name: "Accounts Payable", href: "/financials/ap", icon: CreditCard, section: "Financials" },
  { name: "General Ledger", href: "/gl", icon: BookOpen, section: "Financials" },
  { name: "Pricing", href: "/pricing", icon: Tag, section: "Financials" },
  
  // Operations
  { name: "Catch Weight", href: "/operations/catch-weight", icon: Scale, section: "Operations" },
  
  // Admin
  { name: "Departments", href: "/admin/departments", icon: Building2, section: "Admin" },
  { name: "Roles & Permissions", href: "/admin/roles", icon: ShieldCheck, section: "Admin" },
  { name: "Warehouses", href: "/admin/warehouses", icon: Warehouse, section: "Admin" },
];

export default function Layout({ children }: LayoutProps) {
  const [location, setLocation] = useLocation();
  const [isMobileOpen, setIsMobileOpen] = useState(false);
  const { user, logout } = useAuth();

  const NavContent = () => (
    <div className="flex flex-col h-full">
      <div className="flex items-center h-16 px-6 border-b border-sidebar-border">
        <div className="flex items-center gap-2 font-bold text-xl text-primary">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-primary-foreground">
            FH
          </div>
          FoodHive ERP
        </div>
      </div>
      <ScrollArea className="flex-1 py-4">
        <nav className="px-4 space-y-1">
          {(() => {
            let currentSection = "";
            return navItems.map((item, index) => {
              const isActive = location === item.href || (item.href !== "/" && location.startsWith(item.href));
              const showSection = item.section !== currentSection;
              if (showSection) currentSection = item.section || "";
              
              return (
                <div key={item.href}>
                  {showSection && item.section && (
                    <div className="px-3 pt-4 pb-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                      {item.section}
                    </div>
                  )}
                  <Link href={item.href}>
                    <div
                      className={cn(
                        "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer",
                        isActive
                          ? "bg-sidebar-accent text-sidebar-accent-foreground"
                          : "text-sidebar-foreground hover:bg-sidebar-accent/50 hover:text-sidebar-accent-foreground"
                      )}
                    >
                      <item.icon className="w-4 h-4" />
                      {item.name}
                    </div>
                  </Link>
                </div>
              );
            });
          })()}
        </nav>
      </ScrollArea>
      <div className="p-4 border-t border-sidebar-border">
        <div className="flex items-center gap-3 px-3 py-2">
          <Avatar className="w-8 h-8">
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback>{user?.name?.charAt(0) || 'U'}</AvatarFallback>
          </Avatar>
          <div className="flex flex-col">
            <span className="text-sm font-medium">{user?.name || 'User'}</span>
            <span className="text-xs text-muted-foreground capitalize">{user?.role || 'Guest'}</span>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-background flex">
      {/* Desktop Sidebar */}
      <aside className="hidden lg:block w-64 border-r border-sidebar-border bg-sidebar fixed inset-y-0 z-30">
        <NavContent />
      </aside>

      {/* Mobile Sidebar */}
      <Sheet open={isMobileOpen} onOpenChange={setIsMobileOpen}>
        <SheetContent side="left" className="p-0 w-64 bg-sidebar border-r border-sidebar-border">
          <VisuallyHidden>
            <SheetTitle>Navigation Menu</SheetTitle>
          </VisuallyHidden>
          <NavContent />
        </SheetContent>
      </Sheet>

      {/* Main Content */}
      <div className="flex-1 lg:ml-64 flex flex-col min-h-screen">
        {/* Header */}
        <header className="h-16 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-20 px-6 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" className="lg:hidden" onClick={() => setIsMobileOpen(true)}>
              <Menu className="w-5 h-5" />
            </Button>
            <div className="relative hidden md:block w-96">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <input
                type="text"
                placeholder="Search anything..."
                className="w-full h-9 pl-9 pr-4 rounded-md border border-input bg-background text-sm outline-none focus:ring-2 focus:ring-ring"
              />
            </div>
          </div>

          <div className="flex items-center gap-4">
            <NotificationCenter />
            
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="flex items-center gap-2 pl-2 pr-1">
                  <Avatar className="w-8 h-8">
                    <AvatarImage src="https://github.com/shadcn.png" />
                    <AvatarFallback>{user?.name?.charAt(0) || 'U'}</AvatarFallback>
                  </Avatar>
                  <span className="hidden md:inline-block text-sm font-medium">{user?.name || 'User'}</span>
                  <ChevronDown className="w-4 h-4 text-muted-foreground" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel>My Account</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => setLocation('/profile')}>
                  <User className="mr-2 h-4 w-4" />
                  Profile
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => setLocation('/profile')}>
                  <ShieldCheck className="mr-2 h-4 w-4" />
                  Settings
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem 
                  className="text-destructive focus:text-destructive"
                  onClick={logout}
                >
                  <LogOut className="mr-2 h-4 w-4" />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 p-6">
          <div className="max-w-7xl mx-auto space-y-6">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
