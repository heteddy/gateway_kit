// @Author : detaohe
// @File   : flow
// @Description:
// @Date   : 2022/10/16 19:54

package dao

import (
	"context"
	"errors"
	"gateway_kit/config"
	"gateway_kit/util/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

//type FlowSumDao struct {
//	mongodb.Dao
//}
//
//type FlowSumEntity struct {
//	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
//	Category  string              `json:"category" bson:"category"` // gateway 或者 service
//	Name      string              `json:"name" bson:"name"`         // http服务名词
//	Count     int64               `json:"count" bson:"count"`       //
//	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
//	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
//	DeletedAt int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
//}
//
//func NewFlowSumDao() *FlowSumDao {
//	indices := make(map[string]mongo.IndexModel)
//	idxSvcName := "idx_category_name"
//	indexBackground := true
//	unique := true
//
//	indices[idxSvcName] = mongo.IndexModel{
//		Keys: bson.D{{"name", 1}, {"category", 1}},
//		Options: &options.IndexOptions{
//			Name:       &idxSvcName,
//			Background: &indexBackground,
//			Unique:     &unique,
//		},
//	}
//	return &FlowSumDao{
//		Dao: mongodb.Dao{
//			Client:        config.MongoEngine,
//			Table:         config.All.Tables.Flow,
//			IndexParamMap: indices,
//		},
//	}
//}
//
//func (engine *FlowSumDao) Insert(ctx context.Context, entity *FlowSumEntity) (*FlowSumEntity, error) {
//	entity.CreatedAt = time.Now()
//	entity.UpdatedAt = time.Now()
//	result, err := engine.Collection().InsertOne(ctx, entity)
//	if err != nil {
//		return entity, err
//	} else {
//		if objID, ok := result.InsertedID.(primitive.ObjectID); ok {
//			entity.ID = &objID
//			return entity, nil
//		}
//		return entity, errors.New("mongo _id类型转换错误")
//	}
//}
//
//func (engine *FlowSumDao) Update(ctx context.Context, entity *FlowSumEntity) (*FlowSumEntity, error) {
//	entity.UpdatedAt = time.Now()
//	if entity.Count == 0 {
//		return entity, nil
//	}
//	//opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
//	//var cursor *mongo.Cursor
//	//var err error
//	count, err := engine.Collection().CountDocuments(ctx, bson.D{{"category", entity.Category}, {"name", entity.Name}})
//	if err != nil {
//		return nil, err
//	} else {
//		if count == 0 {
//			return engine.Insert(ctx, entity)
//		} else {
//			ret, err := engine.Collection().UpdateMany(ctx, bson.D{{"category", entity.Category}, {"name", entity.Name}}, bson.M{"$inc": bson.M{"count": entity.Count}})
//			if err != nil {
//				return nil, err
//			} else {
//				config.Logger.Info("update sum count", zap.Int64("modified count", ret.ModifiedCount))
//				return entity, nil
//			}
//		}
//	}
//}
//
//func (engine *FlowSumDao) All(ctx context.Context) (entities []*FlowSumEntity, err error) {
//	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
//	var cursor *mongo.Cursor
//	cursor, err = engine.Collection().Find(ctx, bson.M{"deleted_at": 0}, opt)
//	if err != nil {
//		return
//	} else {
//		err = cursor.All(ctx, &entities)
//		return
//	}
//}
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

type ReqHourDao struct {
	mongodb.Dao
}

func NewReqHourDao() *ReqHourDao {
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
	return &ReqHourDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.RequestHour,
			IndexParamMap: indices,
		},
	}
}

func (engine *ReqHourDao) AllInMonth(ctx context.Context) (entities []*ReqHourEntity, err error) {
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

func (engine *ReqHourDao) Insert(ctx context.Context, entity *ReqHourEntity) (*ReqHourEntity, error) {
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

func (engine *ReqHourDao) Update(ctx context.Context, entity *ReqHourEntity) (*ReqHourEntity, error) {
	if entity.Count == 0 {
		return entity, nil
	}
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
	return &ServiceHourDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.ServiceHour,
			IndexParamMap: indices,
		},
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
