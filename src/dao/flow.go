// @Author : detaohe
// @File   : flow
// @Description:
// @Date   : 2022/10/16 19:54

package dao

import (
	"context"
	"errors"
	"fmt"
	"gateway_kit/config"
	"gateway_kit/util/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type ReqHourEntity struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Service   string              `json:"service" bson:"service"`
	Path      string              `json:"path" bson:"path"`     // http服务名词
	Method    string              `json:"method" bson:"method"` // http服务名词
	Count     int64               `json:"count" bson:"count"`   //
	Hour      int                 `json:"hour " bson:"hour"`
	Date      string              `json:"date" bson:"date"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
}

type RequestHourDao struct {
	mongodb.Dao
}

func NewRequestHourDao() *RequestHourDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_service_path_method_hour_date"
	indexBackground := true
	unique := true

	indices[idxSvcName] = mongo.IndexModel{
		Keys: bson.D{{"service", 1}, {"path", 1}, {"method", 1}, {"hour", 1}, {"date", 1}},
		Options: &options.IndexOptions{
			Name:       &idxSvcName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	idxTTLName := "idx_ttl_created_at"
	indices[idxTTLName] = mongo.IndexModel{
		Keys: bson.D{{"created_at", 1}, {"expireAfterSeconds", 7776000}}, //90*24*60*60
		Options: &options.IndexOptions{
			Name:       &idxTTLName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	return &RequestHourDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.RequestHour,
			IndexParamMap: indices,
		},
	}
}

func (engine *RequestHourDao) AllInMonth(ctx context.Context) (entities []*ReqHourEntity, err error) {
	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
	var cursor *mongo.Cursor
	monthAgo := time.Now().AddDate(0, -1, 0)
	cursor, err = engine.Collection().Find(ctx, bson.D{{"updated_at", 0}, {"created_at", bson.M{"$gt": monthAgo}}}, opt)
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}

func (engine *RequestHourDao) Insert(ctx context.Context, entity *ReqHourEntity) (*ReqHourEntity, error) {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	result, err := engine.Collection().InsertOne(ctx, entity)
	if err != nil {
		return entity, err
	} else {
		if objID, ok := result.InsertedID.(primitive.ObjectID); ok {
			entity.ID = &objID
			return entity, nil
		}
		return entity, errors.New("mongo _id类型转换错误")
	}
}

func (engine *RequestHourDao) Update(ctx context.Context, entity *ReqHourEntity) (*ReqHourEntity, error) {
	if entity.Count == 0 {
		return entity, nil
	}
	// note 是否需要事务或分布式锁
	//session, err := engine.Client.StartSession()
	//if err != nil {
	//	return nil, err
	//}
	//defer session.EndSession(ctx)
	//sessionCtx := mongo.NewSessionContext(ctx, session)
	//if err = session.StartTransaction(); err != nil {
	//	return nil, err
	//}
	var err error
	entity.UpdatedAt = time.Now()
	var count int64
	count, err = engine.Collection().CountDocuments(ctx,
		bson.D{{"service", entity.Service},
			{"path", entity.Path},
			{"method", entity.Method},
			{"hour", entity.Hour},
			{"date", entity.Date}})
	if err != nil {
		return nil, err
	} else {
		//defer func() {
		//	if err2 := session.CommitTransaction(context.Background()); err2 != nil {
		//	} else {
		//		config.Logger.Info("commit success")
		//	}
		//}()
		if count == 0 {
			return engine.Insert(ctx, entity)
		} else {
			ret, err := engine.Collection().UpdateMany(ctx,
				bson.D{{"service", entity.Service},
					{"path", entity.Path},
					{"method", entity.Method},
					{"hour", entity.Hour},
					{"date", entity.Date}},
				bson.M{"$inc": bson.M{"count": entity.Count}})
			if err != nil {
				return nil, err
			} else {
				config.Logger.Info("update sum count", zap.Int64("modified count", ret.ModifiedCount))
				return entity, nil
			}
		}
	}
}

func (engine *RequestHourDao) GetServiceRequestsDetail(ctx context.Context, service string, from, end time.Time) ([]*ReqHourEntity, error) {
	opts := options.Find().SetSort(bson.M{"updated_at": 1})
	cursor, err := engine.Collection().Find(ctx,
		bson.D{
			{"service", service},
			{"updated_at",
				bson.D{
					{"$gte", time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)},
					{"$lt", time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)},
				},
			},
		},
		opts)

	entities := make([]*ReqHourEntity, 0)
	if err = cursor.All(ctx, &entities); err != nil {
		config.Logger.Error("error of decode entities", zap.Error(err))
		return nil, err
	} else {
		for _, e := range entities {
			config.Logger.Info("entity", zap.String("name", e.Service), zap.String("uri", e.Path), zap.Int("hour", e.Hour), zap.Int64("count", e.Count))
		}
	}
	return entities, nil
}

func (engine *RequestHourDao) GetReqSum(ctx context.Context, service, uri, method string, from, end time.Time) {
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	groupStage := bson.D{
		{"$match",
			bson.D{
				{"service", service},
				{"path", uri},
				{"method", method},
				{"updated_at",
					bson.D{
						{"$gte", time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)},
						{"$lt", time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
		{"$group", bson.D{
			{"_id", bson.M{"service": "$service", "path": "$path", "method": "$method"}},
			{"svc_sum",
				bson.D{
					{"$sum", "$count"},
				}},
		}},
	}
	if aggCursor, err := engine.Collection().Aggregate(ctx, mongo.Pipeline{groupStage}, opts); err != nil {
		config.Logger.Error("aggregate error", zap.Error(err))
	} else {
		var results []bson.M
		if err = aggCursor.All(ctx, &results); err != nil {
			config.Logger.Error("aggCursor load error", zap.Error(err))
		}
		for _, result := range results {
			fmt.Printf("category %v, count=%v\n", result["_id"], result["sum_count"])
		}
	}
}

// ServiceHourEntity 按照时间写入
type ServiceHourEntity struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Category  string              `json:"category" bson:"category,omitempty"` //gateway service uri
	Name      string              `json:"name" bson:"name"`                   // http服务名词
	Count     int64               `json:"count" bson:"count"`                 //
	Hour      int                 `json:"hour " bson:"hour"`
	Date      string              `json:"date" bson:"date"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
}

type ServiceHourDao struct {
	mongodb.Dao
}

func NewServiceHourDao() *ServiceHourDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_category_name_hour_date"
	indexBackground := true
	unique := true

	indices[idxSvcName] = mongo.IndexModel{
		Keys: bson.D{{"name", 1}, {"category", 1}, {"hour", 1}, {"date", 1}},
		Options: &options.IndexOptions{
			Name:       &idxSvcName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	idxTTLName := "idx_ttl_created_at"
	indices[idxTTLName] = mongo.IndexModel{
		Keys: bson.D{{"created_at", 1}, {"expireAfterSeconds", 7776000}}, //90*24*60*60
		Options: &options.IndexOptions{
			Name:       &idxTTLName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	return &ServiceHourDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.ServiceHour,
			IndexParamMap: indices,
		},
	}
}

func (engine *ServiceHourDao) GetDetail(ctx context.Context, category, name string, from, end time.Time) (entities []*ServiceHourEntity, err error) {

	opts := options.Find().SetSort(bson.M{"updated_at": 1})
	var cursor *mongo.Cursor
	cursor, err = engine.Collection().Find(ctx,
		bson.D{
			{"category", category},
			{"name", name},
			{"date",
				bson.D{
					{"$gte", from.Format("2006-01-02")},
					{"$lt", end.Format("2006-01-02")},
				},
			},
			//"$and": bson.D{
			//	{"updated_at", bson.M{"$gte": from}},
			//	{"updated_at", bson.M{"lt": end}},
			//},
		},
		opts)

	//entities := make([]*ServiceHourEntity, 0)
	if err = cursor.All(ctx, &entities); err != nil {
		config.Logger.Error("error of decode entities", zap.Error(err))
		return
	} else {
		//
		return
	}
}

type _ID struct {
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
}
type ServiceSumEntity struct {
	ID     _ID   `json:"_id" bson:"_id"`
	SvcSum int64 `json:"svc_sum" bson:"svc_sum"`
}

func (engine *ServiceHourDao) GetSum(ctx context.Context, from, end time.Time) ([]*ServiceSumEntity, error) {
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	matchStage := bson.D{
		{"$match",
			bson.D{
				//{"category", "service"},
				{"date",
					bson.D{
						{"$gte", from.Format("2006-01-02")},
						{"$lt", end.Format("2006-01-02")},
					},
				},
				//{"$and",
				//	bson.D{
				//		{"updated_at", bson.M{"$gte": from}},
				//		{"updated_at", bson.M{"$lt": end}},
				//	},
				//},
			},
		},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.M{"category": "$category", "name": "$name", "date": "$date"}},
			{"svc_sum", bson.D{
				{"$sum", "$count"},
			}},
		}},
	}
	if aggCursor, err := engine.Collection().Aggregate(ctx, mongo.Pipeline{matchStage, groupStage}, opts); err != nil {
		config.Logger.Error("aggregate error", zap.Error(err))
		return nil, err
	} else {
		//var results []bson.M
		var results []*ServiceSumEntity
		if err = aggCursor.All(ctx, &results); err != nil {
			config.Logger.Error("aggCursor load error", zap.Error(err))
			return nil, err
		}
		for _, result := range results {
			//fmt.Printf("category %v, count=%v\n", result["_id"], result[])
			fmt.Printf("result= %v\n", result)
		}

		return results, nil
	}
}

func (engine *ServiceHourDao) AllInMonth(ctx context.Context) (entities []*ServiceHourEntity, err error) {
	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
	var cursor *mongo.Cursor
	monthAgo := time.Now().AddDate(0, -1, 0)
	cursor, err = engine.Collection().Find(ctx, bson.D{{"updated_at", 0}, {"created_at", bson.M{"$gt": monthAgo}}}, opt)
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}

func (engine *ServiceHourDao) Insert(ctx context.Context, entity *ServiceHourEntity) (*ServiceHourEntity, error) {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	result, err := engine.Collection().InsertOne(ctx, entity)
	if err != nil {
		return entity, err
	} else {
		if objID, ok := result.InsertedID.(primitive.ObjectID); ok {
			entity.ID = &objID
			return entity, nil
		}
		return entity, errors.New("mongo _id类型转换错误")
	}
}

func (engine *ServiceHourDao) Update(ctx context.Context, entity *ServiceHourEntity) (*ServiceHourEntity, error) {
	if entity.Count == 0 {
		return entity, nil
	}
	var err error
	//session, err := engine.Client.StartSession()
	//if err != nil {
	//	return nil, err
	//}
	//defer session.EndSession(ctx)
	//sessionCtx := mongo.NewSessionContext(ctx, session)
	//if err = session.StartTransaction(); err != nil {
	//	return nil, err
	//}

	entity.UpdatedAt = time.Now()
	var count int64
	count, err = engine.Collection().CountDocuments(ctx,
		bson.D{{"category", entity.Category},
			{"name", entity.Name},
			{"hour", entity.Hour},
			{"date", entity.Date}})
	if err != nil {
		return nil, err
	} else {
		//defer func() {
		//	if err2 := session.CommitTransaction(context.Background()); err2 != nil {
		//	} else {
		//		config.Logger.Info("commit success")
		//	}
		//}()
		if count == 0 {
			return engine.Insert(ctx, entity)
		} else {
			ret, err := engine.Collection().UpdateMany(ctx,
				bson.D{{"category", entity.Category},
					{"name", entity.Name},
					{"hour", entity.Hour},
					{"date", entity.Date}},
				bson.M{"$inc": bson.M{"count": entity.Count}})
			if err != nil {
				return nil, err
			} else {
				config.Logger.Info("update sum count", zap.Int64("modified count", ret.ModifiedCount))
				return entity, nil
			}
		}
	}
}

// ServiceDayEntity 按照时间写入
type ServiceDayEntity struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Category  string              `json:"category" bson:"category,omitempty"` //gateway service uri
	Name      string              `json:"name" bson:"name"`                   // http服务名词
	Count     int64               `json:"count" bson:"count"`                 //
	Date      string              `json:"date" bson:"date"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
}

type ServiceDayDao struct {
	mongodb.Dao
}

func NewServiceDayDao() *ServiceDayDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_category_name_date"
	indexBackground := true
	unique := true

	indices[idxSvcName] = mongo.IndexModel{
		Keys: bson.D{{"category", 1}, {"name", 1}, {"date", 1}},
		Options: &options.IndexOptions{
			Name:       &idxSvcName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	return &ServiceDayDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.ServiceDay,
			IndexParamMap: indices,
		},
	}
}
func (engine *ServiceDayDao) InsertMany(ctx context.Context, entities []*ServiceDayEntity) error {
	inserts := make([]interface{}, 0, len(entities))
	for _, e := range entities {
		e.CreatedAt = time.Now()
		e.UpdatedAt = time.Now()
		inserts = append(inserts, e)
	}

	_, err := engine.Collection().InsertMany(ctx, inserts)
	return err
}

func (engine *ServiceDayDao) Count(ctx context.Context, category, name string, t time.Time) (int64, error) {
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := today.AddDate(0, 0, 1)
	return engine.Collection().CountDocuments(ctx,
		bson.D{
			{"category", category}, {"name", name},
			{"created_at",
				bson.D{
					{"$gte", today},
					{"$lt", tomorrow},
					//{"$gte", time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)},
					//{"$lt", time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)},
				}},
		})
}

func (engine *ServiceDayDao) GetDetail(ctx context.Context, category, name string, from, end time.Time) (entities []*ServiceDayEntity, err error) {
	opts := options.Find().SetSort(bson.M{"updated_at": 1})
	var cursor *mongo.Cursor
	var matchItems bson.D
	if category != "" {
		matchItems = append(matchItems, bson.E{"category", category})
	}
	if name != "" {
		matchItems = append(matchItems, bson.E{"name", name})
	}
	matchItems = append(matchItems, bson.E{"date",
		bson.D{
			{"$gte", from.Format("2006-01-02")},
			{"$lt", end.Format("2006-01-02")},
		},
	},
	)

	cursor, err = engine.Collection().Find(ctx,
		matchItems,
		opts)

	//entities := make([]*ServiceHourEntity, 0)
	if err = cursor.All(ctx, &entities); err != nil {
		config.Logger.Error("error of decode entities", zap.Error(err))
		return
	} else {
		//
		return
	}
}

// GetSum  定义返回值结构
func (engine *ServiceDayDao) GetSum(ctx context.Context, category, name string, from, end time.Time) ([]*ServiceSumEntity, error) {
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	var matchItems bson.D

	//matchStage := bson.D{
	//	{"$match",
	//		bson.D{
	//			{"category", category},
	//			{"category", name},
	//			{"updated_at", bson.D{{"$gte", from}, {"$lt", end}}},
	//			//{"$and",
	//			//	bson.D{
	//			//		{"updated_at", bson.M{"$gte": from}},
	//			//		{"updated_at", bson.M{"$lt": end}},
	//			//	},
	//			//},
	//		},
	//	},
	//}
	if category != "" {
		matchItems = append(matchItems, bson.E{"category", category})
	}
	if name != "" {
		matchItems = append(matchItems, bson.E{"name", name})
	}
	matchItems = append(matchItems, bson.E{"date",
		bson.D{
			{"$gte", from.Format("2006-01-02")},
			{"$lt", end.Format("2006-01-02")},
		},
	})
	matchStage := bson.D{
		{"$match", matchItems},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.M{"name": "$name"}},
			{"svc_sum", bson.D{
				{"$sum", "$count"},
			}},
		}},
	}
	if aggCursor, err := engine.Collection().Aggregate(ctx, mongo.Pipeline{matchStage, groupStage}, opts); err != nil {
		config.Logger.Error("aggregate error", zap.Error(err))
		return nil, err
	} else {
		var results []*ServiceSumEntity
		if err = aggCursor.All(ctx, &results); err != nil {
			config.Logger.Error("aggCursor load error", zap.Error(err))
			return nil, err
		}
		for _, result := range results {
			//fmt.Printf("category %v, count=%v\n", result["_id"], result[])
			fmt.Printf("result= %v\n", result)
		}

		return results, nil
	}
}

type FlowSummaryTask struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Date      string              `json:"date" bson:"date"`
	Status    string              `json:"status" bson:"status"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
}

type ServiceSummaryTaskDao struct {
	mongodb.Dao
}

func NewServiceSummaryDao() *ServiceSummaryTaskDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_date"
	indexBackground := true
	unique := true

	indices[idxSvcName] = mongo.IndexModel{
		Keys: bson.D{{"date", 1}},
		Options: &options.IndexOptions{
			Name:       &idxSvcName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	idxTTLName := "idx_ttl_created_at"
	indices[idxTTLName] = mongo.IndexModel{
		Keys: bson.D{{"created_at", 1}, {"expireAfterSeconds", 7776000}}, //90*24*60*60
		Options: &options.IndexOptions{
			Name:       &idxTTLName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	return &ServiceSummaryTaskDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.SummaryTask,
			IndexParamMap: indices,
		},
	}
}

func (engine *ServiceSummaryTaskDao) Existed(ctx context.Context, date string) (bool, error) {
	count, err := engine.Collection().CountDocuments(ctx, bson.D{{"date", date}})
	return count > 0, err
}

func (engine *ServiceSummaryTaskDao) Start(ctx context.Context, date string) (*FlowSummaryTask, error) {
	t := FlowSummaryTask{
		Date:   date,
		Status: "running",
	}
	ret, err := engine.Collection().InsertOne(ctx, &t)
	if objID, ok := ret.InsertedID.(primitive.ObjectID); ok {
		t.ID = &objID
		return &t, nil
	}
	return &t, err
}

func (engine *ServiceSummaryTaskDao) Complete(ctx context.Context, date string) error {
	ret := engine.Collection().FindOneAndUpdate(ctx, bson.D{{"date", date}, {"status", "running"}}, bson.M{"$set": bson.M{"status": "complete"}})
	return ret.Err()
}
