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
    $ cp jacobi ~/
    $ cp jacobi_real ~/
    $ cd

### Install Chapel scheduler

    $ wget https://github.com/nqn/mesos-chapel/archive/master.zip
    $ unzip master
    $ cd mesos-chapel-master
    $ export GOPATH=`pwd`
    $ make
    
### Upload assets to HDFS

    $ hadoop fs -mkdir hdfs://54.211.128.164/chapel/
    $ hadoop fs -put mesos-chapel-master/chapel-bootstrap.tgz hdfs://54.211.128.164/chapel/chapel-bootstrap.tgz

### Run sample program

    $ ./mesos-chapel-master/bin/chapel -master ec2-54-81-226-236.compute-1.amazonaws.com:5050 -bootstrap hdfs://54.211.128.164/chapel/chapel-bootstrap.tgz ./jacobi
    
    
