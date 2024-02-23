package repository



func NewDB(cfg config.DB, log *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode
		)
	)
	if err != nil {
		return nil, err
	}
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
