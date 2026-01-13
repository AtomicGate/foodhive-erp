import { useAuth } from "@/contexts/AuthContext";

export const usePermissions = () => {
  const { user } = useAuth();
  
  // In a real app, permissions would come from the user object or a separate API call
  // For now, we'll simulate permissions based on the user's role or a permissions array
  // Assuming user object has a permissions array: user.permissions = ['view_dashboard', 'edit_inventory']
  
  const hasPermission = (permission: string) => {
    if (!user) return false;
    // Admin has all permissions
    if (user.role === 'admin') return true;
    
    // Check if user has the specific permission
    // This assumes the user object structure has been updated to include permissions
    // If not, we might need to fetch them or rely on role-based logic for now
    // @ts-ignore - permissions might not be in the type definition yet
    return user.permissions?.includes(permission) || false;
  };

  const hasAnyPermission = (permissions: string[]) => {
    if (!user) return false;
    if (user.role === 'admin') return true;
    // @ts-ignore
    return permissions.some(p => user.permissions?.includes(p));
  };

  const hasAllPermissions = (permissions: string[]) => {
    if (!user) return false;
    if (user.role === 'admin') return true;
    // @ts-ignore
    return permissions.every(p => user.permissions?.includes(p));
  };

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    // @ts-ignore
    userPermissions: user?.permissions || []
  };
};
