# Installation

download latest binary from here:
 - Linux
   - x32-64 ()
 - Windows
   - asd

# Setting up

to create executable accessible from anywhere type:

```bash
echo -e "\n$(printf 'PATH=\"$PATH:%s\"' $(pwd))\n" >> ~/.bashrc
```

> or change `$(pwd)` to your installation path
> > `$pwd` pastes your current directory 

# Usage

| flags             | data | description       | 
|-------------------|------|-------------------|
| <no flag\>        |      | logs saved data   |
| `-a` _or_ `--add` |      | adds new          |
|                   |      |                   |
