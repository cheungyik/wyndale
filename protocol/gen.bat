@echo off
for %%f in (*.proto) do (
  if exist "%%f" (
    protoc --go_out=..\..\ %%f
    protoc --csharp_out=.\csharp %%f
  )
)
