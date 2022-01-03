package db

func (s *Sqlite) setupDatabase() error {
	stmt, err := s.Prepare(`CREATE TABLE IF NOT EXISTS devices ( 
		id          	text NOT NULL,
		name        	text NOT NULL,
		mac_address   text NOT NULL,
		date_created	text NOT NULL,
		date_updated	text NOT NULL,
		CONSTRAINT 		Pk_devices_id PRIMARY KEY ( id )
	);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = s.Prepare(`CREATE TABLE IF NOT EXISTS accesstokens ( 
		id          	text NOT NULL,
		device_id     text NOT NULL,
		token   			text NOT NULL,
		date_created	text NOT NULL,
		date_updated	text NOT NULL,
		CONSTRAINT 		Pk_accesstokens_id PRIMARY KEY ( id )
		FOREIGN KEY 	( device_id ) REFERENCES devices( id )
	);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
