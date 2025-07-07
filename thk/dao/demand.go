package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"playground/thk/dberror"
	"playground/thk/helpers"
	"playground/thk/models"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type demandDao struct{}

type Demand interface {
	InsertDemand(db *gorm.DB, demand models.Demand, houseId uuid.UUID) (uuid.UUID, error)
	InsertDemandExpression(db *gorm.DB, demandId uuid.UUID, expression models.DemandExpression) error
	InsertDemandConstant(db *gorm.DB, name string, value float64) error
	InsertDemandHouseUnitLink(db *gorm.DB, demandId uuid.UUID, houseUnitId uuid.UUID) error
	InsertDemandRun(db *gorm.DB, run models.DemandRun) (uuid.UUID, error)
	InsertDemandVariableLink(db *gorm.DB, demandId uuid.UUID, demandVariableName string) error
	SelectDemands(db *gorm.DB, filterOptions *DemandFilter) ([]models.Demand, error)
	SelectDemandById(db *gorm.DB, demandId uuid.UUID, houseId uuid.UUID) (models.Demand, error)
	SelectVariableValues(db *gorm.DB, demand models.Demand, demandRun models.DemandRun) (models.Demand, error)
	SelectLatestDemandRun(db *gorm.DB, demandId uuid.UUID) (models.DemandRun, error)
	// SelectDemandExpression(db *gorm.DB, demandId uuid.UUID) (models.DemandExpression, error)
	// SelectDemandVariables(db *gorm.DB, demandId uuid.UUID) ([]models.DemandVariable, error)
	// SelectDemandConstant(db *gorm.DB, houseId uuid.UUID, name string) (float64, error)
	// SelectDemandHouseUnitLinks(db *gorm.DB, demandId uuid.UUID) ([]models.HouseUnit, error)

}

func (d *demandDao) SelectLatestDemandRun(db *gorm.DB, demandId uuid.UUID) (models.DemandRun, error) {
	sql := `
	SELECT
		dr.id,
		dr.demand_id,
		dr.status,
		dr.from_date,
		dr.to_date,
		dr.create_date,
		dr.update_date
	FROM thk.demand_run dr
	WHERE dr.demand_id = @demandId
	ORDER BY dr.create_date DESC
	LIMIT 1
	`
	params := map[string]interface{}{
		"demandId": demandId,
	}
	var res models.DemandRun
	err := db.Raw(sql, params).Scan(&res).Error
	if err != nil {
		return models.DemandRun{}, dberror.DbErrorFromPq(err)
	}
	return res, nil
}

func (d *demandDao) SelectVariableValues(db *gorm.DB, demand models.Demand, demandRun models.DemandRun) (models.Demand, error) {
	if demand.Id == nil {
		return models.Demand{}, fmt.Errorf("demand id is nil")
	}
	demandId := *demand.Id
	var dvt models.DemandVariableType
	areaWithSql := `area_values AS(
		SELECT id AS house_unit_id ,area 
		FROM thk.house_unit hu
		WHERE hu.id in (SELECT house_unit_id FROM thk.demand_house_unit_lnk dhul WHERE demand_id = @demandId)
	)`
	ownershipWithSql := `ownership_values AS(
		SELECT id AS house_unit_id , ownership 
		FROM thk.house_unit hu
		WHERE hu.id in (SELECT house_unit_id FROM thk.demand_house_unit_lnk dhul WHERE demand_id = @demandId)
	)`
	usageWithSql := `%s_date_range_records AS (
		SELECT
			m.house_unit_id,
			m.id,
			m.record,
			m.record_date
		FROM
			(SELECT 
				m.house_unit_id,
				m.id,
				mr.record,
				mr.record_date,
				row_number() OVER (partition by m.id order by mr.record_date desc) AS row_num
			FROM thk.meter m
			INNER JOIN thk.meter_type_lnk mtl 
			ON m.id = mtl.meter_id
			INNER JOIN thk.meter_type mt 
			ON mt.id = mtl.meter_type_id
			INNER JOIN thk.meter_record mr 
			ON mr.meter_id = m.id
			WHERE 1=1
			AND m.house_unit_id in (
				SELECT house_unit_id FROM thk.demand_house_unit_lnk dhul WHERE dhul.demand_id = @demandId
			)
			AND mt.id = @%sMeterTypeId
			AND mr.record_date < @toDate) AS m
		WHERE 1=1
		AND row_num = 1
		),
		%s_last_before AS (
		SELECT
			house_unit_id,
			id,
			record,
			record_date
		FROM (
			SELECT 
				m.house_unit_id,
				m.id,
				mr.record,
				mr.record_date,
				row_number() OVER (partition by m.id order by mr.record_date desc) AS row_num
			FROM thk.meter m
			INNER JOIN thk.meter_type_lnk mtl 
			ON m.id = mtl.meter_id
			INNER JOIN thk.meter_type mt 
			ON mt.id = mtl.meter_type_id
			INNER JOIN thk.meter_record mr 
			ON mr.meter_id = m.id
			WHERE 1=1
			AND m.house_unit_id in (
				SELECT house_unit_id FROM thk.demand_house_unit_lnk dhul WHERE dhul.demand_id = @demandId
			)
			AND mt.id = @%sMeterTypeId
			AND mr.record_date < @fromDate) AS m
		WHERE 1=1
		AND row_num = 1
		),
		%s_meter_usage AS (
		SELECT 
			drr.house_unit_id,
			drr.id,
			drr.record-lb.record AS usage
		FROM %s_date_range_records drr
		INNER JOIN %s_last_before lb
		ON lb.house_unit_id = drr.house_unit_id AND lb.id = drr.id
		),
		%s_unit_usage AS (
		SELECT
			house_unit_id,
			sum(usage) AS usage
		FROM %s_meter_usage
		group by house_unit_id
	)`

	sql := "with "
	middleSql := "SELECT dhul.house_unit_id AS id, json_build_array(\n"
	endSql := ") AS variables FROM thk.demand_house_unit_lnk dhul\n"
	params := map[string]interface{}{}
	var variablesAlreadyAddedToCTE models.DemandVariables
	for i, variable := range demand.Expression.Variables {
		variableType, err := dvt.FromString(variable.Type)
		if err != nil {
			return demand, err
		}
		_, found := helpers.Find(variablesAlreadyAddedToCTE, func(v models.DemandVariable) bool {
			return v.Name == variable.Name
		})
		if !found {
			if i > 0 {
				sql = sql + ",\n"
				middleSql = middleSql + ",\n"
			}
			endSql = endSql + "LEFT JOIN "
		}
		switch variableType {
		case models.DemandVariableTypeConstant:
			middleSql = middleSql + fmt.Sprintf(`json_build_object(
			'name', @name%d,
			'type', @type%d,
			'value', @value%d
		)`, i, i, i)
			params[fmt.Sprintf("name%d", i)] = variable.Name
			params[fmt.Sprintf("type%d", i)] = variableType.String()
			details := variable.Details
			params[fmt.Sprintf("value%d", i)] = details["value"]
		case models.DemandVariableTypeArea:
			if !found {
				sql = sql + areaWithSql
				params["demandId"] = demandId
			}
			// Middle sql
			middleSql = middleSql + fmt.Sprintf(`json_build_object(
			'name', @name%d,
			'type', @type%d,
			'value', av.area
			)`, i, i)
			params[fmt.Sprintf("name%d", i)] = variable.Name
			params[fmt.Sprintf("type%d", i)] = variableType.String()
			if !found {
				// End sql
				endSql = endSql + `area_values av
					ON dhul.house_unit_id = av.house_unit_id
`
			}
		case models.DemandVariableTypeOwnership:
			if !found {
				sql = sql + ownershipWithSql
				params["demandId"] = demandId
			}
			// Middle sql
			middleSql = middleSql + fmt.Sprintf(`json_build_object(
			'name', @name%d,
			'type', @type%d,
			'value', ov.ownership
			)`, i, i)
			params[fmt.Sprintf("name%d", i)] = variable.Name
			params[fmt.Sprintf("type%d", i)] = variableType.String()

			if !found {
				// End sql
				endSql = endSql + `ownership_values ov
				ON dhul.house_unit_id = ov.house_unit_id
`
			}
		case models.DemandVariableTypeUsage:
			if !found {
				sql = sql + fmt.Sprintf(usageWithSql, variable.Name, variable.Name, variable.Name, variable.Name, variable.Name, variable.Name, variable.Name, variable.Name, variable.Name)
				params["demandId"] = demandId
				// params["fromDate"] = demandRun.FromDate.Format("2006-01-02")
				params["fromDate"] = datatypes.Date(demandRun.FromDate)
				// params["toDate"] = demandRun.ToDate.Format("2006-01-02")
				params["toDate"] = datatypes.Date(demandRun.ToDate)
				log.Printf("VARIABLE DETAILS: %+v", variable.Details)
				log.Printf("METER TYPE ID: %+v", variable.Details["meter_type_id"])
				var meterTypeId uuid.UUID
				meterTypeId, err = uuid.Parse(variable.Details["meter_type_id"].(string))
				if err != nil {
					return demand, err
				}
				params[fmt.Sprintf("%sMeterTypeId", variable.Name)] = meterTypeId
			}
			// Middle sql
			middleSql = middleSql + fmt.Sprintf(`json_build_object(
			'name', @name%d,
			'type', @type%d,
			'value', COALESCE(%s_unit_usage.usage, 0)
			)`, i, i, variable.Name)
			params[fmt.Sprintf("name%d", i)] = variable.Name
			params[fmt.Sprintf("type%d", i)] = variableType.String()

			if !found {
				endSql = endSql + fmt.Sprintf(`%s_unit_usage
				ON dhul.house_unit_id = %s_unit_usage.house_unit_id
`, variable.Name, variable.Name)
			}
		}
	}
	log.Printf("PARAMS: %+v", params)
	for _, param := range params {
		log.Printf("PARAM: %+v", param)
		log.Printf("PARAM TYPE: %+v", reflect.TypeOf(param))
	}
	var res models.UnitsWithVariables
	sql = sql + middleSql + endSql
	err := db.Raw(sql, params).Scan(&res).Error
	if err != nil {
		return demand, dberror.DbErrorFromPq(err)
	}
	demand.Units = res
	return demand, nil
}

func (d *demandDao) SelectDemandById(db *gorm.DB, demandId uuid.UUID, houseId uuid.UUID) (models.Demand, error) {
	filterOptions := &DemandFilter{
		Id:      &demandId,
		HouseId: houseId,
	}
	demands, err := d.SelectDemands(db, filterOptions)
	if err != nil {
		return models.Demand{}, err
	}
	if len(demands) == 0 {
		return models.Demand{}, dberror.DbErrorFromPq(sql.ErrNoRows)
	}
	return demands[0], nil
}

func (d *demandDao) SelectDemands(db *gorm.DB, filterOptions *DemandFilter) ([]models.Demand, error) {
	sql := `
	SELECT
		d.id,
		d.name,
		d.description,
		d.start_date ,
		d.end_date ,
		d.period_days ,
		jsonb_build_object(
			'expression', de.expression,
			'variables', de.variables
		) AS expression,
		dhul.units
	FROM thk.demand d 
	LEFT JOIN (
		SELECT
			de.demand_id,
			de.expression,
			json_agg(
				json_build_object(
					'name', dv.name,
					'type', dv.type,
					'details', dv.details
				)
			) AS variables
		FROM thk.demand_expression de
		LEFT JOIN thk.demand_variable_lnk dvl
			ON de.demand_id = dvl.demand_id
		LEFT JOIN thk.demand_variable dv
			ON dvl.demand_variable_name = dv.name
		GROUP BY de.demand_id 
	) AS de
		ON d.id = de.demand_id
	LEFT JOIN (
		SELECT
			dhul.demand_id,
			json_agg(
				json_build_object(
					'id', dhul.house_unit_id 
				)
			) AS units	
			FROM thk.demand_house_unit_lnk dhul
			GROUP BY dhul.demand_id
	) dhul
		ON d.id = dhul.demand_id
	WHERE 1=1`
	params := map[string]interface{}{}
	if filterOptions != nil {
		sql += "\nAND d.house_id = @houseId"
		params["houseId"] = filterOptions.HouseId
		if filterOptions.Id != nil {
			sql += "\nAND d.id = @id"
			params["id"] = *filterOptions.Id
		}
		if filterOptions.Name != nil {
			sql += "\nAND d.name = @name"
			params["name"] = *filterOptions.Name
		}
		if filterOptions.Description != nil {
			sql += "\nAND d.description = @description"
			params["description"] = *filterOptions.Description
		}
		if filterOptions.StartDate != nil {
			sql += fmt.Sprintf("\nAND d.start_date %s @startDate", filterOptions.StartDate.Operator.String())
			params["startDate"] = filterOptions.StartDate.Date
		}
		if filterOptions.EndDate != nil {
			sql += fmt.Sprintf("\nAND d.end_date %s @endDate", filterOptions.EndDate.Operator.String())
			params["endDate"] = filterOptions.EndDate.Date
		}
		if filterOptions.PeriodDays != nil {
			sql += fmt.Sprintf("\nAND d.period_days %s @periodDays", filterOptions.PeriodDays.Operator.String())
			params["periodDays"] = filterOptions.PeriodDays.Days
		}
	}
	var res []models.Demand
	var q *gorm.DB
	if len(params) == 0 {
		q = db.Raw(sql)
	} else {
		q = db.Raw(sql, params)
	}
	err := q.Scan(&res).Error
	if err != nil {
		return nil, dberror.DbErrorFromPq(err)
	}
	return res, nil
}

func (d *demandDao) InsertDemand(db *gorm.DB, demand models.Demand, houseId uuid.UUID) (uuid.UUID, error) {
	sql := `
	INSERT INTO thk.demand (
		house_id,
		name,
		description,
		start_date,
		end_date,
		period_days
	) VALUES (
		@houseId,
		@name,
		@description,
		@startDate,
		@endDate,
		@periodDays
	) RETURNING id
	`
	var res string
	params := map[string]interface{}{
		"houseId":     houseId,
		"name":        demand.Name,
		"description": demand.Description,
		"startDate":   demand.StartDate,
		"endDate":     demand.EndDate,
		"periodDays":  demand.PeriodDays,
	}
	err := db.Raw(sql, params).Scan(&res).Error
	if err != nil {
		return uuid.UUID{}, dberror.DbErrorFromPq(err)
	}
	resUuid, err := uuid.Parse(res)
	if err != nil {
		return uuid.UUID{}, err
	}
	return resUuid, nil
}

func (d *demandDao) InsertDemandExpression(db *gorm.DB, demandId uuid.UUID, expression models.DemandExpression) error {
	sql := `
	INSERT INTO thk.demand_expression (
		demand_id,
		expression
	) VALUES (
		@demandId,
		@expression
	)
	`
	params := map[string]interface{}{
		"demandId":   demandId,
		"expression": expression.Expression,
	}
	err := db.Exec(sql, params).Error
	if err != nil {
		return dberror.DbErrorFromPq(err)
	}
	return nil
}

func (d *demandDao) InsertDemandConstant(db *gorm.DB, name string, value float64) error {
	sql := `
	INSERT INTO thk.demand_variable (
		name,
		type,
		details
	) VALUES (
		@name,
		@type,
		@details
	)
	`
	detailsMap := map[string]interface{}{
		"value": value,
	}
	details, err := json.Marshal(detailsMap)
	if err != nil {
		return err
	}
	params := map[string]interface{}{
		"name":    name,
		"type":    models.DemandVariableTypeConstant.String(),
		"details": details,
	}
	err = db.Exec(sql, params).Error
	if err != nil {
		return dberror.DbErrorFromPq(err)
	}
	return nil
}

func (d *demandDao) InsertDemandHouseUnitLink(db *gorm.DB, demandId uuid.UUID, houseUnitId uuid.UUID) error {
	sql := `
	INSERT INTO thk.demand_house_unit_link (
		demand_id,
		house_unit_id
	) VALUES (
		@demandId,
		@houseUnitId
	)
	`
	params := map[string]interface{}{
		"demandId":    demandId,
		"houseUnitId": houseUnitId,
	}
	err := db.Exec(sql, params).Error
	if err != nil {
		return dberror.DbErrorFromPq(err)
	}
	return nil
}

func (d *demandDao) InsertDemandRun(db *gorm.DB, run models.DemandRun) (uuid.UUID, error) {
	sql := `
	INSERT INTO thk.demand_run (
		demand_id,
		status,
	) VALUES (
		@demandId,
		@status
	) RETURNING id
	`
	var res string
	params := map[string]interface{}{
		"demandId": run.DemandId,
		"status":   run.Status,
	}
	err := db.Raw(sql, params).Scan(&res).Error
	if err != nil {
		return uuid.UUID{}, dberror.DbErrorFromPq(err)
	}
	resUuid, err := uuid.Parse(res)
	if err != nil {
		return uuid.UUID{}, err
	}
	return resUuid, nil
}

func (d *demandDao) InsertDemandVariableLink(db *gorm.DB, demandId uuid.UUID, demandVariableName string) error {
	sql := `
	INSERT INTO thk.demand_variable_lnk (demand_id, demand_variable_name)
	VALUES (@demandId, @demandVariableName)
	`
	params := map[string]interface{}{
		"demandId":         demandId,
		"demandVariableId": demandVariableName,
	}
	err := db.Exec(sql, params).Error
	if err != nil {
		return dberror.DbErrorFromPq(err)
	}
	return nil
}

func NewDemandDAO() Demand {
	return &demandDao{}
}

type DemandFilter struct {
	HouseId     uuid.UUID
	Id          *uuid.UUID
	Name        *string
	Description *string
	StartDate   *struct {
		Operator Operator
		Date     time.Time
	}
	EndDate *struct {
		Operator Operator
		Date     time.Time
	}
	PeriodDays *struct {
		Operator Operator
		Days     int
	}
}

type Operator int

const (
	Equal Operator = iota
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)

var operatorStr = map[Operator]string{
	Equal:              "=",
	NotEqual:           "!=",
	GreaterThan:        ">",
	GreaterThanOrEqual: ">=",
	LessThan:           "<",
	LessThanOrEqual:    "<=",
}

func (o Operator) String() string {
	return operatorStr[o]
}
