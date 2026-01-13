package warehouse

import (
	"context"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

// ============================================
// Service Interface
// ============================================

type WarehouseService interface {
	// Warehouses
	Create(ctx context.Context, req *models.CreateWarehouseRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.WarehouseWithDetails, error)
	GetByCode(ctx context.Context, code string) (*models.Warehouse, error)
	Update(ctx context.Context, id int, req *models.UpdateWarehouseRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters *models.WarehouseListFilters) ([]models.Warehouse, int64, error)

	// Zones
	CreateZone(ctx context.Context, req *models.CreateZoneRequest) (int, error)
	GetZoneByID(ctx context.Context, id int) (*models.WarehouseZone, error)
	GetZonesByWarehouse(ctx context.Context, warehouseID int) ([]models.WarehouseZone, error)
	UpdateZone(ctx context.Context, id int, req *models.CreateZoneRequest) error
	DeleteZone(ctx context.Context, id int) error

	// Locations
	CreateLocation(ctx context.Context, req *models.CreateLocationRequest) (int, error)
	GetLocationByID(ctx context.Context, id int) (*models.WarehouseLocation, error)
	GetLocationByCode(ctx context.Context, code string) (*models.WarehouseLocation, error)
	GetLocationsByWarehouse(ctx context.Context, warehouseID int) ([]models.WarehouseLocation, error)
	GetLocationsByZone(ctx context.Context, zoneID int) ([]models.WarehouseLocation, error)
	UpdateLocation(ctx context.Context, id int, req *models.CreateLocationRequest) error
	DeleteLocation(ctx context.Context, id int) error
}

// ============================================
// Service Implementation
// ============================================

type warehouseServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) WarehouseService {
	return &warehouseServiceImpl{db: db}
}

// ============================================
// Warehouse CRUD
// ============================================

func (s *warehouseServiceImpl) Create(ctx context.Context, req *models.CreateWarehouseRequest) (int, error) {
	query := `
		INSERT INTO warehouses (
			warehouse_code, name, address_line1, address_line2,
			city, state, postal_code, country, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.WarehouseCode, req.Name, req.AddressLine1, req.AddressLine2,
		req.City, req.State, req.PostalCode, req.Country,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create warehouse: %w", err)
	}
	return id, nil
}

func (s *warehouseServiceImpl) GetByID(ctx context.Context, id int) (*models.WarehouseWithDetails, error) {
	warehouse, err := s.getWarehouse(ctx, "id = $1", id)
	if err != nil {
		return nil, err
	}

	result := &models.WarehouseWithDetails{
		Warehouse: *warehouse,
	}

	// Get zones
	zones, _ := s.GetZonesByWarehouse(ctx, id)
	result.Zones = zones

	// Get locations
	locations, _ := s.GetLocationsByWarehouse(ctx, id)
	result.Locations = locations

	return result, nil
}

func (s *warehouseServiceImpl) GetByCode(ctx context.Context, code string) (*models.Warehouse, error) {
	return s.getWarehouse(ctx, "warehouse_code = $1", code)
}

func (s *warehouseServiceImpl) getWarehouse(ctx context.Context, whereClause string, arg interface{}) (*models.Warehouse, error) {
	query := fmt.Sprintf(`
		SELECT id, warehouse_code, name, address_line1, address_line2,
			   city, state, postal_code, country, is_active, created_at
		FROM warehouses
		WHERE %s`, whereClause)

	var w models.Warehouse
	var addr1, addr2, city, state, postal, country *string

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&w.ID, &w.WarehouseCode, &w.Name, &addr1, &addr2,
		&city, &state, &postal, &country, &w.IsActive, &w.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("warehouse not found")
		}
		return nil, fmt.Errorf("failed to get warehouse: %w", err)
	}

	if addr1 != nil {
		w.AddressLine1 = *addr1
	}
	if addr2 != nil {
		w.AddressLine2 = *addr2
	}
	if city != nil {
		w.City = *city
	}
	if state != nil {
		w.State = *state
	}
	if postal != nil {
		w.PostalCode = *postal
	}
	if country != nil {
		w.Country = *country
	}

	return &w, nil
}

func (s *warehouseServiceImpl) Update(ctx context.Context, id int, req *models.UpdateWarehouseRequest) error {
	query := `
		UPDATE warehouses SET
			name = COALESCE($1, name),
			address_line1 = COALESCE($2, address_line1),
			address_line2 = COALESCE($3, address_line2),
			city = COALESCE($4, city),
			state = COALESCE($5, state),
			postal_code = COALESCE($6, postal_code),
			country = COALESCE($7, country),
			is_active = COALESCE($8, is_active)
		WHERE id = $9`

	result, err := s.db.Exec(ctx, query,
		req.Name, req.AddressLine1, req.AddressLine2,
		req.City, req.State, req.PostalCode, req.Country, req.IsActive, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update warehouse: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse not found")
	}

	return nil
}

func (s *warehouseServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE warehouses SET is_active = false WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse not found")
	}
	return nil
}

func (s *warehouseServiceImpl) List(ctx context.Context, filters *models.WarehouseListFilters) ([]models.Warehouse, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filters.Search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR warehouse_code ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filters.Search+"%")
		argNum++
	}

	if filters.IsActive != nil {
		whereClause += fmt.Sprintf(" AND is_active = $%d", argNum)
		args = append(args, *filters.IsActive)
		argNum++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warehouses %s", whereClause)
	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count warehouses: %w", err)
	}

	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT id, warehouse_code, name, address_line1, address_line2,
			   city, state, postal_code, country, is_active, created_at
		FROM warehouses
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		var addr1, addr2, city, state, postal, country *string

		err := rows.Scan(
			&w.ID, &w.WarehouseCode, &w.Name, &addr1, &addr2,
			&city, &state, &postal, &country, &w.IsActive, &w.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan warehouse: %w", err)
		}

		if addr1 != nil {
			w.AddressLine1 = *addr1
		}
		if addr2 != nil {
			w.AddressLine2 = *addr2
		}
		if city != nil {
			w.City = *city
		}
		if state != nil {
			w.State = *state
		}
		if postal != nil {
			w.PostalCode = *postal
		}
		if country != nil {
			w.Country = *country
		}

		warehouses = append(warehouses, w)
	}

	return warehouses, total, nil
}

// ============================================
// Zone CRUD
// ============================================

func (s *warehouseServiceImpl) CreateZone(ctx context.Context, req *models.CreateZoneRequest) (int, error) {
	query := `
		INSERT INTO warehouse_zones (
			warehouse_id, zone_code, name, zone_type,
			temperature_controlled, min_temperature, max_temperature
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.WarehouseID, req.ZoneCode, req.Name, req.ZoneType,
		req.TemperatureControlled, req.MinTemperature, req.MaxTemperature,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create zone: %w", err)
	}
	return id, nil
}

func (s *warehouseServiceImpl) GetZoneByID(ctx context.Context, id int) (*models.WarehouseZone, error) {
	query := `
		SELECT id, warehouse_id, zone_code, name, zone_type,
			   temperature_controlled, min_temperature, max_temperature
		FROM warehouse_zones
		WHERE id = $1`

	var z models.WarehouseZone
	var name, zoneType *string
	var minTemp, maxTemp *float64

	err := s.db.QueryRow(ctx, query, id).Scan(
		&z.ID, &z.WarehouseID, &z.ZoneCode, &name, &zoneType,
		&z.TemperatureControlled, &minTemp, &maxTemp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("zone not found")
		}
		return nil, fmt.Errorf("failed to get zone: %w", err)
	}

	if name != nil {
		z.Name = *name
	}
	if zoneType != nil {
		z.ZoneType = *zoneType
	}
	if minTemp != nil {
		z.MinTemperature = *minTemp
	}
	if maxTemp != nil {
		z.MaxTemperature = *maxTemp
	}

	return &z, nil
}

func (s *warehouseServiceImpl) GetZonesByWarehouse(ctx context.Context, warehouseID int) ([]models.WarehouseZone, error) {
	query := `
		SELECT id, warehouse_id, zone_code, name, zone_type,
			   temperature_controlled, min_temperature, max_temperature
		FROM warehouse_zones
		WHERE warehouse_id = $1
		ORDER BY zone_code ASC`

	rows := s.db.Query(ctx, query, warehouseID)
	defer rows.Close()

	var zones []models.WarehouseZone
	for rows.Next() {
		var z models.WarehouseZone
		var name, zoneType *string
		var minTemp, maxTemp *float64

		err := rows.Scan(
			&z.ID, &z.WarehouseID, &z.ZoneCode, &name, &zoneType,
			&z.TemperatureControlled, &minTemp, &maxTemp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan zone: %w", err)
		}

		if name != nil {
			z.Name = *name
		}
		if zoneType != nil {
			z.ZoneType = *zoneType
		}
		if minTemp != nil {
			z.MinTemperature = *minTemp
		}
		if maxTemp != nil {
			z.MaxTemperature = *maxTemp
		}

		zones = append(zones, z)
	}

	return zones, nil
}

func (s *warehouseServiceImpl) UpdateZone(ctx context.Context, id int, req *models.CreateZoneRequest) error {
	query := `
		UPDATE warehouse_zones SET
			zone_code = $1, name = $2, zone_type = $3,
			temperature_controlled = $4, min_temperature = $5, max_temperature = $6
		WHERE id = $7`

	result, err := s.db.Exec(ctx, query,
		req.ZoneCode, req.Name, req.ZoneType,
		req.TemperatureControlled, req.MinTemperature, req.MaxTemperature, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("zone not found")
	}

	return nil
}

func (s *warehouseServiceImpl) DeleteZone(ctx context.Context, id int) error {
	// Check if zone has locations
	var count int
	s.db.QueryRow(ctx, "SELECT COUNT(*) FROM warehouse_locations WHERE zone_id = $1", id).Scan(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete zone with %d locations", count)
	}

	query := `DELETE FROM warehouse_zones WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete zone: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("zone not found")
	}
	return nil
}

// ============================================
// Location CRUD
// ============================================

func (s *warehouseServiceImpl) CreateLocation(ctx context.Context, req *models.CreateLocationRequest) (int, error) {
	query := `
		INSERT INTO warehouse_locations (
			warehouse_id, zone_id, location_code, aisle, rack, shelf, bin,
			location_type, max_weight, max_volume, is_active, pick_sequence
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true, $11)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.WarehouseID, req.ZoneID, req.LocationCode, req.Aisle, req.Rack,
		req.Shelf, req.Bin, req.LocationType, req.MaxWeight, req.MaxVolume, req.PickSequence,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create location: %w", err)
	}
	return id, nil
}

func (s *warehouseServiceImpl) GetLocationByID(ctx context.Context, id int) (*models.WarehouseLocation, error) {
	return s.getLocation(ctx, "id = $1", id)
}

func (s *warehouseServiceImpl) GetLocationByCode(ctx context.Context, code string) (*models.WarehouseLocation, error) {
	return s.getLocation(ctx, "location_code = $1", code)
}

func (s *warehouseServiceImpl) getLocation(ctx context.Context, whereClause string, arg interface{}) (*models.WarehouseLocation, error) {
	query := fmt.Sprintf(`
		SELECT id, warehouse_id, zone_id, location_code, aisle, rack, shelf, bin,
			   location_type, max_weight, max_volume, is_active, pick_sequence
		FROM warehouse_locations
		WHERE %s`, whereClause)

	var l models.WarehouseLocation
	var aisle, rack, shelf, bin, locType *string
	var maxWeight, maxVolume *float64

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&l.ID, &l.WarehouseID, &l.ZoneID, &l.LocationCode, &aisle, &rack, &shelf, &bin,
		&locType, &maxWeight, &maxVolume, &l.IsActive, &l.PickSequence,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("location not found")
		}
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	if aisle != nil {
		l.Aisle = *aisle
	}
	if rack != nil {
		l.Rack = *rack
	}
	if shelf != nil {
		l.Shelf = *shelf
	}
	if bin != nil {
		l.Bin = *bin
	}
	if locType != nil {
		l.LocationType = *locType
	}
	if maxWeight != nil {
		l.MaxWeight = *maxWeight
	}
	if maxVolume != nil {
		l.MaxVolume = *maxVolume
	}

	return &l, nil
}

func (s *warehouseServiceImpl) GetLocationsByWarehouse(ctx context.Context, warehouseID int) ([]models.WarehouseLocation, error) {
	return s.getLocations(ctx, "warehouse_id = $1", warehouseID)
}

func (s *warehouseServiceImpl) GetLocationsByZone(ctx context.Context, zoneID int) ([]models.WarehouseLocation, error) {
	return s.getLocations(ctx, "zone_id = $1", zoneID)
}

func (s *warehouseServiceImpl) getLocations(ctx context.Context, whereClause string, arg interface{}) ([]models.WarehouseLocation, error) {
	query := fmt.Sprintf(`
		SELECT id, warehouse_id, zone_id, location_code, aisle, rack, shelf, bin,
			   location_type, max_weight, max_volume, is_active, pick_sequence
		FROM warehouse_locations
		WHERE %s
		ORDER BY pick_sequence ASC NULLS LAST, location_code ASC`, whereClause)

	rows := s.db.Query(ctx, query, arg)
	defer rows.Close()

	var locations []models.WarehouseLocation
	for rows.Next() {
		var l models.WarehouseLocation
		var aisle, rack, shelf, bin, locType *string
		var maxWeight, maxVolume *float64

		err := rows.Scan(
			&l.ID, &l.WarehouseID, &l.ZoneID, &l.LocationCode, &aisle, &rack, &shelf, &bin,
			&locType, &maxWeight, &maxVolume, &l.IsActive, &l.PickSequence,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}

		if aisle != nil {
			l.Aisle = *aisle
		}
		if rack != nil {
			l.Rack = *rack
		}
		if shelf != nil {
			l.Shelf = *shelf
		}
		if bin != nil {
			l.Bin = *bin
		}
		if locType != nil {
			l.LocationType = *locType
		}
		if maxWeight != nil {
			l.MaxWeight = *maxWeight
		}
		if maxVolume != nil {
			l.MaxVolume = *maxVolume
		}

		locations = append(locations, l)
	}

	return locations, nil
}

func (s *warehouseServiceImpl) UpdateLocation(ctx context.Context, id int, req *models.CreateLocationRequest) error {
	query := `
		UPDATE warehouse_locations SET
			zone_id = $1, location_code = $2, aisle = $3, rack = $4, shelf = $5, bin = $6,
			location_type = $7, max_weight = $8, max_volume = $9, pick_sequence = $10
		WHERE id = $11`

	result, err := s.db.Exec(ctx, query,
		req.ZoneID, req.LocationCode, req.Aisle, req.Rack, req.Shelf, req.Bin,
		req.LocationType, req.MaxWeight, req.MaxVolume, req.PickSequence, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("location not found")
	}

	return nil
}

func (s *warehouseServiceImpl) DeleteLocation(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE warehouse_locations SET is_active = false WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("location not found")
	}
	return nil
}
