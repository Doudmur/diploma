package repositories

import (
	"database/sql"
	"diploma/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetNotificationByUserID(id int) (*models.Notification, error) {
	row := r.db.QueryRow("SELECT * FROM public.notification WHERE user_id=$1", id)

	var notification models.Notification
	if err := row.Scan(&notification.NotificationId, &notification.UserId, &notification.Message, &notification.Type, &notification.SentAt); err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	err := r.db.QueryRow("INSERT INTO notification(user_id, message, type) VALUES ($1, $2, $3) RETURNING user_id", notification.UserId, notification.Message, notification.Type).Scan(&notification.NotificationId)
	return err
}

func (r *NotificationRepository) DeleteNotification(id int) error {
	_, err := r.db.Exec("DELETE FROM public.notification WHERE notification_id = $1", id)
	return err
}
