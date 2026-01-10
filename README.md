# codegen

[![Release](https://img.shields.io/github/release/khevse/codegen.svg)](https://github.com/khevse/codegen/releases/latest)


# Interface generator

```bash
codegen interface \
--type=github.com/khevse/codegen/tests/mainpkg.StructWithMethods=IStructWithMethods \
--target-dir=./tests/mainpkg \
--suffix=generated
```

## Objects wrapper for tests

```bash
codegen object-test-wrapper \
--interface-type=github.com/khevse/codegen/tests/mainpkg.IFactory=FactoryWrapper \
--object-type=github.com/khevse/codegen/tests/mainpkg.Factory \
--target-dir=./tests/mainpkg  \
--mock-package=github.com/khevse/codegen/tests/mainpkg/mocks \
--suffix=generated
```