package middleware

const (
	authorizationHeaderKey string = "Authorization"
	authorizationType      string = "Bearer"
)

// // Get User Auth middleware
// func (c *middleware) AuthenticateUser() gin.HandlerFunc {
// 	fmt.Println("authenticateUser func middleware")
// 	return c.authorize(token.User)
// }

// // authorize request on request header using user type
// func (c *middleware) authorize(tokenUser token.UserType) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authorizationValues := ctx.GetHeader(authorizationHeaderKey)
// 		fmt.Println("authorize tokenUser middleware", ctx.GetHeader(authorizationValues))

// 		authFields := strings.Fields(authorizationValues)
// 		if len(authFields) < 2 {

// 			err := errors.New("authorization token not provided properly with prefix of Bearer")

// 			response.ErrorResponse(ctx, "Failed to authorize request", err, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		authType := authFields[0]
// 		accessToken := authFields[1]

// 		if !strings.EqualFold(authType, authorizationType) {
// 			err := errors.New("invalid authorization type")
// 			response.ErrorResponse(ctx, "Unauthorized user", err, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		tokenVerifyReq := token.VerifyTokenRequest{
// 			TokenString: accessToken,
// 			UsedFor:     tokenUser,
// 		}

// 		verifyRes, err := c.tokenService.VerifyToken(tokenVerifyReq)

// 		if err != nil {
// 			response.ErrorResponse(ctx, "Unauthorized user", err, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("userId", verifyRes.UserID)
// 	}
// }
