package auth

type Storer interface {
	run() error
}

type AuthHandler struct {
	store Storer
}

type UserStore struct {
	name string
	age  int
}

func (h AuthHandler) AddUser(name string, age int) {
	h.store.run()
}

func HandleGetUser(store Store) func(c *gin.Context) {
    // We return a function that matches the interface needed by the router
	return func(c *gin.Context) {
		username := c.Param("username")
	    // But can still use the database interface within the returned function
		u, err := store.GetUserByUsername(username)
		if err != nil {
			c.JSON(404, gin.H{"error": "Could not find user with this username"})
			return
		}
		c.JSON(200, u)
	}
}
