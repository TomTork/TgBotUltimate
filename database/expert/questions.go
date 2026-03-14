package expert

import (
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Expert"
	"context"
)

func GetQuestions(ctx context.Context, db *Database.DB) ([]Expert.Question, error) {
	rows, err := db.Query(ctx, "SELECT * FROM expert_system")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	questions := make([]Expert.Question, 0)
	for rows.Next() {
		var question Expert.Question
		err = rows.Scan(
			&question.Id,
			&question.Question,
			&question.Variants,
			&question.Results,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}
