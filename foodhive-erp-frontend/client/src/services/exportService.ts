export const exportService = {
  exportToCSV: (data: any[], filename: string, columns?: { key: string; label: string }[]) => {
    if (!data || !data.length) return;

    // Determine columns if not provided
    const headers = columns || Object.keys(data[0]).map(key => ({ key, label: key }));
    
    // Create CSV content
    const csvContent = [
      // Header row
      headers.map(h => `"${h.label}"`).join(','),
      // Data rows
      ...data.map(row => 
        headers.map(h => {
          const value = row[h.key];
          // Handle null/undefined
          if (value === null || value === undefined) return '""';
          // Handle objects/arrays
          if (typeof value === 'object') return `"${JSON.stringify(value).replace(/"/g, '""')}"`;
          // Handle strings with quotes
          return `"${String(value).replace(/"/g, '""')}"`;
        }).join(',')
      )
    ].join('\n');

    // Create download link
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `${filename}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  }
};
