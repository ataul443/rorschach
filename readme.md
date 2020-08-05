# Rorschach

## Overview

It is a file monitoring application, which scans a root directory for any 
changes to files underneath this folder and executes user-specified actions 
on the files based on pattern matching rules.

For instance, suppose rorschach was monitoring the root folder 
`/folder/to/monitor` and a file named `/folder/to/monitor/Foster the People/Helena Beat.mp3` 
was created. A scan of the root folder would detect that this file was created 
and check if there is a pattern rule that matches the name of the file.

For example, suppose there is a pattern rule that looks like this: `CREATE *.mp3 mpv ${FULLPATH}.` 
When the scan detects that the `/folder/to/monitor/Foster the People/Helena Beat.mp3` is created, 
it will see that this file matches the `*.mp3 pattern` and will execute the `mpv ${FULLPATH}` 
command (ie. it will play the song).

### Why we are building it ?

We wanted to learn more about concurrency and optimisations like cpu pooling 
and memory pooling in go. Hence we decided to build this application. We will 
iterate over it's design constantly to make it efficient and fast while learning
a great deal about profiling and testing stuff.

## Usage

Download the git repository somewhere in your local host
```
git clone https://github.com/ataul443/rorschach.git
cd rorschach
git checkout alpha
```

For information on configuration options please look in the `config.toml` file.
Once you are ready with your configuration. Run the following command to build the
rorschach and run it.
```
go build .
./rorschach
```

## Note
The project is in very early stage and heavily work in progress. Right now It is full of bugs,
however we are working hard to improve it day by day.