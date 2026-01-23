import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { 
  DollarSign, 
  TrendingDown, 
  AlertCircle, 
  Clock,
  Plus,
  Loader2,
  Truck,
  FileText,
  Receipt,
  CreditCard,
  Eye,
  Send,
  RefreshCw,
  Download,
  Filter,
  Mail,
  Phone,
  Printer,
  BarChart3,
  ArrowUpRight,
  ArrowDownRight,
  Search,
  CheckCircle,
  XCircle,
  Calendar,
  Wallet,
  Building2,
  Percent,
  CalendarClock
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  PieChart,
  Pie,
  Legend,
  LineChart,
  Line,
  Area,
  AreaChart
} from "recharts";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Separator } from "@/components/ui/separator";
import { toast } from "sonner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Progress } from "@/components/ui/progress";

// Sample data for demo
const sampleAgingData = {
  current: 32000000,
  days_1_30: 18000000,
  days_31_60: 12000000,
  days_61_90: 5000000,
  over_90: 3000000,
  total: 70000000
};

const sampleDueInvoices = [
  { id: 1, invoice_number: "BILL-2026-001", vendor_name: "Lao Fresh Produce", vendor_id: 1, balance: 8500000, total: 8500000, paid: 0, due_date: "2026-01-25", days_until_due: 4, early_discount: 2, discount_deadline: "2026-01-22", email: "fresh@laoproduce.com", phone: "+856 20 555 1001" },
  { id: 2, invoice_number: "BILL-2026-002", vendor_name: "Mekong Seafood Co.", vendor_id: 2, balance: 12750000, total: 12750000, paid: 0, due_date: "2026-01-26", days_until_due: 5, early_discount: 1.5, discount_deadline: "2026-01-23", email: "mekong@seafood.com", phone: "+856 20 555 1002" },
  { id: 3, invoice_number: "BILL-2026-003", vendor_name: "Thai Imports Ltd", vendor_id: 3, balance: 6200000, total: 10000000, paid: 3800000, due_date: "2026-01-28", days_until_due: 7, early_discount: 0, discount_deadline: null, email: "thai@imports.com", phone: "+856 20 555 1003" },
  { id: 4, invoice_number: "BILL-2026-004", vendor_name: "Vietnam Spices", vendor_id: 4, balance: 4800000, total: 4800000, paid: 0, due_date: "2026-01-30", days_until_due: 9, early_discount: 3, discount_deadline: "2026-01-25", email: "vietnam@spices.com", phone: "+856 20 555 1004" },
];

const sampleOverdueInvoices = [
  { id: 5, invoice_number: "BILL-2025-098", vendor_name: "China Packaging", vendor_id: 5, balance: 5500000, total: 8000000, paid: 2500000, due_date: "2026-01-15", days_overdue: 6, email: "china@packaging.com", phone: "+856 20 555 1005" },
  { id: 6, invoice_number: "BILL-2025-095", vendor_name: "Lao Equipment Co", vendor_id: 6, balance: 15200000, total: 15200000, paid: 0, due_date: "2026-01-10", days_overdue: 11, email: "lao@equipment.com", phone: "+856 20 555 1006" },
  { id: 7, invoice_number: "BILL-2025-088", vendor_name: "Quality Supplies", vendor_id: 7, balance: 3800000, total: 7600000, paid: 3800000, due_date: "2026-01-05", days_overdue: 16, email: "quality@supplies.com", phone: "+856 20 555 1007" },
];

const sampleAllInvoices = [
  ...sampleDueInvoices,
  ...sampleOverdueInvoices.map(i => ({ ...i, days_until_due: -i.days_overdue })),
  { id: 8, invoice_number: "BILL-2025-085", vendor_name: "Lao Fresh Produce", vendor_id: 1, balance: 0, total: 5500000, paid: 5500000, due_date: "2026-01-02", days_until_due: 0, status: "PAID", email: "fresh@laoproduce.com", phone: "+856 20 555 1001" },
  { id: 9, invoice_number: "BILL-2026-010", vendor_name: "Mekong Seafood Co.", vendor_id: 2, balance: 9800000, total: 9800000, paid: 0, due_date: "2026-02-05", days_until_due: 15, status: "PENDING", email: "mekong@seafood.com", phone: "+856 20 555 1002" },
];

const sampleRecentPayments = [
  { id: 1, payment_number: "PMT-2026-015", vendor_name: "Lao Fresh Produce", amount: 5500000, payment_date: "2026-01-20", method: "Bank Transfer", invoice: "BILL-2025-085", discount_taken: 110000 },
  { id: 2, payment_number: "PMT-2026-014", vendor_name: "Thai Imports Ltd", amount: 3800000, payment_date: "2026-01-19", method: "Check", invoice: "BILL-2026-003", discount_taken: 0 },
  { id: 3, payment_number: "PMT-2026-013", vendor_name: "Quality Supplies", amount: 3800000, payment_date: "2026-01-18", method: "Bank Transfer", invoice: "BILL-2025-088", discount_taken: 0 },
  { id: 4, payment_number: "PMT-2026-012", vendor_name: "China Packaging", amount: 2500000, payment_date: "2026-01-17", method: "Wire Transfer", invoice: "BILL-2025-098", discount_taken: 0 },
];

const sampleTopVendors = [
  { id: 1, name: "Lao Fresh Produce", balance: 14000000, ytd_purchases: 185000000, last_payment: "2026-01-20", payment_terms: "Net 30", status: "Good" },
  { id: 2, name: "Mekong Seafood Co.", balance: 22550000, ytd_purchases: 220000000, last_payment: "2026-01-15", payment_terms: "Net 15", status: "Warning" },
  { id: 3, name: "Thai Imports Ltd", balance: 6200000, ytd_purchases: 95000000, last_payment: "2026-01-19", payment_terms: "Net 30", status: "Good" },
  { id: 4, name: "Vietnam Spices", balance: 4800000, ytd_purchases: 65000000, last_payment: "2026-01-05", payment_terms: "2/10 Net 30", status: "Good" },
  { id: 5, name: "China Packaging", balance: 5500000, ytd_purchases: 45000000, last_payment: "2026-01-17", payment_terms: "Net 45", status: "Overdue" },
];

const sampleMonthlyTrends = [
  { month: "Aug", purchased: 55000000, paid: 52000000, outstanding: 35000000 },
  { month: "Sep", purchased: 62000000, paid: 58000000, outstanding: 39000000 },
  { month: "Oct", purchased: 58000000, paid: 60000000, outstanding: 37000000 },
  { month: "Nov", purchased: 65000000, paid: 55000000, outstanding: 47000000 },
  { month: "Dec", purchased: 78000000, paid: 65000000, outstanding: 60000000 },
  { month: "Jan", purchased: 45000000, paid: 35000000, outstanding: 70000000 },
];

const sampleVendors = [
  { id: 1, name: "Lao Fresh Produce" },
  { id: 2, name: "Mekong Seafood Co." },
  { id: 3, name: "Thai Imports Ltd" },
  { id: 4, name: "Vietnam Spices" },
  { id: 5, name: "China Packaging" },
  { id: 6, name: "Lao Equipment Co" },
  { id: 7, name: "Quality Supplies" },
];

const sampleScheduledPayments = [
  { id: 1, vendor_name: "Lao Fresh Produce", amount: 8500000, scheduled_date: "2026-01-25", invoices: ["BILL-2026-001"], status: "scheduled" },
  { id: 2, vendor_name: "Mekong Seafood Co.", amount: 12750000, scheduled_date: "2026-01-26", invoices: ["BILL-2026-002"], status: "scheduled" },
  { id: 3, vendor_name: "Vietnam Spices", amount: 4800000, scheduled_date: "2026-01-30", invoices: ["BILL-2026-004"], status: "pending_approval" },
];

export default function APDashboard() {
  // Dialog states
  const [paymentDialogOpen, setPaymentDialogOpen] = useState(false);
  const [billDialogOpen, setBillDialogOpen] = useState(false);
  const [scheduleDialogOpen, setScheduleDialogOpen] = useState(false);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  const [batchPayDialogOpen, setBatchPayDialogOpen] = useState(false);
  
  // Selected items
  const [selectedInvoice, setSelectedInvoice] = useState<any>(null);
  const [selectedVendor, setSelectedVendor] = useState<any>(null);
  const [selectedInvoices, setSelectedInvoices] = useState<number[]>([]);
  
  // Filter states
  const [vendorFilter, setVendorFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [searchTerm, setSearchTerm] = useState("");
  
  // Form states
  const [paymentAmount, setPaymentAmount] = useState("");
  const [paymentMethod, setPaymentMethod] = useState("");
  const [paymentReference, setPaymentReference] = useState("");
  const [paymentDate, setPaymentDate] = useState("");
  
  // Bill form states
  const [billVendor, setBillVendor] = useState("");
  const [billAmount, setBillAmount] = useState("");
  const [billDueDate, setBillDueDate] = useState("");
  const [billDescription, setBillDescription] = useState("");
  const [billInvoiceNumber, setBillInvoiceNumber] = useState("");
  
  // Loading states
  const [isLoading, setIsLoading] = useState(false);

  // Use sample data
  const aging = sampleAgingData;
  const dueInvoices = sampleDueInvoices;
  const overdueInvoices = sampleOverdueInvoices;
  const allInvoices = sampleAllInvoices;
  const recentPayments = sampleRecentPayments;
  const topVendors = sampleTopVendors;
  const monthlyTrends = sampleMonthlyTrends;
  const scheduledPayments = sampleScheduledPayments;

  // Filter invoices
  const filteredInvoices = allInvoices.filter(inv => {
    if (vendorFilter !== "all" && inv.vendor_id.toString() !== vendorFilter) return false;
    if (statusFilter === "overdue" && (inv.days_until_due || 0) >= 0) return false;
    if (statusFilter === "paid" && inv.balance > 0) return false;
    if (statusFilter === "pending" && (inv.balance === 0 || (inv.days_until_due || 0) < 0)) return false;
    if (searchTerm && !inv.invoice_number.toLowerCase().includes(searchTerm.toLowerCase()) && 
        !inv.vendor_name.toLowerCase().includes(searchTerm.toLowerCase())) return false;
    return true;
  });

  const chartData = [
    { name: "Current", value: aging.current, color: "#10b981" },
    { name: "1-30 Days", value: aging.days_1_30, color: "#3b82f6" },
    { name: "31-60 Days", value: aging.days_31_60, color: "#f59e0b" },
    { name: "61-90 Days", value: aging.days_61_90, color: "#f97316" },
    { name: "90+ Days", value: aging.over_90, color: "#ef4444" },
  ];

  const totalPayable = aging.total;
  const totalOverdue = aging.days_1_30 + aging.days_31_60 + aging.days_61_90 + aging.over_90;
  const dueThisWeek = dueInvoices.filter(i => i.days_until_due <= 7).reduce((sum, i) => sum + i.balance, 0);
  const potentialDiscounts = dueInvoices.filter(i => i.early_discount > 0).reduce((sum, i) => sum + (i.balance * i.early_discount / 100), 0);

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'LAK',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const handlePayBill = (invoice: any) => {
    setSelectedInvoice(invoice);
    setPaymentAmount(invoice.balance.toString());
    setPaymentDialogOpen(true);
  };

  const handleViewDetails = (invoice: any) => {
    setSelectedInvoice(invoice);
    setDetailsDialogOpen(true);
  };

  const handleSchedulePayment = (invoice: any) => {
    setSelectedInvoice(invoice);
    setPaymentAmount(invoice.balance.toString());
    setScheduleDialogOpen(true);
  };

  const handleBatchPay = () => {
    if (selectedInvoices.length === 0) {
      toast.error("Please select invoices to pay");
      return;
    }
    setBatchPayDialogOpen(true);
  };

  const toggleInvoiceSelection = (id: number) => {
    setSelectedInvoices(prev => 
      prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]
    );
  };

  const selectAllDue = () => {
    if (selectedInvoices.length === dueInvoices.length) {
      setSelectedInvoices([]);
    } else {
      setSelectedInvoices(dueInvoices.map(i => i.id));
    }
  };

  const selectedTotal = allInvoices
    .filter(i => selectedInvoices.includes(i.id))
    .reduce((sum, i) => sum + i.balance, 0);

  const submitPayment = () => {
    if (!paymentAmount || !paymentMethod) {
      toast.error("Please fill in all required fields");
      return;
    }
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Payment of ${formatCurrency(parseFloat(paymentAmount))} processed successfully`);
      setPaymentDialogOpen(false);
      resetPaymentForm();
      setIsLoading(false);
    }, 1000);
  };

  const submitBill = () => {
    if (!billVendor || !billAmount || !billDueDate) {
      toast.error("Please fill in all required fields");
      return;
    }
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Bill created for ${formatCurrency(parseFloat(billAmount))}`);
      setBillDialogOpen(false);
      resetBillForm();
      setIsLoading(false);
    }, 1000);
  };

  const submitSchedule = () => {
    if (!paymentAmount || !paymentDate) {
      toast.error("Please fill in all required fields");
      return;
    }
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Payment scheduled for ${paymentDate}`);
      setScheduleDialogOpen(false);
      resetPaymentForm();
      setIsLoading(false);
    }, 1000);
  };

  const submitBatchPayment = () => {
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Batch payment of ${formatCurrency(selectedTotal)} for ${selectedInvoices.length} invoices processed`);
      setBatchPayDialogOpen(false);
      setSelectedInvoices([]);
      setIsLoading(false);
    }, 1500);
  };

  const resetPaymentForm = () => {
    setPaymentAmount("");
    setPaymentMethod("");
    setPaymentReference("");
    setPaymentDate("");
    setSelectedInvoice(null);
  };

  const resetBillForm = () => {
    setBillVendor("");
    setBillAmount("");
    setBillDueDate("");
    setBillDescription("");
    setBillInvoiceNumber("");
  };

  const getStatusBadge = (invoice: any) => {
    if (invoice.balance === 0) {
      return <Badge className="bg-emerald-100 text-emerald-800">PAID</Badge>;
    }
    if ((invoice.days_until_due || 0) < 0 || invoice.days_overdue > 0) {
      const daysOver = invoice.days_overdue || Math.abs(invoice.days_until_due);
      if (daysOver > 30) {
        return <Badge variant="destructive">OVERDUE 30+</Badge>;
      }
      return <Badge className="bg-red-100 text-red-800">OVERDUE</Badge>;
    }
    if (invoice.days_until_due <= 7) {
      return <Badge className="bg-amber-100 text-amber-800">DUE SOON</Badge>;
    }
    if (invoice.paid > 0) {
      return <Badge className="bg-blue-100 text-blue-800">PARTIAL</Badge>;
    }
    return <Badge className="bg-gray-100 text-gray-800">PENDING</Badge>;
  };

  const getVendorStatusColor = (status: string) => {
    switch (status) {
      case "Good": return "text-emerald-600";
      case "Warning": return "text-amber-600";
      case "Overdue": return "text-red-600";
      default: return "";
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-orange-600 to-orange-400 bg-clip-text text-transparent">
            Accounts Payable
          </h1>
          <p className="text-muted-foreground">
            Manage vendor bills, payments, and cash flow
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" /> Export
          </Button>
          <Button variant="outline" size="sm">
            <Printer className="mr-2 h-4 w-4" /> Print Report
          </Button>
          <Button variant="outline" size="sm" onClick={handleBatchPay} disabled={selectedInvoices.length === 0}>
            <Wallet className="mr-2 h-4 w-4" /> Batch Pay ({selectedInvoices.length})
          </Button>
          <Button onClick={() => setBillDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" /> Enter Bill
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <Card className="border-l-4 border-l-orange-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Payable</CardTitle>
            <DollarSign className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalPayable)}</div>
            <div className="flex items-center text-xs text-amber-600 mt-1">
              <ArrowUpRight className="h-3 w-3 mr-1" />
              +8% from last month
            </div>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-red-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Overdue</CardTitle>
            <AlertCircle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-destructive">
              {formatCurrency(totalOverdue)}
            </div>
            <p className="text-xs text-muted-foreground">
              {((totalOverdue / totalPayable) * 100).toFixed(1)}% of total • {overdueInvoices.length} bills
            </p>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-amber-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Due This Week</CardTitle>
            <Clock className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">{formatCurrency(dueThisWeek)}</div>
            <p className="text-xs text-muted-foreground">
              {dueInvoices.filter(i => i.days_until_due <= 7).length} bills due
            </p>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-emerald-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Potential Discounts</CardTitle>
            <Percent className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">{formatCurrency(potentialDiscounts)}</div>
            <p className="text-xs text-muted-foreground">
              {dueInvoices.filter(i => i.early_discount > 0).length} bills with early pay discount
            </p>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-blue-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Scheduled Payments</CardTitle>
            <CalendarClock className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(scheduledPayments.reduce((s, p) => s + p.amount, 0))}</div>
            <p className="text-xs text-muted-foreground">
              {scheduledPayments.length} payments queued
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Charts Row */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* Aging Bar Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Aging Overview</CardTitle>
            <CardDescription>Payables by age bucket</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[250px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" horizontal={true} vertical={false} />
                  <XAxis type="number" tickFormatter={(value) => `₭${(value/1000000)}M`} fontSize={10} />
                  <YAxis dataKey="name" type="category" fontSize={10} width={70} />
                  <Tooltip formatter={(value: number) => [formatCurrency(value), 'Amount']} />
                  <Bar dataKey="value" radius={[0, 4, 4, 0]}>
                    {chartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        {/* Payment Trend */}
        <Card>
          <CardHeader>
            <CardTitle>Payment Trend</CardTitle>
            <CardDescription>Monthly purchases vs payments</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[250px]">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={monthlyTrends}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" fontSize={10} />
                  <YAxis tickFormatter={(value) => `${(value/1000000)}M`} fontSize={10} />
                  <Tooltip formatter={(value: number) => formatCurrency(value)} />
                  <Area type="monotone" dataKey="purchased" stackId="1" stroke="#f97316" fill="#f97316" fillOpacity={0.3} name="Purchased" />
                  <Area type="monotone" dataKey="paid" stackId="2" stroke="#10b981" fill="#10b981" fillOpacity={0.5} name="Paid" />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        {/* Outstanding Trend */}
        <Card>
          <CardHeader>
            <CardTitle>Outstanding Balance</CardTitle>
            <CardDescription>Monthly outstanding trend</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[250px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={monthlyTrends}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" fontSize={10} />
                  <YAxis tickFormatter={(value) => `${(value/1000000)}M`} fontSize={10} />
                  <Tooltip formatter={(value: number) => formatCurrency(value)} />
                  <Line type="monotone" dataKey="outstanding" stroke="#f97316" strokeWidth={2} dot={{ fill: '#f97316' }} name="Outstanding" />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <CardTitle>Bill Management</CardTitle>
            <div className="flex flex-wrap items-center gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search bills..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-9 w-[200px]"
                />
              </div>
              <Select value={vendorFilter} onValueChange={setVendorFilter}>
                <SelectTrigger className="w-[150px]">
                  <SelectValue placeholder="All Vendors" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Vendors</SelectItem>
                  {sampleVendors.map(v => (
                    <SelectItem key={v.id} value={v.id.toString()}>{v.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="All Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="overdue">Overdue</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="paid">Paid</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="due" className="space-y-4">
            <TabsList>
              <TabsTrigger value="due" className="gap-2">
                <Clock className="h-4 w-4" />
                Due Soon
                <Badge className="bg-amber-100 text-amber-800 ml-1">{dueInvoices.length}</Badge>
              </TabsTrigger>
              <TabsTrigger value="overdue" className="gap-2">
                <AlertCircle className="h-4 w-4" />
                Overdue
                <Badge variant="destructive" className="ml-1">{overdueInvoices.length}</Badge>
              </TabsTrigger>
              <TabsTrigger value="all" className="gap-2">
                <FileText className="h-4 w-4" />
                All Bills
                <Badge variant="secondary" className="ml-1">{allInvoices.length}</Badge>
              </TabsTrigger>
              <TabsTrigger value="payments" className="gap-2">
                <CreditCard className="h-4 w-4" />
                Payments
              </TabsTrigger>
              <TabsTrigger value="vendors" className="gap-2">
                <Truck className="h-4 w-4" />
                Vendors
              </TabsTrigger>
              <TabsTrigger value="scheduled" className="gap-2">
                <CalendarClock className="h-4 w-4" />
                Scheduled
              </TabsTrigger>
            </TabsList>

            <TabsContent value="due">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                  <Checkbox 
                    checked={selectedInvoices.length === dueInvoices.length && dueInvoices.length > 0}
                    onCheckedChange={selectAllDue}
                  />
                  <span className="text-sm text-muted-foreground">Select All</span>
                </div>
                <p className="text-sm text-muted-foreground">
                  {selectedInvoices.length} selected • Total: {formatCurrency(selectedTotal)}
                </p>
              </div>
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {dueInvoices.map((invoice) => (
                    <div key={invoice.id} className="flex items-center gap-4 p-4 rounded-lg border bg-card hover:bg-muted/50 transition-colors">
                      <Checkbox 
                        checked={selectedInvoices.includes(invoice.id)}
                        onCheckedChange={() => toggleInvoiceSelection(invoice.id)}
                      />
                      <div className="w-10 h-10 rounded-full bg-amber-100 dark:bg-amber-900 flex items-center justify-center shrink-0">
                        <Receipt className="h-5 w-5 text-amber-600" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <p className="font-medium">{invoice.vendor_name}</p>
                          {getStatusBadge(invoice)}
                          {invoice.early_discount > 0 && (
                            <Badge className="bg-emerald-100 text-emerald-800">
                              <Percent className="h-3 w-3 mr-1" />
                              {invoice.early_discount}% discount
                            </Badge>
                          )}
                        </div>
                        <p className="text-sm text-muted-foreground">{invoice.invoice_number}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-bold">{formatCurrency(invoice.balance)}</p>
                        <p className="text-sm text-amber-600">Due in {invoice.days_until_due} days</p>
                        {invoice.early_discount > 0 && invoice.discount_deadline && (
                          <p className="text-xs text-emerald-600">Discount until {invoice.discount_deadline}</p>
                        )}
                      </div>
                      <div className="flex gap-1">
                        <Button variant="ghost" size="icon" onClick={() => handleViewDetails(invoice)}>
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="icon" onClick={() => handleSchedulePayment(invoice)}>
                          <Calendar className="h-4 w-4" />
                        </Button>
                        <Button size="sm" onClick={() => handlePayBill(invoice)}>
                          <Wallet className="h-4 w-4 mr-1" /> Pay
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </TabsContent>

            <TabsContent value="overdue">
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {overdueInvoices.map((invoice) => (
                    <div key={invoice.id} className="flex items-center gap-4 p-4 rounded-lg border border-destructive/50 bg-destructive/5 hover:bg-destructive/10 transition-colors">
                      <div className="w-10 h-10 rounded-full bg-destructive/10 flex items-center justify-center shrink-0">
                        <AlertCircle className="h-5 w-5 text-destructive" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <p className="font-medium">{invoice.vendor_name}</p>
                          {getStatusBadge(invoice)}
                        </div>
                        <p className="text-sm text-muted-foreground">{invoice.invoice_number}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-bold text-destructive">{formatCurrency(invoice.balance)}</p>
                        <p className="text-sm text-destructive">{invoice.days_overdue} days overdue</p>
                      </div>
                      <div className="flex gap-1">
                        <Button variant="ghost" size="icon" onClick={() => handleViewDetails(invoice)}>
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button size="sm" variant="destructive" onClick={() => handlePayBill(invoice)}>
                          <Wallet className="h-4 w-4 mr-1" /> Pay Now
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </TabsContent>

            <TabsContent value="all">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Bill #</TableHead>
                    <TableHead>Vendor</TableHead>
                    <TableHead>Due Date</TableHead>
                    <TableHead className="text-right">Total</TableHead>
                    <TableHead className="text-right">Paid</TableHead>
                    <TableHead className="text-right">Balance</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredInvoices.map((invoice) => (
                    <TableRow key={invoice.id}>
                      <TableCell className="font-mono font-medium">{invoice.invoice_number}</TableCell>
                      <TableCell>{invoice.vendor_name}</TableCell>
                      <TableCell>{invoice.due_date}</TableCell>
                      <TableCell className="text-right">{formatCurrency(invoice.total)}</TableCell>
                      <TableCell className="text-right text-emerald-600">{formatCurrency(invoice.paid)}</TableCell>
                      <TableCell className="text-right font-bold">{formatCurrency(invoice.balance)}</TableCell>
                      <TableCell>{getStatusBadge(invoice)}</TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-1">
                          <Button variant="ghost" size="icon" onClick={() => handleViewDetails(invoice)}>
                            <Eye className="h-4 w-4" />
                          </Button>
                          {invoice.balance > 0 && (
                            <Button variant="ghost" size="icon" onClick={() => handlePayBill(invoice)}>
                              <Wallet className="h-4 w-4" />
                            </Button>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TabsContent>

            <TabsContent value="payments">
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {recentPayments.map((payment) => (
                    <div key={payment.id} className="flex items-center justify-between p-4 rounded-lg border bg-card">
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 rounded-full bg-emerald-100 dark:bg-emerald-900 flex items-center justify-center">
                          <CreditCard className="h-5 w-5 text-emerald-600" />
                        </div>
                        <div>
                          <p className="font-medium">{payment.vendor_name}</p>
                          <p className="text-sm text-muted-foreground">{payment.payment_number} • {payment.invoice}</p>
                        </div>
                      </div>
                      <Badge variant="outline">{payment.method}</Badge>
                      {payment.discount_taken > 0 && (
                        <Badge className="bg-emerald-100 text-emerald-800">
                          <Percent className="h-3 w-3 mr-1" />
                          Saved {formatCurrency(payment.discount_taken)}
                        </Badge>
                      )}
                      <div className="text-right">
                        <p className="font-bold text-emerald-600">{formatCurrency(payment.amount)}</p>
                        <p className="text-sm text-muted-foreground">{payment.payment_date}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </TabsContent>

            <TabsContent value="vendors">
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {topVendors.map((vendor, index) => (
                    <div key={vendor.id} className="flex items-center justify-between p-4 rounded-lg border bg-card">
                      <div className="flex items-center gap-4">
                        <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center font-bold text-primary">
                          {index + 1}
                        </div>
                        <div>
                          <div className="flex items-center gap-2">
                            <p className="font-medium">{vendor.name}</p>
                            <Badge className={`${getVendorStatusColor(vendor.status)} bg-transparent border`}>
                              {vendor.status}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground">
                            {vendor.payment_terms} • Last payment: {vendor.last_payment}
                          </p>
                        </div>
                      </div>
                      <div className="text-center">
                        <p className="text-sm text-muted-foreground">YTD Purchases</p>
                        <p className="font-medium">{formatCurrency(vendor.ytd_purchases)}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-muted-foreground">Current Balance</p>
                        <p className="font-bold text-orange-600">{formatCurrency(vendor.balance)}</p>
                      </div>
                      <Button variant="outline" size="sm">
                        View Statement
                      </Button>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </TabsContent>

            <TabsContent value="scheduled">
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {scheduledPayments.map((payment) => (
                    <div key={payment.id} className="flex items-center justify-between p-4 rounded-lg border bg-card">
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
                          <CalendarClock className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <p className="font-medium">{payment.vendor_name}</p>
                          <p className="text-sm text-muted-foreground">
                            {payment.invoices.join(", ")}
                          </p>
                        </div>
                      </div>
                      <Badge variant={payment.status === "scheduled" ? "secondary" : "outline"}>
                        {payment.status === "scheduled" ? "Scheduled" : "Pending Approval"}
                      </Badge>
                      <div className="text-right">
                        <p className="font-bold">{formatCurrency(payment.amount)}</p>
                        <p className="text-sm text-muted-foreground">{payment.scheduled_date}</p>
                      </div>
                      <div className="flex gap-1">
                        <Button variant="ghost" size="icon">
                          <Eye className="h-4 w-4" />
                        </Button>
                        {payment.status === "pending_approval" && (
                          <Button size="sm" variant="outline">
                            <CheckCircle className="h-4 w-4 mr-1" /> Approve
                          </Button>
                        )}
                        <Button size="sm" variant="ghost">
                          <XCircle className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* Payment Dialog */}
      <Dialog open={paymentDialogOpen} onOpenChange={setPaymentDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Pay Bill</DialogTitle>
            <DialogDescription>
              {selectedInvoice && (
                <span>
                  Paying <strong>{selectedInvoice.invoice_number}</strong>
                  <br />
                  Vendor: {selectedInvoice.vendor_name} • Balance: {formatCurrency(selectedInvoice.balance)}
                  {selectedInvoice.early_discount > 0 && (
                    <span className="text-emerald-600 block mt-1">
                      {selectedInvoice.early_discount}% early payment discount available!
                    </span>
                  )}
                </span>
              )}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="amount" className="text-right">Amount *</Label>
              <Input
                id="amount"
                type="number"
                value={paymentAmount}
                onChange={(e) => setPaymentAmount(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="method" className="text-right">Method *</Label>
              <Select value={paymentMethod} onValueChange={setPaymentMethod}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select payment method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                  <SelectItem value="check">Check</SelectItem>
                  <SelectItem value="wire">Wire Transfer</SelectItem>
                  <SelectItem value="cash">Cash</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="reference" className="text-right">Reference</Label>
              <Input
                id="reference"
                placeholder="Check #, Transaction ID, etc."
                value={paymentReference}
                onChange={(e) => setPaymentReference(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPaymentDialogOpen(false)}>Cancel</Button>
            <Button onClick={submitPayment} disabled={isLoading}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              Process Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Bill Dialog */}
      <Dialog open={billDialogOpen} onOpenChange={setBillDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Enter New Bill</DialogTitle>
            <DialogDescription>Record a new vendor bill</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Vendor *</Label>
              <Select value={billVendor} onValueChange={setBillVendor}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select vendor" />
                </SelectTrigger>
                <SelectContent>
                  {sampleVendors.map(v => (
                    <SelectItem key={v.id} value={v.id.toString()}>{v.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Bill # *</Label>
              <Input
                placeholder="Vendor invoice number"
                value={billInvoiceNumber}
                onChange={(e) => setBillInvoiceNumber(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Amount *</Label>
              <Input
                type="number"
                placeholder="0.00"
                value={billAmount}
                onChange={(e) => setBillAmount(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Due Date *</Label>
              <Input
                type="date"
                value={billDueDate}
                onChange={(e) => setBillDueDate(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Description</Label>
              <Textarea
                placeholder="Bill description..."
                value={billDescription}
                onChange={(e) => setBillDescription(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setBillDialogOpen(false)}>Cancel</Button>
            <Button onClick={submitBill} disabled={isLoading}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              Save Bill
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Schedule Payment Dialog */}
      <Dialog open={scheduleDialogOpen} onOpenChange={setScheduleDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Schedule Payment</DialogTitle>
            <DialogDescription>
              {selectedInvoice && (
                <span>
                  Schedule payment for <strong>{selectedInvoice.invoice_number}</strong>
                </span>
              )}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Amount *</Label>
              <Input
                type="number"
                value={paymentAmount}
                onChange={(e) => setPaymentAmount(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Pay Date *</Label>
              <Input
                type="date"
                value={paymentDate}
                onChange={(e) => setPaymentDate(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Method</Label>
              <Select value={paymentMethod} onValueChange={setPaymentMethod}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select payment method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                  <SelectItem value="check">Check</SelectItem>
                  <SelectItem value="wire">Wire Transfer</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setScheduleDialogOpen(false)}>Cancel</Button>
            <Button onClick={submitSchedule} disabled={isLoading}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              <CalendarClock className="h-4 w-4 mr-2" /> Schedule Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Details Dialog */}
      <Dialog open={detailsDialogOpen} onOpenChange={setDetailsDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>Bill Details</DialogTitle>
          </DialogHeader>
          {selectedInvoice && (
            <div className="space-y-6">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="text-2xl font-bold">{selectedInvoice.invoice_number}</h3>
                  <p className="text-muted-foreground">{selectedInvoice.vendor_name}</p>
                </div>
                {getStatusBadge(selectedInvoice)}
              </div>
              <Separator />
              <div className="grid grid-cols-2 gap-6">
                <div>
                  <Label className="text-muted-foreground">Bill Total</Label>
                  <p className="text-xl font-bold">{formatCurrency(selectedInvoice.total)}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Amount Paid</Label>
                  <p className="text-xl font-bold text-emerald-600">{formatCurrency(selectedInvoice.paid)}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Balance Due</Label>
                  <p className="text-xl font-bold text-orange-600">{formatCurrency(selectedInvoice.balance)}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Due Date</Label>
                  <p className="text-xl font-bold">{selectedInvoice.due_date}</p>
                </div>
              </div>
              {selectedInvoice.early_discount > 0 && (
                <>
                  <Separator />
                  <div className="bg-emerald-50 dark:bg-emerald-950 p-4 rounded-lg">
                    <div className="flex items-center gap-2 mb-2">
                      <Percent className="h-5 w-5 text-emerald-600" />
                      <span className="font-medium text-emerald-700 dark:text-emerald-400">Early Payment Discount Available</span>
                    </div>
                    <p className="text-sm text-emerald-600">
                      Pay by <strong>{selectedInvoice.discount_deadline}</strong> to save <strong>{selectedInvoice.early_discount}%</strong> ({formatCurrency(selectedInvoice.balance * selectedInvoice.early_discount / 100)})
                    </p>
                  </div>
                </>
              )}
              <Separator />
              <div>
                <Label className="text-muted-foreground">Contact Information</Label>
                <div className="mt-2 space-y-2">
                  <div className="flex items-center gap-2">
                    <Mail className="h-4 w-4 text-muted-foreground" />
                    <span>{selectedInvoice.email}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Phone className="h-4 w-4 text-muted-foreground" />
                    <span>{selectedInvoice.phone}</span>
                  </div>
                </div>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setDetailsDialogOpen(false)}>Close</Button>
            {selectedInvoice?.balance > 0 && (
              <>
                <Button variant="outline" onClick={() => { setDetailsDialogOpen(false); handleSchedulePayment(selectedInvoice); }}>
                  <Calendar className="h-4 w-4 mr-2" /> Schedule
                </Button>
                <Button onClick={() => { setDetailsDialogOpen(false); handlePayBill(selectedInvoice); }}>
                  <Wallet className="h-4 w-4 mr-2" /> Pay Now
                </Button>
              </>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Batch Payment Dialog */}
      <Dialog open={batchPayDialogOpen} onOpenChange={setBatchPayDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>Batch Payment</DialogTitle>
            <DialogDescription>
              Process payment for {selectedInvoices.length} selected bills
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="bg-muted p-4 rounded-lg">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total Amount</span>
                <span className="text-2xl font-bold">{formatCurrency(selectedTotal)}</span>
              </div>
            </div>
            <ScrollArea className="h-[200px]">
              <div className="space-y-2">
                {allInvoices.filter(i => selectedInvoices.includes(i.id)).map(inv => (
                  <div key={inv.id} className="flex justify-between items-center p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">{inv.vendor_name}</p>
                      <p className="text-sm text-muted-foreground">{inv.invoice_number}</p>
                    </div>
                    <p className="font-bold">{formatCurrency(inv.balance)}</p>
                  </div>
                ))}
              </div>
            </ScrollArea>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Method *</Label>
              <Select value={paymentMethod} onValueChange={setPaymentMethod}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select payment method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                  <SelectItem value="wire">Wire Transfer</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setBatchPayDialogOpen(false)}>Cancel</Button>
            <Button onClick={submitBatchPayment} disabled={isLoading || !paymentMethod}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              Process {formatCurrency(selectedTotal)}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
