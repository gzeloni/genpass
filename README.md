# genpass
An over engineered password generator to generated passwords that looks "cool". 

## Usage Method
### Clone the repository

```bash
git clone https://github.com/gzeloni/genpass.git
cd genpass
```

### Build

```bash
go build -ldflags="-s -w" -o genpass genpass.go
```

### Run (Local)

```bash
./genpass
```

### Install globally (Optional)

> [!NOTE]
> You can install the binary as a global command:

```bash
mkdir -p ~/.local/bin
mv genpass ~/.local/bin
echo 'export PATH="$HOME/.local/bin:$PATH"' >> .bashrc
source .bashrc
``` 

### Run

```bash
genpass
```

### Examples
> [!IMPORTANT]
> Currently, the project has the following flags available: -d **(delimiter)**, -k **(length)**, --no-symbols.

#### Default (length = 36)

```bash
genpass
```

#### Custom length (must be perfect square)

```bash
genpass 49
```

#### Custom delimiter

```bash
genpass 36 ":"
```

#### Using flags

```bash
genpass -k 64 -d ":" 
```

#### Without symbols

```bash
genpass --no-symbols
```

## Example Output

```bash
G7aKz9-Lm3QxP2-Bv8RtY1-Jk0WsU6
```

## License

This project is licensed under the Apache License.
