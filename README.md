# Khadijah

An easy way to convert some structs into some simple CRUD Cypher queries.

> You can build the complex stuff by hand, this isn't a real query builder

## Usage

Given a struct

```go
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

You can use khadijah to generate `MATCH` `CREATE` or `DELETE` Cypher queries like this:

```go
func main() {
	mark := User{
		ID:    "someID",
		Name:  "emehrkay",
		Email: "spam@aol.com",
	}
	instance := khadijah.New()
	label := "User"
	create := instance.CreateNode(mark, &label, true) // if nil is paseed in for label, the entity's Name is used

	fmt.Printf(`%+v`, create)
}
```

All actions will return a `Maxine` instance:

```
&Maxine{
	Query:CREATE (flava:User {id: $id, name: $name, email: $email}) RETURN flava 
	SetQuery:flava.id = $id, flava.name = $name, flava.email = $email 
	CreateQuery:{id: $id, name: $name, email: $email} 
	Params:map[
		email:spam@aol.com 
		id:someID 
		name:emehrkay
	] 
	TagName:json 
	Variable:flava
	EntityName:User
}
```

Create Node while removing fields

```go
create := instance.CreateNode(mark, nil, true, "id", "email")

// CREATE (flava:User {name: $name}) RETURN flava
```

Update Node

```go
update := instance.UpdateNode(mark, &label, true, "id")

// MERGE (flava:User {id: $id}) SET flava.name = $name, flava.email = $email RETURN flava
```

Detach Delete Node

```go
delete := instance.DetachDeleteNode(mark)

// MATCH (flava {id: $id}) DETACH DELETE flava
```

> these functions are abstracted from a base version which offer more control. Look at the souce

## Extra Recipes 

You can set various attributes when you make an instance of Khadijah. Say you wanted to use a custom tag and not json:

```go

type User struct {
	ID    string `myCoolTag:"id"`
	Name  string `myCoolTag:"name"`
	Email string `myCoolTag:"email"`
}

func main() {
	mark := User{
		ID:    "someID",
		Name:  "emehrkay",
		Email: "spam@aol.com",
	}
	instance := khadijah.New(
		khadijah.SetTagName("myCoolTag"),
	)
	// works the same as above, but using the myCoolTag values for the properties
}
```

## F.A.Q. 

1. What's with the naming?

* Have you seen Living Single? If not, stop reading and go watch it. `Khadijah` runs `Flava` magazine. She is the main character and everything flows through her. `Synclaire`, her cousin and assistant, is quirky and quietly handles things. She is reponsible for connections. `Regine` is their roommate who is constantly dating, that's why she is in charge of single node augmentations. `Maxine` is their boisterous, shoot-from-the-hip neighbor lawer and is in charge of interrogating entities. `Overton` is the handyman in the building that they live in and is reponsible for utilty functionality. `Kyle` is fancy and doesnt have a role in this lib, yet. Im thinking some sort of pretty printer in a future release.

2. Where are the docs?

* I'll put together something soon. Hopefully the interface is simple enough to grok becuase this doesn't really do too much. Just start with `khadijah` and work your way out.

3. Go version support?

* I don't know. I didn't use `any` or any generics because those features arent really necessary and not everyone has 1.18+

4. Why?

* I'm lazy and I wrote a couple of string queries and decided to abstract it to a function. Then I said "this could be a lib," and now we're here.

5. Do these queries work?

* 😬  I don't know. I haven't tested them all. They should though. They should.

## License

MIT
