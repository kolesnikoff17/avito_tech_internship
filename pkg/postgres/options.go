package postgres

// Option is a type of functions-setters
type Option func(*Db)

// MaxConn sets up db's pool max connections
func MaxConn(mc int) Option {
	return func(db *Db) {
		if mc > 0 {
			db.maxConns = mc
		}
	}
}
