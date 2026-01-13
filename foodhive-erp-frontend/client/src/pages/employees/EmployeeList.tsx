import { useState, useMemo, useCallback } from "react";
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
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown,
  Users,
  Loader2,
  Edit,
  Trash2,
  Mail,
  Phone,
  Shield
} from "lucide-react";
import { toast } from "sonner";
import { masterDataService, Employee } from "@/services/masterDataService";

export default function EmployeeList() {
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null);
  const [formData, setFormData] = useState({
    english_name: "",
    arabic_name: "",
    email: "",
    phone: "",
    nationality: "",
    role_id: "",
    password: "",
    status: "CONTINUED"  // Database enum: CONTINUED, RESIGN, ON_LEAVE, SUSPENDED
  });

  const { data: employeesData, isLoading, error, refetch: refetchEmployees } = useQuery({
    queryKey: ['employees'],
    queryFn: () => masterDataService.getEmployees(1, 100),
    staleTime: 30000, // Cache for 30 seconds
    refetchOnWindowFocus: false, // Don't refetch on window focus
    placeholderData: (previousData) => previousData, // Keep previous data during refetch
  });

  const { data: departmentsData } = useQuery({
    queryKey: ['departments'],
    queryFn: () => masterDataService.getDepartments(1, 100)
  });

  const { data: rolesData } = useQuery({
    queryKey: ['roles'],
    queryFn: () => masterDataService.getRoles(1, 100)
  });

  // Ensure arrays are always arrays to prevent .filter errors - memoized to prevent re-renders
  const employees = useMemo(() => 
    Array.isArray(employeesData?.data) ? employeesData.data : [], 
    [employeesData?.data]
  );
  const departments = useMemo(() => 
    Array.isArray(departmentsData?.data) ? departmentsData.data : [], 
    [departmentsData?.data]
  );
  const roles = useMemo(() => 
    Array.isArray(rolesData?.data) ? rolesData.data : [], 
    [rolesData?.data]
  );

  const createMutation = useMutation({
    mutationFn: (data: any) => masterDataService.createEmployee(data),
    onSuccess: () => {
      refetchEmployees();
      setIsDialogOpen(false);
      resetForm();
      toast.success("Employee created successfully");
    },
    onError: (error: any) => {
      const message = error?.response?.data?.error || error?.response?.data?.message || "Failed to create employee";
      toast.error(message);
    }
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => masterDataService.updateEmployee(id, data),
    onSuccess: () => {
      // Close dialog first
      setIsDialogOpen(false);
      setEditingEmployee(null);
      resetForm();
      toast.success("Employee updated successfully");
      // Refetch in background
      refetchEmployees();
    },
    onError: (error: any) => {
      const message = error?.response?.data?.error || error?.response?.data?.message || "Failed to update employee";
      toast.error(message);
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => masterDataService.deleteEmployee(id),
    onSuccess: () => {
      refetchEmployees();
      toast.success("Employee deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete employee");
    }
  });

  const resetForm = () => {
    setFormData({
      english_name: "",
      arabic_name: "",
      email: "",
      phone: "",
      nationality: "",
      role_id: "",
      password: "",
      status: "CONTINUED"
    });
  };

  const handleEdit = useCallback((employee: Employee) => {
    setEditingEmployee(employee);
    setFormData({
      english_name: employee.english_name || "",
      arabic_name: employee.arabic_name || "",
      email: employee.email,
      phone: employee.phone || "",
      nationality: employee.nationality || "",
      role_id: employee.role_id?.toString() || "",
      password: "",
      status: employee.status || "CONTINUED"
    });
    setIsDialogOpen(true);
  }, []);

  const handleSubmit = () => {
    // Validation
    if (!formData.english_name.trim()) {
      toast.error("English name is required");
      return;
    }
    if (!formData.email.trim()) {
      toast.error("Email is required");
      return;
    }
    if (!formData.role_id) {
      toast.error("Role is required");
      return;
    }
    if (!editingEmployee && !formData.password) {
      toast.error("Password is required for new employees");
      return;
    }
    if (!editingEmployee && formData.password.length < 8) {
      toast.error("Password must be at least 8 characters");
      return;
    }
    if (editingEmployee && formData.password && formData.password.length < 8) {
      toast.error("Password must be at least 8 characters");
      return;
    }

    if (editingEmployee) {
      // Update - only send fields that backend accepts
      const updateData: any = {
        english_name: formData.english_name,
        arabic_name: formData.arabic_name || null,
        nationality: formData.nationality || null,
        phone: formData.phone || null,
        role_id: parseInt(formData.role_id),
        status: formData.status,
      };
      console.log("Updating employee with:", updateData);
      updateMutation.mutate({ id: String(editingEmployee.id), data: updateData });
    } else {
      // Create - send all fields including email and password
      const createData = {
        ...formData,
        role_id: parseInt(formData.role_id)
      };
      createMutation.mutate(createData);
    }
  };

  const handleChangeStatus = useCallback((employeeId: number, newStatus: string) => {
    updateMutation.mutate({ 
      id: String(employeeId), 
      data: { status: newStatus } 
    });
  }, [updateMutation]);

  const getInitials = useCallback((name: string) => {
    const parts = (name || '').split(' ');
    return parts.map(p => p.charAt(0)).join('').toUpperCase().slice(0, 2) || 'U';
  }, []);

  const columns: ColumnDef<Employee>[] = useMemo(() => [
    {
      accessorKey: "id",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          ID
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => (
        <div className="font-mono font-medium">{row.getValue("id")}</div>
      ),
    },
    {
      id: "name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const employee = row.original;
        return (
          <div className="flex items-center gap-3">
            <Avatar className="h-8 w-8">
              <AvatarFallback className="bg-primary/10 text-primary text-xs">
                {getInitials(employee.english_name)}
              </AvatarFallback>
            </Avatar>
            <div>
              <div className="font-medium">{employee.english_name}</div>
              {employee.arabic_name && (
                <div className="text-xs text-muted-foreground">{employee.arabic_name}</div>
              )}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "email",
      header: "Email",
      cell: ({ row }) => (
        <div className="flex items-center gap-2 text-muted-foreground">
          <Mail className="h-4 w-4" />
          {row.getValue("email")}
        </div>
      ),
    },
    {
      accessorKey: "phone",
      header: "Phone",
      cell: ({ row }) => {
        const phone = row.getValue("phone") as string;
        return phone ? (
          <div className="flex items-center gap-2 text-muted-foreground">
            <Phone className="h-4 w-4" />
            {phone}
          </div>
        ) : "-";
      },
    },
    {
      accessorKey: "role_id",
      header: "Role",
      cell: ({ row }) => {
        const roleId = row.getValue("role_id") as number;
        const role = roles.find((r: any) => r.id === roleId);
        return role ? (
          <Badge variant="outline" className="flex items-center gap-1 w-fit">
            <Shield className="h-3 w-3" />
            {role.name || role.role_name || 'Unknown'}
          </Badge>
        ) : "-";
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        const statusColors: Record<string, string> = {
          'CONTINUED': 'default',
          'RESIGN': 'secondary',
          'ON_LEAVE': 'outline',
          'SUSPENDED': 'destructive'
        };
        const statusLabels: Record<string, string> = {
          'CONTINUED': 'Active',
          'RESIGN': 'Resigned',
          'ON_LEAVE': 'On Leave',
          'SUSPENDED': 'Suspended'
        };
        return (
          <Badge variant={statusColors[status] as any || "default"}>
            {statusLabels[status] || status || 'Active'}
          </Badge>
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const employee = row.original;
        const isActive = employee.status === 'CONTINUED';
        const isSuspended = employee.status === 'SUSPENDED';
        
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
              <DropdownMenuItem onClick={() => handleEdit(employee)}>
                <Edit className="mr-2 h-4 w-4" /> Edit Employee
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuLabel className="text-xs text-muted-foreground">Change Status</DropdownMenuLabel>
              {!isActive && (
                <DropdownMenuItem 
                  className="text-green-600"
                  onClick={() => handleChangeStatus(employee.id, 'CONTINUED')}
                >
                  <Shield className="mr-2 h-4 w-4" /> Activate
                </DropdownMenuItem>
              )}
              {isActive && (
                <DropdownMenuItem 
                  className="text-orange-600"
                  onClick={() => handleChangeStatus(employee.id, 'SUSPENDED')}
                >
                  <Shield className="mr-2 h-4 w-4" /> Suspend
                </DropdownMenuItem>
              )}
              {!isSuspended && employee.status !== 'ON_LEAVE' && (
                <DropdownMenuItem 
                  onClick={() => handleChangeStatus(employee.id, 'ON_LEAVE')}
                >
                  <Users className="mr-2 h-4 w-4" /> Set On Leave
                </DropdownMenuItem>
              )}
              <DropdownMenuSeparator />
              <DropdownMenuItem 
                className="text-destructive"
                onClick={() => {
                  if (confirm("Are you sure you want to mark this employee as resigned? This action cannot be undone.")) {
                    handleChangeStatus(employee.id, 'RESIGN');
                  }
                }}
              >
                <Trash2 className="mr-2 h-4 w-4" /> Mark as Resigned
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ], [roles, handleEdit, handleChangeStatus, getInitials]);

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
        Failed to load employees
      </div>
    );
  }

  const activeCount = employees.filter((e: Employee) => e.status === 'CONTINUED').length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Employees</h1>
          <p className="text-muted-foreground">
            Manage employee accounts and access
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={(open) => {
          setIsDialogOpen(open);
          if (!open) {
            setEditingEmployee(null);
            resetForm();
          }
        }}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" /> Add Employee
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>
                {editingEmployee ? "Edit Employee" : "Add New Employee"}
              </DialogTitle>
              <DialogDescription>
                {editingEmployee 
                  ? "Update employee information."
                  : "Enter the details for the new employee."}
              </DialogDescription>
            </DialogHeader>
            <form autoComplete="off" onSubmit={(e) => e.preventDefault()}>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="english_name">English Name *</Label>
                  <Input
                    id="english_name"
                    name="emp-english-name"
                    autoComplete="off"
                    placeholder="John Doe"
                    value={formData.english_name}
                    onChange={(e) => setFormData({...formData, english_name: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="arabic_name">Arabic Name</Label>
                  <Input
                    id="arabic_name"
                    name="emp-arabic-name"
                    autoComplete="off"
                    placeholder="الاسم بالعربية"
                    value={formData.arabic_name}
                    onChange={(e) => setFormData({...formData, arabic_name: e.target.value})}
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email *</Label>
                  <Input
                    id="email"
                    name="employee-email-field"
                    type="email"
                    autoComplete="off"
                    placeholder="john.doe@company.com"
                    value={formData.email}
                    onChange={(e) => setFormData({...formData, email: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">Phone</Label>
                  <Input
                    id="phone"
                    name="emp-phone-field"
                    autoComplete="off"
                    placeholder="+1 (555) 000-0000"
                    value={formData.phone}
                    onChange={(e) => setFormData({...formData, phone: e.target.value})}
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="nationality">Nationality</Label>
                  <Input
                    id="nationality"
                    autoComplete="off"
                    placeholder="e.g., American"
                    value={formData.nationality}
                    onChange={(e) => setFormData({...formData, nationality: e.target.value})}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="role">Role *</Label>
                  <Select 
                    value={formData.role_id} 
                    onValueChange={(value) => setFormData({...formData, role_id: value})}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={roles.length === 0 ? "No roles available" : "Select role"} />
                    </SelectTrigger>
                    <SelectContent>
                      {roles.length === 0 ? (
                        <SelectItem value="none" disabled>
                          No roles found - run 004_insert_roles.sql
                        </SelectItem>
                      ) : (
                        roles.map((role: any) => (
                          <SelectItem key={role.id} value={String(role.id)}>
                            {role.name || role.role_name}
                          </SelectItem>
                        ))
                      )}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="status">Status</Label>
                  <Select 
                    value={formData.status} 
                    onValueChange={(value) => setFormData({...formData, status: value})}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="CONTINUED">Active (Continued)</SelectItem>
                      <SelectItem value="ON_LEAVE">On Leave</SelectItem>
                      <SelectItem value="SUSPENDED">Suspended</SelectItem>
                      <SelectItem value="RESIGN">Resigned</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">
                    Password {!editingEmployee && "*"}
                    {editingEmployee && <span className="text-muted-foreground text-xs">(leave blank to keep current)</span>}
                  </Label>
                  <Input
                    id="password"
                    name="new-password-field"
                    type="password"
                    autoComplete="new-password"
                    placeholder={editingEmployee ? "Leave blank to keep current" : "Min 8 characters"}
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                  />
                </div>
              </div>
            </div>
            </form>
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
                {editingEmployee ? "Update" : "Create"}
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
                <p className="text-sm font-medium text-muted-foreground">Total Employees</p>
                <p className="text-2xl font-bold">{employees.length}</p>
              </div>
              <Users className="h-8 w-8 text-muted-foreground" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active</p>
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
                <p className="text-sm font-medium text-muted-foreground">Inactive</p>
                <p className="text-2xl font-bold text-gray-500">{employees.length - activeCount}</p>
              </div>
              <Badge variant="secondary">Inactive</Badge>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Data Table */}
      <DataTable 
        columns={columns} 
        data={employees} 
        searchKey="email"
        searchPlaceholder="Search by email..."
      />
    </div>
  );
}
