# snake.go
### A slither.io inspired game for the terminal

This is game I started writing to learn Go for a new job, as well as learn some 
Github Actions for CI/CD and other things. It was inspired by slither.io and in
hindsight should have probably been named slither.go...

By default the game is in client mode, where you will be prompted to connect to a
server (which I have kept private for obvious reasons) but you can play local games 
with the flag --mode host, where other players will connect to your game game, 
or setup a dedicated server (no graphics) with --mode server.

In practice, never copy/paste commands from github into the terminal😭 
But here are some commands generated by my `make release` command that 
hopefully make installing the lastest version of the game an easier process.

After installing, run in the terminal with `snake-go`

## To install
### Mac
```
curl -L -o /tmp/snake-go.tar.gz "https://github.com/dbrun3/snake-go/releases/download/v1.0.3/snake-go-v1.0.3-darwin-arm64.tar.gz" && sudo xattr -rd com.apple.quarantine /tmp/snake-go.tar.gz && tar -xzf /tmp/snake-go.tar.gz -C /tmp/ && sudo mv /tmp/snake-go /usr/local/bin/ && sudo chmod +x /usr/local/bin/snake-go
```
### Linux
```
curl -L -o /tmp/snake-go.tar.gz "https://github.com/dbrun3/snake-go/releases/download/v1.0.3/snake-go-v1.0.3-linux-amd64.tar.gz" && tar -xzf /tmp/snake-go.tar.gz -C /tmp/ && sudo mv /tmp/snake-go /usr/local/bin/ && sudo chmod +x /usr/local/bin/snake-go
```
### Windows
Using PowerShell:
```
Invoke-WebRequest -Uri "https://github.com/dbrun3/snake-go/releases/download/v1.0.3/snake-go-v1.0.3-windows-amd64.zip" -OutFile "$env:TEMP\snake-go.zip"; Expand-Archive -Path "$env:TEMP\snake-go.zip" -DestinationPath "$env:TEMP\snake-go" -Force; $WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\Snake-Go.lnk"); $Shortcut.TargetPath = "$env:TEMP\snake-go\snake-go.exe"; $Shortcut.Save()
```
