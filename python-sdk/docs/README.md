# Python SDK Documentation

This directory contains the MkDocs documentation for the LLM NPC Python SDK.

## Local Development

### Setup

1. Install uv (if needed):
```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

2. Install documentation dependencies:
```bash
cd python-sdk
uv pip install -e ".[docs]"
```

3. Start the development server:
```bash
mkdocs serve
```

3. Open your browser to `http://localhost:8000`

The site will automatically reload when you make changes to the documentation files.

## Building

To build the static site:

```bash
mkdocs build
```

The built site will be in the `site/` directory.

## Deployment

Documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch:

- **Trigger**: Push to `main` branch with changes in `python-sdk/` directory
- **Workflow**: `.github/workflows/deploy-docs.yml`
- **Deployed to**: `gh-pages` branch
- **URL**: `https://yourusername.github.io/llm-npc-backend/python-sdk/`

### Manual Deployment

You can also deploy manually:

```bash
mkdocs gh-deploy
```

## Structure

```
docs/
├── index.md                    # Home page
├── getting-started.md          # Getting started guide
├── user-guide/                 # User guides
│   ├── tools.md               # Defining tools
│   ├── npcs.md                # Creating NPCs
│   ├── context.md             # Building context
│   ├── responses.md           # Working with responses
│   └── knowledge-graphs.md    # Knowledge graphs
├── api/                        # API reference (auto-generated)
│   ├── client.md
│   ├── decorators.md
│   ├── models.md
│   ├── context.md
│   └── exceptions.md
├── examples/                   # Examples
│   ├── index.md
│   └── simple-game.md
└── advanced/                   # Advanced topics
    ├── error-handling.md
    └── best-practices.md
```

## Writing Documentation

### Markdown Files

Documentation is written in Markdown with some extensions:

- **Code blocks**: Use triple backticks with language
- **Admonitions**: Use `!!!` for callouts (note, warning, tip, etc.)
- **Tabs**: Use `===` for tabbed content
- **Math**: Use `\(` for inline and `\[` for block math

### API Reference

API reference pages use mkdocstrings to auto-generate from docstrings:

```markdown
::: llm_npc.client.NPCClient
    options:
      show_source: false
      members:
        - __init__
        - health_check
```

### Code Examples

Always include complete, runnable examples:

```python
from llm_npc import NPCClient, tool

@tool
def speak(message: str):
    """Make NPC speak"""
    pass

client = NPCClient("http://localhost:8080")
```

## Style Guide

- Use clear, concise language
- Include code examples for every feature
- Add "See Also" sections to link related pages
- Use admonitions for warnings, tips, and notes
- Test all code examples

## Configuration

The documentation is configured in `mkdocs.yml`:

- **Theme**: Material for MkDocs
- **Plugins**: mkdocstrings for API docs, search
- **Extensions**: Code highlighting, admonitions, tabs, etc.

## Troubleshooting

### Build Errors

If the build fails:

1. Check mkdocs.yml syntax
2. Verify all linked files exist
3. Check for invalid markdown
4. Run `mkdocs build --verbose` for details

### Import Errors

If mkdocstrings can't find modules:

1. Make sure the package is installed: `pip install -e .`
2. Check the module paths in API reference pages
3. Verify docstrings are formatted correctly

## Contributing

When adding new features:

1. Update relevant user guide pages
2. Add API reference if new public APIs
3. Include examples
4. Update navigation in `mkdocs.yml` if needed
5. Test the docs locally before pushing

