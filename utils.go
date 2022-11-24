package gosplunk

import "fmt"

func (c Client) getResourcePrefix() string {
	return getResourcePrefix(NameSpace{c.App, c.Username})
}

func getResourcePrefix(ns NameSpace) string {
	if ns.App == "" && ns.User == "" {
		return "/services"
	} else if ns.App == "" && ns.User != "" {
		return fmt.Sprintf("/servicesNS/%s/search", ns.User)
	} else if ns.User == "" && ns.App != "" {
		return fmt.Sprintf("/servicesNS/nobody/%s", ns.App)
	} else {
		return fmt.Sprintf("/servicesNS/%s/%s", ns.User, ns.App)
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
