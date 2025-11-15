# Python SDK Documentation Setup Guide

The Python SDK now has comprehensive documentation built with MkDocs and Material theme!

## 🎉 What Was Created

### Documentation Structure

```
python-sdk/
├── docs/                           # Documentation source
│   ├── index.md                   # Home page
│   ├── getting-started.md         # Getting started guide
│   ├── user-guide/                # Comprehensive user guides
│   │   ├── tools.md              # Defining tools
│   │   ├── npcs.md               # Creating NPCs
│   │   ├── context.md            # Building context
│   │   ├── responses.md          # Working with responses
│   │   └── knowledge-graphs.md   # Knowledge graphs
│   ├── api/                       # Auto-generated API docs
│   │   ├── client.md
│   │   ├── decorators.md
│   │   ├── models.md
│   │   ├── context.md
│   │   └── exceptions.md
│   ├── examples/                  # Examples
│   │   ├── index.md
│   │   └── simple-game.md
│   └── advanced/                  # Advanced topics
│       ├── error-handling.md
│       └── best-practices.md
├── mkdocs.yml                     # MkDocs configuration
└── pyproject.toml                 # Updated with docs dependencies
```

### GitHub Actions

- **File**: `.github/workflows/deploy-docs.yml`
- **Triggers**: Push to `main` branch with changes in `python-sdk/`
- **Deploys to**: GitHub Pages (`gh-pages` branch)

## 🚀 Quick Start

### 1. Install uv (if you haven't already)

```bash
# macOS and Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"

# Or with pip
pip install uv
```

### 2. Install Dependencies

```bash
cd python-sdk
uv pip install -e ".[docs]"
```

This installs:
- `mkdocs` - Documentation generator
- `mkdocs-material` - Material theme
- `mkdocstrings` - Auto API docs from docstrings
- `pymdown-extensions` - Markdown extensions

### 3. Preview Locally

```bash
mkdocs serve
```

Then open `http://localhost:8000` in your browser.

The site will auto-reload when you edit documentation files!

### 4. Build Static Site

```bash
mkdocs build
```

Outputs to `site/` directory.

## 📤 Deploying to GitHub Pages

### Automatic Deployment (Recommended)

Documentation auto-deploys when you push to `main`:

1. Commit your changes:
```bash
git add .
git commit -m "docs: update Python SDK documentation"
git push origin main
```

2. GitHub Actions will:
   - Build the documentation
   - Deploy to `gh-pages` branch
   - Make it available at: `https://piercegov.github.io/llm-npc-backend/python-sdk/`

3. Enable GitHub Pages:
   - Go to your repo settings
   - Navigate to **Pages**
   - Source: `gh-pages` branch, `/ (root)` directory
   - Save

### Manual Deployment

```bash
mkdocs gh-deploy
```

This builds and pushes directly to `gh-pages` branch.

## 🎨 Documentation Features

### Material Theme

- **Dark/light mode** toggle
- **Search** functionality
- **Navigation** tabs and sections
- **Mobile-friendly** responsive design
- **Code** copy buttons

### Auto-Generated API Docs

API reference pages automatically extract documentation from your Python docstrings:

```python
def my_function(param: str):
    """
    Do something useful.
    
    Args:
        param: Description of the parameter
    
    Returns:
        Description of return value
    """
    pass
```

### Rich Content

- ✅ Code syntax highlighting
- ✅ Tabbed code blocks
- ✅ Admonitions (notes, warnings, tips)
- ✅ Cross-references between pages
- ✅ Table of contents
- ✅ Responsive design

## 📝 Editing Documentation

### Adding a New Page

1. Create a markdown file in `docs/`:
```bash
touch docs/user-guide/my-new-feature.md
```

2. Add it to navigation in `mkdocs.yml`:
```yaml
nav:
  - User Guide:
      - My New Feature: user-guide/my-new-feature.md
```

3. Write your content using Markdown

4. Preview with `mkdocs serve`

### Using Admonitions

```markdown
!!! note "Optional Title"
    This is a note block.

!!! warning
    This is a warning.

!!! tip
    This is a helpful tip.
```

### Using Tabs

```markdown
=== "Python"
    ```python
    print("Hello")
    ```

=== "JavaScript"
    ```javascript
    console.log("Hello");
    ```
```

### Code Blocks

```markdown
```python
from llm_npc import NPCClient

client = NPCClient("http://localhost:8080")
```
```

## 🔧 Customization

### Theme Colors

Edit `mkdocs.yml`:

```yaml
theme:
  palette:
    primary: indigo  # Change this
    accent: purple   # Change this
```

### Site Information

Update in `mkdocs.yml`:

```yaml
site_name: LLM NPC Python SDK
site_description: Your description
site_author: Your name
site_url: https://your-url.com
```

### Repository Links

Update in `mkdocs.yml`:

```yaml
repo_name: piercegov/llm-npc-backend
repo_url: https://github.com/piercegov/llm-npc-backend
```

## 🐛 Troubleshooting

### Build Fails

```bash
# Verbose output
mkdocs build --verbose --strict

# Check for broken links
mkdocs build --strict
```

### Can't Import SDK

Make sure it's installed:

```bash
uv pip install -e .
```

### Module Not Found in API Docs

Check that:
1. The module path is correct in the `.md` file
2. The SDK is installed
3. Docstrings are present and properly formatted

## 📚 Documentation Best Practices

1. **Keep it current**: Update docs when you change code
2. **Include examples**: Every feature should have a code example
3. **Link related pages**: Use "See Also" sections
4. **Test examples**: Make sure code examples actually work
5. **Use admonitions**: Highlight important information
6. **Clear headings**: Use descriptive section titles
7. **Search-friendly**: Use keywords users might search for

## 🔗 Useful Links

- **MkDocs Documentation**: https://www.mkdocs.org/
- **Material Theme**: https://squidfunk.github.io/mkdocs-material/
- **mkdocstrings**: https://mkdocstrings.github.io/
- **Markdown Guide**: https://www.markdownguide.org/

## Next Steps

1. **Update URLs**: Replace `piercegov` in:
   - `mkdocs.yml`
   - `pyproject.toml`
   - GitHub Actions workflow

2. **Enable GitHub Pages**: 
   - Push to trigger deployment
   - Configure Pages in repo settings

3. **Customize**: 
   - Change theme colors
   - Add your branding
   - Update social links

4. **Write More Docs**:
   - Add tutorials
   - Document edge cases
   - Include troubleshooting guides

Enjoy your beautiful new documentation! 🎊

