import api from "@/lib/api";

export interface PriceLevel {
  id: number;
  name: string;
  description: string;
  priority: number;
  is_active: boolean;
}

export interface ProductPrice {
  id: number;
  product_id: number;
  product_name?: string;
  price_level_id?: number;
  price: number;
  min_quantity?: number;
  effective_date?: string;
  expiration_date?: string;
}

export interface CustomerPrice {
  id: number;
  customer_id: number;
  product_id: number;
  price: number;
  min_quantity?: number;
  effective_date?: string;
  expiration_date?: string;
}

export interface Contract {
  id: number;
  customer_id: number;
  name: string;
  start_date: string;
  end_date: string;
  is_active: boolean;
  lines: ContractLine[];
}

export interface ContractLine {
  product_id: number;
  price: number;
  min_quantity?: number;
}

export interface Promotion {
  id: number;
  name: string;
  start_date: string;
  end_date: string;
  discount_type: "PERCENT" | "FIXED";
  discount_value: number;
  is_active: boolean;
}

export interface PriceLookupResult {
  product_id: number;
  price: number;
  price_source: string;
  contract_id?: number;
  list_price?: number;
  discount_percent?: number;
  cost?: number;
  margin_percent?: number;
}

export const pricingService = {
  // Price Lookup - Backend: /pricing
  lookupPrice: async (data: { product_id: number; customer_id: number; quantity?: number; as_of_date?: string }): Promise<PriceLookupResult> => {
    const response = await api.post("/pricing/lookup", data);
    return response.data?.data || response.data;
  },

  batchLookupPrices: async (items: Array<{ product_id: number; customer_id: number; quantity?: number }>): Promise<PriceLookupResult[]> => {
    const response = await api.post("/pricing/lookup/batch", { items });
    return response.data?.data || response.data || [];
  },

  checkMargin: async (params: { product_id: number; price: number }) => {
    const response = await api.get("/pricing/check-margin", { params });
    return response.data?.data || response.data;
  },

  // Product Prices
  getProductPrices: async (productId: string): Promise<ProductPrice[]> => {
    const response = await api.get(`/pricing/product/${productId}`);
    return response.data?.data || response.data || [];
  },

  setProductPrice: async (data: Partial<ProductPrice>) => {
    const response = await api.post("/pricing/product/set", data);
    return response.data;
  },

  deleteProductPrice: async (id: string) => {
    const response = await api.delete(`/pricing/product/price/${id}`);
    return response.data;
  },

  // Customer Prices
  getCustomerPrices: async (customerId: string): Promise<CustomerPrice[]> => {
    const response = await api.get(`/pricing/customer/${customerId}`);
    return response.data?.data || response.data || [];
  },

  setCustomerPrice: async (data: Partial<CustomerPrice>) => {
    const response = await api.post("/pricing/customer/set", data);
    return response.data;
  },

  deleteCustomerPrice: async (id: string) => {
    const response = await api.delete(`/pricing/customer/price/${id}`);
    return response.data;
  },

  // Contracts
  getContracts: async (params?: any): Promise<Contract[]> => {
    const response = await api.get("/pricing/contracts/list", { params });
    return response.data?.data || response.data || [];
  },

  getContract: async (id: string): Promise<Contract> => {
    const response = await api.get(`/pricing/contracts/get/${id}`);
    return response.data?.data || response.data;
  },

  createContract: async (data: Partial<Contract>) => {
    const response = await api.post("/pricing/contracts/create", data);
    return response.data;
  },

  deactivateContract: async (id: string) => {
    const response = await api.post(`/pricing/contracts/${id}/deactivate`);
    return response.data;
  },

  // Promotions
  getPromotions: async (params?: any): Promise<Promotion[]> => {
    const response = await api.get("/pricing/promotions/list", { params });
    return response.data?.data || response.data || [];
  },

  getPromotion: async (id: string): Promise<Promotion> => {
    const response = await api.get(`/pricing/promotions/get/${id}`);
    return response.data?.data || response.data;
  },

  createPromotion: async (data: Partial<Promotion>) => {
    const response = await api.post("/pricing/promotions/create", data);
    return response.data;
  },

  deactivatePromotion: async (id: string) => {
    const response = await api.post(`/pricing/promotions/${id}/deactivate`);
    return response.data;
  },

  // Product Costs
  getProductCost: async (productId: string) => {
    const response = await api.get(`/pricing/costs/${productId}`);
    return response.data?.data || response.data;
  },

  updateProductCost: async (data: { product_id: number; cost: number; effective_date?: string }) => {
    const response = await api.post("/pricing/costs/update", data);
    return response.data;
  },

  // Mass Update
  massUpdatePrices: async (data: { product_ids: number[]; adjustment_type: "PERCENT" | "FIXED"; adjustment_value: number; price_level_id?: number }) => {
    const response = await api.post("/pricing/mass-update", data);
    return response.data;
  },

  // Price List
  getPriceList: async (params?: any) => {
    const response = await api.get("/pricing/list", { params });
    return response.data?.data || response.data || [];
  },
};
