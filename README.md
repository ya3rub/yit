# About Yit (Yet another git)

  "Yit" is a lightweight, simplified clone of the popular version control system Git, implemented in the Go programming language.

Features Implemented:

- [x]  Commit
- [x] checkout
- [x] branch
- [x] init
- [x] log
- [x] tag
- [ ] merge
- [ ] push

# Try it !

#### Init

```
  ./yit init
```

#### commit

```
  ./yit commit -m "[MESSAGE]"
```

#### log

```
  ./yit log
```

#### branch  (adding branch)

if start commit not provided, it will look for `HEAD`

```
  ./yit commit -s "[START_COMMIT]" -n "{BRANCH_NAME}"
```

#### Checkout

if start commit not provided, it will look for `HEAD`

> NOTE: dist dir is required for safety reasons..

```
  ./yit checkout -b "[BRANCH_NAME]" -n -d "[DIST_DIR]"
```

#### tag

if start commit not provided, it will look for `HEAD`

```
  ./yit checkout -c "[COMMIT]" -n -t "[TAG]"
```
