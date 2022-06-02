# GraphQL module for the Modulus framework
This is a module for the Modulus framework that lets developers use graphql API in their projects. It is a wrapper for the library https://gqlgen.com/


# Adding the module
To integrate the module follow next steps:
* Add ENV variables from .env.dist to your root .env file
* Copy gqlgen.yml.dist to your project root directory as  gqlgen.yml
* In the gqlgen.yml change package name "boilerplate" to the name of your package
* Copy graph.dist directory to the folder "internal" of your project as graph
* In some of your modules add a graphql file in the root of your module's folder with a schema. For example
```graphql
extend type Query {
    user(id: String!): User
}
extend type Mutation {
    register(email: String!, name: String!): User
}

type User {
    id: String!
    email: String!
    name: String!
}
```
* Call `go run github.com/99designs/gqlgen generate --config gqlgen.yml` to generate Golang code of resolvers and models
* Copy graph/config.go.dist to the graph/config.go. Fix package names in imports, and add this module config to the list of modules of your application.
* Check if the router, for example https://github.com/debugger84/modulus-router-httprouter is added as a module to your application


