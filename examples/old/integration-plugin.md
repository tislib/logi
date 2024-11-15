```logi macro

```

```logi
plugin merchant1 {
    Name "Merchant 1"
    Description "This is merchant 1."
    Integration {
        Type "REST"
        Url "https://merchant1.com"
        Auth {
            Type "Basic"
            Username "merchant1"
            Password "password"
        }
    }
```