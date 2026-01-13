import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import NotFound from "@/pages/NotFound";
import Forbidden from "@/pages/Forbidden";
import { Route, Switch } from "wouter";
import ErrorBoundary from "./components/ErrorBoundary";
import { ThemeProvider } from "./contexts/ThemeContext";
import { AuthProvider } from "./contexts/AuthContext";
import Layout from "./components/Layout";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import SalesOrderList from "./pages/sales/SalesOrderList";
import SalesOrderForm from "./pages/sales/SalesOrderForm";
import PickList from "./pages/sales/PickList";
import Invoice from "./pages/sales/Invoice";
import PurchaseOrderList from "./pages/purchasing/PurchaseOrderList";
import InventoryList from "./pages/inventory/InventoryList";
import ProductList from "./pages/products/ProductList";
import EntityList from "./pages/entities/EntityList";
import MasterData from "./pages/admin/MasterData";
import ARDashboard from "./pages/financials/ARDashboard";
import APDashboard from "./pages/financials/APDashboard";
import GLDashboard from "./pages/financials/GLDashboard";
import ChartOfAccounts from "./pages/financials/ChartOfAccounts";
import JournalEntries from "./pages/financials/JournalEntries";
import TrialBalance from "./pages/financials/TrialBalance";
import CatchWeight from "./pages/operations/CatchWeight";
import PricingManagement from "./pages/pricing/PricingManagement";
import CustomerList from "./pages/customers/CustomerList";
import VendorList from "./pages/vendors/VendorList";
import EmployeeList from "./pages/employees/EmployeeList";
import UserProfile from "./pages/UserProfile";
import ProtectedRoute from "./components/ProtectedRoute";

function Router() {
  return (
    <Switch>
      <Route path="/login" component={Login} />
      <Route path="/403" component={Forbidden} />
      
      {/* Protected Routes wrapped in Layout */}
      <ProtectedRoute 
        path="/" 
        component={() => (
          <Layout>
            <Dashboard />
          </Layout>
        )} 
      />
      
      <ProtectedRoute 
        path="/profile" 
        component={() => (
          <Layout>
            <UserProfile />
          </Layout>
        )} 
      />
      
      {/* Admin Only Routes */}
      <ProtectedRoute 
        path="/admin/:type" 
        allowedRoles={['admin']}
        component={() => (
          <Layout>
            <MasterData />
          </Layout>
        )} 
      />

      {/* Financials - Admin & Finance */}
      <ProtectedRoute 
        path="/financials/ar" 
        allowedRoles={['admin', 'finance']}
        component={() => (
          <Layout>
            <ARDashboard />
          </Layout>
        )} 
      />

      <ProtectedRoute 
        path="/financials/ap" 
        allowedRoles={['admin', 'finance']}
        component={() => (
          <Layout>
            <APDashboard />
          </Layout>
        )} 
      />

      {/* General Ledger Routes */}
      <ProtectedRoute 
        path="/gl" 
        allowedRoles={['admin', 'finance']}
        component={() => (
          <Layout>
            <GLDashboard />
          </Layout>
        )} 
      />

      <ProtectedRoute 
        path="/gl/accounts" 
        allowedRoles={['admin', 'finance']}
        component={() => (
          <Layout>
            <ChartOfAccounts />
          </Layout>
        )} 
      />

      <ProtectedRoute 
        path="/gl/journal-entries" 
        allowedRoles={['admin', 'finance']}
        component={() => (
          <Layout>
            <JournalEntries />
          </Layout>
        )} 
      />

      <ProtectedRoute 
        path="/gl/reports/trial-balance" 
        allowedRoles={['admin', 'finance']}
        component={() => (
          <Layout>
            <TrialBalance />
          </Layout>
        )} 
      />

      {/* Operations - Admin & Warehouse */}
      <ProtectedRoute 
        path="/operations/catch-weight" 
        allowedRoles={['admin', 'warehouse']}
        component={() => (
          <Layout>
            <CatchWeight />
          </Layout>
        )} 
      />

      {/* Pricing - Admin & Sales */}
      <ProtectedRoute 
        path="/pricing" 
        allowedRoles={['admin', 'sales']}
        component={() => (
          <Layout>
            <PricingManagement />
          </Layout>
        )} 
      />
      
      {/* Products - All Roles */}
      <ProtectedRoute 
        path="/products" 
        component={() => (
          <Layout>
            <ProductList />
          </Layout>
        )} 
      />

      {/* Customers */}
      <ProtectedRoute 
        path="/customers" 
        allowedRoles={['admin', 'sales', 'finance']}
        component={() => (
          <Layout>
            <CustomerList />
          </Layout>
        )} 
      />

      {/* Vendors */}
      <ProtectedRoute 
        path="/vendors" 
        allowedRoles={['admin', 'warehouse', 'finance']}
        component={() => (
          <Layout>
            <VendorList />
          </Layout>
        )} 
      />

      {/* Employees */}
      <ProtectedRoute 
        path="/employees" 
        allowedRoles={['admin']}
        component={() => (
          <Layout>
            <EmployeeList />
          </Layout>
        )} 
      />

      {/* Entities (Legacy/Generic - Employees) - All Roles */}
      <ProtectedRoute 
        path="/entities/:type" 
        component={() => (
          <Layout>
            <EntityList />
          </Layout>
        )} 
      />
      
      {/* Inventory - All Roles */}
      <ProtectedRoute 
        path="/inventory" 
        component={() => (
          <Layout>
            <InventoryList />
          </Layout>
        )} 
      />
      
      {/* Purchasing - Admin, Warehouse, Finance */}
      <ProtectedRoute 
        path="/purchase-orders" 
        allowedRoles={['admin', 'warehouse', 'finance']}
        component={() => (
          <Layout>
            <PurchaseOrderList />
          </Layout>
        )} 
      />
      
      {/* Sales - Admin, Sales, Finance */}
      <ProtectedRoute 
        path="/sales-orders" 
        allowedRoles={['admin', 'sales', 'finance']}
        component={() => (
          <Layout>
            <SalesOrderList />
          </Layout>
        )} 
      />
      
      <ProtectedRoute 
        path="/sales-orders/new" 
        allowedRoles={['admin', 'sales']}
        component={() => (
          <Layout>
            <SalesOrderForm />
          </Layout>
        )} 
      />

      <ProtectedRoute 
        path="/sales-orders/picklist/:id" 
        allowedRoles={['admin', 'warehouse']}
        component={() => (
          <Layout>
            <PickList />
          </Layout>
        )} 
      />

      <ProtectedRoute 
        path="/sales-orders/invoice/:id" 
        allowedRoles={['admin', 'finance', 'sales']}
        component={() => (
          <Layout>
            <Invoice />
          </Layout>
        )} 
      />

      <Route path="/404" component={NotFound} />
      <Route component={NotFound} />
    </Switch>
  );
}

const queryClient = new QueryClient();

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
      <ThemeProvider
        defaultTheme="light"
        // switchable
      >
        <AuthProvider>
          <TooltipProvider>
            <Toaster />
            <Router />
        </TooltipProvider>
        </AuthProvider>
      </ThemeProvider>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}
export default App;
