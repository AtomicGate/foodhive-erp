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
  MoreHorizontal, 
  ArrowUpDown, 
  AlertTriangle,
  History,
  ArrowRightLeft
} from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { useQuery } from "@tanstack/react-query";
import { inventoryService } from "@/services/inventoryService";

type InventoryItem = {
  id: string;
  sku: string;
  productName: string;
  category: string;
  onHand: number;
  allocated: number;
  available: number;
  unit: string;
  warehouse: string;
  status: "In Stock" | "Low Stock" | "Out of Stock";
};

export default function InventoryList() {
  const { data: inventory = [], isLoading } = useQuery({
    queryKey: ['inventory'],
    queryFn: () => inventoryService.getInventory(),
  });

  const columns: ColumnDef<InventoryItem>[] = [
    {
      accessorKey: "sku",
      header: "SKU",
      cell: ({ row }) => <div className="font-mono text-xs">{row.getValue("sku")}</div>,
    },
    {
      accessorKey: "productName",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Product
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        )
      },
    },
    {
      accessorKey: "category",
      header: "Category",
    },
    {
      accessorKey: "warehouse",
      header: "Warehouse",
    },
    {
      accessorKey: "available",
      header: () => <div className="text-right">Available</div>,
      cell: ({ row }) => {
        const amount = parseFloat(row.getValue("available"));
        return (
          <div className="text-right font-medium">
            {amount} <span className="text-xs text-muted-foreground">{row.original.unit}</span>
          </div>
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        let variant: "default" | "secondary" | "destructive" | "outline" = "default";
        let className = "";
        
        switch (status) {
          case "In Stock":
            variant = "default";
            className = "bg-emerald-100 text-emerald-800";
            break;
          case "Low Stock":
            variant = "default";
            className = "bg-amber-100 text-amber-800";
            break;
          case "Out of Stock":
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
      id: "actions",
      cell: ({ row }) => {
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
              <DropdownMenuItem>View Details</DropdownMenuItem>
              <DropdownMenuItem>Adjust Stock</DropdownMenuItem>
              <DropdownMenuItem>Transfer Stock</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>View History</DropdownMenuItem>
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
          <h1 className="text-3xl font-bold tracking-tight">Inventory</h1>
          <p className="text-muted-foreground">
            Track stock levels, adjustments, and transfers across warehouses.
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <History className="mr-2 h-4 w-4" /> History
          </Button>
          <Button variant="outline">
            <ArrowRightLeft className="mr-2 h-4 w-4" /> Transfer
          </Button>
          <Button>
            <AlertTriangle className="mr-2 h-4 w-4" /> Adjust Stock
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : (
        <DataTable 
          columns={columns} 
          data={inventory} 
          searchKey="productName"
          searchPlaceholder="Filter products..."
          onExport={() => console.log("Exporting...")}
          onPrint={() => window.print()}
        />
      )}
    </div>
  );
}
