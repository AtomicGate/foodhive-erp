package models

// ============================================
// Warehouse Models
// ============================================

type Warehouse struct {
	ID            int        `json:"id"`
	WarehouseCode string     `json:"warehouse_code"`
	Name          string     `json:"name"`
	AddressLine1  string     `json:"address_line1,omitempty"`
	AddressLine2  string     `json:"address_line2,omitempty"`
	City          string     `json:"city,omitempty"`
	State         string     `json:"state,omitempty"`
	PostalCode    string     `json:"postal_code,omitempty"`
	Country       string     `json:"country,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     CustomDate `json:"created_at"`
}

type WarehouseZone struct {
	ID                    int     `json:"id"`
	WarehouseID           int     `json:"warehouse_id"`
	ZoneCode              string  `json:"zone_code"`
	Name                  string  `json:"name,omitempty"`
	ZoneType              string  `json:"zone_type,omitempty"`
	TemperatureControlled bool    `json:"temperature_controlled"`
	MinTemperature        float64 `json:"min_temperature,omitempty"`
	MaxTemperature        float64 `json:"max_temperature,omitempty"`
}

type WarehouseLocation struct {
	ID           int     `json:"id"`
	WarehouseID  int     `json:"warehouse_id"`
	ZoneID       *int    `json:"zone_id,omitempty"`
	LocationCode string  `json:"location_code"`
	Aisle        string  `json:"aisle,omitempty"`
	Rack         string  `json:"rack,omitempty"`
	Shelf        string  `json:"shelf,omitempty"`
	Bin          string  `json:"bin,omitempty"`
	LocationType string  `json:"location_type,omitempty"`
	MaxWeight    float64 `json:"max_weight,omitempty"`
	MaxVolume    float64 `json:"max_volume,omitempty"`
	IsActive     bool    `json:"is_active"`
	PickSequence *int    `json:"pick_sequence,omitempty"`
}

type WarehouseWithDetails struct {
	Warehouse Warehouse           `json:"warehouse"`
	Zones     []WarehouseZone     `json:"zones,omitempty"`
	Locations []WarehouseLocation `json:"locations,omitempty"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateWarehouseRequest struct {
	WarehouseCode string `json:"warehouse_code"`
	Name          string `json:"name"`
	AddressLine1  string `json:"address_line1,omitempty"`
	AddressLine2  string `json:"address_line2,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Country       string `json:"country,omitempty"`
}

type UpdateWarehouseRequest struct {
	Name         *string `json:"name,omitempty"`
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      *string `json:"country,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type WarehouseListFilters struct {
	Search   string `json:"search,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type CreateZoneRequest struct {
	WarehouseID           int     `json:"warehouse_id"`
	ZoneCode              string  `json:"zone_code"`
	Name                  string  `json:"name,omitempty"`
	ZoneType              string  `json:"zone_type,omitempty"`
	TemperatureControlled bool    `json:"temperature_controlled"`
	MinTemperature        float64 `json:"min_temperature,omitempty"`
	MaxTemperature        float64 `json:"max_temperature,omitempty"`
}

type CreateLocationRequest struct {
	WarehouseID  int     `json:"warehouse_id"`
	ZoneID       *int    `json:"zone_id,omitempty"`
	LocationCode string  `json:"location_code"`
	Aisle        string  `json:"aisle,omitempty"`
	Rack         string  `json:"rack,omitempty"`
	Shelf        string  `json:"shelf,omitempty"`
	Bin          string  `json:"bin,omitempty"`
	LocationType string  `json:"location_type,omitempty"`
	MaxWeight    float64 `json:"max_weight,omitempty"`
	MaxVolume    float64 `json:"max_volume,omitempty"`
	PickSequence *int    `json:"pick_sequence,omitempty"`
}

// ============================================
// Validation
// ============================================

func ValidateWarehouse(v *Validator, req *CreateWarehouseRequest) {
	v.Check(req.WarehouseCode != "", "warehouse_code", "Warehouse code is required")
	v.Check(len(req.WarehouseCode) <= 10, "warehouse_code", "Warehouse code must be 10 characters or less")
	v.Check(req.Name != "", "name", "Warehouse name is required")
}

func ValidateZone(v *Validator, req *CreateZoneRequest) {
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse ID is required")
	v.Check(req.ZoneCode != "", "zone_code", "Zone code is required")
	v.Check(len(req.ZoneCode) <= 10, "zone_code", "Zone code must be 10 characters or less")
	if req.TemperatureControlled {
		v.Check(req.MaxTemperature >= req.MinTemperature, "temperature", "Max temperature must be >= min temperature")
	}
}

func ValidateLocation(v *Validator, req *CreateLocationRequest) {
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse ID is required")
	v.Check(req.LocationCode != "", "location_code", "Location code is required")
	v.Check(len(req.LocationCode) <= 50, "location_code", "Location code must be 50 characters or less")
	if req.MaxWeight != 0 {
		v.Check(req.MaxWeight > 0, "max_weight", "Max weight must be greater than 0")
	}
	if req.MaxVolume != 0 {
		v.Check(req.MaxVolume > 0, "max_volume", "Max volume must be greater than 0")
	}
}
