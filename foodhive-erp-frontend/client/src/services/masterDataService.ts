import api from "@/lib/api";

export interface Department {
  id: number;
  name: string;
  code: string;
  manager_id?: number;
  manager_name?: string;
  description?: string;
  is_active: boolean;
}

export interface Role {
  id: number;
  name?: string;  // Frontend name
  role_name?: string;  // Backend field name
  description: string;
  is_active?: boolean;
  permissions?: Permission[];
}

export interface Permission {
  page_id: number;
  page_name: string;
  can_create: boolean;
  can_view: boolean;
  can_update: boolean;
  can_delete: boolean;
}

export interface Page {
  id: number;
  name: string;
  code: string;
  module: string;
  description?: string;
}

export interface Warehouse {
  id: number;
  name: string;
  code: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  is_active: boolean;
}

export interface Employee {
  id: number;
  email: string;
  english_name: string;
  arabic_name?: string;
  phone?: string;
  nationality?: string;
  role_id?: number;
  status: string;
  created_at?: string;
  updated_at?: string;
}

export interface Customer {
  id: number;
  customer_code: string;
  name: string;
  email?: string;
  phone?: string;
  credit_limit: number;
  payment_terms_days: number;
  sales_rep_id?: number;
  is_active: boolean;
}

export interface Vendor {
  id: number;
  vendor_code: string;
  name: string;
  email?: string;
  phone?: string;
  payment_terms_days: number;
  is_active: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

// Helper to wrap array responses - always returns a valid PaginatedResponse
// Handles double-wrapped responses: {data: {data: [...], pagination: {...}}}
const wrapResponse = <T>(data: any, page: number, limit: number): PaginatedResponse<T> => {
  // Handle double-wrapped response from SuccessResponse helper
  const innerData = data?.data?.data || data?.data || data || [];
  const pagination = data?.data?.pagination || data?.pagination;
  
  // Always return a valid PaginatedResponse, even if items is not an array
  const safeItems = Array.isArray(innerData) ? innerData : [];
  return {
    data: safeItems,
    total: pagination?.total_items || safeItems.length,
    page,
    limit
  };
};

export const masterDataService = {
  // ===== Departments - Backend: /departments =====
  getDepartments: async (page = 1, limit = 10): Promise<PaginatedResponse<Department>> => {
    const response = await api.get(`/departments/list`, { params: { page, per_page: limit } });
    return wrapResponse<Department>(response.data, page, limit);
  },
  getDepartment: async (id: string): Promise<Department> => {
    const response = await api.get(`/departments/get/${id}`);
    return response.data?.data || response.data;
  },
  createDepartment: async (data: Partial<Department>) => {
    const response = await api.post("/departments/create", data);
    return response.data;
  },
  updateDepartment: async (id: string, data: Partial<Department>) => {
    const response = await api.put(`/departments/update/${id}`, data);
    return response.data;
  },
  deleteDepartment: async (id: string) => {
    const response = await api.delete(`/departments/delete/${id}`);
    return response.data;
  },

  // ===== Roles - Backend: /roles =====
  getRoles: async (page = 1, limit = 10): Promise<PaginatedResponse<Role>> => {
    const response = await api.get(`/roles/list`, { params: { page, per_page: limit } });
    // Backend wraps in {data: [...]} via SuccessResponse
    const roles = Array.isArray(response.data) ? response.data : (response.data?.data || []);
    return {
      data: roles,
      total: roles.length,
      page,
      limit
    };
  },
  getRole: async (id: string): Promise<Role> => {
    const response = await api.get(`/roles/get/${id}`);
    return response.data?.data || response.data;
  },
  createRole: async (data: Partial<Role>) => {
    const response = await api.post("/roles/create", data);
    return response.data;
  },
  updateRole: async (id: string, data: Partial<Role>) => {
    const response = await api.put(`/roles/update/${id}`, data);
    return response.data;
  },
  deleteRole: async (id: string) => {
    const response = await api.delete(`/roles/delete/${id}`);
    return response.data;
  },
  
  // Role Permissions
  setRolePermissions: async (roleId: string, permissions: Permission[]) => {
    const response = await api.post(`/roles/${roleId}/permissions`, { permissions });
    return response.data;
  },

  // Pages (Modules)
  getPages: async (): Promise<Page[]> => {
    const response = await api.get("/roles/pages/list");
    return response.data?.data || response.data || [];
  },

  // ===== Warehouses - Backend: /warehouses =====
  getWarehouses: async (page = 1, limit = 10): Promise<PaginatedResponse<Warehouse>> => {
    const response = await api.get(`/warehouses/list`, { params: { page, per_page: limit } });
    return wrapResponse<Warehouse>(response.data, page, limit);
  },
  getWarehouse: async (id: string): Promise<Warehouse> => {
    const response = await api.get(`/warehouses/get/${id}`);
    return response.data?.data || response.data;
  },
  createWarehouse: async (data: Partial<Warehouse>) => {
    const response = await api.post("/warehouses/create", data);
    return response.data;
  },
  updateWarehouse: async (id: string, data: Partial<Warehouse>) => {
    const response = await api.put(`/warehouses/update/${id}`, data);
    return response.data;
  },
  deleteWarehouse: async (id: string) => {
    const response = await api.delete(`/warehouses/delete/${id}`);
    return response.data;
  },

  // ===== Employees - Backend: /employees =====
  getEmployees: async (page = 1, limit = 10): Promise<PaginatedResponse<Employee>> => {
    const response = await api.get(`/employees/list`, { params: { page, page_size: limit } });
    
    // Backend returns: { data: { data: [...], pagination: {...} } }
    const innerData = response.data?.data?.data || response.data?.data || response.data || [];
    const employees = Array.isArray(innerData) 
      ? innerData.map((item: any) => item.employee || item)
      : [];
    
    const pagination = response.data?.data?.pagination;
    
    return {
      data: employees,
      total: pagination?.total_items || employees.length,
      page: pagination?.page || page,
      limit: pagination?.page_size || limit
    };
  },
  getEmployee: async (id: string): Promise<Employee> => {
    const response = await api.get(`/employees/get/${id}`);
    const data = response.data?.data || response.data;
    // Unwrap nested employee if present
    return data?.employee || data;
  },
  createEmployee: async (data: Partial<Employee> & { password?: string }) => {
    const response = await api.post("/employees/create", data);
    return response.data;
  },
  updateEmployee: async (id: string, data: Partial<Employee>) => {
    const response = await api.put(`/employees/update/${id}`, data);
    return response.data;
  },
  deleteEmployee: async (id: string) => {
    const response = await api.delete(`/employees/delete/${id}`);
    return response.data;
  },

  // ===== Customers - Backend: /customers =====
  getCustomers: async (page = 1, limit = 10): Promise<PaginatedResponse<Customer>> => {
    const response = await api.get(`/customers/list`, { params: { page, per_page: limit } });
    console.log('Raw API response:', response.data);
    const result = wrapResponse<Customer>(response.data, page, limit);
    console.log('Wrapped response:', result);
    return result;
  },
  getCustomer: async (id: string): Promise<Customer> => {
    const response = await api.get(`/customers/get/${id}`);
    return response.data?.data || response.data;
  },
  createCustomer: async (data: Partial<Customer>) => {
    const response = await api.post("/customers/create", data);
    return response.data;
  },
  updateCustomer: async (id: string, data: Partial<Customer>) => {
    const response = await api.put(`/customers/update/${id}`, data);
    return response.data;
  },
  deleteCustomer: async (id: string) => {
    const response = await api.delete(`/customers/delete/${id}`);
    return response.data;
  },

  // ===== Vendors - Backend: /vendors =====
  getVendors: async (page = 1, limit = 10): Promise<PaginatedResponse<Vendor>> => {
    const response = await api.get(`/vendors/list`, { params: { page, per_page: limit } });
    return wrapResponse<Vendor>(response.data, page, limit);
  },
  getVendor: async (id: string): Promise<Vendor> => {
    const response = await api.get(`/vendors/get/${id}`);
    return response.data?.data || response.data;
  },
  createVendor: async (data: Partial<Vendor>) => {
    const response = await api.post("/vendors/create", data);
    return response.data;
  },
  updateVendor: async (id: string, data: Partial<Vendor>) => {
    const response = await api.put(`/vendors/update/${id}`, data);
    return response.data;
  },
  deleteVendor: async (id: string) => {
    const response = await api.delete(`/vendors/delete/${id}`);
    return response.data;
  },
};
