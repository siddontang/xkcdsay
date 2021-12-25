xkcdsay is a simple application just for fun. Once again, Just for fun.

I like `cowsay`, `ponysay` or other similar programs very much, and I also like accessing [xkcd](https://xkcd.com/) to view the comics too. So I build the `xkcdsay` which can let me view the xkcd comic in the terminal. 

All commic data of `xkcdsay` is downloaded from XKCD and is saved to a [TiDB](https://github.com/pingcap/tidb) DevTier cluster hosted on [TiDB cloud](https://tidbcloud.com/). 

[**You can singup and try TiDB cloud for Free**](https://tidbcloud.com/signup)

## Build from the source

```bash
git clone https://github.com/siddontang/xkcdsay.git
cd xkcdsay
make

# the xkcdsay binary will be installed in the current ./bin/xkcdsay
```

## Install with Homebrew

```bash
brew install siddontang/brew/xkcdsay
```

## Usage

```bash
# Randomly see a comic 
xkcdsay 
```

![image](https://user-images.githubusercontent.com/1080370/147331905-4247319c-340d-45bf-938e-fd04eef779b5.png)

```bash
# See the 1st comic
xkcdsay -n 1
```


![image](https://user-images.githubusercontent.com/1080370/147331792-1f6b2769-ddf1-4e11-9afa-d7e623f7b32d.png)

## TODO - Need help

- [ ] support comic cache for xkcdsay
- [ ] add an ASCII art support like cowsay for the saying words
