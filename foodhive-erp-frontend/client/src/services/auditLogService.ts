import api from '@/lib/api';

export interface AuditLog {
  id: string;
  entityType: string;
  entityId: string;
  action: 'CREATE' | 'UPDATE' | 'DELETE';
  userId: string;
  userName: string;
  timestamp: string;
  details?: string;
  changes?: Record<string, { old: any; new: any }>;
}

export interface PaginatedAuditLogs {
  items: AuditLog[];
  total: number;
}

export const auditLogService = {
  getLogs: async (
    page: number = 1, 
    limit: number = 10, 
    entityType?: string,
    entityId?: string
  ): Promise<PaginatedAuditLogs> => {
    // In a real implementation, this would call the backend API
    // const response = await api.get('/audit-logs', { params: { page, limit, entityType, entityId } });
    // return response.data;

    // Mock data for demonstration
    return new Promise((resolve) => {
      setTimeout(() => {
        const logs: AuditLog[] = Array.from({ length: limit }).map((_, i) => ({
          id: `log-${Date.now()}-${i}`,
          entityType: entityType || 'Department',
          entityId: entityId || `dept-${i}`,
          action: ['CREATE', 'UPDATE', 'DELETE'][Math.floor(Math.random() * 3)] as any,
          userId: 'user-1',
          userName: 'Admin User',
          timestamp: new Date(Date.now() - Math.random() * 1000000000).toISOString(),
          details: 'Performed an action on the entity',
        }));
        
        resolve({
          items: logs,
          total: 100
        });
      }, 500);
    });
  },

  logAction: async (
    entityType: string,
    entityId: string,
    action: 'CREATE' | 'UPDATE' | 'DELETE',
    details?: string,
    changes?: Record<string, { old: any; new: any }>
  ) => {
    // await api.post('/audit-logs', { entityType, entityId, action, details, changes });
    console.log('Audit Log:', { entityType, entityId, action, details, changes });
  }
};
