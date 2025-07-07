package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DemandVariableType int

const (
	DemandVariableTypeConstant DemandVariableType = iota
	DemandVariableTypeArea
	DemandVariableTypeOwnership
	DemandVariableTypeUsage
	DemandVariableTypeSumUsage
)

var demandVariableTypeName = map[DemandVariableType]string{
	DemandVariableTypeConstant:  "constant",
	DemandVariableTypeArea:      "area",
	DemandVariableTypeOwnership: "ownership",
	DemandVariableTypeUsage:     "usage",
	DemandVariableTypeSumUsage:  "sum_usage",
}

func (t DemandVariableType) String() string {
	return demandVariableTypeName[t]
}

func (t DemandVariableType) FromString(s string) (DemandVariableType, error) {
	for k, v := range demandVariableTypeName {
		if v == s {
			return k, nil
		}
	}
	return DemandVariableTypeConstant, fmt.Errorf("invalid demand variable type: %s", s)
}

type DemandVariableWithValue struct {
	Name    string                 `gorm:"column:name" json:"name"`
	Type    string                 `gorm:"column:type" json:"type"`
	Details map[string]interface{} `gorm:"column:details" json:"details"`
	Value   float64                `gorm:"column:value" json:"value"`
}

func (p *DemandVariableWithValue) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type DemandVariableWithValues []DemandVariableWithValue

func (p *DemandVariableWithValues) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type DemandVariable struct {
	Name    string                 `gorm:"column:name" json:"name"`
	Type    string                 `gorm:"column:type" json:"type"`
	Details map[string]interface{} `gorm:"column:details" json:"details"`
}

func (p *DemandVariable) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type DemandVariables []DemandVariable

func (p *DemandVariables) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type DemandExpression struct {
	Expression string          `gorm:"column:expression" json:"expression"`
	Variables  DemandVariables `gorm:"column:variables;type:DemandVariables" json:"variables"`
}

func (p *DemandExpression) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type UnitWithVariables struct {
	Id        uuid.UUID                 `gorm:"column:id;type:uuid" json:"id"`
	Variables *DemandVariableWithValues `gorm:"column:variables;type:DemandVariableWithValues" json:"variables"`
}

func (p *UnitWithVariables) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), &p.Id)
}

type UnitsWithVariables []UnitWithVariables

func (p *UnitsWithVariables) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type Demand struct {
	Id          *uuid.UUID         `gorm:"column:id;type:uuid" json:"id"`
	Name        string             `gorm:"column:name" json:"name"`
	Description string             `gorm:"column:description" json:"description"`
	StartDate   time.Time          `gorm:"column:start_date;type:time" json:"start_date"`
	EndDate     time.Time          `gorm:"column:end_date;type:time" json:"end_date"`
	PeriodDays  int                `gorm:"column:period_days" json:"period_days"`
	Expression  DemandExpression   `gorm:"column:expression;type:DemandExpression" json:"expression"`
	Units       UnitsWithVariables `gorm:"column:units;type:UnitIds" json:"units"`
}

func (p *Demand) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}

type DemandRunStatus int

const (
	DemandRunStatusPending DemandRunStatus = iota
	DemandRunStatusRunning
	DemandRunStatusFinished
	DemandRunStatusFailed
)

var demandRunStatusName = map[DemandRunStatus]string{
	DemandRunStatusPending:  "pending",
	DemandRunStatusRunning:  "running",
	DemandRunStatusFinished: "finished",
	DemandRunStatusFailed:   "failed",
}

func (t DemandRunStatus) String() string {
	return demandRunStatusName[t]
}

func (t DemandRunStatus) FromString(s string) (DemandRunStatus, error) {
	for k, v := range demandRunStatusName {
		if v == s {
			return k, nil
		}
	}
	return DemandRunStatusPending, fmt.Errorf("invalid demand run status: %s", s)
}

type DemandRun struct {
	Id         *uuid.UUID `gorm:"column:id;type:uuid" json:"id"`
	DemandId   uuid.UUID  `gorm:"column:demand_id;type:uuid" json:"demand_id"`
	Status     string     `gorm:"column:status" json:"status"`
	FromDate   time.Time  `gorm:"column:from_date;type:time" json:"from_date"`
	ToDate     time.Time  `gorm:"column:to_date;type:time" json:"to_date"`
	CreateDate time.Time  `gorm:"column:create_date;type:time" json:"create_date"`
	UpdateDate time.Time  `gorm:"column:update_date;type:time" json:"update_date"`
}

func (p *DemandRun) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), p)
}
