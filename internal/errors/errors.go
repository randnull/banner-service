package errors

import "errors"


// Ошибки доступа
var InvalidToken = errors.New("Invalid token")
var NotRoots = errors.New("No roots")


// Ошибки запроса
var InvalidRequestBody = errors.New("Invalid request body")
var BadBannerIdParam = errors.New("Banner Id not correct")
var BadTagIdParam = errors.New("Tag Id not correct")
var BadFeatureIdParam = errors.New("Feature Id not correct")
var BadUpdateData = errors.New("Data not correct")
var UserNotFound = errors.New("Wrong username or password")
var UsernameAlreadyTaken = errors.New("Username already taken")
var BannerNotFound = errors.New("Banner not found")
var BannerAlreadyExist = errors.New("Banner with couple (feature_id, tag_id) already exist")
var InternalError = errors.New("Internal Error")
