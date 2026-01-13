import { useState } from "react";
import { Link } from "wouter";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { DataTable } from "@/components/ui/data-table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  Plus, 
  MoreHorizontal, 
  ArrowUpDown,
  FileText,
  Loader2,
  Eye,
  CheckCircle,
  XCircle,
  RotateCcw,
  ChevronRight
} from "lucide-react";
import { toast } from "sonner";
import { financialService } from "@/services/financialService";

interface JournalEntry {
  id: number;
  entry_number: string;
  entry_date: string;
  entry_type: string;
  description: string;
  total_debit: number;
  total_credit: number;
  status: string;
  created_by: number;
  posted_at?: string;
  posted_by?: number;
}

export default function JournalEntries() {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState("all");

  const { data: entries, isLoading, error } = useQuery({
    queryKey: ['journalEntries', statusFilter],
    queryFn: () => financialService.getJournalEntries(
      statusFilter !== 'all' ? { status: statusFilter } : {}
    )
  });

  const postMutation = useMutation({
    mutationFn: (id: string) => financialService.postJournalEntry(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['journalEntries'] });
      toast.success("Journal entry posted successfully");
    },
    onError: () => {
      toast.error("Failed to post journal entry");
    }
  });

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value);
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "DRAFT":
        return <Badge variant="secondary" className="bg-gray-100 text-gray-800">Draft</Badge>;
      case "POSTED":
        return <Badge variant="default" className="bg-emerald-100 text-emerald-800">Posted</Badge>;
      case "VOIDED":
        return <Badge variant="destructive">Voided</Badge>;
      case "REVERSED":
        return <Badge variant="outline" className="bg-amber-100 text-amber-800">Reversed</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const getEntryTypeBadge = (type: string) => {
    switch (type) {
      case "STANDARD":
        return <Badge variant="outline">Standard</Badge>;
      case "ADJUSTING":
        return <Badge variant="outline" className="bg-blue-100 text-blue-800">Adjusting</Badge>;
      case "CLOSING":
        return <Badge variant="outline" className="bg-purple-100 text-purple-800">Closing</Badge>;
      case "REVERSING":
        return <Badge variant="outline" className="bg-amber-100 text-amber-800">Reversing</Badge>;
      default:
        return <Badge variant="outline">{type}</Badge>;
    }
  };

  const columns: ColumnDef<JournalEntry>[] = [
    {
      accessorKey: "entry_number",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Entry #
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => (
        <div className="font-mono font-medium">{row.getValue("entry_number")}</div>
      ),
    },
    {
      accessorKey: "entry_date",
      header: "Date",
      cell: ({ row }) => {
        const date = row.getValue("entry_date") as string;
        return new Date(date).toLocaleDateString();
      },
    },
    {
      accessorKey: "entry_type",
      header: "Type",
      cell: ({ row }) => getEntryTypeBadge(row.getValue("entry_type")),
    },
    {
      accessorKey: "description",
      header: "Description",
      cell: ({ row }) => (
        <div className="max-w-[300px] truncate">
          {row.getValue("description")}
        </div>
      ),
    },
    {
      accessorKey: "total_debit",
      header: () => <div className="text-right">Debit</div>,
      cell: ({ row }) => (
        <div className="text-right font-medium">
          {formatCurrency(row.getValue("total_debit"))}
        </div>
      ),
    },
    {
      accessorKey: "total_credit",
      header: () => <div className="text-right">Credit</div>,
      cell: ({ row }) => (
        <div className="text-right font-medium">
          {formatCurrency(row.getValue("total_credit"))}
        </div>
      ),
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => getStatusBadge(row.getValue("status")),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const entry = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <Link href={`/gl/journal-entries/${entry.id}`}>
                <DropdownMenuItem>
                  <Eye className="mr-2 h-4 w-4" /> View Details
                </DropdownMenuItem>
              </Link>
              {entry.status === "DRAFT" && (
                <DropdownMenuItem onClick={() => postMutation.mutate(String(entry.id))}>
                  <CheckCircle className="mr-2 h-4 w-4" /> Post Entry
                </DropdownMenuItem>
              )}
              {entry.status === "POSTED" && (
                <>
                  <DropdownMenuItem>
                    <RotateCcw className="mr-2 h-4 w-4" /> Reverse Entry
                  </DropdownMenuItem>
                  <DropdownMenuItem className="text-destructive">
                    <XCircle className="mr-2 h-4 w-4" /> Void Entry
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8 text-center text-red-500">
        Failed to load journal entries
      </div>
    );
  }

  const draftCount = entries?.filter((e: JournalEntry) => e.status === "DRAFT").length || 0;
  const postedCount = entries?.filter((e: JournalEntry) => e.status === "POSTED").length || 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-1">
            <Link href="/gl" className="hover:text-foreground">General Ledger</Link>
            <ChevronRight className="h-4 w-4" />
            <span>Journal Entries</span>
          </div>
          <h1 className="text-3xl font-bold tracking-tight">Journal Entries</h1>
          <p className="text-muted-foreground">
            Create and manage journal entries
          </p>
        </div>
        <Link href="/gl/journal-entries/new">
          <Button>
            <Plus className="mr-2 h-4 w-4" /> New Entry
          </Button>
        </Link>
      </div>

      {/* Status Tabs */}
      <Tabs value={statusFilter} onValueChange={setStatusFilter}>
        <TabsList>
          <TabsTrigger value="all">
            All Entries
          </TabsTrigger>
          <TabsTrigger value="DRAFT" className="flex items-center gap-2">
            Draft
            {draftCount > 0 && (
              <Badge variant="secondary" className="ml-1">{draftCount}</Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="POSTED">Posted</TabsTrigger>
          <TabsTrigger value="VOIDED">Voided</TabsTrigger>
        </TabsList>

        <TabsContent value={statusFilter} className="mt-4">
          <DataTable 
            columns={columns} 
            data={entries || []} 
            searchKey="description"
            searchPlaceholder="Search entries..."
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
