import api from '@/lib/api';

export const inventoryService = {
  // Inventory Query - Backend: /inventory
  getInventory: async (params?: any) => {
    const response = await api.get('/inventory/list', { params });
    return response.data?.data || response.data || [];
  },

  getInventoryItem: async (id: string) => {
    const response = await api.get(`/inventory/get/${id}`);
    return response.data?.data || response.data;
  },

  getInventoryByProduct: async (productId: string) => {
    const response = await api.get(`/inventory/product/${productId}`);
    return response.data?.data || response.data || [];
  },

  getInventoryByWarehouse: async (warehouseId: string) => {
    const response = await api.get(`/inventory/warehouse/${warehouseId}`);
    return response.data?.data || response.data || [];
  },

  getInventoryByLot: async (lotNumber: string) => {
    const response = await api.get(`/inventory/lot/${lotNumber}`);
    return response.data?.data || response.data || [];
  },

  // Summary
  getProductSummary: async (productId: string) => {
    const response = await api.get(`/inventory/summary/product/${productId}`);
    return response.data?.data || response.data;
  },

  getExpiringInventory: async (params?: { days?: number }) => {
    const response = await api.get('/inventory/expiring', { params });
    return response.data?.data || response.data || [];
  },

  // Operations
  receiveInventory: async (data: any) => {
    const response = await api.post('/inventory/receive', data);
    return response.data;
  },

  adjustInventory: async (data: any) => {
    const response = await api.post('/inventory/adjust', data);
    return response.data;
  },

  transferInventory: async (data: any) => {
    const response = await api.post('/inventory/transfer', data);
    return response.data;
  },

  // Transaction History
  getTransactions: async (params?: any) => {
    const response = await api.get('/inventory/transactions', { params });
    return response.data?.data || response.data || [];
  },
};
