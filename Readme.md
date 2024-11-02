Logi - A Language for Building Extensible Systems
=================================================

# Overview

Logi is a purpose-built language designed to facilitate the creation of modular, adaptable, and extensible systems. The
focus is on providing flexibility for developers to expand and customize systems with minimal effort, ensuring
scalability, and encouraging code reusability.

# Motivation

With help of Logi, you can simply define a declarative system that can be easily written and read by devs and AI.

## Examples

### Example 1: Simple API definition

```logi
api MyAPI {
    endpoint /hello {
        method GET {
            response {
                message: "Hello, World!"
            }
        }
    }
    
    endpoint /add-two-numbers {
        method POST {
            response (param1 int4, param2 int8) int8 {
                    return param1 + param2
            }
        }
    }
    
    endpoint /users {
        inject(UserService) as userService

        method POST {
            response (user User) User {
                return userService.createUser(user)
            }
        }
    }
}
```

### Example 2: Repository definition

```logi
entity User {
    id: int [id, required]
    name: string [required, default: "John Doe"]
    email: string [transient]
}

repository UserRepository User {
    query getUserById(id int) `SELECT * FROM users WHERE id = $1`
    query getUsers() `SELECT * FROM users`
    query getUsersByName(name string) queryBuilder(User).where('name', name).get()
}
```

### Example 3: Service definition

```logi
service UserService {
    inject UserRepository as userRepository

    method getUserById(id int) User {
        return userRepository.getUserById(id)
    }
    
    method getUsers() []User {
        return userRepository.getUsers()
    }
    
    method getUsersByName(name string) []User {
        return userRepository.getUsersByName(name)
    }
}
```

### Example 4: Macro definition

Create a macro to handle events(which does not exists in language)

```logi-macro
macro eventHandler {
    kind: Syntax
    args (name Name)
 
    syntax {
       on <eventName Name>() { }
       before <eventName Name>() <returnType Type> { }
       before <eventName Name>() { }
    }
}
```

Define user creation event handler

```logi
eventHandler UserCreated {
    inject(EmailService) as emailService

    on userCreated (user User) {
        emailService.sendEmail(user.email, "Welcome to our platform!")
    }
}
```

# Key Features

## 1. Declarative Syntax

• Logi uses a declarative approach to define system behaviors, making it more readable and maintainable.
• The syntax is optimized for configuration-based programming, enabling rapid changes without complex coding.

Idea: Allow for “self-documenting code” through declarative structures that describe the functionality directly within
the syntax.

## 2. Adaptive Design with Macros

In Logi, macros are used to adjust language to be flexible and adaptable to different use cases.
Macros can be used to define custom elements, properties, and behaviors, allowing developers to create systems that are
tailored to their specific needs.

## 3. Modular Architecture

Logi promotes a modular architecture that allows developers to build systems in a structured and organized way.

## 4. Ready for Data representation

Logi is designed to work with data structures, making it easy to represent and manipulate data in a variety of formats.
It supports to easily define data values also to replace json, xml, yaml, etc.

## 5. AI Ready

Logi language is supposed to be ready to be used by AI to integrate it to other systems.
WIth help of Logi,
you can define a specific declarative language for specific system where logi can be used as communication interface.

## 6. Highly extensible

With help of macros and modular architecture, Logi is highly extensible and can be used to build complex systems.

## 7. Functional Programming

Logi is designed to support functional programming paradigms, making it easier to write clean, concise, and efficient
code.

# Declaration

In Logi language, there are different type of elements.

## Common syntax

The common syntax for declaring an element is:

For definition of data

```
<elementType> <elementName> {
    <property1> <value1>
    <property2> <value2>
    <property3> {
        <nestedProperty1> <nestedValue1>
        <nestedProperty2> <nestedValue2>
    }
}

struct User {
    id int
    name string
    email string
    tags string[]
}
```

For definition of a data type

```
User user1 {
    id  33
    name "John Doe"
    email "AbcX"
    tags ["tag1", "tag2"]
}
```

For definition of code block

```
<codeType> <codeName> (param1 type1, param2 type2) returnType {
    ...code...
    return <value>
}
```

### Full definition

#### Rules
1. Each Statement should be on a single line. This is for both logi and logi-macro files.

```
<macroName> <definitionName> 
```

## Macros

Macros are the heart of Logi language. They are used to define custom elements, properties, and behaviors.

There are different types of macros:

1. Syntax macros
2. Helper macros
3. Marker macros

### Syntax macros

Syntax macros are top level macros that define new syntax elements in the language.

#### Examples:

##### Entity macro

Macro definition:

```logi
macro entity {
    kind Syntax
    
    syntax {
       <property Name> <type Type> [id, required, default: <defaultValue>, unique]
    }
}
```

Usage:

```logi
entity User {
    id: int [id, required]
    name: string [required, default: "John Doe"]
    email: string [transient]
}

entity Product {
    id: int [id, required]
    name: string [required]
    price: float [required]
    owner: User [required]
    description: string
}
```

##### Formula macro

Let's imagine, we want to store formulas in our system, later we can use them in different places.

```logi
macro formula {
    kind Syntax
    
    syntax {
       args (...[<argName Name> <Type<Number>]>)
       expr Expression
    }
}
```

Usage:

```logi
formula AddTwoNumbers {
    args (a int, b int)
    expr a + b
}
```

##### AI model macro

Let's imagine, we want to store AI models in our system, later we can use them in different places.

```logi
macro model {
    kind Syntax
    
    defintion {
       ModelKind <enum[Sequential, Recurrent]>
       Activation <enum[Relu, Sigmoid, Softmax]>
       
       Dense (<length int>, <inputShape int[]>, <activation Activation>)
       Conv2D (<filters int>, <kernelSize int[]>, <activation Activation>)
       MaxPooling2D (<poolSize int[]>)
       
       Layer <oneOf[Dense, Conv2D, MaxPooling2D]>
    }
    
    syntax {
       kind ModelKind
       layers Layer[]
       output Layer
    }
}
```

Usage:

```logi
model MyModel {
    kind Sequential
    layers [
        Dense(128, [784], RELU),
        Dense(10, [128], SOFTMAX)
    ]
    output Dense(10, [128], SOFTMAX)
}
```

##### Html macro
You can define html elements with help of macros in Logi.

```logi
macro html {
    kind Syntax
    
    defintion {
        Head {
            Title <string>
            Meta {
                Charset <string>
                Name <string>
                Content <string>
            }
        }

        Body {
          ...Element
        }
        
        H1 {
            ...Element
        }
        
        H1 <string>
        
        // Elements
        RootElement <oneOf[Head, Body]>
        Element <oneOf[H1, P]>
    }
    
    syntax {
       Head [required]
       Body [required]
    }
}
```

Usage:

```logi
html MyPage {
    Head {
        Title "My Page"
        Meta {
            Charset "UTF-8"
            Name "description"
            Content "This is my page"
        }
    }
    Body {
        H1 "Hello, World!"
        P "This is a paragraph."
        
        P {
            "This is another paragraph."
            Div {
                "This is a div."
            }
        }
    }
}
```

#### General Principles for Syntax Macros
                                                                                                                                       

