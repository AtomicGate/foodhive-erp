package models

// ============================================
// Catch Weight Models
// ============================================
// Used for products where actual weight varies from standard weight
// Example: A box of chicken ordered as "10 pieces" but actual weight is 17.55kg
// Each piece can have individual weights: 1.7kg, 1.8kg, 1.65kg, etc.

// WeightUOM - Weight Unit of Measure
type WeightUOM string

const (
	WeightUOMKG WeightUOM = "KG"
	WeightUOMLB WeightUOM = "LB"
	WeightUOMGR WeightUOM = "GR"
	WeightUOMOZ WeightUOM = "OZ"
)

// ============================================
// Catch Weight Entry (Header)
// ============================================

type CatchWeightEntry struct {
	ID              int            `json:"id"`
	ProductID       int            `json:"product_id"`
	ReferenceType   string         `json:"reference_type"` // "RECEIVING", "SALES", "PICKING", "ADJUSTMENT"
	ReferenceID     int            `json:"reference_id"`   // receiving_id, sales_order_id, pick_list_id
	ReferenceNumber string         `json:"reference_number"`
	LotNumber       string         `json:"lot_number,omitempty"`
	ExpectedWeight  float64        `json:"expected_weight"` // Standard/ordered weight
	ActualWeight    float64        `json:"actual_weight"`   // Total actual weight captured
	WeightUOM       WeightUOM      `json:"weight_uom"`
	PieceCount      int            `json:"piece_count"`      // Number of pieces
	Variance        float64        `json:"variance"`         // Actual - Expected
	VariancePercent float64        `json:"variance_percent"` // (Variance / Expected) * 100
	IsBilled        bool           `json:"is_billed"`        // Has this been invoiced?
	CapturedBy      int            `json:"captured_by"`
	CapturedAt      CustomDateTime `json:"captured_at"`
	Notes           string         `json:"notes,omitempty"`
}

// ============================================
// Catch Weight Piece (Individual Item Weights)
// ============================================

type CatchWeightPiece struct {
	ID           int            `json:"id"`
	EntryID      int            `json:"entry_id"`     // Foreign key to CatchWeightEntry
	PieceNumber  int            `json:"piece_number"` // 1, 2, 3, etc.
	Weight       float64        `json:"weight"`
	WeightUOM    WeightUOM      `json:"weight_uom"`
	Barcode      string         `json:"barcode,omitempty"` // Individual piece barcode if any
	TagNumber    string         `json:"tag_number,omitempty"`
	QualityGrade string         `json:"quality_grade,omitempty"` // A, B, C grade
	Temperature  float64        `json:"temperature,omitempty"`   // For cold chain
	CapturedAt   CustomDateTime `json:"captured_at"`
	Notes        string         `json:"notes,omitempty"`
}

// ============================================
// Catch Weight Entry with Pieces
// ============================================

type CatchWeightEntryWithPieces struct {
	Entry        CatchWeightEntry   `json:"entry"`
	Pieces       []CatchWeightPiece `json:"pieces"`
	ProductSKU   string             `json:"product_sku"`
	ProductName  string             `json:"product_name"`
	AveragePiece float64            `json:"average_piece_weight"`
	MinPiece     float64            `json:"min_piece_weight"`
	MaxPiece     float64            `json:"max_piece_weight"`
}

// ============================================
// Catch Weight Product Configuration
// ============================================

type CatchWeightConfig struct {
	ProductID           int       `json:"product_id"`
	StandardPieceWeight float64   `json:"standard_piece_weight"` // Expected weight per piece
	WeightUOM           WeightUOM `json:"weight_uom"`
	MinWeight           float64   `json:"min_weight"`            // Minimum acceptable weight
	MaxWeight           float64   `json:"max_weight"`            // Maximum acceptable weight
	VarianceTolerance   float64   `json:"variance_tolerance"`    // Allowed variance % (e.g., 5%)
	RequirePieceWeights bool      `json:"require_piece_weights"` // Must capture each piece?
	PricingMethod       string    `json:"pricing_method"`        // "ACTUAL_WEIGHT", "STANDARD_WEIGHT", "CATCH_UP"
}

// ============================================
// Request Types
// ============================================

// Used when receiving goods or shipping
type CaptureCatchWeightRequest struct {
	ProductID       int                         `json:"product_id"`
	ReferenceType   string                      `json:"reference_type"` // "RECEIVING", "SALES", "PICKING"
	ReferenceID     int                         `json:"reference_id"`
	ReferenceNumber string                      `json:"reference_number"`
	LotNumber       string                      `json:"lot_number,omitempty"`
	ExpectedWeight  float64                     `json:"expected_weight"`
	WeightUOM       WeightUOM                   `json:"weight_uom"`
	Pieces          []CapturePieceWeightRequest `json:"pieces"`
	Notes           string                      `json:"notes,omitempty"`
}

type CapturePieceWeightRequest struct {
	Weight       float64 `json:"weight"`
	Barcode      string  `json:"barcode,omitempty"`
	TagNumber    string  `json:"tag_number,omitempty"`
	QualityGrade string  `json:"quality_grade,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	Notes        string  `json:"notes,omitempty"`
}

// Quick capture - just total weight without piece details
type QuickCatchWeightRequest struct {
	ProductID       int       `json:"product_id"`
	ReferenceType   string    `json:"reference_type"`
	ReferenceID     int       `json:"reference_id"`
	ReferenceNumber string    `json:"reference_number"`
	LotNumber       string    `json:"lot_number,omitempty"`
	ExpectedWeight  float64   `json:"expected_weight"`
	ActualWeight    float64   `json:"actual_weight"`
	PieceCount      int       `json:"piece_count"`
	WeightUOM       WeightUOM `json:"weight_uom"`
	Notes           string    `json:"notes,omitempty"`
}

// Update product catch weight configuration
type UpdateCatchWeightConfigRequest struct {
	StandardPieceWeight float64   `json:"standard_piece_weight"`
	WeightUOM           WeightUOM `json:"weight_uom"`
	MinWeight           float64   `json:"min_weight"`
	MaxWeight           float64   `json:"max_weight"`
	VarianceTolerance   float64   `json:"variance_tolerance"`
	RequirePieceWeights bool      `json:"require_piece_weights"`
	PricingMethod       string    `json:"pricing_method"`
}

// ============================================
// Catch Weight Reports
// ============================================

type CatchWeightVarianceReport struct {
	ProductID       int       `json:"product_id"`
	ProductSKU      string    `json:"product_sku"`
	ProductName     string    `json:"product_name"`
	TotalEntries    int       `json:"total_entries"`
	TotalExpected   float64   `json:"total_expected"`
	TotalActual     float64   `json:"total_actual"`
	TotalVariance   float64   `json:"total_variance"`
	VariancePercent float64   `json:"variance_percent"`
	WeightUOM       WeightUOM `json:"weight_uom"`
}

type CatchWeightSummaryByLot struct {
	LotNumber     string    `json:"lot_number"`
	ProductID     int       `json:"product_id"`
	ProductSKU    string    `json:"product_sku"`
	TotalPieces   int       `json:"total_pieces"`
	TotalWeight   float64   `json:"total_weight"`
	AverageWeight float64   `json:"average_weight"`
	MinWeight     float64   `json:"min_weight"`
	MaxWeight     float64   `json:"max_weight"`
	WeightUOM     WeightUOM `json:"weight_uom"`
}

type CatchWeightListFilters struct {
	ProductID     *int   `json:"product_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	LotNumber     string `json:"lot_number,omitempty"`
	DateFrom      string `json:"date_from,omitempty"`
	DateTo        string `json:"date_to,omitempty"`
	HasVariance   bool   `json:"has_variance,omitempty"` // Only show entries with variance
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
}

// ============================================
// Catch Weight Billing
// ============================================

type CatchWeightBillingAdjustment struct {
	InvoiceID        int     `json:"invoice_id"`
	InvoiceLineID    int     `json:"invoice_line_id"`
	ProductID        int     `json:"product_id"`
	StandardWeight   float64 `json:"standard_weight"`
	ActualWeight     float64 `json:"actual_weight"`
	UnitPrice        float64 `json:"unit_price"`
	StandardAmount   float64 `json:"standard_amount"`
	ActualAmount     float64 `json:"actual_amount"`
	AdjustmentAmount float64 `json:"adjustment_amount"`
}

// ============================================
// Validation
// ============================================

func ValidateCatchWeightCapture(v *Validator, req *CaptureCatchWeightRequest) {
	v.Check(req.ProductID > 0, "product_id", "Product is required")
	v.Check(req.ReferenceType != "", "reference_type", "Reference type is required")
	v.Check(req.ReferenceID > 0, "reference_id", "Reference ID is required")
	v.Check(req.ExpectedWeight > 0, "expected_weight", "Expected weight must be positive")
	v.Check(len(req.Pieces) > 0, "pieces", "At least one piece weight is required")

	for i, piece := range req.Pieces {
		v.Check(piece.Weight > 0, "pieces", "Piece "+string(rune(i+1))+" weight must be positive")
	}
}

func ValidateQuickCatchWeight(v *Validator, req *QuickCatchWeightRequest) {
	v.Check(req.ProductID > 0, "product_id", "Product is required")
	v.Check(req.ReferenceType != "", "reference_type", "Reference type is required")
	v.Check(req.ReferenceID > 0, "reference_id", "Reference ID is required")
	v.Check(req.ActualWeight > 0, "actual_weight", "Actual weight must be positive")
	v.Check(req.PieceCount > 0, "piece_count", "Piece count must be positive")
}

// ============================================
// Weight Conversion Helpers
// ============================================

// ConvertWeight converts weight between units
func ConvertWeight(weight float64, fromUOM, toUOM WeightUOM) float64 {
	// First convert to grams (base unit)
	var grams float64
	switch fromUOM {
	case WeightUOMKG:
		grams = weight * 1000
	case WeightUOMLB:
		grams = weight * 453.592
	case WeightUOMGR:
		grams = weight
	case WeightUOMOZ:
		grams = weight * 28.3495
	default:
		grams = weight * 1000 // Default assume KG
	}

	// Then convert to target unit
	switch toUOM {
	case WeightUOMKG:
		return grams / 1000
	case WeightUOMLB:
		return grams / 453.592
	case WeightUOMGR:
		return grams
	case WeightUOMOZ:
		return grams / 28.3495
	default:
		return grams / 1000 // Default return KG
	}
}
