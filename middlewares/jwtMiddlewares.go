package middlewares

// func SetClientJWTmiddlewares(g *echo.Group) {
// 	jwtConfig := kc.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", kc.App.ENV))

// 	middlewareConfig := middleware.JWTConfig{
// 		Claims:        &handlers.JWTclaims{},
// 		SigningMethod: "HS512",
// 		SigningKey:    []byte(jwtConfig["jwt_secret"].(string)),
// 	}
// 	g.Use(middleware.JWTWithConfig(middlewareConfig))

// 	g.Use(validateJWTclient)
// }

// func validateJWTclient(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		user := c.Get("user")
// 		token := user.(*jwt.Token)

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 			if claims["role"] == "client" {
// 				return next(c)
// 			} else {
// 				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
// 			}
// 		}

// 		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
// 	}
// }
