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
  TrendingUp, 
  AlertCircle, 
  Clock,
  Plus,
  Loader2,
  Users,
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
  Calendar
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

// Sample data for demo
const sampleAgingData = {
  current: 45000000,
  days_1_30: 28000000,
  days_31_60: 15000000,
  days_61_90: 8000000,
  over_90: 4000000,
  total: 100000000
};

const sampleOverdueInvoices = [
  { id: 1, invoice_number: "INV-2026-001", customer_name: "ABC Restaurant", customer_id: 1, balance: 12500000, total: 15000000, paid: 2500000, due_date: "2026-01-10", days_overdue: 11, email: "abc@restaurant.com", phone: "+856 20 555 0001" },
  { id: 2, invoice_number: "INV-2026-002", customer_name: "XYZ Food Court", customer_id: 2, balance: 8750000, total: 8750000, paid: 0, due_date: "2026-01-05", days_overdue: 16, email: "xyz@foodcourt.com", phone: "+856 20 555 0002" },
  { id: 3, invoice_number: "INV-2025-098", customer_name: "Fresh Mart", customer_id: 3, balance: 6200000, total: 10000000, paid: 3800000, due_date: "2025-12-28", days_overdue: 24, email: "fresh@mart.com", phone: "+856 20 555 0003" },
  { id: 4, invoice_number: "INV-2025-095", customer_name: "Golden Wok", customer_id: 4, balance: 4800000, total: 4800000, paid: 0, due_date: "2025-12-20", days_overdue: 32, email: "golden@wok.com", phone: "+856 20 555 0004" },
  { id: 5, invoice_number: "INV-2025-088", customer_name: "Ocean Seafood", customer_id: 5, balance: 3500000, total: 7000000, paid: 3500000, due_date: "2025-12-15", days_overdue: 37, email: "ocean@seafood.com", phone: "+856 20 555 0005" },
];

const sampleAllInvoices = [
  ...sampleOverdueInvoices,
  { id: 6, invoice_number: "INV-2026-010", customer_name: "Happy Kitchen", customer_id: 6, balance: 0, total: 5500000, paid: 5500000, due_date: "2026-01-25", days_overdue: 0, status: "PAID", email: "happy@kitchen.com", phone: "+856 20 555 0006" },
  { id: 7, invoice_number: "INV-2026-011", customer_name: "Spice Garden", customer_id: 7, balance: 9800000, total: 9800000, paid: 0, due_date: "2026-01-28", days_overdue: 0, status: "POSTED", email: "spice@garden.com", phone: "+856 20 555 0007" },
  { id: 8, invoice_number: "INV-2026-012", customer_name: "Noodle House", customer_id: 8, balance: 7200000, total: 7200000, paid: 0, due_date: "2026-01-30", days_overdue: 0, status: "DRAFT", email: "noodle@house.com", phone: "+856 20 555 0008" },
];

const sampleRecentPayments = [
  { id: 1, payment_number: "PAY-2026-015", customer_name: "Happy Kitchen", amount: 5500000, payment_date: "2026-01-20", method: "Bank Transfer", invoice: "INV-2026-010" },
  { id: 2, payment_number: "PAY-2026-014", customer_name: "Spice Garden", amount: 3200000, payment_date: "2026-01-19", method: "Check", invoice: "INV-2026-005" },
  { id: 3, payment_number: "PAY-2026-013", customer_name: "Noodle House", amount: 7800000, payment_date: "2026-01-18", method: "Cash", invoice: "INV-2026-003" },
  { id: 4, payment_number: "PAY-2026-012", customer_name: "Pizza Palace", amount: 4100000, payment_date: "2026-01-17", method: "Bank Transfer", invoice: "INV-2025-095" },
  { id: 5, payment_number: "PAY-2026-011", customer_name: "Fresh Mart", amount: 3800000, payment_date: "2026-01-16", method: "Bank Transfer", invoice: "INV-2025-098" },
];

const sampleTopCustomers = [
  { id: 1, name: "ABC Restaurant", balance: 25000000, credit_limit: 50000000, credit_used: 50, last_payment: "2026-01-15", status: "Good" },
  { id: 2, name: "XYZ Food Court", balance: 18500000, credit_limit: 30000000, credit_used: 62, last_payment: "2026-01-10", status: "Warning" },
  { id: 3, name: "Fresh Mart", balance: 15200000, credit_limit: 25000000, credit_used: 61, last_payment: "2026-01-16", status: "Good" },
  { id: 4, name: "Golden Wok", balance: 12800000, credit_limit: 20000000, credit_used: 64, last_payment: "2025-12-28", status: "At Risk" },
  { id: 5, name: "Ocean Seafood", balance: 9500000, credit_limit: 15000000, credit_used: 63, last_payment: "2026-01-05", status: "Good" },
];

const sampleMonthlyTrends = [
  { month: "Aug", invoiced: 85000000, collected: 78000000, outstanding: 45000000 },
  { month: "Sep", invoiced: 92000000, collected: 85000000, outstanding: 52000000 },
  { month: "Oct", invoiced: 88000000, collected: 90000000, outstanding: 50000000 },
  { month: "Nov", invoiced: 95000000, collected: 82000000, outstanding: 63000000 },
  { month: "Dec", invoiced: 110000000, collected: 95000000, outstanding: 78000000 },
  { month: "Jan", invoiced: 75000000, collected: 53000000, outstanding: 100000000 },
];

const sampleCustomers = [
  { id: 1, name: "ABC Restaurant" },
  { id: 2, name: "XYZ Food Court" },
  { id: 3, name: "Fresh Mart" },
  { id: 4, name: "Golden Wok" },
  { id: 5, name: "Ocean Seafood" },
  { id: 6, name: "Happy Kitchen" },
  { id: 7, name: "Spice Garden" },
  { id: 8, name: "Noodle House" },
];

export default function ARDashboard() {
  // Dialog states
  const [paymentDialogOpen, setPaymentDialogOpen] = useState(false);
  const [invoiceDialogOpen, setInvoiceDialogOpen] = useState(false);
  const [statementDialogOpen, setStatementDialogOpen] = useState(false);
  const [reminderDialogOpen, setReminderDialogOpen] = useState(false);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  
  // Selected items
  const [selectedInvoice, setSelectedInvoice] = useState<any>(null);
  const [selectedCustomer, setSelectedCustomer] = useState<any>(null);
  const [selectedInvoices, setSelectedInvoices] = useState<number[]>([]);
  
  // Filter states
  const [customerFilter, setCustomerFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [searchTerm, setSearchTerm] = useState("");
  
  // Form states
  const [paymentAmount, setPaymentAmount] = useState("");
  const [paymentMethod, setPaymentMethod] = useState("");
  const [paymentReference, setPaymentReference] = useState("");
  
  // Invoice form states
  const [invoiceCustomer, setInvoiceCustomer] = useState("");
  const [invoiceAmount, setInvoiceAmount] = useState("");
  const [invoiceDueDate, setInvoiceDueDate] = useState("");
  const [invoiceDescription, setInvoiceDescription] = useState("");
  
  // Loading states
  const [isLoading, setIsLoading] = useState(false);

  // Use sample data
  const aging = sampleAgingData;
  const overdueInvoices = sampleOverdueInvoices;
  const allInvoices = sampleAllInvoices;
  const recentPayments = sampleRecentPayments;
  const topCustomers = sampleTopCustomers;
  const monthlyTrends = sampleMonthlyTrends;

  // Filter invoices
  const filteredInvoices = allInvoices.filter(inv => {
    if (customerFilter !== "all" && inv.customer_id.toString() !== customerFilter) return false;
    if (statusFilter === "overdue" && inv.days_overdue <= 0) return false;
    if (statusFilter === "paid" && inv.balance > 0) return false;
    if (statusFilter === "pending" && (inv.balance === 0 || inv.days_overdue > 0)) return false;
    if (searchTerm && !inv.invoice_number.toLowerCase().includes(searchTerm.toLowerCase()) && 
        !inv.customer_name.toLowerCase().includes(searchTerm.toLowerCase())) return false;
    return true;
  });

  const chartData = [
    { name: "Current", value: aging.current, color: "#10b981" },
    { name: "1-30 Days", value: aging.days_1_30, color: "#3b82f6" },
    { name: "31-60 Days", value: aging.days_31_60, color: "#f59e0b" },
    { name: "61-90 Days", value: aging.days_61_90, color: "#f97316" },
    { name: "90+ Days", value: aging.over_90, color: "#ef4444" },
  ];

  const pieData = [
    { name: "Current", value: aging.current, color: "#10b981" },
    { name: "Overdue", value: aging.total - aging.current, color: "#ef4444" },
  ];

  const totalReceivable = aging.total;
  const totalOverdue = aging.days_1_30 + aging.days_31_60 + aging.days_61_90 + aging.over_90;
  const collectionRate = ((aging.current / aging.total) * 100).toFixed(1);
  const avgDSO = 24; // Days Sales Outstanding

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'LAK',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const handleRecordPayment = (invoice: any) => {
    setSelectedInvoice(invoice);
    setPaymentAmount(invoice.balance.toString());
    setPaymentDialogOpen(true);
  };

  const handleViewDetails = (invoice: any) => {
    setSelectedInvoice(invoice);
    setDetailsDialogOpen(true);
  };

  const handleViewStatement = (customer: any) => {
    setSelectedCustomer(customer);
    setStatementDialogOpen(true);
  };

  const handleSendReminder = (invoice: any) => {
    setSelectedInvoice(invoice);
    setReminderDialogOpen(true);
  };

  const handleBulkReminder = () => {
    if (selectedInvoices.length === 0) {
      toast.error("Please select invoices to send reminders");
      return;
    }
    toast.success(`Reminders sent to ${selectedInvoices.length} customers`);
    setSelectedInvoices([]);
  };

  const toggleInvoiceSelection = (id: number) => {
    setSelectedInvoices(prev => 
      prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]
    );
  };

  const selectAllOverdue = () => {
    if (selectedInvoices.length === overdueInvoices.length) {
      setSelectedInvoices([]);
    } else {
      setSelectedInvoices(overdueInvoices.map(i => i.id));
    }
  };

  const submitPayment = () => {
    if (!paymentAmount || !paymentMethod) {
      toast.error("Please fill in all required fields");
      return;
    }
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Payment of ${formatCurrency(parseFloat(paymentAmount))} recorded successfully`);
      setPaymentDialogOpen(false);
      resetPaymentForm();
      setIsLoading(false);
    }, 1000);
  };

  const submitInvoice = () => {
    if (!invoiceCustomer || !invoiceAmount || !invoiceDueDate) {
      toast.error("Please fill in all required fields");
      return;
    }
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Invoice created for ${formatCurrency(parseFloat(invoiceAmount))}`);
      setInvoiceDialogOpen(false);
      resetInvoiceForm();
      setIsLoading(false);
    }, 1000);
  };

  const submitReminder = () => {
    setIsLoading(true);
    setTimeout(() => {
      toast.success(`Payment reminder sent to ${selectedInvoice?.customer_name}`);
      setReminderDialogOpen(false);
      setIsLoading(false);
    }, 1000);
  };

  const resetPaymentForm = () => {
    setPaymentAmount("");
    setPaymentMethod("");
    setPaymentReference("");
    setSelectedInvoice(null);
  };

  const resetInvoiceForm = () => {
    setInvoiceCustomer("");
    setInvoiceAmount("");
    setInvoiceDueDate("");
    setInvoiceDescription("");
  };

  const getStatusBadge = (invoice: any) => {
    if (invoice.balance === 0) {
      return <Badge className="bg-emerald-100 text-emerald-800">PAID</Badge>;
    }
    if (invoice.days_overdue > 30) {
      return <Badge variant="destructive">OVERDUE 30+</Badge>;
    }
    if (invoice.days_overdue > 0) {
      return <Badge className="bg-orange-100 text-orange-800">OVERDUE</Badge>;
    }
    if (invoice.paid > 0) {
      return <Badge className="bg-blue-100 text-blue-800">PARTIAL</Badge>;
    }
    return <Badge className="bg-gray-100 text-gray-800">PENDING</Badge>;
  };

  const getCreditStatusColor = (status: string) => {
    switch (status) {
      case "Good": return "text-emerald-600";
      case "Warning": return "text-amber-600";
      case "At Risk": return "text-red-600";
      default: return "";
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-emerald-600 to-emerald-400 bg-clip-text text-transparent">
            Accounts Receivable
          </h1>
          <p className="text-muted-foreground">
            Monitor receivables, aging, and customer collections
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" /> Export
          </Button>
          <Button variant="outline" size="sm">
            <Printer className="mr-2 h-4 w-4" /> Print Report
          </Button>
          <Button variant="outline" size="sm" onClick={handleBulkReminder}>
            <Send className="mr-2 h-4 w-4" /> Send Reminders ({selectedInvoices.length})
          </Button>
          <Button onClick={() => setInvoiceDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" /> New Invoice
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <Card className="border-l-4 border-l-emerald-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Receivable</CardTitle>
            <DollarSign className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalReceivable)}</div>
            <div className="flex items-center text-xs text-emerald-600 mt-1">
              <ArrowUpRight className="h-3 w-3 mr-1" />
              +12% from last month
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
              {((totalOverdue / totalReceivable) * 100).toFixed(1)}% of total • {overdueInvoices.length} invoices
            </p>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-blue-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">DSO (Days Sales Outstanding)</CardTitle>
            <Clock className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{avgDSO} Days</div>
            <div className="flex items-center text-xs text-emerald-600 mt-1">
              <ArrowDownRight className="h-3 w-3 mr-1" />
              -2 days improvement
            </div>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-purple-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Collection Rate</CardTitle>
            <BarChart3 className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">{collectionRate}%</div>
            <div className="w-full h-2 bg-muted rounded-full overflow-hidden mt-2">
              <div className="h-full bg-purple-500 rounded-full" style={{ width: `${collectionRate}%` }} />
            </div>
          </CardContent>
        </Card>

        <Card className="border-l-4 border-l-amber-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Collected This Month</CardTitle>
            <TrendingUp className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(53000000)}</div>
            <p className="text-xs text-muted-foreground">
              Target: {formatCurrency(75000000)}
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
            <CardDescription>Receivables by age bucket</CardDescription>
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

        {/* Collection Trend */}
        <Card>
          <CardHeader>
            <CardTitle>Collection Trend</CardTitle>
            <CardDescription>Monthly invoiced vs collected</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[250px]">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={monthlyTrends}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" fontSize={10} />
                  <YAxis tickFormatter={(value) => `${(value/1000000)}M`} fontSize={10} />
                  <Tooltip formatter={(value: number) => formatCurrency(value)} />
                  <Area type="monotone" dataKey="invoiced" stackId="1" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.3} name="Invoiced" />
                  <Area type="monotone" dataKey="collected" stackId="2" stroke="#10b981" fill="#10b981" fillOpacity={0.5} name="Collected" />
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
                  <Line type="monotone" dataKey="outstanding" stroke="#ef4444" strokeWidth={2} dot={{ fill: '#ef4444' }} name="Outstanding" />
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
            <CardTitle>Invoice Management</CardTitle>
            <div className="flex flex-wrap items-center gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search invoices..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-9 w-[200px]"
                />
              </div>
              <Select value={customerFilter} onValueChange={setCustomerFilter}>
                <SelectTrigger className="w-[150px]">
                  <SelectValue placeholder="All Customers" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Customers</SelectItem>
                  {sampleCustomers.map(c => (
                    <SelectItem key={c.id} value={c.id.toString()}>{c.name}</SelectItem>
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
          <Tabs defaultValue="overdue" className="space-y-4">
            <TabsList>
              <TabsTrigger value="overdue" className="gap-2">
                <AlertCircle className="h-4 w-4" />
                Overdue
                <Badge variant="destructive" className="ml-1">{overdueInvoices.length}</Badge>
              </TabsTrigger>
              <TabsTrigger value="all" className="gap-2">
                <FileText className="h-4 w-4" />
                All Invoices
                <Badge variant="secondary" className="ml-1">{allInvoices.length}</Badge>
              </TabsTrigger>
              <TabsTrigger value="payments" className="gap-2">
                <CreditCard className="h-4 w-4" />
                Payments
              </TabsTrigger>
              <TabsTrigger value="customers" className="gap-2">
                <Users className="h-4 w-4" />
                Customers
              </TabsTrigger>
            </TabsList>

            <TabsContent value="overdue">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                  <Checkbox 
                    checked={selectedInvoices.length === overdueInvoices.length}
                    onCheckedChange={selectAllOverdue}
                  />
                  <span className="text-sm text-muted-foreground">Select All</span>
                </div>
                <p className="text-sm text-muted-foreground">
                  {selectedInvoices.length} selected
                </p>
              </div>
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {overdueInvoices.map((invoice) => (
                    <div key={invoice.id} className="flex items-center gap-4 p-4 rounded-lg border bg-card hover:bg-muted/50 transition-colors">
                      <Checkbox 
                        checked={selectedInvoices.includes(invoice.id)}
                        onCheckedChange={() => toggleInvoiceSelection(invoice.id)}
                      />
                      <div className="w-10 h-10 rounded-full bg-destructive/10 flex items-center justify-center shrink-0">
                        <Receipt className="h-5 w-5 text-destructive" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <p className="font-medium">{invoice.customer_name}</p>
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
                        <Button variant="ghost" size="icon" onClick={() => handleSendReminder(invoice)}>
                          <Mail className="h-4 w-4" />
                        </Button>
                        <Button size="sm" onClick={() => handleRecordPayment(invoice)}>
                          <CreditCard className="h-4 w-4 mr-1" /> Pay
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
                    <TableHead>Invoice #</TableHead>
                    <TableHead>Customer</TableHead>
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
                      <TableCell>{invoice.customer_name}</TableCell>
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
                            <Button variant="ghost" size="icon" onClick={() => handleRecordPayment(invoice)}>
                              <CreditCard className="h-4 w-4" />
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
                          <p className="font-medium">{payment.customer_name}</p>
                          <p className="text-sm text-muted-foreground">{payment.payment_number} • {payment.invoice}</p>
                        </div>
                      </div>
                      <Badge variant="outline">{payment.method}</Badge>
                      <div className="text-right">
                        <p className="font-bold text-emerald-600">{formatCurrency(payment.amount)}</p>
                        <p className="text-sm text-muted-foreground">{payment.payment_date}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </TabsContent>

            <TabsContent value="customers">
              <ScrollArea className="h-[400px]">
                <div className="space-y-3">
                  {topCustomers.map((customer, index) => (
                    <div key={customer.id} className="flex items-center justify-between p-4 rounded-lg border bg-card">
                      <div className="flex items-center gap-4">
                        <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center font-bold text-primary">
                          {index + 1}
                        </div>
                        <div>
                          <div className="flex items-center gap-2">
                            <p className="font-medium">{customer.name}</p>
                            <Badge className={`${getCreditStatusColor(customer.status)} bg-transparent border`}>
                              {customer.status}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground">
                            Last payment: {customer.last_payment}
                          </p>
                        </div>
                      </div>
                      <div className="text-center">
                        <p className="text-sm text-muted-foreground">Credit Used</p>
                        <div className="w-24 h-2 bg-muted rounded-full overflow-hidden mt-1">
                          <div 
                            className={`h-full rounded-full ${customer.credit_used > 80 ? 'bg-red-500' : customer.credit_used > 60 ? 'bg-amber-500' : 'bg-emerald-500'}`}
                            style={{ width: `${customer.credit_used}%` }}
                          />
                        </div>
                        <p className="text-xs mt-1">{customer.credit_used}%</p>
                      </div>
                      <div className="text-right">
                        <p className="font-bold">{formatCurrency(customer.balance)}</p>
                        <p className="text-sm text-muted-foreground">of {formatCurrency(customer.credit_limit)}</p>
                      </div>
                      <Button variant="outline" size="sm" onClick={() => handleViewStatement(customer)}>
                        View Statement
                      </Button>
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
            <DialogTitle>Record Payment</DialogTitle>
            <DialogDescription>
              {selectedInvoice && (
                <span>
                  Recording payment for <strong>{selectedInvoice.invoice_number}</strong>
                  <br />
                  Customer: {selectedInvoice.customer_name} • Balance: {formatCurrency(selectedInvoice.balance)}
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
                  <SelectItem value="cash">Cash</SelectItem>
                  <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                  <SelectItem value="check">Check</SelectItem>
                  <SelectItem value="credit_card">Credit Card</SelectItem>
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
              Record Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Invoice Dialog */}
      <Dialog open={invoiceDialogOpen} onOpenChange={setInvoiceDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Create New Invoice</DialogTitle>
            <DialogDescription>Create a new invoice for a customer</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Customer *</Label>
              <Select value={invoiceCustomer} onValueChange={setInvoiceCustomer}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select customer" />
                </SelectTrigger>
                <SelectContent>
                  {sampleCustomers.map(c => (
                    <SelectItem key={c.id} value={c.id.toString()}>{c.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Amount *</Label>
              <Input
                type="number"
                placeholder="0.00"
                value={invoiceAmount}
                onChange={(e) => setInvoiceAmount(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Due Date *</Label>
              <Input
                type="date"
                value={invoiceDueDate}
                onChange={(e) => setInvoiceDueDate(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Description</Label>
              <Textarea
                placeholder="Invoice description..."
                value={invoiceDescription}
                onChange={(e) => setInvoiceDescription(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setInvoiceDialogOpen(false)}>Cancel</Button>
            <Button onClick={submitInvoice} disabled={isLoading}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              Create Invoice
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reminder Dialog */}
      <Dialog open={reminderDialogOpen} onOpenChange={setReminderDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Send Payment Reminder</DialogTitle>
            <DialogDescription>
              {selectedInvoice && (
                <span>
                  Send reminder to <strong>{selectedInvoice.customer_name}</strong> for {selectedInvoice.invoice_number}
                </span>
              )}
            </DialogDescription>
          </DialogHeader>
          {selectedInvoice && (
            <div className="space-y-4 py-4">
              <div className="flex items-center gap-4 p-3 bg-muted rounded-lg">
                <Mail className="h-5 w-5 text-muted-foreground" />
                <span>{selectedInvoice.email}</span>
              </div>
              <div className="flex items-center gap-4 p-3 bg-muted rounded-lg">
                <Phone className="h-5 w-5 text-muted-foreground" />
                <span>{selectedInvoice.phone}</span>
              </div>
              <Separator />
              <div className="bg-amber-50 dark:bg-amber-950 p-4 rounded-lg">
                <p className="text-sm">
                  <strong>Amount Due:</strong> {formatCurrency(selectedInvoice.balance)}
                </p>
                <p className="text-sm">
                  <strong>Days Overdue:</strong> {selectedInvoice.days_overdue} days
                </p>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setReminderDialogOpen(false)}>Cancel</Button>
            <Button onClick={submitReminder} disabled={isLoading}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              <Send className="h-4 w-4 mr-2" /> Send Reminder
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Details Dialog */}
      <Dialog open={detailsDialogOpen} onOpenChange={setDetailsDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>Invoice Details</DialogTitle>
          </DialogHeader>
          {selectedInvoice && (
            <div className="space-y-6">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="text-2xl font-bold">{selectedInvoice.invoice_number}</h3>
                  <p className="text-muted-foreground">{selectedInvoice.customer_name}</p>
                </div>
                {getStatusBadge(selectedInvoice)}
              </div>
              <Separator />
              <div className="grid grid-cols-2 gap-6">
                <div>
                  <Label className="text-muted-foreground">Invoice Total</Label>
                  <p className="text-xl font-bold">{formatCurrency(selectedInvoice.total)}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Amount Paid</Label>
                  <p className="text-xl font-bold text-emerald-600">{formatCurrency(selectedInvoice.paid)}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Balance Due</Label>
                  <p className="text-xl font-bold text-destructive">{formatCurrency(selectedInvoice.balance)}</p>
                </div>
                <div>
                  <Label className="text-muted-foreground">Due Date</Label>
                  <p className="text-xl font-bold">{selectedInvoice.due_date}</p>
                </div>
              </div>
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
                <Button variant="outline" onClick={() => { setDetailsDialogOpen(false); handleSendReminder(selectedInvoice); }}>
                  <Send className="h-4 w-4 mr-2" /> Send Reminder
                </Button>
                <Button onClick={() => { setDetailsDialogOpen(false); handleRecordPayment(selectedInvoice); }}>
                  <CreditCard className="h-4 w-4 mr-2" /> Record Payment
                </Button>
              </>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Statement Dialog */}
      <Dialog open={statementDialogOpen} onOpenChange={setStatementDialogOpen}>
        <DialogContent className="sm:max-w-[700px]">
          <DialogHeader>
            <DialogTitle>Customer Statement</DialogTitle>
            <DialogDescription>
              {selectedCustomer?.name}
            </DialogDescription>
          </DialogHeader>
          {selectedCustomer && (
            <div className="space-y-4">
              <div className="grid grid-cols-3 gap-4">
                <Card>
                  <CardContent className="pt-4">
                    <p className="text-sm text-muted-foreground">Credit Limit</p>
                    <p className="text-xl font-bold">{formatCurrency(selectedCustomer.credit_limit)}</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-4">
                    <p className="text-sm text-muted-foreground">Current Balance</p>
                    <p className="text-xl font-bold text-destructive">{formatCurrency(selectedCustomer.balance)}</p>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="pt-4">
                    <p className="text-sm text-muted-foreground">Available Credit</p>
                    <p className="text-xl font-bold text-emerald-600">{formatCurrency(selectedCustomer.credit_limit - selectedCustomer.balance)}</p>
                  </CardContent>
                </Card>
              </div>
              <div className="border rounded-lg p-4">
                <h4 className="font-medium mb-3">Recent Transactions</h4>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Date</TableHead>
                      <TableHead>Description</TableHead>
                      <TableHead className="text-right">Debit</TableHead>
                      <TableHead className="text-right">Credit</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow>
                      <TableCell>2026-01-15</TableCell>
                      <TableCell>Invoice INV-2026-001</TableCell>
                      <TableCell className="text-right">{formatCurrency(15000000)}</TableCell>
                      <TableCell className="text-right">-</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2026-01-10</TableCell>
                      <TableCell>Payment PAY-2026-008</TableCell>
                      <TableCell className="text-right">-</TableCell>
                      <TableCell className="text-right text-emerald-600">{formatCurrency(2500000)}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2026-01-05</TableCell>
                      <TableCell>Invoice INV-2025-098</TableCell>
                      <TableCell className="text-right">{formatCurrency(10000000)}</TableCell>
                      <TableCell className="text-right">-</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setStatementDialogOpen(false)}>Close</Button>
            <Button variant="outline">
              <Printer className="h-4 w-4 mr-2" /> Print Statement
            </Button>
            <Button variant="outline">
              <Download className="h-4 w-4 mr-2" /> Download PDF
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
