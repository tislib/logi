Logi — A language for abstraction
=================================================

# Overview

Logi is a language for abstraction that allows developers to define systems declaratively.
It is designed to be simple, flexible, and extensible, making it easy to create complex systems with minimal effort.

# Motivation

With the help of Logi, you can create a custom language for your specific system.

Logi is designed to be a language for abstraction that allows developers to define systems declaratively.

You can download binaries from release page, or you can build from source.

## Example 1. Define a DSL for a credit rules

credit-rule.lgm

```logi
macro creditRule {
    kind Syntax
    
    syntax {
        creditScore <min int> <max int>
        income <min int> <max int>
        age <min int> <max int>
    }
}
```

Now let's define a credit rule using the macro:

credit-rule.lg

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

Now you can easily compile the Logi file and use the definitions in your application.

```shell
logi compile -i engine-config.lg
```

This will parse macro and logi definition file, and will compile it to json, which can be used by your application.

[see result](examples/credit-rule/credit-rule.json)

See examples folder for all examples.

## Example 2. Define a DSL for a chatbot

```logi
macro chatbot {
    kind Syntax
    
    syntax {
        intent <name Name> {
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

Logi language also have Virtual Machine where you can load, read, execute the Logi content.

[see main.go](examples/chat-bot/main.go)

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
  "age": {
    "age": 30
  },
  "name": {
    "name": "John"
  }
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

There are different types of syntax elements in Logi:

1. **Keyword**: A keyword is an identifier that is used to define a syntax element.
2. **VariableKeyword**: A variable keyword is a keyword that can be used as a variable in the syntax. It is for matching
   dynamic information in definition
3. **ParameterList**: A parameter list is a list of parameters. It is used to define a list of parameters in a syntax
   element.
4. Scope: A scope is a block of code enclosed in curly braces. It is used to define a block of code in a syntax element.
5. **ArgumentList**: An argument list is a list of arguments. It is used to define a list of arguments in a syntax
   element.
   syntax element.
6. **AttributeList**: An attribute list is a list of attributes. It is used to define a list of attributes in a syntax
   element.
7. **TypeReference**: A type reference is to include a statement from *types* section into syntax.
8. **Combination**: A combination is a combination of multiple syntax elements. It is like Or statement

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

it will be translated as:

```json
{
  "name": {
    "name": "John",
    "knownAs": "Johnny"
  },
  "age": {
    "age": 30
  }
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
  "paramLogLevel": {
    "paramName": "LogLevel",
    "paramValue": "info"
  },
  "paramPort": {
    "paramName": "Port",
    "paramValue": "8080"
  }
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
        property <fieldName Name> <fieldType Type>
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
  "propertyAge": {
    "fieldName": "age",
    "fieldType": "int"
  },
  "propertyName": {
    "fieldName": "name",
    "fieldType": "string"
  }
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

#### Scope
Scopes are for defining syntax for nested blocks of code. And reusing them in syntax.

An example of scope is:

```logi
circuit simple1 {
    components {
        Led 	yellowLed 5
        Led 	redLed 6
        Led 	blueLed 13
        Button 	button1 17
        Button 	button2 19
    }

    actions {
        on(yellowLed)
        on(redLed)

        on_click(button1) {
            if (status(button2) == 'on') {
                on(blueLed)
                on(yellowLed)
                on(redLed)
            } else {
                off(blueLed)
            }
        }

        on_click(button2) {
            off(blueLed)
            on(yellowLed)
            on(redLed)
        }
    }
}
```
In this logi file, you can see nested blocks of code. `components` and `actions` are scopes. 
Now, let's define macro file for it.

```logi
macro circuit {
    kind Syntax

    syntax {
        components { components }
        actions { command | handler }
    }

    scopes {
        components {
            Led 	<component Name> <pin int>
            Button 	<component Name> <pin int>
        }
        command {
            // Basic commands
            on(<component Name>)
            off(<component Name>)
            blink(<component Name>, <count int>, <seconds float>)
            wait(<seconds float>)
            brightness(<component Name>, <value float>)
            fade_in(<component Name>, <seconds float>)
            fade_out(<component Name>, <seconds float>)
            // Conditional commands
            if (<condition bool>) { command | handler }
            if (<condition bool>) { command | handler } else { command | handler }
        }
        handler {
            // Event handlers
            on_click(<component Name>) { command }
            on_click(<component Name>, <count int>) { command }
            on_press(<component Name>, <count int>) { command }
            on_release(<component Name>, <count int>) { command }
            while_held(<component Name>) { command }
        }
    }
}
```
As you can see from this code, scopes are defined in scopes block.
```logi
...
syntax {
   {command} { scopeName }
}
scopes {
   {scopeName} {
       {Statement1}
       {Statement2}
       ...  
   }
}
```
And it can be used in anyplace with { scopeName } syntax element.

Later to execute this code, you either need to execute it from its json compiled form: [see result](examples/circuit/circuit-1.json)
Or in golang, you can use virtual machine and implementer: [see main.go](pkg/vm/vm_test.go)


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

