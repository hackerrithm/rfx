package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/authenticating"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/listing"
)

// Users are a list of users
var Users []authenticating.User

func RandNumberRunes() int {
	s := "9780486653556"
	var factor, sum1, sum2 int
	for i, c := range s[:12] {
		if i%2 == 0 {
			factor = 1
		} else {
			factor = 3
		}
		buf := make([]byte, 1)
		_ = utf8.EncodeRune(buf, c)
		value, _ := strconv.Atoi(string(buf))
		sum1 += value * factor
		sum2 += (int(c) - '0') * factor

	}

	return sum1 + sum2
}

// define custom GraphQL ObjectType `userType` for our Golang struct `User`
// Note that
// - the fields in our userType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"uid": &graphql.Field{
			Type: graphql.String,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"firstname": &graphql.Field{
			Type: graphql.String,
		},
		"lastname": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
		"gender": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{createUser(text:"My+new+user"){id,text,done}}'
		*/
		"createUser": &graphql.Field{
			Type:        userType, // the return type for this field
			Description: "Create new user",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"firstname": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"lastname": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"gender": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// marshall and cast the argument value
				username, _ := params.Args["username"].(string)
				firstname, _ := params.Args["firstname"].(string)
				lastanme, _ := params.Args["lastanme"].(string)
				password, _ := params.Args["password"].(string)
				gender, _ := params.Args["gnder"].(string)

				// figure out new id
				newID := RandNumberRunes()

				// perform mutation operation here
				// for e.g. create a User and save to DB.
				newUser := authenticating.User{
					UID:       newID,
					UserName:  username,
					FirstName: firstname,
					LastName:  lastanme,
					Password:  password,
					Gender:    gender,
				}

				Users = append(Users, newUser)

				var a authenticating.Service
				a.AddUser(newUser)

				// return the new User object that we supposedly save to DB
				// Note here that
				// - we are returning a `User` struct instance here
				// - we previously specified the return Type to be `userType`
				// - `User` struct maps to `userType`, as defined in `userType` ObjectConfig`
				return Users, nil
			},
		},
	},
})

// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/graphql?query={lastUser{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		/*
		   curl -g 'http://localhost:8080/graphql?query={user(id:"b"){id,text,done}}'
		*/
		"user": &graphql.Field{
			Type:        userType,
			Description: "Get single user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var foundUser listing.User
				idQuery, isOK := params.Args["id"].(int)
				if isOK {
					var l listing.Service

					foundUser, _ = l.GetUser(idQuery)
				}

				return foundUser, nil
			},
		},

		// "lastUser": &graphql.Field{
		// 	Type:        userType,
		// 	Description: "Last user added",
		// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// 		return Users[len(Users)-1], nil
		// 	},
		// },

		/*
		   curl -g 'http://localhost:8080/graphql?query={userList{id,text,done}}'
		*/
		"userList": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "List of users",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var l listing.Service
				var users = l.GetUsers()
				fmt.Println("recovery game")
				return users, nil
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

type key int

const (
	header key = iota
	// ...
)

func httpHeaderMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), header, r.Header)

		next.ContextHandler(ctx, w, r)
	})
}

// SetupMux set up
func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()

	// graphql Handler
	graphqlHandler := http.HandlerFunc(graphqlHandlerFunc)

	// add in addContext middlware
	mux.Handle("/graphql", addContext(graphqlHandler))

	return mux
}

func graphqlHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// get query
	opts := handler.NewRequestOptions(r)

	// execute graphql query
	params := graphql.Params{
		Schema:         schema, // defined in another file
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        r.Context(),
	}
	result := graphql.Do(params)

	// output JSON
	var buff []byte
	w.WriteHeader(http.StatusOK)
	buff, _ = json.Marshal(result)

	w.Write(buff)
}

func addContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.WithValue(r.Context(), header, r.Header)

		// next.ContextHandler(ctx, w, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
