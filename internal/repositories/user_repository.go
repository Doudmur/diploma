package repositories

import (
	"database/sql"
	"diploma/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	rows, err := r.db.Query("SELECT * FROM public.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	row := r.db.QueryRow("SELECT * FROM public.user WHERE user_id=$1", id)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByIin(iin int) (*models.User, error) {
	row := r.db.QueryRow("SELECT * FROM public.user WHERE iin=$1", iin)

	var user models.User
	if err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Iin, &user.Role, &user.BiometricDataHash, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM public.user WHERE user_id = $1", id)
	return err
}

//func (r *BookRepository) CreateBook(book *models.Book) error {
//	err := r.db.QueryRow("INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id", book.Title, book.Author).Scan(&book.ID)
//	return err
//}
//
//func (r *BookRepository) UpdateBook(book *models.Book) error {
//	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2 WHERE id = $3", book.Title, book.Author, book.ID)
//	return err
//}
//
