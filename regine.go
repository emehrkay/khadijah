package khadijah

import "fmt"

func newRegine(matchClause string) *regine {
	return &regine{
		matchClause: matchClause,
	}
}

type regine struct {
	matchClause string
}

func (r *regine) nodeWithProperties(entity interface{}, label string) *Maxine {
	maxx := rootMaxx.Parse(entity)
	maxx.Query = fmt.Sprintf(`(%s:%s %s)`, maxx.Variable, label, maxx.CreateQuery)

	return maxx
}

func (r *regine) matchNodeWithMatch(entity interface{}, label, matchClause string, withReturn bool) *Maxine {
	maxx := rootMaxx.Parse(entity)
	maxx.Query = fmt.Sprintf(`MATCH (%s:%s %s)`, maxx.Variable, label, matchClause)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s`, maxx.Query, maxx.Variable)
	}

	return maxx
}

func (r *regine) matchNode(entity interface{}, label string, withReturn bool) *Maxine {
	return r.matchNodeWithMatch(entity, label, r.matchClause, withReturn)
}

// CREATE (x:Label {param: $param}) RETURN x
func (r *regine) createNode(entity interface{}, label string, withReturn bool, excludes ...string) *Maxine {
	maxx := rootMaxx.Parse(entity, excludes...)
	maxx.Query = fmt.Sprintf(`CREATE (%s:%s %s)`, maxx.Variable, label, maxx.CreateQuery)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s`, maxx.Query, maxx.Variable)
	}

	return maxx
}

// MERGE (x:Label {param: $param}) SET param1 = $param1 RETURN x
func (r *regine) updateNodeWithMatch(entity interface{}, label, matchClause string, withReturn bool, excludes ...string) *Maxine {
	maxx := rootMaxx.Parse(entity, excludes...)
	maxx.Query = fmt.Sprintf(`MERGE (%s:%s %s) SET %s`, maxx.Variable, label, matchClause, maxx.SetQuery)

	if withReturn {
		maxx.Query = fmt.Sprintf(`%s RETURN %s`, maxx.Query, maxx.Variable)
	}

	return maxx
}

// MERGE (x:Label {id: $id}) SET param1 = $param1 RETURN x
func (r *regine) updateNode(entity interface{}, label string, withReturn bool, excludes ...string) *Maxine {
	return r.updateNodeWithMatch(entity, label, r.matchClause, withReturn, excludes...)
}

// MATCH (x {param: $param}) [DETACH] DELETE x
func (r *regine) deleteNodeWithMatch(entity interface{}, detach bool, matchClause string) *Maxine {
	detachClause := " "
	if detach {
		detachClause = " DETACH "
	}

	maxx := rootMaxx.Parse(entity)
	maxx.Query = fmt.Sprintf(`MATCH (%s %s)%sDELETE %s`, maxx.Variable, matchClause, detachClause, maxx.Variable)
	return maxx
}

// MATCH (x {param: $param}) [DETACH] DELETE x
func (r *regine) detachDeleteNodeWithMatch(entity interface{}, matchClause string) *Maxine {
	return r.deleteNodeWithMatch(entity, true, matchClause)
}

// MATCH (x {param: $param}) DETACH DELETE x
func (r *regine) detachDeleteNode(entity interface{}) *Maxine {
	return r.deleteNodeWithMatch(entity, true, r.matchClause)
}

// MATCH (x {param: $param}) DELETE x
func (r *regine) deleteNode(entity interface{}) *Maxine {
	return r.deleteNodeWithMatch(entity, false, r.matchClause)
}
