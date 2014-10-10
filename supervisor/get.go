package supervisor

// Get gets the file or directory associated with the given key.
// If the key points to a directory, files and directories under
// it will be returned in sorted or unsorted order, depending on
// the sort flag.
// If recursive is set to false, contents under child directories
// will not be returned.
// If recursive is set to true, all the contents will be returned.
func (e *EtcdClient) Get(key string, sort, recursive bool) (*Response, error) {
	raw, err := e.RawGet(key, sort, recursive)

	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

func (e *EtcdClient) RawGet(key string, sort, recursive bool) (*RawResponse, error) {
	ops := Options{
		"recursive": recursive,
		"sorted":    sort,
	}

	return e.get(key, ops)
}
