import api from '@/lib/api';

export interface CatchWeightEntry {
  id: number;
  product_id: number;
  product_name?: string;
  reference_type: 'RECEIVING' | 'SALES_ORDER' | 'PICK_LIST' | 'ADJUSTMENT';
  reference_id: number;
  lot_number?: string;
  expected_weight: number;
  actual_weight: number;
  variance_weight: number;
  variance_percent: number;
  unit: string;
  captured_by: number;
  captured_at: string;
  is_billed: boolean;
}

export interface CatchWeightPiece {
  id: number;
  entry_id: number;
  piece_number: number;
  weight: number;
  barcode?: string;
  notes?: string;
}

export interface CatchWeightConfig {
  product_id: number;
  product_name?: string;
  is_catch_weight: boolean;
  catch_weight_unit: string;
  average_weight?: number;
  min_weight?: number;
  max_weight?: number;
  tolerance_percent: number;
  requires_piece_weights: boolean;
}

export interface CaptureRequest {
  product_id: number;
  reference_type: 'RECEIVING' | 'SALES_ORDER' | 'PICK_LIST' | 'ADJUSTMENT';
  reference_id: number;
  lot_number?: string;
  expected_weight: number;
  actual_weight: number;
  pieces?: Array<{
    piece_number: number;
    weight: number;
    barcode?: string;
    notes?: string;
  }>;
}

export interface VarianceReport {
  product_id: number;
  product_name: string;
  total_expected: number;
  total_actual: number;
  total_variance: number;
  avg_variance_percent: number;
  entry_count: number;
}

export const catchWeightService = {
  // Weight Capture - Backend: /catch-weight
  captureWeight: async (data: CaptureRequest) => {
    const response = await api.post('/catch-weight/capture', data);
    return response.data;
  },

  quickCapture: async (data: Omit<CaptureRequest, 'pieces'>) => {
    const response = await api.post('/catch-weight/capture/quick', data);
    return response.data;
  },

  // Piece Management
  addPiece: async (entryId: string, data: { piece_number: number; weight: number; barcode?: string; notes?: string }) => {
    const response = await api.post(`/catch-weight/entries/${entryId}/pieces`, data);
    return response.data;
  },

  updatePiece: async (pieceId: string, data: Partial<CatchWeightPiece>) => {
    const response = await api.put(`/catch-weight/pieces/${pieceId}`, data);
    return response.data;
  },

  deletePiece: async (pieceId: string) => {
    const response = await api.delete(`/catch-weight/pieces/${pieceId}`);
    return response.data;
  },

  // Entry Retrieval
  getEntry: async (id: string): Promise<CatchWeightEntry> => {
    const response = await api.get(`/catch-weight/entries/${id}`);
    return response.data?.data || response.data;
  },

  getEntries: async (params?: any): Promise<CatchWeightEntry[]> => {
    const response = await api.get('/catch-weight/entries', { params });
    return response.data?.data || response.data || [];
  },

  getPiecesByEntry: async (entryId: string): Promise<CatchWeightPiece[]> => {
    const response = await api.get(`/catch-weight/entries/${entryId}/pieces`);
    return response.data?.data || response.data || [];
  },

  getEntryByReference: async (refType: string, refId: string, productId: string): Promise<CatchWeightEntry> => {
    const response = await api.get(`/catch-weight/reference/${refType}/${refId}/product/${productId}`);
    return response.data?.data || response.data;
  },

  // Product Configuration
  getProductConfig: async (productId: string): Promise<CatchWeightConfig> => {
    const response = await api.get(`/catch-weight/products/${productId}/config`);
    return response.data?.data || response.data;
  },

  updateProductConfig: async (productId: string, data: Partial<CatchWeightConfig>) => {
    const response = await api.put(`/catch-weight/products/${productId}/config`, data);
    return response.data;
  },

  // Reports
  getVarianceReport: async (params?: { product_id?: number; start_date?: string; end_date?: string }): Promise<VarianceReport[]> => {
    const response = await api.get('/catch-weight/reports/variance', { params });
    return response.data?.data || response.data || [];
  },

  getLotSummary: async (productId: string, lotNumber: string) => {
    const response = await api.get(`/catch-weight/reports/lot/${productId}/${lotNumber}`);
    return response.data?.data || response.data;
  },

  // Billing
  calculateBillingAdjustment: async (entryId: string) => {
    const response = await api.post('/catch-weight/billing/adjustment', { entry_id: parseInt(entryId) });
    return response.data?.data || response.data;
  },

  markAsBilled: async (entryId: string) => {
    const response = await api.post(`/catch-weight/entries/${entryId}/mark-billed`);
    return response.data;
  },

  // Validation
  validateWeight: async (data: { product_id: number; weight: number; expected_weight?: number }) => {
    const response = await api.post('/catch-weight/validate-weight', data);
    return response.data?.data || response.data;
  },
};
