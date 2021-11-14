package server

import (
	"net/http"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/0xdod/go-realworld/pkg"
)

func (s *Server) createArticle() http.HandlerFunc {
	type Input struct {
		Article struct {
			Title       string   `json:"title" validate:"required"`
			Description string   `json:"description"`
			Body        string   `json:"body" validate:"required"`
			Tags        []string `json:"tagList"`
		} `json:"article"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.Article); err != nil {
			validationError(w, err)
			return
		}

		article := conduit.Article{
			Title:       input.Article.Title,
			Body:        input.Article.Body,
			Slug:        pkg.Slugify(input.Article.Title),
			Description: input.Article.Description,
		}

		article.AddTags(input.Article.Tags...)
		user := userFromContext(r.Context())
		article.Author = user

		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		if err := s.articleService.CreateArticle(r.Context(), &article); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"article": article})
	}
}
