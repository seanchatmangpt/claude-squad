---
description: Generate or update project documentation from code
allowed-tools: [Bash, Read, Glob, Edit]
model: claude-sonnet-4-5-20250929
---

# Generate Documentation

## Documentation Tasks

1. **Generate Godoc**
   - Review all public functions and types
   - Ensure each has descriptive comments
   - Generate: `godoc -html github.com/smtg-ai/claude-squad > docs/api.html`

2. **API Documentation**
   - Extract public API from main modules
   - Create structured API reference
   - Document parameters, return values, errors

3. **Architecture Documentation**
   - Identify major components
   - Document data flow
   - Explain key design decisions

4. **Example Documentation**
   - Create runnable examples
   - Document common workflows
   - Include expected outputs

5. **Update README**
   - Add installation instructions
   - Include quick start guide
   - List key features
   - Add troubleshooting section

6. **Changelog**
   - Document recent changes
   - List breaking changes
   - Note migration paths

7. **Configuration Reference**
   - Document all config options
   - Provide examples
   - Explain defaults

## Output

- Update relevant markdown files
- Ensure all links work
- Verify code examples are current
- Check formatting consistency
