import { Link } from "wouter";
import { 
  ColumnDef, 
} from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown
} from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { useQuery } from "@tanstack/react-query";
import { salesService } from "@/services/salesService";

// Mock Data Type
type SalesOrder = {
  id: string;
  orderNumber: string;
  customer: string;
  date: string;
  total: number;
  status: "Draft" | "Confirmed" | "Picking" | "Shipped" | "Invoiced" | "Cancelled";
  paymentStatus: "Paid" | "Unpaid" | "Partial";
};

export default function SalesOrderList() {
  const { data: orders = [], isLoading } = useQuery({
    queryKey: ['sales-orders'],
    queryFn: () => salesService.getOrders(),
  });

  const columns: ColumnDef<SalesOrder>[] = [
    {
      accessorKey: "orderNumber",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Order #
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        )
      },
      cell: ({ row }) => <div className="font-medium">{row.getValue("orderNumber")}</div>,
    },
    {
      accessorKey: "customer",
      header: "Customer",
    },
    {
      accessorKey: "date",
      header: "Date",
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        let variant: "default" | "secondary" | "destructive" | "outline" = "default";
        let className = "";
        
        switch (status) {
          case "Draft":
            variant = "secondary";
            className = "bg-gray-100 text-gray-800 hover:bg-gray-200";
            break;
          case "Confirmed":
            variant = "default";
            className = "bg-blue-100 text-blue-800 hover:bg-blue-200";
            break;
          case "Picking":
            variant = "default";
            className = "bg-amber-100 text-amber-800 hover:bg-amber-200";
            break;
          case "Shipped":
            variant = "default";
            className = "bg-purple-100 text-purple-800 hover:bg-purple-200";
            break;
          case "Invoiced":
            variant = "default";
            className = "bg-emerald-100 text-emerald-800 hover:bg-emerald-200";
            break;
          case "Cancelled":
            variant = "destructive";
            break;
        }

        return (
          <Badge variant={variant} className={className}>
            {status}
          </Badge>
        );
      },
    },
    {
      accessorKey: "total",
      header: () => <div className="text-right">Total</div>,
      cell: ({ row }) => {
        const amount = parseFloat(row.getValue("total"));
        const formatted = new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "USD",
        }).format(amount);
        return <div className="text-right font-medium">{formatted}</div>;
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const order = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => navigator.clipboard.writeText(order.id)}>
                Copy Order ID
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>View Details</DropdownMenuItem>
              <DropdownMenuItem>Edit Order</DropdownMenuItem>
              <DropdownMenuSeparator />
              <Link href={`/sales-orders/picklist/${order.id}`}>
                <DropdownMenuItem>Generate Pick List</DropdownMenuItem>
              </Link>
              <Link href={`/sales-orders/invoice/${order.id}`}>
                <DropdownMenuItem>Create Invoice</DropdownMenuItem>
              </Link>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Sales Orders</h1>
          <p className="text-muted-foreground">
            Manage your sales orders, track status, and process shipments.
          </p>
        </div>
        <Link href="/sales-orders/new">
          <Button>
            <Plus className="mr-2 h-4 w-4" /> Create Order
          </Button>
        </Link>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : (
        <DataTable 
          columns={columns} 
          data={orders} 
          searchKey="customer"
          onExport={() => console.log("Exporting...")}
          onPrint={() => window.print()}
        />
      )}
    </div>
  );
}
