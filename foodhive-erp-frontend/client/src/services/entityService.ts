import api from '@/lib/api';

// Entity type to backend endpoint mapping
const entityEndpoints: Record<string, string> = {
  customers: '/customers',
  vendors: '/vendors',
  employees: '/employees',
  departments: '/departments',
  warehouses: '/warehouses',
  roles: '/roles',
  products: '/products',
};

// Generic entity service that maps to specific backend endpoints
export const entityService = {
  getEntities: async (type: string, params?: any) => {
    const endpoint = entityEndpoints[type] || `/${type}`;
    const response = await api.get(`${endpoint}/list`, { params });
    return response.data?.data || response.data || [];
  },

  getEntity: async (type: string, id: string) => {
    const endpoint = entityEndpoints[type] || `/${type}`;
    const response = await api.get(`${endpoint}/get/${id}`);
    return response.data?.data || response.data;
  },

  createEntity: async (type: string, data: any) => {
    const endpoint = entityEndpoints[type] || `/${type}`;
    const response = await api.post(`${endpoint}/create`, data);
    return response.data;
  },

  updateEntity: async (type: string, id: string, data: any) => {
    const endpoint = entityEndpoints[type] || `/${type}`;
    const response = await api.put(`${endpoint}/update/${id}`, data);
    return response.data;
  },

  deleteEntity: async (type: string, id: string) => {
    const endpoint = entityEndpoints[type] || `/${type}`;
    const response = await api.delete(`${endpoint}/delete/${id}`);
    return response.data;
  },

  // Helper to check if entity type is supported
  isSupported: (type: string): boolean => {
    return type in entityEndpoints;
  },
};
