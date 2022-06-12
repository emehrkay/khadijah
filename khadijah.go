package khadijah

var (
	DefaultTagName       = "json"
	DefaultVariable      = "flava"
	DefaultStartVariable = "start"
	DefaultEndVariable   = "end"
	DefaultMatchClause   = "{id: $id}"
	DefaultSettings      = []KhadijahSetting{
		SetTagName(DefaultTagName),
		SetVariable(DefaultVariable),
		SetStartVariable(DefaultStartVariable),
		SetEndVariable(DefaultEndVariable),
		SetMatchClause(DefaultMatchClause),
	}
	rootMaxx *maxine
)

// KhadijahSetting type that defines a setting for khadijah
type KhadijahSetting func(instance *khadijah)

// SetTagName will set khadijah.TagName
func SetTagName(tagName string) KhadijahSetting {
	return func(instance *khadijah) {
		instance.TagName = tagName
	}
}

// SetVariable will set khadijah.Variable
func SetVariable(variable string) KhadijahSetting {
	return func(instance *khadijah) {
		instance.Variable = variable
	}
}

// SetMatchClause will set khadijah.matchCaluse
func SetMatchClause(matchClause string) KhadijahSetting {
	return func(instance *khadijah) {
		instance.MatchClause = matchClause
	}
}

// SetMatchClause will set khadijah.StartVariable
func SetStartVariable(startVariable string) KhadijahSetting {
	return func(instance *khadijah) {
		instance.StartVariable = startVariable
	}
}

// SetMatchClause will set khadijah.matchCaluse
func SetEndVariable(endVariable string) KhadijahSetting {
	return func(instance *khadijah) {
		instance.EndVariable = endVariable
	}
}

// New creates an instance of khadijah with "json" as the default tag name
// used to pull values from the passed in structs and "flava" as the default
// variable that is used in the returned queries
func New(settings ...KhadijahSetting) *khadijah {
	// set defaults and override them
	settings = append(DefaultSettings, settings...)
	instance := &khadijah{}
	instance.Apply(settings...)
	rootMaxx = NewMaxine(instance.TagName, instance.Variable)

	return instance
}

type khadijah struct {
	TagName       string
	Variable      string
	StartVariable string
	EndVariable   string
	MatchClause   string
}

// Apply will set some properties on the instance
func (k *khadijah) Apply(settings ...KhadijahSetting) {
	for _, setFn := range settings {
		setFn(k)
	}
}

// NodeWithProperties creates a simple (var:label {propts}) string
func (k *khadijah) NodeWithProperties(entity interface{}, label string) *maxine {
	reg := newRegine(k.MatchClause)

	return reg.nodeWithProperties(entity, label)
}

// MatchNode creates a simple Match (var:label {props}) cypther query
func (k *khadijah) MatchNode(entity interface{}, label string, withReturn bool) *maxine {
	reg := newRegine(k.MatchClause)

	return reg.matchNode(entity, label, withReturn)
}

// CreateNode builds a simple cypher CREATE query that looks like:
//     CREATE (x:Label {param: $param}) RETURN x
func (k *khadijah) CreateNode(entity interface{}, label string, withReturn bool, excludes ...string) *maxine {
	reg := newRegine(k.MatchClause)

	return reg.createNode(entity, label, withReturn, excludes...)
}

// UpdateNodeWithMatch builds a simpole cyper Merge ... SET query that looks like:
//		MERGE (x:Label {param: $param}) SET param1 = $param1 RETURN x
func (k *khadijah) UpdateNodeWithMatch(entity interface{}, label, matchClause string, withReturn bool, excludes ...string) *maxine {
	reg := newRegine(k.MatchClause)

	return reg.updateNodeWithMatch(entity, label, matchClause, withReturn, excludes...)
}

// UpdateNode works like UpdateNodeWithMatch, but defaults the matchClause to {id: $id}
// creates a query that looks like:
//		MATCH (x:Label {id: $id}) SET param1 = $param1 RETURN x
func (k *khadijah) UpdateNode(entity interface{}, label string, withReturn bool, excludes ...string) *maxine {
	return k.UpdateNodeWithMatch(entity, label, k.MatchClause, withReturn, excludes...)
}

// DeleteNodeWithMatch builds a cypher MATCH .. DELETE quer that looks like:
//		MATCH (x {param: $param}) [DETACH] DELETE x
func (k *khadijah) DeleteNodeWithMatch(entity interface{}, detach bool, matchClause string) *maxine {
	reg := newRegine(k.MatchClause)

	return reg.deleteNodeWithMatch(entity, detach, matchClause)
}

// DetachDeleteNodeWithMatch build a MATCH ... DETACH DELETE cypher query using
// the provided matching clause
//		MATCH (x {param: $param}) [DETACH] DELETE x
func (k *khadijah) DetachDeleteNodeWithMatch(entity interface{}, matchClause string) *maxine {
	return k.DeleteNodeWithMatch(entity, true, matchClause)
}

// DetachDeleteNode build a MATCH ... DETACH DELETE cypher query using the default
// matching clause
//		MATCH (x {param: $param}) [DETACH] DELETE x
func (k *khadijah) DetachDeleteNode(entity interface{}) *maxine {
	return k.DeleteNodeWithMatch(entity, true, k.MatchClause)
}

// CreateEdge builds a complex MATCh (nodeA), (nodeB) CREATE query
//		MATCH (start:Lable {matches}), (end:Label {props}) CREATE (start)-[edge:label {matches}]->(end) RETURN start, end, edge
func (k *khadijah) CreateEdge(start, end, edge interface{}, direction, startLabel, endLabel, edgeLabel string, withReturn bool, excldues ...string) *maxine {
	return k.CreateEdgeWithMatches(start, startLabel, DefaultMatchClause, direction, end, endLabel, DefaultMatchClause, edge, edgeLabel, withReturn, excldues...)
}

// CreateEdgeWithMatches a complex MATCh (nodeA), (nodeB) CREATE query
//		MATCH (start:Lable {matches}), (end:Label {props}) CREATE (start)-[edge:label {matches}]->(end) RETURN start, end, edge
func (k *khadijah) CreateEdgeWithMatches(start interface{}, startLabel, startMatchClause, direction string, end interface{}, endLabel, endMatchClause string, edge interface{}, edgeLabel string, withReturn bool, excldues ...string) *maxine {
	syn := newSynclaire(k.MatchClause, k.StartVariable, k.EndVariable)

	return syn.createEdgeWithMatches(start, startLabel, startMatchClause, direction, end, endLabel, endMatchClause, edge, edgeLabel, withReturn, excldues...)
}

func (k *khadijah) UpdateEdgeWithMatches(start interface{}, startLabel, startMatchClause, direction string, end interface{}, endLabel, endMatchClause string, edge interface{}, edgeLabel, edgeMatchClause string, withReturn bool, excldues ...string) *maxine {
	syn := newSynclaire(k.MatchClause, k.StartVariable, k.EndVariable)

	return syn.updateEdgeWithMatches(start, startLabel, startMatchClause, direction, end, endLabel, endMatchClause, edge, edgeLabel, edgeMatchClause, withReturn, excldues...)
}

func (k *khadijah) UpdateEdge(start interface{}, startLabel, direction string, end interface{}, endLabel string, edge interface{}, edgeLabel string, withReturn bool, excldues ...string) *maxine {
	return k.UpdateEdgeWithMatches(start, startLabel, DefaultMatchClause, direction, end, endLabel, DefaultMatchClause, edge, edgeLabel, DefaultMatchClause, withReturn, excldues...)
}
