package service

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/randnull/banner-service/pkg/models"
)

type BannerRepostory interface {
	AddBaner(banner *models.Banner) (int, error)
	DeleteBanner(banner_id int) error
	GetBanner(tag_id int, feature_id int) (*models.Content, error)
	UpdateBanner(banner_id int, banner *models.UpdateBanner) error
	GetAllBanners(tag_id int, feature_id int, limit int, offset int) ([]*models.BannerDB, error)
}

type BannerService struct {
	repo BannerRepostory
	c	 *cache.Cache
}

func NewBannerSevice(repo BannerRepostory) *BannerService {
	return &BannerService{
		repo:	repo,
		c: 	 	cache.New(45 * time.Minute, 100 * time.Minute),
	}
}

func (service *BannerService) CreateBanner(banner *models.Banner) (int, error) {
	id, err := service.repo.AddBaner(banner)

	if err != nil {
		return -1, err
	}

	if banner.IsActive {
		for _, tag_id := range banner.TagIds {
			str_cache := fmt.Sprintf("tag_id%vfeat_id%v", tag_id, banner.FeatureId)

			service.c.Set(str_cache, &banner.Content, cache.DefaultExpiration)
		}
	}

	return id, nil
}

func (service *BannerService) DeleteBanner(banner_id int) error {
	err := service.repo.DeleteBanner(banner_id)

	if err != nil {
		return err
	}

	return nil
}


func (service *BannerService) GetBanner(tag_id int, feature_id int, use_last_version bool) (*models.Content, error) {
	str_cache := fmt.Sprintf("tag_id%vfeat_id%v", tag_id, feature_id)

	if !use_last_version {
		result, found := service.c.Get(str_cache)

		if found {			
			answer, _ := result.(*models.Content)

			return answer, nil
		}
	}

	content, err := service.repo.GetBanner(tag_id, feature_id)

	service.c.Set(str_cache, content, cache.DefaultExpiration)

	if err != nil {
		return nil, err
	}

	return content, nil
}


func (service *BannerService) UpdateBanner(banner_id int, banner *models.UpdateBanner) error {
	err := service.repo.UpdateBanner(banner_id, banner)

	if err != nil {
		return err
	}

	return nil
}


func (service *BannerService) GetAllBanners(tag_id int, feature_id int, limit int, offset int) ([]*models.BannerDB, error) {
	banners, err := service.repo.GetAllBanners(tag_id, feature_id, limit, offset)

	if err != nil {
		return nil, err
	}

	return banners, nil
}
