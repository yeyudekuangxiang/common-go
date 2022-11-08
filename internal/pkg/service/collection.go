package service

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"time"
)

/*
desc: 收藏
*/

type (
	CollectionService interface {
		TopicCollections(openId string, limit, offset int) ([]*entity.Topic, int64, error) //我的收藏
		Collection(objId int64, objType int, openId string) error                          //收藏
		CancelCollection(objId int64, objType int, openId string) error                    //取消收藏
		Collections(openId string, objType, limit, offset int) []int64                     //收藏数据
		CollectionV2(objId int64, objType int, openId string) (bool, error)                //收藏
		FindOneByTopic(topicId int64, openId string) (*entity.Collection, error)
	}

	defaultCollectionService struct {
		ctx             *mioContext.MioContext
		collectionModel repository.CollectionModel
		topicModel      repository.TopicModel
	}
)

func (d defaultCollectionService) FindOneByTopic(topicId int64, openId string) (*entity.Collection, error) {
	ojb, err := d.collectionModel.FindOneByObj(topicId, 0, openId)
	if err != nil {
		return nil, err
	}
	return ojb, nil
}

func (d defaultCollectionService) Collections(openId string, objType, limit, offset int) []int64 {
	return d.getCollections(objType, openId, limit, offset)
}

func (d defaultCollectionService) TopicCollections(openId string, limit, offset int) ([]*entity.Topic, int64, error) {
	//cond type; get ids
	ids := d.getCollections(0, openId, 0, 0)
	//查找文章
	if len(ids) <= 0 {
		return []*entity.Topic{}, 0, nil
	}

	list, total, err := d.topicModel.GetTopicList(repository.GetTopicPageListBy{
		Ids:    ids,
		Status: 3,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d defaultCollectionService) Collection(objId int64, objType int, openId string) error {
	result, err := d.collectionModel.FindOneByObj(objId, objType, openId)
	if err != nil {
		if err == entity.ErrNotFount {
			//insert
			data := &entity.Collection{
				ObjId:     objId,
				ObjType:   objType,
				Status:    1,
				OpenId:    openId,
				CreatedAt: time.Now(),
			}
			_, err = d.collectionModel.Insert(data)
			if err != nil {
				return err
			}
			err = d.incrTopicCollections(objId)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	//update
	if result.Status == 2 {
		result.Status = 1
		err = d.collectionModel.Update(result)
		if err != nil {
			return err
		}

		err = d.incrTopicCollections(objId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d defaultCollectionService) CollectionV2(objId int64, objType int, openId string) (bool, error) {
	var isFirst bool
	err := d.ctx.Transaction(func(ctx *mioContext.MioContext) error {
		collectionModel := repository.NewCollectionRepository(ctx)
		topicModel := repository.NewTopicModel(ctx)

		result, err := collectionModel.FindOneByObj(objId, objType, openId)
		if err != nil {
			if err == entity.ErrNotFount {
				isFirst = true
				//insert
				data := &entity.Collection{
					ObjId:     objId,
					ObjType:   objType,
					Status:    1,
					OpenId:    openId,
					CreatedAt: time.Now(),
				}
				_, err = collectionModel.Insert(data)
				if err != nil {
					return err
				}
				err = topicModel.ChangeTopicCollectionCount(objId, "collection_count", 1)
				if err != nil {
					return err
				}
				return nil
			}
			return err
		}

		if result.Status == 2 {
			result.Status = 1
			err = collectionModel.Update(result)
			if err != nil {
				return err
			}

			err = topicModel.UpdateColumn(objId, "collection_count", 1)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return isFirst, err
	}

	return isFirst, nil
}

func (d defaultCollectionService) CancelCollection(objId int64, objType int, openId string) error {
	err := d.ctx.Transaction(func(ctx *mioContext.MioContext) error {
		collectionModel := repository.NewCollectionRepository(ctx)
		topicModel := repository.NewTopicModel(ctx)

		result, err := collectionModel.FindOneByObj(objId, objType, openId)
		if err != nil {
			return err
		}
		if result.Status == 2 {
			return nil
		}
		result.Status = 2
		err = collectionModel.Update(result)
		if err != nil {
			return err
		}
		err = topicModel.ChangeTopicCollectionCount(objId, "collection_count", -1)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (d defaultCollectionService) getCollections(objType int, openId string, limit, offset int) []int64 {
	var objIds []int64
	collections, _, err := d.collectionModel.FindAllByOpenId(objType, openId, limit, offset)
	if err != nil {
		return objIds
	}

	for _, collection := range collections {
		objIds = append(objIds, collection.ObjId)
	}
	return objIds
}

func (d defaultCollectionService) incrTopicCollections(objId int64) error {
	return d.topicModel.ChangeTopicCollectionCount(objId, "collection_count", 1)
}

func (d defaultCollectionService) decrTopicCollections(objId int64) error {
	return d.topicModel.ChangeTopicCollectionCount(objId, "collection_count", -1)
}

func NewCollectionService(ctx *mioContext.MioContext) CollectionService {
	return &defaultCollectionService{
		ctx:             ctx,
		collectionModel: repository.NewCollectionRepository(ctx),
		topicModel:      repository.NewTopicModel(ctx),
	}
}
