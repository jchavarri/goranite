# Goranite Example Site

This directory contains a simple example site that demonstrates the capabilities of Goranite.

## Structure

```
example/
├── config.json          # Site configuration
├── content/
│   └── posts/           # Blog posts in Markdown
├── static/
│   ├── css/            # Custom CSS files
│   └── images/         # Static images
└── public/             # Generated output (created after build)
```

## Building the Example

From the goranite directory, run:

```bash
./goranite -build -site example
```

## Serving the Example

To start a development server with live reload:

```bash
./goranite -serve -site example
```

Then visit http://localhost:8080 to see the site.

## Customizing

- Edit `config.json` to change site metadata
- Add new posts to `content/posts/`
- Modify `static/css/custom.css` to customize the styling
- Add images to `static/images/`

## Post Format

Posts should be written in Markdown with JSON frontmatter:

```markdown
---
{
  "title": "Your Post Title",
  "date": "2024-01-15T00:00:00Z",
  "tags": ["tag1", "tag2"],
  "summary": "A brief summary of your post"
}
---

# Your Post Content

Write your content here in Markdown...
``` 