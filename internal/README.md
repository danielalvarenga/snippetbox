In `internal` folder we have all Go code that can be shared among the applications inside the project, like validations, helpers, database models, etc.

In Go, the code inside the `internal` folder can be imported only by the code inside the parent folder of the `internal` folder, even the project is public on github.