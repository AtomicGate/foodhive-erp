import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { 
  DollarSign, 
  ShoppingCart, 
  Activity,
  ArrowUpRight,
  ArrowDownRight,
  Package,
  AlertTriangle,
  Truck,
  Plus,
  X
} from "lucide-react";
import { 
  Area, 
  AreaChart, 
  ResponsiveContainer, 
  Tooltip, 
  XAxis, 
  YAxis,
  CartesianGrid,
} from "recharts";
import { useQuery } from "@tanstack/react-query";
import { dashboardService } from "@/services/dashboardService";
import { useAuth } from "@/contexts/AuthContext";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";

export default function Dashboard() {
  const { user, hasRole } = useAuth();
  const [visibleWidgets, setVisibleWidgets] = useState<string[]>([
    "revenue", "orders", "inventory", "employees", "chart", "recentSales"
  ]);
  
  type TrendType = 'up' | 'down';
  type StatData = { total: number; change: number; trend: TrendType };
  type StatsType = { revenue: StatData; orders: StatData; inventory: StatData; employees: StatData };
  
  const { data: stats, isLoading: statsLoading, isError: statsError } = useQuery<StatsType>({
    queryKey: ['dashboard-stats'],
    queryFn: dashboardService.getStats,
    retry: 1,
    staleTime: 30000,
  });

  const { data: recentSales = [], isLoading: salesLoading } = useQuery({
    queryKey: ['dashboard-recent-sales'],
    queryFn: dashboardService.getRecentSales,
    retry: 1,
  });

  const { data: chartData = [], isLoading: chartLoading } = useQuery({
    queryKey: ['dashboard-chart'],
    queryFn: () => dashboardService.getRevenueChart('year'),
    retry: 1,
  });

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value);
  };

  const toggleWidget = (widgetId: string) => {
    setVisibleWidgets(prev => 
      prev.includes(widgetId) 
        ? prev.filter(id => id !== widgetId)
        : [...prev, widgetId]
    );
  };

  // Safe access to stats properties with proper typing
  const defaultStat: StatData = { total: 0, change: 0, trend: 'up' };
  const revenue = stats?.revenue || defaultStat;
  const orders = stats?.orders || defaultStat;
  const inventory = stats?.inventory || defaultStat;
  const employees = stats?.employees || { total: 0, change: 0, trend: 'up' };

  // Ensure recentSales is an array
  const safeRecentSales = Array.isArray(recentSales) ? recentSales : [];

  // Ensure chartData is an array
  const safeChartData = Array.isArray(chartData) ? chartData : [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Welcome back, {user?.name} ({user?.role})
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground mr-2">Last updated: Today at 09:00 AM</span>
          
          <Dialog>
            <DialogTrigger asChild>
              <Button variant="outline" size="sm">
                <Plus className="mr-2 h-4 w-4" />
                Customize Widgets
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Customize Dashboard</DialogTitle>
                <DialogDescription>
                  Select which widgets you want to see on your dashboard.
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="widget-revenue" 
                    checked={visibleWidgets.includes("revenue")}
                    onCheckedChange={() => toggleWidget("revenue")}
                  />
                  <Label htmlFor="widget-revenue">Total Revenue</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="widget-orders" 
                    checked={visibleWidgets.includes("orders")}
                    onCheckedChange={() => toggleWidget("orders")}
                  />
                  <Label htmlFor="widget-orders">Active Orders</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="widget-inventory" 
                    checked={visibleWidgets.includes("inventory")}
                    onCheckedChange={() => toggleWidget("inventory")}
                  />
                  <Label htmlFor="widget-inventory">Inventory Items</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="widget-employees" 
                    checked={visibleWidgets.includes("employees")}
                    onCheckedChange={() => toggleWidget("employees")}
                  />
                  <Label htmlFor="widget-employees">Active Employees</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="widget-chart" 
                    checked={visibleWidgets.includes("chart")}
                    onCheckedChange={() => toggleWidget("chart")}
                  />
                  <Label htmlFor="widget-chart">Revenue Overview Chart</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="widget-recentSales" 
                    checked={visibleWidgets.includes("recentSales")}
                    onCheckedChange={() => toggleWidget("recentSales")}
                  />
                  <Label htmlFor="widget-recentSales">Recent Sales List</Label>
                </div>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Stats Cards - Role Based Visibility */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {/* Revenue - Visible to Admin, Finance, Sales */}
        {hasRole(['admin', 'finance', 'sales']) && visibleWidgets.includes("revenue") && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatCurrency(revenue.total)}</div>
              <p className="text-xs text-muted-foreground flex items-center mt-1">
                {revenue.trend === 'up' ? (
                  <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
                ) : (
                  <ArrowDownRight className="h-3 w-3 text-rose-500 mr-1" />
                )}
                <span className={revenue.trend === 'up' ? "text-emerald-500 font-medium" : "text-rose-500 font-medium"}>
                  {revenue.change > 0 ? '+' : ''}{revenue.change}%
                </span>
                <span className="ml-1">from last month</span>
              </p>
            </CardContent>
          </Card>
        )}

        {/* Active Orders - Visible to Admin, Sales, Warehouse */}
        {hasRole(['admin', 'sales', 'warehouse']) && visibleWidgets.includes("orders") && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Orders</CardTitle>
              <ShoppingCart className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">+{orders.total}</div>
              <p className="text-xs text-muted-foreground flex items-center mt-1">
                {orders.trend === 'up' ? (
                  <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
                ) : (
                  <ArrowDownRight className="h-3 w-3 text-rose-500 mr-1" />
                )}
                <span className={orders.trend === 'up' ? "text-emerald-500 font-medium" : "text-rose-500 font-medium"}>
                  {orders.change > 0 ? '+' : ''}{orders.change}%
                </span>
                <span className="ml-1">from last month</span>
              </p>
            </CardContent>
          </Card>
        )}

        {/* Inventory - Visible to Admin, Warehouse, Sales */}
        {hasRole(['admin', 'warehouse', 'sales']) && visibleWidgets.includes("inventory") && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Inventory Items</CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{inventory.total.toLocaleString()}</div>
              <p className="text-xs text-muted-foreground flex items-center mt-1">
                {inventory.trend === 'up' ? (
                  <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
                ) : (
                  <ArrowDownRight className="h-3 w-3 text-rose-500 mr-1" />
                )}
                <span className={inventory.trend === 'up' ? "text-emerald-500 font-medium" : "text-rose-500 font-medium"}>
                  {inventory.change > 0 ? '+' : ''}{inventory.change}%
                </span>
                <span className="ml-1">from last month</span>
              </p>
            </CardContent>
          </Card>
        )}

        {/* Employees - Visible to Admin only */}
        {hasRole('admin') && visibleWidgets.includes("employees") && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Employees</CardTitle>
              <Activity className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{employees.total}</div>
              <p className="text-xs text-muted-foreground flex items-center mt-1">
                {employees.trend === 'up' ? (
                  <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
                ) : (
                  <ArrowDownRight className="h-3 w-3 text-rose-500 mr-1" />
                )}
                <span className={employees.trend === 'up' ? "text-emerald-500 font-medium" : "text-rose-500 font-medium"}>
                  {employees.change > 0 ? '+' : ''}{employees.change}
                </span>
                <span className="ml-1">since last hour</span>
              </p>
            </CardContent>
          </Card>
        )}

        {/* Warehouse Specific: Pending Shipments */}
        {hasRole('warehouse') && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Pending Shipments</CardTitle>
              <Truck className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">12</div>
              <p className="text-xs text-muted-foreground flex items-center mt-1">
                <span className="text-amber-500 font-medium">3 urgent</span>
                <span className="ml-1">need attention</span>
              </p>
            </CardContent>
          </Card>
        )}

        {/* Finance Specific: Overdue Invoices */}
        {hasRole('finance') && (
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Overdue Invoices</CardTitle>
              <AlertTriangle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">5</div>
              <p className="text-xs text-muted-foreground flex items-center mt-1">
                <span className="text-rose-500 font-medium">$12,450</span>
                <span className="ml-1">total outstanding</span>
              </p>
            </CardContent>
          </Card>
        )}
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Main Chart - Visible to Admin, Finance, Sales */}
        {hasRole(['admin', 'finance', 'sales']) && visibleWidgets.includes("chart") && (
          <Card className="col-span-4">
            <CardHeader>
              <CardTitle>Overview</CardTitle>
              <CardDescription>
                Monthly revenue and sales performance for the current year.
              </CardDescription>
            </CardHeader>
            <CardContent className="pl-2">
              <div className="h-[350px]">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={safeChartData}>
                    <defs>
                      <linearGradient id="colorTotal" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="var(--primary)" stopOpacity={0.3}/>
                        <stop offset="95%" stopColor="var(--primary)" stopOpacity={0}/>
                      </linearGradient>
                    </defs>
                    <XAxis 
                      dataKey="name" 
                      stroke="#888888" 
                      fontSize={12} 
                      tickLine={false} 
                      axisLine={false} 
                    />
                    <YAxis 
                      stroke="#888888" 
                      fontSize={12} 
                      tickLine={false} 
                      axisLine={false} 
                      tickFormatter={(value) => `$${value}`} 
                    />
                    <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="var(--border)" />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'var(--background)', 
                        borderColor: 'var(--border)',
                        borderRadius: 'var(--radius)'
                      }}
                    />
                    <Area 
                      type="monotone" 
                      dataKey="total" 
                      stroke="var(--primary)" 
                      fillOpacity={1} 
                      fill="url(#colorTotal)" 
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Recent Sales - Visible to Admin, Sales */}
        {hasRole(['admin', 'sales']) && visibleWidgets.includes("recentSales") && (
          <Card className="col-span-3">
            <CardHeader>
              <CardTitle>Recent Sales</CardTitle>
              <CardDescription>
                You made {safeRecentSales.length} sales this month.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-8">
                {safeRecentSales.map((sale: any, index: number) => (
                  <div key={index} className="flex items-center">
                    <div className="h-9 w-9 rounded-full bg-primary/10 flex items-center justify-center text-primary font-medium">
                      {sale.name.charAt(0)}{sale.name.split(" ")[1]?.charAt(0)}
                    </div>
                    <div className="ml-4 space-y-1">
                      <p className="text-sm font-medium leading-none">{sale.name}</p>
                      <p className="text-xs text-muted-foreground">{sale.email}</p>
                    </div>
                    <div className="ml-auto font-medium">{sale.amount}</div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
