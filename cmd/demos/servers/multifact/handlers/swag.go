package handlers

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
func _() {} // "GET /swagger"

// @Summary Hello World
// @Description Returns a hello message in different formats based on the "Accept" header.
// @Tags Hello
// @Accept json
// @Accept text/html
// @Produce json
// @Produce text/html
// @Success 200 {object} map[string]interface{} "JSON response with a hello message"
// @Failure 406 {object} map[string]string "Invalid Accept Header"
// @Router / [get]
func _() {} // "GET /"

// RegisterUsernamePassword handles user registration with username and password.
// @Summary Register a new user
// @Description Registers a user by accepting a username and password in JSON format.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RequestBody true "Username and Password"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {object} ErrorResponseBody "Insufficient username or password"
// @Failure 406 {string} string "Invalid Accept Header"
// @Failure 500 {object} ErrorResponseBody "Internal Server Error"
// @Router /register/ [POST]
func _() {} // "POST /register/"

// AuthByUsernamePassword authenticates a user using their username and password.
// @Summary Authenticate user
// @Description Authenticates a user with username and password. Depending on the authentication mode, a strong token or OTP process is initiated.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body AuthUsernameResourceBody true "Username and Password"
// @Success 200 {object} TokenResponseData "Authentication successful"
// @Failure 400 {object} ErrorResponseBody "Insufficient username or password"
// @Failure 401 {object} ErrorResponseBody "Unauthorized - invalid credentials"
// @Failure 406 {string} string "Invalid Accept Header"
// @Failure 500 {object} ErrorResponseBody "Internal Server Error"
// @Router /token/ [POST]
func _() {} // "POST /token/"

// SubmitOtp verifies the OTP and generates a strong token for the user.
// @Summary Submit OTP for Authentication
// @Description Verifies the provided OTP and token, and generates a strong authentication token upon successful validation.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body OTP true "OTP and Token"
// @Success 200 {object} SubmitOtpResponseBody "OTP successfully verified, returns a strong token"
// @Failure 400 {object} ErrorResponseBody "Insufficient OTP or token"
// @Failure 401 {object} ErrorResponseBody "Unauthorized - invalid OTP or token"
// @Failure 406 {string} string "Invalid Accept Header"
// @Failure 500 {object} ErrorResponseBody "Internal Server Error"
// @Router /otp/ [POST]
func _() {} // "POST /otp/"

// PatchUser modifies user details based on the provided patch operation.
// @Summary Modify User Details
// @Description Allows modifying user details, such as authentication mode and email, through a patch operation.
// @Tags User
// @Accept json
// @Produce json
// @Param request body PatchRequestBody true "Patch operation details"
// @Success 200 {string} string "User modified successfully"
// @Failure 400 {object} ErrorResponseBody "Bad Request - Invalid input or multiple patches provided"
// @Failure 406 {string} string "Invalid Accept Header"
// @Router /user/ [PATCH]
func _() {} // "PATCH /user/"
