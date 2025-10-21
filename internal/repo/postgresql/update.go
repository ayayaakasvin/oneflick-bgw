package postgresql

import (
	"context"
	"math"

	"github.com/ayayaakasvin/trends-updater/internal/models"
	"github.com/ayayaakasvin/trends-updater/script"
)

func (p *PostgreSQL) FetchUpdateTrending(ctx context.Context) ([]models.EventStats, error) {
	rows, err := p.conn.QueryContext(ctx, script.Top10Script)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []models.EventStats
	for rows.Next() {
		var e models.EventStats
		if err := rows.Scan(
			&e.EventUUID, &e.Title, &e.Description,
			&e.StartingTime, &e.EndingTime, &e.Status,
			&e.Capacity, &e.ImageURL, &e.CategoryName,
			&e.OrganizerUsername, &e.TotalTicketsSold, &e.Rank,
		); err != nil {
			return nil, err
		}

		if e.Capacity > 0 {
			e.FillRate = math.Round((float64(e.TotalTicketsSold) / float64(e.Capacity)) * 100)
		} else {
			e.FillRate = 0
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (p *PostgreSQL) ArchiveOldEvents(ctx context.Context) (int, error) {
	var archivedCount int
	err := p.conn.QueryRowContext(ctx, script.ArchiveOldEventsScript).Scan(&archivedCount)
	if err != nil {
		return 0, err
	}
	return archivedCount, nil
}

func (p *PostgreSQL) UpdateEventStatuses(ctx context.Context) (int64, error) {
	result, err := p.conn.ExecContext(ctx, script.UpdateStatusesScript)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}