import api from "@/lib/api";

export interface ARAging {
  current: number;
  days_1_30: number;
  days_31_60: number;
  days_61_90: number;
  over_90: number;
  total: number;
}

export interface ARAgingReport {
  customer_id: number;
  customer_name: string;
  current: number;
  days_1_30: number;
  days_31_60: number;
  days_61_90: number;
  over_90: number;
  total: number;
}

export interface ARInvoice {
  id: number;
  invoice_number: string;
  customer_id: number;
  customer_name?: string;
  invoice_date: string;
  due_date: string;
  total_amount: number;
  paid_amount: number;
  balance: number;
  status: "DRAFT" | "POSTED" | "PARTIAL" | "PAID" | "VOIDED";
}

export interface ARPayment {
  id: number;
  payment_number: string;
  customer_id: number;
  payment_date: string;
  amount: number;
  payment_method: string;
}

export interface APInvoice {
  id: number;
  invoice_number: string;
  vendor_id: number;
  vendor_name?: string;
  invoice_date: string;
  due_date: string;
  total_amount: number;
  paid_amount: number;
  balance: number;
  status: "PENDING" | "APPROVED" | "PARTIAL" | "PAID" | "VOIDED";
}

export const financialService = {
  // ===== Accounts Receivable (AR) =====
  
  // AR Aging
  getARAging: async (customerId?: string): Promise<ARAging> => {
    if (customerId) {
      const response = await api.get(`/ar/aging/${customerId}`);
      return response.data?.data || response.data;
    }
    const response = await api.get("/ar/aging/report");
    // Aggregate totals from report
    const data = response.data?.data || response.data || [];
    if (Array.isArray(data)) {
      return data.reduce((acc: ARAging, item: ARAgingReport) => ({
        current: acc.current + (item.current || 0),
        days_1_30: acc.days_1_30 + (item.days_1_30 || 0),
        days_31_60: acc.days_31_60 + (item.days_31_60 || 0),
        days_61_90: acc.days_61_90 + (item.days_61_90 || 0),
        over_90: acc.over_90 + (item.over_90 || 0),
        total: acc.total + (item.total || 0),
      }), { current: 0, days_1_30: 0, days_31_60: 0, days_61_90: 0, over_90: 0, total: 0 });
    }
    return data;
  },

  getARAgingReport: async (): Promise<ARAgingReport[]> => {
    const response = await api.get("/ar/aging/report");
    return response.data?.data || response.data || [];
  },

  getCustomerAging: async (customerId: string): Promise<ARAging> => {
    const response = await api.get(`/ar/aging/${customerId}`);
    return response.data?.data || response.data;
  },

  // AR Invoices
  getARInvoices: async (params?: any): Promise<ARInvoice[]> => {
    const response = await api.get("/ar/invoices/list", { params });
    return response.data?.data || response.data || [];
  },

  getARInvoice: async (id: string): Promise<ARInvoice> => {
    const response = await api.get(`/ar/invoices/get/${id}`);
    return response.data?.data || response.data;
  },

  getOverdueInvoices: async (): Promise<ARInvoice[]> => {
    const response = await api.get("/ar/overdue");
    return response.data?.data || response.data || [];
  },

  createARInvoice: async (data: any) => {
    const response = await api.post("/ar/invoices/create", data);
    return response.data;
  },

  postARInvoice: async (id: string) => {
    const response = await api.post(`/ar/invoices/${id}/post`);
    return response.data;
  },

  voidARInvoice: async (id: string) => {
    const response = await api.post(`/ar/invoices/${id}/void`);
    return response.data;
  },

  // AR Payments
  getARPayments: async (params?: any): Promise<ARPayment[]> => {
    const response = await api.get("/ar/payments/list", { params });
    return response.data?.data || response.data || [];
  },

  getARPayment: async (id: string): Promise<ARPayment> => {
    const response = await api.get(`/ar/payments/get/${id}`);
    return response.data?.data || response.data;
  },

  createARPayment: async (data: any) => {
    const response = await api.post("/ar/payments/create", data);
    return response.data;
  },

  // AR Credit
  getCustomerCredit: async (customerId: string) => {
    const response = await api.get(`/ar/credit/${customerId}`);
    return response.data?.data || response.data;
  },

  checkCredit: async (customerId: string, amount: number) => {
    const response = await api.get(`/ar/credit/${customerId}/check`, { params: { amount } });
    return response.data?.data || response.data;
  },

  updateCreditLimit: async (customerId: string, limit: number) => {
    const response = await api.put(`/ar/credit/${customerId}/limit`, { credit_limit: limit });
    return response.data;
  },

  // AR Statement
  getCustomerStatement: async (customerId: string, params?: any) => {
    const response = await api.get(`/ar/statement/${customerId}`, { params });
    return response.data?.data || response.data;
  },

  // ===== Accounts Payable (AP) =====

  // AP Invoices
  getAPInvoices: async (params?: any): Promise<APInvoice[]> => {
    const response = await api.get("/ap/invoices/list", { params });
    return response.data?.data || response.data || [];
  },

  getAPInvoice: async (id: string): Promise<APInvoice> => {
    const response = await api.get(`/ap/invoices/get/${id}`);
    return response.data?.data || response.data;
  },

  createAPInvoice: async (data: any) => {
    const response = await api.post("/ap/invoices/create", data);
    return response.data;
  },

  approveAPInvoice: async (id: string) => {
    const response = await api.post(`/ap/invoices/${id}/approve`);
    return response.data;
  },

  voidAPInvoice: async (id: string) => {
    const response = await api.post(`/ap/invoices/${id}/void`);
    return response.data;
  },

  // AP Payments
  getAPPayments: async (params?: any) => {
    const response = await api.get("/ap/payments/list", { params });
    return response.data?.data || response.data || [];
  },

  createAPPayment: async (data: any) => {
    const response = await api.post("/ap/payments/create", data);
    return response.data;
  },

  voidAPPayment: async (id: string) => {
    const response = await api.post(`/ap/payments/${id}/void`);
    return response.data;
  },

  // AP Aging
  getVendorBalance: async (vendorId: string) => {
    const response = await api.get(`/ap/balance/${vendorId}`);
    return response.data?.data || response.data;
  },

  getVendorAging: async (vendorId: string) => {
    const response = await api.get(`/ap/aging/${vendorId}`);
    return response.data?.data || response.data;
  },

  getAPAgingReport: async () => {
    const response = await api.get("/ap/aging/report");
    return response.data?.data || response.data || [];
  },

  getDueInvoices: async () => {
    const response = await api.get("/ap/due");
    return response.data?.data || response.data || [];
  },

  getOverdueAPInvoices: async () => {
    const response = await api.get("/ap/overdue");
    return response.data?.data || response.data || [];
  },

  // ===== General Ledger (GL) =====

  // Chart of Accounts
  getAccounts: async (params?: any) => {
    const response = await api.get("/gl/accounts", { params });
    return response.data?.data || response.data || [];
  },

  getAccount: async (id: string) => {
    const response = await api.get(`/gl/accounts/${id}`);
    return response.data?.data || response.data;
  },

  getChartOfAccounts: async () => {
    const response = await api.get("/gl/chart-of-accounts");
    return response.data?.data || response.data || [];
  },

  createAccount: async (data: any) => {
    const response = await api.post("/gl/accounts", data);
    return response.data;
  },

  updateAccount: async (id: string, data: any) => {
    const response = await api.put(`/gl/accounts/${id}`, data);
    return response.data;
  },

  deleteAccount: async (id: string) => {
    const response = await api.delete(`/gl/accounts/${id}`);
    return response.data;
  },

  // Fiscal Years & Periods
  getFiscalYears: async () => {
    const response = await api.get("/gl/fiscal-years");
    return response.data?.data || response.data || [];
  },

  getCurrentFiscalYear: async () => {
    const response = await api.get("/gl/fiscal-years/current");
    return response.data?.data || response.data;
  },

  getPeriods: async (fiscalYearId: string) => {
    const response = await api.get(`/gl/fiscal-years/${fiscalYearId}/periods`);
    return response.data?.data || response.data || [];
  },

  getCurrentPeriod: async () => {
    const response = await api.get("/gl/periods/current");
    return response.data?.data || response.data;
  },

  // Journal Entries
  getJournalEntries: async (params?: any) => {
    const response = await api.get("/gl/journal-entries", { params });
    return response.data?.data || response.data || [];
  },

  getJournalEntry: async (id: string) => {
    const response = await api.get(`/gl/journal-entries/${id}`);
    return response.data?.data || response.data;
  },

  createJournalEntry: async (data: any) => {
    const response = await api.post("/gl/journal-entries", data);
    return response.data;
  },

  postJournalEntry: async (id: string) => {
    const response = await api.post(`/gl/journal-entries/${id}/post`);
    return response.data;
  },

  // Financial Reports
  getTrialBalance: async (params?: any) => {
    const response = await api.get("/gl/reports/trial-balance", { params });
    return response.data?.data || response.data;
  },

  getIncomeStatement: async (params?: any) => {
    const response = await api.get("/gl/reports/income-statement", { params });
    return response.data?.data || response.data;
  },

  getBalanceSheet: async (params?: any) => {
    const response = await api.get("/gl/reports/balance-sheet", { params });
    return response.data?.data || response.data;
  },

  getAccountActivity: async (accountId: string, params?: any) => {
    const response = await api.get(`/gl/reports/account-activity/${accountId}`, { params });
    return response.data?.data || response.data;
  },
};
