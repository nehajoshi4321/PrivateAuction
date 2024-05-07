# PrivateAuction

## 1.Install PBC

```sh
$ sudo apt-get install -y build-essential flex bison libgmp3-dev

The PBC source can be compiled and installed using the usual GNU Build System:

$ wget -c https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz -O - | tar -xz
$ cd pbc-0.5.14
$ ./configure
$ make
$ sudo make install
$ export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/usr/local/lib/
$ sudo vi /etc/profile ; # add above path

``````````````````

After installing, you may need to rebuild the search path for libraries.

**NOTE: the PBC library is installed to /usr/local/lib so you may need to add ```export LD_LIBRARY_PATH=/usr/local/lib/``` to your .profile or equivalent**

**## 2.Install BGN re-req on Ubuntu**
We may need Go to be installed - ref https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-20-04
```sh
$  curl  -OL https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
$  sha256sum go1.22.2.linux-amd64.tar.gz 
$  sudo tar -C /usr/local -xvf go1.22.2.linux-amd64.tar.gz 
$  export PATH=$PATH:/usr/local/go/bin
$  sudo vi /etc/profile ; # add above paths
```
alternate
```sh
sudo snap install go
```

## 3.Running BGN
```sh
$ cd bgn
$ make install && make build && make run
``````````````````
## Testing BGN
``````````````````
$ cd bgn
$ make install && make build ; # has errors on duplicate globals and clash in name space
$ go test
$ go test -bench=.
``````````````````
## Testing the auction code
``````````````````
$ cd bgn/cmd/
$ go run test2.go <Iteration:> <keyBitLength> <msgSpace> <msgSpace> <numBidders> <numBidders> <randPercent> <maxBid> ;  //to test the encrypted message
$ go run test2.go  5 1000  100000 4 5 100
``````````````````
