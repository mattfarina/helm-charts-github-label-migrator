# Helm Charts GitHub Label Migrator

Problem: When we first added the stale issues to the helm charts repo it added the "wontfix" label. This is not the label we use for that (oops).

Solution: The project uses the GitHub API to add the right label.

It also provides a template for working with the GitHub API for issues.