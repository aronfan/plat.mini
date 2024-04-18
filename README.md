# plat.mini
golang backend module for miniprogram and minigame

## Installation
* go get
```bash
go get github.com/aronfan/plat.mini
```
* git submodule (recommended)
```bash
cd /path/to/proj
mkdir sub
cd sub
git submodule add https://github.com/aronfan/plat.mini
git submodule update --remote
git add plat.mini
git commit -m "add plat.mini"
git push
```
add the following line to your go.mod if you want to make local modifications
```bash
replace github.com/aronfan/plat.mini => ./sub/plat.mini
```