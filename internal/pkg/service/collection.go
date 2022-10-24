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
	}

	defaultCollectionService struct {
		collectionModel repository.CollectionModel
		topicModel      repository.TopicModel
	}
)

func (d defaultCollectionService) TopicCollections(openId string, limit, offset int) ([]*entity.Topic, int64, error) {
	//cond type; get ids
	ids := d.getCollections(0, openId, 0, 0)
	//查找文章
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
	result, err := d.collectionModel.FindOneByOjb(objId, objType, openId)
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
			return err
		}
		return err
	}
	//update
	if result.Status == 2 {
		result.Status = 1
		return d.collectionModel.Update(result)
	}
	return nil
}

func (d defaultCollectionService) CancelCollection(objId int64, objType int, openId string) error {
	result, err := d.collectionModel.FindOneByOjb(objId, objType, openId)
	if err != nil {
		return err
	}
	if result.Status == 2 {
		return nil
	}
	result.Status = 2
	return d.collectionModel.Update(result)
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

func NewCollectionService(ctx *mioContext.MioContext) CollectionService {
	return &defaultCollectionService{
		collectionModel: repository.NewCollectionRepository(ctx),
		topicModel:      repository.NewTopicRepository(ctx),
	}
}
