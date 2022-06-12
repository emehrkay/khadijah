package khadijah

import (
	"fmt"
	"reflect"
	"strings"
)

// NewMaxine will create a new instance of maxine with a given tagName and variable
func NewMaxine(tagName, variable string) *maxine {
	maxx := &maxine{
		Params:   M{},
		TagName:  tagName,
		Variable: variable,
	}

	return maxx
}

type maxine struct {
	// hods the findal quuery. This isn't prdouced in maxine, but is callers
	Query string `json:"query"`

	// holds the set information: "var.name = $name, var.age = $age, var.paramN = $paramN"
	SetQuery string `json:"setQuery"`

	// holds the create params information: name: $name, age: $age, paramN: $paramN
	CreateQuery string `json:"createQuery"`

	// holds the params to be used in the final query: M{name: entity.Name, paramN: entity.ParamN}
	Params M `json:"params"`

	// defines which tag should be used
	TagName string `json:"tagName"`

	// defines the prefix value used in the SetQuery
	Variable string `json:"variable"`
}

// Parse does the work of converting a struct to query placeloders and
// a params map. This will exclude any fields passed in as exclude
// it uses the tag name to matach the field name to the cypher property
// a new maxine instance is created on every call allowing reuse of previously
// defined properties
// example return:
// maxine{
//     CreateQuery: "{email: $email, username: $username, password: $password}",
// 	   SetQuery: "u.email = $email, u.usrname = $username, u.password = $password",
//     Parms: M{"email": entity.Email, "username": entity.Username, "password": entity.Password},
// }
func (m *maxine) Parse(entity interface{}, exclude ...string) *maxine {
	maxx := NewMaxine(m.TagName, m.Variable)
	queryParams := []string{}
	setParams := []string{}
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	exclude = append(exclude, "-")

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag, ok := field.Tag.Lookup(m.TagName)

		// only add the param if it is not in the exclude list
		if ok && !Contains(exclude, tag) {
			queryParams = append(queryParams, fmt.Sprintf(`%s: $%s`, tag, tag))
			setParams = append(setParams, fmt.Sprintf(`%s.%s = $%s`, maxx.Variable, tag, tag))
		}

		fieldValue := entityValue.FieldByName(field.Name)
		value := fieldValue.Interface()
		maxx.Params[tag] = value
	}

	maxx.CreateQuery = fmt.Sprintf(`{%s}`, strings.Join(queryParams, ", "))
	maxx.SetQuery = strings.Join(setParams, ", ")

	return maxx
}
