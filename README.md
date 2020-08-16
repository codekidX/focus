## 🎯 focus

> Never leave the terminal and focus on what is important for your project.

Focus is a cli tool that provides the all the information and task needed for you
to focus on the development of the tasks and issues at hand.

### Installing focus

**Install script is WIP!!** - for now do:

```bash
go install -u github.com/codekidX/focus
```

### Using focus

Navigate to your GitHub project on your terminal and run:

- list

```bash
focus
```

> this will list all your issues of your current git repository

- on

```bash
focus on 1
```

> this will focus on the issue number **1** of your repository, and will show the description and info

- open

```bash
focus open 1
```

> this will open issue number **1** on your default browser

- help

```bash
focus (help, h, ?)
```

- page

```bash
focus page 2
```

> shows page 2 of the issues list (pagination)

Supports the following functionalities:

- List issues
- View an issue
- Open issue
- Create an issue
- Shows basic info on the repository
