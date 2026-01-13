package gl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound           = errors.New("GL entity not found")
	ErrAccountNotFound    = errors.New("GL account not found")
	ErrPeriodNotFound     = errors.New("GL period not found")
	ErrPeriodClosed       = errors.New("period is closed")
	ErrJournalNotFound    = errors.New("journal entry not found")
	ErrJournalPosted      = errors.New("journal entry already posted")
	ErrUnbalancedEntry    = errors.New("journal entry is unbalanced")
	ErrDuplicateCode      = errors.New("account code already exists")
	ErrAccountNotPostable = errors.New("account is not postable")
	ErrFiscalYearNotFound = errors.New("fiscal year not found")
)

type GLService interface {
	// Account Management (Chart of Accounts)
	CreateAccount(ctx context.Context, req models.CreateGLAccountRequest) (int, error)
	GetAccountByID(ctx context.Context, id int) (*models.GLAccount, error)
	GetAccountByCode(ctx context.Context, code string) (*models.GLAccount, error)
	UpdateAccount(ctx context.Context, id int, req models.UpdateGLAccountRequest) error
	DeleteAccount(ctx context.Context, id int) error
	ListAccounts(ctx context.Context, filters models.GLAccountListFilters) ([]models.GLAccount, int64, error)
	GetChartOfAccounts(ctx context.Context) ([]models.GLAccountWithChildren, error)

	// Fiscal Year & Period Management
	CreateFiscalYear(ctx context.Context, req models.CreateFiscalYearRequest) (int, error)
	GetFiscalYearByID(ctx context.Context, id int) (*models.GLFiscalYear, error)
	GetCurrentFiscalYear(ctx context.Context) (*models.GLFiscalYear, error)
	ListFiscalYears(ctx context.Context) ([]models.GLFiscalYear, error)
	CloseFiscalYear(ctx context.Context, id int, closedBy int) error
	GetPeriodByID(ctx context.Context, id int) (*models.GLPeriod, error)
	GetCurrentPeriod(ctx context.Context) (*models.GLPeriod, error)
	ListPeriods(ctx context.Context, fiscalYearID int) ([]models.GLPeriod, error)
	ClosePeriod(ctx context.Context, id int, closedBy int) error
	ReopenPeriod(ctx context.Context, id int) error

	// Journal Entry Management
	CreateJournalEntry(ctx context.Context, req models.CreateJournalEntryRequest, createdBy int) (int, error)
	GetJournalEntryByID(ctx context.Context, id int) (*models.GLJournalEntryWithLines, error)
	UpdateJournalEntry(ctx context.Context, id int, req models.CreateJournalEntryRequest) error
	DeleteJournalEntry(ctx context.Context, id int) error
	ListJournalEntries(ctx context.Context, filters models.JournalEntryListFilters) ([]models.GLJournalEntry, int64, error)
	PostJournalEntry(ctx context.Context, id int, postedBy int) error
	ReverseJournalEntry(ctx context.Context, id int, reversalDate string, createdBy int) (int, error)
	VoidJournalEntry(ctx context.Context, id int) error

	// Recurring Entries
	CreateRecurringEntry(ctx context.Context, req models.CreateRecurringEntryRequest, createdBy int) (int, error)
	GetRecurringEntryByID(ctx context.Context, id int) (*models.GLRecurringEntry, error)
	ListRecurringEntries(ctx context.Context) ([]models.GLRecurringEntry, error)
	ProcessRecurringEntries(ctx context.Context, asOfDate string, createdBy int) (int, error)

	// Reports
	GetTrialBalance(ctx context.Context, filters models.ReportFilters) (*models.TrialBalanceReport, error)
	GetIncomeStatement(ctx context.Context, filters models.ReportFilters) (*models.IncomeStatementReport, error)
	GetBalanceSheet(ctx context.Context, filters models.ReportFilters) (*models.BalanceSheetReport, error)
	GetAccountActivity(ctx context.Context, accountID int, dateFrom, dateTo string) (*models.AccountActivityReport, error)

	// Integration (for posting from other modules)
	PostFromAR(ctx context.Context, invoiceID int, createdBy int) (int, error)
	PostFromAP(ctx context.Context, invoiceID int, createdBy int) (int, error)
}

type glServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) GLService {
	return &glServiceImpl{db: db}
}

// ============================================
// Account Management
// ============================================

func (s *glServiceImpl) CreateAccount(ctx context.Context, req models.CreateGLAccountRequest) (int, error) {
	// Check for duplicate code
	var exists bool
	err := s.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM gl_accounts WHERE account_code = $1)
	`, req.AccountCode).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("checking duplicate: %w", err)
	}
	if exists {
		return 0, ErrDuplicateCode
	}

	var id int
	err = s.db.QueryRow(ctx, `
		INSERT INTO gl_accounts (
			account_code, account_name, account_type, account_sub_type,
			parent_id, description, currency, is_postable, is_bank_account,
			bank_account_id, normal_balance, opening_balance, current_balance,
			department_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12, $13)
		RETURNING id
	`,
		req.AccountCode, req.AccountName, req.AccountType, req.AccountSubType,
		req.ParentID, req.Description, req.Currency, req.IsPostable, req.IsBankAccount,
		req.BankAccountID, req.NormalBalance, req.OpeningBalance, req.DepartmentID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting account: %w", err)
	}

	return id, nil
}

func (s *glServiceImpl) GetAccountByID(ctx context.Context, id int) (*models.GLAccount, error) {
	account := &models.GLAccount{}
	var subType, description *string
	var parentID, bankAccountID, deptID *int

	err := s.db.QueryRow(ctx, `
		SELECT id, account_code, account_name, account_type, account_sub_type,
		       parent_id, description, currency, is_active, is_postable,
		       is_bank_account, bank_account_id, normal_balance, opening_balance,
		       current_balance, budget_amount, department_id, created_at, updated_at
		FROM gl_accounts WHERE id = $1
	`, id).Scan(
		&account.ID, &account.AccountCode, &account.AccountName, &account.AccountType, &subType,
		&parentID, &description, &account.Currency, &account.IsActive, &account.IsPostable,
		&account.IsBankAccount, &bankAccountID, &account.NormalBalance, &account.OpeningBalance,
		&account.CurrentBalance, &account.BudgetAmount, &deptID, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("getting account: %w", err)
	}

	if subType != nil {
		account.AccountSubType = models.GLAccountSubType(*subType)
	}
	if description != nil {
		account.Description = *description
	}
	if parentID != nil {
		account.ParentID = parentID
	}
	if bankAccountID != nil {
		account.BankAccountID = bankAccountID
	}
	if deptID != nil {
		account.DepartmentID = deptID
	}

	return account, nil
}

func (s *glServiceImpl) GetAccountByCode(ctx context.Context, code string) (*models.GLAccount, error) {
	var id int
	err := s.db.QueryRow(ctx, `SELECT id FROM gl_accounts WHERE account_code = $1`, code).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("finding account: %w", err)
	}
	return s.GetAccountByID(ctx, id)
}

func (s *glServiceImpl) UpdateAccount(ctx context.Context, id int, req models.UpdateGLAccountRequest) error {
	result, err := s.db.Exec(ctx, `
		UPDATE gl_accounts SET
			account_name = COALESCE($1, account_name),
			account_sub_type = COALESCE($2, account_sub_type),
			parent_id = COALESCE($3, parent_id),
			description = COALESCE($4, description),
			is_postable = COALESCE($5, is_postable),
			is_active = COALESCE($6, is_active),
			budget_amount = COALESCE($7, budget_amount),
			department_id = COALESCE($8, department_id),
			updated_at = NOW()
		WHERE id = $9
	`,
		req.AccountName, req.AccountSubType, req.ParentID, req.Description,
		req.IsPostable, req.IsActive, req.BudgetAmount, req.DepartmentID, id,
	)
	if err != nil {
		return fmt.Errorf("updating account: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrAccountNotFound
	}
	return nil
}

func (s *glServiceImpl) DeleteAccount(ctx context.Context, id int) error {
	// Check if account has transactions
	var hasTransactions bool
	err := s.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM gl_journal_lines WHERE account_id = $1)
	`, id).Scan(&hasTransactions)
	if err != nil {
		return fmt.Errorf("checking transactions: %w", err)
	}
	if hasTransactions {
		// Soft delete - just deactivate
		_, err = s.db.Exec(ctx, `UPDATE gl_accounts SET is_active = false WHERE id = $1`, id)
		return err
	}

	result, err := s.db.Exec(ctx, `DELETE FROM gl_accounts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting account: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrAccountNotFound
	}
	return nil
}

func (s *glServiceImpl) ListAccounts(ctx context.Context, filters models.GLAccountListFilters) ([]models.GLAccount, int64, error) {
	query := `
		SELECT id, account_code, account_name, account_type, account_sub_type,
		       parent_id, description, currency, is_active, is_postable,
		       is_bank_account, bank_account_id, normal_balance, opening_balance,
		       current_balance, budget_amount, department_id, created_at, updated_at
		FROM gl_accounts WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM gl_accounts WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if filters.AccountType != nil {
		query += fmt.Sprintf(" AND account_type = $%d", argNum)
		countQuery += fmt.Sprintf(" AND account_type = $%d", argNum)
		args = append(args, *filters.AccountType)
		argNum++
	}

	if filters.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argNum)
		countQuery += fmt.Sprintf(" AND is_active = $%d", argNum)
		args = append(args, *filters.IsActive)
		argNum++
	}

	if filters.IsPostable != nil {
		query += fmt.Sprintf(" AND is_postable = $%d", argNum)
		countQuery += fmt.Sprintf(" AND is_postable = $%d", argNum)
		args = append(args, *filters.IsPostable)
		argNum++
	}

	if filters.ParentID != nil {
		query += fmt.Sprintf(" AND parent_id = $%d", argNum)
		countQuery += fmt.Sprintf(" AND parent_id = $%d", argNum)
		args = append(args, *filters.ParentID)
		argNum++
	}

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (account_code ILIKE $%d OR account_name ILIKE $%d)", argNum, argNum)
		countQuery += fmt.Sprintf(" AND (account_code ILIKE $%d OR account_name ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filters.Search+"%")
		argNum++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting accounts: %w", err)
	}

	if filters.PageSize == 0 {
		filters.PageSize = 50
	}
	if filters.Page == 0 {
		filters.Page = 1
	}
	offset := (filters.Page - 1) * filters.PageSize

	query += fmt.Sprintf(" ORDER BY account_code LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var accounts []models.GLAccount
	for rows.Next() {
		var a models.GLAccount
		var subType, description *string
		var parentID, bankAccountID, deptID *int

		err := rows.Scan(
			&a.ID, &a.AccountCode, &a.AccountName, &a.AccountType, &subType,
			&parentID, &description, &a.Currency, &a.IsActive, &a.IsPostable,
			&a.IsBankAccount, &bankAccountID, &a.NormalBalance, &a.OpeningBalance,
			&a.CurrentBalance, &a.BudgetAmount, &deptID, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning account: %w", err)
		}

		if subType != nil {
			a.AccountSubType = models.GLAccountSubType(*subType)
		}
		if description != nil {
			a.Description = *description
		}
		if parentID != nil {
			a.ParentID = parentID
		}
		if bankAccountID != nil {
			a.BankAccountID = bankAccountID
		}
		if deptID != nil {
			a.DepartmentID = deptID
		}

		accounts = append(accounts, a)
	}

	return accounts, total, nil
}

func (s *glServiceImpl) GetChartOfAccounts(ctx context.Context) ([]models.GLAccountWithChildren, error) {
	// Get all active accounts ordered by code
	rows := s.db.Query(ctx, `
		SELECT id, account_code, account_name, account_type, account_sub_type,
		       parent_id, description, currency, is_active, is_postable,
		       is_bank_account, bank_account_id, normal_balance, opening_balance,
		       current_balance, budget_amount, department_id, created_at, updated_at
		FROM gl_accounts
		WHERE is_active = true
		ORDER BY account_code
	`)
	defer rows.Close()

	accountMap := make(map[int]*models.GLAccountWithChildren)
	var rootAccounts []*models.GLAccountWithChildren

	for rows.Next() {
		var a models.GLAccount
		var subType, description *string
		var parentID, bankAccountID, deptID *int

		err := rows.Scan(
			&a.ID, &a.AccountCode, &a.AccountName, &a.AccountType, &subType,
			&parentID, &description, &a.Currency, &a.IsActive, &a.IsPostable,
			&a.IsBankAccount, &bankAccountID, &a.NormalBalance, &a.OpeningBalance,
			&a.CurrentBalance, &a.BudgetAmount, &deptID, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning account: %w", err)
		}

		if subType != nil {
			a.AccountSubType = models.GLAccountSubType(*subType)
		}
		if description != nil {
			a.Description = *description
		}
		if parentID != nil {
			a.ParentID = parentID
		}
		if bankAccountID != nil {
			a.BankAccountID = bankAccountID
		}
		if deptID != nil {
			a.DepartmentID = deptID
		}

		node := &models.GLAccountWithChildren{Account: a}
		accountMap[a.ID] = node

		if parentID == nil {
			rootAccounts = append(rootAccounts, node)
		}
	}

	// Build hierarchy
	for _, node := range accountMap {
		if node.Account.ParentID != nil {
			if parent, ok := accountMap[*node.Account.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	// Calculate levels
	var setLevels func(nodes []models.GLAccountWithChildren, level int) []models.GLAccountWithChildren
	setLevels = func(nodes []models.GLAccountWithChildren, level int) []models.GLAccountWithChildren {
		for i := range nodes {
			nodes[i].Level = level
			if len(nodes[i].Children) > 0 {
				nodes[i].Children = setLevels(nodes[i].Children, level+1)
			}
		}
		return nodes
	}

	result := make([]models.GLAccountWithChildren, len(rootAccounts))
	for i, r := range rootAccounts {
		result[i] = *r
	}
	return setLevels(result, 0), nil
}

// ============================================
// Fiscal Year & Period Management
// ============================================

func (s *glServiceImpl) CreateFiscalYear(ctx context.Context, req models.CreateFiscalYearRequest) (int, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return 0, fmt.Errorf("parsing start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return 0, fmt.Errorf("parsing end date: %w", err)
	}

	var fiscalYearID int
	err = s.db.QueryRow(ctx, `
		INSERT INTO gl_fiscal_years (year_code, start_date, end_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.YearCode, startDate, endDate).Scan(&fiscalYearID)
	if err != nil {
		return 0, fmt.Errorf("creating fiscal year: %w", err)
	}

	// Create 12 periods + 1 adjustment period
	current := startDate
	for i := 1; i <= 13; i++ {
		var periodEnd time.Time
		var periodName string
		var isAdjustment bool

		if i == 13 {
			// Adjustment period
			periodName = "Period 13 - Adjustments"
			current = endDate
			periodEnd = endDate
			isAdjustment = true
		} else {
			periodName = current.Format("January 2006")
			// Move to end of month
			periodEnd = current.AddDate(0, 1, -1)
			if periodEnd.After(endDate) {
				periodEnd = endDate
			}
		}

		_, err = s.db.Exec(ctx, `
			INSERT INTO gl_periods (fiscal_year_id, period_number, period_name, start_date, end_date, is_adjustment)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, fiscalYearID, i, periodName, current, periodEnd, isAdjustment)
		if err != nil {
			return 0, fmt.Errorf("creating period %d: %w", i, err)
		}

		if i < 13 {
			current = periodEnd.AddDate(0, 0, 1)
		}
	}

	return fiscalYearID, nil
}

func (s *glServiceImpl) GetFiscalYearByID(ctx context.Context, id int) (*models.GLFiscalYear, error) {
	fy := &models.GLFiscalYear{}
	var closedBy *int
	var closedAt *time.Time

	err := s.db.QueryRow(ctx, `
		SELECT id, year_code, start_date, end_date, is_current, is_closed, closed_by, closed_at, created_at
		FROM gl_fiscal_years WHERE id = $1
	`, id).Scan(
		&fy.ID, &fy.YearCode, &fy.StartDate, &fy.EndDate, &fy.IsCurrent, &fy.IsClosed,
		&closedBy, &closedAt, &fy.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFiscalYearNotFound
		}
		return nil, fmt.Errorf("getting fiscal year: %w", err)
	}

	if closedBy != nil {
		fy.ClosedBy = closedBy
	}

	return fy, nil
}

func (s *glServiceImpl) GetCurrentFiscalYear(ctx context.Context) (*models.GLFiscalYear, error) {
	var id int
	err := s.db.QueryRow(ctx, `SELECT id FROM gl_fiscal_years WHERE is_current = true LIMIT 1`).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFiscalYearNotFound
		}
		return nil, fmt.Errorf("finding current fiscal year: %w", err)
	}
	return s.GetFiscalYearByID(ctx, id)
}

func (s *glServiceImpl) ListFiscalYears(ctx context.Context) ([]models.GLFiscalYear, error) {
	rows := s.db.Query(ctx, `
		SELECT id, year_code, start_date, end_date, is_current, is_closed, closed_by, closed_at, created_at
		FROM gl_fiscal_years ORDER BY start_date DESC
	`)
	defer rows.Close()

	var years []models.GLFiscalYear
	for rows.Next() {
		var fy models.GLFiscalYear
		var closedBy *int
		var closedAt *time.Time

		err := rows.Scan(
			&fy.ID, &fy.YearCode, &fy.StartDate, &fy.EndDate, &fy.IsCurrent, &fy.IsClosed,
			&closedBy, &closedAt, &fy.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning fiscal year: %w", err)
		}
		if closedBy != nil {
			fy.ClosedBy = closedBy
		}
		years = append(years, fy)
	}

	return years, nil
}

func (s *glServiceImpl) CloseFiscalYear(ctx context.Context, id int, closedBy int) error {
	// Check all periods are closed
	var openPeriods int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM gl_periods WHERE fiscal_year_id = $1 AND status = 'OPEN'
	`, id).Scan(&openPeriods)
	if err != nil {
		return fmt.Errorf("checking open periods: %w", err)
	}
	if openPeriods > 0 {
		return fmt.Errorf("cannot close fiscal year with %d open periods", openPeriods)
	}

	_, err = s.db.Exec(ctx, `
		UPDATE gl_fiscal_years SET is_closed = true, closed_by = $1, closed_at = NOW() WHERE id = $2
	`, closedBy, id)
	return err
}

func (s *glServiceImpl) GetPeriodByID(ctx context.Context, id int) (*models.GLPeriod, error) {
	p := &models.GLPeriod{}
	var closedBy *int
	var closedAt *time.Time

	err := s.db.QueryRow(ctx, `
		SELECT id, fiscal_year_id, period_number, period_name, start_date, end_date,
		       status, is_adjustment, closed_by, closed_at
		FROM gl_periods WHERE id = $1
	`, id).Scan(
		&p.ID, &p.FiscalYearID, &p.PeriodNumber, &p.PeriodName, &p.StartDate, &p.EndDate,
		&p.Status, &p.IsAdjustment, &closedBy, &closedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPeriodNotFound
		}
		return nil, fmt.Errorf("getting period: %w", err)
	}

	if closedBy != nil {
		p.ClosedBy = closedBy
	}

	return p, nil
}

func (s *glServiceImpl) GetCurrentPeriod(ctx context.Context) (*models.GLPeriod, error) {
	var id int
	now := time.Now()
	err := s.db.QueryRow(ctx, `
		SELECT p.id FROM gl_periods p
		JOIN gl_fiscal_years fy ON p.fiscal_year_id = fy.id
		WHERE fy.is_current = true AND p.status = 'OPEN' AND p.is_adjustment = false
		  AND $1 BETWEEN p.start_date AND p.end_date
		LIMIT 1
	`, now).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPeriodNotFound
		}
		return nil, fmt.Errorf("finding current period: %w", err)
	}
	return s.GetPeriodByID(ctx, id)
}

func (s *glServiceImpl) ListPeriods(ctx context.Context, fiscalYearID int) ([]models.GLPeriod, error) {
	rows := s.db.Query(ctx, `
		SELECT id, fiscal_year_id, period_number, period_name, start_date, end_date,
		       status, is_adjustment, closed_by, closed_at
		FROM gl_periods WHERE fiscal_year_id = $1 ORDER BY period_number
	`, fiscalYearID)
	defer rows.Close()

	var periods []models.GLPeriod
	for rows.Next() {
		var p models.GLPeriod
		var closedBy *int
		var closedAt *time.Time

		err := rows.Scan(
			&p.ID, &p.FiscalYearID, &p.PeriodNumber, &p.PeriodName, &p.StartDate, &p.EndDate,
			&p.Status, &p.IsAdjustment, &closedBy, &closedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning period: %w", err)
		}
		if closedBy != nil {
			p.ClosedBy = closedBy
		}
		periods = append(periods, p)
	}

	return periods, nil
}

func (s *glServiceImpl) ClosePeriod(ctx context.Context, id int, closedBy int) error {
	_, err := s.db.Exec(ctx, `
		UPDATE gl_periods SET status = 'CLOSED', closed_by = $1, closed_at = NOW() WHERE id = $2
	`, closedBy, id)
	return err
}

func (s *glServiceImpl) ReopenPeriod(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, `
		UPDATE gl_periods SET status = 'OPEN', closed_by = NULL, closed_at = NULL WHERE id = $1
	`, id)
	return err
}

// ============================================
// Journal Entry Management
// ============================================

func (s *glServiceImpl) CreateJournalEntry(ctx context.Context, req models.CreateJournalEntryRequest, createdBy int) (int, error) {
	entryDate, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		return 0, fmt.Errorf("parsing entry date: %w", err)
	}

	// Find the period for this date
	var periodID int
	err = s.db.QueryRow(ctx, `
		SELECT p.id FROM gl_periods p
		JOIN gl_fiscal_years fy ON p.fiscal_year_id = fy.id
		WHERE $1 BETWEEN p.start_date AND p.end_date AND p.status = 'OPEN' AND p.is_adjustment = false
		LIMIT 1
	`, entryDate).Scan(&periodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrPeriodClosed
		}
		return 0, fmt.Errorf("finding period: %w", err)
	}

	// Calculate totals
	var totalDebit, totalCredit float64
	for _, line := range req.Lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}

	// Generate journal number
	var journalNum string
	err = s.db.QueryRow(ctx, `
		SELECT 'JE-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD(COALESCE(MAX(id), 0)::TEXT, 4, '0')
		FROM gl_journal_entries
	`).Scan(&journalNum)
	if err != nil {
		journalNum = fmt.Sprintf("JE-%s-0001", time.Now().Format("20060102"))
	}

	if req.ExchangeRate == 0 {
		req.ExchangeRate = 1.0
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Insert journal entry
	var entryID int
	err = s.db.QueryRow(ctx, `
		INSERT INTO gl_journal_entries (
			journal_number, entry_date, posting_date, period_id, entry_type, status,
			description, reference, total_debit, total_credit, currency, exchange_rate,
			auto_reverse, auto_reverse_date, created_by
		) VALUES ($1, $2, $2, $3, $4, 'DRAFT', $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`,
		journalNum, entryDate, periodID, req.EntryType, req.Description, req.Reference,
		totalDebit, totalCredit, req.Currency, req.ExchangeRate,
		req.AutoReverse, req.AutoReverseDate, createdBy,
	).Scan(&entryID)
	if err != nil {
		return 0, fmt.Errorf("inserting journal entry: %w", err)
	}

	// Insert lines
	for i, line := range req.Lines {
		// Verify account is postable
		var isPostable bool
		err = s.db.QueryRow(ctx, `SELECT is_postable FROM gl_accounts WHERE id = $1`, line.AccountID).Scan(&isPostable)
		if err != nil {
			return 0, fmt.Errorf("checking account %d: %w", line.AccountID, err)
		}
		if !isPostable {
			return 0, ErrAccountNotPostable
		}

		_, err = s.db.Exec(ctx, `
			INSERT INTO gl_journal_lines (
				journal_id, line_number, account_id, description, debit_amount, credit_amount,
				department_id, project_id, reference
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			entryID, i+1, line.AccountID, line.Description, line.DebitAmount, line.CreditAmount,
			line.DepartmentID, line.ProjectID, line.Reference,
		)
		if err != nil {
			return 0, fmt.Errorf("inserting line %d: %w", i+1, err)
		}
	}

	return entryID, nil
}

func (s *glServiceImpl) GetJournalEntryByID(ctx context.Context, id int) (*models.GLJournalEntryWithLines, error) {
	entry := &models.GLJournalEntryWithLines{}
	var ref, srcDoc, srcMod *string
	var srcID, recID, revID, postedBy *int
	var autoRevDate *time.Time
	var postedAt *time.Time

	err := s.db.QueryRow(ctx, `
		SELECT je.id, je.journal_number, je.entry_date, je.posting_date, je.period_id, je.entry_type,
		       je.status, je.description, je.reference, je.source_document, je.source_module, je.source_id,
		       je.total_debit, je.total_credit, je.currency, je.exchange_rate, je.is_recurring, je.recurring_id,
		       je.reversed_entry_id, je.auto_reverse, je.auto_reverse_date, je.created_by, je.posted_by,
		       je.posted_at, je.created_at, p.period_name, e.full_name
		FROM gl_journal_entries je
		JOIN gl_periods p ON je.period_id = p.id
		JOIN employees e ON je.created_by = e.id
		WHERE je.id = $1
	`, id).Scan(
		&entry.Entry.ID, &entry.Entry.JournalNumber, &entry.Entry.EntryDate, &entry.Entry.PostingDate,
		&entry.Entry.PeriodID, &entry.Entry.EntryType, &entry.Entry.Status, &entry.Entry.Description,
		&ref, &srcDoc, &srcMod, &srcID, &entry.Entry.TotalDebit, &entry.Entry.TotalCredit,
		&entry.Entry.Currency, &entry.Entry.ExchangeRate, &entry.Entry.IsRecurring, &recID,
		&revID, &entry.Entry.AutoReverse, &autoRevDate, &entry.Entry.CreatedBy, &postedBy,
		&postedAt, &entry.Entry.CreatedAt, &entry.PeriodName, &entry.CreatorName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJournalNotFound
		}
		return nil, fmt.Errorf("getting journal entry: %w", err)
	}

	if ref != nil {
		entry.Entry.Reference = *ref
	}
	if srcDoc != nil {
		entry.Entry.SourceDocument = *srcDoc
	}
	if srcMod != nil {
		entry.Entry.SourceModule = *srcMod
	}
	if srcID != nil {
		entry.Entry.SourceID = srcID
	}
	if recID != nil {
		entry.Entry.RecurringID = recID
	}
	if revID != nil {
		entry.Entry.ReversedEntryID = revID
	}
	if postedBy != nil {
		entry.Entry.PostedBy = postedBy
	}

	// Get lines
	rows := s.db.Query(ctx, `
		SELECT id, journal_id, line_number, account_id, description, debit_amount, credit_amount,
		       department_id, project_id, reference, created_at
		FROM gl_journal_lines WHERE journal_id = $1 ORDER BY line_number
	`, id)
	defer rows.Close()

	for rows.Next() {
		var line models.GLJournalLine
		var desc, lineRef *string
		var deptID, projID *int

		err := rows.Scan(
			&line.ID, &line.JournalID, &line.LineNumber, &line.AccountID, &desc,
			&line.DebitAmount, &line.CreditAmount, &deptID, &projID, &lineRef, &line.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning line: %w", err)
		}

		if desc != nil {
			line.Description = *desc
		}
		if lineRef != nil {
			line.Reference = *lineRef
		}
		if deptID != nil {
			line.DepartmentID = deptID
		}
		if projID != nil {
			line.ProjectID = projID
		}

		entry.Lines = append(entry.Lines, line)
	}

	return entry, nil
}

func (s *glServiceImpl) UpdateJournalEntry(ctx context.Context, id int, req models.CreateJournalEntryRequest) error {
	// Check if entry is still draft
	var status models.JournalEntryStatus
	err := s.db.QueryRow(ctx, `SELECT status FROM gl_journal_entries WHERE id = $1`, id).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJournalNotFound
		}
		return fmt.Errorf("checking status: %w", err)
	}
	if status != models.JournalStatusDraft {
		return ErrJournalPosted
	}

	// Delete existing lines
	_, err = s.db.Exec(ctx, `DELETE FROM gl_journal_lines WHERE journal_id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting lines: %w", err)
	}

	// Calculate totals
	var totalDebit, totalCredit float64
	for _, line := range req.Lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}

	// Update entry
	_, err = s.db.Exec(ctx, `
		UPDATE gl_journal_entries SET
			entry_date = $1, description = $2, reference = $3,
			total_debit = $4, total_credit = $5, auto_reverse = $6, auto_reverse_date = $7
		WHERE id = $8
	`, req.EntryDate, req.Description, req.Reference, totalDebit, totalCredit,
		req.AutoReverse, req.AutoReverseDate, id)
	if err != nil {
		return fmt.Errorf("updating entry: %w", err)
	}

	// Insert new lines
	for i, line := range req.Lines {
		_, err = s.db.Exec(ctx, `
			INSERT INTO gl_journal_lines (
				journal_id, line_number, account_id, description, debit_amount, credit_amount,
				department_id, project_id, reference
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			id, i+1, line.AccountID, line.Description, line.DebitAmount, line.CreditAmount,
			line.DepartmentID, line.ProjectID, line.Reference,
		)
		if err != nil {
			return fmt.Errorf("inserting line %d: %w", i+1, err)
		}
	}

	return nil
}

func (s *glServiceImpl) DeleteJournalEntry(ctx context.Context, id int) error {
	var status models.JournalEntryStatus
	err := s.db.QueryRow(ctx, `SELECT status FROM gl_journal_entries WHERE id = $1`, id).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJournalNotFound
		}
		return fmt.Errorf("checking status: %w", err)
	}
	if status != models.JournalStatusDraft {
		return ErrJournalPosted
	}

	_, err = s.db.Exec(ctx, `DELETE FROM gl_journal_lines WHERE journal_id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting lines: %w", err)
	}
	_, err = s.db.Exec(ctx, `DELETE FROM gl_journal_entries WHERE id = $1`, id)
	return err
}

func (s *glServiceImpl) ListJournalEntries(ctx context.Context, filters models.JournalEntryListFilters) ([]models.GLJournalEntry, int64, error) {
	query := `
		SELECT id, journal_number, entry_date, posting_date, period_id, entry_type, status,
		       description, reference, source_document, source_module, source_id,
		       total_debit, total_credit, currency, exchange_rate, is_recurring, recurring_id,
		       reversed_entry_id, auto_reverse, auto_reverse_date, created_by, posted_by, posted_at, created_at
		FROM gl_journal_entries WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM gl_journal_entries WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if filters.PeriodID != nil {
		query += fmt.Sprintf(" AND period_id = $%d", argNum)
		countQuery += fmt.Sprintf(" AND period_id = $%d", argNum)
		args = append(args, *filters.PeriodID)
		argNum++
	}

	if filters.EntryType != nil {
		query += fmt.Sprintf(" AND entry_type = $%d", argNum)
		countQuery += fmt.Sprintf(" AND entry_type = $%d", argNum)
		args = append(args, *filters.EntryType)
		argNum++
	}

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		countQuery += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, *filters.Status)
		argNum++
	}

	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND entry_date >= $%d", argNum)
		countQuery += fmt.Sprintf(" AND entry_date >= $%d", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}

	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND entry_date <= $%d", argNum)
		countQuery += fmt.Sprintf(" AND entry_date <= $%d", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting entries: %w", err)
	}

	if filters.PageSize == 0 {
		filters.PageSize = 50
	}
	if filters.Page == 0 {
		filters.Page = 1
	}
	offset := (filters.Page - 1) * filters.PageSize

	query += fmt.Sprintf(" ORDER BY entry_date DESC, id DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var entries []models.GLJournalEntry
	for rows.Next() {
		var e models.GLJournalEntry
		var ref, srcDoc, srcMod *string
		var srcID, recID, revID, postedBy *int
		var autoRevDate, postedAt *time.Time

		err := rows.Scan(
			&e.ID, &e.JournalNumber, &e.EntryDate, &e.PostingDate, &e.PeriodID, &e.EntryType, &e.Status,
			&e.Description, &ref, &srcDoc, &srcMod, &srcID, &e.TotalDebit, &e.TotalCredit,
			&e.Currency, &e.ExchangeRate, &e.IsRecurring, &recID, &revID, &e.AutoReverse,
			&autoRevDate, &e.CreatedBy, &postedBy, &postedAt, &e.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning entry: %w", err)
		}

		if ref != nil {
			e.Reference = *ref
		}
		if srcDoc != nil {
			e.SourceDocument = *srcDoc
		}
		if srcMod != nil {
			e.SourceModule = *srcMod
		}
		if srcID != nil {
			e.SourceID = srcID
		}
		if recID != nil {
			e.RecurringID = recID
		}
		if revID != nil {
			e.ReversedEntryID = revID
		}
		if postedBy != nil {
			e.PostedBy = postedBy
		}

		entries = append(entries, e)
	}

	return entries, total, nil
}

func (s *glServiceImpl) PostJournalEntry(ctx context.Context, id int, postedBy int) error {
	// Get entry
	var status models.JournalEntryStatus
	var totalDebit, totalCredit float64
	err := s.db.QueryRow(ctx, `
		SELECT status, total_debit, total_credit FROM gl_journal_entries WHERE id = $1
	`, id).Scan(&status, &totalDebit, &totalCredit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJournalNotFound
		}
		return fmt.Errorf("getting entry: %w", err)
	}

	if status != models.JournalStatusDraft && status != models.JournalStatusPending {
		return ErrJournalPosted
	}

	// Verify balanced
	if totalDebit != totalCredit {
		return ErrUnbalancedEntry
	}

	// Update account balances
	rows := s.db.Query(ctx, `
		SELECT account_id, SUM(debit_amount) as debit, SUM(credit_amount) as credit
		FROM gl_journal_lines WHERE journal_id = $1
		GROUP BY account_id
	`, id)
	defer rows.Close()

	for rows.Next() {
		var accountID int
		var debit, credit float64
		if err := rows.Scan(&accountID, &debit, &credit); err != nil {
			return fmt.Errorf("scanning line: %w", err)
		}

		// Get account normal balance
		var normalBalance string
		err = s.db.QueryRow(ctx, `SELECT normal_balance FROM gl_accounts WHERE id = $1`, accountID).Scan(&normalBalance)
		if err != nil {
			return fmt.Errorf("getting account: %w", err)
		}

		var balanceChange float64
		if normalBalance == "DEBIT" {
			balanceChange = debit - credit
		} else {
			balanceChange = credit - debit
		}

		_, err = s.db.Exec(ctx, `
			UPDATE gl_accounts SET current_balance = current_balance + $1 WHERE id = $2
		`, balanceChange, accountID)
		if err != nil {
			return fmt.Errorf("updating account balance: %w", err)
		}
	}

	// Update entry status
	_, err = s.db.Exec(ctx, `
		UPDATE gl_journal_entries SET status = 'POSTED', posting_date = NOW(), posted_by = $1, posted_at = NOW()
		WHERE id = $2
	`, postedBy, id)
	return err
}

func (s *glServiceImpl) ReverseJournalEntry(ctx context.Context, id int, reversalDate string, createdBy int) (int, error) {
	entry, err := s.GetJournalEntryByID(ctx, id)
	if err != nil {
		return 0, err
	}

	if entry.Entry.Status != models.JournalStatusPosted {
		return 0, fmt.Errorf("can only reverse posted entries")
	}

	// Create reversal entry
	req := models.CreateJournalEntryRequest{
		EntryDate:    reversalDate,
		EntryType:    entry.Entry.EntryType,
		Description:  "Reversal of " + entry.Entry.JournalNumber,
		Reference:    entry.Entry.Reference,
		Currency:     entry.Entry.Currency,
		ExchangeRate: entry.Entry.ExchangeRate,
	}

	// Swap debits and credits
	for _, line := range entry.Lines {
		req.Lines = append(req.Lines, models.CreateJournalLineRequest{
			AccountID:    line.AccountID,
			Description:  "Reversal: " + line.Description,
			DebitAmount:  line.CreditAmount, // Swapped
			CreditAmount: line.DebitAmount,  // Swapped
			DepartmentID: line.DepartmentID,
			ProjectID:    line.ProjectID,
		})
	}

	reversalID, err := s.CreateJournalEntry(ctx, req, createdBy)
	if err != nil {
		return 0, err
	}

	// Link the entries
	_, err = s.db.Exec(ctx, `UPDATE gl_journal_entries SET reversed_entry_id = $1 WHERE id = $2`, reversalID, id)
	if err != nil {
		return 0, err
	}

	_, err = s.db.Exec(ctx, `UPDATE gl_journal_entries SET status = 'REVERSED' WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	return reversalID, nil
}

func (s *glServiceImpl) VoidJournalEntry(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, `UPDATE gl_journal_entries SET status = 'VOIDED' WHERE id = $1 AND status = 'DRAFT'`, id)
	return err
}

// ============================================
// Recurring Entries (stubs - implement as needed)
// ============================================

func (s *glServiceImpl) CreateRecurringEntry(ctx context.Context, req models.CreateRecurringEntryRequest, createdBy int) (int, error) {
	// TODO: Implement recurring entry creation
	return 0, fmt.Errorf("not implemented")
}

func (s *glServiceImpl) GetRecurringEntryByID(ctx context.Context, id int) (*models.GLRecurringEntry, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *glServiceImpl) ListRecurringEntries(ctx context.Context) ([]models.GLRecurringEntry, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *glServiceImpl) ProcessRecurringEntries(ctx context.Context, asOfDate string, createdBy int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// ============================================
// Reports
// ============================================

func (s *glServiceImpl) GetTrialBalance(ctx context.Context, filters models.ReportFilters) (*models.TrialBalanceReport, error) {
	report := &models.TrialBalanceReport{}

	query := `
		SELECT a.id, a.account_code, a.account_name, a.account_type,
		       a.opening_balance,
		       COALESCE(SUM(jl.debit_amount), 0) as period_debit,
		       COALESCE(SUM(jl.credit_amount), 0) as period_credit
		FROM gl_accounts a
		LEFT JOIN gl_journal_lines jl ON a.id = jl.account_id
		LEFT JOIN gl_journal_entries je ON jl.journal_id = je.id AND je.status = 'POSTED'
		WHERE a.is_active = true
	`
	args := []interface{}{}
	argNum := 1

	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND (je.posting_date >= $%d OR je.id IS NULL)", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}

	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND (je.posting_date <= $%d OR je.id IS NULL)", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	query += " GROUP BY a.id, a.account_code, a.account_name, a.account_type, a.opening_balance ORDER BY a.account_code"

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	for rows.Next() {
		var row models.TrialBalanceRow
		var openingBalance float64

		err := rows.Scan(
			&row.AccountID, &row.AccountCode, &row.AccountName, &row.AccountType,
			&openingBalance, &row.PeriodDebit, &row.PeriodCredit,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		// Set opening balances based on account type
		if row.AccountType == models.GLAccountTypeAsset || row.AccountType == models.GLAccountTypeExpense {
			row.OpeningDebit = openingBalance
		} else {
			row.OpeningCredit = openingBalance
		}

		// Calculate closing
		row.ClosingDebit = row.OpeningDebit + row.PeriodDebit - row.PeriodCredit
		row.ClosingCredit = row.OpeningCredit + row.PeriodCredit - row.PeriodDebit

		if row.ClosingDebit < 0 {
			row.ClosingCredit = -row.ClosingDebit
			row.ClosingDebit = 0
		}
		if row.ClosingCredit < 0 {
			row.ClosingDebit = -row.ClosingCredit
			row.ClosingCredit = 0
		}

		report.Rows = append(report.Rows, row)
		report.TotalDebit += row.ClosingDebit
		report.TotalCredit += row.ClosingCredit
	}

	return report, nil
}

func (s *glServiceImpl) GetIncomeStatement(ctx context.Context, filters models.ReportFilters) (*models.IncomeStatementReport, error) {
	report := &models.IncomeStatementReport{}

	// Get revenue and expense accounts with their activity
	query := `
		SELECT a.id, a.account_code, a.account_name, a.account_type,
		       COALESCE(SUM(jl.credit_amount) - SUM(jl.debit_amount), 0) as net_amount
		FROM gl_accounts a
		LEFT JOIN gl_journal_lines jl ON a.id = jl.account_id
		LEFT JOIN gl_journal_entries je ON jl.journal_id = je.id AND je.status = 'POSTED'
		WHERE a.is_active = true AND a.account_type IN ('REVENUE', 'EXPENSE')
	`
	args := []interface{}{}
	argNum := 1

	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND (je.posting_date >= $%d OR je.id IS NULL)", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}

	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND (je.posting_date <= $%d OR je.id IS NULL)", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	query += " GROUP BY a.id, a.account_code, a.account_name, a.account_type ORDER BY a.account_type DESC, a.account_code"

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	for rows.Next() {
		var row models.IncomeStatementRow
		err := rows.Scan(&row.AccountID, &row.AccountCode, &row.AccountName, &row.AccountType, &row.Amount)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		if row.AccountType == models.GLAccountTypeRevenue {
			report.Revenue = append(report.Revenue, row)
			report.TotalRevenue += row.Amount
		} else {
			row.Amount = -row.Amount // Expenses are shown as positive
			report.Expenses = append(report.Expenses, row)
			report.TotalExpenses += row.Amount
		}
	}

	report.NetIncome = report.TotalRevenue - report.TotalExpenses

	return report, nil
}

func (s *glServiceImpl) GetBalanceSheet(ctx context.Context, filters models.ReportFilters) (*models.BalanceSheetReport, error) {
	report := &models.BalanceSheetReport{}

	query := `
		SELECT a.id, a.account_code, a.account_name, a.account_type, a.current_balance
		FROM gl_accounts a
		WHERE a.is_active = true AND a.account_type IN ('ASSET', 'LIABILITY', 'EQUITY')
		ORDER BY a.account_type, a.account_code
	`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	for rows.Next() {
		var row models.BalanceSheetRow
		err := rows.Scan(&row.AccountID, &row.AccountCode, &row.AccountName, &row.AccountType, &row.Balance)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		switch row.AccountType {
		case models.GLAccountTypeAsset:
			report.Assets = append(report.Assets, row)
			report.TotalAssets += row.Balance
		case models.GLAccountTypeLiability:
			report.Liabilities = append(report.Liabilities, row)
			report.TotalLiabilities += row.Balance
		case models.GLAccountTypeEquity:
			report.Equity = append(report.Equity, row)
			report.TotalEquity += row.Balance
		}
	}

	return report, nil
}

func (s *glServiceImpl) GetAccountActivity(ctx context.Context, accountID int, dateFrom, dateTo string) (*models.AccountActivityReport, error) {
	report := &models.AccountActivityReport{AccountID: accountID}

	// Get account info
	account, err := s.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	report.AccountCode = account.AccountCode
	report.AccountName = account.AccountName

	// Get opening balance (sum of all transactions before dateFrom)
	err = s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(jl.debit_amount) - SUM(jl.credit_amount), 0)
		FROM gl_journal_lines jl
		JOIN gl_journal_entries je ON jl.journal_id = je.id
		WHERE jl.account_id = $1 AND je.status = 'POSTED' AND je.posting_date < $2
	`, accountID, dateFrom).Scan(&report.OpeningBalance)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("getting opening balance: %w", err)
	}
	report.OpeningBalance += account.OpeningBalance

	// Get activity
	rows := s.db.Query(ctx, `
		SELECT je.posting_date, je.journal_number, je.description, je.reference,
		       jl.debit_amount, jl.credit_amount, je.source_module
		FROM gl_journal_lines jl
		JOIN gl_journal_entries je ON jl.journal_id = je.id
		WHERE jl.account_id = $1 AND je.status = 'POSTED'
		  AND je.posting_date >= $2 AND je.posting_date <= $3
		ORDER BY je.posting_date, je.id
	`, accountID, dateFrom, dateTo)
	defer rows.Close()

	balance := report.OpeningBalance
	for rows.Next() {
		var row models.AccountActivityRow
		var ref, srcMod *string

		err := rows.Scan(&row.Date, &row.JournalNum, &row.Description, &ref, &row.Debit, &row.Credit, &srcMod)
		if err != nil {
			return nil, fmt.Errorf("scanning activity: %w", err)
		}

		if ref != nil {
			row.Reference = *ref
		}
		if srcMod != nil {
			row.SourceModule = *srcMod
		}

		balance += row.Debit - row.Credit
		row.Balance = balance

		report.Activity = append(report.Activity, row)
	}

	report.ClosingBalance = balance

	return report, nil
}

// ============================================
// Integration
// ============================================

func (s *glServiceImpl) PostFromAR(ctx context.Context, invoiceID int, createdBy int) (int, error) {
	// TODO: Implement AR to GL posting
	return 0, fmt.Errorf("not implemented")
}

func (s *glServiceImpl) PostFromAP(ctx context.Context, invoiceID int, createdBy int) (int, error) {
	// TODO: Implement AP to GL posting
	return 0, fmt.Errorf("not implemented")
}
