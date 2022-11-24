package gosplunk

func (c Client) LookupRetrieve(search string) ([]map[string]interface{}, error) {
	lookupSearch := NewSearch{
		Search:   search,
		Earliest: "1",
		Latest:   "now",
	}
	res, err := c.SearchExport(lookupSearch)
	if err != nil {
		c.Logger.Error(err)
		return nil, err
	}
	var Map []map[string]interface{}
	for _, item := range res {
		Map = append(Map, item.Result)
	}
	return Map, nil
}

// returns map based on search results of unique keys and their values
func (c Client) LookupRetrieveKv(search string) (map[string]string, error) {
	lookupSearch := NewSearch{
		Search:   search,
		Earliest: "1",
		Latest:   "now",
	}
	res, err := c.SearchExport(lookupSearch)
	if err != nil {
		c.Logger.Error(err)
		return nil, err
	}
	Map := make(map[string]string)
	for _, item := range res {
		value, ok := item.Result["value"].(string)
		if !ok {
			c.Logger.Debugw("Could not convert value to string", "value", item.Result["value"])
			continue
		}
		key, ok := item.Result["key"].(string)
		if !ok {
			c.Logger.Debugw("Could not convert key to string", "key", item.Result["key"])
			continue
		}
		Map[key] = value
	}
	return Map, nil
}
