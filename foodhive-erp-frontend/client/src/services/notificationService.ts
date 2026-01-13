import api from '@/lib/api';

export interface Notification {
  id: string;
  title: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
  read: boolean;
  timestamp: string;
  link?: string;
}

export const notificationService = {
  getNotifications: async (): Promise<Notification[]> => {
    // In a real app, this would fetch from the backend
    // const response = await api.get('/notifications');
    // return response.data;

    // Mock data
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([
          {
            id: '1',
            title: 'Low Inventory Alert',
            message: 'Product "Ground Beef" is below safety stock level.',
            type: 'warning',
            read: false,
            timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString(), // 30 mins ago
            link: '/inventory'
          },
          {
            id: '2',
            title: 'New Order Received',
            message: 'Order #ORD-2024-001 has been placed by ABC Restaurant.',
            type: 'success',
            read: false,
            timestamp: new Date(Date.now() - 1000 * 60 * 60).toISOString(), // 1 hour ago
            link: '/sales-orders'
          },
          {
            id: '3',
            title: 'System Maintenance',
            message: 'Scheduled maintenance tonight at 2:00 AM.',
            type: 'info',
            read: true,
            timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(), // 1 day ago
          }
        ]);
      }, 500);
    });
  },

  markAsRead: async (id: string) => {
    // await api.put(`/notifications/${id}/read`);
    console.log(`Notification ${id} marked as read`);
  },

  markAllAsRead: async () => {
    // await api.put('/notifications/read-all');
    console.log('All notifications marked as read');
  }
};
