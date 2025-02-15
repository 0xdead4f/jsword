![image_2025-02-14_18-46-19](https://github.com/user-attachments/assets/8327071c-5024-4bea-9c1b-190eca4b8b36)

# Jsword
Build custom wordlist by extracting word from remote/local javascript or html page.

## Installation

```bash
go install github.com/0xdead4f/jsword@latest
```

## Usage

```bash
jsword <path|url>
```

Examples:
```bash
jsword script.js
jsword https://example.com/script.js
```

## It Covers

JavaScript Detection:
- Function declarations (function myFunc)
- Class declarations (class MyClass)
- Method declarations (myMethod() {})
- Object property access (obj.property)
- Property names (key: value)
- Const declarations (const myVar = ...)

HTML/Template Detection:
- ID attributes with quotes (id="myId")
- ID attributes without quotes (id=myId)
- Name attributes with quotes (name="myName")
- Name attributes without quotes (name=myInput)

Template Literal Support:
- Handlebars/Angular: {{variable}}
- React/JSX: {variable}
- ES6 Template: ${variable}
- JSP: <%variable%>
- Other common syntaxes: [[variable]], @{variable}, ((variable))

Object Pattern Detection:
- Property names in objects ({property: value})
- Nested object identifiers
- Property-value pairs
- JSX attribute patterns

