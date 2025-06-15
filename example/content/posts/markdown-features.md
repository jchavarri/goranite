---
{
  "title": "Markdown Features Demo",
  "date": "2024-01-20T00:00:00Z",
  "tags": ["markdown", "demo", "features"],
  "summary": "A comprehensive demonstration of various Markdown features supported by Goranite."
}
---

# Markdown Features Demo

This post demonstrates various Markdown features supported by Goranite.

## Text Formatting

You can use **bold text**, *italic text*, and even ***bold italic text***.

You can also use `inline code` and ~~strikethrough text~~.

## Lists

### Unordered Lists

- First item
- Second item
  - Nested item
  - Another nested item
- Third item

### Ordered Lists

1. First step
2. Second step
   1. Sub-step A
   2. Sub-step B
3. Third step

## Code Blocks

### JavaScript

```javascript
function greet(name) {
    console.log(`Hello, ${name}!`);
}

greet('Goranite');
```

### Python

```python
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(10))
```

### Bash

```bash
#!/bin/bash
echo "Building with Goranite..."
./goranite -build
echo "Build complete!"
```

## Blockquotes

> This is a blockquote. It can contain multiple paragraphs and other Markdown elements.
> 
> Like this second paragraph with **bold text**.

## Links

Check out the [Goranite repository](https://github.com/jchavarri/goranite) for more information.

## Tables

| Feature | Supported | Notes |
|---------|-----------|-------|
| Markdown | ✅ | Full CommonMark support |
| Syntax Highlighting | ✅ | Via Chroma |
| Live Reload | ✅ | Development server |
| Themes | 🚧 | Coming soon |

## Horizontal Rule

---

That's all for this demo! 