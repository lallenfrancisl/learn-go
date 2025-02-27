package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/lallenfrancisl/gopi"
	"github.com/lallenfrancisl/greenlight-api/internal/data"
)

func (app *application) GetDocs(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/openapi.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Greenlight Movies DB API",
		},
		DarkMode: true,
	})

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	fmt.Fprintln(w, htmlContent)
}

func Docs() {
	docs.
		DefineTag(gopi.TagDef{
			Name:        "Healthcheck",
			Description: "Healthcheck routes",
		}).
		DefineTag(gopi.TagDef{
			Name:        "Movies",
			Description: "APIs for managing movies",
		}).
		DefineTag(gopi.TagDef{
			Name:        "Users",
			Description: "User managment APIs",
		})

	docs.Route("/v1/healthcheck").Get().
		Tags([]string{"Healthcheck"}).
		Response(http.StatusOK, envelope{
			"status": "",
			"system_info": map[string]string{
				"status":      "available",
				"environment": "",
				"version":     "",
			},
		})

	docs.Route("/v1/movies").Post().
		Summary("Create a new movie").
		Tags([]string{"Movies"}).
		Body(&createMoviePayload{}).
		Response(
			http.StatusOK,
			envelope{"movie": &data.Movie{}},
		)

	docs.Route("/v1/movies/{id}").Get().
		Summary("Get a movie by id").
		Params(
			gopi.PathParam("id", 0).
				Description("Id of the movie").
				Required(),
		).
		Tags([]string{"Movies"}).
		Response(http.StatusOK, envelope{"movie": &data.Movie{}})

	docs.Route("/v1/movies/{id}").Patch().
		Summary("Update a movie by id").
		Tags([]string{"Movies"}).
		Params(
			gopi.PathParam("id", 0).
				Description("Id of the movie").
				Required(),
		).
		Body(updateMoviePayload{}).
		Response(http.StatusOK, envelope{"movie": data.Movie{}})

	docs.Route("/v1/movies/{id}").Delete().
		Summary("Delete a movie by id").
		Tags([]string{"Movies"}).
		Params(
			gopi.PathParam("id", 0).
				Description("Id of the movie").
				Required(),
		).
		Response(http.StatusOK, envelope{"message": ""})

	docs.Route("/v1/movies").Get().
		Summary("List all the movies").
		Tags([]string{"Movies"}).
		Params(
			gopi.QueryParam("title", "").
				Description("Search by title"),
			gopi.QueryParam("genres", []string{}).
				Description("Filter by list of genres"),
			gopi.QueryParam("page", 0).
				Description("Page number"),
			gopi.QueryParam("page_size", 0).
				Description("Number of items in each page"),
			gopi.QueryParam("sort", "").
				Description("Sort by given field name and direction"),
		).
		Response(
			http.StatusOK,
			envelope{"movies": []data.Movie{}, "metadata": data.Metadata{}},
		)

	docs.Route("/v1/users").Post().
		Summary("Create a new user").
		Tags([]string{"Users"}).
		Description("Create a new user and start the activation flow").
		Body(RegisterUserPayload{}).
		Response(http.StatusOK, envelope{"user": data.User{}})

	docs.Route("/v1/users/{id}/activate").Put().
		Summary("Activate user").
		Tags([]string{"Users"}).
		Body(activateUserPayload{}).
		Response(http.StatusOK, envelope{"user": data.User{}})

	docs.Route("/v1/users/login").Post().
		Summary("Login").
		Tags([]string{"Users"}).
		Body(loginPayload{}).
		Response(http.StatusOK, envelope{"credentials": ""})
}

func writeDocsFile(docs *gopi.Gopi) {
	Docs()

	js, err := docs.MarshalJSONIndent("", "    ")
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	err = os.WriteFile("./docs/openapi.json", js, os.FileMode(os.O_TRUNC))
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}
