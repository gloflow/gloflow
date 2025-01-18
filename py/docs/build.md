


build package and publish to PyPi:
```
python3 -m build
twine upload dist/*
```


update local gloflow package from PyPi public repo:
```
pip3 install --upgrade gloflow
```


install localy in `editable` mode.
links the package directly to the source code directory
```
pip3 install -e .
```