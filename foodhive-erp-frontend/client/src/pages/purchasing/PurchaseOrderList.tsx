import { Link } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
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
  ArrowUpDown, 
  PackageCheck,
  Loader2
} from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { purchasingService, PurchaseOrder } from "@/services/purchasingService";

export default function PurchaseOrderList() {
  const { data: purchaseOrders, isLoading, error } = useQuery({
    queryKey: ['purchaseOrders'],
    queryFn: purchasingService.getPurchaseOrders
  });

  const columns: ColumnDef<PurchaseOrder>[] = [
    {
      accessorKey: "poNumber",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            PO #
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        )
      },
      cell: ({ row }) => <div className="font-medium">{row.getValue("poNumber")}</div>,
    },
    {
      accessorKey: "vendorName",
      header: "Vendor",
    },
    {
      accessorKey: "date",
      header: "Date",
      cell: ({ row }) => new Date(row.getValue("date")).toLocaleDateString(),
    },
    {
      accessorKey: "expectedDate",
      header: "Expected",
      cell: ({ row }) => new Date(row.getValue("expectedDate")).toLocaleDateString(),
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
            className = "bg-gray-100 text-gray-800";
            break;
          case "Sent":
            variant = "default";
            className = "bg-blue-100 text-blue-800";
            break;
          case "Received":
            variant = "default";
            className = "bg-emerald-100 text-emerald-800";
            break;
          case "Partial":
            variant = "default";
            className = "bg-amber-100 text-amber-800";
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
      accessorKey: "totalAmount",
      header: () => <div className="text-right">Total</div>,
      cell: ({ row }) => {
        const amount = parseFloat(row.getValue("totalAmount"));
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
        const po = row.original;
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
              <DropdownMenuItem onClick={() => navigator.clipboard.writeText(po.id)}>
                Copy PO ID
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>View Details</DropdownMenuItem>
              <DropdownMenuItem>Edit PO</DropdownMenuItem>
              <DropdownMenuSeparator />
              <Link href={`/purchase-orders/receive/${po.id}`}>
                <DropdownMenuItem>
                  <PackageCheck className="mr-2 h-4 w-4" /> Receive Items
                </DropdownMenuItem>
              </Link>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  if (isLoading) return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  if (error) return <div className="p-8 text-center text-red-500">Failed to load purchase orders</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Purchase Orders</h1>
          <p className="text-muted-foreground">
            Manage vendor orders and incoming shipments.
          </p>
        </div>
        <Link href="/purchase-orders/new">
          <Button>
            <Plus className="mr-2 h-4 w-4" /> Create PO
          </Button>
        </Link>
      </div>

      <DataTable 
        columns={columns} 
        data={purchaseOrders || []} 
        searchKey="vendorName"
        searchPlaceholder="Filter vendors..."
      />
    </div>
  );
}
