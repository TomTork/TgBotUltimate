package users

import (
	"TgBotUltimate/types/Database"
	"context"
)

func GetExpertSystemAnswers(ctx context.Context, db *Database.DB, userID int64) ([]Database.ExpertSystemAnswer, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, user_tg_id, question_id, variant_index
		FROM user_expert_system_answers
		WHERE user_tg_id = $1
		ORDER BY question_id ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answers := make([]Database.ExpertSystemAnswer, 0)
	for rows.Next() {
		var answer Database.ExpertSystemAnswer
		if err := rows.Scan(
			&answer.ID,
			&answer.UserTgID,
			&answer.QuestionID,
			&answer.VariantIndex,
		); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

func SaveExpertSystemAnswer(ctx context.Context, db *Database.DB, answer Database.ExpertSystemAnswer) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO user_expert_system_answers (user_tg_id, question_id, variant_index)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_tg_id, question_id)
		DO UPDATE SET variant_index = EXCLUDED.variant_index`,
		answer.UserTgID,
		answer.QuestionID,
		answer.VariantIndex,
	)
	return err
}

func ResetExpertSystemAnswers(ctx context.Context, db *Database.DB, userID int64) error {
	_, err := db.Exec(
		ctx,
		`DELETE FROM user_expert_system_answers WHERE user_tg_id = $1`,
		userID,
	)
	return err
}
