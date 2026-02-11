package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"context"
)

func GetTags(ctx context.Context, db *Database.DB) ([]Database.ITag, error) {
	rows, err := db.Query(ctx, queries.GetAll("tags"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]Database.ITag, 0)
	for rows.Next() {
		var tag Database.ITag
		err = rows.Scan(
			&tag.Id,
			&tag.Code,
			&tag.FlatCode,
			&tag.Name,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagByCode(ctx context.Context, db *Database.DB, code string) (*Database.ITag, error) {
	tag := Database.ITag{}
	err := db.QueryRow(ctx, queries.GetS("tags", "code", code)).Scan(
		&tag.Id,
		&tag.Code,
		&tag.FlatCode,
		&tag.Name,
	)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func CreateTag(ctx context.Context, db *Database.DB, tag Database.ITag) error {
	existsTag, _ := GetTagByCode(ctx, db, *tag.Code)
	if existsTag == nil {
		err := db.QueryRow(ctx, queries.Create("tags", queries.TagsFields, queries.TagsValues(tag))).Scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateTag(ctx context.Context, db *Database.DB, tag Database.ITag) error {
	err := db.QueryRow(ctx, queries.UpdateS("tags", "code", *tag.Code, queries.TagsFields, queries.TagsValues(tag))).Scan()
	if err != nil {
		return err
	}
	return nil
}
