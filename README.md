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
	create := instance.CreateNode(mark, "User", true)

	fmt.Printf(`%+v`, create)
}
```

This will return a `maxine` instance that has a create query

```
CREATE (flava:User {id: $id, name: $name, email: $email}) RETURN flava
```

## F.A.Q. 

1. What's with the naming?

* Have you seen Living Single? If not, stop reading and go watch it. `Khadijah` runs `Flava` magazine. She is the main character and everything flows through her. `Synclaire`, her cousin and assistant, is quirky and quietly handles things. She is reponsible for connections. `Regine` is their roommate who is constantly dating, that's why she is in charge of single node augmentations. `Maxine` is their boisterous, shoot-from-the-hip neighbor lawer and is in charge of interrogating entities. `Overton` is the handyman in the building that they live in and is reponsible for utilty functionality. `Kyle` is fancy and doesnt have a role in this lib, yet. Im thinking some sort of pretty printer in a future release.

2. Where are the docs?

* I'll put together something soon. Hopefully the interface is simple enough to grok becuase this doesn't really do too much. Just start with `khadijah` and work your way out.

## License

MIT
