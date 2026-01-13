import React, { createContext, useContext, useState, useEffect, useRef, useCallback } from 'react';
import { authService } from '@/services/authService';
import { toast } from 'sonner';

export type UserRole = 'admin' | 'sales' | 'warehouse' | 'finance';

interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (credentials: any) => Promise<void>;
  logout: () => void;
  hasRole: (role: UserRole | UserRole[]) => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Session timeout in milliseconds (e.g., 30 minutes)
const SESSION_TIMEOUT = 30 * 60 * 1000; 

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const logout = useCallback(() => {
    authService.logout();
    setUser(null);
    localStorage.removeItem('user');
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
  }, []);

  const resetTimeout = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    
    if (user) {
      timeoutRef.current = setTimeout(() => {
        toast.error("Session expired due to inactivity");
        logout();
      }, SESSION_TIMEOUT);
    }
  }, [user, logout]);

  useEffect(() => {
    // Check for existing token and load user
    const token = localStorage.getItem('token');
    if (token) {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        setUser(JSON.parse(storedUser));
      } else {
        // Default fallback for demo
        setUser({
          id: '1',
          name: 'John Doe',
          email: 'john@foodhive.com',
          role: 'admin'
        });
      }
    }
    setIsLoading(false);
  }, []);

  // Setup activity listeners
  useEffect(() => {
    if (!user) return;

    const events = ['mousedown', 'keydown', 'scroll', 'touchstart'];
    
    const handleActivity = () => {
      resetTimeout();
    };

    // Initial reset
    resetTimeout();

    // Add listeners
    events.forEach(event => {
      window.addEventListener(event, handleActivity);
    });

    return () => {
      // Cleanup
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      events.forEach(event => {
        window.removeEventListener(event, handleActivity);
      });
    };
  }, [user, resetTimeout]);

  const login = async (credentials: any) => {
    const data = await authService.login(credentials);
    localStorage.setItem('token', data.token);
    
    // Map backend role to frontend role type
    const roleMapping: Record<string, UserRole> = {
      'Admin': 'admin',
      'Administrator': 'admin',
      'admin': 'admin',
      'Sales': 'sales',
      'Sales Rep': 'sales',
      'sales': 'sales',
      'Warehouse': 'warehouse',
      'Warehouse Manager': 'warehouse',
      'warehouse': 'warehouse',
      'Finance': 'finance',
      'Accountant': 'finance',
      'finance': 'finance',
    };
    
    const userData: User = {
      id: String(data.user?.id || '1'),
      name: data.user?.name || 'User',
      email: data.user?.email || credentials.email,
      role: roleMapping[data.user?.role] || 'admin'
    };
    
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
  };

  const hasRole = (requiredRole: UserRole | UserRole[]) => {
    if (!user) return false;
    if (user.role === 'admin') return true; 
    
    if (Array.isArray(requiredRole)) {
      return requiredRole.includes(user.role);
    }
    return user.role === requiredRole;
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, logout, hasRole }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
