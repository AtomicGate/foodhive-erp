import { useState } from "react";
import { 
  Scale, 
  Search, 
  ArrowRight, 
  CheckCircle, 
  AlertTriangle,
  XCircle
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function CatchWeight() {
  const [scannedItem, setScannedItem] = useState<any>(null);
  const [weight, setWeight] = useState("");

  const handleScan = () => {
    // Simulate scan
    setScannedItem({
      id: 1,
      sku: "BEEF001",
      name: "Ground Beef",
      expectedWeight: 50,
      unit: "LB",
      tolerance: 0.1 // 10%
    });
  };

  const getVarianceStatus = () => {
    if (!scannedItem || !weight) return null;
    const actual = parseFloat(weight);
    const expected = scannedItem.expectedWeight;
    const variance = Math.abs((actual - expected) / expected);
    
    if (variance <= 0.05) return "success"; // Within 5%
    if (variance <= scannedItem.tolerance) return "warning"; // Within tolerance
    return "error"; // Out of tolerance
  };

  const status = getVarianceStatus();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Catch Weight Capture</h1>
          <p className="text-muted-foreground">
            Record actual weights for variable weight products.
          </p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Item Entry</CardTitle>
            <CardDescription>Scan barcode or enter SKU manually</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Input placeholder="Scan barcode or enter SKU..." />
              <Button onClick={handleScan}>
                <Search className="h-4 w-4" />
              </Button>
            </div>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Reference Document</label>
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="Select Order / PO / Pick List" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="so-1">SO-2024-001</SelectItem>
                  <SelectItem value="pl-1">PL-2024-001</SelectItem>
                  <SelectItem value="po-1">PO-2024-001</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {scannedItem && (
          <Card className={
            status === "success" ? "border-emerald-500 bg-emerald-50/50" :
            status === "warning" ? "border-amber-500 bg-amber-50/50" :
            status === "error" ? "border-red-500 bg-red-50/50" : ""
          }>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>{scannedItem.name}</CardTitle>
                <Badge variant="outline">{scannedItem.sku}</Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Expected Weight:</span>
                <span className="font-medium">{scannedItem.expectedWeight} {scannedItem.unit}</span>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Actual Weight</label>
                <div className="flex gap-2">
                  <div className="relative flex-1">
                    <Input 
                      type="number" 
                      className="pl-10 text-lg font-bold" 
                      value={weight}
                      onChange={(e) => setWeight(e.target.value)}
                      autoFocus
                    />
                    <Scale className="absolute left-3 top-3 h-5 w-5 text-muted-foreground" />
                  </div>
                  <div className="flex items-center px-3 border rounded-md bg-muted font-medium">
                    {scannedItem.unit}
                  </div>
                </div>
              </div>

              {weight && (
                <div className="flex items-center gap-2">
                  {status === "success" && (
                    <>
                      <CheckCircle className="h-5 w-5 text-emerald-600" />
                      <span className="text-emerald-700 font-medium">Weight Verified</span>
                    </>
                  )}
                  {status === "warning" && (
                    <>
                      <AlertTriangle className="h-5 w-5 text-amber-600" />
                      <span className="text-amber-700 font-medium">Variance Detected (Within Tolerance)</span>
                    </>
                  )}
                  {status === "error" && (
                    <>
                      <XCircle className="h-5 w-5 text-red-600" />
                      <span className="text-red-700 font-medium">Out of Tolerance! Approval Required.</span>
                    </>
                  )}
                </div>
              )}

              <Button className="w-full" disabled={!weight || status === "error"}>
                Confirm Weight <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
