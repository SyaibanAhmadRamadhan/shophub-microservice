package presentations

import userservice "user-service/internal/services/user_service"

type Dependency struct {
	UserService userservice.UserService
}
