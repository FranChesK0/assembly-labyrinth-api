# [Assembly Labyrinth](https://assembly-labyrinth.franchesko.top)
## _Your Path to Assembly Mastery_

Assembly Labyrinth is a web-app, where you can test your knowledge of basic assembly.
- Learn about basic assembly
- Test your knowledge
- Add your own levels

## Manual

### 1. Architecture

The Execution Node executes a program specified in the Instruction Set.
A Node program specifies computational operations to perform. Operations
are performed sequentially, beginning with the first instruction program.
After executing the last instruction of the program, exectution automatically
continues to the first instruction.

All registers store integer values between -999 and 999 (inclusive).

#### 1.1. ACC
_Type: Internal_

ACC is the primary storage register for a Node. ACC is used as the implicit
source or destiation operand of many instructions.

#### 1.2. BAK
_Type: Internal (non-addressable)_

BAK is temporary storage for values in ACC. It is only accessible through the
SAV and SWP instructions, and cannot be read or written directly.

#### 1.3. LEFT, RIGHT, UP, DOWN
_Type: Port_

The four communication registers UP, DOWN, LEFT and RIGHT correspond to the four
ports that all Nodes use to communicate with topologically adjacent nodes. Some
ports will be disconnected on certain Nodes within the hardware, and will block
indefinitely if a READ or WRITE command is issued.

### 2. Instruction Set

<SRC> and <DST> instruction parameters may specify a port or internal register.
Any use of a port will block until the corresponding node connected to that
port completes the communication by reading or writing a value. Additionally,
a <SRC> parameter may be a literal integer value between -999 and 999(inclusive).

<LABEL> parameters are arbitrary textual names used to specify jump targets.

#### 2.1. Labels
_Syntax: <LABEL>:_

Labels are used to identify targets for jump instructions. When used as jump
target, the instruction following the label will be executed next.

Examples:
```
LOOP:       This label is on a line by itself.
    MOV 8 ACC
```
#### 2.2. MOV
_Syntax: MOV <SRC> <DST>_

<SRC> is read and the resulting value is written to <DST>.

Examples:
```
MOV 8 ACC           The literal value 8 is written to the ACC register.
MOV LEFT RIGHT      A value is read from the LEFT port, and then written to RIGHT.
```

#### 2.3. SWP
_Syntax: SWP_

The values of ACC and BAK are exchanged.

#### 2.4. SAV
_Syntax: SAV_

The value of ACC is written to BAK.

#### 2.5. ADD
_Syntax: ADD <SRC>_

The value of <SRC> is added to the value of ACC and the result is stored to ACC.

Examples:
```
ADD 16      The literal value 16 is added to the value in the ACC register.
ADD LEFT    A value is read from the LEFT port, and then added to ACC.
```

#### 2.6. SUB
_Syntax: SUB <SRC>_

The value of <SRC> is subtracted from the value of ACC and the result is stored
to ACC.

Examples:
```
SUB 16      The literal value of 16 is subtracted from the value in the ACC.
SUB LEFT    A value is read from the LEFT port, and then subtracted from ACC.
```

#### 2.7. NEG
_Syntax: NEG_

The value of ACC is arithmetically negated. A value of zero remains the same.

#### 2.8. JMP
_Syntax: JMP <LABEL>_

The instruction after the label <LABEL> will be executed next.

#### 2.9. JEZ
_Syntax: JEZ <LABEL>_

The instruction after the label <LABEL> will be executed next if the value
of ACC is zero.

#### 2.10. JNZ
_Syntax: JNZ <LABEL>_

The instruction after the label <LABEL> will be executed next if the value
of ACC is not zero.

#### 2.11. JGZ
_Syntax: JGZ <LABEL>_

The instruction after the label <LABEL> will be executed next if the value
of ACC is greater than zero.

#### 2.12. JLZ
_Syntax: JLZ <LABEL>_

The instruction after the label <LABEL> will be executed next if the value
of ACC is less than zero.

#### 2.13. JRO
_Syntax: JRO <SRC>_

The instruction at the offset specified by <SRC> relative to the current
instruction will be executed next.

### 3. Example Programs
The following sample program reads a sequence of values from the
LEFT port, doubling each value read and writing that to the RIGHT
port. Because of the automatic looping behavior, it continues to
the first instruction after executing the last instruction.

```
MOV LEFT ACC        Read a value from the LEFT port into the ACC register.
ADD ACC             Add the value in ACC to itself, doubling it.
MOV ACC RIGHT       Write the value in the ACC register to the RIGHT port.
```

The following sample program reads a sequence of values from the
UP port, writing positive values to the RIGHT port and negative
values to the LEFT port. Zero values are discarded.

```
START:
    MOV UP ACC      Read a value from the UP port into the ACC register.
    JGZ POSITIVE    If the value in ACC is greater than zero, jump to "POSITIVE".
    JLZ NEGATIVE    If the value in ACC is less than zero, jump to "NEGATIVE".
    JMP START       The value was neither positive nor negative, so jump to "START".
POSITIVE:
    MOV ACC RIGHT   Write the value in the ACC register to the RIGHT port.
    JMP START       Jump to "START".
NEGATIVE:
    MOV ACC LEFT    Write the value in the ACC register to the LEFT port.
    JMP START       Jump to "START".
```

## Docker

Assembly Labyrinth is very easy to install and deploy in a Docker container.

Clone repositories: [API](https://github.com/FranChesK0/assembly-labyrinth) and [WEB](https://github.com/FranChesK0/assembly-labyrinth-web)

```sh
git clone https://github.com/FranChesK0/assembly-labyrinth
git clone https://github.com/FranChesK0/assembly-labyrinth-web
```

In web add api.js file.

```sh
cd assembly-labyrinth-web
touch api.js
```

This file should like this:

```js
const api = "https://your-domain-with-api/levels"
```

When ready, simply use the Dockerfile to build the image.

```sh
cd assembly-labyrinth
docker-compose build
```
```sh
cd assembly-labyrinth-web
docker-compose build
```

This will create the images and pull in the necessary dependencies.

By default, the Docker will expose ports 80 and 8082, so change this within the
compose.yaml if necessary. Once done, run the Docker image:

```sh
docker-compose up -d
```

Verify the deployment by navigating to your server address in
your preferred browser.

## Authors

[FranChesK0](https://github.com/FranChesK0) - [API](https://github.com/FranChesK0/assembly-labyrinth)

[MaFFioZZo](https://github.com/MaFFioZZo) - [WEB](https://github.com/MaFFioZZo/assembly-labyrinth)

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

**Free Software, Hell Yeah!**