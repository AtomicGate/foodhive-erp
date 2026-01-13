import api from '@/lib/api';

// Dashboard service - aggregates data from multiple backend endpoints
// TODO: Backend could add /dashboard endpoints for optimized aggregation

export const dashboardService = {
  getStats: async () => {
    try {
      // Try to fetch real data from backend endpoints
      const [salesOrders, inventory, employees] = await Promise.allSettled([
        api.get('/sales-orders/list', { params: { per_page: 1 } }),
        api.get('/inventory/list', { params: { per_page: 1 } }),
        api.get('/employees/list', { params: { per_page: 1 } }),
      ]);

      // Calculate stats from responses
      const orderCount = salesOrders.status === 'fulfilled' 
        ? (salesOrders.value.data?.pagination?.total_items || salesOrders.value.data?.total || salesOrders.value.data?.data?.length || 0)
        : 0;
      
      const inventoryCount = inventory.status === 'fulfilled'
        ? (inventory.value.data?.pagination?.total_items || inventory.value.data?.total || inventory.value.data?.data?.length || 0)
        : 0;

      const employeeCount = employees.status === 'fulfilled'
        ? (employees.value.data?.pagination?.total_items || employees.value.data?.total || employees.value.data?.data?.length || 0)
        : 0;

      return {
        revenue: { total: 0, change: 0, trend: 'up' as const },
        orders: { total: orderCount, change: 0, trend: 'up' as const },
        inventory: { total: inventoryCount, change: 0, trend: 'up' as const },
        employees: { total: employeeCount, change: 0, trend: 'up' as const }
      };
    } catch (error) {
      // Return zeros if API fails
      return {
        revenue: { total: 0, change: 0, trend: 'up' as const },
        orders: { total: 0, change: 0, trend: 'up' as const },
        inventory: { total: 0, change: 0, trend: 'up' as const },
        employees: { total: 0, change: 0, trend: 'up' as const }
      };
    }
  },

  getRecentSales: async () => {
    try {
      // Fetch recent sales orders
      const response = await api.get('/sales-orders/list', { 
        params: { per_page: 5, sort: 'created_at', order: 'desc' } 
      });
      
      const orders = response.data?.data || response.data || [];
      
      if (orders.length > 0) {
        return orders.map((order: any) => ({
          name: order.customer_name || 'Customer',
          email: order.customer_email || 'customer@email.com',
          amount: `+$${(order.total_amount || 0).toFixed(2)}`,
          status: order.status || 'Pending'
        }));
      }
      
      // Return empty array if no orders
      return [];
    } catch (error) {
      // Return empty array if API fails
      return [];
    }
  },

  getRevenueChart: async (period: 'day' | 'week' | 'month' | 'year' = 'month') => {
    // TODO: Backend should implement /dashboard/revenue-chart endpoint
    // For now, return empty chart data - real data will come from actual sales
    const emptyData: Record<string, { name: string; total: number }[]> = {
      day: [
        { name: "00:00", total: 0 },
        { name: "04:00", total: 0 },
        { name: "08:00", total: 0 },
        { name: "12:00", total: 0 },
        { name: "16:00", total: 0 },
        { name: "20:00", total: 0 },
      ],
      week: [
        { name: "Mon", total: 0 },
        { name: "Tue", total: 0 },
        { name: "Wed", total: 0 },
        { name: "Thu", total: 0 },
        { name: "Fri", total: 0 },
        { name: "Sat", total: 0 },
        { name: "Sun", total: 0 },
      ],
      month: [
        { name: "Week 1", total: 0 },
        { name: "Week 2", total: 0 },
        { name: "Week 3", total: 0 },
        { name: "Week 4", total: 0 },
      ],
      year: [
        { name: "Jan", total: 0 },
        { name: "Feb", total: 0 },
        { name: "Mar", total: 0 },
        { name: "Apr", total: 0 },
        { name: "May", total: 0 },
        { name: "Jun", total: 0 },
        { name: "Jul", total: 0 },
        { name: "Aug", total: 0 },
        { name: "Sep", total: 0 },
        { name: "Oct", total: 0 },
        { name: "Nov", total: 0 },
        { name: "Dec", total: 0 },
      ],
    };
    
    return emptyData[period] || emptyData.month;
  }
};
