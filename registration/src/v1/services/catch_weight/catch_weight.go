package catch_weight

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound         = errors.New("catch weight entry not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrNotCatchWeight   = errors.New("product is not configured for catch weight")
	ErrVarianceExceeded = errors.New("weight variance exceeds tolerance")
	ErrWeightOutOfRange = errors.New("piece weight is outside acceptable range")
	ErrDuplicateEntry   = errors.New("catch weight already captured for this reference")
)

type CatchWeightService interface {
	// Capture Operations
	CaptureCatchWeight(ctx context.Context, req models.CaptureCatchWeightRequest, capturedBy int) (int, error)
	QuickCaptureCatchWeight(ctx context.Context, req models.QuickCatchWeightRequest, capturedBy int) (int, error)
	AddPieceWeight(ctx context.Context, entryID int, piece models.CapturePieceWeightRequest) (int, error)
	UpdatePieceWeight(ctx context.Context, pieceID int, weight float64) error
	DeletePieceWeight(ctx context.Context, pieceID int) error

	// Retrieval
	GetEntryByID(ctx context.Context, id int) (*models.CatchWeightEntryWithPieces, error)
	GetEntryByReference(ctx context.Context, refType string, refID int, productID int) (*models.CatchWeightEntryWithPieces, error)
	ListEntries(ctx context.Context, filters models.CatchWeightListFilters) ([]models.CatchWeightEntry, int64, error)
	GetPiecesByEntry(ctx context.Context, entryID int) ([]models.CatchWeightPiece, error)

	// Product Configuration
	GetProductConfig(ctx context.Context, productID int) (*models.CatchWeightConfig, error)
	UpdateProductConfig(ctx context.Context, productID int, req models.UpdateCatchWeightConfigRequest) error

	// Reports
	GetVarianceReport(ctx context.Context, productID *int, dateFrom, dateTo string) ([]models.CatchWeightVarianceReport, error)
	GetLotSummary(ctx context.Context, productID int, lotNumber string) (*models.CatchWeightSummaryByLot, error)

	// Billing Integration
	CalculateBillingAdjustment(ctx context.Context, invoiceID int, invoiceLineID int, productID int, standardWeight float64) (*models.CatchWeightBillingAdjustment, error)
	MarkAsBilled(ctx context.Context, entryID int) error

	// Validation
	ValidateWeight(ctx context.Context, productID int, weight float64) error
}

type catchWeightServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) CatchWeightService {
	return &catchWeightServiceImpl{db: db}
}

// ============================================
// Capture Operations
// ============================================

func (s *catchWeightServiceImpl) CaptureCatchWeight(ctx context.Context, req models.CaptureCatchWeightRequest, capturedBy int) (int, error) {
	// Verify product exists and is catch weight
	var isCatchWeight bool
	err := s.db.QueryRow(ctx, `
		SELECT is_catch_weight FROM products WHERE id = $1
	`, req.ProductID).Scan(&isCatchWeight)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrProductNotFound
		}
		return 0, fmt.Errorf("checking product: %w", err)
	}
	if !isCatchWeight {
		return 0, ErrNotCatchWeight
	}

	// Check for duplicate
	var exists bool
	err = s.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM catch_weight_entries 
			WHERE reference_type = $1 AND reference_id = $2 AND product_id = $3
		)
	`, req.ReferenceType, req.ReferenceID, req.ProductID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("checking duplicate: %w", err)
	}
	if exists {
		return 0, ErrDuplicateEntry
	}

	// Calculate totals from pieces
	var totalWeight float64
	for _, piece := range req.Pieces {
		totalWeight += piece.Weight
	}

	variance := totalWeight - req.ExpectedWeight
	variancePercent := 0.0
	if req.ExpectedWeight > 0 {
		variancePercent = (variance / req.ExpectedWeight) * 100
	}

	// Insert entry
	var entryID int
	err = s.db.QueryRow(ctx, `
		INSERT INTO catch_weight_entries (
			product_id, reference_type, reference_id, reference_number,
			lot_number, expected_weight, actual_weight, weight_uom,
			piece_count, variance, variance_percent, captured_by, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`,
		req.ProductID, req.ReferenceType, req.ReferenceID, req.ReferenceNumber,
		req.LotNumber, req.ExpectedWeight, totalWeight, req.WeightUOM,
		len(req.Pieces), variance, variancePercent, capturedBy, req.Notes,
	).Scan(&entryID)
	if err != nil {
		return 0, fmt.Errorf("inserting entry: %w", err)
	}

	// Insert individual pieces
	for i, piece := range req.Pieces {
		_, err = s.db.Exec(ctx, `
			INSERT INTO catch_weight_pieces (
				entry_id, piece_number, weight, weight_uom,
				barcode, tag_number, quality_grade, temperature, notes
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			entryID, i+1, piece.Weight, req.WeightUOM,
			piece.Barcode, piece.TagNumber, piece.QualityGrade, piece.Temperature, piece.Notes,
		)
		if err != nil {
			return 0, fmt.Errorf("inserting piece %d: %w", i+1, err)
		}
	}

	return entryID, nil
}

func (s *catchWeightServiceImpl) QuickCaptureCatchWeight(ctx context.Context, req models.QuickCatchWeightRequest, capturedBy int) (int, error) {
	// Verify product exists and is catch weight
	var isCatchWeight bool
	err := s.db.QueryRow(ctx, `
		SELECT is_catch_weight FROM products WHERE id = $1
	`, req.ProductID).Scan(&isCatchWeight)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrProductNotFound
		}
		return 0, fmt.Errorf("checking product: %w", err)
	}
	if !isCatchWeight {
		return 0, ErrNotCatchWeight
	}

	variance := req.ActualWeight - req.ExpectedWeight
	variancePercent := 0.0
	if req.ExpectedWeight > 0 {
		variancePercent = (variance / req.ExpectedWeight) * 100
	}

	// Insert entry (no individual pieces)
	var entryID int
	err = s.db.QueryRow(ctx, `
		INSERT INTO catch_weight_entries (
			product_id, reference_type, reference_id, reference_number,
			lot_number, expected_weight, actual_weight, weight_uom,
			piece_count, variance, variance_percent, captured_by, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (reference_type, reference_id, product_id) 
		DO UPDATE SET 
			actual_weight = EXCLUDED.actual_weight,
			piece_count = EXCLUDED.piece_count,
			variance = EXCLUDED.variance,
			variance_percent = EXCLUDED.variance_percent,
			captured_by = EXCLUDED.captured_by,
			captured_at = NOW()
		RETURNING id
	`,
		req.ProductID, req.ReferenceType, req.ReferenceID, req.ReferenceNumber,
		req.LotNumber, req.ExpectedWeight, req.ActualWeight, req.WeightUOM,
		req.PieceCount, variance, variancePercent, capturedBy, req.Notes,
	).Scan(&entryID)
	if err != nil {
		return 0, fmt.Errorf("inserting entry: %w", err)
	}

	return entryID, nil
}

func (s *catchWeightServiceImpl) AddPieceWeight(ctx context.Context, entryID int, piece models.CapturePieceWeightRequest) (int, error) {
	// Get current max piece number
	var maxPieceNum int
	var weightUOM models.WeightUOM
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(p.piece_number), 0), e.weight_uom
		FROM catch_weight_entries e
		LEFT JOIN catch_weight_pieces p ON e.id = p.entry_id
		WHERE e.id = $1
		GROUP BY e.weight_uom
	`, entryID).Scan(&maxPieceNum, &weightUOM)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("getting entry: %w", err)
	}

	var pieceID int
	err = s.db.QueryRow(ctx, `
		INSERT INTO catch_weight_pieces (
			entry_id, piece_number, weight, weight_uom,
			barcode, tag_number, quality_grade, temperature, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		entryID, maxPieceNum+1, piece.Weight, weightUOM,
		piece.Barcode, piece.TagNumber, piece.QualityGrade, piece.Temperature, piece.Notes,
	).Scan(&pieceID)
	if err != nil {
		return 0, fmt.Errorf("inserting piece: %w", err)
	}

	// Update entry totals
	err = s.updateEntryTotals(ctx, entryID)
	if err != nil {
		return 0, err
	}

	return pieceID, nil
}

func (s *catchWeightServiceImpl) UpdatePieceWeight(ctx context.Context, pieceID int, weight float64) error {
	var entryID int
	err := s.db.QueryRow(ctx, `
		UPDATE catch_weight_pieces 
		SET weight = $1, captured_at = NOW()
		WHERE id = $2
		RETURNING entry_id
	`, weight, pieceID).Scan(&entryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("updating piece: %w", err)
	}

	return s.updateEntryTotals(ctx, entryID)
}

func (s *catchWeightServiceImpl) DeletePieceWeight(ctx context.Context, pieceID int) error {
	var entryID int
	err := s.db.QueryRow(ctx, `
		DELETE FROM catch_weight_pieces WHERE id = $1 RETURNING entry_id
	`, pieceID).Scan(&entryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("deleting piece: %w", err)
	}

	return s.updateEntryTotals(ctx, entryID)
}

func (s *catchWeightServiceImpl) updateEntryTotals(ctx context.Context, entryID int) error {
	_, err := s.db.Exec(ctx, `
		UPDATE catch_weight_entries e
		SET 
			actual_weight = COALESCE((SELECT SUM(weight) FROM catch_weight_pieces WHERE entry_id = e.id), 0),
			piece_count = COALESCE((SELECT COUNT(*) FROM catch_weight_pieces WHERE entry_id = e.id), 0),
			variance = COALESCE((SELECT SUM(weight) FROM catch_weight_pieces WHERE entry_id = e.id), 0) - e.expected_weight,
			variance_percent = CASE 
				WHEN e.expected_weight > 0 THEN 
					((COALESCE((SELECT SUM(weight) FROM catch_weight_pieces WHERE entry_id = e.id), 0) - e.expected_weight) / e.expected_weight) * 100
				ELSE 0
			END
		WHERE id = $1
	`, entryID)
	return err
}

// ============================================
// Retrieval Operations
// ============================================

func (s *catchWeightServiceImpl) GetEntryByID(ctx context.Context, id int) (*models.CatchWeightEntryWithPieces, error) {
	entry := &models.CatchWeightEntryWithPieces{}

	err := s.db.QueryRow(ctx, `
		SELECT e.id, e.product_id, e.reference_type, e.reference_id, e.reference_number,
		       e.lot_number, e.expected_weight, e.actual_weight, e.weight_uom,
		       e.piece_count, e.variance, e.variance_percent, e.is_billed,
		       e.captured_by, e.captured_at, e.notes,
		       p.sku, p.name
		FROM catch_weight_entries e
		JOIN products p ON e.product_id = p.id
		WHERE e.id = $1
	`, id).Scan(
		&entry.Entry.ID, &entry.Entry.ProductID, &entry.Entry.ReferenceType,
		&entry.Entry.ReferenceID, &entry.Entry.ReferenceNumber, &entry.Entry.LotNumber,
		&entry.Entry.ExpectedWeight, &entry.Entry.ActualWeight, &entry.Entry.WeightUOM,
		&entry.Entry.PieceCount, &entry.Entry.Variance, &entry.Entry.VariancePercent,
		&entry.Entry.IsBilled, &entry.Entry.CapturedBy, &entry.Entry.CapturedAt, &entry.Entry.Notes,
		&entry.ProductSKU, &entry.ProductName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting entry: %w", err)
	}

	// Get pieces
	rows := s.db.Query(ctx, `
		SELECT id, entry_id, piece_number, weight, weight_uom,
		       barcode, tag_number, quality_grade, temperature, captured_at, notes
		FROM catch_weight_pieces
		WHERE entry_id = $1
		ORDER BY piece_number
	`, id)
	defer rows.Close()

	var minWeight, maxWeight, totalWeight float64
	minWeight = math.MaxFloat64

	for rows.Next() {
		var piece models.CatchWeightPiece
		var barcode, tagNum, grade, notes *string

		err := rows.Scan(
			&piece.ID, &piece.EntryID, &piece.PieceNumber, &piece.Weight, &piece.WeightUOM,
			&barcode, &tagNum, &grade, &piece.Temperature, &piece.CapturedAt, &notes,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning piece: %w", err)
		}

		if barcode != nil {
			piece.Barcode = *barcode
		}
		if tagNum != nil {
			piece.TagNumber = *tagNum
		}
		if grade != nil {
			piece.QualityGrade = *grade
		}
		if notes != nil {
			piece.Notes = *notes
		}

		entry.Pieces = append(entry.Pieces, piece)

		totalWeight += piece.Weight
		if piece.Weight < minWeight {
			minWeight = piece.Weight
		}
		if piece.Weight > maxWeight {
			maxWeight = piece.Weight
		}
	}

	if len(entry.Pieces) > 0 {
		entry.AveragePiece = totalWeight / float64(len(entry.Pieces))
		entry.MinPiece = minWeight
		entry.MaxPiece = maxWeight
	}

	return entry, nil
}

func (s *catchWeightServiceImpl) GetEntryByReference(ctx context.Context, refType string, refID int, productID int) (*models.CatchWeightEntryWithPieces, error) {
	var entryID int
	err := s.db.QueryRow(ctx, `
		SELECT id FROM catch_weight_entries
		WHERE reference_type = $1 AND reference_id = $2 AND product_id = $3
	`, refType, refID, productID).Scan(&entryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("finding entry: %w", err)
	}
	return s.GetEntryByID(ctx, entryID)
}

func (s *catchWeightServiceImpl) ListEntries(ctx context.Context, filters models.CatchWeightListFilters) ([]models.CatchWeightEntry, int64, error) {
	query := `
		SELECT id, product_id, reference_type, reference_id, reference_number,
		       lot_number, expected_weight, actual_weight, weight_uom,
		       piece_count, variance, variance_percent, is_billed,
		       captured_by, captured_at, notes
		FROM catch_weight_entries
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM catch_weight_entries WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if filters.ProductID != nil {
		query += fmt.Sprintf(" AND product_id = $%d", argNum)
		countQuery += fmt.Sprintf(" AND product_id = $%d", argNum)
		args = append(args, *filters.ProductID)
		argNum++
	}

	if filters.ReferenceType != "" {
		query += fmt.Sprintf(" AND reference_type = $%d", argNum)
		countQuery += fmt.Sprintf(" AND reference_type = $%d", argNum)
		args = append(args, filters.ReferenceType)
		argNum++
	}

	if filters.LotNumber != "" {
		query += fmt.Sprintf(" AND lot_number = $%d", argNum)
		countQuery += fmt.Sprintf(" AND lot_number = $%d", argNum)
		args = append(args, filters.LotNumber)
		argNum++
	}

	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND captured_at >= $%d", argNum)
		countQuery += fmt.Sprintf(" AND captured_at >= $%d", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}

	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND captured_at <= $%d", argNum)
		countQuery += fmt.Sprintf(" AND captured_at <= $%d", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	if filters.HasVariance {
		query += " AND variance != 0"
		countQuery += " AND variance != 0"
	}

	// Get total count
	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting entries: %w", err)
	}

	// Pagination
	if filters.PageSize == 0 {
		filters.PageSize = 50
	}
	if filters.Page == 0 {
		filters.Page = 1
	}
	offset := (filters.Page - 1) * filters.PageSize

	query += fmt.Sprintf(" ORDER BY captured_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var entries []models.CatchWeightEntry
	for rows.Next() {
		var e models.CatchWeightEntry
		var lotNum, notes *string

		err := rows.Scan(
			&e.ID, &e.ProductID, &e.ReferenceType, &e.ReferenceID, &e.ReferenceNumber,
			&lotNum, &e.ExpectedWeight, &e.ActualWeight, &e.WeightUOM,
			&e.PieceCount, &e.Variance, &e.VariancePercent, &e.IsBilled,
			&e.CapturedBy, &e.CapturedAt, &notes,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning entry: %w", err)
		}

		if lotNum != nil {
			e.LotNumber = *lotNum
		}
		if notes != nil {
			e.Notes = *notes
		}

		entries = append(entries, e)
	}

	return entries, total, nil
}

func (s *catchWeightServiceImpl) GetPiecesByEntry(ctx context.Context, entryID int) ([]models.CatchWeightPiece, error) {
	rows := s.db.Query(ctx, `
		SELECT id, entry_id, piece_number, weight, weight_uom,
		       barcode, tag_number, quality_grade, temperature, captured_at, notes
		FROM catch_weight_pieces
		WHERE entry_id = $1
		ORDER BY piece_number
	`, entryID)
	defer rows.Close()

	var pieces []models.CatchWeightPiece
	for rows.Next() {
		var piece models.CatchWeightPiece
		var barcode, tagNum, grade, notes *string

		err := rows.Scan(
			&piece.ID, &piece.EntryID, &piece.PieceNumber, &piece.Weight, &piece.WeightUOM,
			&barcode, &tagNum, &grade, &piece.Temperature, &piece.CapturedAt, &notes,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning piece: %w", err)
		}

		if barcode != nil {
			piece.Barcode = *barcode
		}
		if tagNum != nil {
			piece.TagNumber = *tagNum
		}
		if grade != nil {
			piece.QualityGrade = *grade
		}
		if notes != nil {
			piece.Notes = *notes
		}

		pieces = append(pieces, piece)
	}

	return pieces, nil
}

// ============================================
// Product Configuration
// ============================================

func (s *catchWeightServiceImpl) GetProductConfig(ctx context.Context, productID int) (*models.CatchWeightConfig, error) {
	config := &models.CatchWeightConfig{ProductID: productID}

	err := s.db.QueryRow(ctx, `
		SELECT standard_piece_weight, weight_uom, min_weight, max_weight,
		       variance_tolerance, require_piece_weights, pricing_method
		FROM catch_weight_config
		WHERE product_id = $1
	`, productID).Scan(
		&config.StandardPieceWeight, &config.WeightUOM, &config.MinWeight, &config.MaxWeight,
		&config.VarianceTolerance, &config.RequirePieceWeights, &config.PricingMethod,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return default config
			return &models.CatchWeightConfig{
				ProductID:           productID,
				WeightUOM:           models.WeightUOMKG,
				VarianceTolerance:   5.0, // 5% default
				RequirePieceWeights: false,
				PricingMethod:       "ACTUAL_WEIGHT",
			}, nil
		}
		return nil, fmt.Errorf("getting config: %w", err)
	}

	return config, nil
}

func (s *catchWeightServiceImpl) UpdateProductConfig(ctx context.Context, productID int, req models.UpdateCatchWeightConfigRequest) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO catch_weight_config (
			product_id, standard_piece_weight, weight_uom, min_weight, max_weight,
			variance_tolerance, require_piece_weights, pricing_method
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (product_id) DO UPDATE SET
			standard_piece_weight = EXCLUDED.standard_piece_weight,
			weight_uom = EXCLUDED.weight_uom,
			min_weight = EXCLUDED.min_weight,
			max_weight = EXCLUDED.max_weight,
			variance_tolerance = EXCLUDED.variance_tolerance,
			require_piece_weights = EXCLUDED.require_piece_weights,
			pricing_method = EXCLUDED.pricing_method
	`,
		productID, req.StandardPieceWeight, req.WeightUOM, req.MinWeight, req.MaxWeight,
		req.VarianceTolerance, req.RequirePieceWeights, req.PricingMethod,
	)
	return err
}

// ============================================
// Reports
// ============================================

func (s *catchWeightServiceImpl) GetVarianceReport(ctx context.Context, productID *int, dateFrom, dateTo string) ([]models.CatchWeightVarianceReport, error) {
	query := `
		SELECT e.product_id, p.sku, p.name,
		       COUNT(*) as total_entries,
		       SUM(e.expected_weight) as total_expected,
		       SUM(e.actual_weight) as total_actual,
		       SUM(e.variance) as total_variance,
		       e.weight_uom
		FROM catch_weight_entries e
		JOIN products p ON e.product_id = p.id
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if productID != nil {
		query += fmt.Sprintf(" AND e.product_id = $%d", argNum)
		args = append(args, *productID)
		argNum++
	}

	if dateFrom != "" {
		query += fmt.Sprintf(" AND e.captured_at >= $%d", argNum)
		args = append(args, dateFrom)
		argNum++
	}

	if dateTo != "" {
		query += fmt.Sprintf(" AND e.captured_at <= $%d", argNum)
		args = append(args, dateTo)
		argNum++
	}

	query += " GROUP BY e.product_id, p.sku, p.name, e.weight_uom ORDER BY total_variance DESC"

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var reports []models.CatchWeightVarianceReport
	for rows.Next() {
		var r models.CatchWeightVarianceReport
		err := rows.Scan(
			&r.ProductID, &r.ProductSKU, &r.ProductName,
			&r.TotalEntries, &r.TotalExpected, &r.TotalActual,
			&r.TotalVariance, &r.WeightUOM,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning report: %w", err)
		}

		if r.TotalExpected > 0 {
			r.VariancePercent = (r.TotalVariance / r.TotalExpected) * 100
		}

		reports = append(reports, r)
	}

	return reports, nil
}

func (s *catchWeightServiceImpl) GetLotSummary(ctx context.Context, productID int, lotNumber string) (*models.CatchWeightSummaryByLot, error) {
	summary := &models.CatchWeightSummaryByLot{
		ProductID: productID,
		LotNumber: lotNumber,
	}

	err := s.db.QueryRow(ctx, `
		SELECT p.sku,
		       COALESCE(SUM(e.piece_count), 0) as total_pieces,
		       COALESCE(SUM(e.actual_weight), 0) as total_weight,
		       COALESCE(AVG(cwp.weight), 0) as avg_weight,
		       COALESCE(MIN(cwp.weight), 0) as min_weight,
		       COALESCE(MAX(cwp.weight), 0) as max_weight,
		       e.weight_uom
		FROM catch_weight_entries e
		JOIN products p ON e.product_id = p.id
		LEFT JOIN catch_weight_pieces cwp ON e.id = cwp.entry_id
		WHERE e.product_id = $1 AND e.lot_number = $2
		GROUP BY p.sku, e.weight_uom
	`, productID, lotNumber).Scan(
		&summary.ProductSKU, &summary.TotalPieces, &summary.TotalWeight,
		&summary.AverageWeight, &summary.MinWeight, &summary.MaxWeight, &summary.WeightUOM,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting lot summary: %w", err)
	}

	return summary, nil
}

// ============================================
// Billing Integration
// ============================================

func (s *catchWeightServiceImpl) CalculateBillingAdjustment(ctx context.Context, invoiceID int, invoiceLineID int, productID int, standardWeight float64) (*models.CatchWeightBillingAdjustment, error) {
	// Get the actual weight from catch weight entries for this invoice/line
	var actualWeight float64
	var unitPrice float64

	// Try to find catch weight entry linked to this line
	err := s.db.QueryRow(ctx, `
		SELECT cwe.actual_weight
		FROM catch_weight_entries cwe
		WHERE cwe.reference_type = 'SALES' 
		  AND cwe.reference_id = $1 
		  AND cwe.product_id = $2
	`, invoiceID, productID).Scan(&actualWeight)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No catch weight entry, use standard weight
			actualWeight = standardWeight
		} else {
			return nil, fmt.Errorf("getting catch weight: %w", err)
		}
	}

	// Get unit price from invoice line
	err = s.db.QueryRow(ctx, `
		SELECT unit_price FROM ar_invoice_lines WHERE id = $1
	`, invoiceLineID).Scan(&unitPrice)
	if err != nil {
		// Fallback: get from product pricing
		err = s.db.QueryRow(ctx, `
			SELECT unit_price FROM product_prices WHERE product_id = $1 AND is_active = true LIMIT 1
		`, productID).Scan(&unitPrice)
		if err != nil {
			unitPrice = 0
		}
	}

	standardAmount := standardWeight * unitPrice
	actualAmount := actualWeight * unitPrice

	return &models.CatchWeightBillingAdjustment{
		InvoiceID:        invoiceID,
		InvoiceLineID:    invoiceLineID,
		ProductID:        productID,
		StandardWeight:   standardWeight,
		ActualWeight:     actualWeight,
		UnitPrice:        unitPrice,
		StandardAmount:   standardAmount,
		ActualAmount:     actualAmount,
		AdjustmentAmount: actualAmount - standardAmount,
	}, nil
}

func (s *catchWeightServiceImpl) MarkAsBilled(ctx context.Context, entryID int) error {
	result, err := s.db.Exec(ctx, `
		UPDATE catch_weight_entries SET is_billed = true WHERE id = $1
	`, entryID)
	if err != nil {
		return fmt.Errorf("marking as billed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// ============================================
// Validation
// ============================================

func (s *catchWeightServiceImpl) ValidateWeight(ctx context.Context, productID int, weight float64) error {
	config, err := s.GetProductConfig(ctx, productID)
	if err != nil {
		return err
	}

	if config.MinWeight > 0 && weight < config.MinWeight {
		return fmt.Errorf("%w: weight %.2f is below minimum %.2f", ErrWeightOutOfRange, weight, config.MinWeight)
	}

	if config.MaxWeight > 0 && weight > config.MaxWeight {
		return fmt.Errorf("%w: weight %.2f is above maximum %.2f", ErrWeightOutOfRange, weight, config.MaxWeight)
	}

	return nil
}

// Helper to calculate variance percentage
func (s *catchWeightServiceImpl) calculateVariance(expected, actual float64) float64 {
	if expected == 0 {
		return 0
	}
	return ((actual - expected) / expected) * 100
}
