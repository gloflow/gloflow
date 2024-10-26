from setuptools import setup, find_packages

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

with open(f'{modd_str}/requirements.txt') as f:
    required_packages = f.read().splitlines()

authors_lst = [
    "Ivan Trajkovic"
]

setup(
    name="gloflow",
    version="0.1.11",
    author=",".join(authors_lst),
    author_email="glofloworg@gmail.com",
    description="""
Py package for interacting with the gloflow platform API's.
""",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/gloflow/gloflow",
    packages=find_packages(
        where="src",
        exclude=[
            "src/deprecated"
        ]),
    classifiers=[
        "Programming Language :: Python :: 3",
        "Operating System :: POSIX :: Linux",
    ],
    python_requires=">=3.6",

    install_requires=required_packages,
    # install_requires=[
    #    "requests",
    # ],

    package_dir={"": "src"},  # Base directory for packages is src/
)
