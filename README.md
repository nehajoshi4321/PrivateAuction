# PrivateAuction

## 1.Install PBC

```sh
$ sudo apt-get install build-essential flex bison

The PBC source can be compiled and installed using the usual GNU Build System:

$ wget -c https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz -O - | tar -xz
$ cd pbc-0.5.14
``````````
cd pbc-0.5.14
./configure
make
sudo make install
``````````````````

After installing, you may need to rebuild the search path for libraries.

**NOTE: the PBC library is installed to /usr/local/lib so you may need to add ```export LD_LIBRARY_PATH=/usr/local/lib/``` to your .profile or equivalent**

**## 2.Install BGN**
Most systems include a package for GMP. To install GMP in Debian / Ubuntu:

```sh
$ sudo apt-get install libgmp-dev
```
For an RPM installation with YUM:
```sh
$ sudo yum install gmp-devel
```
For installation with Fink (http://www.finkproject.org/) on Mac OS X:
```sh
$ sudo fink install gmp gmp-shlibs
```
For more information or to compile from source, visit https://gmplib.org/

## 3.Running BGN
```sh
$ cd bgn
$ make install && make build && make run

## Testing BGN

$ cd bgn
$ make install && make build
$ go test
$ go test -bench=.

## Testing the auction code
``````````````````
$ cd bgn/cmd/
$ go run test2.go "\nIteration:", i, "keyBitLength:", keyBitLength, "\t msgSpace:", msgSpace, "\t numBidders:", numBidders, randPercent:", randPercent, "\t maxBid:", maxBid //to test the encrypted message
``````````````````
