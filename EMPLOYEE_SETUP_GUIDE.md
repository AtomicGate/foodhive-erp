# Employee Setup Guide

## Issues Found and Fixed

### 1. ✅ Backend Employee Create Function
**Problem:** The `Create` function was trying to insert into columns that don't exist in the simplified schema:
- `account_status` (should use `status`)
- `date_of_birth` (removed)
- `contract_id` (removed)
- `department_id` (removed)
- `warehouse_id` (removed)
- `created_by` (removed)

**Fixed:** Updated `registration/src/v1/services/employee/employee.go` to match the simplified schema:
```sql
INSERT INTO employees (
    email, password, english_name, arabic_name, 
    nationality, phone, role_id, status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
```

### 2. ⚠️ Roles Dropdown Empty
**Problem:** Roles dropdown in Employee form shows no options.

**Possible Causes:**
1. No roles exist in the database
2. Roles endpoint returns empty array
3. Frontend not parsing roles correctly

**Solution Steps:**

#### Step 1: Check if Roles Exist
Run this SQL in pgAdmin:
```sql
SELECT * FROM roles;
```

If empty, run the roles SQL file:
```sql
-- Run the file: sql/004_insert_roles.sql
-- Or copy/paste the INSERT statements from that file
```

Or create roles manually:
```sql
INSERT INTO roles (role_name, description, is_active) VALUES
    ('Super Admin', 'Full system access', true),
    ('Manager', 'Department manager', true),
    ('Employee', 'Standard employee', true),
    ('Sales Rep', 'Sales representative', true),
    ('Buyer', 'Purchasing department', true);
```

#### Step 2: Verify Roles Endpoint
Test the endpoint:
```powershell
# After logging in, get token first
$token = "YOUR_JWT_TOKEN"
Invoke-RestMethod -Uri "http://localhost:8080/v1/roles/list" -Headers @{Authorization="Bearer $token"}
```

Should return array of roles like:
```json
[
  {"id": 1, "role_name": "Super Admin", "description": "...", "is_active": true}
]
```

#### Step 3: Check Frontend Parsing
The frontend service at `foodhive-erp-frontend/client/src/services/masterDataService.ts` expects:
- Backend returns: `[{id, role_name, description, is_active}]`
- Frontend maps: `role.name` or `role.role_name`

**Fix:** Update `EmployeeList.tsx` line ~442 to use `role.name || role.role_name`

### 3. 📍 Where to Add Roles

**Option 1: Master Data Page (Recommended)**
Navigate to: `/admin/roles`
- This page allows creating, editing, and managing roles
- Located at: `foodhive-erp-frontend/client/src/pages/admin/MasterData.tsx`

**Option 2: Direct SQL**
Run the SQL insert statements above in pgAdmin

**Option 3: API Call**
```powershell
$body = @{
    role_name = "Manager"
    description = "Department manager"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/v1/roles/create" `
    -Method POST `
    -ContentType "application/json" `
    -Headers @{Authorization="Bearer $token"} `
    -Body $body
```

## Required Fields for Employee Creation

Based on the simplified schema and validation:

### Required:
- ✅ `email` - Must be valid email format
- ✅ `password` - Must be at least 8 characters
- ✅ `english_name` - Employee's English name
- ✅ `role_id` - Must be > 0 (must exist in roles table)

### Optional:
- `arabic_name` - Employee's Arabic name
- `phone` - Phone number
- `nationality` - Nationality
- `status` - Defaults to "ACTIVE" if not provided

## Testing Employee Creation

1. **Ensure roles exist:**
   ```sql
   SELECT id, role_name FROM roles WHERE is_active = true;
   ```

2. **Create employee via frontend:**
   - Go to Employees page
   - Click "Add Employee"
   - Fill required fields
   - Select a role from dropdown
   - Submit

3. **Verify creation:**
   ```sql
   SELECT * FROM employees ORDER BY id DESC LIMIT 1;
   ```

## Next Steps

1. ✅ Fixed backend Create function
2. ⏳ Add roles to database (use SQL or Master Data page)
3. ⏳ Test employee creation
4. ⏳ Verify roles dropdown populates correctly

## Master Data Management Pages

Access these pages to manage system data:

- **Roles:** `/admin/roles` - Create and manage user roles
- **Departments:** `/admin/departments` - Create and manage departments  
- **Warehouses:** `/admin/warehouses` - Create and manage warehouses

These pages are located at: `foodhive-erp-frontend/client/src/pages/admin/MasterData.tsx`
