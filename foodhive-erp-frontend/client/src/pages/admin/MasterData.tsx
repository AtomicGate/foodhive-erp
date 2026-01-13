import { useState, useEffect } from "react";
import { useRoute } from "wouter";
import { useQuery, useMutation, useQueryClient, keepPreviousData } from "@tanstack/react-query";
import { 
  ColumnDef, 
  flexRender, 
  getCoreRowModel, 
  useReactTable, 
  getPaginationRowModel,
  getSortedRowModel,
  SortingState,
  getFilteredRowModel,
  ColumnFiltersState,
  PaginationState,
  RowSelectionState
} from "@tanstack/react-table";
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown, 
  Search, 
  Building2,
  Shield,
  Warehouse,
  Loader2,
  Pencil,
  Trash2,
  ChevronLeft,
  ChevronRight,
  AlertTriangle,
  CheckSquare,
  Square,
  Download
} from "lucide-react";
import { masterDataService, Department, Role, Warehouse as WarehouseType } from "@/services/masterDataService";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Checkbox } from "@/components/ui/checkbox";
import { AuditLogViewer } from "@/components/AuditLogViewer";
import { auditLogService } from "@/services/auditLogService";
import { exportService } from "@/services/exportService";

type MasterData = {
  id: string;
  name: string;
  description?: string;
  status?: "Active" | "Inactive";
  count?: number;
  location?: string;
  manager?: string;
  code?: string;
  capacity?: string;
  permissions?: string[];
};

// Form Schemas
const departmentSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  code: z.string().min(2, "Code must be at least 2 characters"),
  manager: z.string().min(2, "Manager name is required"),
});

const roleSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  description: z.string().min(5, "Description must be at least 5 characters"),
  permissions: z.array(z.string()).optional(),
});

const warehouseSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  location: z.string().min(5, "Location must be at least 5 characters"),
  capacity: z.string().min(1, "Capacity is required"),
  manager: z.string().min(2, "Manager name is required"),
});

const AVAILABLE_PERMISSIONS = [
  { id: "view_dashboard", label: "View Dashboard" },
  { id: "view_sales", label: "View Sales" },
  { id: "create_sales", label: "Create Sales" },
  { id: "view_inventory", label: "View Inventory" },
  { id: "edit_inventory", label: "Edit Inventory" },
  { id: "view_finance", label: "View Finance" },
  { id: "manage_users", label: "Manage Users" },
];

export default function MasterData() {
  const [match, params] = useRoute("/admin/:type");
  const type = params?.type || "departments";
  const title = type.charAt(0).toUpperCase() + type.slice(1);
  const queryClient = useQueryClient();
  
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isBulkDeleteDialogOpen, setIsBulkDeleteDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<MasterData | null>(null);
  const [itemToDelete, setItemToDelete] = useState<MasterData | null>(null);
  
  // Pagination State
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  });

  // Reset pagination and selection when type changes
  useEffect(() => {
    setPagination({ pageIndex: 0, pageSize: 10 });
    setRowSelection({});
  }, [type]);

  // Data Fetching
  const { data, isLoading, error } = useQuery({
    queryKey: ['masterData', type, pagination.pageIndex, pagination.pageSize],
    queryFn: async () => {
      const page = pagination.pageIndex + 1;
      const limit = pagination.pageSize;
      
      switch (type) {
        case "departments": {
          const res = await masterDataService.getDepartments(page, limit);
          return {
            items: res.data.map((d: any) => ({
              ...d,
              id: String(d.id),
              description: `Managed by ${d.manager_name || 'N/A'}`,
              status: d.is_active ? "Active" : "Inactive",
              count: d.employee_count || 0,
              manager: d.manager_name || ''
            })),
            total: res.total
          };
        }
        case "roles": {
          const res = await masterDataService.getRoles(page, limit);
          return {
            items: res.data.map((r: any) => ({
              ...r,
              id: String(r.id),
              status: r.is_active !== false ? "Active" : "Inactive",
              count: r.user_count || 0
            })),
            total: res.total
          };
        }
        case "warehouses": {
          const res = await masterDataService.getWarehouses(page, limit);
          return {
            items: res.data.map((w: any) => ({
              ...w,
              id: String(w.id),
              description: [w.city, w.state, w.country].filter(Boolean).join(', ') || w.address || 'N/A',
              status: w.is_active ? "Active" : "Inactive",
              count: 0, // Capacity is not in the API response
              location: [w.city, w.state].filter(Boolean).join(', '),
              capacity: '0'
            })),
            total: res.total
          };
        }
        default:
          return { items: [], total: 0 };
      }
    },
    placeholderData: keepPreviousData
  });

  // Mutations
  const createMutation = useMutation({
    mutationFn: async (values: any) => {
      let result;
      switch (type) {
        case "departments": result = await masterDataService.createDepartment(values); break;
        case "roles": result = await masterDataService.createRole(values); break;
        case "warehouses": result = await masterDataService.createWarehouse(values); break;
      }
      // Log audit
      await auditLogService.logAction(title.slice(0, -1), result?.id || 'unknown', 'CREATE', `Created new ${title.slice(0, -1)}: ${values.name}`);
      return result;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['masterData', type] });
      setIsDialogOpen(false);
      toast.success(`${title.slice(0, -1)} created successfully`);
    },
    onError: () => toast.error(`Failed to create ${title.slice(0, -1)}`)
  });

  const updateMutation = useMutation({
    mutationFn: async (values: any) => {
      if (!editingItem) return;
      let result;
      switch (type) {
        case "departments": result = await masterDataService.updateDepartment(editingItem.id, values); break;
        case "roles": result = await masterDataService.updateRole(editingItem.id, values); break;
        case "warehouses": result = await masterDataService.updateWarehouse(editingItem.id, values); break;
      }
      // Log audit
      await auditLogService.logAction(title.slice(0, -1), editingItem.id, 'UPDATE', `Updated ${title.slice(0, -1)}: ${values.name}`);
      return result;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['masterData', type] });
      setIsDialogOpen(false);
      setEditingItem(null);
      toast.success(`${title.slice(0, -1)} updated successfully`);
    },
    onError: () => toast.error(`Failed to update ${title.slice(0, -1)}`)
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      switch (type) {
        case "departments": await masterDataService.deleteDepartment(id); break;
        case "roles": await masterDataService.deleteRole(id); break;
        case "warehouses": await masterDataService.deleteWarehouse(id); break;
      }
      // Log audit
      await auditLogService.logAction(title.slice(0, -1), id, 'DELETE', `Deleted ${title.slice(0, -1)}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['masterData', type] });
      setIsDeleteDialogOpen(false);
      setItemToDelete(null);
      toast.success(`${title.slice(0, -1)} deleted successfully`);
    },
    onError: () => toast.error(`Failed to delete ${title.slice(0, -1)}`)
  });

  const bulkDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => {
      // In a real app, we'd have a bulk delete endpoint
      // For now, we'll execute them in parallel
      const promises = ids.map(id => {
        switch (type) {
          case "departments": return masterDataService.deleteDepartment(id);
          case "roles": return masterDataService.deleteRole(id);
          case "warehouses": return masterDataService.deleteWarehouse(id);
          default: return Promise.resolve();
        }
      });
      await Promise.all(promises);
      
      // Log audit for each
      for (const id of ids) {
        await auditLogService.logAction(title.slice(0, -1), id, 'DELETE', `Bulk deleted ${title.slice(0, -1)}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['masterData', type] });
      setIsBulkDeleteDialogOpen(false);
      setRowSelection({});
      toast.success(`${Object.keys(rowSelection).length} items deleted successfully`);
    },
    onError: () => toast.error(`Failed to delete selected items`)
  });

  // Form Handling
  const form = useForm({
    resolver: zodResolver(
      type === "departments" ? departmentSchema :
      type === "roles" ? roleSchema :
      warehouseSchema
    ),
    defaultValues: {
      name: "",
      code: "",
      manager: "",
      description: "",
      location: "",
      capacity: "",
      permissions: [] as string[],
    }
  });

  const onSubmit = (values: any) => {
    if (editingItem) {
      updateMutation.mutate(values);
    } else {
      createMutation.mutate(values);
    }
  };

  const handleEdit = (item: MasterData) => {
    setEditingItem(item);
    form.reset({
      name: item.name,
      code: item.code || "",
      manager: item.manager || "",
      description: item.description || "",
      location: item.location || "",
      capacity: item.capacity || "",
      permissions: item.permissions || [],
    });
    setIsDialogOpen(true);
  };

  const handleCreate = () => {
    setEditingItem(null);
    form.reset({
      name: "",
      code: "",
      manager: "",
      description: "",
      location: "",
      capacity: "",
      permissions: [],
    });
    setIsDialogOpen(true);
  };

  const handleDeleteClick = (item: MasterData) => {
    setItemToDelete(item);
    setIsDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (itemToDelete) {
      deleteMutation.mutate(itemToDelete.id);
    }
  };

  const confirmBulkDelete = () => {
    const ids = Object.keys(rowSelection);
    if (ids.length > 0) {
      bulkDeleteMutation.mutate(ids);
    }
  };

  const handleExport = () => {
    if (!data?.items) return;
    
    const exportData = data.items.map((item: any) => {
      // Customize export data based on type
      if (type === 'departments') {
        return {
          Name: item.name,
          Code: item.code,
          Manager: item.manager,
          'Employee Count': item.count,
          Status: item.status
        };
      } else if (type === 'roles') {
        return {
          Name: item.name,
          Description: item.description,
          'User Count': item.count,
          Status: item.status
        };
      } else {
        return {
          Name: item.name,
          Location: item.location,
          Capacity: item.capacity,
          Manager: item.manager,
          Status: item.status
        };
      }
    });
    
    exportService.exportToCSV(exportData, `${type}_export_${new Date().toISOString().split('T')[0]}`);
    toast.success('Export started');
  };

  const columns: ColumnDef<MasterData>[] = [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && "indeterminate")}
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
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
      cell: ({ row }) => <div className="font-medium">{row.getValue("name")}</div>,
    },
    {
      accessorKey: "description",
      header: "Description",
    },
    {
      accessorKey: "count",
      header: () => <div className="text-right">{type === "warehouses" ? "Capacity" : "Users"}</div>,
      cell: ({ row }) => <div className="text-right">{row.getValue("count")}</div>,
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
          <Badge variant={status === "Active" ? "default" : "secondary"}>
            {status || "Active"}
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
              <DropdownMenuItem onClick={() => handleEdit(row.original)}>
                <Pencil className="mr-2 h-4 w-4" /> Edit
              </DropdownMenuItem>
              {type === "roles" && <DropdownMenuItem>Permissions</DropdownMenuItem>}
              <DropdownMenuSeparator />
              <DropdownMenuItem 
                className="text-destructive"
                onClick={() => handleDeleteClick(row.original)}
              >
                <Trash2 className="mr-2 h-4 w-4" /> Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const table = useReactTable({
    data: (data?.items as MasterData[]) || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    onPaginationChange: setPagination,
    onRowSelectionChange: setRowSelection,
    manualPagination: true,
    rowCount: data?.total || 0,
    getRowId: (row) => row.id, // Use ID for selection
    state: {
      sorting,
      columnFilters,
      pagination,
      rowSelection,
    },
  });

  const getIcon = () => {
    switch (type) {
      case 'departments': return <Building2 className="h-6 w-6" />;
      case 'roles': return <Shield className="h-6 w-6" />;
      case 'warehouses': return <Warehouse className="h-6 w-6" />;
      default: return <Building2 className="h-6 w-6" />;
    }
  };

  if (isLoading && !data) return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  if (error) return <div className="p-8 text-center text-red-500">Failed to load {type}</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-lg text-primary">
            {getIcon()}
          </div>
          <div>
            <h1 className="text-3xl font-bold tracking-tight capitalize">{title}</h1>
            <p className="text-muted-foreground">
              Manage system {type} and configurations.
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <AuditLogViewer entityType={title.slice(0, -1)} />
          
          <Button variant="outline" onClick={handleExport}>
            <Download className="mr-2 h-4 w-4" />
            Export
          </Button>
          
          {Object.keys(rowSelection).length > 0 && (
            <Button 
              variant="destructive" 
              onClick={() => setIsBulkDeleteDialogOpen(true)}
            >
              <Trash2 className="mr-2 h-4 w-4" /> 
              Delete Selected ({Object.keys(rowSelection).length})
            </Button>
          )}
          
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button onClick={handleCreate}>
                <Plus className="mr-2 h-4 w-4" /> Add {title.slice(0, -1)}
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>{editingItem ? 'Edit' : 'Add'} {title.slice(0, -1)}</DialogTitle>
                <DialogDescription>
                  {editingItem ? 'Update the details below.' : 'Fill in the details to create a new item.'}
                </DialogDescription>
              </DialogHeader>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                  {type === "departments" && (
                    <>
                      <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Name</FormLabel>
                            <FormControl>
                              <Input placeholder="Department Name" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="code"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Code</FormLabel>
                            <FormControl>
                              <Input placeholder="DEPT-001" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="manager"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Manager</FormLabel>
                            <FormControl>
                              <Input placeholder="Manager Name" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </>
                  )}

                  {type === "roles" && (
                    <>
                      <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Role Name</FormLabel>
                            <FormControl>
                              <Input placeholder="Admin" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="description"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Description</FormLabel>
                            <FormControl>
                              <Input placeholder="Role description" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="permissions"
                        render={() => (
                          <FormItem>
                            <div className="mb-4">
                              <FormLabel className="text-base">Permissions</FormLabel>
                              <DialogDescription>
                                Select the permissions for this role.
                              </DialogDescription>
                            </div>
                            <div className="grid grid-cols-2 gap-4 border rounded-md p-4">
                              {AVAILABLE_PERMISSIONS.map((permission) => (
                                <FormField
                                  key={permission.id}
                                  control={form.control}
                                  name="permissions"
                                  render={({ field }) => {
                                    return (
                                      <FormItem
                                        key={permission.id}
                                        className="flex flex-row items-start space-x-3 space-y-0"
                                      >
                                        <FormControl>
                                          <Checkbox
                                            checked={field.value?.includes(permission.id)}
                                            onCheckedChange={(checked) => {
                                              return checked
                                                ? field.onChange([...(field.value || []), permission.id])
                                                : field.onChange(
                                                    field.value?.filter(
                                                      (value) => value !== permission.id
                                                    )
                                                  )
                                            }}
                                          />
                                        </FormControl>
                                        <FormLabel className="font-normal">
                                          {permission.label}
                                        </FormLabel>
                                      </FormItem>
                                    )
                                  }}
                                />
                              ))}
                            </div>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </>
                  )}

                  {type === "warehouses" && (
                    <>
                      <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Warehouse Name</FormLabel>
                            <FormControl>
                              <Input placeholder="Main Warehouse" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="location"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Location</FormLabel>
                            <FormControl>
                              <Input placeholder="123 Main St, City" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="capacity"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Capacity</FormLabel>
                            <FormControl>
                              <Input placeholder="10000" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="manager"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Manager</FormLabel>
                            <FormControl>
                              <Input placeholder="Manager Name" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </>
                  )}

                  <DialogFooter>
                    <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                      {createMutation.isPending || updateMutation.isPending ? (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      ) : null}
                      {editingItem ? 'Update' : 'Create'}
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </DialogContent>
          </Dialog>
        </div>

        <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-destructive">
                <AlertTriangle className="h-5 w-5" />
                Confirm Deletion
              </DialogTitle>
              <DialogDescription>
                Are you sure you want to delete <span className="font-semibold text-foreground">{itemToDelete?.name}</span>? 
                This action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>Cancel</Button>
              <Button 
                variant="destructive" 
                onClick={confirmDelete}
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={isBulkDeleteDialogOpen} onOpenChange={setIsBulkDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-destructive">
                <AlertTriangle className="h-5 w-5" />
                Confirm Bulk Deletion
              </DialogTitle>
              <DialogDescription>
                Are you sure you want to delete <span className="font-semibold text-foreground">{Object.keys(rowSelection).length}</span> items? 
                This action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsBulkDeleteDialogOpen(false)}>Cancel</Button>
              <Button 
                variant="destructive" 
                onClick={confirmBulkDelete}
                disabled={bulkDeleteMutation.isPending}
              >
                {bulkDeleteMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Delete All
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="flex items-center py-4 gap-2">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder={`Filter ${type}...`}
            value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
            onChange={(event) =>
              table.getColumn("name")?.setFilterValue(event.target.value)
            }
            className="pl-9"
          />
        </div>
      </div>

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      
      <div className="flex items-center justify-between py-4">
        <div className="flex items-center gap-2">
          <p className="text-sm text-muted-foreground">Rows per page</p>
          <Select
            value={`${table.getState().pagination.pageSize}`}
            onValueChange={(value) => {
              table.setPageSize(Number(value));
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={table.getState().pagination.pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 30, 40, 50].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        
        <div className="flex items-center space-x-2">
          <div className="text-sm text-muted-foreground mr-4">
            Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
            {Object.keys(rowSelection).length > 0 && (
              <span className="ml-2 text-primary">
                ({Object.keys(rowSelection).length} selected)
              </span>
            )}
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            <ChevronLeft className="h-4 w-4" />
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
