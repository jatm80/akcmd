# Akash Command Center

A lite weight `client` for communicating with the Akash Network. There are no `provider`, or `validator` bits in this project.

# Install

### Linux or Mac

Install using `curl` or `wget` on POSIX systems

```
curl -o- https://raw.githubusercontent.com/ovrclk/akcmd/spike/install.sh | bash && /tmp/akcmd_installer
```

```
wget -qO- https://raw.githubusercontent.com/ovrclk/akcmd/spike/install.sh | bash && /tmp/akcmd_installer
```

### Windows

Install using Powershell

Ensure you user has a script policy set

```
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Default
```

Run the install script

```
cd $Env:USERPROFILE;
Invoke-WebRequest https://raw.githubusercontent.com/ovrclk/akcmd/spike/install.ps1 -OutFile install.ps1;
.\install.ps1;
del install.ps1
```

# Contributing

### Building

Produces binary in the `out` folder

- golang 1.16.x

Generate the client app

```
make clean && make client
```

Generate the installer app

```
make clean && make installer
```
