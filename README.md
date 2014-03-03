Cray Chapel Mesos scheduler
========

Mesos framework scheduler for running The Chapel Parallel Programming Language (http://chapel.cray.com).
Usage:

    $ ./bin/chapel -locales <number of locales i.e. nodes> -master <leading-master> <chapel program (without -nl)>


This is still an experimental framework and any participation and feedback is appreciated.

## Installation on an Elastic Mesos Cluster

First off, go ahead and launch a cluster at http://elastic.mesosphere.io.
Then log into one of the master nodes.

    $ sudo aptitude install make g++ libprotobuf-dev mercurial golang mpich2

### Install Chapel scheduler

    $ wget https://github.com/nqn/mesos-chapel/archive/master.zip
    $ unzip master
    $ cd mesos-chapel-master
    $ export GOPATH=`pwd`
    $ make

### Install Cray Chapel

    $ wget http://gasnet.lbl.gov/GASNet-1.22.0.tar.gz
    $ tar -xvzf GASNet-1.22.0.tar.gz
    $ cd GASNet-1.22.0/
    $ ./configure
    $ sudo make install
    $ cd ../
    $ wget http://downloads.sourceforge.net/project/chapel/chapel/1.8.0/chapel-1.8.0.tar.gz
    $ tar -xvzf chapel-1.8.0.tar.gz
    $ cd chapel-1.8.0
    $ export CHPL_COMM=gasnet
    $ make
    $ echo 'export PATH=$PATH:~/chapel-1.8.0/bin/linux64/' >> ~/.bashrc
    $ source ~/.bashrc

### Compile sample chapel program

    $ cd examples/programs
    $ chpl jacobi.chpl -o jacobi
    $ cp jacobi ~/mesos-chapel-master/
    $ cp jacobi_real ~/mesos-chapel-master/
    $ cd
    
### Upload assets to HDFS

    $ cd mesos-chapel-master
    $ hadoop fs -mkdir hdfs://54.211.128.164/chapel/
    $ hadoop fs -put chapel-bootstrap.tgz hdfs://54.211.128.164/chapel/chapel-bootstrap.tgz
    $ hadoop fs -put jacobi_real hdfs://54.211.128.164/chapel/jacobi_real

### Run sample program

    $ ./bin/chapel -master ec2-54-81-226-236.compute-1.amazonaws.com:5050 -name-node hdfs://54.211.128.164 -locales 3 ./jacobi
    I0303 01:03:45.841869  4199 sched.cpp:218] No credentials provided. Attempting to register without authentication
    I0303 01:03:45.842119  4199 sched.cpp:230] Detecting new master
    [ 1 / 3 ] Setting up locale..
    [ 2 / 3 ] Setting up locale..
    [ 3 / 3 ] Setting up locale..
    Jacobi computation complete.
    Delta is 9.92124e-06 (< epsilon = 1e-05)
    # of iterations: 60
    
    
