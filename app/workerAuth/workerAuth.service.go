package workerAuth

// Login service logs in a user
//func Login(c *fiber.Ctx) error {
//	b := new(LoginDto)
//
//	if err := utils.ParseBodyAndValidate(c, b); err != nil {
//		panic(fmt.Sprintf("error parsing login dto:%v", err))
//	}
//
//	sw, err := supplyWorkers.GetSupplyWorkerByPhoneNumber(b.PhoneNumber)
//	if err != nil {
//		return fiber.NewError(fiber.StatusUnauthorized,
//			fmt.Sprintf("PhoneNumber: %v not found", b.PhoneNumber),
//		)
//	}
//
//	if sw.FirstTimeLogin {
//		if b.Password == "" {
//			return c.JSON(fiber.Map{
//				"name":           sw.FullName,
//				"firstTimeLogin": true,
//			})
//		} else {
//			hash := password.Generate(b.Password)
//			sw.HashPassword = hash
//			payload := jwt.TokenPayload{Id: sw.ID.String()}
//			accessToken := jwt.GenerateAccessToken(&payload)
//			refreshToken := jwt.GenerateRefreshToken(&payload)
//			hashRefreshToken := password.Generate(refreshToken)
//			sw.HashRefreshToken = hashRefreshToken
//			sw.FirstTimeLogin = false
//			_ = database.UpdateDocument(sw.Key, sw, "supplyWorkers")
//			sw.HashPassword = ""
//			sw.HashRefreshToken = ""
//			return c.JSON(fiber.Map{
//				"accessToken":  accessToken,
//				"refreshToken": refreshToken,
//				"user":         sw,
//			})
//
//		}
//	} else {
//		if b.Password == "" {
//			return fiber.NewError(fiber.StatusUnauthorized,
//				fmt.Sprintf("Password not found"),
//			)
//		} else {
//			match := password.Verify(sw.HashPassword, b.Password)
//			if !match {
//				return fiber.NewError(fiber.StatusUnauthorized,
//					fmt.Sprintf("wrong PhoneNumber or Password"),
//				)
//			}
//			payload := jwt.TokenPayload{Id: sw.ID.String()}
//			accessToken := jwt.GenerateAccessToken(&payload)
//			refreshToken := jwt.GenerateRefreshToken(&payload)
//			hashRefreshToken := password.Generate(refreshToken)
//			sw.HashRefreshToken = hashRefreshToken
//			sw.FirstTimeLogin = false
//			_ = database.UpdateDocument(sw.Key, sw, "supplyWorkers")
//			sw.HashPassword = ""
//			sw.HashRefreshToken = ""
//			return c.JSON(fiber.Map{
//				"accessToken":  accessToken,
//				"refreshToken": refreshToken,
//				"user":         sw,
//			})
//
//		}
//	}
//
//	return fiber.NewError(fiber.StatusInternalServerError,
//		fmt.Sprintf("something unaccepted happened"),
//	)
//}

//func refreshToken(c *fiber.Ctx) error {
//	b := new(refreshTokenDto)
//
//	if err := utils.ParseBodyAndValidate(c, b); err != nil {
//		panic(fmt.Sprintf("error parsing refresh token dto:%v", err))
//	}
//
//	payload, err := jwt.VerifyRefreshToken(b.RefreshToken)
//	if err != nil {
//		return utils.CustomErrorResponse(401, 4012, err, "", c)
//	}
//	splitedId := strings.Split(payload.Id, "/")
//	key := splitedId[1]
//
//	sw := supplyWorkers.GetSupplyWorkerByKey(key)
//	match := password.Verify(sw.HashRefreshToken, b.RefreshToken)
//	if !match {
//		return fiber.NewError(fiber.StatusUnauthorized,
//			fmt.Sprintf("wrong PhoneNumber or Password"),
//		)
//	}
//	accessToken := jwt.GenerateAccessToken(payload)
//	sw.HashPassword = ""
//	sw.HashRefreshToken = ""
//	return c.JSON(fiber.Map{
//		"accessToken": accessToken,
//		"user":        sw,
//	})
//
//}

////// Signup service creates a user
////func Signup(ctx *fiber.Ctx) error {
////	b := new(.SignupDTO)
////
////	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
////		return err
////	}
////
////	err := dal.FindUserByEmail(&struct{ ID string }{}, b.Email).Error
////
////	// If email already exists, return
////	if !errors.Is(err, gorm.ErrRecordNotFound) {
////		return fiber.NewError(fiber.StatusConflict, "Email already exists")
////	}
////
////	user := &dal.User{
////		Name:     b.Name,
////		Password: password.Generate(b.Password),
////		Email:    b.Email,
////	}
////
////	// Create a user, if error return
////	if err := dal.CreateUser(user); err.Error != nil {
////		return fiber.NewError(fiber.StatusConflict, err.Error.Error())
////	}
////
////	// generate access token
////	t := jwt.Generate(&jwt.TokenPayload{
////		ID: user.ID,
////	})
////
////	return ctx.JSON(&types.AuthResponse{
////		User: &types.UserResponse{
////			ID:    user.ID,
////			Name:  user.Name,
////			Email: user.Email,
////		},
////		Auth: &types.AccessResponse{
////			Token: t,
////		},
////	})
////}
