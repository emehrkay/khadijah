package khadijah_test

import (
	"fmt"
	"reflect"
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
		matchClause   k.M
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
				k.SetMatchClause(k.M{"mymatch": "clause"}),
			},
			k.DefaultTagName,
			k.DefaultVariable,
			k.DefaultStartVariable,
			k.DefaultEndVariable,
			k.M{"mymatch": "clause"},
		},
		{
			"all custom",
			[]k.KhadijahSetting{
				k.SetTagName("mytag"),
				k.SetVariable("hellox"),
				k.SetMatchClause(k.M{"mymatch": "clause"}),
				k.SetStartVariable("startNew"),
				k.SetEndVariable("endNew"),
			},
			"mytag",
			"hellox",
			"startNew",
			"endNew",
			k.M{"mymatch": "clause"},
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

			if !reflect.DeepEqual(instance.MatchClause, test.matchClause) {
				t.Errorf(`got %v for MatchClause, but expected %v`, instance.MatchClause, test.matchClause)
			}
		})
	}
}

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
	user2    interface{}
}

type Knows struct {
	_ struct{}
}

type Follows struct {
	_     struct{}
	Since string `json:"since"`
}

var (
	ul                = "user"
	userLabel *string = &ul
	userJ             = TestJsonUser{
		ID:    "someeyedeeJSON",
		Name:  "somenameJSON",
		Email: "emailTestJSON",
	}

	userC = TestCustomUser{
		ID:    "someeyedeeCUSTOM",
		Name:  "somenameCUSTOM",
		Email: "emailTestCUSTOM",
	}

	knows   = Knows{}
	follows = Follows{
		Since: "yesterday",
	}

	cases = []InstanceTypes{
		{
			"with default settings",
			[]k.KhadijahSetting{},
			userJ,
			userJ,
		},
		{
			"with custom tag setting",
			[]k.KhadijahSetting{
				k.SetTagName("custom"),
			},
			userC,
			userJ,
		},
		{
			"with custom tag and variable settings",
			[]k.KhadijahSetting{
				k.SetTagName("custom"),
				k.SetVariable("xxxyyyzzz"),
			},
			userC,
			userJ,
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
					fmt.Sprintf("CREATE (%s:%s ) RETURN %s", instance.Variable, *userLabel, instance.Variable),
					true,
					[]string{"id", "name", "email"},
				},
				{
					"create while excluding everything and without a return",
					fmt.Sprintf("CREATE (%s:%s )", instance.Variable, *userLabel),
					false,
					[]string{"id", "name", "email"},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					maxx := instance.CreateNode(user, userLabel, test.withReturn, test.excludes...)

					if maxx.Query != test.expected {
						t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
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
		matchClause k.M
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
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%v) = $id SET %s, %s, %s RETURN %s`,
						instance.Variable,
						*userLabel,
						instance.Variable,
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
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%v) = $id SET %s, %s, %s`,
						instance.Variable,
						*userLabel,
						instance.Variable,
						aliasField(instance.Variable, "id"),
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email")),
					k.DefaultMatchClause,
					false,
					[]string{},
				},
				{
					"should update with default match clause while ignoring id and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%v) = $id SET %s, %s RETURN %s`,
						instance.Variable,
						*userLabel,
						instance.Variable,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email"),
						instance.Variable),
					k.DefaultMatchClause,
					true,
					[]string{"id"},
				},
				{
					"should update with default match clause while ignoring id and without a return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%v) = $id SET %s, %s`,
						instance.Variable,
						*userLabel,
						instance.Variable,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email")),
					k.DefaultMatchClause,
					false,
					[]string{"id"},
				},
				{
					"should update with custom match clause while ignoring id and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE custom = $custom SET %s, %s RETURN %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email"),
						instance.Variable),
					k.M{"custom": "custom"},
					true,
					[]string{"id"},
				},
				{
					"should update with custom match clause while ignoring id and without a return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE custom = $custom SET %s, %s`,
						instance.Variable,
						*userLabel,
						aliasField(instance.Variable, "name"),
						aliasField(instance.Variable, "email")),
					k.M{"custom": "custom"},
					false,
					[]string{"id"},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					t.Run("UpdateNodeWithMatch", func(t *testing.T) {
						maxx := instance.UpdateNodeWithMatch(user, userLabel, test.matchClause, test.withReturn, test.excludes...)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
						}
					})

					t.Run("UpdateNode", func(t *testing.T) {
						// skip tests where the match clause is not the default one
						// beacuase the UpdateNode function only uses the default match clause
						if !reflect.DeepEqual(test.matchClause, k.DefaultMatchClause) {
							t.Skip("augmenting the matchClause is not relevant to this UpdateNode")
							return
						}

						maxx := instance.UpdateNode(user, userLabel, test.withReturn, test.excludes...)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
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
		matchClause k.M
		detach      bool
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := k.New(testCase.settings...)
			user := testCase.user
			tests := []Delete{
				{
					"can delete node with custom matching",
					fmt.Sprintf(`MATCH (%s) WHERE abcDEFG = $a DELETE %s`, instance.Variable, instance.Variable),
					k.M{"abcDEFG": "a"},
					false,
				},
				{
					"can detach delete node with custom matching",
					fmt.Sprintf(`MATCH (%s) WHERE abcDEFG = $a DETACH DELETE %s`, instance.Variable, instance.Variable),
					k.M{"abcDEFG": "a"},
					true,
				},
				{
					"can delete node with default matching",
					fmt.Sprintf(`MATCH (%s) WHERE %s DELETE %s`, instance.Variable, instance.RootMaxx.MatchClause, instance.Variable),
					instance.MatchClause,
					false,
				},
				{
					"can detach delete node with default matching",
					fmt.Sprintf(`MATCH (%s) WHERE %s DETACH DELETE %s`, instance.Variable, instance.RootMaxx.MatchClause, instance.Variable),
					instance.MatchClause,
					true,
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					t.Run("DeleteNodeWithMatch", func(t *testing.T) {
						maxx := instance.DeleteNodeWithMatch(user, test.detach, test.matchClause)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
						}
					})

					t.Run("DetachDeleteNodeWithMatch", func(t *testing.T) {
						// ignore tests without detach
						if !test.detach {
							t.Skip("tests wihtout detach do not apply to DetachDeleteNodeWithMatch")
						}

						maxx := instance.DetachDeleteNodeWithMatch(user, test.matchClause)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
						}
					})

					t.Run("DetachDeleteNode", func(t *testing.T) {
						// ignore tests without detach
						if !test.detach || !reflect.DeepEqual(test.matchClause, k.DefaultMatchClause) {
							t.Skip("tests wihtout detach or use the default match clause do not apply to DetachDeleteNode")
						}

						maxx := instance.DetachDeleteNode(user)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
						}
					})
				})
			}
		})
	}
}

func TestCreateEdgeSuite(t *testing.T) {

	type Create struct {
		name             string
		expected         string
		startMatchClause k.M
		direction        string
		endMatchClasue   k.M
		relationship     interface{}
		withReturn       bool
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := k.New(testCase.settings...)
			start := testCase.user
			end := testCase.user2
			startVar := instance.StartVariable
			endVar := instance.EndVariable
			edgeVar := instance.Variable
			m := k.Maxine{}
			sm := m.Parse(start)
			startLabel := sm.EntityName
			em := m.Parse(end)
			endLabel := em.EntityName

			tests := []Create{
				{
					"create out edge with default matches and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%s) = $start_id MATCH (%s:%s) WHERE id(%s) = $end_id CREATE (%s)-[%s:Knows %s]->(%s) RETURN %s, %s, %s`,
						startVar,
						startLabel,
						startVar,
						endVar,
						endLabel,
						endVar,
						startVar,
						edgeVar,
						instance.RootMaxx.CreateQuery,
						endVar,
						startVar,
						edgeVar,
						endVar),
					k.DefaultMatchClause,
					"out",
					k.DefaultMatchClause,
					knows,
					true,
				},
				{
					"create in edge with default matches and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%s) = $start_id MATCH (%s:%s) WHERE id(%s) = $end_id CREATE (%s)<-[%s:Knows %s]-(%s) RETURN %s, %s, %s`,
						startVar,
						startLabel,
						startVar,
						endVar,
						endLabel,
						endVar,
						startVar,
						edgeVar,
						instance.RootMaxx.CreateQuery,
						endVar,
						startVar,
						edgeVar,
						endVar),
					k.DefaultMatchClause,
					"in",
					k.DefaultMatchClause,
					knows,
					true,
				},
				{
					"create undirected edge with default matches and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%s) = $start_id MATCH (%s:%s) WHERE id(%s) = $end_id CREATE (%s)-[%s:Knows %s]-(%s) RETURN %s, %s, %s`,
						startVar,
						startLabel,
						startVar,
						endVar,
						endLabel,
						endVar,
						startVar,
						edgeVar,
						instance.RootMaxx.CreateQuery,
						endVar,
						startVar,
						edgeVar,
						endVar),
					k.DefaultMatchClause,
					"",
					k.DefaultMatchClause,
					knows,
					true,
				},
				{
					"create out edge with default matches and no return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%s) = $start_id MATCH (%s:%s) WHERE id(%s) = $end_id CREATE (%s)-[%s:Knows %s]->(%s)`,
						startVar,
						startLabel,
						startVar,
						endVar,
						endLabel,
						endVar,
						startVar,
						edgeVar,
						instance.RootMaxx.CreateQuery,
						endVar),
					k.DefaultMatchClause,
					"out",
					k.DefaultMatchClause,
					knows,
					false,
				},
				{
					"create in edge with default matches and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%s) = $start_id MATCH (%s:%s) WHERE id(%s) = $end_id CREATE (%s)<-[%s:Knows %s]-(%s)`,
						startVar,
						startLabel,
						startVar,
						endVar,
						endLabel,
						endVar,
						startVar,
						edgeVar,
						instance.RootMaxx.CreateQuery,
						endVar),
					k.DefaultMatchClause,
					"in",
					k.DefaultMatchClause,
					knows,
					false,
				},
				{
					"create undirected edge with default matches and return",
					fmt.Sprintf(`MATCH (%s:%s) WHERE id(%s) = $start_id MATCH (%s:%s) WHERE id(%s) = $end_id CREATE (%s)-[%s:Knows %s]-(%s)`,
						startVar,
						startLabel,
						startVar,
						endVar,
						endLabel,
						endVar,
						startVar,
						edgeVar,
						instance.RootMaxx.CreateQuery,
						endVar),
					k.DefaultMatchClause,
					"",
					k.DefaultMatchClause,
					knows,
					false,
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					t.Run("CreateEdgeWithMatches", func(t *testing.T) {
						maxx := instance.CreateEdgeWithMatches(start, nil, test.startMatchClause, test.direction, end, nil, test.endMatchClasue, test.relationship, nil, test.withReturn)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
						}
					})

					t.Run("CreateEdgeWithMatches", func(t *testing.T) {
						maxx := instance.CreateEdge(start, end, test.relationship, test.direction, nil, nil, nil, test.withReturn)

						if maxx.Query != test.expected {
							t.Errorf("\nexpected: \n\t%s \nbut got: \n\t%v\n", test.expected, maxx.Query)
						}
					})
				})
			}

		})
	}
}

func TestUpdateEdgeSuite(t *testing.T) {

}
