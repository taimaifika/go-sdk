data "external_schema" "gorm" {
  program = ["go", "run", "./migrate.go"]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/dev"
  url = "postgres://postgres:1zMPM23hFq1@localhost:5432/todo-list?sslmode=disable"
  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
