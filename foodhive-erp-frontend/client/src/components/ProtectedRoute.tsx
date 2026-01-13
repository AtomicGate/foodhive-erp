import { useAuth, UserRole } from "@/contexts/AuthContext";
import { Redirect, Route, RouteProps } from "wouter";

interface ProtectedRouteProps extends RouteProps {
  component: React.ComponentType<any>;
  allowedRoles?: UserRole[];
}

export default function ProtectedRoute({ 
  component: Component, 
  allowedRoles, 
  ...rest 
}: ProtectedRouteProps) {
  const { user, isLoading, hasRole } = useAuth();

  // Show nothing while loading auth state
  if (isLoading) {
    return null; 
  }

  // If not logged in, redirect to login
  if (!user) {
    return <Route {...rest} component={() => <Redirect to="/login" />} />;
  }

  // If roles are specified and user doesn't have permission, redirect to forbidden
  if (allowedRoles && !hasRole(allowedRoles)) {
    return <Route {...rest} component={() => <Redirect to="/403" />} />;
  }

  // Otherwise render the component
  return <Route {...rest} component={Component} />;
}
