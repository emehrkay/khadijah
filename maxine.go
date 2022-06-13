package khadijah

import (
	"fmt"
	"reflect"
	"strings"
)

// NewMaxine will create a new instance of Maxine with a given tagName and variable
func NewMaxine(tagName, variable, paramPrefix string) *Maxine {
	maxx := &Maxine{
		Params:     M{},
		TagName:    tagName,
		Variable:   variable,
		ParamPefix: paramPrefix,
	}

	return maxx
}

type Maxine struct {
	// hods the findal quuery. This isn't prdouced in Maxine, but is callers
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

	// the prefix for the params
	ParamPefix string `json:"paramPrefix"`

	// Holds the name of the entity for use in situations where a label isn't provided
	EntityName string `json:"entityName"`
}

// Parse does the work of converting a struct to query placeloders and
// a params map. This will exclude any fields passed in as exclude
// it uses the tag name to matach the field name to the cypher property
// a new Maxine instance is created on every call allowing reuse of previously
// defined properties
// example return:
// Maxine{
//     CreateQuery: "{email: $email, username: $username, password: $password}",
// 	   SetQuery: "u.email = $email, u.usrname = $username, u.password = $password",
//     Parms: M{"email": entity.Email, "username": entity.Username, "password": entity.Password},
// }
func (m *Maxine) Parse(entity interface{}, exclude ...string) *Maxine {
	maxx := NewMaxine(m.TagName, m.Variable, m.ParamPefix)
	queryParams := []string{}
	setParams := []string{}
	entityType := reflect.TypeOf(entity)

	// resolve the entity type and name
	for entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	maxx.EntityName = entityType.Name()
	entityValue := reflect.ValueOf(entity)
	exclude = append(exclude, "-")

	var getTag = func(tag string) string {
		return fmt.Sprintf(`%s%s`, m.ParamPefix, tag)
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag, ok := field.Tag.Lookup(m.TagName)
		tagFixed := getTag(tag)

		// only add the param if it is not in the exclude list
		if ok && !Contains(exclude, tag) {
			queryParams = append(queryParams, fmt.Sprintf(`%s: $%s`, tag, tagFixed))
			setParams = append(setParams, fmt.Sprintf(`%s.%s = $%s`, maxx.Variable, tag, tagFixed))
		}

		fieldValue := entityValue.FieldByName(field.Name)
		value := fieldValue.Interface()
		maxx.Params[tagFixed] = value
	}

	maxx.CreateQuery = fmt.Sprintf(`{%s}`, strings.Join(queryParams, ", "))
	maxx.SetQuery = strings.Join(setParams, ", ")

	return maxx
}
