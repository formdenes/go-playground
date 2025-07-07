package main

import (
	"fmt"
	"log"
	"playground/thk/config"
	"playground/thk/dao"
	"playground/thk/models"
	"time"

	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/showa-93/go-mask"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	houseIdString string = "6a6ceedb-2155-46f9-9101-8634398b9b7f"
	demandId      string = "0ce9bd28-e692-40bf-bb0d-9a962bf3ef5e"
)

var demandDao dao.Demand

func GetDatabase(cfg config.Config) *gorm.DB {
	log.Printf("Initializing DB connection")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User)

	// db, err := dbx.MustOpen("postgres", psqlInfo)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{TranslateError: false})

	if err != nil {
		log.Printf("Could not connect to db, error=%v", err)
	}
	log.Printf("Connected to DB")

	return db

}

func main() {
	var cfg config.Config

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		log.Printf("failed to get config, error: %v", err)
	}
	maskedConfig, _ := mask.Mask(cfg)
	log.Printf("Start foresight data service with config %+v", maskedConfig)

	db := GetDatabase(cfg)
	log.Printf("DB connection established %+v", db)
	log.Println("")
	log.Println("")
	log.Println("")

	demandDao = dao.NewDemandDAO()

	demand, err := GetDemand(db)
	if err != nil {
		log.Printf("Could not insert demand, error=%v", err)
	}
	demandRun, err := GetDemandRun(db)
	if err != nil {
		log.Printf("Could not get demand run, error=%v", err)
	}
	// log.Printf("Inserted demand %+v", demand)
	log.Printf("Inserted demand run %+v", demandRun)
	log.Printf("")
	log.Printf("")
	log.Printf("")
	variableValues, err := GetDemandVariableValues(db, demand, demandRun)
	if err != nil {
		log.Printf("Could not get variable values, error=%v", err)
	}
	log.Printf("Variable values %+v", variableValues)
}

func GetDemandRun(db *gorm.DB) (models.DemandRun, error) {
	demandId := uuid.MustParse(demandId)
	return demandDao.SelectLatestDemandRun(db, demandId)
}

func GetDemandVariableValues(db *gorm.DB, demand models.Demand, demandRun models.DemandRun) (models.Demand, error) {
	return demandDao.SelectVariableValues(db, demand, demandRun)
}

func GetDemand(db *gorm.DB) (models.Demand, error) {
	houseId := uuid.MustParse(houseIdString)
	demandId := uuid.MustParse(demandId)

	return demandDao.SelectDemandById(db, demandId, houseId)
}

func InsertDemand(db *gorm.DB) (uuid.UUID, error) {
	houseId := uuid.MustParse(houseIdString)
	demand := models.Demand{
		Name:        "test demand",
		Description: "test demand description",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 24),
		PeriodDays:  30,
		Expression: models.DemandExpression{
			Expression: "ter*(tul)/gfogy",
			Variables: models.DemandVariables{
				models.DemandVariable{
					Name:    "ter",
					Type:    models.DemandVariableTypeArea.String(),
					Details: map[string]interface{}{},
				},
				models.DemandVariable{
					Name:    "tul",
					Type:    models.DemandVariableTypeOwnership.String(),
					Details: map[string]interface{}{},
				},
				models.DemandVariable{
					Name:    "gfogy",
					Type:    models.DemandVariableTypeUsage.String(),
					Details: map[string]interface{}{},
				},
			},
		},
		Units: []models.UnitWithVariables{},
	}

	demandId, err := demandDao.InsertDemand(db, demand, houseId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return demandId, nil
}
