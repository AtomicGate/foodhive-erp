import React from 'react';
import { usePermissions } from '@/hooks/usePermissions';
import { useLocation } from 'wouter';
import { useAuth } from '@/contexts/AuthContext';

interface PermissionGuardProps {
  permission: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
  redirectTo?: string;
}

export const PermissionGuard: React.FC<PermissionGuardProps> = ({ 
  permission, 
  children, 
  fallback = null,
  redirectTo
}) => {
  const { hasPermission } = usePermissions();
  const [, setLocation] = useLocation();

  if (hasPermission(permission)) {
    return <>{children}</>;
  }

  if (redirectTo) {
    setLocation(redirectTo);
    return null;
  }

  return <>{fallback}</>;
};

interface RoleGuardProps {
  role: string | string[];
  children: React.ReactNode;
  fallback?: React.ReactNode;
  redirectTo?: string;
}

export const RoleGuard: React.FC<RoleGuardProps> = ({ 
  role, 
  children, 
  fallback = null,
  redirectTo
}) => {
  const { hasRole } = useAuth();
  const [, setLocation] = useLocation();

  // @ts-ignore - hasRole expects UserRole type but we're passing string
  if (hasRole(role)) {
    return <>{children}</>;
  }

  if (redirectTo) {
    setLocation(redirectTo);
    return null;
  }

  return <>{fallback}</>;
};
