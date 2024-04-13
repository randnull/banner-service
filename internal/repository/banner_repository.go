package repository

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"time"

	"github.com/jmoiron/sqlx"
	pq "github.com/lib/pq"

	"github.com/randnull/banner-service/internal/config"
	"github.com/randnull/banner-service/internal/errors"
	"github.com/randnull/banner-service/pkg/models"
	"github.com/randnull/banner-service/internal/repository/cast"
)


type Repository struct {
	db *sqlx.DB
}


func NewRepository(cfg *config.Config) *Repository {
	time.Sleep(5 * time.Second) // Ожидание Postgres

	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.UserDB, cfg.PasswordDB, cfg.HostDB, cfg.PortDB, cfg.NameDB)

	db, err := sqlx.Open("postgres", link)

	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Database is ready")

	return &Repository{
		db: db,
	}
}


func (storage *Repository) CheckIfExist(tags_ids []int, feature_id int) bool {
	query := `SELECT COUNT(*) FROM banners WHERE ($1) = ANY(tags_ids) AND feature_id = $2`

	var count int

	for _, tag_id := range tags_ids {
		count = 0
		err := storage.db.Get(&count, query, tag_id, feature_id)

		if err != nil {
			return true
		}

		if count > 0 {
			return true
		}
	}

	return false
}


func (storage *Repository) CheckIfUpdateValid(tags_ids []int, feature_id int, banner_id int) bool {
	query := `SELECT COUNT(*) FROM banners WHERE ($1) = ANY(tags_ids) AND feature_id = $2 AND id != $3`

	var count int

	for _, tag_id := range tags_ids {
		count = 0
		err := storage.db.Get(&count, query, tag_id, feature_id, banner_id)

		if err != nil {
			return true
		}

		if count > 0 {
			return true
		}
	}

	return false
}


func (storage *Repository) AddBaner(banner *models.Banner) (int, error) {
	is_exist := storage.CheckIfExist(banner.TagIds, banner.FeatureId)

	if is_exist {
		return -1, errors.BannerAlreadyExist
	}

	query :=`INSERT INTO banners
				(tags_ids, feature_id, title, text, url, is_active, create_datetime, update_datetime) 
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`

	var id int

	err := storage.db.QueryRow(query,
							pq.Array(banner.TagIds),
							banner.FeatureId,
							banner.Content.Title,
							banner.Content.Text,
							banner.Content.Url,
							banner.IsActive,
							time.Now(),
							time.Now()).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}


func (storage *Repository) DeleteBanner(banner_id int) error {
	query := `DELETE FROM banners WHERE id = $1`

	_, err := storage.db.Exec(query, banner_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.BannerNotFound
		}
		return err
	}

	return nil
}


func (storage *Repository) GetBanner(tag_id int, feature_id int) (*models.Content, error) {
	query := `SELECT title, text, url FROM banners WHERE ($1) = ANY(tags_ids) AND feature_id = $2 AND is_active = true`

	var content models.Content

	err := storage.db.QueryRowx(query, tag_id, feature_id).StructScan(&content)

    if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.BannerNotFound
		}
    	return nil, err
    }
	
	return &content, nil
}


func (storage *Repository) UpdateBanner(banner_id int, banner *models.UpdateBanner) error {
	var tags_ids []int
	var feature_id int

	var tags_ids_from_db pq.Int64Array
	var feature_id_from_db int

	if banner.TagIds == nil || banner.FeatureId == nil {
		check_query := `SELECT tags_ids, feature_id FROM banners WHERE id = $1`

		_ = storage.db.QueryRow(check_query, banner_id).Scan(&tags_ids_from_db, &feature_id_from_db)
	}

	if banner.TagIds == nil {
		tags_ids = cast.CastPqArrayToInt(tags_ids_from_db)
	} else {
		tags_ids = *banner.TagIds
	}

	if banner.FeatureId == nil {
		feature_id = feature_id_from_db
	} else {
		feature_id = *banner.FeatureId
	}

	is_exist := storage.CheckIfUpdateValid(tags_ids, feature_id, banner_id)

	if is_exist {
		return errors.BannerAlreadyExist
	}

	TagsIds := cast.ConvertArrayToInterface(banner.TagIds)
	FeatureId := cast.ConvertIntToInterface(banner.FeatureId)
	Title := cast.ConvertStringToInterface(banner.Content.Title)
	Text := cast.ConvertStringToInterface(banner.Content.Text)
	Url := cast.ConvertStringToInterface(banner.Content.Url)
	IsActive := cast.ConvertBoolToInterface(banner.IsActive)
	
	query := `UPDATE banners SET
				tags_ids = COALESCE($1, tags_ids),
				feature_id = COALESCE($2, feature_id),
				title = COALESCE($3, title),
				text = COALESCE($4, text),
				url = COALESCE($5, url),
				is_active = COALESCE($6, is_active),
				update_datetime = $7
			WHERE id = $8`

	_, err := storage.db.Exec(query, TagsIds, FeatureId, Title, Text, Url, IsActive, time.Now(), banner_id)

	if err != nil {
		return err
	} 

	return nil
}


func (storage *Repository) GetAllBanners(tag_id int, feature_id int, limit int, offset int) ([]*models.BannerDB, error) {
	var rows *sqlx.Rows
	var err error

	var banners []*models.BannerDB

	var args = map[string]interface{}{}

	var query string

	if (tag_id == -1) && (feature_id == -1) {
		query = "SELECT * FROM banners"
	} else {
		query = "SELECT * FROM banners WHERE true"
	}

	if tag_id != -1 {
		query += " AND :tag_id = ANY(tags_ids)"
		args["tag_id"] = tag_id
	}

	if feature_id != -1 {
		query += " AND feature_id = :feature_id"
		args["feature_id"] = feature_id
	}

	if limit != -1 {
		query += " LIMIT :limit"
		args["limit"] = limit
	}

	if offset != -1 {
		query += " OFFSET :offset"
		args["offset"] = offset
	}

	rows, err = storage.db.NamedQuery(query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var banner models.BannerDB
		err := rows.StructScan(&banner)
		if err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}

	return banners, nil
}
