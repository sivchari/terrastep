# terrastep
terrastep is a library which manages execution order of terraform.

```yml
steps:
  - name: 'run a'
    tactics:
      - validate
      - fmt
      - plan
      - apply
    steps:
      - './a'
      - './b'

  - name: 'run c'
    tactics:
      - validate
      - fmt
      - plan
      - apply
    steps:
      - './c'
      - './d'
```
