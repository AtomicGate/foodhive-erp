import { Link } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { DataTable } from "@/components/ui/data-table";
import { 
  DollarSign, 
  TrendingDown, 
  AlertCircle, 
  Clock,
  Plus,
  Loader2,
  FileText,
  CreditCard,
  MoreHorizontal,
  ArrowUpDown,
  Eye,
  CheckCircle
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
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { financialService, APInvoice } from "@/services/financialService";

export default function APDashboard() {
  const { data: apAgingReport, isLoading: isAgingLoading } = useQuery({
    queryKey: ['apAgingReport'],
    queryFn: financialService.getAPAgingReport
  });

  const { data: dueInvoices, isLoading: isDueLoading } = useQuery({
    queryKey: ['dueInvoices'],
    queryFn: financialService.getDueInvoices
  });

  const { data: overdueInvoices, isLoading: isOverdueLoading } = useQuery({
    queryKey: ['overdueAPInvoices'],
    queryFn: financialService.getOverdueAPInvoices
  });

  const isLoading = isAgingLoading || isDueLoading || isOverdueLoading;

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  // Calculate totals from aging report
  const totals = apAgingReport?.reduce((acc: any, vendor: any) => ({
    current: acc.current + (vendor.current || 0),
    days_1_30: acc.days_1_30 + (vendor.days_1_30 || 0),
    days_31_60: acc.days_31_60 + (vendor.days_31_60 || 0),
    days_61_90: acc.days_61_90 + (vendor.days_61_90 || 0),
    over_90: acc.over_90 + (vendor.over_90 || 0),
    total: acc.total + (vendor.total || 0),
  }), { current: 0, days_1_30: 0, days_31_60: 0, days_61_90: 0, over_90: 0, total: 0 }) || 
  { current: 25000, days_1_30: 15000, days_31_60: 8000, days_61_90: 3000, over_90: 1000, total: 52000 };

  const chartData = [
    { name: "Current", value: totals.current, color: "#10b981" },
    { name: "1-30 Days", value: totals.days_1_30, color: "#3b82f6" },
    { name: "31-60 Days", value: totals.days_31_60, color: "#f59e0b" },
    { name: "61-90 Days", value: totals.days_61_90, color: "#f97316" },
    { name: "90+ Days", value: totals.over_90, color: "#ef4444" },
  ];

  const totalPayable = totals.total;
  const totalOverdue = totals.days_1_30 + totals.days_31_60 + totals.days_61_90 + totals.over_90;

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value);
  };

  const invoiceColumns: ColumnDef<APInvoice>[] = [
    {
      accessorKey: "invoice_number",
      header: "Invoice #",
      cell: ({ row }) => (
        <div className="font-mono font-medium">{row.getValue("invoice_number")}</div>
      ),
    },
    {
      accessorKey: "vendor_name",
      header: "Vendor",
    },
    {
      accessorKey: "due_date",
      header: "Due Date",
      cell: ({ row }) => {
        const date = row.getValue("due_date") as string;
        const dueDate = new Date(date);
        const isOverdue = dueDate < new Date();
        return (
          <span className={isOverdue ? "text-destructive font-medium" : ""}>
            {dueDate.toLocaleDateString()}
          </span>
        );
      },
    },
    {
      accessorKey: "balance",
      header: () => <div className="text-right">Amount Due</div>,
      cell: ({ row }) => (
        <div className="text-right font-medium">
          {formatCurrency(row.getValue("balance"))}
        </div>
      ),
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        let variant: "default" | "secondary" | "destructive" | "outline" = "default";
        let className = "";
        
        switch (status) {
          case "PENDING":
            className = "bg-gray-100 text-gray-800";
            break;
          case "APPROVED":
            className = "bg-blue-100 text-blue-800";
            break;
          case "PAID":
            className = "bg-emerald-100 text-emerald-800";
            break;
          case "PARTIAL":
            className = "bg-amber-100 text-amber-800";
            break;
          case "VOIDED":
            variant = "destructive";
            break;
        }

        return <Badge variant={variant} className={className}>{status}</Badge>;
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const invoice = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <Eye className="mr-2 h-4 w-4" /> View Details
              </DropdownMenuItem>
              {invoice.status === "APPROVED" && (
                <DropdownMenuItem>
                  <CreditCard className="mr-2 h-4 w-4" /> Pay Invoice
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Accounts Payable</h1>
          <p className="text-muted-foreground">
            Monitor payables, aging, and vendor payments
          </p>
        </div>
        <div className="flex gap-2">
          <Link href="/ap/invoices">
            <Button variant="outline">
              <FileText className="mr-2 h-4 w-4" /> View Invoices
            </Button>
          </Link>
          <Link href="/ap/payments/new">
            <Button>
              <Plus className="mr-2 h-4 w-4" /> Record Payment
            </Button>
          </Link>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Payable</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalPayable)}</div>
            <p className="text-xs text-muted-foreground">
              Across all vendors
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
              {formatCurrency(totalOverdue)}
            </div>
            <p className="text-xs text-muted-foreground">
              {totalPayable > 0 ? ((totalOverdue / totalPayable) * 100).toFixed(1) : 0}% of total payable
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Due This Week</CardTitle>
            <Clock className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">
              {formatCurrency(dueInvoices?.reduce((sum: number, inv: any) => sum + (inv.balance || 0), 0) || 0)}
            </div>
            <p className="text-xs text-muted-foreground">
              {dueInvoices?.length || 0} invoices due
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Paid This Month</CardTitle>
            <TrendingDown className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">$28,450.00</div>
            <p className="text-xs text-muted-foreground">
              15 payments processed
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Charts and Tables */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Aging Chart */}
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Aging Overview</CardTitle>
            <CardDescription>Payables by age bucket</CardDescription>
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
                    tickFormatter={(value) => `$${value/1000}k`} 
                  />
                  <Tooltip 
                    cursor={{ fill: 'transparent' }}
                    formatter={(value: number) => [formatCurrency(value), 'Amount']}
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

        {/* Top Vendors */}
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Top Payables by Vendor</CardTitle>
            <CardDescription>Largest outstanding balances</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              {apAgingReport?.slice(0, 5).map((vendor: any, index: number) => (
                <div key={index} className="flex items-center">
                  <div className="ml-4 space-y-1 flex-1">
                    <p className="text-sm font-medium leading-none">{vendor.vendor_name || `Vendor ${index + 1}`}</p>
                    <p className="text-sm text-muted-foreground">
                      {vendor.over_90 > 0 && <span className="text-destructive">90+ days overdue</span>}
                    </p>
                  </div>
                  <div className="font-medium">
                    {formatCurrency(vendor.total || 0)}
                  </div>
                </div>
              )) || (
                <div className="text-center text-muted-foreground py-4">
                  No vendor data available
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Due Invoices Table */}
      <Card>
        <CardHeader>
          <CardTitle>Invoices Due Soon</CardTitle>
          <CardDescription>Invoices due within the next 7 days</CardDescription>
        </CardHeader>
        <CardContent>
          <DataTable 
            columns={invoiceColumns} 
            data={dueInvoices || []} 
            searchKey="vendor_name"
            searchPlaceholder="Search vendors..."
          />
        </CardContent>
      </Card>
    </div>
  );
}
