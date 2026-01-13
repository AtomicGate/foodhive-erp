import api from '@/lib/api';

export const salesService = {
  // Sales Orders - Backend: /sales-orders
  getOrders: async (params?: any) => {
    const response = await api.get('/sales-orders/list', { params });
    return response.data?.data || response.data || [];
  },
  getOrder: async (id: string) => {
    const response = await api.get(`/sales-orders/get/${id}`);
    return response.data?.data || response.data;
  },
  getOrderByNumber: async (orderNumber: string) => {
    const response = await api.get(`/sales-orders/number/${orderNumber}`);
    return response.data?.data || response.data;
  },
  createOrder: async (data: any) => {
    const response = await api.post('/sales-orders/create', data);
    return response.data;
  },
  updateOrder: async (id: string, data: any) => {
    const response = await api.put(`/sales-orders/update/${id}`, data);
    return response.data;
  },
  deleteOrder: async (id: string) => {
    const response = await api.delete(`/sales-orders/delete/${id}`);
    return response.data;
  },
  confirmOrder: async (id: string) => {
    const response = await api.post(`/sales-orders/confirm/${id}`);
    return response.data;
  },
  cancelOrder: async (id: string) => {
    const response = await api.post(`/sales-orders/cancel/${id}`);
    return response.data;
  },
  shipOrder: async (id: string) => {
    const response = await api.post(`/sales-orders/ship/${id}`);
    return response.data;
  },

  // Order Lines
  addLine: async (orderId: string, data: any) => {
    const response = await api.post(`/sales-orders/${orderId}/lines`, data);
    return response.data;
  },
  updateLine: async (lineId: string, data: any) => {
    const response = await api.put(`/sales-orders/lines/${lineId}`, data);
    return response.data;
  },
  deleteLine: async (lineId: string) => {
    const response = await api.delete(`/sales-orders/lines/${lineId}`);
    return response.data;
  },

  // Order Guide
  getOrderGuide: async (customerId: string) => {
    const response = await api.get(`/sales-orders/order-guide/${customerId}`);
    return response.data?.data || response.data || [];
  },

  // Lost Sales
  recordLostSale: async (data: any) => {
    const response = await api.post('/sales-orders/lost-sale', data);
    return response.data;
  },
  getLostSales: async (params?: any) => {
    const response = await api.get('/sales-orders/lost-sales', { params });
    return response.data?.data || response.data || [];
  },
  
  // Pick List Methods - Backend: /picking
  getPickList: async (id: string) => {
    const response = await api.get(`/picking/get/${id}`);
    return response.data?.data || response.data;
  },
  getPickLists: async (params?: any) => {
    const response = await api.get('/picking/list', { params });
    return response.data?.data || response.data || [];
  },
  createPickList: async (data: any) => {
    const response = await api.post('/picking/create', data);
    return response.data;
  },
  generatePickList: async (data: any) => {
    const response = await api.post('/picking/generate', data);
    return response.data;
  },
  startPicking: async (id: string) => {
    const response = await api.post(`/picking/${id}/start`);
    return response.data;
  },
  completePickList: async (id: string) => {
    const response = await api.post(`/picking/${id}/complete`);
    return response.data;
  },
  updatePickList: async (id: string, data: any) => {
    const response = await api.put(`/picking/update/${id}`, data);
    return response.data;
  },
  cancelPickList: async (id: string) => {
    const response = await api.post(`/picking/${id}/cancel`);
    return response.data;
  },
  getPickLines: async (id: string) => {
    const response = await api.get(`/picking/${id}/lines`);
    return response.data?.data || response.data || [];
  },
  confirmPickLine: async (lineId: string, data: any) => {
    const response = await api.post(`/picking/lines/${lineId}/confirm`, data);
    return response.data;
  },

  // Invoice Methods - Backend: /ar/invoices
  getInvoice: async (id: string) => {
    const response = await api.get(`/ar/invoices/get/${id}`);
    return response.data?.data || response.data;
  },
  getInvoices: async (params?: any) => {
    const response = await api.get('/ar/invoices/list', { params });
    return response.data?.data || response.data || [];
  },
  getInvoiceByNumber: async (number: string) => {
    const response = await api.get(`/ar/invoices/number/${number}`);
    return response.data?.data || response.data;
  },
  createInvoice: async (data: any) => {
    const response = await api.post('/ar/invoices/create', data);
    return response.data;
  },
  createInvoiceFromOrder: async (orderId: string) => {
    const response = await api.post(`/ar/invoices/from-order/${orderId}`);
    return response.data;
  },
  postInvoice: async (id: string) => {
    const response = await api.post(`/ar/invoices/${id}/post`);
    return response.data;
  },
  voidInvoice: async (id: string) => {
    const response = await api.post(`/ar/invoices/${id}/void`);
    return response.data;
  }
};
