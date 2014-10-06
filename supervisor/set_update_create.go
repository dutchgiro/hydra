package supervisor

// Set sets the given key to the given value.
// It will create a new key value pair or replace the old one.
// It will not replace a existing directory.
func (e *EtcdClient) Set(key string, value string, ttl uint64) (*Response, error) {
	raw, err := e.RawSet(key, value, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

// Set sets the given key to a directory.
// It will create a new directory or replace the old key value pair by a directory.
// It will not replace a existing directory.
func (e *EtcdClient) SetDir(key string, ttl uint64) (*Response, error) {
	raw, err := e.RawSetDir(key, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

// CreateDir creates a directory. It succeeds only if
// the given key does not yet exist.
func (e *EtcdClient) CreateDir(key string, ttl uint64) (*Response, error) {
	raw, err := e.RawCreateDir(key, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

// UpdateDir updates the given directory. It succeeds only if the
// given key already exists.
func (e *EtcdClient) UpdateDir(key string, ttl uint64) (*Response, error) {
	raw, err := e.RawUpdateDir(key, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

// Create creates a file with the given value under the given key.  It succeeds
// only if the given key does not yet exist.
func (e *EtcdClient) Create(key string, value string, ttl uint64) (*Response, error) {
	raw, err := e.RawCreate(key, value, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

// CreateInOrder creates a file with a key that's guaranteed to be higher than other
// keys in the given directory. It is useful for creating queues.
func (e *EtcdClient) CreateInOrder(dir string, value string, ttl uint64) (*Response, error) {
	raw, err := e.RawCreateInOrder(dir, value, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

// Update updates the given key to the given value.  It succeeds only if the
// given key already exists.
func (e *EtcdClient) Update(key string, value string, ttl uint64) (*Response, error) {
	raw, err := e.RawUpdate(key, value, ttl)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

func (e *EtcdClient) RawUpdateDir(key string, ttl uint64) (*RawResponse, error) {
	ops := Options{
		"prevExist": true,
		"dir":       true,
	}

	return e.put(key, "", ttl, ops)
}

func (e *EtcdClient) RawCreateDir(key string, ttl uint64) (*RawResponse, error) {
	ops := Options{
		"prevExist": false,
		"dir":       true,
	}

	return e.put(key, "", ttl, ops)
}

func (e *EtcdClient) RawSet(key string, value string, ttl uint64) (*RawResponse, error) {
	return e.put(key, value, ttl, nil)
}

func (e *EtcdClient) RawSetDir(key string, ttl uint64) (*RawResponse, error) {
	ops := Options{
		"dir": true,
	}

	return e.put(key, "", ttl, ops)
}

func (e *EtcdClient) RawUpdate(key string, value string, ttl uint64) (*RawResponse, error) {
	ops := Options{
		"prevExist": true,
	}

	return e.put(key, value, ttl, ops)
}

func (e *EtcdClient) RawCreate(key string, value string, ttl uint64) (*RawResponse, error) {
	ops := Options{
		"prevExist": false,
	}

	return e.put(key, value, ttl, ops)
}

func (e *EtcdClient) RawCreateInOrder(dir string, value string, ttl uint64) (*RawResponse, error) {
	return e.post(dir, value, ttl)
}
