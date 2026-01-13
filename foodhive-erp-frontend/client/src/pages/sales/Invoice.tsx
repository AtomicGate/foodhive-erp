import { Link, useRoute } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { 
  Printer, 
  ArrowLeft, 
  Download,
  Mail,
  CreditCard,
  Building2,
  Phone,
  Globe,
  Loader2
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { salesService } from "@/services/salesService";

export default function Invoice() {
  const [, params] = useRoute("/invoices/:id");
  const id = params?.id;

  const { data: invoice, isLoading, error } = useQuery({
    queryKey: ['invoice', id],
    queryFn: () => salesService.getInvoice(id!),
    enabled: !!id
  });

  if (isLoading) return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  if (error || !invoice) return <div className="p-8 text-center text-red-500">Failed to load invoice</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between print:hidden">
        <div className="flex items-center gap-4">
          <Link href="/sales-orders">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-3xl font-bold tracking-tight">Invoice #{invoice.invoiceNumber || id}</h1>
              <Badge variant={invoice.status === 'Paid' ? 'default' : 'destructive'} className="ml-2">
                {invoice.status}
              </Badge>
            </div>
            <p className="text-muted-foreground">
              Issued on {new Date(invoice.date).toLocaleDateString()} • Due on {new Date(invoice.dueDate).toLocaleDateString()}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline">
            <Mail className="mr-2 h-4 w-4" /> Email
          </Button>
          <Button variant="outline">
            <Download className="mr-2 h-4 w-4" /> PDF
          </Button>
          <Button onClick={() => window.print()}>
            <Printer className="mr-2 h-4 w-4" /> Print
          </Button>
          <Button disabled={invoice.status === 'Paid'}>
            <CreditCard className="mr-2 h-4 w-4" /> Record Payment
          </Button>
        </div>
      </div>

      <Card className="max-w-4xl mx-auto print:shadow-none print:border-none">
        <CardHeader className="space-y-6">
          <div className="flex justify-between items-start">
            <div className="flex items-center gap-2">
              <div className="w-10 h-10 rounded-lg bg-primary flex items-center justify-center text-primary-foreground font-bold text-xl">
                FH
              </div>
              <div>
                <h2 className="text-2xl font-bold text-primary">FoodHive ERP</h2>
                <p className="text-sm text-muted-foreground">Wholesale Food Distribution</p>
              </div>
            </div>
            <div className="text-right space-y-1">
              <h3 className="text-xl font-bold text-muted-foreground">INVOICE</h3>
              <p className="font-medium">{invoice.invoiceNumber || id}</p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-8 text-sm">
            <div className="space-y-2">
              <p className="font-semibold text-muted-foreground">Bill To:</p>
              <p className="font-bold text-lg">{invoice.customerName}</p>
              <p className="text-muted-foreground w-48">{invoice.billingAddress}</p>
              <div className="flex items-center gap-2 text-muted-foreground pt-2">
                <Mail className="h-3 w-3" /> {invoice.customerEmail}
              </div>
              <div className="flex items-center gap-2 text-muted-foreground">
                <Phone className="h-3 w-3" /> {invoice.customerPhone}
              </div>
            </div>
            <div className="space-y-2 text-right">
              <p className="font-semibold text-muted-foreground">Pay To:</p>
              <p className="font-bold text-lg">FoodHive Inc.</p>
              <p className="text-muted-foreground">456 Distribution Blvd</p>
              <p className="text-muted-foreground">Chicago, IL 60601</p>
              <div className="flex items-center justify-end gap-2 text-muted-foreground pt-2">
                <Globe className="h-3 w-3" /> www.foodhive.com
              </div>
              <div className="flex items-center justify-end gap-2 text-muted-foreground">
                <Building2 className="h-3 w-3" /> Bank of America
              </div>
              <p className="text-muted-foreground">Acct: **** 1234</p>
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Description</TableHead>
                <TableHead className="text-right">Quantity</TableHead>
                <TableHead className="text-right">Unit Price</TableHead>
                <TableHead className="text-right">Total</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {invoice.items?.map((item: any) => (
                <TableRow key={item.id}>
                  <TableCell className="font-medium">
                    {item.productName}
                    <span className="block text-xs text-muted-foreground">Unit: {item.unit}</span>
                  </TableCell>
                  <TableCell className="text-right">{item.quantity}</TableCell>
                  <TableCell className="text-right">${Number(item.unitPrice).toFixed(2)}</TableCell>
                  <TableCell className="text-right">${Number(item.total).toFixed(2)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          <div className="flex justify-end">
            <div className="w-64 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Subtotal</span>
                <span>${Number(invoice.subtotal).toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Tax ({invoice.taxRate || 0}%)</span>
                <span>${Number(invoice.taxAmount).toFixed(2)}</span>
              </div>
              <Separator />
              <div className="flex justify-between font-bold text-lg">
                <span>Total Due</span>
                <span>${Number(invoice.totalAmount).toFixed(2)}</span>
              </div>
            </div>
          </div>
        </CardContent>

        <CardFooter className="bg-muted/50 p-6">
          <div className="w-full space-y-2">
            <p className="font-semibold text-sm">Payment Terms</p>
            <p className="text-sm text-muted-foreground">
              {invoice.paymentTerms || 'Net 30'}. Please include invoice number on your check.
              Thank you for your business!
            </p>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
