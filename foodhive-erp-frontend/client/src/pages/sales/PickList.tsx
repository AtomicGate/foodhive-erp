import { useState, useEffect } from "react";
import { Link, useRoute } from "wouter";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { 
  CheckSquare, 
  Printer, 
  ArrowLeft, 
  Package, 
  MapPin,
  AlertCircle,
  Loader2
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Separator } from "@/components/ui/separator";
import { Progress } from "@/components/ui/progress";
import { toast } from "sonner";
import { salesService } from "@/services/salesService";

export default function PickList() {
  const [, params] = useRoute("/pick-list/:id");
  const id = params?.id;
  const queryClient = useQueryClient();
  
  // Local state for tracking picked items before submission
  const [pickedItems, setPickedItems] = useState<Record<string, boolean>>({});

  const { data: pickList, isLoading, error } = useQuery({
    queryKey: ['pickList', id],
    queryFn: () => salesService.getPickList(id!),
    enabled: !!id
  });

  // Initialize picked state when data loads
  useEffect(() => {
    if (pickList?.items) {
      const initialPicked: Record<string, boolean> = {};
      pickList.items.forEach((item: any) => {
        initialPicked[item.id] = item.picked || false;
      });
      setPickedItems(initialPicked);
    }
  }, [pickList]);

  const updatePickStatusMutation = useMutation({
    mutationFn: (data: { id: string, items: any[] }) => 
      salesService.updatePickList(data.id, { items: data.items }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pickList', id] });
      toast.success("Pick list updated successfully");
    },
    onError: () => {
      toast.error("Failed to update pick list");
    }
  });

  const completePickListMutation = useMutation({
    mutationFn: (id: string) => salesService.completePickList(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pickList', id] });
      toast.success("Pick list completed successfully!");
    }
  });

  if (isLoading) return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  if (error || !pickList) return <div className="p-8 text-center text-red-500">Failed to load pick list</div>;

  const items = pickList.items || [];
  const totalItems = items.length;
  const pickedCount = Object.values(pickedItems).filter(Boolean).length;
  const progress = totalItems > 0 ? (pickedCount / totalItems) * 100 : 0;

  const toggleItem = (itemId: string) => {
    const newPickedState = { ...pickedItems, [itemId]: !pickedItems[itemId] };
    setPickedItems(newPickedState);
    
    // Optimistically update backend (or you could save only on "Complete")
    // For this implementation, we'll save on complete or add a "Save" button
    // But to keep it responsive, let's just update local state for now
  };

  const handleComplete = () => {
    if (progress < 100) {
      toast.error("Please pick all items before completing.");
      return;
    }
    completePickListMutation.mutate(id!);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/sales-orders">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-3xl font-bold tracking-tight">Pick List #{pickList.pickListNumber || id}</h1>
              <Badge variant="outline" className="ml-2">{pickList.status}</Badge>
            </div>
            <p className="text-muted-foreground">
              Order: {pickList.orderNumber} • {pickList.customerName}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => window.print()}>
            <Printer className="mr-2 h-4 w-4" /> Print
          </Button>
          <Button onClick={handleComplete} disabled={progress < 100 || completePickListMutation.isPending}>
            {completePickListMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <CheckSquare className="mr-2 h-4 w-4" />
            )}
            Complete Picking
          </Button>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>Items to Pick</CardTitle>
            <CardDescription>
              {totalItems - pickedCount} items remaining
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              {items.map((item: any) => (
                <div 
                  key={item.id} 
                  className={`flex items-start space-x-4 p-4 rounded-lg border transition-colors ${
                    pickedItems[item.id] ? "bg-muted/50 border-transparent" : "bg-card border-border"
                  }`}
                >
                  <Checkbox 
                    checked={pickedItems[item.id] || false} 
                    onCheckedChange={() => toggleItem(item.id)}
                    className="mt-1"
                  />
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center justify-between">
                      <p className={`font-medium ${pickedItems[item.id] ? "text-muted-foreground line-through" : ""}`}>
                        {item.productName}
                      </p>
                      <Badge variant={pickedItems[item.id] ? "secondary" : "default"}>
                        {item.quantity} {item.unit}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Package className="h-3 w-3" /> {item.sku}
                      </span>
                      <span className="flex items-center gap-1 text-primary font-medium">
                        <MapPin className="h-3 w-3" /> {item.location || 'Unassigned'}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Progress</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex items-center justify-between text-sm">
                  <span>Completion</span>
                  <span className="font-medium">{Math.round(progress)}%</span>
                </div>
                <Progress value={progress} className="h-2" />
              </div>
              <div className="grid grid-cols-2 gap-4 pt-4">
                <div className="space-y-1">
                  <span className="text-xs text-muted-foreground">Total Items</span>
                  <p className="text-2xl font-bold">{totalItems}</p>
                </div>
                <div className="space-y-1">
                  <span className="text-xs text-muted-foreground">Picked</span>
                  <p className="text-2xl font-bold text-primary">
                    {pickedCount}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Location Info</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-start gap-3">
                <MapPin className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="font-medium">{pickList.warehouseName || 'Main Warehouse'}</p>
                  <p className="text-sm text-muted-foreground">{pickList.zone || 'General Zone'}</p>
                </div>
              </div>
              <Separator />
              <div className="flex items-start gap-3">
                <AlertCircle className="h-5 w-5 text-amber-500 mt-0.5" />
                <div>
                  <p className="font-medium">Special Instructions</p>
                  <p className="text-sm text-muted-foreground">
                    {pickList.notes || 'No special instructions.'}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
