import { useState } from "react";
import { 
  ColumnDef, 
} from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown, 
  Tags,
  Scale,
  Loader2,
  Edit,
  Package,
  Trash2,
  X
} from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { productService } from "@/services/productService";
import { toast } from "sonner";
import { Link } from "wouter";

type Product = {
  id: number;
  sku: string;
  name: string;
  description?: string;
  category_id?: number;
  base_unit: string;
  is_catch_weight: boolean;
  catch_weight_unit?: string;
  shelf_life_days?: number;
  country_of_origin?: string;
  is_lot_tracked?: boolean;
  qc_required?: boolean;
  is_active: boolean;
};

type ProductUnit = {
  id: number;
  product_id: number;
  unit_name: string;
  description?: string;
  conversion_factor: number;
  barcode?: string;
  weight?: number;
  is_purchase_unit: boolean;
  is_sales_unit: boolean;
};

export default function ProductList() {
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isUnitsDialogOpen, setIsUnitsDialogOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [formData, setFormData] = useState({
    sku: "",
    name: "",
    description: "",
    base_unit: "EA",
    is_catch_weight: false,
    catch_weight_unit: "",
    shelf_life_days: 0,
    country_of_origin: "",
    is_lot_tracked: false,
    qc_required: false,
  });
  const [unitFormData, setUnitFormData] = useState({
    unit_name: "",
    description: "",
    conversion_factor: 1,
    barcode: "",
    weight: 0,
    is_purchase_unit: false,
    is_sales_unit: true,
  });
  const [editingUnit, setEditingUnit] = useState<ProductUnit | null>(null);

  const { data: productsData = [], isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => productService.getProducts(),
  });

  const { data: productUnits = [], refetch: refetchUnits } = useQuery({
    queryKey: ['product-units', selectedProduct?.id],
    queryFn: () => selectedProduct ? productService.getProductUnits(String(selectedProduct.id)) : [],
    enabled: !!selectedProduct,
  });

  // Ensure products is always an array
  const products = Array.isArray(productsData) ? productsData : [];
  const units = Array.isArray(productUnits) ? productUnits : [];

  const createMutation = useMutation({
    mutationFn: (data: any) => productService.createProduct(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      setIsDialogOpen(false);
      resetForm();
      toast.success("Product created successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to create product");
    }
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => productService.updateProduct(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      setIsDialogOpen(false);
      setEditingProduct(null);
      resetForm();
      toast.success("Product updated successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to update product");
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => productService.deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      toast.success("Product deactivated successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to deactivate product");
    }
  });

  const reactivateMutation = useMutation({
    mutationFn: (id: string) => productService.updateProduct(id, { is_active: true }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      toast.success("Product reactivated successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to reactivate product");
    }
  });

  const addUnitMutation = useMutation({
    mutationFn: (data: any) => productService.addUnit(data),
    onSuccess: () => {
      refetchUnits();
      setUnitFormData({
        unit_name: "",
        description: "",
        conversion_factor: 1,
        barcode: "",
        weight: 0,
        is_purchase_unit: false,
        is_sales_unit: true,
      });
      toast.success("Unit added successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to add unit");
    }
  });

  const deleteUnitMutation = useMutation({
    mutationFn: (id: string) => productService.deleteUnit(id),
    onSuccess: () => {
      refetchUnits();
      toast.success("Unit deleted successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to delete unit");
    }
  });

  const updateUnitMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => productService.updateUnit(id, data),
    onSuccess: () => {
      refetchUnits();
      setEditingUnit(null);
      resetUnitForm();
      toast.success("Unit updated successfully");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to update unit");
    }
  });

  const resetForm = () => {
    setFormData({
      sku: "",
      name: "",
      description: "",
      base_unit: "EA",
      is_catch_weight: false,
      catch_weight_unit: "",
      shelf_life_days: 0,
      country_of_origin: "",
      is_lot_tracked: false,
      qc_required: false,
    });
  };

  const handleEdit = (product: Product) => {
    setEditingProduct(product);
    setFormData({
      sku: product.sku,
      name: product.name,
      description: product.description || "",
      base_unit: product.base_unit,
      is_catch_weight: product.is_catch_weight,
      catch_weight_unit: product.catch_weight_unit || "",
      shelf_life_days: product.shelf_life_days || 0,
      country_of_origin: product.country_of_origin || "",
      is_lot_tracked: product.is_lot_tracked || false,
      qc_required: product.qc_required || false,
    });
    setIsDialogOpen(true);
  };

  const handleManageUnits = (product: Product) => {
    setSelectedProduct(product);
    setIsUnitsDialogOpen(true);
  };

  const handleSubmit = () => {
    if (editingProduct) {
      // Backend accepts these fields for update
      const updateData = {
        name: formData.name,
        description: formData.description || null,
        base_unit: formData.base_unit,
        shelf_life_days: formData.shelf_life_days || null,
        country_of_origin: formData.country_of_origin || null,
        qc_required: formData.qc_required,
      };
      console.log('Sending update:', updateData);
      updateMutation.mutate({ id: String(editingProduct.id), data: updateData });
    } else {
      createMutation.mutate(formData);
    }
  };

  const handleAddUnit = () => {
    if (!selectedProduct) return;
    addUnitMutation.mutate({
      product_id: selectedProduct.id,
      ...unitFormData,
    });
  };

  const resetUnitForm = () => {
    setUnitFormData({
      unit_name: "",
      description: "",
      conversion_factor: 1,
      barcode: "",
      weight: 0,
      is_purchase_unit: false,
      is_sales_unit: true,
    });
  };

  const handleEditUnit = (unit: ProductUnit) => {
    setEditingUnit(unit);
    setUnitFormData({
      unit_name: unit.unit_name,
      description: unit.description || "",
      conversion_factor: unit.conversion_factor,
      barcode: unit.barcode || "",
      weight: unit.weight || 0,
      is_purchase_unit: unit.is_purchase_unit,
      is_sales_unit: unit.is_sales_unit,
    });
  };

  const handleSaveUnit = () => {
    if (editingUnit) {
      updateUnitMutation.mutate({
        id: String(editingUnit.id),
        data: unitFormData,
      });
    } else {
      handleAddUnit();
    }
  };

  const handleCancelEditUnit = () => {
    setEditingUnit(null);
    resetUnitForm();
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    setEditingProduct(null);
    resetForm();
  };

  const columns: ColumnDef<Product>[] = [
    {
      accessorKey: "sku",
      header: "SKU",
      cell: ({ row }) => <div className="font-mono text-xs">{row.getValue("sku")}</div>,
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
      cell: ({ row }) => <div className="text-muted-foreground truncate max-w-[200px]">{row.getValue("description") || "-"}</div>,
    },
    {
      accessorKey: "base_unit",
      header: "Unit",
      cell: ({ row }) => <Badge variant="outline">{row.getValue("base_unit") || "EA"}</Badge>,
    },
    {
      accessorKey: "is_catch_weight",
      header: "Type",
      cell: ({ row }) => (
        row.getValue("is_catch_weight") ? 
        <Badge variant="secondary" className="flex w-fit items-center gap-1"><Scale className="h-3 w-3" /> Catch Weight</Badge> : 
        <Badge variant="outline">Standard</Badge>
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
        const product = row.original;
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
              <DropdownMenuItem onClick={() => handleEdit(product)}>
                <Edit className="mr-2 h-4 w-4" /> Edit Product
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleManageUnits(product)}>
                <Package className="mr-2 h-4 w-4" /> Manage Units
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href={`/inventory?product_id=${product.id}`}>
                  <Package className="mr-2 h-4 w-4" /> View Inventory
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {product.is_active ? (
                <DropdownMenuItem 
                  className="text-destructive"
                  onClick={() => {
                    if (confirm("Are you sure you want to deactivate this product?")) {
                      deleteMutation.mutate(String(product.id));
                    }
                  }}
                >
                  <Trash2 className="mr-2 h-4 w-4" /> Deactivate
                </DropdownMenuItem>
              ) : (
                <DropdownMenuItem 
                  className="text-green-600"
                  onClick={() => {
                    if (confirm("Reactivate this product?")) {
                      reactivateMutation.mutate(String(product.id));
                    }
                  }}
                >
                  <Plus className="mr-2 h-4 w-4" /> Reactivate
                </DropdownMenuItem>
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
          <h1 className="text-3xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">
            Manage your product catalog, categories, and units of measure.
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <Tags className="mr-2 h-4 w-4" /> Categories
          </Button>
          <Dialog open={isDialogOpen} onOpenChange={(open) => {
            if (!open) handleCloseDialog();
            else setIsDialogOpen(true);
          }}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 h-4 w-4" /> Add Product
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>{editingProduct ? "Edit Product" : "Add New Product"}</DialogTitle>
                <DialogDescription>
                  {editingProduct ? "Update product details." : "Create a new product in your catalog."}
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="sku">SKU *</Label>
                    <Input
                      id="sku"
                      placeholder="PROD001"
                      value={formData.sku}
                      onChange={(e) => setFormData({...formData, sku: e.target.value})}
                      disabled={!!editingProduct}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="base_unit">Base Unit *</Label>
                    <Input
                      id="base_unit"
                      placeholder="EA, LB, KG, CS"
                      value={formData.base_unit}
                      onChange={(e) => setFormData({...formData, base_unit: e.target.value})}
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="name">Product Name *</Label>
                  <Input
                    id="name"
                    placeholder="Product Name"
                    value={formData.name}
                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Input
                    id="description"
                    placeholder="Product description"
                    value={formData.description}
                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="shelf_life_days">Shelf Life (Days)</Label>
                    <Input
                      id="shelf_life_days"
                      type="number"
                      placeholder="30"
                      value={formData.shelf_life_days || ""}
                      onChange={(e) => setFormData({...formData, shelf_life_days: parseInt(e.target.value) || 0})}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="country_of_origin">Country of Origin</Label>
                    <Input
                      id="country_of_origin"
                      placeholder="USA"
                      maxLength={3}
                      value={formData.country_of_origin}
                      onChange={(e) => setFormData({...formData, country_of_origin: e.target.value.toUpperCase()})}
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="is_catch_weight"
                      checked={formData.is_catch_weight}
                      onCheckedChange={(checked) => setFormData({...formData, is_catch_weight: checked as boolean})}
                    />
                    <Label htmlFor="is_catch_weight">Catch Weight Item</Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="is_lot_tracked"
                      checked={formData.is_lot_tracked}
                      onCheckedChange={(checked) => setFormData({...formData, is_lot_tracked: checked as boolean})}
                    />
                    <Label htmlFor="is_lot_tracked">Lot Tracked</Label>
                  </div>
                </div>
                {formData.is_catch_weight && (
                  <div className="space-y-2">
                    <Label htmlFor="catch_weight_unit">Catch Weight Unit *</Label>
                    <Input
                      id="catch_weight_unit"
                      placeholder="LB, KG"
                      value={formData.catch_weight_unit}
                      onChange={(e) => setFormData({...formData, catch_weight_unit: e.target.value})}
                    />
                  </div>
                )}
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="qc_required"
                    checked={formData.qc_required}
                    onCheckedChange={(checked) => setFormData({...formData, qc_required: checked as boolean})}
                  />
                  <Label htmlFor="qc_required">QC Required</Label>
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={handleCloseDialog}>
                  Cancel
                </Button>
                <Button onClick={handleSubmit} disabled={createMutation.isPending || updateMutation.isPending}>
                  {(createMutation.isPending || updateMutation.isPending) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  {editingProduct ? "Update Product" : "Create Product"}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Units Management Dialog */}
      <Dialog open={isUnitsDialogOpen} onOpenChange={(open) => {
        setIsUnitsDialogOpen(open);
        if (!open) {
          setEditingUnit(null);
          resetUnitForm();
        }
      }}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>Manage Units - {selectedProduct?.name}</DialogTitle>
            <DialogDescription>
              Add and manage units of measure for this product. Base unit: {selectedProduct?.base_unit}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            {/* Add/Edit Unit Form */}
            <div className="border p-3 rounded-lg bg-muted/50">
              <div className="flex items-center justify-between mb-2">
                <Label className="text-sm font-medium">
                  {editingUnit ? `Editing: ${editingUnit.unit_name}` : "Add New Unit"}
                </Label>
                {editingUnit && (
                  <Button variant="ghost" size="sm" onClick={handleCancelEditUnit}>
                    <X className="h-4 w-4 mr-1" /> Cancel
                  </Button>
                )}
              </div>
              <div className="grid grid-cols-6 gap-2 items-end">
                <div className="space-y-1">
                  <Label className="text-xs">Unit Name</Label>
                  <Input
                    placeholder="CS"
                    value={unitFormData.unit_name}
                    onChange={(e) => setUnitFormData({...unitFormData, unit_name: e.target.value})}
                  />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Conversion</Label>
                  <Input
                    type="number"
                    step="0.01"
                    placeholder="12"
                    value={unitFormData.conversion_factor}
                    onChange={(e) => setUnitFormData({...unitFormData, conversion_factor: parseFloat(e.target.value) || 1})}
                  />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Barcode</Label>
                  <Input
                    placeholder="Barcode"
                    value={unitFormData.barcode}
                    onChange={(e) => setUnitFormData({...unitFormData, barcode: e.target.value})}
                  />
                </div>
                <div className="flex items-center space-x-1">
                  <Checkbox
                    id="is_purchase_unit"
                    checked={unitFormData.is_purchase_unit}
                    onCheckedChange={(checked) => setUnitFormData({...unitFormData, is_purchase_unit: checked as boolean})}
                  />
                  <Label htmlFor="is_purchase_unit" className="text-xs">Buy</Label>
                </div>
                <div className="flex items-center space-x-1">
                  <Checkbox
                    id="is_sales_unit"
                    checked={unitFormData.is_sales_unit}
                    onCheckedChange={(checked) => setUnitFormData({...unitFormData, is_sales_unit: checked as boolean})}
                  />
                  <Label htmlFor="is_sales_unit" className="text-xs">Sell</Label>
                </div>
                <Button 
                  size="sm" 
                  onClick={handleSaveUnit} 
                  disabled={addUnitMutation.isPending || updateUnitMutation.isPending}
                  variant={editingUnit ? "default" : "outline"}
                >
                  {(addUnitMutation.isPending || updateUnitMutation.isPending) ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : editingUnit ? (
                    <>Save</>
                  ) : (
                    <Plus className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>

            {/* Units Table */}
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Unit</TableHead>
                  <TableHead>Conversion Factor</TableHead>
                  <TableHead>Barcode</TableHead>
                  <TableHead>Purchase</TableHead>
                  <TableHead>Sales</TableHead>
                  <TableHead className="w-[80px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {units.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                      No additional units defined. Add units above.
                    </TableCell>
                  </TableRow>
                ) : (
                  units.map((unit: ProductUnit) => (
                    <TableRow key={unit.id} className={editingUnit?.id === unit.id ? "bg-muted" : ""}>
                      <TableCell className="font-medium">{unit.unit_name}</TableCell>
                      <TableCell>{unit.conversion_factor} {selectedProduct?.base_unit}</TableCell>
                      <TableCell className="font-mono text-xs">{unit.barcode || "-"}</TableCell>
                      <TableCell>{unit.is_purchase_unit ? <Badge>Yes</Badge> : "-"}</TableCell>
                      <TableCell>{unit.is_sales_unit ? <Badge>Yes</Badge> : "-"}</TableCell>
                      <TableCell>
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleEditUnit(unit)}
                            disabled={editingUnit?.id === unit.id}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              if (confirm("Delete this unit?")) {
                                deleteUnitMutation.mutate(String(unit.id));
                              }
                            }}
                          >
                            <X className="h-4 w-4 text-destructive" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsUnitsDialogOpen(false)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : (
        <DataTable 
          columns={columns} 
          data={products} 
          searchKey="name"
          searchPlaceholder="Filter products..."
        />
      )}
    </div>
  );
}
