import { useState } from "react";
import { Link } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  BookOpen, 
  FileText, 
  TrendingUp, 
  DollarSign,
  Plus,
  Loader2,
  Calendar,
  BarChart3,
  PieChart,
  ArrowUpRight,
  ArrowDownRight
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart as RechartsPieChart,
  Pie,
  Cell,
  Legend
} from "recharts";
import { financialService } from "@/services/financialService";

const COLORS = ['#10b981', '#3b82f6', '#f59e0b', '#ef4444', '#8b5cf6'];

export default function GLDashboard() {
  const [selectedPeriod, setSelectedPeriod] = useState("current");

  const { data: currentPeriod, isLoading: isPeriodLoading } = useQuery({
    queryKey: ['currentPeriod'],
    queryFn: financialService.getCurrentPeriod
  });

  const { data: trialBalance, isLoading: isTBLoading } = useQuery({
    queryKey: ['trialBalance'],
    queryFn: () => financialService.getTrialBalance()
  });

  const { data: incomeStatement, isLoading: isISLoading } = useQuery({
    queryKey: ['incomeStatement'],
    queryFn: () => financialService.getIncomeStatement()
  });

  const isLoading = isPeriodLoading || isTBLoading || isISLoading;

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  // Mock summary data (replace with real calculations)
  const summaryData = {
    totalAssets: 1250000,
    totalLiabilities: 450000,
    totalEquity: 800000,
    totalRevenue: 520000,
    totalExpenses: 380000,
    netIncome: 140000
  };

  const accountTypeData = [
    { name: 'Assets', value: summaryData.totalAssets, color: '#10b981' },
    { name: 'Liabilities', value: summaryData.totalLiabilities, color: '#ef4444' },
    { name: 'Equity', value: summaryData.totalEquity, color: '#3b82f6' },
  ];

  const monthlyData = [
    { month: 'Jul', revenue: 45000, expenses: 32000 },
    { month: 'Aug', revenue: 52000, expenses: 38000 },
    { month: 'Sep', revenue: 48000, expenses: 35000 },
    { month: 'Oct', revenue: 61000, expenses: 42000 },
    { month: 'Nov', revenue: 55000, expenses: 39000 },
    { month: 'Dec', revenue: 58000, expenses: 41000 },
  ];

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">General Ledger</h1>
          <p className="text-muted-foreground">
            Financial overview and accounting management
          </p>
        </div>
        <div className="flex gap-2">
          <Link href="/gl/journal-entries">
            <Button variant="outline">
              <FileText className="mr-2 h-4 w-4" /> Journal Entries
            </Button>
          </Link>
          <Link href="/gl/journal-entries/new">
            <Button>
              <Plus className="mr-2 h-4 w-4" /> New Entry
            </Button>
          </Link>
        </div>
      </div>

      {/* Current Period Banner */}
      {currentPeriod && (
        <Card className="bg-primary/5 border-primary/20">
          <CardContent className="flex items-center justify-between py-4">
            <div className="flex items-center gap-3">
              <Calendar className="h-5 w-5 text-primary" />
              <div>
                <p className="font-medium">Current Period: {currentPeriod.name || 'Period 1'}</p>
                <p className="text-sm text-muted-foreground">
                  {currentPeriod.start_date ? new Date(currentPeriod.start_date).toLocaleDateString() : 'N/A'} - 
                  {currentPeriod.end_date ? new Date(currentPeriod.end_date).toLocaleDateString() : 'N/A'}
                </p>
              </div>
            </div>
            <Badge variant={currentPeriod.status === 'OPEN' ? 'default' : 'secondary'}>
              {currentPeriod.status || 'OPEN'}
            </Badge>
          </CardContent>
        </Card>
      )}

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Assets</CardTitle>
            <TrendingUp className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">
              {formatCurrency(summaryData.totalAssets)}
            </div>
            <p className="text-xs text-muted-foreground flex items-center mt-1">
              <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
              <span className="text-emerald-500">+5.2%</span> from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Liabilities</CardTitle>
            <DollarSign className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {formatCurrency(summaryData.totalLiabilities)}
            </div>
            <p className="text-xs text-muted-foreground flex items-center mt-1">
              <ArrowDownRight className="h-3 w-3 text-emerald-500 mr-1" />
              <span className="text-emerald-500">-2.1%</span> from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Equity</CardTitle>
            <BookOpen className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {formatCurrency(summaryData.totalEquity)}
            </div>
            <p className="text-xs text-muted-foreground flex items-center mt-1">
              <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
              <span className="text-emerald-500">+8.4%</span> from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Net Income</CardTitle>
            <BarChart3 className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">
              {formatCurrency(summaryData.netIncome)}
            </div>
            <p className="text-xs text-muted-foreground flex items-center mt-1">
              <ArrowUpRight className="h-3 w-3 text-emerald-500 mr-1" />
              <span className="text-emerald-500">+12.3%</span> from last period
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Charts */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Revenue vs Expenses</CardTitle>
            <CardDescription>Monthly comparison for current fiscal year</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={monthlyData}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="month" stroke="#888888" fontSize={12} tickLine={false} axisLine={false} />
                  <YAxis stroke="#888888" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(v) => `$${v/1000}k`} />
                  <Tooltip formatter={(value: number) => formatCurrency(value)} />
                  <Legend />
                  <Bar dataKey="revenue" fill="#10b981" name="Revenue" radius={[4, 4, 0, 0]} />
                  <Bar dataKey="expenses" fill="#ef4444" name="Expenses" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Balance Sheet Overview</CardTitle>
            <CardDescription>Assets, Liabilities & Equity distribution</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <RechartsPieChart>
                  <Pie
                    data={accountTypeData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={100}
                    paddingAngle={5}
                    dataKey="value"
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  >
                    {accountTypeData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value: number) => formatCurrency(value)} />
                </RechartsPieChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Links */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Link href="/gl/accounts">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-blue-100 rounded-lg">
                <BookOpen className="h-6 w-6 text-blue-600" />
              </div>
              <div>
                <p className="font-medium">Chart of Accounts</p>
                <p className="text-sm text-muted-foreground">Manage GL accounts</p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/gl/reports/trial-balance">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-emerald-100 rounded-lg">
                <BarChart3 className="h-6 w-6 text-emerald-600" />
              </div>
              <div>
                <p className="font-medium">Trial Balance</p>
                <p className="text-sm text-muted-foreground">View account balances</p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/gl/reports/income-statement">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-purple-100 rounded-lg">
                <TrendingUp className="h-6 w-6 text-purple-600" />
              </div>
              <div>
                <p className="font-medium">Income Statement</p>
                <p className="text-sm text-muted-foreground">Profit & Loss report</p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/gl/reports/balance-sheet">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-amber-100 rounded-lg">
                <PieChart className="h-6 w-6 text-amber-600" />
              </div>
              <div>
                <p className="font-medium">Balance Sheet</p>
                <p className="text-sm text-muted-foreground">Financial position</p>
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  );
}
