package persistency

type Persistency interface {
	// Open DB instance
	OpenDB() error
	// Close current DB instance
	CloseDB()
	// Read data from current DB instance
	ReadData(key string) ([]byte, error)
	// Write data to current DB instance
	WriteData(key string, data []byte) error
}
