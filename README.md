# iospeedtest

This application makes bits fly from one location to another and measures how fast they go.

## Usage

`iospeedtest -srcdir PATH -dstdir PATH [-size 1] [-streams 1] [-cleanup false]`

#### srcdir

Directory where source files will be written to. These files are created with size in GB equal to the `size` parameter, or `1` by default. Files are generated with random sequence of alphanumeric characters.

#### dstdir

Directory where the source files will be copied to in order to measure speed. This application was written with a mounted file share in mind for this.

#### size

Size of each test file in GB. Default 1.

#### streams

How many transfers to perform simultaneously. This will create multiple source files. I could probably just create one source file and read it n times but whatever.

#### cleanup

Delete the source and dest files after program completes. Default false.
