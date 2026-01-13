import { useState } from "react";
import { Link } from "wouter";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { DataTable } from "@/components/ui/data-table";
import { Card, CardContent } from "@/components/ui/card";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown,
  Building2,
  Loader2,
  Edit,
  Trash2,
  Eye,
  DollarSign,
  FileText,
  Package
} from "lucide-react";
import { toast } from "sonner";
import { masterDataService, Vendor } from "@/services/masterDataService";

export default function VendorList() {
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingVendor, setEditingVendor] = useState<Vendor | null>(null);
  const [formData, setFormData] = useState({
    vendor_code: "",
    name: "",
    email: "",
    phone: "",
    payment_terms_days: 30,
    is_active: true
  });

  const { data: vendorsData, isLoading, error } = useQuery({
    queryKey: ['vendors'],
    queryFn: () => masterDataService.getVendors(1, 100)
  });

  // Ensure vendors is always an array to prevent .filter errors
  const rawVendors = vendorsData?.data;
  const vendors = Array.isArray(rawVendors) ? rawVendors : [];

  const createMutation = useMutation({
    mutationFn: (data: any) => masterDataService.createVendor(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendors'] });
      setIsDialogOpen(false);
      resetForm();
      toast.success("Vendor created successfully");
    },
    onError: () => {
      toast.error("Failed to create vendor");
    }
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => masterDataService.updateVendor(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendors'] });
      setIsDialogOpen(false);
      setEditingVendor(null);
      resetForm();
      toast.success("Vendor updated successfully");
    },
    onError: () => {
      toast.error("Failed to update vendor");
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => masterDataService.deleteVendor(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendors'] });
      toast.success("Vendor deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete vendor");
    }
  });

  const resetForm = () => {
    setFormData({
      vendor_code: "",
      name: "",
      email: "",
      phone: "",
      payment_terms_days: 30,
      is_active: true
    });
  };

  const handleEdit = (vendor: Vendor) => {
    setEditingVendor(vendor);
    setFormData({
      vendor_code: vendor.vendor_code,
      name: vendor.name,
      email: vendor.email || "",
      phone: vendor.phone || "",
      payment_terms_days: vendor.payment_terms_days,
      is_active: vendor.is_active
    });
    setIsDialogOpen(true);
  };

  const handleSubmit = () => {
    if (editingVendor) {
      updateMutation.mutate({ id: String(editingVendor.id), data: formData });
    } else {
      createMutation.mutate(formData);
    }
  };

  const columns: ColumnDef<Vendor>[] = [
    {
      accessorKey: "vendor_code",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Code
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => (
        <div className="font-mono font-medium">{row.getValue("vendor_code")}</div>
      ),
    },
    {
      accessorKey: "name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Vendor Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => (
        <div className="font-medium">{row.getValue("name")}</div>
      ),
    },
    {
      accessorKey: "email",
      header: "Email",
      cell: ({ row }) => (
        <div className="text-muted-foreground">{row.getValue("email") || "-"}</div>
      ),
    },
    {
      accessorKey: "phone",
      header: "Phone",
      cell: ({ row }) => (
        <div className="text-muted-foreground">{row.getValue("phone") || "-"}</div>
      ),
    },
    {
      accessorKey: "payment_terms_days",
      header: "Payment Terms",
      cell: ({ row }) => (
        <Badge variant="outline">{row.getValue("payment_terms_days")} days</Badge>
      ),
    },
    {
      accessorKey: "is_active",
      header: "Status",
      cell: ({ row }) => {
        const isActive = row.getValue("is_active") as boolean;
        return (
          <Badge variant={isActive ? "default" : "secondary"}>
            {isActive ? "Active" : "Inactive"}
          </Badge>
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const vendor = row.original;
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
              <DropdownMenuSeparator />
              <Link href={`/vendors/${vendor.id}`}>
                <DropdownMenuItem>
                  <Eye className="mr-2 h-4 w-4" /> View Details
                </DropdownMenuItem>
              </Link>
              <DropdownMenuItem onClick={() => handleEdit(vendor)}>
                <Edit className="mr-2 h-4 w-4" /> Edit
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <Link href={`/ap/aging/${vendor.id}`}>
                <DropdownMenuItem>
                  <DollarSign className="mr-2 h-4 w-4" /> View AP
                </DropdownMenuItem>
              </Link>
              <Link href={`/purchase-orders?vendor=${vendor.id}`}>
                <DropdownMenuItem>
                  <FileText className="mr-2 h-4 w-4" /> Purchase Orders
                </DropdownMenuItem>
              </Link>
              <Link href={`/products?vendor=${vendor.id}`}>
                <DropdownMenuItem>
                  <Package className="mr-2 h-4 w-4" /> Products
                </DropdownMenuItem>
              </Link>
              <DropdownMenuSeparator />
              <DropdownMenuItem 
                className="text-destructive"
                onClick={() => deleteMutation.mutate(String(vendor.id))}
              >
                <Trash2 className="mr-2 h-4 w-4" /> Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

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
        Failed to load vendors
      </div>
    );
  }

  const activeCount = vendors.filter((v: Vendor) => v.is_active).length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Vendors</h1>
          <p className="text-muted-foreground">
            Manage your vendor accounts and suppliers
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={(open) => {
          setIsDialogOpen(open);
          if (!open) {
            setEditingVendor(null);
            resetForm();
          }
        }}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" /> Add Vendor
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {editingVendor ? "Edit Vendor" : "Add New Vendor"}
              </DialogTitle>
              <DialogDescription>
                {editingVendor 
                  ? "Update vendor information."
                  : "Enter the details for the new vendor."}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="vendor_code">Vendor Code</Label>
                  <Input
                    id="vendor_code"
                    placeholder="VEND001"
                    value={formData.vendor_code}
                    onChange={(e) => setFormData({...formData, vendor_code: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="payment_terms">Payment Terms (days)</Label>
                  <Input
                    id="payment_terms"
                    type="number"
                    placeholder="30"
                    value={formData.payment_terms_days}
                    onChange={(e) => setFormData({...formData, payment_terms_days: parseInt(e.target.value) || 0})}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="name">Vendor Name</Label>
                <Input
                  id="name"
                  placeholder="ABC Suppliers Inc."
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="orders@supplier.com"
                    value={formData.email}
                    onChange={(e) => setFormData({...formData, email: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">Phone</Label>
                  <Input
                    id="phone"
                    placeholder="+1 (555) 000-0000"
                    value={formData.phone}
                    onChange={(e) => setFormData({...formData, phone: e.target.value})}
                  />
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                Cancel
              </Button>
              <Button 
                onClick={handleSubmit}
                disabled={createMutation.isPending || updateMutation.isPending}
              >
                {(createMutation.isPending || updateMutation.isPending) && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                {editingVendor ? "Update" : "Create"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Vendors</p>
                <p className="text-2xl font-bold">{vendors.length}</p>
              </div>
              <Building2 className="h-8 w-8 text-muted-foreground" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Vendors</p>
                <p className="text-2xl font-bold text-emerald-600">{activeCount}</p>
              </div>
              <Badge className="bg-emerald-100 text-emerald-800">Active</Badge>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Inactive Vendors</p>
                <p className="text-2xl font-bold text-gray-500">{vendors.length - activeCount}</p>
              </div>
              <Badge variant="secondary">Inactive</Badge>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Data Table */}
      <DataTable 
        columns={columns} 
        data={vendors} 
        searchKey="name"
        searchPlaceholder="Search vendors..."
      />
    </div>
  );
}
