xkcdsay is a simple application just for fun. Once again, Just for fun.

I like `cowsay`, `ponysay` or other similar programs very much, and I also like accessing [xkcd](https://xkcd.com/) to view the comics too. So I build the `xkcdsay` which can let me view the xkcd comic in the terminal. 

All commic data of `xkcdsay` is downloaded from XKCD and is saved to a [TiDB](https://github.com/pingcap/tidb) DevTier cluster hosted on [TiDB cloud](https://tidbcloud.com/). 

[**You can singup and try TiDB cloud for Free**](https://tidbcloud.com/signup)

## Build

```bash
make
```

## Usage

```bash
# Randomly see a comic 
./bin/xkcdsay 

# See the newest one currently
./bin/xkcdsay -n -1

# See the 10th comic
./bin/xkcdsay -n 10
```