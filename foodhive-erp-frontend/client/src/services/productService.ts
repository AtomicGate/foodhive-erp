import api from '@/lib/api';

// Handle double-wrapped response from backend
const unwrapData = (data: any) => {
  const innerData = data?.data?.data || data?.data || data || [];
  return Array.isArray(innerData) ? innerData : [];
};

export const productService = {
  // Products - Backend: /products
  getProducts: async (params?: any) => {
    const response = await api.get('/products/list', { params });
    console.log('Products raw response:', response.data);
    return unwrapData(response.data);
  },

  getProduct: async (id: string) => {
    const response = await api.get(`/products/get/${id}`);
    return response.data?.data || response.data;
  },

  getProductBySku: async (sku: string) => {
    const response = await api.get(`/products/sku/${sku}`);
    return response.data?.data || response.data;
  },

  getProductByBarcode: async (barcode: string) => {
    const response = await api.get(`/products/barcode/${barcode}`);
    return response.data?.data || response.data;
  },

  createProduct: async (data: any) => {
    const response = await api.post('/products/create', data);
    return response.data;
  },

  updateProduct: async (id: string, data: any) => {
    console.log('Updating product:', id, 'with data:', data);
    const response = await api.put(`/products/update/${id}`, data);
    console.log('Update response:', response.data);
    return response.data;
  },

  deleteProduct: async (id: string) => {
    const response = await api.delete(`/products/delete/${id}`);
    return response.data;
  },

  // Categories - Backend: /products/categories
  getCategories: async (params?: any) => {
    const response = await api.get('/products/categories/list', { params });
    return unwrapData(response.data);
  },

  getCategory: async (id: string) => {
    const response = await api.get(`/products/categories/get/${id}`);
    return response.data?.data || response.data;
  },

  createCategory: async (data: any) => {
    const response = await api.post('/products/categories/create', data);
    return response.data;
  },

  updateCategory: async (id: string, data: any) => {
    const response = await api.put(`/products/categories/update/${id}`, data);
    return response.data;
  },

  deleteCategory: async (id: string) => {
    const response = await api.delete(`/products/categories/delete/${id}`);
    return response.data;
  },

  // Units of Measure - Backend: /products/units
  getProductUnits: async (productId: string) => {
    const response = await api.get(`/products/units/product/${productId}`);
    return unwrapData(response.data);
  },

  addUnit: async (data: any) => {
    const response = await api.post('/products/units/add', data);
    return response.data;
  },

  updateUnit: async (id: string, data: any) => {
    const response = await api.put(`/products/units/update/${id}`, data);
    return response.data;
  },

  deleteUnit: async (id: string) => {
    const response = await api.delete(`/products/units/delete/${id}`);
    return response.data;
  },
};
