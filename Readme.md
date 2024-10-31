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
    id: int
    name: string
    email: string
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
    syntax {
        args (name Name) 
        block {
            on <eventName Name>() CodeBlock
        }
    }
    
    target {
        allow: [root]
    }
    
    handler {
        var definition Definition
        
        definition.canInstantiate = false
        definition.name = name
        definition.package = this.package
        
        for (onBlock in this.block.on) {
            var event Event
            event.name = onBlock.eventName
            event.code = onBlock.code
            definition.events.add(event)
        }
        
        this.syntax.addDefinition(definition)
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
    <property1>: <value1>
    <property2>: <value2>
    <property3>: {
        <nestedProperty1>: <nestedValue1>
        <nestedProperty2>: <nestedValue2>
    }
}
```

For definition of data type

```
<elementType> <elementName> {
    <property1>: <type1>
    <property2>: <type2>
    <property3>: {
        <nestedProperty1>: <type3>
        <nestedProperty2>: <type4>
    }
}
```

For definition of code block

```
<codeType> <codeName> (param1 type1, param2 type2) returnType {
    ...code...
    return <value>
}
```

## Macros

Macros are the heart of Logi language. They are used to define custom elements, properties, and behaviors.

```
macro 




