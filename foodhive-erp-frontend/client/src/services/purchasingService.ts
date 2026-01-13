import api from '@/lib/api';

export interface PurchaseOrder {
  id: string;
  poNumber: string;
  vendorId: string;
  vendorName: string;
  date: string;
  expectedDate: string;
  status: 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'PARTIAL' | 'RECEIVED' | 'CANCELLED';
  totalAmount: number;
  items: PurchaseOrderItem[];
}

export interface PurchaseOrderItem {
  id: string;
  productId: string;
  productName: string;
  quantity: number;
  receivedQuantity: number;
  unitPrice: number;
  total: number;
}

export interface Receiving {
  id: string;
  poId: string;
  receivedDate: string;
  receivedBy: number;
  lines: ReceivingLine[];
}

export interface ReceivingLine {
  poLineId: string;
  quantityReceived: number;
  lotNumber?: string;
  expirationDate?: string;
}

export const purchasingService = {
  // Purchase Orders - Backend: /purchase-orders
  getPurchaseOrders: async (params?: any) => {
    const response = await api.get('/purchase-orders/list', { params });
    return response.data?.data || response.data || [];
  },

  getPurchaseOrder: async (id: string) => {
    const response = await api.get(`/purchase-orders/get/${id}`);
    return response.data?.data || response.data;
  },

  getPurchaseOrderByNumber: async (poNumber: string) => {
    const response = await api.get(`/purchase-orders/number/${poNumber}`);
    return response.data?.data || response.data;
  },

  createPurchaseOrder: async (data: Partial<PurchaseOrder>) => {
    const response = await api.post('/purchase-orders/create', data);
    return response.data;
  },

  updatePurchaseOrder: async (id: string, data: Partial<PurchaseOrder>) => {
    const response = await api.put(`/purchase-orders/update/${id}`, data);
    return response.data;
  },

  deletePurchaseOrder: async (id: string) => {
    const response = await api.delete(`/purchase-orders/delete/${id}`);
    return response.data;
  },

  submitPurchaseOrder: async (id: string) => {
    const response = await api.post(`/purchase-orders/submit/${id}`);
    return response.data;
  },

  cancelPurchaseOrder: async (id: string) => {
    const response = await api.post(`/purchase-orders/cancel/${id}`);
    return response.data;
  },

  // PO Lines
  addLine: async (poId: string, data: any) => {
    const response = await api.post(`/purchase-orders/${poId}/lines`, data);
    return response.data;
  },

  updateLine: async (lineId: string, data: any) => {
    const response = await api.put(`/purchase-orders/lines/${lineId}`, data);
    return response.data;
  },

  deleteLine: async (lineId: string) => {
    const response = await api.delete(`/purchase-orders/lines/${lineId}`);
    return response.data;
  },
  
  // Receiving
  createReceiving: async (data: Partial<Receiving>) => {
    const response = await api.post('/purchase-orders/receive', data);
    return response.data;
  },

  getReceiving: async (id: string) => {
    const response = await api.get(`/purchase-orders/receiving/${id}`);
    return response.data?.data || response.data;
  },

  getReceivings: async (params?: any) => {
    const response = await api.get('/purchase-orders/receivings', { params });
    return response.data?.data || response.data || [];
  }
};
