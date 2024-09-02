package repositories

import (
	"database/sql"
	"errors"
	"web_storage/internal/models"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (fr *FileRepository) CreateFileRepo(file *models.FileMeta) error {
	query := `INSERT INTO files (name, size, link) VALUES ($1, $2, $3) RETURNING id`
	err := fr.db.QueryRow(query, file.Name, file.Size, file.Link).Scan(&file.Id)
	if err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) GetAllFilesRepo() ([]*models.FileMeta, error) {
	query := `SELECT id, name, size, link FROM files`
	rows, err := fr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.FileMeta
	for rows.Next() {
		var fileMeta models.FileMeta
		if err := rows.Scan(&fileMeta.Id, &fileMeta.Name, &fileMeta.Size, &fileMeta.Link); err != nil {
			return nil, err
		}
		files = append(files, &fileMeta)
	}

	return files, nil
}

func (fr *FileRepository) GetFileByIdRepo(id int) (*models.FileMeta, error) {
	var fileMeta models.FileMeta
	query := `SELECT * from files WHERE id=$1`
	row := fr.db.QueryRow(query, id)
	err := row.Scan(&fileMeta.Id, &fileMeta.Name, &fileMeta.Size, &fileMeta.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &fileMeta, nil
}

func (fr *FileRepository) DeleteFileByIdRepo(id int) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := fr.db.Exec(query, id)
	return err
}
