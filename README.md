# codegen

[![Release](https://img.shields.io/github/release/khevse/codegen.svg)](https://github.com/khevse/codegen/releases/latest)


# Interface generator

```bash
bin/codegen interface \
--type=github.com/khevse/codegen/tests/mainpkg.StructWithMethods=IStructWithMethods \
--target-dir=./internal/command/interface_creator \
--suffix=_generated
```

## Objects wrapper for tests

```bash
bin/codegen object-test-wrapper \
--interface-type=github.com/khevse/codegen/tests/mainpkg.IFactory=FactoryWrapper \
--target-dir=./internal/command/object_test_wrapper  \
--mock-package=github.com/khevse/codegen/tests/mainpkg/mocks \
--suffix=_generated
```