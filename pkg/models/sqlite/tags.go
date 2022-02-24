package sqlite

// CreateTags - adds tags to the database
func (m *ForumModel) CreateTags(tags []string) error {
	stmt := `INSERT INTO "main"."tags" (Tag) VALUES (?)`

	for _, tag := range tags {
		_, err := m.DB.Exec(stmt, tag)
		if err != nil {
			continue
		}
	}
	return nil
}

// GetAllTags - retrieves all tags from the database
func (m *ForumModel) GetAllTags() ([]string, error) {
	stmt := `SELECT tag FROM tags`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var tag string
	var tags []string
	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
