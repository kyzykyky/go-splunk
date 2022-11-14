package gosplunk

func (c Client) RetrieveLookup(search string) ([]map[string]interface{}, error) {
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
