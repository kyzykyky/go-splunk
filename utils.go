package gosplunk

import "fmt"

func (c Client) getResourcePrefix() string {
	if c.App == "" && c.Username == "" {
		return "/services"
	} else if c.App == "" && c.Username != "" {
		return fmt.Sprintf("/servicesNS/%s/search", c.Username)
	} else if c.Username == "" && c.App != "" {
		return fmt.Sprintf("/servicesNS/-/%s", c.App)
	} else {
		return fmt.Sprintf("/servicesNS/%s/%s", c.Username, c.App)
	}
}

// Converts default MV mixed string and []string to []string
func multivalueConverter(asset interface{}) []string {
	field := make([]string, 0)
	r, ok := asset.([]interface{})
	if !ok {
		ast, ok := asset.(string)
		if ok {
			field = append(field, ast)
		}
	} else {
		for _, a := range r {
			ast, ok := a.(string)
			if ok {
				field = append(field, ast)
			}
		}
	}
	return field
}
