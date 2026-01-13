import { useRoute } from "wouter";
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
  ArrowUpDown, 
  Mail,
  Phone,
} from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { useQuery } from "@tanstack/react-query";
import { entityService } from "@/services/entityService";

type Entity = {
  id: string;
  name: string;
  code: string;
  email: string;
  phone: string;
  type: "Customer" | "Vendor" | "Employee";
  status: "Active" | "Inactive";
  balance?: number;
};

export default function EntityList() {
  const [match, params] = useRoute("/entities/:type");
  const type = params?.type || "customers";
  const entityType = type.charAt(0).toUpperCase() + type.slice(1, -1); // customers -> Customer
  
  const { data: entities = [], isLoading } = useQuery({
    queryKey: ['entities', type],
    queryFn: () => entityService.getEntities(type),
  });

  const columns: ColumnDef<Entity>[] = [
    {
      accessorKey: "code",
      header: "Code",
      cell: ({ row }) => <div className="font-mono text-xs">{row.getValue("code")}</div>,
    },
    {
      accessorKey: "name",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Name
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        )
      },
      cell: ({ row }) => (
        <div>
          <div className="font-medium">{row.getValue("name")}</div>
          <div className="text-xs text-muted-foreground flex items-center gap-1 mt-0.5">
            <Mail className="h-3 w-3" /> {row.original.email}
          </div>
        </div>
      ),
    },
    {
      accessorKey: "phone",
      header: "Phone",
      cell: ({ row }) => (
        <div className="flex items-center gap-2 text-sm">
          <Phone className="h-3 w-3 text-muted-foreground" />
          {row.getValue("phone")}
        </div>
      ),
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
          <Badge variant={status === "Active" ? "default" : "secondary"}>
            {status}
          </Badge>
        );
      },
    },
    ...(type !== "employees" ? [{
      accessorKey: "balance",
      header: () => <div className="text-right">Balance</div>,
      cell: ({ row }: { row: any }) => {
        const amount = parseFloat(row.getValue("balance"));
        const formatted = new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "USD",
        }).format(amount);
        return (
          <div className={`text-right font-medium ${amount > 0 ? "text-destructive" : "text-muted-foreground"}`}>
            {formatted}
          </div>
        );
      },
    }] : []),
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
              <DropdownMenuItem>Edit {entityType}</DropdownMenuItem>
              <DropdownMenuItem>View Details</DropdownMenuItem>
              {type !== "employees" && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>View Orders</DropdownMenuItem>
                  <DropdownMenuItem>Statement</DropdownMenuItem>
                </>
              )}
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
          <h1 className="text-3xl font-bold tracking-tight capitalize">{type}</h1>
          <p className="text-muted-foreground">
            Manage your {type} directory and information.
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" /> Add {entityType}
        </Button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : (
        <DataTable 
          columns={columns} 
          data={entities} 
          searchKey="name"
          searchPlaceholder={`Filter ${type}...`}
          onExport={() => console.log("Exporting...")}
          onPrint={() => window.print()}
        />
      )}
    </div>
  );
}
