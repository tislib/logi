Logi — A language for abstraction
=================================================

# Overview

Logi is a language for abstraction that allows developers to define systems declaratively.
It is designed to be simple, flexible, and extensible, making it easy to create complex systems with minimal effort.

# Motivation

With the help of Logi, you can create a custom language for your specific system.

Logi is designed to be a language for abstraction that allows developers to define systems declaratively.

## Example 1. Define a DSL for a credit rules

```logi
macro creditRule {
    kind Syntax
    
    syntax {
        creditScore <minCreditScore int> <maxCreditScore int>
        income <minIncome int> <maxIncome int>
        age <minAge int> <maxAge int>
    }
}
```

Now let's define a credit rule using the macro:

```logi
creditRule Rule1 {
    creditScore 500 600
    income 20000 30000
    age 18 65
}

creditRule Rule2 {
    creditScore 600 700
    income 30000 40000
    age 18 65
}
```

Now in golang, you can read the rules and apply them to the user data.

```go
var g, err = New(tt.options...)
if err != nil {
log.Fatal(err)
}

g.LoadMacroFile("creditRule.lgm")
definitions,err := g.LoadLogiFile("creditRules.logi")
if err != nil {
log.Fatal(err)
}

// definitions is a list of credit rules, which can be used by your app to apply the rules to the user data.
```

## Example 2. Define a DSL for a chatbot

```logi
macro chatbot {
    kind Syntax
    
    syntax {
        intent <name string> {
            pattern <pattern string>
            response <response string>
        }
    }
}
```

Now let's define a chatbot using the macro:

```logi
chatbot MyChatbot {
    intent Greeting {
        pattern "Hello"
        response "Hi there!"
    }
    
    intent Farewell {
        pattern "Goodbye"
        response "See you later!"
    }
}
```

Now in golang, you can read the chatbot definition and use it to create a chatbot.

```go
var g, err = New(tt.options...)
if err != nil {
log.Fatal(err)
}

g.LoadMacroFile("chatbot.lgm")
definitions,err := g.LoadLogiFile("chatbot.logi")
if err != nil {
log.Fatal(err)
}

// definitions is a list of chatbot intents, which can be used by your app to create a chatbot.
```

# Syntax

In Logi, there are two main elements: macros and definitions. They are defined in separate files, and separately loaded.

## Overview

There are different kind of macros in Logi:

1. Syntax: Defines the syntax of a new language element.
2. Rule: Defines a rule that can be applied to a language element.
3. Transform: Defines a transformation that can be applied to a language element.

## Syntax Macros

### Macro Definition

#### A Macro

```logi
macro {MacroName} {
    kind Syntax
    
    types {
        {TypeName1} {SyntaxStatement1}
        {TypeName2} {SyntaxStatement2}
    }
    
    syntax {
        {SyntaxStatement1}
        {SyntaxStatement2}
    }
}
```

#### A Definition

```logi
{MacroName} {DefinitionName} {
   {SyntaxStatement1}
   {SyntaxStatement2}
}
```

In Macro definition, dynamic parts are enclosed in curly braces. For example, {MacroName} is the name of the macro,
{TypeName1} is the name of the type, and {SyntaxStatement1} is the syntax statement.

The main thing in macro definition is syntax, which defines how its definitions will be.

Types are optional, it is used to define complex types which can be used in syntax or types

**Simple Macro definition**

Let's define a simple macro that defines a person with a name and age.

```logi
macro person {
    kind Syntax
    
    syntax {
        name <name string>
        age <age int>
    }
}
```

Now, let's define a person definition using the macro:

```logi
person John {
    name "John"
    age 30
}
```

It will be matched by **_Logi engine_** and will be translated into a definition like this:

```json
{
  "name": "John",
  "age": 30
}
```

### Syntax Statements

Syntax statements are the building blocks of a macro. They define the structure of the macro definition.

Each **_syntax statement_** is white space separated **_syntax element_**(s).

```logi
macro {MacroName} {
    kind Syntax
    
    syntax {
        {SyntaxElement1} {SyntaxElement2} {SyntaxElement3}
        {SyntaxElement4} {SyntaxElement5}
    }
}
```

When Logi engine reads a definition, it matches each statement of definition to the macro statement.
Statements are unordered, it means that first statement of macro can match with any statement of definition.
Statements are also optional, it means that if a statement is not present in definition, it will be ignored.

In other hand, syntax elements are by default required and ordered. It means that first element of macro statement will
match with first element of definition statement.

Each statement itself must be written in a single line.

Person Macro:

```logi
macro person {
    kind Syntax
    
    syntax {
        firstName <firstName string> lastName <lastName string>
        age <age int>
    }
}
```

Person Definitions:

Correct:

```logi
person JohnDoe {
    firstName "John" lastName "Doe"
    age 30
}
```

Correct, but unordered:

```logi
person JohnDoe {
    lastName "Doe" firstName "John"
    age 30
}
```

Incorrect, wrong order of elements:

```logi
person JohnDoe {
    firstName "John" lastName "Doe"
    age 30
}
```

Incorrect, missing element:

```logi
person JohnDoe {
    firstName "John"
    age 30
}
```

Correct, missing statement:

```logi
person JohnDoe {
    firstName "John" lastName "Doe"
}
```

Incorrect, line separated statement:

```logi
person JohnDoe {
    firstName "John"
    lastName "Doe"
    age 30
}
```

### Syntax Elements

There are different type of syntax elements in Logi:

1. **Keyword**: A keyword is an identifier that is used to define a syntax element.
2. **VariableKeyword**: A variable keyword is a keyword that can be used as a variable in the syntax. It is for matching
   dynamic information in definition
3. **ParameterList**: A parameter list is a list of parameters. It is used to define a list of parameters in a syntax
   element.
4. **ArgumentList**: An argument list is a list of arguments. It is used to define a list of arguments in a syntax
   element.
5. **CodeBlock**: A code block is a block of code. It is used to define a block of code in a syntax element.
6. **ExpressionBlock**: An expression block is a block of expression. It is used to define a block of expression in a
   syntax element.
7. **AttributeList**: An attribute list is a list of attributes. It is used to define a list of attributes in a syntax
   element.
8. **TypeReference**: A type reference is to include a statement from *types* section into syntax.
9. **Combination**: A combination is a combination of multiple syntax elements. It is like Or statement
10. **Structure**: A structure is a combination of multiple syntax elements. It is for making nested definition
    statements.

#### Keyword

Keywords are used to define syntax structure. They are just plain identifiers and play roles in matching definition

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        name <name string> as known as <knownAs string>
        age <age int> years old
    }
}
```

Definition

```logi
person John {
    name "John" as known as "Johnny"
    age 30 years old
}
```

As you see, `name`, `as`, `known`, `age`, `years`, `old` are keywords. They are used to form the syntax structure.

#### VariableKeyword

Variable keywords are used to match dynamic information in definition.

Variable keywords are defined as `<variableName type>`. It is a matching variable in the definition based on a type.

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        name <name string>
        age <age int>
    }
}
```

Definition

```logi
person John {
    name "John"
    age 30
}
```

In this example, `name` and `age` are variable keywords. They are used to match dynamic information in the definition.
Based on their types, string is matching "John" and int is matching 30.

List of Variable Keywords types:

##### Primitive types

1. `string`: It is used to match a string value. Example: `name <name string>` will match with `name "John"`
2. `int`: It is used to match an integer value. Example: `age <age int>` will match with `age 30`
3. `float`: It is used to match a float value. Example: `weight <weight float>` will match with `weight 70.5`
4. `bool`: It is used to match a boolean value. Example: `isMarried <isMarried bool>` will match with `isMarried true`
5. `date`: It is used to match a date value. Example: `birthDate <birthDate date>` will match with
   `birthDate "1990-01-01"`
6. `time`: It is used to match a time value. Example: `birthTime <birthTime time>` will match with
   `birthTime "12:00:00"`
7. `datetime`: It is used to match a datetime value. Example: `birthDateTime <birthDateTime datetime>` will match with
   `birthDateTime "1990-01-01T12:00:00"`
8. `duration`: It is used to match a duration value. Example: `workDuration <workDuration duration>` will match with
   `workDuration "1h30m"`
9. `money`: It is used to match a money value. Example: `salary <salary money>` will match with `salary "1000.50 USD"`
10. `unit`: It is used to match a unit value. Example: `weight <weight unit>` will match with `weight "70.5 kg"`

##### Special Types

###### Name type

`Name`: It is used to match a name value.
Example: `name <name Name>` will match with `name John`, It is used to
match a single keyword (without double quotes)

Example:

```logi
macro config {
    kind Syntax
    
    syntax {
        Param <paramName Name> <paramValue string>
    }
}
```

```logi
config EngineConfig {
    Param LogLevel "info"
    Param Port "8080"
}
```

It will be matched by **_Logi engine_** and will be translated into a definition like this:

```json
{
  "LogLevel": "info",
  "Port": "8080"
}
```

###### Type type

`Type`: It is used to match a type value.
Example: `type <typeName Type>` will match with `type int`, It is used to match type as value.

Example:

```logi
macro entity {
    kind Syntax
    
    syntax {
        property <fieldName string> <fieldType Type>
    }
}
```

```logi
entity User {
    property name string
    property age int
}
```

It will be matched by **_Logi engine_** and will be translated into a definition like this:

```json
{
  "name": "string",
  "age": "int"
}
```

###### Types defined in types section

Types can be used in variable keywords. It is used to define complex types which can be used in syntax or types.

Example:

```logi
macro entity {
    kind Syntax
    
    types {
        FieldType Name
    }
    
    syntax {
        field <fieldName FieldType> <fieldValue string>
    }
}
```

```logi
entity User {
    field name "John"
    field age "30"
}
```

It will be matched by **_Logi engine_** and will be translated into a definition like this:

```json
{
  "fieldName": "John",
  "fieldAge": "30"
}
```

Another example:

```logi
macro entity {
    kind Syntax
    
    types {
        Property <name Name> <type Type> <size int>
    }
    
    syntax {
        Prop <prop Property>
    }
}
```

```logi
entity User {
    Prop name string 50
    Prop age int 4
}
```

It will be matched by **_Logi engine_** and will be translated into a definition like this:

```json
{
  "propName": {
    "name": "name",
    "type": "string",
    "size": 50
  },
  "propAge": {
    "name": "age",
    "type": "int",
    "size": 4
  }
}
```

#### ParameterList

Parameter list is used to define a list of parameters inside parentheses.

Example: `Personal (<name String>, <age int>)` will match with `Personal ("John", 30)`,
It is like variable keyword, but working in a list and required to have parenthesis.

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        name <name string>
        age <age int>
        Personal (<address string>, <phone string>)
    }
}
```

Definition

```logi
person John {
    name "John"
    age 30
    Personal ("New York", "1234567890")
}
```

This will be translated as:

```json
{
  "name": "John",
  "age": 30,
  "Personal": {
    "address": "New York",
    "phone": "1234567890"
  }
}
```

Parameter list can also be used in types section.

Macro

```logi
macro person {
    kind Syntax
    
    types {
        Address (<city string>, <country string>) at (<street string>, <zip int>)
    }
    
    syntax {
        name <name string>
        age <age int>
        Personal <address Address>
    }
}
```

Definition

```logi
person John {
    name "John"
    age 30
    Personal Address ("New York", "USA") at ("Wall Street", 10005)
}
```

This will be translated as:

```json
{
  "name": "John",
  "age": 30,
  "Personal": {
    "city": "New York",
    "country": "USA",
    "street": "Wall Street",
    "zip": 10005
  }
}
```

#### ArgumentList

Argument list is used to define a list of arguments inside parentheses.
The difference between parameters and arguments are,
parameters are used to define values, arguments are used to define properties.

Example: `(...[<args Type<string>>])` will match with `(name string, age int)`

Macro

```logi
macro simpleInterface {
     kind Syntax

     syntax {
         <methodName Name> (...[<args Type<string>>]) <returnType Type>
     }
 }
```

Definition

```logi
 simpleInterface UserService {
     createUser (name string, age int) User
 }
```

This will be translated as:

```json
{
  "methodName": "createUser",
  "args": [
    {
      "name": "name",
      "type": "string"
    },
    {
      "name": "age",
      "type": "int"
    }
  ],
  "returnType": "User"
}
```

#### AttributeList

Attribute list is used to define a list of attributes inside square brackets.

Example: `[required boolean]` will match with `[required]`

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        <property Name> <type Type> [required boolean, default string]
    }
}
```

Definition

```logi
person John {
    name string [required]
    age int [required, default "30"]
}
```

This will be translated as:

```json
{
  "name": {
    "type": "string",
    "required": true
  },
  "age": {
    "type": "int",
    "required": true,
    "default": "30"
  }
}
```

#### CodeBlock

Code block is used to define a block of code inside curly braces.

Example: `{ code }` will match with `{return "Hello"}`

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        name <name string>
        age <age int>
        greet { code }
    }
}
```

Definition

```logi
person John {
    name "John"
    age 30
    greet {
      return "Hello"
    }
}
```

This will be translated as:

code block will not be translated to definition immediately, but you can execute them inside VM.

#### ExpressionBlock

Expression block is used to define a block of expression inside curly braces.

Example: `{ expression }` will match with `{1 + 2}`

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        name <name string>
        age <age int>
        greet { expr }
    }
}
```

Definition

```logi
person John {
    name "John"
    age 30
    greet {1 + 2}
}
```

The difference between code block and expression is, inside expression no need to write `return` keyword.

#### TypeReference

Type reference is used to include a statement from *types* section into syntax.

Example: `<TypeName1>` will match TypeName1 from types section

Macro

```logi
macro person {
    kind Syntax
    
    types {
        Name <firstName string> <lastName string>
    }
    
    syntax {
        name <Name>
        age <age int>
    }
}
```

Definition

```logi
person JohnDoe {
    name "John" "Doe"
    age 30
}
```

This will be translated as:

```json
{
  "name": {
    "firstName": "John",
    "lastName": "Doe"
  },
  "age": 30
}
```

TypeReference can be used to simplify complex structures.

#### Combination

Combination is a combination of multiple syntax elements. It is like Or statement.

Example: `(<stringValue string> | <numberValue int>)` will match with `"Hello"` or `123` but not `true`

Macro

```logi
macro person {
    kind Syntax
    
    syntax {
        name (<firstName string> | <firstName Name>) (<lastName string> | <lastName Name>)
        age <age int>
    }
}
```

Definition

```logi
person JohnDoe {
    name "John" "Doe"
    age 30
}
```

Another Definition

```logi
person JohnDoe {
    name John Doe
    age 30
}
```

Both will be translated as:

```json
{
  "name": {
    "firstName": "John",
    "lastName": "Doe"
  },
  "age": 30
}
```

#### Structure

Structure is used to define nested statements.

Example: `{<name string> <age int> Nested2 { <nestedName string> <nestedAge int> }}` will match with
`{"John" 30 { "Jane" 25 }}`

Macro

```logi
macro user {
    kind Syntax
    
    syntax {
        name <name string>
        age <age int>
        Auth {
            username <username string>
            password <password string>
            Token {
                accessToken <accessToken string>
                refreshToken <refreshToken string>
            }
        }
    }
}
```

Definition

```logi
user JohnDoe {
    name "John"
    age 30
    Auth {
        username "john.doe"
        password "password"
        Token {
            accessToken "access token"
            refreshToken "refresh token"
        }
    }
}
```

This will be translated as:

```json
{
  "name": "John",
  "age": 30,
  "Auth": {
    "username": "john.doe",
    "password": "password",
    "Token": {
      "accessToken": "access token",
      "refreshToken": "refresh token"
    }
 }
}
```

### Comments

Comments are used to write notes in the code. They are not executed by the engine.

There are two types of comments in Logi:
1. Single line comment: It starts with `//` and ends with the end of the line.
2. Multi-line comment: It starts with `/*` and ends with `*/`.

Example:

```logi
// This is a single line comment

/* This is a multi-line comment
   It can be used to write multiple lines of comments
*/
```

Comments can be used in both macros and logi files.

