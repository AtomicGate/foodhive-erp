import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { 
  DollarSign, 
  Search, 
  Filter, 
  Download, 
  Upload,
  History,
  Loader2
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { pricingService } from "@/services/pricingService";
import { toast } from "sonner";

export default function PricingManagement() {
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState("");

  const { data: productPrices, isLoading: isProductsLoading } = useQuery({
    queryKey: ['priceList'],
    queryFn: () => pricingService.getPriceList()
  });

  const updatePriceMutation = useMutation({
    mutationFn: ({ productId, price }: { productId: number; price: number }) => 
      pricingService.setProductPrice({ product_id: productId, price }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['priceList'] });
      toast.success("Price updated successfully");
    },
    onError: () => {
      toast.error("Failed to update price");
    }
  });

  const filteredProducts = productPrices?.filter((p: any) => 
    (p.product_name || '').toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  if (isProductsLoading) {
    return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Pricing Management</h1>
          <p className="text-muted-foreground">
            Manage product base prices, customer contracts, and promotions.
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <History className="mr-2 h-4 w-4" /> History
          </Button>
          <Button variant="outline">
            <Upload className="mr-2 h-4 w-4" /> Import
          </Button>
          <Button variant="outline">
            <Download className="mr-2 h-4 w-4" /> Export
          </Button>
        </div>
      </div>

      <Tabs defaultValue="products" className="space-y-4">
        <TabsList>
          <TabsTrigger value="products">Product Base Pricing</TabsTrigger>
          <TabsTrigger value="customers">Customer Specific</TabsTrigger>
          <TabsTrigger value="promotions">Promotions</TabsTrigger>
        </TabsList>

        <TabsContent value="products" className="space-y-4">
          <div className="flex items-center gap-2">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input 
                placeholder="Search products..." 
                className="pl-9" 
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
            <Button variant="outline" size="icon">
              <Filter className="h-4 w-4" />
            </Button>
          </div>

          <div className="rounded-md border bg-card">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product Name</TableHead>
                  <TableHead className="text-right">Min Price</TableHead>
                  <TableHead className="text-right">Base Price</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredProducts.map((product: any) => (
                  <TableRow key={product.id}>
                    <TableCell className="font-medium">{product.product_name || product.productName}</TableCell>
                    <TableCell className="text-right text-muted-foreground">
                      {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(product.min_quantity || 0)}
                    </TableCell>
                    <TableCell className="text-right font-bold">
                      {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(product.price || 0)}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button 
                        variant="ghost" 
                        size="sm"
                        onClick={() => {
                          const newPrice = prompt("Enter new base price:", (product.price || 0).toString());
                          if (newPrice && !isNaN(parseFloat(newPrice))) {
                            updatePriceMutation.mutate({ productId: product.product_id, price: parseFloat(newPrice) });
                          }
                        }}
                      >
                        Edit
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
                {filteredProducts.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={4} className="h-24 text-center">
                      No products found.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </TabsContent>

        <TabsContent value="customers" className="space-y-4">
          <div className="flex items-center justify-center h-48 text-muted-foreground border rounded-md bg-muted/10">
            Customer specific pricing module coming soon...
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
