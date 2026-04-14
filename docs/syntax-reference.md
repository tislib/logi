# Logi Syntax Reference

This document is the complete syntax reference for the Logi language. For a quick introduction, see the [README](../README.md).

## Overview

In Logi, there are two main elements: **macros** and **definitions**. They are defined in separate files (`.lgm` and `.lg` respectively) and loaded separately.

### Macro Kinds

1. **Syntax** -- Defines the syntax of a new language element
2. **Rule** -- Defines a rule that can be applied to a language element
3. **Transform** -- Defines a transformation that can be applied to a language element

---

## Macro Definition

### Structure

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

A corresponding definition:

```logi
{MacroName} {DefinitionName} {
   {SyntaxStatement1}
   {SyntaxStatement2}
}
```

Dynamic parts are enclosed in curly braces. The `syntax` block defines how definitions will be structured. The `types` block is optional and defines complex types that can be used in syntax or other types.

---

## Syntax Statements

Syntax statements are the building blocks of a macro. Each statement is a whitespace-separated sequence of syntax elements.

```logi
macro {MacroName} {
    kind Syntax

    syntax {
        {SyntaxElement1} {SyntaxElement2} {SyntaxElement3}
        {SyntaxElement4} {SyntaxElement5}
    }
}
```

### Matching Rules

- **Statements are unordered** -- the first statement of a macro can match any statement of a definition
- **Statements are optional** -- if a statement is not present in a definition, it is ignored
- **Elements are ordered and required** -- within a statement, the first element of the macro matches the first element of the definition
- **Each statement must be written on a single line**

### Example

```logi
macro person {
    kind Syntax

    syntax {
        firstName <firstName string> lastName <lastName string>
        age <age int>
    }
}
```

Valid definitions:

```logi
// Correct
person JohnDoe {
    firstName "John" lastName "Doe"
    age 30
}

// Correct (statements are unordered)
person JohnDoe {
    age 30
    firstName "John" lastName "Doe"
}

// Correct (statements are optional)
person JohnDoe {
    firstName "John" lastName "Doe"
}
```

Invalid definitions:

```logi
// Wrong: missing required element (lastName)
person JohnDoe {
    firstName "John"
    age 30
}

// Wrong: elements split across lines
person JohnDoe {
    firstName "John"
    lastName "Doe"
    age 30
}
```

---

## Syntax Elements

### 1. Keyword

Keywords define the structural tokens of your DSL. They are plain identifiers used for matching.

```logi
macro person {
    kind Syntax

    syntax {
        name <name string> as known as <knownAs string>
        age <age int> years old
    }
}
```

```logi
person John {
    name "John" as known as "Johnny"
    age 30 years old
}
```

Here `name`, `as`, `known`, `age`, `years`, `old` are keywords that form the syntax structure.

### 2. VariableKeyword

Variable keywords capture dynamic values from definitions. Syntax: `<variableName type>`.

```logi
macro person {
    kind Syntax

    syntax {
        name <name string>
        age <age int>
    }
}
```

```logi
person John {
    name "John"
    age 30
}
```

#### Primitive Types

| Type | Matches | Example |
|------|---------|---------|
| `string` | String values | `name "John"` |
| `int` | Integer values | `age 30` |
| `float` | Float values | `weight 70.5` |
| `bool` | Boolean values | `isMarried true` |
| `date` | Date values | `birthDate "1990-01-01"` |
| `time` | Time values | `birthTime "12:00:00"` |
| `datetime` | Datetime values | `birthDateTime "1990-01-01T12:00:00"` |
| `duration` | Duration values | `workDuration "1h30m"` |
| `money` | Money values | `salary "1000.50 USD"` |
| `unit` | Unit values | `weight "70.5 kg"` |

#### Special Types

**`Name`** -- matches a single unquoted identifier:

```logi
macro config {
    kind Syntax
    syntax {
        Param <paramName Name> <paramValue string>
    }
}

config EngineConfig {
    Param LogLevel "info"
    Param Port "8080"
}
```

**`Type`** -- matches a type as a value:

```logi
macro entity {
    kind Syntax
    syntax {
        property <fieldName Name> <fieldType Type>
    }
}

entity User {
    property name string
    property age int
}
```

**Custom types** -- defined in the `types` section:

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

entity User {
    Prop name string 50
    Prop age int 4
}
```

### 3. ParameterList

Defines a comma-separated list of values inside parentheses.

```logi
macro person {
    kind Syntax

    syntax {
        name <name string>
        Personal (<address string>, <phone string>)
    }
}

person John {
    name "John"
    Personal ("New York", "1234567890")
}
```

Parameter lists can also be used in the `types` section:

```logi
macro person {
    kind Syntax

    types {
        Address (<city string>, <country string>) at (<street string>, <zip int>)
    }

    syntax {
        Personal <address Address>
    }
}

person John {
    Personal Address ("New York", "USA") at ("Wall Street", 10005)
}
```

### 4. Scope

Scopes define syntax for nested blocks of code. They are declared in a `scopes` block and referenced in `syntax` with `{ scopeName }`.

```logi
macro circuit {
    kind Syntax

    syntax {
        components { components }
        actions { command | handler }
    }

    scopes {
        components {
            Led     <component Name> <pin int>
            Button  <component Name> <pin int>
        }
        command {
            on(<component Name>)
            off(<component Name>)
            if (<condition bool>) { command | handler }
            if (<condition bool>) { command | handler } else { command | handler }
        }
        handler {
            on_click(<component Name>) { command }
            while_held(<component Name>) { command }
        }
    }
}
```

```logi
circuit simple1 {
    components {
        Led     yellowLed 5
        Button  button1 17
    }

    actions {
        on(yellowLed)

        on_click(button1) {
            if (status(button2) == 'on') {
                on(yellowLed)
            } else {
                off(yellowLed)
            }
        }
    }
}
```

### 5. ArgumentList

Defines a property list inside parentheses (names + types, not values).

```logi
macro simpleInterface {
    kind Syntax

    syntax {
        <methodName Name> (...[<args Type<string>>]) <returnType Type>
    }
}

simpleInterface UserService {
    createUser (name string, age int) User
}
```

### 6. AttributeList

Defines optional attributes inside square brackets.

```logi
macro person {
    kind Syntax

    syntax {
        <property Name> <type Type> [required boolean, default string]
    }
}

person John {
    name string [required]
    age int [required, default "30"]
}
```

### 7. TypeReference

Includes a type from the `types` section inline in the syntax.

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

person JohnDoe {
    name "John" "Doe"
    age 30
}
```

### 8. Combination

An OR of multiple syntax elements.

```logi
macro person {
    kind Syntax

    syntax {
        name (<firstName string> | <firstName Name>) (<lastName string> | <lastName Name>)
        age <age int>
    }
}

// Both are valid:
person JohnDoe { name "John" "Doe" ... }
person JohnDoe { name John Doe ... }
```

---

## Comments

```logi
// Single line comment

/* Multi-line comment
   spanning multiple lines */
```

Comments can be used in both macro and definition files.
