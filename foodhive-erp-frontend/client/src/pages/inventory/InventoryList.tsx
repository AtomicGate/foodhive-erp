import { useState } from "react";
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
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  MoreHorizontal, 
  ArrowUpDown, 
  AlertTriangle,
  History,
  ArrowRightLeft,
  Package,
  Warehouse,
  TrendingDown,
  Clock,
  Plus,
  RefreshCw,
  Eye,
  Download,
  Filter,
  PackagePlus,
  CalendarDays,
  MapPin,
  DollarSign,
  Layers,
  ArrowDownRight,
  ArrowUpRight,
  RotateCcw
} from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { inventoryService } from "@/services/inventoryService";
import { toast } from "sonner";
import { Separator } from "@/components/ui/separator";
import { ScrollArea } from "@/components/ui/scroll-area";

// Types matching backend response
type InventoryItem = {
  inventory: {
    id: number;
    product_id: number;
    warehouse_id: number;
    location_code: string;
    lot_number: string;
    production_date: string;
    expiry_date: string;
    quantity_on_hand: number;
    quantity_allocated: number;
    quantity_on_order: number;
    quantity_available: number;
    last_cost: number;
    average_cost: number;
    created_at: string;
    updated_at: string;
  };
  product_name: string;
  product_sku: string;
  warehouse_name: string;
  days_to_expiry: number;
  age_in_days: number;
};

// Flattened type for table display
type FlatInventoryItem = {
  id: number;
  sku: string;
  productName: string;
  warehouseName: string;
  locationCode: string;
  lotNumber: string;
  onHand: number;
  allocated: number;
  available: number;
  onOrder: number;
  averageCost: number;
  lastCost: number;
  expiryDate: string;
  productionDate: string;
  daysToExpiry: number;
  ageInDays: number;
  status: "In Stock" | "Low Stock" | "Out of Stock" | "Expiring Soon";
  productId: number;
  warehouseId: number;
  inventoryValue: number;
};

// Transaction type
type InventoryTransaction = {
  id: number;
  product_id: number;
  warehouse_id: number;
  transaction_type: string;
  quantity: number;
  lot_number: string;
  unit_cost: number;
  reference_type: string;
  reference_number: string;
  notes: string;
  created_at: string;
};

// Stats type
type InventoryStats = {
  totalProducts: number;
  totalValue: number;
  lowStockItems: number;
  expiringItems: number;
  totalQuantity: number;
};

export default function InventoryList() {
  const queryClient = useQueryClient();
  
  // State for dialogs
  const [adjustDialogOpen, setAdjustDialogOpen] = useState(false);
  const [transferDialogOpen, setTransferDialogOpen] = useState(false);
  const [receiveDialogOpen, setReceiveDialogOpen] = useState(false);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  const [historyDialogOpen, setHistoryDialogOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<FlatInventoryItem | null>(null);
  
  // Filter state
  const [warehouseFilter, setWarehouseFilter] = useState<string>("all");
  const [activeTab, setActiveTab] = useState("all");
  
  // Adjustment form state
  const [adjustmentQty, setAdjustmentQty] = useState("");
  const [adjustmentReason, setAdjustmentReason] = useState("");
  const [adjustmentNotes, setAdjustmentNotes] = useState("");
  
  // Transfer form state
  const [transferQty, setTransferQty] = useState("");
  const [transferToWarehouse, setTransferToWarehouse] = useState("");
  const [transferNotes, setTransferNotes] = useState("");
  
  // Receive form state
  const [receiveProductId, setReceiveProductId] = useState("");
  const [receiveWarehouseId, setReceiveWarehouseId] = useState("");
  const [receiveQty, setReceiveQty] = useState("");
  const [receiveCost, setReceiveCost] = useState("");
  const [receiveLotNumber, setReceiveLotNumber] = useState("");
  const [receiveExpiryDate, setReceiveExpiryDate] = useState("");
  const [receiveNotes, setReceiveNotes] = useState("");

  // Fetch inventory data
  const { data: rawInventory = [], isLoading, refetch } = useQuery({
    queryKey: ['inventory'],
    queryFn: () => inventoryService.getInventory({ page_size: 500 }),
  });

  // Fetch warehouses for filter and transfer
  const { data: warehousesData = [] } = useQuery({
    queryKey: ['warehouses'],
    queryFn: async () => {
      try {
        const response = await fetch('/api/warehouses/list');
        const data = await response.json();
        return data.data || data || [];
      } catch {
        return [];
      }
    },
  });
  const warehouses = Array.isArray(warehousesData) ? warehousesData : [];

  // Fetch products for receive dialog
  const { data: productsData = [] } = useQuery({
    queryKey: ['products'],
    queryFn: async () => {
      try {
        const response = await fetch('/api/products/list');
        const data = await response.json();
        return data.data || data || [];
      } catch {
        return [];
      }
    },
  });
  const products = Array.isArray(productsData) ? productsData : [];

  // Fetch expiring inventory
  const { data: expiringItems = [] } = useQuery({
    queryKey: ['inventory-expiring'],
    queryFn: () => inventoryService.getExpiringInventory({ days: 7 }),
  });

  // Fetch transaction history for selected item
  const { data: transactionHistory = [], isLoading: isLoadingHistory } = useQuery({
    queryKey: ['inventory-transactions', selectedItem?.productId],
    queryFn: () => inventoryService.getTransactions({ 
      product_id: selectedItem?.productId,
      limit: 50
    }),
    enabled: !!selectedItem && historyDialogOpen,
  });

  // Adjust mutation
  const adjustMutation = useMutation({
    mutationFn: (data: any) => inventoryService.adjustInventory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
      setAdjustDialogOpen(false);
      resetAdjustForm();
      toast.success("Inventory adjusted successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to adjust inventory");
    },
  });

  // Transfer mutation
  const transferMutation = useMutation({
    mutationFn: (data: any) => inventoryService.transferInventory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
      setTransferDialogOpen(false);
      resetTransferForm();
      toast.success("Inventory transferred successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to transfer inventory");
    },
  });

  // Receive mutation
  const receiveMutation = useMutation({
    mutationFn: (data: any) => inventoryService.receiveInventory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
      setReceiveDialogOpen(false);
      resetReceiveForm();
      toast.success("Inventory received successfully");
    },
    onError: (error: any) => {
      toast.error(error.message || "Failed to receive inventory");
    },
  });

  // Transform backend data to flat format for table
  const inventory: FlatInventoryItem[] = (Array.isArray(rawInventory) ? rawInventory : rawInventory?.data || []).map((item: any) => {
    const inv = item.inventory || item;
    const available = inv.quantity_available ?? ((inv.quantity_on_hand || 0) - (inv.quantity_allocated || 0));
    const daysToExpiry = item.days_to_expiry ?? inv.days_to_expiry ?? 999;
    const avgCost = inv.average_cost || 0;
    const onHand = inv.quantity_on_hand || 0;
    
    let status: FlatInventoryItem["status"] = "In Stock";
    if (available <= 0) {
      status = "Out of Stock";
    } else if (daysToExpiry <= 7 && daysToExpiry >= 0) {
      status = "Expiring Soon";
    } else if (available < 20) {
      status = "Low Stock";
    }

    return {
      id: inv.id,
      sku: item.product_sku || inv.product_sku || 'N/A',
      productName: item.product_name || inv.product_name || 'Unknown Product',
      warehouseName: item.warehouse_name || inv.warehouse_name || 'Unknown',
      locationCode: inv.location_code || '',
      lotNumber: inv.lot_number || '',
      onHand: onHand,
      allocated: inv.quantity_allocated || 0,
      available: available || 0,
      onOrder: inv.quantity_on_order || 0,
      averageCost: avgCost,
      lastCost: inv.last_cost || 0,
      expiryDate: inv.expiry_date || '',
      productionDate: inv.production_date || '',
      daysToExpiry: daysToExpiry,
      ageInDays: item.age_in_days || inv.age_in_days || 0,
      status,
      productId: inv.product_id,
      warehouseId: inv.warehouse_id,
      inventoryValue: onHand * avgCost,
    };
  });

  // Filter inventory based on warehouse and tab
  const filteredInventory = inventory.filter(item => {
    // Warehouse filter
    if (warehouseFilter !== "all" && item.warehouseId.toString() !== warehouseFilter) {
      return false;
    }
    // Tab filter
    if (activeTab === "low-stock" && item.status !== "Low Stock" && item.status !== "Out of Stock") {
      return false;
    }
    if (activeTab === "expiring" && item.status !== "Expiring Soon") {
      return false;
    }
    return true;
  });

  // Calculate stats
  const stats: InventoryStats = {
    totalProducts: new Set(inventory.map(i => i.productId)).size,
    totalValue: inventory.reduce((sum, i) => sum + i.inventoryValue, 0),
    lowStockItems: inventory.filter(i => i.status === "Low Stock" || i.status === "Out of Stock").length,
    expiringItems: inventory.filter(i => i.status === "Expiring Soon").length,
    totalQuantity: inventory.reduce((sum, i) => sum + i.onHand, 0),
  };

  const resetAdjustForm = () => {
    setAdjustmentQty("");
    setAdjustmentReason("");
    setAdjustmentNotes("");
    setSelectedItem(null);
  };

  const resetTransferForm = () => {
    setTransferQty("");
    setTransferToWarehouse("");
    setTransferNotes("");
    setSelectedItem(null);
  };

  const resetReceiveForm = () => {
    setReceiveProductId("");
    setReceiveWarehouseId("");
    setReceiveQty("");
    setReceiveCost("");
    setReceiveLotNumber("");
    setReceiveExpiryDate("");
    setReceiveNotes("");
  };

  const handleViewDetails = (item: FlatInventoryItem) => {
    setSelectedItem(item);
    setDetailsDialogOpen(true);
  };

  const handleViewHistory = (item: FlatInventoryItem) => {
    setSelectedItem(item);
    setHistoryDialogOpen(true);
  };

  const handleAdjustStock = (item: FlatInventoryItem) => {
    setSelectedItem(item);
    setAdjustDialogOpen(true);
  };

  const handleTransferStock = (item: FlatInventoryItem) => {
    setSelectedItem(item);
    setTransferDialogOpen(true);
  };

  const submitAdjustment = () => {
    if (!selectedItem || !adjustmentQty || !adjustmentReason) return;
    
    adjustMutation.mutate({
      product_id: selectedItem.productId,
      warehouse_id: selectedItem.warehouseId,
      location_code: selectedItem.locationCode,
      lot_number: selectedItem.lotNumber,
      quantity: parseFloat(adjustmentQty),
      reason: adjustmentReason,
      notes: adjustmentNotes,
    });
  };

  const submitTransfer = () => {
    if (!selectedItem || !transferQty || !transferToWarehouse) return;
    
    transferMutation.mutate({
      product_id: selectedItem.productId,
      from_warehouse_id: selectedItem.warehouseId,
      to_warehouse_id: parseInt(transferToWarehouse),
      from_location_code: selectedItem.locationCode,
      lot_number: selectedItem.lotNumber,
      quantity: parseFloat(transferQty),
      notes: transferNotes,
    });
  };

  const submitReceive = () => {
    if (!receiveProductId || !receiveWarehouseId || !receiveQty) return;
    
    receiveMutation.mutate({
      product_id: parseInt(receiveProductId),
      warehouse_id: parseInt(receiveWarehouseId),
      quantity: parseFloat(receiveQty),
      unit_cost: parseFloat(receiveCost) || 0,
      lot_number: receiveLotNumber,
      expiry_date: receiveExpiryDate || null,
      notes: receiveNotes,
    });
  };

  const columns: ColumnDef<FlatInventoryItem>[] = [
    {
      accessorKey: "sku",
      header: "SKU",
      cell: ({ row }) => (
        <div className="font-mono text-xs font-semibold text-primary">
          {row.getValue("sku")}
        </div>
      ),
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
      cell: ({ row }) => (
        <div>
          <div className="font-medium">{row.getValue("productName")}</div>
          {row.original.lotNumber && (
            <div className="text-xs text-muted-foreground flex items-center gap-1">
              <Layers className="h-3 w-3" />
              Lot: {row.original.lotNumber}
            </div>
          )}
        </div>
      ),
    },
    {
      accessorKey: "warehouseName",
      header: "Location",
      cell: ({ row }) => (
        <div>
          <div className="font-medium flex items-center gap-1">
            <Warehouse className="h-3 w-3 text-muted-foreground" />
            {row.getValue("warehouseName")}
          </div>
          {row.original.locationCode && (
            <div className="text-xs text-muted-foreground flex items-center gap-1">
              <MapPin className="h-3 w-3" />
              {row.original.locationCode}
            </div>
          )}
        </div>
      ),
    },
    {
      accessorKey: "onHand",
      header: () => <div className="text-right">On Hand</div>,
      cell: ({ row }) => {
        const amount = row.getValue("onHand") as number;
        return (
          <div className="text-right font-medium tabular-nums">
            {amount.toLocaleString()}
          </div>
        );
      },
    },
    {
      accessorKey: "available",
      header: () => <div className="text-right">Available</div>,
      cell: ({ row }) => {
        const amount = row.getValue("available") as number;
        return (
          <div className={`text-right font-bold tabular-nums ${amount <= 0 ? 'text-destructive' : amount < 20 ? 'text-amber-600' : 'text-emerald-600'}`}>
            {amount.toLocaleString()}
          </div>
        );
      },
    },
    {
      accessorKey: "inventoryValue",
      header: () => <div className="text-right">Value</div>,
      cell: ({ row }) => {
        const value = row.getValue("inventoryValue") as number;
        return (
          <div className="text-right text-muted-foreground tabular-nums">
            ₭{value.toLocaleString()}
          </div>
        );
      },
    },
    {
      accessorKey: "expiryDate",
      header: "Expiry",
      cell: ({ row }) => {
        const expiry = row.original.expiryDate;
        const daysToExpiry = row.original.daysToExpiry;
        
        if (!expiry) return <span className="text-muted-foreground">-</span>;
        
        const date = new Date(expiry);
        const isExpiringSoon = daysToExpiry <= 7 && daysToExpiry >= 0;
        const isExpired = daysToExpiry < 0;
        
        return (
          <div className={`text-sm ${isExpired ? 'text-destructive' : isExpiringSoon ? 'text-orange-600' : ''}`}>
            <div>{date.toLocaleDateString()}</div>
            {isExpiringSoon && (
              <div className="text-xs font-medium">{daysToExpiry}d left</div>
            )}
            {isExpired && (
              <div className="text-xs font-medium">Expired!</div>
            )}
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
        let icon = null;
        
        switch (status) {
          case "In Stock":
            className = "bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-200";
            break;
          case "Low Stock":
            className = "bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200";
            icon = <TrendingDown className="h-3 w-3 mr-1" />;
            break;
          case "Out of Stock":
            variant = "destructive";
            break;
          case "Expiring Soon":
            className = "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200";
            icon = <Clock className="h-3 w-3 mr-1" />;
            break;
        }

        return (
          <Badge variant={variant} className={`${className} flex items-center w-fit`}>
            {icon}
            {status}
          </Badge>
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const item = row.original;
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
              <DropdownMenuItem onClick={() => handleViewDetails(item)}>
                <Eye className="h-4 w-4 mr-2" />
                View Details
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleViewHistory(item)}>
                <History className="h-4 w-4 mr-2" />
                View History
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => handleAdjustStock(item)}>
                <Plus className="h-4 w-4 mr-2" />
                Adjust Stock
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleTransferStock(item)}>
                <ArrowRightLeft className="h-4 w-4 mr-2" />
                Transfer Stock
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'LAK',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString();
  };

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case 'RECEIVE':
        return <ArrowDownRight className="h-4 w-4 text-emerald-500" />;
      case 'SHIP':
        return <ArrowUpRight className="h-4 w-4 text-blue-500" />;
      case 'ADJUST_IN':
        return <Plus className="h-4 w-4 text-emerald-500" />;
      case 'ADJUST_OUT':
        return <TrendingDown className="h-4 w-4 text-amber-500" />;
      case 'TRANSFER_IN':
        return <ArrowDownRight className="h-4 w-4 text-blue-500" />;
      case 'TRANSFER_OUT':
        return <ArrowUpRight className="h-4 w-4 text-purple-500" />;
      case 'RETURN':
        return <RotateCcw className="h-4 w-4 text-orange-500" />;
      default:
        return <History className="h-4 w-4 text-muted-foreground" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            Inventory Management
          </h1>
          <p className="text-muted-foreground">
            Track stock levels, adjustments, and transfers across warehouses.
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" onClick={() => refetch()}>
            <RefreshCw className="mr-2 h-4 w-4" /> Refresh
          </Button>
          <Button variant="outline" onClick={() => setReceiveDialogOpen(true)}>
            <PackagePlus className="mr-2 h-4 w-4" /> Receive Stock
          </Button>
          <Button variant="default">
            <Download className="mr-2 h-4 w-4" /> Export
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <Card className="border-l-4 border-l-primary">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Products</CardTitle>
            <Package className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalProducts}</div>
            <p className="text-xs text-muted-foreground">Unique SKUs in stock</p>
          </CardContent>
        </Card>
        
        <Card className="border-l-4 border-l-blue-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Quantity</CardTitle>
            <Layers className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalQuantity.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">Units in all warehouses</p>
          </CardContent>
        </Card>
        
        <Card className="border-l-4 border-l-emerald-500">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Value</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(stats.totalValue)}</div>
            <p className="text-xs text-muted-foreground">At average cost</p>
          </CardContent>
        </Card>
        
        <Card className="border-l-4 border-l-amber-500 cursor-pointer hover:shadow-md transition-shadow" onClick={() => setActiveTab("low-stock")}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Low Stock</CardTitle>
            <TrendingDown className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">{stats.lowStockItems}</div>
            <p className="text-xs text-muted-foreground">Items need reorder</p>
          </CardContent>
        </Card>
        
        <Card className="border-l-4 border-l-orange-500 cursor-pointer hover:shadow-md transition-shadow" onClick={() => setActiveTab("expiring")}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Expiring Soon</CardTitle>
            <AlertTriangle className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">{stats.expiringItems}</div>
            <p className="text-xs text-muted-foreground">Within 7 days</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters & Tabs */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full md:w-auto">
              <TabsList>
                <TabsTrigger value="all" className="gap-2">
                  <Package className="h-4 w-4" />
                  All Items
                  <Badge variant="secondary" className="ml-1">{inventory.length}</Badge>
                </TabsTrigger>
                <TabsTrigger value="low-stock" className="gap-2">
                  <TrendingDown className="h-4 w-4" />
                  Low Stock
                  <Badge variant="secondary" className="ml-1 bg-amber-100 text-amber-800">{stats.lowStockItems}</Badge>
                </TabsTrigger>
                <TabsTrigger value="expiring" className="gap-2">
                  <Clock className="h-4 w-4" />
                  Expiring
                  <Badge variant="secondary" className="ml-1 bg-orange-100 text-orange-800">{stats.expiringItems}</Badge>
                </TabsTrigger>
              </TabsList>
            </Tabs>
            
            <div className="flex items-center gap-2">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <Select value={warehouseFilter} onValueChange={setWarehouseFilter}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="All Warehouses" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Warehouses</SelectItem>
                  {warehouses.map((wh: any) => (
                    <SelectItem key={wh.id} value={wh.id.toString()}>
                      {wh.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : (
            <DataTable 
              columns={columns} 
              data={filteredInventory} 
              searchKey="productName"
              searchPlaceholder="Search products..."
            />
          )}
        </CardContent>
      </Card>

      {/* Details Dialog */}
      <Dialog open={detailsDialogOpen} onOpenChange={setDetailsDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Inventory Details
            </DialogTitle>
            <DialogDescription>
              Complete information for this inventory item
            </DialogDescription>
          </DialogHeader>
          {selectedItem && (
            <div className="space-y-6">
              {/* Product Info */}
              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-2">Product Information</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-xs text-muted-foreground">SKU</Label>
                    <p className="font-mono font-medium">{selectedItem.sku}</p>
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Product Name</Label>
                    <p className="font-medium">{selectedItem.productName}</p>
                  </div>
                </div>
              </div>
              
              <Separator />
              
              {/* Location Info */}
              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-2">Location</h4>
                <div className="grid grid-cols-3 gap-4">
                  <div>
                    <Label className="text-xs text-muted-foreground">Warehouse</Label>
                    <p className="font-medium">{selectedItem.warehouseName}</p>
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Location Code</Label>
                    <p className="font-medium">{selectedItem.locationCode || '-'}</p>
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Lot Number</Label>
                    <p className="font-medium">{selectedItem.lotNumber || '-'}</p>
                  </div>
                </div>
              </div>
              
              <Separator />
              
              {/* Quantities */}
              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-2">Quantities</h4>
                <div className="grid grid-cols-4 gap-4">
                  <div className="bg-muted/50 rounded-lg p-3 text-center">
                    <p className="text-2xl font-bold">{selectedItem.onHand}</p>
                    <p className="text-xs text-muted-foreground">On Hand</p>
                  </div>
                  <div className="bg-muted/50 rounded-lg p-3 text-center">
                    <p className="text-2xl font-bold text-amber-600">{selectedItem.allocated}</p>
                    <p className="text-xs text-muted-foreground">Allocated</p>
                  </div>
                  <div className="bg-emerald-50 dark:bg-emerald-950 rounded-lg p-3 text-center">
                    <p className="text-2xl font-bold text-emerald-600">{selectedItem.available}</p>
                    <p className="text-xs text-muted-foreground">Available</p>
                  </div>
                  <div className="bg-muted/50 rounded-lg p-3 text-center">
                    <p className="text-2xl font-bold text-blue-600">{selectedItem.onOrder}</p>
                    <p className="text-xs text-muted-foreground">On Order</p>
                  </div>
                </div>
              </div>
              
              <Separator />
              
              {/* Dates & Cost */}
              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-2">Dates & Cost</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-xs text-muted-foreground">Production Date</Label>
                    <p className="font-medium">{formatDate(selectedItem.productionDate)}</p>
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Expiry Date</Label>
                    <p className={`font-medium ${selectedItem.daysToExpiry <= 7 ? 'text-orange-600' : ''}`}>
                      {formatDate(selectedItem.expiryDate)}
                      {selectedItem.daysToExpiry <= 7 && selectedItem.daysToExpiry >= 0 && (
                        <span className="ml-2 text-xs">({selectedItem.daysToExpiry} days left)</span>
                      )}
                    </p>
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Average Cost</Label>
                    <p className="font-medium">{formatCurrency(selectedItem.averageCost)}</p>
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Inventory Value</Label>
                    <p className="font-medium text-emerald-600">{formatCurrency(selectedItem.inventoryValue)}</p>
                  </div>
                </div>
              </div>
            </div>
          )}
          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setDetailsDialogOpen(false)}>
              Close
            </Button>
            <Button onClick={() => { setDetailsDialogOpen(false); handleViewHistory(selectedItem!); }}>
              <History className="h-4 w-4 mr-2" /> View History
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Transaction History Dialog */}
      <Dialog open={historyDialogOpen} onOpenChange={setHistoryDialogOpen}>
        <DialogContent className="sm:max-w-[700px] max-h-[80vh]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <History className="h-5 w-5" />
              Transaction History
            </DialogTitle>
            <DialogDescription>
              {selectedItem?.productName} - Recent inventory movements
            </DialogDescription>
          </DialogHeader>
          <ScrollArea className="h-[400px] pr-4">
            {isLoadingHistory ? (
              <div className="flex items-center justify-center h-32">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
              </div>
            ) : Array.isArray(transactionHistory) && transactionHistory.length > 0 ? (
              <div className="space-y-3">
                {transactionHistory.map((tx: InventoryTransaction) => (
                  <div key={tx.id} className="flex items-start gap-4 p-3 rounded-lg border bg-card">
                    <div className="mt-1">
                      {getTransactionIcon(tx.transaction_type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between">
                        <span className="font-medium text-sm">{tx.transaction_type.replace('_', ' ')}</span>
                        <span className={`font-bold ${tx.quantity > 0 ? 'text-emerald-600' : 'text-red-600'}`}>
                          {tx.quantity > 0 ? '+' : ''}{tx.quantity}
                        </span>
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        {tx.reference_type && tx.reference_number && (
                          <span className="mr-3">Ref: {tx.reference_type} #{tx.reference_number}</span>
                        )}
                        {tx.lot_number && <span className="mr-3">Lot: {tx.lot_number}</span>}
                        {tx.unit_cost > 0 && <span>Cost: {formatCurrency(tx.unit_cost)}</span>}
                      </div>
                      {tx.notes && (
                        <p className="text-xs text-muted-foreground mt-1 italic">{tx.notes}</p>
                      )}
                      <p className="text-xs text-muted-foreground mt-2">
                        <CalendarDays className="h-3 w-3 inline mr-1" />
                        {new Date(tx.created_at).toLocaleString()}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-32 text-muted-foreground">
                <History className="h-8 w-8 mb-2" />
                <p>No transaction history found</p>
              </div>
            )}
          </ScrollArea>
          <DialogFooter>
            <Button variant="outline" onClick={() => setHistoryDialogOpen(false)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Adjust Stock Dialog */}
      <Dialog open={adjustDialogOpen} onOpenChange={setAdjustDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Adjust Stock</DialogTitle>
            <DialogDescription>
              {selectedItem && (
                <span>
                  Adjusting <strong>{selectedItem.productName}</strong> at {selectedItem.warehouseName}
                  <br />
                  Current quantity: <strong>{selectedItem.onHand}</strong>
                </span>
              )}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="quantity" className="text-right">
                Quantity
              </Label>
              <Input
                id="quantity"
                type="number"
                placeholder="Use negative for decrease"
                value={adjustmentQty}
                onChange={(e) => setAdjustmentQty(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="reason" className="text-right">
                Reason
              </Label>
              <Select value={adjustmentReason} onValueChange={setAdjustmentReason}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select reason" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="CYCLE_COUNT">Cycle Count</SelectItem>
                  <SelectItem value="DAMAGE">Damage</SelectItem>
                  <SelectItem value="EXPIRED">Expired</SelectItem>
                  <SelectItem value="THEFT">Theft/Loss</SelectItem>
                  <SelectItem value="FOUND">Found</SelectItem>
                  <SelectItem value="OTHER">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="notes" className="text-right">
                Notes
              </Label>
              <Textarea
                id="notes"
                placeholder="Additional notes..."
                value={adjustmentNotes}
                onChange={(e) => setAdjustmentNotes(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setAdjustDialogOpen(false)}>
              Cancel
            </Button>
            <Button 
              onClick={submitAdjustment} 
              disabled={!adjustmentQty || !adjustmentReason || adjustMutation.isPending}
            >
              {adjustMutation.isPending ? "Adjusting..." : "Confirm Adjustment"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Transfer Stock Dialog */}
      <Dialog open={transferDialogOpen} onOpenChange={setTransferDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Transfer Stock</DialogTitle>
            <DialogDescription>
              {selectedItem && (
                <span>
                  Transfer <strong>{selectedItem.productName}</strong> from {selectedItem.warehouseName}
                  <br />
                  Available quantity: <strong>{selectedItem.available}</strong>
                </span>
              )}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="transfer-qty" className="text-right">
                Quantity
              </Label>
              <Input
                id="transfer-qty"
                type="number"
                min="1"
                max={selectedItem?.available || 0}
                placeholder="Quantity to transfer"
                value={transferQty}
                onChange={(e) => setTransferQty(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="to-warehouse" className="text-right">
                To Warehouse
              </Label>
              <Select value={transferToWarehouse} onValueChange={setTransferToWarehouse}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select destination" />
                </SelectTrigger>
                <SelectContent>
                  {warehouses
                    .filter((w: any) => w.id !== selectedItem?.warehouseId)
                    .map((warehouse: any) => (
                      <SelectItem key={warehouse.id} value={warehouse.id.toString()}>
                        {warehouse.name}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="transfer-notes" className="text-right">
                Notes
              </Label>
              <Textarea
                id="transfer-notes"
                placeholder="Transfer notes..."
                value={transferNotes}
                onChange={(e) => setTransferNotes(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setTransferDialogOpen(false)}>
              Cancel
            </Button>
            <Button 
              onClick={submitTransfer} 
              disabled={!transferQty || !transferToWarehouse || transferMutation.isPending}
            >
              {transferMutation.isPending ? "Transferring..." : "Confirm Transfer"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Receive Stock Dialog */}
      <Dialog open={receiveDialogOpen} onOpenChange={setReceiveDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <PackagePlus className="h-5 w-5" />
              Receive Stock
            </DialogTitle>
            <DialogDescription>
              Add new inventory to a warehouse
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Product</Label>
              <Select value={receiveProductId} onValueChange={setReceiveProductId}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select product" />
                </SelectTrigger>
                <SelectContent>
                  {products.map((p: any) => (
                    <SelectItem key={p.id} value={p.id.toString()}>
                      {p.sku} - {p.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Warehouse</Label>
              <Select value={receiveWarehouseId} onValueChange={setReceiveWarehouseId}>
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select warehouse" />
                </SelectTrigger>
                <SelectContent>
                  {warehouses.map((w: any) => (
                    <SelectItem key={w.id} value={w.id.toString()}>
                      {w.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Quantity</Label>
              <Input
                type="number"
                min="1"
                placeholder="Quantity"
                value={receiveQty}
                onChange={(e) => setReceiveQty(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Unit Cost</Label>
              <Input
                type="number"
                min="0"
                step="0.01"
                placeholder="Cost per unit"
                value={receiveCost}
                onChange={(e) => setReceiveCost(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Lot Number</Label>
              <Input
                placeholder="Optional lot number"
                value={receiveLotNumber}
                onChange={(e) => setReceiveLotNumber(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Expiry Date</Label>
              <Input
                type="date"
                value={receiveExpiryDate}
                onChange={(e) => setReceiveExpiryDate(e.target.value)}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Notes</Label>
              <Textarea
                placeholder="Optional notes..."
                value={receiveNotes}
                onChange={(e) => setReceiveNotes(e.target.value)}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setReceiveDialogOpen(false)}>
              Cancel
            </Button>
            <Button 
              onClick={submitReceive} 
              disabled={!receiveProductId || !receiveWarehouseId || !receiveQty || receiveMutation.isPending}
            >
              {receiveMutation.isPending ? "Receiving..." : "Receive Stock"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
