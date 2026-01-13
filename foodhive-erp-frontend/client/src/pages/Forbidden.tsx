import { Button } from "@/components/ui/button";
import { ShieldAlert } from "lucide-react";
import { Link } from "wouter";

export default function Forbidden() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-background p-4">
      <div className="text-center space-y-6 max-w-md">
        <div className="flex justify-center">
          <div className="h-24 w-24 rounded-full bg-rose-100 flex items-center justify-center">
            <ShieldAlert className="h-12 w-12 text-rose-600" />
          </div>
        </div>
        
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Access Denied</h1>
          <p className="text-muted-foreground">
            You don't have permission to access this page. Please contact your administrator if you believe this is a mistake.
          </p>
        </div>

        <div className="flex justify-center gap-4">
          <Link href="/">
            <Button variant="default" size="lg">
              Return to Dashboard
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
