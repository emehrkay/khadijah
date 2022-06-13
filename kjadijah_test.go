package khadijah_test

import (
	"fmt"
	"testing"

	k "github.com/emehrkay/khadijah"
)

func TestNew(t *testing.T) {
	type NewTest struct {
		name          string
		settings      []k.KhadijahSetting
		tag           string
		variable      string
		startVariable string
		endVariable   string
		matchClause   string
	}

	tests := []NewTest{
		{
			"all default settings",
			[]k.KhadijahSetting{},
			k.DefaultTagName,
			k.DefaultVariable,
			k.DefaultStartVariable,
			k.DefaultEndVariable,
			k.DefaultMatchClause,
		},
		{
			"custom tag, rest default",
			[]k.KhadijahSetting{
				k.SetTagName("mytag"),
			},
			"mytag",
			k.DefaultVariable,
			k.DefaultStartVariable,
			k.DefaultEndVariable,
			k.DefaultMatchClause,
		},
		{
			"custom variable, rest default",
			[]k.KhadijahSetting{
				k.SetVariable("hello"),
			},
			k.DefaultTagName,
			"hello",
			k.DefaultStartVariable,
			k.DefaultEndVariable,
			k.DefaultMatchClause,
		},
		{
			"custom match cluase, rest default",
			[]k.KhadijahSetting{
				k.SetMatchClause("mymatchclause"),
			},
			k.DefaultTagName,
			k.DefaultVariable,
			k.DefaultStartVariable,
			k.DefaultEndVariable,
			"mymatchclause",
		},
		{
			"all custom",
			[]k.KhadijahSetting{
				k.SetTagName("mytag"),
				k.SetVariable("hellox"),
				k.SetMatchClause("mymatchclause"),
				k.SetStartVariable("startNew"),
				k.SetEndVariable("endNew"),
			},
			"mytag",
			"hellox",
			"startNew",
			"endNew",
			"mymatchclause",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			instance := k.New(test.settings...)

			if instance.TagName != test.tag {
				t.Errorf(`got %v for Tagname, but expected %v`, instance.TagName, test.tag)
			}

			if instance.Variable != test.variable {
				t.Errorf(`got %v for Variable, but expected %v`, instance.Variable, test.variable)
			}

			if instance.StartVariable != test.startVariable {
				t.Errorf(`got %v for StartVariable, but expected %v`, instance.Variable, test.startVariable)
			}

			if instance.EndVariable != test.endVariable {
				t.Errorf(`got %v for EndVariable, but expected %v`, instance.Variable, test.endVariable)
			}

			if instance.MatchClause != test.matchClause {
				t.Errorf(`got %v for MatchClause, but expected %v`, instance.MatchClause, test.matchClause)
			}
		})
	}
}

var ul = "user"
var userLabel *string = &ul

type TestJsonUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TestCustomUser struct {
	ID    string `custom:"id"`
	Name  string `custom:"name"`
	Email string `custom:"email"`
}

type InstanceTypes struct {
	name     string
	settings []k.KhadijahSetting
	user     interface{}
}

var (
	userJ = TestJsonUser{
		ID:    "someeyedee",
		Name:  "somename",
		Email: "emailTest",
	}

	userC = TestCustomUser{
		ID:    "someeyedee",
		Name:  "somename",
		Email: "emailTest",
	}

	cases = []InstanceTypes{
		{
			"with default settings",
			[]k.KhadijahSetting{},
			userJ,
		},
		{
			"with custom tag setting",
			[]k.KhadijahSetting{
				k.SetTagName("custom"),
			},
			userC,
		},
		{
			"with custom tag and variable settings",
			[]k.KhadijahSetting{
				k.SetTagName("custom"),
				k.SetVariable("xxxyyyzzz"),
			},
			userC,
		},
	}
)

func aliasField(alias, field string) string {
	return fmt.Sprintf(`%s.%s = $%s`, alias, field, field)
}

func TestCreateNode(t *testing.T) {
	type Create struct {
		name       string
		expected   string
		withReturn bool
		excludes   []string
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := k.New(testCase.settings...)
			user := testCase.user
			tests := []Create{
				{
					"create without any excludes and a return",
					fmt.Sprintf("CREATE (%s:%s {id: $id, name: $name, email: $email}) RETURN %s", instance.Variable, *userLabel, instance.Variable),
					true,
					[]string{},
				},
				{
					"create without any excludes and without a return",
					fmt.Sprintf("CREATE (%s:%s {id: $id, name: $name, email: $email})", instance.Variable, *userLabel),
					false,
					[]string{},
				},
				{
					"create while excluding id and a return",
					fmt.Sprintf("CREATE (%s:%s {name: $name, email: $email}) RETURN %s", instance.Variable, *userLabel, instance.Variable),
					true,
					[]string{"id"},
				},
				{
					"create while excluding id and without a return",
					fmt.Sprintf("CREATE (%s:%s {name: $name, email: $email})", instance.Variable, *userLabel),
					false,
					[]string{"id"},
				},
				{
					"create while excluding everything and with a return",
					fmt.Sprintf("CREATE (%s:%s {}) RETURN %s", instance.Variable, *userLabel, instance.Variable),
					true,
					[]string{"id", "name", "email"},
				},
				{
					"create while excluding everything and without a return",
					fmt.Sprintf("CREATE (%s:%s {})", instance.Variable, *userLabel),
					false,
					[]string{"id", "name", "email"},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					maxx := instance.CreateNode(user, userLabel, test.withReturn, test.excludes...)

					if maxx.Query != test.expected {
						t.Errorf(`expected: "%s" but got: "%v"`, test.expected, maxx.Query)
					}
				})
			}
		})
	}
}

func TestUpdateNodeSuite(t *testing.T) {
	type Update struct {
		name        string
		expected    string
		matchClause string
		withReturn  bool
		excludes    []string
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := k.New(testCase.settings...)
			user := testCase.user
			tests := []Update{
				{
					"should update with default match clause and return",
					fmt.Sprintf(`MERGE (%s:%s {id: $id}) SET %s, %s, %s RETURN %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "id"),
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email"),
						instance.Variable),
					k.DefaultMatchClause,
					true,
					[]string{},
				},
				{
					"should update with default match clause and without a return",
					fmt.Sprintf(`MERGE (%s:%s {id: $id}) SET %s, %s, %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "id"),
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email")),
					k.DefaultMatchClause,
					false,
					[]string{},
				},
				{
					"should update with default match clause while ignoring id and return",
					fmt.Sprintf(`MERGE (%s:%s {id: $id}) SET %s, %s RETURN %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email"),
						instance.Variable),
					k.DefaultMatchClause,
					true,
					[]string{"id"},
				},
				{
					"should update with default match clause while ignoring id and without a return",
					fmt.Sprintf(`MERGE (%s:%s {id: $id}) SET %s, %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email")),
					k.DefaultMatchClause,
					false,
					[]string{"id"},
				},
				{
					"should update with custom match clause while ignoring id and return",
					fmt.Sprintf(`MERGE (%s:%s {custom: $custom}) SET %s, %s RETURN %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email"),
						instance.Variable),
					"{custom: $custom}",
					true,
					[]string{"id"},
				},
				{
					"should update with custom match clause while ignoring id and without a return",
					fmt.Sprintf(`MERGE (%s:%s {custom: $custom}) SET %s, %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email")),
					"{custom: $custom}",
					false,
					[]string{"id"},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					t.Run("UpdateNodeWithMatch", func(t *testing.T) {
						maxx := instance.UpdateNodeWithMatch(user, userLabel, test.matchClause, test.withReturn, test.excludes...)

						if maxx.Query != test.expected {
							t.Errorf(`expected: "%s" but got: "%v"`, test.expected, maxx.Query)
						}
					})

					t.Run("UpdateNode", func(t *testing.T) {
						// skip tests where the match clause is not the default one
						// beacuase the UpdateNode function only uses the default match clause
						if test.matchClause != k.DefaultMatchClause {
							t.Skip("augmenting the matchClause is not relevant to this UpdateNode")
							return
						}

						maxx := instance.UpdateNode(user, userLabel, test.withReturn, test.excludes...)

						if maxx.Query != test.expected {
							t.Errorf(`expected: "%s" but got: "%v"`, test.expected, maxx.Query)
						}
					})
				})
			}
		})
	}
}

func TestDeleteNodeSuite(t *testing.T) {
	type Delete struct {
		name        string
		expected    string
		matchClause string
		detach      bool
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := k.New(testCase.settings...)
			user := testCase.user
			tests := []Delete{
				{
					"can delete node with custom matching",
					fmt.Sprintf(`MATCH (%s {abcDEFG=$a}) DELETE %s`, instance.Variable, instance.Variable),
					"{abcDEFG=$a}",
					false,
				},
				{
					"can detach delete node with custom matching",
					fmt.Sprintf(`MATCH (%s {abcDEFG=$a}) DETACH DELETE %s`, instance.Variable, instance.Variable),
					"{abcDEFG=$a}",
					true,
				},
				{
					"can delete node with default matching",
					fmt.Sprintf(`MATCH (%s %s) DELETE %s`, instance.Variable, instance.MatchClause, instance.Variable),
					instance.MatchClause,
					false,
				},
				{
					"can detach delete node with default matching",
					fmt.Sprintf(`MATCH (%s %s) DETACH DELETE %s`, instance.Variable, instance.MatchClause, instance.Variable),
					instance.MatchClause,
					true,
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					t.Run("DeleteNodeWithMatch", func(t *testing.T) {
						maxx := instance.DeleteNodeWithMatch(user, test.detach, test.matchClause)

						if maxx.Query != test.expected {
							t.Errorf(`expected: "%s" but got: "%s"`, test.expected, maxx.Query)
						}
					})

					t.Run("DetachDeleteNodeWithMatch", func(t *testing.T) {
						// ignore tests without detach
						if !test.detach {
							t.Skip("tests wihtout detach do not apply to DetachDeleteNodeWithMatch")
						}

						maxx := instance.DetachDeleteNodeWithMatch(user, test.matchClause)

						if maxx.Query != test.expected {
							t.Errorf(`expected: "%s" but got: "%v"`, test.expected, maxx.Query)
						}
					})

					t.Run("DetachDeleteNode", func(t *testing.T) {
						// ignore tests without detach
						if !test.detach || test.matchClause != k.DefaultMatchClause {
							t.Skip("tests wihtout detach or use the default match clause do not apply to DetachDeleteNode")
						}

						maxx := instance.DetachDeleteNode(user)

						if maxx.Query != test.expected {
							t.Errorf(`expected: "%s" but got: "%v"`, test.expected, maxx.Query)
						}
					})
				})
			}
		})
	}
}

func TestCreateEdgeSuite(t *testing.T) {

}

func TestUpdateEdgeSuite(t *testing.T) {

}
