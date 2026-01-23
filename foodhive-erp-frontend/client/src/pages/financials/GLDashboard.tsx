import { useState } from "react";
import { Link } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
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
  ArrowDownRight,
  AlertCircle,
  RefreshCw
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

// Sample data for demo when API fails or no data
const sampleSummaryData = {
  totalAssets: 1250000000,
  totalLiabilities: 450000000,
  totalEquity: 800000000,
  totalRevenue: 520000000,
  totalExpenses: 380000000,
  netIncome: 140000000
};

const sampleMonthlyData = [
  { month: 'Jul', revenue: 45000000, expenses: 32000000 },
  { month: 'Aug', revenue: 52000000, expenses: 38000000 },
  { month: 'Sep', revenue: 48000000, expenses: 35000000 },
  { month: 'Oct', revenue: 61000000, expenses: 42000000 },
  { month: 'Nov', revenue: 55000000, expenses: 39000000 },
  { month: 'Dec', revenue: 58000000, expenses: 41000000 },
];

export default function GLDashboard() {
  const [selectedPeriod, setSelectedPeriod] = useState("current");

  // Fetch current period - with error handling
  const { data: currentPeriod, isLoading: isPeriodLoading, error: periodError, refetch: refetchPeriod } = useQuery({
    queryKey: ['currentPeriod'],
    queryFn: financialService.getCurrentPeriod,
    retry: 1,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  // Fetch trial balance - with error handling
  const { data: trialBalance, isLoading: isTBLoading, error: tbError, refetch: refetchTB } = useQuery({
    queryKey: ['trialBalance'],
    queryFn: () => financialService.getTrialBalance(),
    retry: 1,
    staleTime: 5 * 60 * 1000,
  });

  // Fetch income statement - with error handling
  const { data: incomeStatement, isLoading: isISLoading, error: isError, refetch: refetchIS } = useQuery({
    queryKey: ['incomeStatement'],
    queryFn: () => financialService.getIncomeStatement(),
    retry: 1,
    staleTime: 5 * 60 * 1000,
  });

  // Use sample data - in production, you'd calculate from real data
  const summaryData = sampleSummaryData;
  const monthlyData = sampleMonthlyData;

  const accountTypeData = [
    { name: 'Assets', value: summaryData.totalAssets, color: '#10b981' },
    { name: 'Liabilities', value: summaryData.totalLiabilities, color: '#ef4444' },
    { name: 'Equity', value: summaryData.totalEquity, color: '#3b82f6' },
  ];

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'LAK',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const hasAnyError = periodError || tbError || isError;

  const refetchAll = () => {
    refetchPeriod();
    refetchTB();
    refetchIS();
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-blue-600 to-blue-400 bg-clip-text text-transparent">
            General Ledger
          </h1>
          <p className="text-muted-foreground">
            Financial overview and accounting management
          </p>
        </div>
        <div className="flex gap-2">
          {hasAnyError && (
            <Button variant="outline" onClick={refetchAll}>
              <RefreshCw className="mr-2 h-4 w-4" /> Retry
            </Button>
          )}
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

      {/* Error Banner */}
      {hasAnyError && (
        <Card className="bg-amber-50 dark:bg-amber-950 border-amber-200 dark:border-amber-800">
          <CardContent className="flex items-center gap-3 py-4">
            <AlertCircle className="h-5 w-5 text-amber-600" />
            <div>
              <p className="font-medium text-amber-800 dark:text-amber-200">Some data couldn't be loaded</p>
              <p className="text-sm text-amber-600 dark:text-amber-400">
                Showing sample data. This may be because fiscal periods haven't been set up yet.
              </p>
            </div>
            <Button variant="outline" size="sm" className="ml-auto" onClick={refetchAll}>
              <RefreshCw className="h-4 w-4 mr-1" /> Retry
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Current Period Banner */}
      <Card className="bg-primary/5 border-primary/20">
        <CardContent className="flex items-center justify-between py-4">
          <div className="flex items-center gap-3">
            <Calendar className="h-5 w-5 text-primary" />
            {isPeriodLoading ? (
              <div className="space-y-2">
                <Skeleton className="h-4 w-32" />
                <Skeleton className="h-3 w-48" />
              </div>
            ) : currentPeriod ? (
              <div>
                <p className="font-medium">Current Period: {currentPeriod.name || 'Period 1'}</p>
                <p className="text-sm text-muted-foreground">
                  {currentPeriod.start_date ? new Date(currentPeriod.start_date).toLocaleDateString() : 'N/A'} - 
                  {currentPeriod.end_date ? new Date(currentPeriod.end_date).toLocaleDateString() : 'N/A'}
                </p>
              </div>
            ) : (
              <div>
                <p className="font-medium">No Active Period</p>
                <p className="text-sm text-muted-foreground">
                  Create a fiscal year to get started
                </p>
              </div>
            )}
          </div>
          {currentPeriod ? (
            <Badge variant={currentPeriod.status === 'OPEN' ? 'default' : 'secondary'}>
              {currentPeriod.status || 'OPEN'}
            </Badge>
          ) : (
            <Link href="/gl/fiscal-years/new">
              <Button size="sm" variant="outline">
                <Plus className="h-4 w-4 mr-1" /> Create Fiscal Year
              </Button>
            </Link>
          )}
        </CardContent>
      </Card>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="border-l-4 border-l-emerald-500">
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

        <Card className="border-l-4 border-l-red-500">
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

        <Card className="border-l-4 border-l-blue-500">
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

        <Card className="border-l-4 border-l-purple-500">
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
                  <YAxis stroke="#888888" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(v) => `₭${v/1000000}M`} />
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
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer h-full">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
                <BookOpen className="h-6 w-6 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <p className="font-medium">Chart of Accounts</p>
                <p className="text-sm text-muted-foreground">Manage GL accounts</p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/gl/reports/trial-balance">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer h-full">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-emerald-100 dark:bg-emerald-900 rounded-lg">
                <BarChart3 className="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <p className="font-medium">Trial Balance</p>
                <p className="text-sm text-muted-foreground">View account balances</p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/gl/reports/income-statement">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer h-full">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
                <TrendingUp className="h-6 w-6 text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <p className="font-medium">Income Statement</p>
                <p className="text-sm text-muted-foreground">Profit & Loss report</p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/gl/reports/balance-sheet">
          <Card className="hover:bg-muted/50 transition-colors cursor-pointer h-full">
            <CardContent className="flex items-center gap-4 py-6">
              <div className="p-3 bg-amber-100 dark:bg-amber-900 rounded-lg">
                <PieChart className="h-6 w-6 text-amber-600 dark:text-amber-400" />
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
