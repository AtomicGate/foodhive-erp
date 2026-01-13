import { useState } from "react";
import { Link } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  TableFooter,
} from "@/components/ui/table";
import { 
  Download,
  Printer,
  Loader2,
  ChevronRight,
  Calendar
} from "lucide-react";
import { financialService } from "@/services/financialService";

interface TrialBalanceEntry {
  account_id: number;
  account_code: string;
  account_name: string;
  account_type: string;
  debit_balance: number;
  credit_balance: number;
}

export default function TrialBalance() {
  const [selectedPeriod, setSelectedPeriod] = useState("current");

  const { data: trialBalance, isLoading, error } = useQuery({
    queryKey: ['trialBalance', selectedPeriod],
    queryFn: () => financialService.getTrialBalance()
  });

  const { data: currentPeriod } = useQuery({
    queryKey: ['currentPeriod'],
    queryFn: financialService.getCurrentPeriod
  });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8 text-center text-red-500">
        Failed to load trial balance
      </div>
    );
  }

  const entries: TrialBalanceEntry[] = trialBalance?.accounts || trialBalance || [];
  
  const totals = entries.reduce((acc, entry) => ({
    debit: acc.debit + (entry.debit_balance || 0),
    credit: acc.credit + (entry.credit_balance || 0)
  }), { debit: 0, credit: 0 });

  const formatCurrency = (value: number) => {
    if (value === 0) return "-";
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value);
  };

  const getAccountTypeColor = (type: string) => {
    switch (type) {
      case "ASSET": return "bg-emerald-100 text-emerald-800";
      case "LIABILITY": return "bg-red-100 text-red-800";
      case "EQUITY": return "bg-blue-100 text-blue-800";
      case "REVENUE": return "bg-purple-100 text-purple-800";
      case "EXPENSE": return "bg-amber-100 text-amber-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const isBalanced = Math.abs(totals.debit - totals.credit) < 0.01;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-1">
            <Link href="/gl" className="hover:text-foreground">General Ledger</Link>
            <ChevronRight className="h-4 w-4" />
            <Link href="/gl/reports" className="hover:text-foreground">Reports</Link>
            <ChevronRight className="h-4 w-4" />
            <span>Trial Balance</span>
          </div>
          <h1 className="text-3xl font-bold tracking-tight">Trial Balance</h1>
          <p className="text-muted-foreground">
            Summary of all account balances for the period
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
            <SelectTrigger className="w-[200px]">
              <Calendar className="mr-2 h-4 w-4" />
              <SelectValue placeholder="Select period" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="current">Current Period</SelectItem>
              <SelectItem value="previous">Previous Period</SelectItem>
              <SelectItem value="ytd">Year to Date</SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" /> Export
          </Button>
          <Button variant="outline" size="sm">
            <Printer className="mr-2 h-4 w-4" /> Print
          </Button>
        </div>
      </div>

      {/* Period Info */}
      {currentPeriod && (
        <Card className="bg-muted/50">
          <CardContent className="flex items-center justify-between py-4">
            <div>
              <p className="font-medium">Period: {currentPeriod.name || 'Current Period'}</p>
              <p className="text-sm text-muted-foreground">
                {currentPeriod.start_date ? new Date(currentPeriod.start_date).toLocaleDateString() : ''} - 
                {currentPeriod.end_date ? new Date(currentPeriod.end_date).toLocaleDateString() : ''}
              </p>
            </div>
            <Badge variant={isBalanced ? "default" : "destructive"}>
              {isBalanced ? "Balanced" : "Out of Balance"}
            </Badge>
          </CardContent>
        </Card>
      )}

      {/* Trial Balance Table */}
      <Card>
        <CardHeader>
          <CardTitle>Account Balances</CardTitle>
          <CardDescription>
            Showing {entries.length} accounts with balances
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[120px]">Account Code</TableHead>
                <TableHead>Account Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead className="text-right">Debit</TableHead>
                <TableHead className="text-right">Credit</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {entries.length > 0 ? (
                entries.map((entry) => (
                  <TableRow key={entry.account_id}>
                    <TableCell className="font-mono">{entry.account_code}</TableCell>
                    <TableCell className="font-medium">{entry.account_name}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className={getAccountTypeColor(entry.account_type)}>
                        {entry.account_type}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {formatCurrency(entry.debit_balance)}
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {formatCurrency(entry.credit_balance)}
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                    No account balances found for this period
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
            <TableFooter>
              <TableRow className="bg-muted/50 font-bold">
                <TableCell colSpan={3}>Total</TableCell>
                <TableCell className="text-right">{formatCurrency(totals.debit)}</TableCell>
                <TableCell className="text-right">{formatCurrency(totals.credit)}</TableCell>
              </TableRow>
              {!isBalanced && (
                <TableRow className="bg-destructive/10">
                  <TableCell colSpan={3} className="text-destructive font-medium">
                    Difference (Out of Balance)
                  </TableCell>
                  <TableCell colSpan={2} className="text-right text-destructive font-bold">
                    {formatCurrency(Math.abs(totals.debit - totals.credit))}
                  </TableCell>
                </TableRow>
              )}
            </TableFooter>
          </Table>
        </CardContent>
      </Card>

      {/* Summary by Account Type */}
      <div className="grid gap-4 md:grid-cols-5">
        {['ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE'].map((type) => {
          const typeEntries = entries.filter(e => e.account_type === type);
          const typeDebit = typeEntries.reduce((sum, e) => sum + (e.debit_balance || 0), 0);
          const typeCredit = typeEntries.reduce((sum, e) => sum + (e.credit_balance || 0), 0);
          const netBalance = typeDebit - typeCredit;
          
          return (
            <Card key={type}>
              <CardContent className="pt-6">
                <div className="space-y-2">
                  <Badge variant="outline" className={getAccountTypeColor(type)}>
                    {type}
                  </Badge>
                  <div>
                    <p className="text-sm text-muted-foreground">Net Balance</p>
                    <p className={`text-xl font-bold ${netBalance < 0 ? 'text-red-600' : ''}`}>
                      {formatCurrency(Math.abs(netBalance))}
                      {netBalance < 0 ? ' (Cr)' : netBalance > 0 ? ' (Dr)' : ''}
                    </p>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {typeEntries.length} accounts
                  </p>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
