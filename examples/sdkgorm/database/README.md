# Database chema as code with GORM
## Inspecting the schema
```bash
atlas schema inspect -u [uri] > schema.hcl
```
Example:
```bash
atlas schema inspect -u "postgres://postgres:1zMPM23hFq1@localhost:5432/todo-list?sslmode=disable" > schema.hcl
```

## Apply the schema
```bash
atlas schema apply -u [uri] -f schema.hcl
```
Example:
```bash
atlas schema apply -u "postgres://postgres:1zMPM23hFq1@localhost:5432/todo-list?sslmode=disable" -f schema.hcl
```

## Generate migration diff
```bash
atlas migrate diff -u [uri] -f schema.hcl
```
Example:
```bash
atlas migrate diff -u "postgres://postgres:1zMPM23hFq1@localhost:5432/todo-list?sslmode=disable" -f schema.hcl
```
