Usage of Logi language in testing.
===============================

# Api Integration Testing

```logi
apiHelper AuthenticateApiHelper {
    exec() {
       // Authenticate user
    }
}

apiTest UserApiTest {
    userHelper AuthenticateApiHelper
    
    test createAndGetUser() {
       var user = http.post<User>("/users", {name: "John Doe", email: "john@doe.com"})
       
        var user2 = http.get<User>("/users/" + user.id)
        
        assertEqual(user, user2)
    }
}
```