# Commit Message Instructions

- Use conventional commit message format.
- The commit message should have just a short description (80 characters or less) as a title.
- The short description should be in the format: `<type>(<scope>):<icon> <short description>`
  - `type`: The type of change (e.g., feat, fix, docs, style, refactor, test, chore).
    - `feat`: ✨ A new feature
    - `fix`: 🐛 A bug fix
    - `docs`: 📝 Documentation only changes
    - `style`: 💄 Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
    - `refactor`: ♻️ A code change that neither fixes a bug nor adds a feature
    - `test`: ✅ Adding missing tests or correcting existing tests
    - `chore`: 🔧 Changes to the build process or auxiliary tools and libraries such as documentation generation
    - `perf`: ⚡️ A code change that improves performance
    - `ci`: 👷 Changes to CI configuration files and scripts
    - `build`: 🏗️ Changes that affect the build system or external dependencies
    - `revert`: ⏪ Reverts a previous commit
    - `wip`: 🚧 Work in progress
    - `security`: 🔒 Security-related changes
    - `i18n`: 🌐 Internationalization and localization
    - `a11y`: ♿ Accessibility improvements
    - `ux`: 🎨 User experience improvements
    - `ui`: 🖌️ User interface changes
    - `config`: 🔧 Configuration file changes
    - `deps`: 📦 Dependency updates
    - `infra`: 🌐 Infrastructure changes
    - `init`: 🎉 Initial commit
    - `analytics`: 📈 Analytics or tracking code
    - `seo`: 🔍 SEO improvements
    - `legal`: ⚖️ Licensing or legal changes
    - `typo`: ✏️ Typo fixes
    - `comment`: 💬 Adding or updating comments in the code
    - `example`: 💡 Adding or updating examples
    - `mock`: 🤖 Adding or updating mocks
    - `hotfix`: 🚑 Critical hotfix
    - `merge`: 🔀 Merging branches
    - `cleanup`: 🧹 Code cleanup
    - `deprecate`: 🗑️ Deprecating code or features
    - `move`: 🚚 Moving or renaming files
    - `rename`: ✏️ Renaming files or variables
    - `split`: ✂️ Splitting files or functions
    - `combine`: 🧬 Combining files or functions
    - `add`: ➕ Adding files or features
    - `remove`: ➖ Removing files or features
    - `update`: ⬆️ Updating files or features
    - `downgrade`: ⬇️ Downgrading files or features
    - `patch`: 🩹 Applying patches
    - `optimize`: 🛠️ Optimizing code
    - `docs`: 📝 Documentation changes
    - `test`: ✅ Adding or updating tests
    - `fix`: 🐛 Bug fixes
    - `style`: 💄 Code style changes (formatting, etc.)
    - `refactor`: ♻️ Code refactoring
    - `perf`: ⚡️ Performance improvements
    - `ci`: 👷 Continuous integration changes
    - `build`: 🏗️ Build system changes
    - `revert`: ⏪ Reverting changes
    - `wip`: 🚧 Work in progress
    - `security`: 🔒 Security improvements
    - `i18n`: 🌐 Internationalization changes
    - `a11y`: ♿ Accessibility improvements
    - `ux`: 🎨 User experience improvements
    - `ui`: 🖌️ User interface changes
    - `config`: 🔧 Configuration changes
    - `deps`: 📦 Dependency updates
    - `infra`: 🌐 Infrastructure changes
    - `init`: 🎉 Initial commit
    - `analytics`: 📈 Analytics changes
    - `seo`: 🔍 SEO improvements
    - `legal`: ⚖️ Legal changes
    - `typo`: ✏️ Typo fixes
    - `comment`: 💬 Comment changes
    - `example`: 💡 Example changes
    - `mock`: 🤖 Mock changes
    - `hotfix`: 🚑 Hotfix changes
    - `merge`: 🔀 Merge changes
    - `cleanup`: 🧹 Cleanup changes
    - `deprecate`: 🗑️ Deprecation changes
    - `move`: 🚚 Move changes
    - `rename`: ✏️ Rename changes
    - `split`: ✂️ Split changes
    - `combine`: 🧬 Combine changes
    - `add`: ➕ Add changes
    - `remove`: ➖ Remove changes
    - `update`: ⬆️ Update changes
    - `downgrade`: ⬇️ Downgrade changes
    - `patch`: 🩹 Patch changes
    - `optimize`: 🛠️ Optimize changes
  - `scope`: The scope of the change (e.g., component or file name). Include this if the change is specific to a particular part of the codebase.
- `short description`: A brief summary of the change.

### Commit Message Example

```
feat(auth): ✨ Add user authentication

```

## Pull Request Description Instructions

- Use a clear and concise title following the commit message format.
- Provide a detailed description including:
  - What changes were made and why.
  - List of key changes (use bullet points).
  - Reference related issues (e.g., Closes #123).
  - Testing instructions if applicable.
  - Screenshots or examples if relevant.

### Pull Request Description Example

```
## Title
feat(pallet): ✨ Implement QR code generation for pallets

## Description
Add QR code functionality to improve inventory tracking by allowing quick access to pallet details via scanning.

## Changes
- Add QR code service for generating codes
- Store QR codes in static directory
- Update pallet creation endpoint to generate QR

## Testing
- Create a new pallet and verify QR code is generated
- Scan the QR code to ensure it links correctly

Closes #45
```
