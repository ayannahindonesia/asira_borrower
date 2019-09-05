package models

import (
	"asira_borrower/asira"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jinzhu/gorm"
)

type (
	BaseModel struct {
		ID          uint64    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key,column:id"`
		CreatedTime time.Time `json:"created_time" gorm:"column:created_time" sql:"DEFAULT:current_timestamp"`
		UpdatedTime time.Time `json:"updated_time" gorm:"column:updated_time" sql:"DEFAULT:current_timestamp"`
	}

	DBFunc func(tx *gorm.DB) error

	PagedSearchResult struct {
		TotalData   int         `json:"total_data"`   // matched datas
		Rows        int         `json:"rows"`         // shown datas per page
		CurrentPage int         `json:"current_page"` // current page
		LastPage    int         `json:"last_page"`
		From        int         `json:"from"` // offset, starting index of data shown in current page
		To          int         `json:"to"`   // last index of data shown in current page
		Data        interface{} `json:"data"`
	}
	SearchResult struct {
		TotalData int         `json:"total_data"` // matched datas
		Data      interface{} `json:"data"`
	}
)

// helper for inserting data using gorm.DB functions
func WithinTransaction(fn DBFunc) (err error) {
	tx := asira.App.DB.Begin()
	defer tx.Commit()
	err = fn(tx)

	return err
}

// inserts a row into db.
func Create(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if !asira.App.DB.NewRecord(i) {
			return err
		}
		if err = tx.Create(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Update row in db.
func Save(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if asira.App.DB.NewRecord(i) {
			return err
		}
		if err = tx.Save(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Delete row in db.
func Delete(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if err = tx.Delete(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Find by id.
func FindbyID(i interface{}, id int) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Last(i, id).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

func FilterSearchSingle(i interface{}, filter interface{}) (err error) {
	db := asira.App.DB

	// filtering
	refFilter := reflect.ValueOf(filter).Elem()
	refType := refFilter.Type()
	for x := 0; x < refFilter.NumField(); x++ {
		field := refFilter.Field(x)
		if field.Interface() != "" {
			db = db.Where(fmt.Sprintf("%s = ?", refType.Field(x).Tag.Get("json")), field.Interface())
		}
	}

	if err = db.Last(i).Error; err != nil {
		db.Rollback()
		return err
	}

	return nil
}

func PagedFilterSearch(i interface{}, page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	if page <= 0 {
		page = 1
	}

	if rows <= 0 {
		rows = 25 // default row is 25 per page
	}

	db := asira.App.DB

	// filtering
	refFilter := reflect.ValueOf(filter).Elem()
	refType := refFilter.Type()
	for x := 0; x < refFilter.NumField(); x++ {
		field := refFilter.Field(x)
		if field.Interface() != "" {
			db = db.Where(fmt.Sprintf("%s = ?", refType.Field(x).Tag.Get("json")), field.Interface())
		}
	}

	// ordering and sorting
	if orderby != "" {
		orders := strings.Split(orderby, ",")
		sort := strings.Split(sort, ",")

		for k, v := range orders {
			e := v
			if len(sort) > k {
				value := sort[k]
				if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
					e = v + " " + strings.ToUpper(value)
				}
			}
			db = db.Order(e)
		}
	}

	tempDB := db
	var (
		total_rows int
		lastPage   int = 1 // default 1
	)

	tempDB.Find(i).Count(&total_rows)

	offset := (page * rows) - rows
	lastPage = int(math.Ceil(float64(total_rows) / float64(rows)))

	db.Limit(rows).Offset(offset).Find(i)

	result = PagedSearchResult{
		TotalData:   total_rows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        &i,
	}

	return result, err
}

func KafkaSubmitModel(i interface{}, model string) (err error) {
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics.produces", asira.App.ENV))

	var payload interface{}
	payload = kafkaPayloadBuilder(i, model)

	jMarshal, _ := json.Marshal(payload)

	kafkaProducer, err := sarama.NewAsyncProducer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topics["for_lender"].(string),
		Value: sarama.StringEncoder(strings.TrimSuffix(model, "_delete") + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", topics["for_lender"].(string))
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", topics["for_lender"].(string), err)
	}

	return nil
}

func kafkaPayloadBuilder(i interface{}, model string) (payload interface{}) {
	switch model {
	default:
		if strings.HasSuffix(model, "_delete") {
			type ModelDelete struct {
				ID     float64 `json:"id"`
				Model  string  `json:"model"`
				Delete bool    `json:"delete"`
			}
			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(i)
			json.Unmarshal(inrec, &inInterface)
			if modelID, ok := inInterface["id"].(float64); ok {
				payload = ModelDelete{
					ID:     modelID,
					Model:  strings.TrimSuffix(model, "_delete"),
					Delete: true,
				}
			}
		} else {
			payload = i
		}
		break
	}

	return payload
}

func FilterSearch(i interface{}, filter interface{}) (result SearchResult, err error) {
	db := asira.App.DB

	// filtering
	refFilter := reflect.ValueOf(filter).Elem()
	refType := refFilter.Type()
	for x := 0; x < refFilter.NumField(); x++ {
		field := refFilter.Field(x)
		if field.Interface() != "" {
			switch refType.Field(x).Tag.Get("condition") {
			case "OR":
				var e []string
				for _, filter := range field.Interface().([]string) {
					e = append(e, refType.Field(x).Tag.Get("json")+" = '"+filter+"' ")
				}
				db = db.Where(strings.Join(e, " OR "))
			}
		}
	}

	var total_rows int
	db.Find(i).Count(&total_rows)

	result = SearchResult{
		TotalData: total_rows,
		Data:      &i,
	}

	return result, err
}
