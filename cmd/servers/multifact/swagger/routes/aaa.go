package mux

// ShowAccount godoc
// @Summary      swagger api
// @Description  swagger docs
// @Tags         accounts
// @Produce      json
// @Param        id   query      string  true  "file name" 	Enums(index.html, doc.json)
// @Success 200 {string} string "ok, html or json"
// @Header       200              {string}  Content-Type  "content type"
// @Failure      404  {object}  int
// @Failure      500  {object}  int
// @Router       /swagger [get]
func SwagHandler(w http.ResponseWriter, r *http.Request) {
	swagger.Handler()(w, r)
}
