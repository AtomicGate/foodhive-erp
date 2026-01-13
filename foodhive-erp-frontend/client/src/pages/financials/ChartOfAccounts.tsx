import { useState } from "react";
import { Link } from "wouter";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
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
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { DataTable } from "@/components/ui/data-table";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown,
  BookOpen,
  Loader2,
  Edit,
  Trash2,
  Eye,
  ChevronRight
} from "lucide-react";
import { toast } from "sonner";
import { financialService } from "@/services/financialService";

interface GLAccount {
  id: number;
  account_code: string;
  account_name: string;
  account_type: string;
  parent_id?: number;
  description?: string;
  is_active: boolean;
  balance?: number;
}

export default function ChartOfAccounts() {
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<GLAccount | null>(null);
  const [formData, setFormData] = useState({
    account_code: "",
    account_name: "",
    account_type: "ASSET",
    parent_id: "",
    description: "",
    is_active: true
  });

  const { data: accounts, isLoading, error } = useQuery({
    queryKey: ['glAccounts'],
    queryFn: financialService.getAccounts
  });

  const createMutation = useMutation({
    mutationFn: (data: any) => financialService.createAccount(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['glAccounts'] });
      setIsDialogOpen(false);
      resetForm();
      toast.success("Account created successfully");
    },
    onError: () => {
      toast.error("Failed to create account");
    }
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => financialService.updateAccount(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['glAccounts'] });
      setIsDialogOpen(false);
      setEditingAccount(null);
      resetForm();
      toast.success("Account updated successfully");
    },
    onError: () => {
      toast.error("Failed to update account");
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => financialService.deleteAccount(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['glAccounts'] });
      toast.success("Account deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete account");
    }
  });

  const resetForm = () => {
    setFormData({
      account_code: "",
      account_name: "",
      account_type: "ASSET",
      parent_id: "",
      description: "",
      is_active: true
    });
  };

  const handleEdit = (account: GLAccount) => {
    setEditingAccount(account);
    setFormData({
      account_code: account.account_code,
      account_name: account.account_name,
      account_type: account.account_type,
      parent_id: account.parent_id?.toString() || "",
      description: account.description || "",
      is_active: account.is_active
    });
    setIsDialogOpen(true);
  };

  const handleSubmit = () => {
    const data = {
      ...formData,
      parent_id: formData.parent_id ? parseInt(formData.parent_id) : null
    };

    if (editingAccount) {
      updateMutation.mutate({ id: String(editingAccount.id), data });
    } else {
      createMutation.mutate(data);
    }
  };

  const getAccountTypeColor = (type: string) => {
    switch (type) {
      case "ASSET": return "bg-emerald-100 text-emerald-800";
      case "LIABILITY": return "bg-red-100 text-red-800";
      case "EQUITY": return "bg-blue-100 text-blue-800";
      case "REVENUE": return "bg-purple-100 text-purple-800";
      case "EXPENSE": return "bg-amber-100 text-amber-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const formatCurrency = (value: number | undefined) => {
    if (value === undefined) return "-";
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value);
  };

  const columns: ColumnDef<GLAccount>[] = [
    {
      accessorKey: "account_code",
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
        <div className="font-mono font-medium">{row.getValue("account_code")}</div>
      ),
    },
    {
      accessorKey: "account_name",
      header: "Account Name",
      cell: ({ row }) => (
        <div className="font-medium">{row.getValue("account_name")}</div>
      ),
    },
    {
      accessorKey: "account_type",
      header: "Type",
      cell: ({ row }) => {
        const type = row.getValue("account_type") as string;
        return (
          <Badge variant="outline" className={getAccountTypeColor(type)}>
            {type}
          </Badge>
        );
      },
    },
    {
      accessorKey: "balance",
      header: () => <div className="text-right">Balance</div>,
      cell: ({ row }) => {
        const balance = row.getValue("balance") as number | undefined;
        return (
          <div className="text-right font-medium">
            {formatCurrency(balance)}
          </div>
        );
      },
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
        const account = row.original;
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
              <Link href={`/gl/accounts/${account.id}/activity`}>
                <DropdownMenuItem>
                  <Eye className="mr-2 h-4 w-4" /> View Activity
                </DropdownMenuItem>
              </Link>
              <DropdownMenuItem onClick={() => handleEdit(account)}>
                <Edit className="mr-2 h-4 w-4" /> Edit
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem 
                className="text-destructive"
                onClick={() => deleteMutation.mutate(String(account.id))}
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
        Failed to load chart of accounts
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-1">
            <Link href="/gl" className="hover:text-foreground">General Ledger</Link>
            <ChevronRight className="h-4 w-4" />
            <span>Chart of Accounts</span>
          </div>
          <h1 className="text-3xl font-bold tracking-tight">Chart of Accounts</h1>
          <p className="text-muted-foreground">
            Manage your general ledger accounts
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={(open) => {
          setIsDialogOpen(open);
          if (!open) {
            setEditingAccount(null);
            resetForm();
          }
        }}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" /> Add Account
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {editingAccount ? "Edit Account" : "Create New Account"}
              </DialogTitle>
              <DialogDescription>
                {editingAccount 
                  ? "Update the account details below."
                  : "Add a new account to your chart of accounts."}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="account_code">Account Code</Label>
                  <Input
                    id="account_code"
                    placeholder="1000"
                    value={formData.account_code}
                    onChange={(e) => setFormData({...formData, account_code: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="account_type">Account Type</Label>
                  <Select 
                    value={formData.account_type} 
                    onValueChange={(value) => setFormData({...formData, account_type: value})}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="ASSET">Asset</SelectItem>
                      <SelectItem value="LIABILITY">Liability</SelectItem>
                      <SelectItem value="EQUITY">Equity</SelectItem>
                      <SelectItem value="REVENUE">Revenue</SelectItem>
                      <SelectItem value="EXPENSE">Expense</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="account_name">Account Name</Label>
                <Input
                  id="account_name"
                  placeholder="Cash in Bank"
                  value={formData.account_name}
                  onChange={(e) => setFormData({...formData, account_name: e.target.value})}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="parent_id">Parent Account (Optional)</Label>
                <Select 
                  value={formData.parent_id} 
                  onValueChange={(value) => setFormData({...formData, parent_id: value})}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select parent account" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">None</SelectItem>
                    {accounts?.filter((a: GLAccount) => !editingAccount || a.id !== editingAccount.id)
                      .map((account: GLAccount) => (
                        <SelectItem key={account.id} value={String(account.id)}>
                          {account.account_code} - {account.account_name}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  placeholder="Account description..."
                  value={formData.description}
                  onChange={(e) => setFormData({...formData, description: e.target.value})}
                />
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
                {editingAccount ? "Update" : "Create"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-5">
        {['ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE'].map((type) => {
          const count = accounts?.filter((a: GLAccount) => a.account_type === type).length || 0;
          return (
            <Card key={type}>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">{type}</p>
                    <p className="text-2xl font-bold">{count}</p>
                  </div>
                  <Badge variant="outline" className={getAccountTypeColor(type)}>
                    {type.charAt(0)}
                  </Badge>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Data Table */}
      <DataTable 
        columns={columns} 
        data={accounts || []} 
        searchKey="account_name"
        searchPlaceholder="Search accounts..."
      />
    </div>
  );
}
