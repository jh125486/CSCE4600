# Project 2: Shell Builtins

## Description

For this project we'll be adding commands to a simple shell. 

The shell is already written, but you will choose five (5) shell builtins (or shell-adjacent) commands to rewrite into Go, and integrate into the Go shell.

There are many builtins or shell-adjacent commands to pick from: 
[Bourne Shell Builtins](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html), 
[Bash Builtins](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html,), and 
[Built-in csh and tcsh Commands](https://docstore.mik.ua/orelly/linux/lnut/ch08_09.htm).

Feel free to pick from `sh`, `bash`, `csh`, `tcsh`, `ksh` or `zsh` builtins... or if you have something else in mind, ping me and we'll work it out.

As an example, two shell builtins have already been added to the package builtins:

- `cd`
- `env`

## Steps

1. Clone down the example input/output and skeleton `main.go`:

    `git clone https://github.com/jh125486/CSCE4600`
 
2. Copy the `Project2` files to your own git project.

Start editing the `main.go` command switch (lines 57-64) and the package `builtins` with your chosen commands.

## Grading

Code must compile and run.

Each type is worth different points:

- 10 points for each command implemented.
- 50 peer points (points per peer adjusted for group size)

## Deliverables

A GitHub link to your project which includes:

- `README.md` <- describes anything needed to build (optional)
- `main.go` <- your shell
- `builtins package` <- each command should have it's own file (for readability)
