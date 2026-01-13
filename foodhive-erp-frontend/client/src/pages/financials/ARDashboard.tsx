import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { 
  DollarSign, 
  TrendingUp, 
  AlertCircle, 
  Clock,
  Plus,
  Loader2
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell
} from "recharts";
import { financialService } from "@/services/financialService";

export default function ARDashboard() {
  const { data: agingData, isLoading: isAgingLoading } = useQuery({
    queryKey: ['arAging'],
    queryFn: () => financialService.getARAging()
  });

  const { data: overdueInvoices, isLoading: isOverdueLoading } = useQuery({
    queryKey: ['overdueInvoices'],
    queryFn: financialService.getOverdueInvoices
  });

  if (isAgingLoading || isOverdueLoading) {
    return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  }

  // Handle both snake_case (from API) and camelCase property names
  const aging = agingData || { current: 0, days_1_30: 0, days_31_60: 0, days_61_90: 0, over_90: 0, total: 0 };
  
  const chartData = [
    { name: "Current", value: aging.current || 0, color: "#10b981" },
    { name: "1-30 Days", value: aging.days_1_30 || 0, color: "#3b82f6" },
    { name: "31-60 Days", value: aging.days_31_60 || 0, color: "#f59e0b" },
    { name: "61-90 Days", value: aging.days_61_90 || 0, color: "#f97316" },
    { name: "90+ Days", value: aging.over_90 || 0, color: "#ef4444" },
  ];

  const totalReceivable = aging.total || 0;
  const totalOverdue = (aging.days_1_30 || 0) + (aging.days_31_60 || 0) + (aging.days_61_90 || 0) + (aging.over_90 || 0);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Accounts Receivable</h1>
          <p className="text-muted-foreground">
            Monitor receivables, aging, and collections.
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">View Invoices</Button>
          <Button>
            <Plus className="mr-2 h-4 w-4" /> Record Payment
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Receivable</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(totalReceivable)}
            </div>
            <p className="text-xs text-muted-foreground">
              +12% from last month
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Overdue</CardTitle>
            <AlertCircle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-destructive">
              {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(totalOverdue)}
            </div>
            <p className="text-xs text-muted-foreground">
              {totalReceivable > 0 ? ((totalOverdue / totalReceivable) * 100).toFixed(1) : 0}% of total receivable
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg. Collection Days</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">24 Days</div>
            <p className="text-xs text-muted-foreground">
              -2 days from last month
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Collected This Month</CardTitle>
            <TrendingUp className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$32,450.00</div>
            <p className="text-xs text-muted-foreground">
              On track for monthly goal
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Aging Overview</CardTitle>
          </CardHeader>
          <CardContent className="pl-2">
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
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
                  <Tooltip 
                    cursor={{ fill: 'transparent' }}
                    formatter={(value: number) => [`$${value.toLocaleString()}`, 'Amount']}
                  />
                  <Bar dataKey="value" radius={[4, 4, 0, 0]}>
                    {chartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Top Overdue Accounts</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-8">
              {overdueInvoices?.slice(0, 5).map((invoice: any) => (
                <div key={invoice.id} className="flex items-center">
                  <div className="ml-4 space-y-1">
                    <p className="text-sm font-medium leading-none">{invoice.customer_name || invoice.customerName}</p>
                    <p className="text-sm text-muted-foreground">Inv #{invoice.invoice_number || invoice.invoiceNumber}</p>
                  </div>
                  <div className="ml-auto font-medium text-destructive">
                    +{new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(invoice.balance || 0)}
                  </div>
                </div>
              ))}
              {(!overdueInvoices || overdueInvoices.length === 0) && (
                <div className="text-center text-muted-foreground py-8">
                  No overdue invoices found.
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
